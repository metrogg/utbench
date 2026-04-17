import os
import sys

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from eval.dataset import load_code_files
from eval.evaluator import UnitTestEvaluator, print_report_summary

MODEL = "qwen2.5-coder:7b"
BASE_URL = "http://localhost:11434/v1"

def main():
    print("Loading code files from data/basic_tests ...")
    code_files = load_code_files("data/basic_tests")
    print(f"Loaded {len(code_files)} code files:")
    for cf in code_files:
        print(f"  {cf.language:6s} | {cf.category:20s} | {cf.filename}")

    print(f"\nStarting evaluation with {MODEL} (Ollama) ...")
    evaluator = UnitTestEvaluator(api_key="ollama", model=MODEL, base_url=BASE_URL, output_dir="eval/results")
    report = evaluator.run_evaluation(code_files)

    print_report_summary(report)

if __name__ == "__main__":
    main()
