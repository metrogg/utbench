import java.util.ArrayList;
import java.util.List;

class PaginationManager {
    /**
     * Calculates pagination details for a given dataset.
     * 
     * @param totalItems Total number of items in the dataset
     * @param pageSize Number of items per page
     * @param currentPage The current page number (1-based index)
     * @return PaginationResult containing all pagination details
     * @throws IllegalArgumentException if any parameter is invalid
     */
    public PaginationResult calculatePagination(int totalItems, int pageSize, int currentPage) {
        if (totalItems < 0) {
            throw new IllegalArgumentException("Total items cannot be negative");
        }
        if (pageSize <= 0) {
            throw new IllegalArgumentException("Page size must be positive");
        }
        if (currentPage <= 0) {
            throw new IllegalArgumentException("Current page must be positive");
        }

        int totalPages = (int) Math.ceil((double) totalItems / pageSize);
        currentPage = Math.min(currentPage, totalPages);
        
        int startItem = (currentPage - 1) * pageSize + 1;
        int endItem = Math.min(currentPage * pageSize, totalItems);
        
        List<Integer> visiblePages = calculateVisiblePages(currentPage, totalPages);
        
        return new PaginationResult(
            totalItems,
            pageSize,
            currentPage,
            totalPages,
            startItem,
            endItem,
            visiblePages
        );
    }
    
    private List<Integer> calculateVisiblePages(int currentPage, int totalPages) {
        List<Integer> pages = new ArrayList<>();
        int range = 2; // Number of pages to show before and after current page
        
        // Always add first page
        if (currentPage > range + 1 && totalPages > range * 2 + 1) {
            pages.add(1);
            if (currentPage > range + 2) {
                pages.add(-1); // -1 represents ellipsis
            }
        }
        
        // Add pages around current page
        int start = Math.max(1, currentPage - range);
        int end = Math.min(totalPages, currentPage + range);
        
        for (int i = start; i <= end; i++) {
            pages.add(i);
        }
        
        // Always add last page
        if (currentPage + range < totalPages && totalPages > range * 2 + 1) {
            if (currentPage + range + 1 < totalPages) {
                pages.add(-1); // -1 represents ellipsis
            }
            pages.add(totalPages);
        }
        
        return pages;
    }
    
    public static class PaginationResult {
        public final int totalItems;
        public final int pageSize;
        public final int currentPage;
        public final int totalPages;
        public final int startItem;
        public final int endItem;
        public final List<Integer> visiblePages;
        
        public PaginationResult(int totalItems, int pageSize, int currentPage, 
                              int totalPages, int startItem, int endItem, 
                              List<Integer> visiblePages) {
            this.totalItems = totalItems;
            this.pageSize = pageSize;
            this.currentPage = currentPage;
            this.totalPages = totalPages;
            this.startItem = startItem;
            this.endItem = endItem;
            this.visiblePages = visiblePages;
        }
    }
}
