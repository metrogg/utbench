# Python Eval Dataset Handoff

## Recommended dataset

Use the self-contained Python dataset for single-file sandbox evaluation:

- `C:\shijian_project\data\python_ut_dataset\python_dataset_self_contained.json`
- `C:\shijian_project\data\python_code_files_self_contained`

This dataset contains 200 Python samples with 4 balanced categories:

- `simple_function`
- `boundary`
- `complex_dependency`
- `interface_mock`

Each sample is filtered to satisfy the basic single-file evaluation contract:

- no relative imports
- no package-internal imports
- stdlib-only imports

## Why this replaces the old recommendation

The previous `realworld` dataset was built from module files inside installed packages and open-source repos.
That selection logic optimized for complexity, but it did not enforce single-file executability.
As a result, many samples depended on:

- relative imports such as `from . import ...`
- package-internal modules such as `from pandas.core...`
- third-party packages such as `numpy`, `pandas`, `datasets`, or `pyarrow`

Those files are not appropriate for isolated sandbox evaluation unless the full package tree and environment are recreated.

## Realworld dataset status

`build_realworld_module_dataset.py` now includes the self-contained precheck.
If the source pool cannot provide enough valid candidates, the build should fail instead of emitting a misleading dataset.

Keep the realworld dataset only for project-level or environment-rich experiments where package context is intentionally available.

## Build commands

Generate the self-contained dataset:

```powershell
python build_self_contained_python_dataset.py
```

Validate a dataset file:

```powershell
python verify_dataset.py
python verify_dataset.py data/python_ut_dataset/python_dataset_realworld.json
```

## Validation rules

Validation is implemented in `python_dataset_precheck.py` and currently flags:

- `relative_import`
- `internal_package_import`
- `third_party_import`

## Output shape

The self-contained dataset keeps the same basic schema used by the other Python dataset builders:

- `samples[*].code`
- `samples[*].category`
- `samples[*].metrics`
- `samples[*].features`
- `samples[*].category_scores`

It also records `self_contained_check` for each sample.
