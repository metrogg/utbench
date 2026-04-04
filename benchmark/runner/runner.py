from __future__ import annotations

import concurrent.futures
import hashlib
import json
import logging
import os
import re
import threading
import time
from dataclasses import asdict, dataclass
from datetime import datetime
from pathlib import Path
from typing import Any
from urllib import error, request

import yaml

from .prompt_builder import PromptBuilder


DEFAULT_API_TIMEOUT_MS = 60_000


@dataclass
class ModelConfig:
    model_name: str
    provider: str
    endpoint: str
    model: str
    api_key_env: str
    parameters: dict[str, Any]
    timeout_ms: int


@dataclass
class ErrorInfo:
    kind: str
    message: str
    retryable: bool
    status_code: int | None = None


@dataclass
class TestResult:
    success: bool
    model: str
    language: str
    sample_id: str
    sample_path: str
    generated_test_path: str | None
    metadata_path: str | None
    latency_ms: int | None
    prompt_tokens: int | None
    completion_tokens: int | None
    total_tokens: int | None
    error: dict[str, Any] | None = None


class Runner:
    """Read dataset samples, call model APIs, and store generated tests."""

    LANGUAGE_EXT = {
        "java": "java",
        "python": "py",
        "go": "go",
        "cpp": "cpp",
        "javascript": "js",
    }

    def __init__(
        self,
        config_path: str | Path = "benchmark/config/models.yaml",
        dataset_root: str | Path = "dataset",
        results_root: str | Path = "results",
        retries: int = 3,
        backoff_seconds: float = 2.0,
        logger: logging.Logger | None = None,
    ) -> None:
        self.config_path = Path(config_path)
        self.dataset_root = Path(dataset_root)
        self.results_root = Path(results_root)
        self.retries = max(1, retries)
        self.backoff_seconds = max(0.0, backoff_seconds)
        self.logger = logger or logging.getLogger("utbench.runner")
        self._config = self._load_config()
        self._api_timeout_ms = (
            self._config.get("benchmark", {})
            .get("timeouts", {})
            .get("api_call", DEFAULT_API_TIMEOUT_MS)
        )
        parallel_cfg = self._config.get("benchmark", {}).get("parallel", {})
        self._max_concurrent_models = int(parallel_cfg.get("max_concurrent_models", 1))
        self._max_concurrent_samples = int(parallel_cfg.get("max_concurrent_samples", 1))
        self._max_workers = max(
            1,
            self._max_concurrent_models * self._max_concurrent_samples,
        )
        self._prompt_builder = PromptBuilder(self._config.get("benchmark", {}))

    def run(
        self,
        model_names: list[str] | None = None,
        languages: list[str] | None = None,
        max_samples: int | None = None,
        dry_run: bool = False,
        resume: bool = True,
        reset_checkpoint: bool = False,
        max_workers: int | None = None,
    ) -> dict[str, Any]:
        models = self._select_models(model_names)
        samples = self._collect_samples(languages, max_samples)
        run_id = datetime.utcnow().strftime("%Y%m%dT%H%M%S%fZ")
        selected_models = [m.model_name for m in models]

        checkpoint_file = self._checkpoint_path(selected_models, languages, max_samples)
        if reset_checkpoint and checkpoint_file.exists():
            checkpoint_file.unlink()
        checkpoint = self._load_checkpoint(checkpoint_file) if resume else {"completed": []}
        completed = set(checkpoint.get("completed", []))

        tasks: list[tuple[ModelConfig, Path, str, str]] = []
        for model in models:
            for sample_path in samples:
                language = sample_path.parent.name.lower()
                sample_id = sample_path.stem
                key = self._task_key(model.model_name, language, sample_id)
                if resume and key in completed:
                    continue
                tasks.append((model, sample_path, language, key))
        tasks_total_by_model: dict[str, int] = {}
        for model, _, _, _ in tasks:
            tasks_total_by_model[model.model_name] = (
                tasks_total_by_model.get(model.model_name, 0) + 1
            )
        progress_by_model: dict[str, dict[str, int]] = {
            name: {"total": total, "done": 0, "success": 0, "failed": 0}
            for name, total in tasks_total_by_model.items()
        }
        progress_lock = threading.Lock()

        self.logger.info(
            "Runner start: models=%s, total_samples=%s, pending_tasks=%s, dry_run=%s, resume=%s",
            selected_models,
            len(samples),
            len(tasks),
            dry_run,
            resume,
        )

        results: list[TestResult] = []
        workers = max_workers if max_workers and max_workers > 0 else self._max_workers
        workers = max(1, workers)

        with concurrent.futures.ThreadPoolExecutor(max_workers=workers) as executor:
            future_map: dict[
                concurrent.futures.Future[TestResult],
                tuple[ModelConfig, str],
            ] = {}
            for model, sample_path, language, _ in tasks:
                fut = executor.submit(
                    self.generate_test,
                    sample_path=sample_path,
                    model_config=model,
                    language=language,
                    dry_run=dry_run,
                )
                future_map[fut] = (model, language)

            for fut in concurrent.futures.as_completed(future_map):
                model, language = future_map[fut]
                result = fut.result()
                results.append(result)
                key = self._task_key(result.model, result.language, result.sample_id)
                if result.success and resume:
                    completed.add(key)
                    self._save_checkpoint(checkpoint_file, sorted(completed))

                if result.success:
                    self.logger.info(
                        "OK model=%s lang=%s sample=%s",
                        model.model_name,
                        language,
                        result.sample_id,
                    )
                else:
                    self.logger.error(
                        "FAILED model=%s lang=%s sample=%s err=%s",
                        model.model_name,
                        language,
                        result.sample_id,
                        result.error,
                    )
                with progress_lock:
                    state = progress_by_model.get(
                        result.model,
                        {"total": 0, "done": 0, "success": 0, "failed": 0},
                    )
                    state["done"] += 1
                    if result.success:
                        state["success"] += 1
                    else:
                        state["failed"] += 1
                    progress_by_model[result.model] = state
                    self.logger.info(
                        "PROGRESS model=%s done=%s/%s success=%s failed=%s",
                        result.model,
                        state["done"],
                        state["total"],
                        state["success"],
                        state["failed"],
                    )

        summary = self._build_summary(run_id, results, dry_run)
        summary["resume"] = resume
        summary["checkpoint_file"] = str(checkpoint_file.resolve())
        summary["pending_tasks"] = len(tasks)
        summary["skipped_by_checkpoint"] = (len(selected_models) * len(samples)) - len(tasks)
        summary["max_workers"] = workers
        summary["progress_by_model"] = progress_by_model
        summary_file = self.results_root / f"runner_summary_{run_id}.json"
        summary_file.parent.mkdir(parents=True, exist_ok=True)
        summary_file.write_text(
            json.dumps(summary, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )
        self.logger.info("Summary written to %s", summary_file)
        return summary

    def generate_test(
        self,
        sample_path: str | Path,
        model_config: ModelConfig,
        language: str,
        dry_run: bool = False,
    ) -> TestResult:
        """Generate tests for one sample and persist artifacts."""
        sample_input = self._prompt_builder.read_sample(
            sample_path=sample_path,
            default_language=language,
        )
        sample = sample_input.sample_path
        sample_id = sample_input.sample_id
        language = sample_input.language
        prompt = self._prompt_builder.build_prompt(sample_input)

        latency_ms: int | None = None
        prompt_tokens: int | None = None
        completion_tokens: int | None = None
        total_tokens: int | None = None
        raw_response: dict[str, Any] | None = None

        try:
            if dry_run:
                test_code = self._make_dry_run_content(sample_id, language)
            else:
                raw_response, latency_ms = self.call_model_api(prompt, model_config)
                usage = self._extract_usage(raw_response)
                prompt_tokens = usage.get("prompt_tokens")
                completion_tokens = usage.get("completion_tokens")
                total_tokens = usage.get("total_tokens")
                content = self._extract_response_text(raw_response, model_config.provider)
                test_code = self._extract_code(content, language)

            test_file, metadata_file, response_file = self._write_outputs(
                model=model_config.model_name,
                language=language,
                sample_id=sample_id,
                sample_path=sample,
                prompt=prompt,
                generated_test=test_code,
                raw_response=raw_response,
                latency_ms=latency_ms,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=total_tokens,
            )

            return TestResult(
                success=True,
                model=model_config.model_name,
                language=language,
                sample_id=sample_id,
                sample_path=str(sample.resolve()),
                generated_test_path=str(test_file.resolve()),
                metadata_path=str(metadata_file.resolve()),
                latency_ms=latency_ms,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=total_tokens,
                error=None,
            )
        except Exception as exc:  # pylint: disable=broad-except
            err = self.handle_error(exc)
            self._write_failure_metadata(
                model=model_config.model_name,
                language=language,
                sample_id=sample_id,
                sample_path=sample,
                error_info=err,
            )
            return TestResult(
                success=False,
                model=model_config.model_name,
                language=language,
                sample_id=sample_id,
                sample_path=str(sample.resolve()),
                generated_test_path=None,
                metadata_path=None,
                latency_ms=latency_ms,
                prompt_tokens=prompt_tokens,
                completion_tokens=completion_tokens,
                total_tokens=total_tokens,
                error=asdict(err),
            )

    def call_model_api(
        self, prompt: str, model_config: ModelConfig
    ) -> tuple[dict[str, Any], int]:
        """Call model API and return (response_json, latency_ms)."""
        api_key = os.getenv(model_config.api_key_env)
        if not api_key:
            raise RuntimeError(
                f"Missing API key env var: {model_config.api_key_env} "
                f"(model={model_config.model_name})"
            )

        endpoint = self._resolve_endpoint(model_config)
        payload = self._build_payload(prompt, model_config)
        payload_bytes = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        timeout = max(1.0, model_config.timeout_ms / 1000.0)

        last_exc: Exception | None = None
        for attempt in range(1, self.retries + 1):
            started = time.perf_counter()
            req = request.Request(
                endpoint,
                data=payload_bytes,
                headers={
                    "Content-Type": "application/json",
                    "Authorization": f"Bearer {api_key}",
                },
                method="POST",
            )

            try:
                with request.urlopen(req, timeout=timeout) as resp:
                    raw = resp.read().decode("utf-8")
                    latency_ms = int((time.perf_counter() - started) * 1000)
                    parsed = json.loads(raw)
                    return parsed, latency_ms
            except error.HTTPError as exc:
                body = exc.read().decode("utf-8", errors="replace")
                last_exc = RuntimeError(
                    f"HTTP {exc.code} calling {endpoint}: {body[:500]}"
                )
                err_info = self.handle_error(last_exc, status_code=exc.code)
                if not err_info.retryable or attempt >= self.retries:
                    raise last_exc
            except error.URLError as exc:
                last_exc = RuntimeError(f"Network error calling {endpoint}: {exc.reason}")
                err_info = self.handle_error(last_exc)
                if not err_info.retryable or attempt >= self.retries:
                    raise last_exc
            except TimeoutError as exc:
                last_exc = RuntimeError(f"Timeout calling {endpoint}: {exc}")
                err_info = self.handle_error(last_exc)
                if not err_info.retryable or attempt >= self.retries:
                    raise last_exc

            sleep_seconds = self.backoff_seconds * (2 ** (attempt - 1))
            if sleep_seconds > 0:
                time.sleep(sleep_seconds)

        if last_exc is None:
            raise RuntimeError("Unknown API error")
        raise last_exc

    def handle_error(self, exc: Exception, status_code: int | None = None) -> ErrorInfo:
        """Classify errors and decide retryability."""
        msg = str(exc)
        retryable = False
        kind = "unknown_error"

        if status_code is None:
            http_match = re.search(r"\bHTTP\s+(\d{3})\b", msg)
            if http_match:
                status_code = int(http_match.group(1))

        if status_code is not None:
            if status_code in (408, 429) or status_code >= 500:
                kind = "http_retryable"
                retryable = True
            else:
                kind = "http_non_retryable"
                retryable = False
            return ErrorInfo(
                kind=kind,
                message=msg,
                retryable=retryable,
                status_code=status_code,
            )

        lowered = msg.lower()
        if "timeout" in lowered:
            kind = "timeout"
            retryable = True
        elif "network error" in lowered or "temporarily unavailable" in lowered:
            kind = "network_error"
            retryable = True
        elif "missing api key env var" in lowered:
            kind = "auth_config_error"
            retryable = False
        elif "json" in lowered:
            kind = "response_parse_error"
            retryable = False

        return ErrorInfo(kind=kind, message=msg, retryable=retryable)

    def _load_config(self) -> dict[str, Any]:
        if not self.config_path.exists():
            raise FileNotFoundError(f"Config file not found: {self.config_path}")
        with self.config_path.open("r", encoding="utf-8") as f:
            config = yaml.safe_load(f)
        if not isinstance(config, dict):
            raise ValueError(f"Invalid config content: {self.config_path}")
        return config

    def _select_models(self, model_names: list[str] | None) -> list[ModelConfig]:
        model_map = self._config.get("models", {})
        if not isinstance(model_map, dict):
            raise ValueError("Invalid models section in config")

        names = [n.strip() for n in model_names or [] if n.strip()]
        selected: list[ModelConfig] = []
        for name, item in model_map.items():
            if names and name not in names:
                continue
            if not item.get("enabled", False):
                continue
            cfg = item.get("config", {})
            selected.append(
                ModelConfig(
                    model_name=name,
                    provider=str(item.get("provider", "")).lower(),
                    endpoint=str(cfg.get("api_endpoint", "")).rstrip("/"),
                    model=str(cfg.get("model", "")),
                    api_key_env=str(cfg.get("api_key_env", "")),
                    parameters=dict(cfg.get("parameters", {})),
                    timeout_ms=int(self._api_timeout_ms),
                )
            )

        if not selected:
            raise ValueError(f"No enabled models selected. input={model_names}")
        return selected

    def _collect_samples(
        self, languages: list[str] | None, max_samples: int | None
    ) -> list[Path]:
        if not self.dataset_root.exists():
            raise FileNotFoundError(f"Dataset dir not found: {self.dataset_root}")

        language_set = {
            lang.strip().lower() for lang in (languages or []) if lang.strip()
        }
        samples: list[Path] = []
        for lang_dir in sorted(self.dataset_root.iterdir()):
            if not lang_dir.is_dir():
                continue
            language = lang_dir.name.lower()
            if language_set and language not in language_set:
                continue

            files = sorted([p for p in lang_dir.rglob("*") if p.is_file()])
            if max_samples is not None and max_samples > 0:
                files = files[:max_samples]
            samples.extend(files)

        if not samples:
            raise ValueError(f"No dataset samples found. languages={languages}")
        return samples

    def _resolve_endpoint(self, model_config: ModelConfig) -> str:
        base = model_config.endpoint.rstrip("/")
        if model_config.provider == "dashscope":
            if "compatible-mode" in base:
                return f"{base}/chat/completions"
            if base.endswith("/api/v1"):
                return f"{base}/services/aigc/text-generation/generation"
            return f"{base}/services/aigc/text-generation/generation"
        return f"{base}/chat/completions"

    def _build_payload(self, prompt: str, model_config: ModelConfig) -> dict[str, Any]:
        params = dict(model_config.parameters)
        if model_config.provider == "dashscope":
            if "compatible-mode" in model_config.endpoint:
                payload = {
                    "model": model_config.model,
                    "messages": [
                        {
                            "role": "system",
                            "content": "You generate high-quality unit tests.",
                        },
                        {"role": "user", "content": prompt},
                    ],
                    "stream": False,
                }
                payload.update(params)
                return payload
            return {
                "model": model_config.model,
                "input": {"messages": [{"role": "user", "content": prompt}]},
                "parameters": params,
            }

        payload = {
            "model": model_config.model,
            "messages": [
                {
                    "role": "system",
                    "content": "You generate high-quality unit tests.",
                },
                {"role": "user", "content": prompt},
            ],
            "stream": False,
        }
        payload.update(params)
        return payload

    def _extract_response_text(self, response: dict[str, Any], provider: str) -> str:
        if provider == "dashscope":
            output = response.get("output", {})
            if isinstance(output, dict):
                if isinstance(output.get("text"), str):
                    return output["text"]
                choices = output.get("choices", [])
                if choices:
                    message = choices[0].get("message", {})
                    content = message.get("content")
                    if isinstance(content, str):
                        return content
                    if isinstance(content, list):
                        text_parts = [
                            x.get("text", "") for x in content if isinstance(x, dict)
                        ]
                        if text_parts:
                            return "\n".join(text_parts)

        choices = response.get("choices", [])
        if not choices:
            raise ValueError(f"No choices in model response: {response}")

        first = choices[0]
        message = first.get("message", {})
        content = message.get("content")
        if isinstance(content, str):
            return content
        if isinstance(content, list):
            text_parts = [x.get("text", "") for x in content if isinstance(x, dict)]
            if text_parts:
                return "\n".join(text_parts)
        if isinstance(first.get("text"), str):
            return first["text"]
        raise ValueError(f"Unable to extract text from model response: {response}")

    def _extract_usage(self, response: dict[str, Any]) -> dict[str, int | None]:
        usage = response.get("usage", {})
        if not isinstance(usage, dict):
            usage = {}
        return {
            "prompt_tokens": usage.get("prompt_tokens") or usage.get("input_tokens"),
            "completion_tokens": usage.get("completion_tokens")
            or usage.get("output_tokens"),
            "total_tokens": usage.get("total_tokens"),
        }

    def _extract_code(self, content: str, language: str) -> str:
        blocks = re.findall(r"```([a-zA-Z0-9_]*)\n(.*?)```", content, flags=re.DOTALL)
        if not blocks:
            return content.strip()

        lang = language.lower()
        aliases = {lang}
        if lang == "python":
            aliases.add("py")
        if lang == "javascript":
            aliases.add("js")
        if lang == "cpp":
            aliases.add("c++")

        for tag, code in blocks:
            if tag.strip().lower() in aliases:
                return code.strip()
        return blocks[0][1].strip()

    def _write_outputs(
        self,
        model: str,
        language: str,
        sample_id: str,
        sample_path: Path,
        prompt: str,
        generated_test: str,
        raw_response: dict[str, Any] | None,
        latency_ms: int | None,
        prompt_tokens: int | None,
        completion_tokens: int | None,
        total_tokens: int | None,
    ) -> tuple[Path, Path, Path]:
        timestamp = datetime.utcnow().strftime("%Y%m%d%H%M%S")
        ext = self.LANGUAGE_EXT.get(language.lower(), "txt")
        base_name = f"{model}_{language}_{sample_id}_{timestamp}"

        model_root = self.results_root / model
        tests_dir = model_root / "tests"
        reports_dir = model_root / "reports"
        artifacts_dir = model_root / "artifacts"
        tests_dir.mkdir(parents=True, exist_ok=True)
        reports_dir.mkdir(parents=True, exist_ok=True)
        artifacts_dir.mkdir(parents=True, exist_ok=True)

        test_file = tests_dir / f"{base_name}.test.{ext}"
        metadata_file = reports_dir / f"{base_name}.metadata.json"
        response_file = artifacts_dir / f"{base_name}.response.json"

        test_file.write_text(generated_test, encoding="utf-8")

        metadata = {
            "model": model,
            "language": language,
            "sample_id": sample_id,
            "sample_path": str(sample_path.resolve()),
            "generated_test_path": str(test_file.resolve()),
            "created_at_utc": datetime.utcnow().isoformat() + "Z",
            "latency_ms": latency_ms,
            "tokens": {
                "prompt_tokens": prompt_tokens,
                "completion_tokens": completion_tokens,
                "total_tokens": total_tokens,
            },
            "prompt_preview": prompt[:400],
            "success": True,
        }
        metadata_file.write_text(
            json.dumps(metadata, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

        response_payload = raw_response if raw_response is not None else {"dry_run": True}
        response_file.write_text(
            json.dumps(response_payload, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )
        return test_file, metadata_file, response_file

    def _write_failure_metadata(
        self,
        model: str,
        language: str,
        sample_id: str,
        sample_path: Path,
        error_info: ErrorInfo,
    ) -> None:
        timestamp = datetime.utcnow().strftime("%Y%m%d%H%M%S")
        model_root = self.results_root / model
        failed_dir = model_root / "artifacts" / "failed"
        failed_dir.mkdir(parents=True, exist_ok=True)
        failure_file = failed_dir / f"{model}_{language}_{sample_id}_{timestamp}.error.json"
        payload = {
            "model": model,
            "language": language,
            "sample_id": sample_id,
            "sample_path": str(sample_path.resolve()),
            "created_at_utc": datetime.utcnow().isoformat() + "Z",
            "success": False,
            "error": asdict(error_info),
        }
        failure_file.write_text(
            json.dumps(payload, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )

    def _make_dry_run_content(self, sample_id: str, language: str) -> str:
        if language == "python":
            return (
                "import pytest\n\n"
                f"def test_{sample_id}_placeholder():\n"
                "    assert True\n"
            )
        if language == "java":
            return (
                "import org.junit.Test;\n"
                "import static org.junit.Assert.assertTrue;\n\n"
                "public class GeneratedTest {\n"
                "    @Test\n"
                "    public void placeholder() {\n"
                "        assertTrue(true);\n"
                "    }\n"
                "}\n"
            )
        if language == "go":
            return (
                "package main\n\n"
                "import \"testing\"\n\n"
                "func TestPlaceholder(t *testing.T) {\n"
                "}\n"
            )
        return f"// dry-run placeholder test for {sample_id} ({language})\n"

    def _build_summary(
        self, run_id: str, results: list[TestResult], dry_run: bool
    ) -> dict[str, Any]:
        total = len(results)
        success = sum(1 for x in results if x.success)
        failed = total - success
        return {
            "run_id": run_id,
            "created_at_utc": datetime.utcnow().isoformat() + "Z",
            "dry_run": dry_run,
            "total": total,
            "success": success,
            "failed": failed,
            "results": [asdict(x) for x in results],
        }

    def _task_key(self, model: str, language: str, sample_id: str) -> str:
        return f"{model}|{language}|{sample_id}"

    def _checkpoint_path(
        self,
        model_names: list[str],
        languages: list[str] | None,
        max_samples: int | None,
    ) -> Path:
        models_norm = ",".join(sorted(model_names))
        langs_norm = ",".join(sorted([x.lower() for x in (languages or [])])) or "all"
        max_samples_norm = str(max_samples) if max_samples is not None else "all"
        scope = f"models={models_norm};langs={langs_norm};max={max_samples_norm}"
        scope_hash = hashlib.sha1(scope.encode("utf-8")).hexdigest()[:12]
        ckpt_dir = self.results_root / "checkpoints"
        ckpt_dir.mkdir(parents=True, exist_ok=True)
        return ckpt_dir / f"runner_{scope_hash}.checkpoint.json"

    def _load_checkpoint(self, checkpoint_path: Path) -> dict[str, Any]:
        if not checkpoint_path.exists():
            return {"completed": []}
        try:
            data = json.loads(checkpoint_path.read_text(encoding="utf-8"))
            if not isinstance(data, dict):
                return {"completed": []}
            completed = data.get("completed", [])
            if not isinstance(completed, list):
                completed = []
            return {"completed": completed}
        except Exception:
            return {"completed": []}

    def _save_checkpoint(self, checkpoint_path: Path, completed: list[str]) -> None:
        payload = {
            "updated_at_utc": datetime.utcnow().isoformat() + "Z",
            "completed": completed,
        }
        checkpoint_path.write_text(
            json.dumps(payload, ensure_ascii=False, indent=2),
            encoding="utf-8",
        )
