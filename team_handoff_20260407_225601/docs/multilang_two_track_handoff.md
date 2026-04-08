# Multilingual Two-Track Dataset Handoff

## Overview

The workspace now follows a two-track dataset strategy:

- `self_contained`: single-file evaluation in an isolated sandbox
- `module_level`: project-linked or module-level evaluation with repository/package context

## Recommended datasets

### Python

- Self-contained:
  - `C:\shijian_project\data\python_ut_dataset\python_dataset_self_contained.json`
  - `C:\shijian_project\data\python_code_files_self_contained`
- Module-level:
  - `C:\shijian_project\data\python_ut_dataset\python_dataset_module_level.json`
  - `C:\shijian_project\data\python_code_files_module_level`

### Go

- Self-contained:
  - `C:\shijian_project\data\go_ut_dataset\go_dataset_self_contained.json`
  - `C:\shijian_project\data\go_code_files_self_contained`
- Module-level:
  - `C:\shijian_project\data\go_ut_dataset\go_dataset_module_level.json`
  - `C:\shijian_project\data\go_code_files_module_level`

### Java

- Self-contained:
  - `C:\shijian_project\data\java_ut_dataset\java_dataset_self_contained.json`
  - `C:\shijian_project\data\java_code_files_self_contained`
- Module-level:
  - `C:\shijian_project\data\java_ut_dataset\java_dataset_module_level.json`
  - `C:\shijian_project\data\java_code_files_module_level`

### C++

- Self-contained:
  - `C:\shijian_project\data\cpp_ut_dataset\cpp_dataset_self_contained.json`
  - `C:\shijian_project\data\cpp_code_files_self_contained`
- Module-level:
  - `C:\shijian_project\data\cpp_ut_dataset\cpp_dataset_module_level.json`
  - `C:\shijian_project\data\cpp_code_files_module_level`

### JavaScript

- Self-contained:
  - `C:\shijian_project\data\javascript_ut_dataset\javascript_dataset_self_contained.json`
  - `C:\shijian_project\data\javascript_code_files_self_contained`
- Module-level:
  - `C:\shijian_project\data\javascript_ut_dataset\javascript_dataset_module_level.json`
  - `C:\shijian_project\data\javascript_code_files_module_level`

## Source summary

- Python self-contained: BigCodeBench complete solutions, filtered to stdlib-only single-file inputs
- Python module-level: real project modules from installed packages and local repos
- Go self-contained: stdlib-only real functions from `data/go/go_all.jsonl`
- Go module-level: project-linked real functions from `data/go/go_all.jsonl`
- Java self-contained: HumanEval-X plus AutoCodeBenchmark canonical solutions synced into `data/external_sources/java_cpp_self_contained`
- Java module-level: full Java source files from `data/bcb/code/hm-dianping`
- C++ self-contained: HumanEval-X plus AutoCodeBenchmark canonical solutions synced into `data/external_sources/java_cpp_self_contained`
- C++ module-level: project-linked modules from `data/cpp_ut/cpp_ut_all.jsonl`
- JavaScript self-contained: lodash helper closures inlined into standalone bundles
- JavaScript module-level: project-linked function slices from `js_sources/lodash.js`

## External source sync

Java and C++ self-contained datasets now depend on a synced local source cache:

- `data/external_sources/java_cpp_self_contained/humanevalx_java.jsonl`
- `data/external_sources/java_cpp_self_contained/humanevalx_cpp.jsonl`
- `data/external_sources/java_cpp_self_contained/autocodebenchmark_java.jsonl`
- `data/external_sources/java_cpp_self_contained/autocodebenchmark_cpp.jsonl`

These files are created by:

```powershell
python sync_external_self_contained_sources.py
```

## Build commands

Build everything:

```powershell
python build_multilang_two_track_datasets.py
```

Build individual tracks:

```powershell
python sync_external_self_contained_sources.py
python build_self_contained_python_dataset.py
python build_python_module_level_dataset.py
python build_go_two_track_datasets.py
python build_java_module_level_dataset.py
python build_cpp_module_level_dataset.py
python build_javascript_two_track_datasets.py
```

## Evaluation guidance

- Use `self_contained` datasets for reproducible single-file benchmarking
- Use `module_level` datasets for realistic project-context evaluation
- Do not compare self-contained and module-level scores as if they measure the same difficulty

## Project-Level Readiness

The `module_level` datasets are ready for all five languages, but the execution readiness is different:

| Language | `module_level` dataset ready | Project restoration ready | Direct project-level evaluation status |
|----------|------------------------------|---------------------------|----------------------------------------|
| Python | Yes | Partial | Demo-ready with `python_project_module_eval.py` |
| Go | Yes | No | Needs repository/environment restoration work |
| Java | Yes | No | Needs Maven/Spring project restoration work |
| C++ | Yes | No | Needs build-system and dependency restoration work |
| JavaScript | Yes | No | Needs package/module restoration work |

Important:

- `module_level` means the evaluation target is a real project module
- It does not mean the full original repository has already been packaged for every language
- For Python, a runnable restoration/evaluation example already exists
- For the other four languages, teammates still need to implement restoration and execution scripts before full project-level evaluation
