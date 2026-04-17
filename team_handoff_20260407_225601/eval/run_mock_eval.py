import os
import sys
import random

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from eval.dataset import load_code_files, EvalResult, ModelEvalReport
from eval.evaluator import print_report_summary
from eval.scorer import (
    evaluate_python_test,
    evaluate_java_test,
    evaluate_cpp_test,
    evaluate_go_test
)

MODEL_NAME = "gpt-4o"

MOCK_TESTS = {
    "simple_function_python_simple_function_0.py": '''```python
import unittest
from typing import List

def has_close_elements(numbers: List[float], threshold: float) -> bool:
    for idx, elem in enumerate(numbers):
        for idx2, elem2 in enumerate(numbers):
            if idx != idx2:
                distance = abs(elem - elem2)
                if distance < threshold:
                    return True
    return False

class TestHasCloseElements(unittest.TestCase):
    def test_normal_case_close_elements(self):
        self.assertTrue(has_close_elements([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3))

    def test_normal_case_no_close_elements(self):
        self.assertFalse(has_close_elements([1.0, 2.0, 3.0], 0.5))

    def test_empty_list(self):
        self.assertFalse(has_close_elements([], 0.5))

    def test_single_element(self):
        self.assertFalse(has_close_elements([1.0], 0.5))

    def test_two_elements_close(self):
        self.assertTrue(has_close_elements([1.0, 1.1], 0.2))

    def test_two_elements_far(self):
        self.assertFalse(has_close_elements([1.0, 5.0], 2.0))

    def test_threshold_zero(self):
        self.assertFalse(has_close_elements([1.0, 2.0, 3.0], 0.0))

    def test_negative_numbers(self):
        self.assertTrue(has_close_elements([-1.0, -0.9, 5.0], 0.2))

    def test_large_threshold(self):
        self.assertTrue(has_close_elements([1.0, 100.0, 200.0], 1000.0))

    def test_duplicate_elements(self):
        self.assertTrue(has_close_elements([1.0, 1.0, 2.0], 0.1))

if __name__ == "__main__":
    unittest.main()
```''',
    "boundary_python_boundary_0.py": '''```python
import unittest
from typing import List

def below_zero(operations: List[int]) -> bool:
    balance = 0
    for op in operations:
        balance += op
        if balance < 0:
            return True
    return False

class TestBelowZero(unittest.TestCase):
    def test_normal_positive_operations(self):
        self.assertFalse(below_zero([1, 2, 3]))

    def test_normal_negative_operations(self):
        self.assertTrue(below_zero([1, 2, -4, 5]))

    def test_empty_operations(self):
        self.assertFalse(below_zero([]))

    def test_single_deposit(self):
        self.assertFalse(below_zero([100]))

    def test_single_withdrawal_positive(self):
        self.assertFalse(below_zero([50]))

    def test_single_withdrawal_negative(self):
        self.assertTrue(below_zero([-1]))

    def test_boundary_exactly_zero(self):
        self.assertFalse(below_zero([1, -1]))

    def test_boundary_goes_negative_then_positive(self):
        self.assertTrue(below_zero([1, -2, 5]))

    def test_large_values(self):
        self.assertFalse(below_zero([1000000, -500000, -400000]))

    def test_all_withdrawals(self):
        self.assertTrue(below_zero([-1, -2, -3]))

    def test_alternating_operations(self):
        self.assertFalse(below_zero([10, -5, 10, -5, -5]))

    def test_boundary_balance_never_negative(self):
        self.assertFalse(below_zero([5, -3, -2]))

if __name__ == "__main__":
    unittest.main()
```''',
    "interface_mock_python_interface_mock_0.py": '''```python
import unittest
import itertools
from random import shuffle

def task_func(numbers=list(range(1, 3))):
    permutations = list(itertools.permutations(numbers))
    sum_diffs = 0
    for perm in permutations:
        perm = list(perm)
        shuffle(perm)
        diffs = [abs(perm[i] - perm[i+1]) for i in range(len(perm)-1)]
        sum_diffs += sum(diffs)
    avg_sum_diffs = sum_diffs / len(permutations)
    return avg_sum_diffs

class TestTaskFunc(unittest.TestCase):
    def test_default_input(self):
        result = task_func()
        self.assertIsInstance(result, float)
        self.assertGreater(result, 0)

    def test_small_list(self):
        result = task_func([1, 2])
        self.assertIsInstance(result, float)

    def test_larger_list(self):
        result = task_func([1, 2, 3, 4])
        self.assertIsInstance(result, float)
        self.assertGreater(result, 0)

    def test_single_element(self):
        result = task_func([1])
        self.assertEqual(result, 0)

    def test_negative_numbers(self):
        result = task_func([-1, -2, -3])
        self.assertIsInstance(result, float)

    def test_mixed_numbers(self):
        result = task_func([-5, 0, 5])
        self.assertIsInstance(result, float)

if __name__ == "__main__":
    unittest.main()
```''',
    "complex_dependency_python_complex_dependency_0.py": '''```python
import unittest
import random
import statistics

def task_func(LETTERS):
    random_dict = {k: [random.randint(0, 100) for _ in range(random.randint(1, 10))] for k in LETTERS}
    sorted_dict = dict(sorted(random_dict.items(), key=lambda item: statistics.mean(item[1]), reverse=True))
    return sorted_dict

class TestTaskFunc(unittest.TestCase):
    def setUp(self):
        random.seed(42)

    def test_basic_functionality(self):
        LETTERS = ['A', 'B', 'C']
        result = task_func(LETTERS)
        self.assertIsInstance(result, dict)
        self.assertEqual(len(result), 3)

    def test_sorted_by_mean_descending(self):
        LETTERS = ['A', 'B', 'C']
        result = task_func(LETTERS)
        means = [statistics.mean(v) for v in result.values()]
        for i in range(len(means) - 1):
            self.assertGreaterEqual(means[i], means[i+1])

    def test_empty_letters(self):
        result = task_func([])
        self.assertEqual(result, {})

    def test_single_letter(self):
        LETTERS = ['X']
        result = task_func(LETTERS)
        self.assertEqual(len(result), 1)

    def test_many_letters(self):
        LETTERS = list('ABCDEFGHIJKLMNOPQRSTUVWXYZ')
        result = task_func(LETTERS)
        self.assertEqual(len(result), 26)

if __name__ == "__main__":
    unittest.main()
```''',
    "simple_function_cpp_simple_function_0.cpp": '''```cpp
#include <gtest/gtest.h>
#include <vector>
#include <cmath>

bool has_close_elements(std::vector<float> numbers, float threshold) {
    for (size_t i = 0; i < numbers.size(); i++) {
        for (size_t j = i + 1; j < numbers.size(); j++) {
            if (std::abs(numbers[i] - numbers[j]) < threshold) {
                return true;
            }
        }
    }
    return false;
}

TEST(HasCloseElementsTest, NormalCaseCloseElements) {
    EXPECT_TRUE(has_close_elements({1.0f, 2.8f, 3.0f, 4.0f, 5.0f, 2.0f}, 0.3f));
}

TEST(HasCloseElementsTest, NormalCaseNoCloseElements) {
    EXPECT_FALSE(has_close_elements({1.0f, 2.0f, 3.0f}, 0.5f));
}

TEST(HasCloseElementsTest, EmptyVector) {
    EXPECT_FALSE(has_close_elements({}, 0.5f));
}

TEST(HasCloseElementsTest, SingleElement) {
    EXPECT_FALSE(has_close_elements({1.0f}, 0.5f));
}

TEST(HasCloseElementsTest, TwoElementsClose) {
    EXPECT_TRUE(has_close_elements({1.0f, 1.1f}, 0.2f));
}

TEST(HasCloseElementsTest, TwoElementsFar) {
    EXPECT_FALSE(has_close_elements({1.0f, 5.0f}, 2.0f));
}

TEST(HasCloseElementsTest, ThresholdZero) {
    EXPECT_FALSE(has_close_elements({1.0f, 2.0f, 3.0f}, 0.0f));
}

TEST(HasCloseElementsTest, NegativeNumbers) {
    EXPECT_TRUE(has_close_elements({-1.0f, -0.9f, 5.0f}, 0.2f));
}

TEST(HasCloseElementsTest, DuplicateElements) {
    EXPECT_TRUE(has_close_elements({1.0f, 1.0f, 2.0f}, 0.1f));
}
```''',
    "boundary_cpp_boundary_0.cpp": '''```cpp
#include <gtest/gtest.h>
#include <vector>

bool below_zero(std::vector<long> operations) {
    long balance = 0;
    for (long op : operations) {
        balance += op;
        if (balance < 0) {
            return true;
        }
    }
    return false;
}

TEST(BelowZeroTest, NormalPositiveOperations) {
    EXPECT_FALSE(below_zero({1, 2, 3}));
}

TEST(BelowZeroTest, NormalNegativeOperations) {
    EXPECT_TRUE(below_zero({1, 2, -4, 5}));
}

TEST(BelowZeroTest, EmptyOperations) {
    EXPECT_FALSE(below_zero({}));
}

TEST(BelowZeroTest, SingleDeposit) {
    EXPECT_FALSE(below_zero({100}));
}

TEST(BelowZeroTest, SingleNegativeWithdrawal) {
    EXPECT_TRUE(below_zero({-1}));
}

TEST(BelowZeroTest, BoundaryExactlyZero) {
    EXPECT_FALSE(below_zero({1, -1}));
}

TEST(BelowZeroTest, BoundaryGoesNegativeThenPositive) {
    EXPECT_TRUE(below_zero({1, -2, 5}));
}

TEST(BelowZeroTest, AllWithdrawals) {
    EXPECT_TRUE(below_zero({-1, -2, -3}));
}

TEST(BelowZeroTest, LargeValues) {
    EXPECT_FALSE(below_zero({1000000, -500000, -400000}));
}
```''',
    "complex_dependency_cpp_complex_dependency_0.cpp": '''```cpp
#include <gtest/gtest.h>
#include <string>
#include <memory>

namespace google::api::expr::codelab {
absl::StatusOr<std::string> ParseAndEvaluate(absl::string_view cel_expr) {
    if (cel_expr.empty()) {
        return absl::InvalidArgumentError("Empty expression");
    }
    return std::string(cel_expr);
}
}

TEST(ParseAndEvaluateTest, ValidExpression) {
    auto result = google::api::expr::codelab::ParseAndEvaluate("1 + 1");
    ASSERT_TRUE(result.ok());
    EXPECT_EQ(result.value(), "1 + 1");
}

TEST(ParseAndEvaluateTest, EmptyExpression) {
    auto result = google::api::expr::codelab::ParseAndEvaluate("");
    EXPECT_FALSE(result.ok());
}

TEST(ParseAndEvaluateTest, ComplexExpression) {
    auto result = google::api::expr::codelab::ParseAndEvaluate("a && b || c");
    ASSERT_TRUE(result.ok());
}
```''',
    "simple_function_go_simple_function_0.go": '''```go
package main

import (
	"testing"
)

func freeTo(in *inflights, to uint64) {
	if in.count == 0 || to < in.buffer[in.start] {
		return
	}
	idx := in.start
	var i int
	for i = 0; i < in.count; i++ {
		if to < in.buffer[idx] {
			break
		}
		size := in.size
		if idx++; idx >= size {
			idx -= size
		}
	}
	in.count -= i
	in.start = idx
	if in.count == 0 {
		in.start = 0
	}
}

type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

func TestFreeTo(t *testing.T) {
	t.Run("empty inflights", func(t *testing.T) {
		in := &inflights{buffer: make([]uint64, 10), count: 0, start: 0, size: 10}
		freeTo(in, 5)
		if in.count != 0 {
			t.Errorf("expected count 0, got %d", in.count)
		}
	})

	t.Run("free all inflights", func(t *testing.T) {
		in := &inflights{buffer: []uint64{1, 2, 3}, count: 3, start: 0, size: 3}
		freeTo(in, 10)
		if in.count != 0 {
			t.Errorf("expected count 0, got %d", in.count)
		}
	})

	t.Run("free partial inflights", func(t *testing.T) {
		in := &inflights{buffer: []uint64{1, 2, 3, 4, 5}, count: 5, start: 0, size: 5}
		freeTo(in, 2)
		if in.count != 3 {
			t.Errorf("expected count 3, got %d", in.count)
		}
	})
}
```''',
    "boundary_go_boundary_0.go": '''```go
package main

import (
	"testing"
)

func maybeDecrTo(pr *Progress, rejected, last uint64) bool {
	if pr.State == ProgressStateReplicate {
		if rejected <= pr.Match {
			return false
		}
		pr.Next = pr.Match + 1
		return true
	}
	if pr.Next-1 != rejected {
		return false
	}
	if pr.Next = min(rejected, last+1); pr.Next < 1 {
		pr.Next = 1
	}
	pr.resume()
	return true
}

type Progress struct {
	State ProgressStateType
	Match uint64
	Next  uint64
}

type ProgressStateType int

const (
	ProgressStateReplicate ProgressStateType = iota
	ProgressStateProbe
)

func (pr *Progress) resume() {}

func TestMaybeDecrTo(t *testing.T) {
	t.Run("replicate state stale rejection", func(t *testing.T) {
		pr := &Progress{State: ProgressStateReplicate, Match: 10, Next: 15}
		result := maybeDecrTo(pr, 5, 20)
		if result != false {
			t.Errorf("expected false, got %v", result)
		}
	})

	t.Run("replicate state valid rejection", func(t *testing.T) {
		pr := &Progress{State: ProgressStateReplicate, Match: 5, Next: 15}
		result := maybeDecrTo(pr, 10, 20)
		if result != true {
			t.Errorf("expected true, got %v", result)
		}
		if pr.Next != 6 {
			t.Errorf("expected Next=6, got %d", pr.Next)
		}
	})

	t.Run("probe state stale rejection", func(t *testing.T) {
		pr := &Progress{State: ProgressStateProbe, Match: 5, Next: 10}
		result := maybeDecrTo(pr, 5, 20)
		if result != false {
			t.Errorf("expected false, got %v", result)
		}
	})
}
```''',
    "complex_dependency_go_complex_dependency_0.go": '''```go
package main

import (
	"testing"
)

func Read(lg *Logger, snapname string) ([]byte, error) {
	if snapname == "" {
		return nil, ErrEmptySnapshot
	}
	return []byte(snapname), nil
}

type Logger struct{}
type ErrEmptySnapshot struct{}

func (e ErrEmptySnapshot) Error() string { return "empty snapshot" }

func TestRead(t *testing.T) {
	t.Run("empty snapshot name", func(t *testing.T) {
		lg := &Logger{}
		_, err := Read(lg, "")
		if err == nil {
			t.Error("expected error for empty snapshot name")
		}
	})

	t.Run("valid snapshot name", func(t *testing.T) {
		lg := &Logger{}
		data, err := Read(lg, "test.snap")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if string(data) != "test.snap" {
			t.Errorf("expected 'test.snap', got '%s'", string(data))
		}
	})
}
```''',
    "interface_mock_go_interface_mock_0.go": '''```go
package main

import (
	"testing"
)

func (in *inflights) freeTo(to uint64) {
	if in.count == 0 || to < in.buffer[in.start] {
		return
	}
	idx := in.start
	var i int
	for i = 0; i < in.count; i++ {
		if to < in.buffer[idx] {
			break
		}
		size := in.size
		if idx++; idx >= size {
			idx -= size
		}
	}
	in.count -= i
	in.start = idx
	if in.count == 0 {
		in.start = 0
	}
}

type inflights struct {
	buffer []uint64
	start  int
	count  int
	size   int
}

func TestInflightsFreeTo(t *testing.T) {
	t.Run("empty inflights", func(t *testing.T) {
		in := &inflights{buffer: make([]uint64, 10), count: 0, start: 0, size: 10}
		in.freeTo(5)
		if in.count != 0 {
			t.Errorf("expected count 0, got %d", in.count)
		}
	})

	t.Run("free all", func(t *testing.T) {
		in := &inflights{buffer: []uint64{1, 2, 3}, count: 3, start: 0, size: 3}
		in.freeTo(10)
		if in.count != 0 {
			t.Errorf("expected count 0, got %d", in.count)
		}
	})

	t.Run("free partial", func(t *testing.T) {
		in := &inflights{buffer: []uint64{1, 2, 3, 4, 5}, count: 5, start: 0, size: 5}
		in.freeTo(2)
		if in.count != 3 {
			t.Errorf("expected count 3, got %d", in.count)
		}
	})

	t.Run("out of left side", func(t *testing.T) {
		in := &inflights{buffer: []uint64{5, 6, 7}, count: 3, start: 0, size: 3}
		in.freeTo(1)
		if in.count != 3 {
			t.Errorf("expected count 3, got %d", in.count)
		}
	})
}
```''',
    "simple_function_java_simple_function_0.java": '''```java
import org.junit.Test;
import static org.junit.Assert.*;
import java.util.*;

public class NameParserTest {

    @Test
    public void testParseSimplePath() {
        NameParser parser = new NameParser();
        Collection<NameToken> tokens = parser.parse("/media/aap");
        assertEquals(2, tokens.size());
    }

    @Test
    public void testParsePathWithNumericRange() {
        NameParser parser = new NameParser();
        Collection<NameToken> tokens = parser.parse("/media/aap/[1..100]");
        assertEquals(3, tokens.size());
    }

    @Test
    public void testParseEmptyString() {
        NameParser parser = new NameParser();
        Collection<NameToken> tokens = parser.parse("");
        assertEquals(0, tokens.size());
    }

    @Test
    public void testParseOnlySlashes() {
        NameParser parser = new NameParser();
        Collection<NameToken> tokens = parser.parse("///");
        assertEquals(0, tokens.size());
    }
}
```''',
    "boundary_java_boundary_0.java": '''```java
import org.junit.Test;
import static org.junit.Assert.*;
import java.util.*;

public class ThreadDeadlockDetectorTest {

    @Test
    public void returnsAnEmptySetIfNoThreadsAreDeadlocked() {
        ThreadMXBean mockThreads = mock(ThreadMXBean.class);
        when(mockThreads.findDeadlockedThreads()).thenReturn(null);
        ThreadDeadlockDetector detector = new ThreadDeadlockDetector(mockThreads);
        Set<String> result = detector.getDeadlockedThreads();
        assertTrue(result.isEmpty());
    }

    @Test
    public void returnsDeadlockedThreadNames() {
        ThreadMXBean mockThreads = mock(ThreadMXBean.class);
        when(mockThreads.findDeadlockedThreads()).thenReturn(new long[]{1L, 2L});
        ThreadDeadlockDetector detector = new ThreadDeadlockDetector(mockThreads);
        Set<String> result = detector.getDeadlockedThreads();
        assertEquals(2, result.size());
    }
}
```''',
    "interface_mock_java_interface_mock_0.java": '''```java
import org.junit.Test;
import static org.junit.Assert.*;
import java.io.*;

public class XlsLegendParserTest {

    @Test(expected = IllegalArgumentException.class)
    public void testGlobalLegendParsingRegionalSheet() throws Exception {
        XlsLegendParser parser = new XlsLegendParser();
        InputStream inputStream = new ByteArrayInputStream(new byte[0]);
        parser.parse(inputStream, false);
    }

    @Test
    public void testParseRegionalLegend() throws Exception {
        XlsLegendParser parser = new XlsLegendParser();
        InputStream inputStream = getClass().getResourceAsStream("/test_regional.xls");
        if (inputStream != null) {
            LegendClass[] classes = parser.parse(inputStream, true);
            assertNotNull(classes);
        }
    }

    @Test
    public void testParseGlobalLegend() throws Exception {
        XlsLegendParser parser = new XlsLegendParser();
        InputStream inputStream = getClass().getResourceAsStream("/test_global.xls");
        if (inputStream != null) {
            LegendClass[] classes = parser.parse(inputStream, false);
            assertNotNull(classes);
        }
    }
}
```''',
    "complex_dependency_java_complex_dependency_0.java": '''```java
import org.junit.Test;
import static org.junit.Assert.*;
import java.io.*;

public class LocalAndroidPlatformsTest {

    @Test
    public void testFindLocalJavaSdk() {
        try {
            File sdk = LocalAndroidPlatforms.findLocalJavaSdk();
            if (sdk != null) {
                assertTrue(sdk.exists());
            }
        } catch (FileNotFoundException e) {
            assertTrue(e.getMessage().contains("Unable to find"));
        } catch (IOException e) {
            fail("Unexpected IOException: " + e.getMessage());
        }
    }

    @Test
    public void testFindLocalJavaSdkWithEnvVariable() {
        String originalHome = System.getenv("ANDROID_HOME");
        try {
            File sdk = LocalAndroidPlatforms.findLocalJavaSdk();
        } catch (IOException e) {
        }
    }
}
```'''
}


def generate_mock_test(filename: str) -> str:
    return MOCK_TESTS.get(filename, f"// Mock test for {filename}\n@Test\npublic void testPlaceholder() {{\n    // TODO: implement\n}}")


def main():
    print("Loading code files from data/basic_tests ...")
    code_files = load_code_files("data/basic_tests")
    print(f"Loaded {len(code_files)} code files:")
    for cf in code_files:
        print(f"  {cf.language:6s} | {cf.category:20s} | {cf.filename}")

    print(f"\nStarting MOCK evaluation with {MODEL_NAME} ...")
    print("(Using pre-generated test cases for demonstration)")
    print()

    report = ModelEvalReport(model_name=MODEL_NAME)

    scorers = {
        "python": evaluate_python_test,
        "java": evaluate_java_test,
        "cpp": evaluate_cpp_test,
        "go": evaluate_go_test
    }

    for i, cf in enumerate(code_files):
        print(f"[{i+1}/{len(code_files)}] Evaluating {cf.language}/{cf.filename} ...")
        generated_test = generate_mock_test(cf.filename)
        scorer = scorers.get(cf.language)
        if scorer:
            scores = scorer(cf.code, generated_test)
        else:
            scores = {"syntax_validity": 0.8, "test_coverage": 0.7, "assertion_quality": 0.7,
                      "edge_case_handling": 0.6, "code_structure": 0.8, "total": 0.72}

        total_score = scores.get("total", 0.0)
        result = EvalResult(
            language=cf.language,
            category=cf.category,
            filename=cf.filename,
            generated_test=generated_test,
            scores={k: v for k, v in scores.items() if k != "total"},
            total_score=total_score,
            model_name=MODEL_NAME
        )
        report.results.append(result)
        print(f"  Score: {total_score:.2f}/1.00")

    print_report_summary(report)

    import json
    from datetime import datetime
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    report_data = {
        "model_name": report.model_name,
        "total_files": len(report.results),
        "avg_total_score": round(report.avg_total_score, 4),
        "avg_by_language": {k: round(v, 4) for k, v in report.avg_by_language.items()},
        "avg_by_category": {k: round(v, 4) for k, v in report.avg_by_category.items()},
        "avg_by_dimension": {k: round(v, 4) for k, v in report.avg_by_dimension.items()},
        "results": [
            {
                "language": r.language,
                "category": r.category,
                "filename": r.filename,
                "total_score": r.total_score,
                "scores": r.scores,
                "generated_test": r.generated_test[:1000]
            }
            for r in report.results
        ]
    }
    os.makedirs("eval/results", exist_ok=True)
    filepath = f"eval/results/{MODEL_NAME.replace('-', '_')}_mock_{timestamp}.json"
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(report_data, f, indent=2, ensure_ascii=False)
    print(f"\nJSON report saved to: {filepath}")


if __name__ == "__main__":
    main()
