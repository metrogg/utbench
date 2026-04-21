import java.util.Arrays;
import java.util.List;
import java.util.ArrayList;

class AdditionChainFinder {
    private int[] shortestChain;
    private int shortestLength;
    private int target;

    /**
     * Finds the shortest addition chain for a given target number.
     * An addition chain is a sequence where each number is the sum of two previous numbers.
     *
     * @param target The target number to reach in the addition chain
     * @return List of integers representing the shortest addition chain
     * @throws IllegalArgumentException if target is less than 1
     */
    public List<Integer> findShortestAdditionChain(int target) {
        if (target < 1) {
            throw new IllegalArgumentException("Target must be at least 1");
        }

        this.target = target;
        this.shortestLength = Integer.MAX_VALUE;
        this.shortestChain = new int[32]; // Sufficient size for practical numbers

        int[] currentChain = new int[32];
        currentChain[0] = 1; // Addition chain always starts with 1

        depthFirstSearch(currentChain, 1);

        // Convert array to list and trim to actual length
        List<Integer> result = new ArrayList<>();
        for (int i = 0; i < shortestLength; i++) {
            result.add(shortestChain[i]);
        }
        return result;
    }

    private void depthFirstSearch(int[] currentChain, int currentLength) {
        if (currentLength >= shortestLength) {
            return; // No need to proceed if we already have a shorter chain
        }

        int lastNumber = currentChain[currentLength - 1];
        if (lastNumber == target) {
            // Found a new shortest chain
            shortestLength = currentLength;
            System.arraycopy(currentChain, 0, shortestChain, 0, currentLength);
            return;
        }

        // Try all possible pairs of previous numbers to generate the next number
        for (int i = currentLength - 1; i >= 0; i--) {
            for (int j = i; j >= 0; j--) {
                int nextNumber = currentChain[i] + currentChain[j];
                if (nextNumber > lastNumber && nextNumber <= target) {
                    currentChain[currentLength] = nextNumber;
                    depthFirstSearch(currentChain, currentLength + 1);
                }
            }
        }
    }

    /**
     * Validates if a given sequence is a valid addition chain for the target number.
     *
     * @param chain The sequence to validate
     * @param target The target number to reach
     * @return true if valid addition chain, false otherwise
     */
    public static boolean validateAdditionChain(List<Integer> chain, int target) {
        if (chain.isEmpty() || chain.get(0) != 1 || chain.get(chain.size() - 1) != target) {
            return false;
        }

        for (int i = 1; i < chain.size(); i++) {
            boolean valid = false;
            int current = chain.get(i);
            
            // Check all possible pairs of previous numbers
            for (int j = 0; j < i; j++) {
                for (int k = j; k < i; k++) {
                    if (chain.get(j) + chain.get(k) == current) {
                        valid = true;
                        break;
                    }
                }
                if (valid) break;
            }
            
            if (!valid) return false;
        }
        return true;
    }
}
