import java.util.Arrays;
import java.util.List;

class SequenceAnalyzer {
    
    /**
     * Finds the length of the longest consecutive sequence that can be formed
     * using the given numbers, where zeros can be used as wildcards (jokers).
     * 
     * @param numbers List of integers including potential jokers (zeros)
     * @return Length of the longest consecutive sequence that can be formed
     */
    public int findLongestConsecutiveSequence(List<Integer> numbers) {
        if (numbers.isEmpty()) {
            return 0;
        }
        
        // Convert to array and sort
        int[] data = new int[numbers.size()];
        for (int i = 0; i < numbers.size(); i++) {
            data[i] = numbers.get(i);
        }
        Arrays.sort(data);
        
        int jokerCount = 0;
        for (int num : data) {
            if (num == 0) {
                jokerCount++;
            }
        }
        
        if (jokerCount == data.length) {
            return data.length;
        }
        
        int maxLength = 0;
        
        for (int i = jokerCount; i < data.length; i++) {
            int currentJokers = jokerCount;
            int prev = data[i];
            int currentLength = 1;
            boolean firstGap = true;
            int k = i + 1;
            
            while (k < data.length) {
                int curr = data[k];
                if (prev == curr) {
                    k++;
                } else if (prev + 1 == curr) {
                    currentLength++;
                    prev = curr;
                    k++;
                } else {
                    if (firstGap) {
                        i = k - 1;
                        firstGap = false;
                    }
                    if (currentJokers == 0) break;
                    currentJokers--;
                    currentLength++;
                    prev++;
                }
            }
            
            maxLength = Math.max(maxLength, currentLength + currentJokers);
        }
        
        return maxLength;
    }
}
