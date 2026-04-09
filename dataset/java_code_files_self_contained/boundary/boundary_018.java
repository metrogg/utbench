import java.util.Arrays;

class PermutationGenerator {

    /**
     * Generates the next lexicographical permutation of an array of integers.
     * If the array is already in descending order, it will be sorted in ascending order.
     * 
     * @param nums The array of integers to permute (modified in place)
     * @return true if a next permutation exists and was generated, false if the array was reset to first permutation
     * @throws IllegalArgumentException if the input array is null
     */
    public boolean nextPermutation(int[] nums) {
        if (nums == null) {
            throw new IllegalArgumentException("Input array cannot be null");
        }
        
        if (nums.length <= 1) {
            return false;
        }

        // Find the first decreasing element from the end
        int i = nums.length - 2;
        while (i >= 0 && nums[i] >= nums[i + 1]) {
            i--;
        }

        boolean hasNextPermutation = true;
        
        if (i >= 0) {
            // Find the smallest number greater than nums[i] to the right of i
            int j = nums.length - 1;
            while (nums[j] <= nums[i]) {
                j--;
            }
            swap(nums, i, j);
        } else {
            hasNextPermutation = false;
        }

        // Reverse the suffix after i
        reverse(nums, i + 1, nums.length - 1);
        
        return hasNextPermutation;
    }

    /**
     * Swaps two elements in an array using XOR swap algorithm
     * 
     * @param nums The array containing elements to swap
     * @param i First index
     * @param j Second index
     */
    private void swap(int[] nums, int i, int j) {
        if (i != j) {
            nums[i] ^= nums[j];
            nums[j] ^= nums[i];
            nums[i] ^= nums[j];
        }
    }

    /**
     * Reverses a portion of an array between two indices
     * 
     * @param nums The array to reverse
     * @param start Starting index (inclusive)
     * @param end Ending index (inclusive)
     */
    private void reverse(int[] nums, int start, int end) {
        while (start < end) {
            swap(nums, start++, end--);
        }
    }

    /**
     * Generates all permutations of an array in lexicographical order
     * 
     * @param nums The initial array (will be modified)
     * @return List of all permutations in order
     */
    public int[][] generateAllPermutations(int[] nums) {
        if (nums == null || nums.length == 0) {
            return new int[0][];
        }

        // Sort the array to start with the first permutation
        Arrays.sort(nums);
        
        // Calculate factorial to determine array size
        int factorial = 1;
        for (int i = 2; i <= nums.length; i++) {
            factorial *= i;
        }

        int[][] result = new int[factorial][];
        result[0] = Arrays.copyOf(nums, nums.length);
        
        for (int i = 1; i < factorial; i++) {
            if (!nextPermutation(nums)) {
                break;
            }
            result[i] = Arrays.copyOf(nums, nums.length);
        }

        return result;
    }
}
