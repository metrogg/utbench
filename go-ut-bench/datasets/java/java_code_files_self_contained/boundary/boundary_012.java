// Converted Java method
import java.util.List;

class PaginationUtils {
    
    /**
     * Paginates a list of items based on limit and offset parameters.
     * Handles edge cases like invalid inputs and returns appropriate results.
     * 
     * @param <T> The type of items in the list
     * @param items The complete list of items to paginate
     * @param limit Maximum number of items to return (must be positive)
     * @param offset Starting index (must be non-negative)
     * @return Sublist representing the paginated results
     * @throws IllegalArgumentException if limit or offset are invalid
     */
    public static <T> List<T> paginate(List<T> items, int limit, int offset) {
        if (limit <= 0) {
            throw new IllegalArgumentException("Limit must be positive");
        }
        if (offset < 0) {
            throw new IllegalArgumentException("Offset cannot be negative");
        }
        if (items == null || items.isEmpty()) {
            return List.of();
        }
        
        int fromIndex = offset;
        int toIndex = Math.min(offset + limit, items.size());
        
        if (fromIndex >= items.size()) {
            return List.of();
        }
        
        return items.subList(fromIndex, toIndex);
    }
    
    /**
     * Calculates the total number of pages for a given list size and items per page.
     * 
     * @param totalItems Total number of items
     * @param itemsPerPage Number of items per page
     * @return Total number of pages needed
     */
    public static int calculateTotalPages(int totalItems, int itemsPerPage) {
        if (itemsPerPage <= 0) {
            throw new IllegalArgumentException("Items per page must be positive");
        }
        return (int) Math.ceil((double) totalItems / itemsPerPage);
    }
}
