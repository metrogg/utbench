// Converted Java method
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;

class BookmarkAnalyzer {
    
    /**
     * Analyzes a list of bookmarks and provides various statistics and filtering options.
     * 
     * @param bookmarks List of bookmark URLs
     * @param minLength Minimum length for a bookmark to be considered valid
     * @return A BookmarkAnalysisResult object containing various statistics about the bookmarks
     * @throws IllegalArgumentException if bookmarks list is null
     */
    public BookmarkAnalysisResult analyzeBookmarks(List<String> bookmarks, int minLength) {
        if (bookmarks == null) {
            throw new IllegalArgumentException("Bookmarks list cannot be null");
        }
        
        BookmarkAnalysisResult result = new BookmarkAnalysisResult();
        
        // Basic statistics
        result.totalCount = bookmarks.size();
        result.validCount = (int) bookmarks.stream()
                .filter(url -> url != null && url.length() >= minLength)
                .count();
        
        // Filter valid bookmarks
        List<String> validBookmarks = bookmarks.stream()
                .filter(url -> url != null && url.length() >= minLength)
                .collect(Collectors.toList());
        
        // Calculate average length
        if (!validBookmarks.isEmpty()) {
            result.averageLength = validBookmarks.stream()
                    .mapToInt(String::length)
                    .average()
                    .getAsDouble();
        }
        
        // Find longest and shortest
        if (!validBookmarks.isEmpty()) {
            result.longestBookmark = validBookmarks.stream()
                    .max((a, b) -> Integer.compare(a.length(), b.length()))
                    .get();
            result.shortestBookmark = validBookmarks.stream()
                    .min((a, b) -> Integer.compare(a.length(), b.length()))
                    .get();
        }
        
        // Count secure (https) bookmarks
        result.secureCount = (int) validBookmarks.stream()
                .filter(url -> url.startsWith("https://"))
                .count();
        
        return result;
    }
    
    public static class BookmarkAnalysisResult {
        public int totalCount;
        public int validCount;
        public int secureCount;
        public double averageLength;
        public String longestBookmark;
        public String shortestBookmark;
        
        @Override
        public String toString() {
            return String.format(
                "BookmarkAnalysisResult{totalCount=%d, validCount=%d, secureCount=%d, " +
                "averageLength=%.2f, longestBookmark='%s', shortestBookmark='%s'}",
                totalCount, validCount, secureCount, averageLength, 
                longestBookmark, shortestBookmark
            );
        }
    }
}
