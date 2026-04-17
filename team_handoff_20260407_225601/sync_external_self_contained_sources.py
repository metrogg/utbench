import json
from pathlib import Path

from huggingface_hub import hf_hub_download

from dataset_builder_utils import ensure_dir, write_json


BASE_DIR = Path(__file__).resolve().parent
OUTPUT_DIR = BASE_DIR / "data" / "external_sources" / "java_cpp_self_contained"
MANIFEST_PATH = OUTPUT_DIR / "manifest.json"


def copy_humaneval_x(language: str) -> dict:
    remote_filename = f"data/{language}/data/humaneval.jsonl"
    local_path = Path(
        hf_hub_download(
            repo_id="THUDM/humaneval-x",
            filename=remote_filename,
            repo_type="dataset",
        )
    )
    output_path = OUTPUT_DIR / f"humanevalx_{language}.jsonl"
    ensure_dir(output_path.parent)
    output_path.write_text(local_path.read_text(encoding="utf-8"), encoding="utf-8")
    line_count = sum(1 for _ in output_path.open("r", encoding="utf-8"))
    return {
        "repo_id": "THUDM/humaneval-x",
        "remote_filename": remote_filename,
        "local_path": str(output_path.relative_to(BASE_DIR)).replace("\\", "/"),
        "records": line_count,
    }


def filter_autocodebenchmark(language: str) -> dict:
    remote_filename = "autocodebench.jsonl"
    local_path = Path(
        hf_hub_download(
            repo_id="tencent/AutoCodeBenchmark",
            filename=remote_filename,
            repo_type="dataset",
        )
    )
    output_path = OUTPUT_DIR / f"autocodebenchmark_{language}.jsonl"
    records = 0
    with local_path.open("r", encoding="utf-8") as src, output_path.open("w", encoding="utf-8") as dst:
        for line in src:
            item = json.loads(line)
            if item.get("language") != language:
                continue
            dst.write(json.dumps(item, ensure_ascii=False) + "\n")
            records += 1
    return {
        "repo_id": "tencent/AutoCodeBenchmark",
        "remote_filename": remote_filename,
        "language_filter": language,
        "local_path": str(output_path.relative_to(BASE_DIR)).replace("\\", "/"),
        "records": records,
    }


def main() -> None:
    ensure_dir(OUTPUT_DIR)
    manifest = {
        "description": "Downloaded external self-contained sources for Java and C++ dataset construction",
        "sources": {
            "humanevalx_java": copy_humaneval_x("java"),
            "humanevalx_cpp": copy_humaneval_x("cpp"),
            "autocodebenchmark_java": filter_autocodebenchmark("java"),
            "autocodebenchmark_cpp": filter_autocodebenchmark("cpp"),
        },
    }
    write_json(MANIFEST_PATH, manifest)
    print(f"External sources synced to: {OUTPUT_DIR}")


if __name__ == "__main__":
    main()
