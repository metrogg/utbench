// Converted Java method
import java.util.*;
import java.util.stream.Collectors;

class DataMatcher {
    
    /**
     * Enhanced data matching with multiple matching strategies and result analysis
     * 
     * @param dataset1 First dataset to match
     * @param dataset2 Second dataset to match
     * @param matchStrategy Strategy for matching (1: exact match, 2: contains, 3: fuzzy match)
     * @param keyIndices Array of indices to use for matching [dataset1 index, dataset2 index]
     * @return Map containing matched results, unmatched from dataset1, and unmatched from dataset2
     * @throws IllegalArgumentException if inputs are invalid
     */
    public Map<String, List<String[]>> enhancedMatch(List<String[]> dataset1, 
                                                   List<String[]> dataset2, 
                                                   int matchStrategy, 
                                                   int[] keyIndices) {
        
        // Validate inputs
        if (dataset1 == null || dataset2 == null) {
            throw new IllegalArgumentException("Input datasets cannot be null");
        }
        if (matchStrategy < 1 || matchStrategy > 3) {
            throw new IllegalArgumentException("Invalid match strategy");
        }
        if (keyIndices == null || keyIndices.length != 2) {
            throw new IllegalArgumentException("Key indices must be an array of two elements");
        }
        
        List<String[]> matched = new ArrayList<>();
        List<String[]> unmatched1 = new ArrayList<>(dataset1);
        List<String[]> unmatched2 = new ArrayList<>(dataset2);
        
        // Find maximum lengths for padding
        int maxLen1 = dataset1.stream().mapToInt(a -> a.length).max().orElse(0);
        int maxLen2 = dataset2.stream().mapToInt(a -> a.length).max().orElse(0);
        
        // Match records based on strategy
        Iterator<String[]> iter1 = unmatched1.iterator();
        while (iter1.hasNext()) {
            String[] record1 = iter1.next();
            if (record1.length <= keyIndices[0]) continue;
            
            String key1 = record1[keyIndices[0]].trim();
            
            Iterator<String[]> iter2 = unmatched2.iterator();
            while (iter2.hasNext()) {
                String[] record2 = iter2.next();
                if (record2.length <= keyIndices[1]) continue;
                
                String key2 = record2[keyIndices[1]].trim();
                boolean isMatch = false;
                
                switch (matchStrategy) {
                    case 1: // Exact match
                        isMatch = key1.equals(key2);
                        break;
                    case 2: // Contains
                        isMatch = key1.contains(key2) || key2.contains(key1);
                        break;
                    case 3: // Fuzzy match (simplified)
                        isMatch = key1.toLowerCase().contains(key2.toLowerCase()) || 
                                 key2.toLowerCase().contains(key1.toLowerCase());
                        break;
                }
                
                if (isMatch) {
                    // Combine matched records
                    String[] combined = new String[record1.length + record2.length];
                    System.arraycopy(record1, 0, combined, 0, record1.length);
                    System.arraycopy(record2, 0, combined, record1.length, record2.length);
                    matched.add(combined);
                    
                    iter1.remove();
                    iter2.remove();
                    break;
                }
            }
        }
        
        // Prepare results with padding for unmatched records
        Map<String, List<String[]>> results = new HashMap<>();
        results.put("matched", matched);
        
        // Pad unmatched records with empty strings
        results.put("unmatched1", unmatched1.stream()
            .map(record -> {
                String[] padded = new String[maxLen1];
                Arrays.fill(padded, "");
                System.arraycopy(record, 0, padded, 0, record.length);
                return padded;
            })
            .collect(Collectors.toList()));
        
        results.put("unmatched2", unmatched2.stream()
            .map(record -> {
                String[] padded = new String[maxLen2];
                Arrays.fill(padded, "");
                System.arraycopy(record, 0, padded, 0, record.length);
                return padded;
            })
            .collect(Collectors.toList()));
        
        return results;
    }
}
