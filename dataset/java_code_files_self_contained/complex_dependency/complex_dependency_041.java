// Converted Java method
import java.util.List;
import java.util.ArrayList;
import java.util.Collections;

class PaginationProcessor<T> {
    private List<T> fullDataset;
    private int defaultPageSize;
    
    /**
     * Creates a new PaginationProcessor with the given dataset and default page size.
     * 
     * @param dataset The complete list of items to paginate
     * @param defaultPageSize The number of items per page (default)
     */
    public PaginationProcessor(List<T> dataset, int defaultPageSize) {
        this.fullDataset = new ArrayList<>(dataset);
        this.defaultPageSize = defaultPageSize;
    }
    
    /**
     * Gets a paginated subset of the data with metadata about the pagination.
     * 
     * @param page The requested page number (1-based index)
     * @param itemsPerPage Number of items per page (if null, uses default)
     * @return PaginationResult containing the page data and metadata
     * @throws IllegalArgumentException if page or itemsPerPage is invalid
     */
    public PaginationResult<T> getPaginatedData(int page, Integer itemsPerPage) {
        int effectiveItemsPerPage = itemsPerPage != null ? itemsPerPage : defaultPageSize;
        
        if (page < 1) {
            throw new IllegalArgumentException("Page number must be positive");
        }
        if (effectiveItemsPerPage < 1) {
            throw new IllegalArgumentException("Items per page must be positive");
        }
        
        int totalItems = fullDataset.size();
        int totalPages = (int) Math.ceil((double) totalItems / effectiveItemsPerPage);
        
        if (totalItems == 0) {
            return new PaginationResult<>(Collections.emptyList(), page, effectiveItemsPerPage, totalPages, totalItems);
        }
        
        if (page > totalPages) {
            page = totalPages;
        }
        
        int fromIndex = (page - 1) * effectiveItemsPerPage;
        int toIndex = Math.min(fromIndex + effectiveItemsPerPage, totalItems);
        
        List<T> pageData = fullDataset.subList(fromIndex, toIndex);
        
        return new PaginationResult<>(pageData, page, effectiveItemsPerPage, totalPages, totalItems);
    }
    
    /**
     * Inner class representing the pagination result with metadata.
     */
    public static class PaginationResult<T> {
        private final List<T> items;
        private final int currentPage;
        private final int itemsPerPage;
        private final int totalPages;
        private final int totalItems;
        
        public PaginationResult(List<T> items, int currentPage, int itemsPerPage, int totalPages, int totalItems) {
            this.items = new ArrayList<>(items);
            this.currentPage = currentPage;
            this.itemsPerPage = itemsPerPage;
            this.totalPages = totalPages;
            this.totalItems = totalItems;
        }
        
        // Getters
        public List<T> getItems() { return items; }
        public int getCurrentPage() { return currentPage; }
        public int getItemsPerPage() { return itemsPerPage; }
        public int getTotalPages() { return totalPages; }
        public int getTotalItems() { return totalItems; }
        
        @Override
        public String toString() {
            return String.format("Page %d of %d (items %d-%d of %d)",
                    currentPage, totalPages,
                    (currentPage - 1) * itemsPerPage + 1,
                    Math.min(currentPage * itemsPerPage, totalItems),
                    totalItems);
        }
    }
}
