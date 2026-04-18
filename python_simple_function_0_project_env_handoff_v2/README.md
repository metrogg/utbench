# Python Project-Level Sample Handoff

- Generated at: `2026-04-17T19:43:50`
- Package name: `python_simple_function_0_project_env_handoff_v2`
- Sample ID: `python_simple_function_0`
- Category: `simple_function`
- Package: `dateutil`
- Target module: `venv/Lib/site-packages/dateutil/relativedelta.py`
- Module import path: `dateutil.relativedelta`

## Contents

- `workspace/`: restored project workspace ready to run
- `results/`: saved report, coverage, and pytest XML from the successful run
- `env/requirements_python_project_eval.txt`: Python dependencies for rerun
- `sample_metadata.json`: original dataset entry for this sample

## Current Result

- Compile check: `passed`
- Pytest: `7/7` passed
- failures: `0`
- errors: `0`
- coverage: `58%`

## How To Use On Another Machine

1. Install Python 3.11+.
2. Create a virtual environment in this package root.
3. Install `env/requirements_python_project_eval.txt`.
4. Run pytest or coverage against the restored workspace.

```powershell
python -m venv venv
.\venv\Scripts\python.exe -m pip install --upgrade pip
.\venv\Scripts\python.exe -m pip install -r .\env\requirements_python_project_eval.txt
$env:PYTHONPATH = (Resolve-Path .\workspace)
.\venv\Scripts\python.exe -m pytest .\workspace\tests\test_generated_relativedelta.py -q
.\venv\Scripts\python.exe -m coverage run --branch --source dateutil -m pytest .\workspace\tests\test_generated_relativedelta.py -q
.\venv\Scripts\python.exe -m coverage json -o .\results\coverage_rerun.json
```

## Notes

- This package already includes the restored `dateutil` package structure.
- The generated test file is `workspace/tests/test_generated_relativedelta.py`.
- If you only want the run result, open `results/report.json`.

