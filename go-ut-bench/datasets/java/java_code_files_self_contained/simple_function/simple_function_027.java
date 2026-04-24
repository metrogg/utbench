// Converted Java method
import java.util.ArrayList;
import java.util.Collections;

class CargoOptimizer {
    
    /**
     * Calculates the maximum number of trips possible given weights and a capacity.
     * The algorithm:
     * 1. First counts all elements that can make a trip alone (weight >= capacity)
     * 2. Then combines remaining elements in optimal groups
     * 3. Each group must satisfy: (max_weight * group_size) >= capacity
     * 
     * @param weights List of weights to be transported
     * @param capacity The capacity constraint for each trip
     * @return Maximum number of possible trips
     */
    public static int calculateMaxTrips(ArrayList<Integer> weights, int capacity) {
        if (weights == null || weights.isEmpty()) {
            return 0;
        }
        
        // Make a copy to avoid modifying original list
        ArrayList<Integer> remainingWeights = new ArrayList<>(weights);
        int tripCount = 0;
        
        // First pass: count items that can go alone
        for (int i = 0; i < remainingWeights.size(); i++) {
            if (remainingWeights.get(i) >= capacity) {
                tripCount++;
                remainingWeights.remove(i);
                i--; // Adjust index after removal
            }
        }
        
        // Sort remaining weights in descending order
        remainingWeights.sort(Collections.reverseOrder());
        
        // Process remaining items
        while (!remainingWeights.isEmpty()) {
            int maxWeight = remainingWeights.get(0);
            int groupSize = 1;
            
            // Calculate minimum group size needed
            while (maxWeight * groupSize < capacity && 
                   groupSize < remainingWeights.size()) {
                groupSize++;
            }
            
            if (maxWeight * groupSize >= capacity) {
                tripCount++;
                // Remove the used items (always taking largest remaining)
                for (int i = 0; i < groupSize && !remainingWeights.isEmpty(); i++) {
                    remainingWeights.remove(0);
                }
            } else {
                break; // No more valid groups possible
            }
        }
        
        return tripCount;
    }
}
