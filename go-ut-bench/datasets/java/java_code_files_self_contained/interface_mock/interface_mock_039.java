// Converted Java method
import java.util.ArrayList;
import java.util.List;

class UnitTestRunner {
    private List<UnitTest> tests = new ArrayList<>();
    private int passedTests = 0;
    private int failedTests = 0;

    /**
     * Adds a test to the test suite
     * @param test The UnitTest to add
     */
    public void addTest(UnitTest test) {
        tests.add(test);
    }

    /**
     * Runs all tests in the test suite and returns detailed results
     * @return TestResult object containing summary and detailed results
     */
    public TestResult runAllTests() {
        passedTests = 0;
        failedTests = 0;
        List<String> details = new ArrayList<>();

        for (UnitTest test : tests) {
            try {
                test.runTest();
                passedTests++;
                details.add("PASSED: " + test.getTestName());
            } catch (AssertionError e) {
                failedTests++;
                details.add("FAILED: " + test.getTestName() + " - " + e.getMessage());
            } catch (Exception e) {
                failedTests++;
                details.add("ERROR: " + test.getTestName() + " - " + e.getMessage());
            }
        }

        return new TestResult(passedTests, failedTests, details);
    }

    /**
     * Inner class to hold test results
     */
    public static class TestResult {
        public final int passed;
        public final int failed;
        public final List<String> details;

        public TestResult(int passed, int failed, List<String> details) {
            this.passed = passed;
            this.failed = failed;
            this.details = details;
        }
    }
}

/**
 * Interface for unit tests
 */
interface UnitTest {
    void runTest() throws Exception;
    String getTestName();
}
