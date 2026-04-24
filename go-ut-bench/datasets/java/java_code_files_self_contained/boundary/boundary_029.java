// Converted Java method
import java.util.Arrays;

class EggDropSolver {
    /**
     * Solves the generalized egg drop problem with multiple solutions:
     * 1. Basic DP solution (original)
     * 2. Mathematical solution
     * 3. Binary search solution
     * 
     * @param k Number of eggs
     * @param n Number of floors
     * @return Map containing all solutions and their computation times
     */
    public static EggDropResult solveEggDrop(int k, int n) {
        if (k <= 0 || n <= 0) {
            throw new IllegalArgumentException("Eggs and floors must be positive values");
        }

        long startTime, endTime;
        int dpSolution, mathSolution, binarySearchSolution;

        // DP Solution
        startTime = System.nanoTime();
        dpSolution = dpEggDrop(k, n);
        endTime = System.nanoTime();
        long dpTime = endTime - startTime;

        // Mathematical Solution
        startTime = System.nanoTime();
        mathSolution = mathEggDrop(k, n);
        endTime = System.nanoTime();
        long mathTime = endTime - startTime;

        // Binary Search Solution
        startTime = System.nanoTime();
        binarySearchSolution = binarySearchEggDrop(k, n);
        endTime = System.nanoTime();
        long binarySearchTime = endTime - startTime;

        return new EggDropResult(
            dpSolution, mathSolution, binarySearchSolution,
            dpTime, mathTime, binarySearchTime
        );
    }

    // Original DP solution
    private static int dpEggDrop(int k, int n) {
        int[] dp = new int[k + 1];
        int ans = 0;
        while (dp[k] < n) {
            for (int i = k; i > 0; i--) {
                dp[i] = dp[i] + dp[i - 1] + 1;
            }
            ans++;
        }
        return ans;
    }

    // Mathematical solution using binomial coefficients
    private static int mathEggDrop(int k, int n) {
        int low = 1, high = n;
        while (low < high) {
            int mid = low + (high - low) / 2;
            if (binomialCoefficientSum(mid, k, n) < n) {
                low = mid + 1;
            } else {
                high = mid;
            }
        }
        return low;
    }

    private static int binomialCoefficientSum(int x, int k, int n) {
        int sum = 0, term = 1;
        for (int i = 1; i <= k; i++) {
            term *= x - i + 1;
            term /= i;
            sum += term;
            if (sum > n) break;
        }
        return sum;
    }

    // Binary search solution
    private static int binarySearchEggDrop(int k, int n) {
        int[][] memo = new int[k + 1][n + 1];
        for (int[] row : memo) {
            Arrays.fill(row, -1);
        }
        return bsHelper(k, n, memo);
    }

    private static int bsHelper(int k, int n, int[][] memo) {
        if (n == 0 || n == 1) return n;
        if (k == 1) return n;
        if (memo[k][n] != -1) return memo[k][n];

        int low = 1, high = n, result = n;
        while (low <= high) {
            int mid = (low + high) / 2;
            int broken = bsHelper(k - 1, mid - 1, memo);
            int notBroken = bsHelper(k, n - mid, memo);
            if (broken < notBroken) {
                low = mid + 1;
                result = Math.min(result, notBroken + 1);
            } else {
                high = mid - 1;
                result = Math.min(result, broken + 1);
            }
        }
        memo[k][n] = result;
        return result;
    }

    static class EggDropResult {
        final int dpSolution;
        final int mathSolution;
        final int binarySearchSolution;
        final long dpTime;
        final long mathTime;
        final long binarySearchTime;

        EggDropResult(int dpSolution, int mathSolution, int binarySearchSolution,
                     long dpTime, long mathTime, long binarySearchTime) {
            this.dpSolution = dpSolution;
            this.mathSolution = mathSolution;
            this.binarySearchSolution = binarySearchSolution;
            this.dpTime = dpTime;
            this.mathTime = mathTime;
            this.binarySearchTime = binarySearchTime;
        }

        @Override
        public String toString() {
            return String.format(
                "DP Solution: %d (%,d ns)%nMath Solution: %d (%,d ns)%nBinary Search Solution: %d (%,d ns)",
                dpSolution, dpTime, mathSolution, mathTime, binarySearchSolution, binarySearchTime
            );
        }
    }
}
