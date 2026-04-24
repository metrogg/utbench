// Converted Java method
import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

class PaginationService<T> {
    
    /**
     * Paginates and sorts a list of items with optional filtering
     * 
     * @param items The full list of items to paginate
     * @param page The page number (1-based)
     * @param rows Number of items per page
     * @param sortBy Property name to sort by (null for no sorting)
     * @param desc Whether to sort in descending order
     * @param key Optional filter key to search in items
     * @return PaginatedResult containing the page data and total count
     */
    public PaginatedResult<T> paginate(List<T> items, int page, int rows, String sortBy, boolean desc, String key) {
        // Filter items if key is provided
        List<T> filteredItems = filterItems(items, key);
        
        // Sort items if sortBy is provided
        if (sortBy != null) {
            filteredItems = sortItems(filteredItems, sortBy, desc);
        }
        
        // Calculate pagination boundaries
        int total = filteredItems.size();
        int fromIndex = (page - 1) * rows;
        if (fromIndex >= total) {
            return new PaginatedResult<>(new ArrayList<>(), total);
        }
        
        int toIndex = Math.min(fromIndex + rows, total);
        List<T> pageItems = filteredItems.subList(fromIndex, toIndex);
        
        return new PaginatedResult<>(pageItems, total);
    }
    
    private List<T> filterItems(List<T> items, String key) {
        if (key == null || key.isEmpty()) {
            return new ArrayList<>(items);
        }
        
        return items.stream()
            .filter(item -> item.toString().toLowerCase().contains(key.toLowerCase()))
            .collect(Collectors.toList());
    }
    
    @SuppressWarnings("unchecked")
    private List<T> sortItems(List<T> items, String sortBy, boolean desc) {
        Comparator<T> comparator = (a, b) -> {
            try {
                Object valA = a.getClass().getField(sortBy).get(a);
                Object valB = b.getClass().getField(sortBy).get(b);
                return ((Comparable<Object>)valA).compareTo(valB);
            } catch (Exception e) {
                return 0;
            }
        };
        
        if (desc) {
            comparator = comparator.reversed();
        }
        
        return items.stream()
            .sorted(comparator)
            .collect(Collectors.toList());
    }
    
    public static class PaginatedResult<T> {
        private final List<T> items;
        private final int total;
        
        public PaginatedResult(List<T> items, int total) {
            this.items = items;
            this.total = total;
        }
        
        public List<T> getItems() {
            return items;
        }
        
        public int getTotal() {
            return total;
        }
    }
}
