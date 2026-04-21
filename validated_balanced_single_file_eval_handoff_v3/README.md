# Validated Balanced Single-File Eval Dataset

- Generated at: `2026-04-19T13:51:41`
- Target: `50` samples per language per category
- Validation: real compile and/or execute check on every selected sample before packaging

## Validation Rules

- Python: `python sample.py`
- JavaScript: `node --check sample.js` + `node sample.js`
- Java: `javac` + generated `Runner` using `Class.forName(...)` + `java Runner`
- C++: generated `runner.cpp` including the sample + `g++` + run executable
- Go: generated `go.mod` and optional empty `main()` wrapper + `go run .`

## Summary

- `python`: selected `200` / target `200`; attempted validations `200`; status `OK`
- `javascript`: selected `200` / target `200`; attempted validations `203`; status `OK`
- `java`: selected `200` / target `200`; attempted validations `202`; status `OK`
- `cpp`: selected `200` / target `200`; attempted validations `258`; status `OK`
- `go`: selected `200` / target `200`; attempted validations `463`; status `OK`