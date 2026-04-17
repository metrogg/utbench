import argparse
import json
import shutil
import subprocess
from datetime import datetime
from pathlib import Path


BASE_DIR = Path(__file__).resolve().parent
DIST_DIR = BASE_DIR / "dist"
PYTHON_EXE = BASE_DIR / "venv" / "Scripts" / "python.exe"


SCRIPT_FILES = [
    "package_team_handoff.py",
    "build_multilang_two_track_datasets.py",
    "build_self_contained_python_dataset.py",
    "build_python_module_level_dataset.py",
    "build_go_two_track_datasets.py",
    "build_java_module_level_dataset.py",
    "build_cpp_module_level_dataset.py",
    "build_javascript_two_track_datasets.py",
    "sync_external_self_contained_sources.py",
    "python_dataset_precheck.py",
    "dataset_builder_utils.py",
    "verify_dataset.py",
    "python_project_module_eval.py",
]

EVAL_FILES = [
    "eval/dataset.py",
    "eval/evaluator.py",
    "eval/llm_client.py",
    "eval/scorer.py",
    "eval/run_eval.py",
    "eval/run_mock_eval.py",
    "eval/run_ollama_eval.py",
]

DOC_FILES = [
    "multilang_two_track_handoff.md",
    "python_eval_dataset_handoff.md",
    "TEAM_HANDOFF_GUIDE.md",
    "setup_python_handoff_env.ps1",
]

DATA_MAPPINGS = [
    ("data/python_ut_dataset/python_dataset_self_contained.json", "data/python_ut_dataset/python_dataset_self_contained.json"),
    ("data/python_code_files_self_contained", "data/python_code_files_self_contained"),
    ("data/python_ut_dataset/python_dataset_module_level.json", "data/python_ut_dataset/python_dataset_module_level.json"),
    ("data/python_code_files_module_level", "data/python_code_files_module_level"),
    ("data/go_ut_dataset/go_dataset_self_contained.json", "data/go_ut_dataset/go_dataset_self_contained.json"),
    ("data/go_code_files_self_contained", "data/go_code_files_self_contained"),
    ("data/go_ut_dataset/go_dataset_module_level.json", "data/go_ut_dataset/go_dataset_module_level.json"),
    ("data/go_code_files_module_level", "data/go_code_files_module_level"),
    ("data/javascript_ut_dataset/javascript_dataset_self_contained.json", "data/javascript_ut_dataset/javascript_dataset_self_contained.json"),
    ("data/javascript_code_files_self_contained", "data/javascript_code_files_self_contained"),
    ("data/javascript_ut_dataset/javascript_dataset_module_level.json", "data/javascript_ut_dataset/javascript_dataset_module_level.json"),
    ("data/javascript_code_files_module_level", "data/javascript_code_files_module_level"),
    ("data/java_ut_dataset/java_dataset_self_contained.json", "data/java_ut_dataset/java_dataset_self_contained.json"),
    ("data/java_ut_dataset/java_dataset_module_level.json", "data/java_ut_dataset/java_dataset_module_level.json"),
    ("data/java_code_files_module_level", "data/java_code_files_module_level"),
    ("data/java_code_files_self_contained", "data/java_code_files_self_contained"),
    ("data/cpp_ut_dataset/cpp_dataset_self_contained.json", "data/cpp_ut_dataset/cpp_dataset_self_contained.json"),
    ("data/cpp_ut_dataset/cpp_dataset_module_level.json", "data/cpp_ut_dataset/cpp_dataset_module_level.json"),
    ("data/cpp_code_files_module_level", "data/cpp_code_files_module_level"),
    ("data/cpp_code_files_self_contained", "data/cpp_code_files_self_contained"),
    ("data/external_sources/java_cpp_self_contained", "data/external_sources/java_cpp_self_contained"),
]


def copy_path(src: Path, dst: Path) -> None:
    if src.is_dir():
        shutil.copytree(
            src,
            dst,
            dirs_exist_ok=True,
            ignore=shutil.ignore_patterns("__pycache__", "*.pyc", "*.pyo"),
        )
    else:
        dst.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(src, dst)


def write_requirements(output_path: Path) -> None:
    proc = subprocess.run(
        [str(PYTHON_EXE), "-m", "pip", "freeze"],
        cwd=BASE_DIR,
        capture_output=True,
        text=True,
        check=True,
    )
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(proc.stdout.strip() + "\n", encoding="utf-8")


def latest_example_report() -> Path | None:
    candidates = sorted(
        BASE_DIR.glob("eval/project_level_runs/python_simple_function_0/*/report_*.json"),
        key=lambda path: path.stat().st_mtime,
        reverse=True,
    )
    return candidates[0] if candidates else None


def build_manifest(package_root: Path) -> None:
    manifest = {
        "generated_at": datetime.now().isoformat(timespec="seconds"),
        "package_root": str(package_root),
        "datasets": {},
    }

    data_root = package_root / "data"
    for language_dir in sorted(data_root.iterdir()):
        if not language_dir.is_dir():
            continue
        manifest["datasets"][language_dir.name] = {}
        files = []
        for path in language_dir.rglob("*"):
            if path.is_file():
                files.append(str(path.relative_to(package_root)).replace("\\", "/"))
        manifest["datasets"][language_dir.name] = {
            "file_count": len(files),
            "files": files[:20],
        }

    manifest_path = package_root / "MANIFEST.json"
    manifest_path.write_text(json.dumps(manifest, ensure_ascii=False, indent=2), encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser(description="Package the cleaned two-track datasets for teammate handoff.")
    parser.add_argument("--name", default=None, help="Optional output folder name under dist/")
    parser.add_argument("--zip", action="store_true", help="Also create a zip archive.")
    args = parser.parse_args()

    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    folder_name = args.name or f"team_handoff_{timestamp}"
    package_root = DIST_DIR / folder_name

    if package_root.exists():
        shutil.rmtree(package_root)

    (package_root / "docs").mkdir(parents=True, exist_ok=True)
    (package_root / "eval").mkdir(parents=True, exist_ok=True)
    (package_root / "env").mkdir(parents=True, exist_ok=True)
    (package_root / "examples" / "python_project_eval").mkdir(parents=True, exist_ok=True)

    for relative in SCRIPT_FILES:
        copy_path(BASE_DIR / relative, package_root / Path(relative).name)

    for relative in EVAL_FILES:
        destination = package_root / "eval" / Path(relative).name
        copy_path(BASE_DIR / relative, destination)

    for relative in DOC_FILES:
        copy_path(BASE_DIR / relative, package_root / "docs" / Path(relative).name)

    for src_rel, dst_rel in DATA_MAPPINGS:
        copy_path(BASE_DIR / src_rel, package_root / dst_rel)

    example_test = BASE_DIR / "eval" / "generated_tests" / "python_project_demo_dateutil_relativedelta.py"
    if example_test.exists():
        copy_path(
            example_test,
            package_root / "examples" / "python_project_eval" / "generated_test_dateutil_relativedelta.py",
        )

    example_report = latest_example_report()
    if example_report is not None:
        copy_path(
            example_report,
            package_root / "examples" / "python_project_eval" / "example_report.json",
        )

    write_requirements(package_root / "env" / "requirements_python_project_eval.txt")
    build_manifest(package_root)

    if args.zip:
        archive_path = shutil.make_archive(str(package_root), "zip", root_dir=package_root.parent, base_dir=package_root.name)
        print(f"ZIP: {archive_path}")

    print(f"PACKAGE: {package_root}")


if __name__ == "__main__":
    main()
