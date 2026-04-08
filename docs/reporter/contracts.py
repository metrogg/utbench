from __future__ import annotations

_LANGUAGE_TOKENS = {"python", "java", "go", "cpp", "javascript"}

_COLUMN_LABELS_ZH = {
    "model": "妯″瀷",
    "language": "璇█",
    "total_samples": "鏍锋湰鏁?",
    "compile_pass_rate": "缂栬瘧閫氳繃鐜?",
    "test_pass_rate": "娴嬭瘯閫氳繃鐜?",
    "avg_line_coverage": "骞冲潎琛岃鐩栫巼",
    "avg_branch_coverage": "骞冲潎鍒嗘敮瑕嗙洊鐜?",
    "avg_function_coverage": "骞冲潎鍑芥暟瑕嗙洊鐜?",
    "avg_mutation_score": "骞冲潎鍙樺紓寰楀垎",
    "avg_generation_latency_ms": "骞冲潎鐢熸垚鑰楁椂(ms)",
    "avg_total_tokens": "骞冲潎鎬籘oken",
    "avg_eval_runtime_ms": "骞冲潎璇勬祴鑰楁椂(ms)",
    "avg_prompt_tokens": "骞冲潎杈撳叆Token",
    "avg_completion_tokens": "骞冲潎杈撳嚭Token",
    "sample_id": "鏍锋湰ID",
    "complexity": "澶嶆潅搴?",
    "scenario": "鍦烘櫙",
    "compile_pass": "缂栬瘧閫氳繃",
    "test_pass": "娴嬭瘯閫氳繃",
    "line_coverage": "琛岃鐩栫巼",
    "branch_coverage": "鍒嗘敮瑕嗙洊鐜?",
    "mutation_score": "鍙樺紓寰楀垎",
    "mutation_total": "鍙樺紓浣撴€绘暟",
    "mutation_killed": "宸叉潃姝?",
    "mutation_survived": "瀛樻椿",
    "mutation_no_tests": "鏈鐩?no_tests)",
    "mutation_not_checked": "鏈墽琛?not_checked)",
    "mutation_timeout": "瓒呮椂(timeout)",
    "mutation_skipped": "璺宠繃(skipped)",
    "mutation_suspicious": "鍙枒(suspicious)",
    "generation_latency_ms": "鐢熸垚鑰楁椂(ms)",
    "total_tokens": "鎬籘oken",
    "generation_metrics_source": "鐢熸垚鎸囨爣鏉ユ簮",
    "generation_metrics_reason": "鐢熸垚鎸囨爣缂哄け鍘熷洜",
    "stage": "闃舵",
    "error_type": "閿欒绫诲瀷",
    "count": "鏁伴噺",
    "example_model": "绀轰緥妯″瀷",
    "example_sample": "绀轰緥鏍锋湰",
    "example_message": "绀轰緥淇℃伅",
    "test_error": "娴嬭瘯閿欒",
    "coverage_error": "瑕嗙洊鐜囬敊璇?",
    "mutation_error": "鍙樺紓閿欒",
}

_GEN_METRICS_SOURCE_LABELS_ZH = {
    "metadata": "metadata",
    "runner_summary": "runner_summary 鍥炲～",
    "missing": "缂哄け",
}

_GEN_METRICS_REASON_LABELS_ZH = {
    "metadata_missing_fallback_runner_summary": "metadata 缂哄け锛屽凡鍥炲～ runner_summary",
    "metadata_and_runner_summary_missing": "metadata 涓?runner_summary 鍧囩己澶?",
}

_STAGE_LABELS_ZH = {
    "compile": "缂栬瘧",
    "test": "娴嬭瘯",
    "coverage": "瑕嗙洊鐜?",
    "mutation": "鍙樺紓",
}

_ERROR_TYPE_LABELS_ZH = {
    "module_not_found": "妯″潡缂哄け",
    "name_error": "鍚嶇О閿欒",
    "assertion_failure": "鏂█澶辫触",
    "import_error": "瀵煎叆閿欒",
    "syntax_error": "璇硶閿欒",
    "timeout": "瓒呮椂",
    "stop_iteration": "杩唬鍋滄",
    "recursion_error": "閫掑綊閿欒",
    "other": "鍏朵粬",
}
