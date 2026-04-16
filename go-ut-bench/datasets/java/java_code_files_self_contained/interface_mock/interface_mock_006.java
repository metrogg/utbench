// Converted Java method
import java.util.List;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;

class EnhancedProductManager {
    private Map<Integer, Product> productDatabase;
    private LoggerService loggerService;
    private List<Integer> restrictedCategories;

    /**
     * Initializes the EnhancedProductManager with required dependencies
     * @param loggerService Service for logging operations
     * @param restrictedCategories List of category IDs that are restricted
     */
    public EnhancedProductManager(LoggerService loggerService, List<Integer> restrictedCategories) {
        this.productDatabase = new HashMap<>();
        this.loggerService = loggerService;
        this.restrictedCategories = new ArrayList<>(restrictedCategories);
    }

    /**
     * Adds a product to the system after validation
     * @param product Product to be added
     * @return Map containing operation status and messages
     */
    public Map<String, String> addProduct(Product product) {
        Map<String, String> result = new HashMap<>();
        
        // Validate product
        if (product == null) {
            result.put("status", "error");
            result.put("message", "Product cannot be null");
            return result;
        }

        if (product.getName() == null || product.getName().isEmpty()) {
            result.put("status", "error");
            result.put("message", "Product name cannot be empty");
            return result;
        }

        if (restrictedCategories.contains(product.getCategoryId())) {
            result.put("status", "error");
            result.put("message", "Products in category " + product.getCategoryId() + " are not allowed");
            return result;
        }

        if (productDatabase.containsKey(product.getId())) {
            result.put("status", "error");
            result.put("message", "Product with ID " + product.getId() + " already exists");
            return result;
        }

        // Add product
        productDatabase.put(product.getId(), product);
        loggerService.logToSystem("Product added: " + product.getName());
        
        result.put("status", "success");
        result.put("message", "Product added successfully");
        return result;
    }

    /**
     * Retrieves all products in the system
     * @return List of all products
     */
    public List<Product> getAllProducts() {
        return new ArrayList<>(productDatabase.values());
    }

    /**
     * Finds a product by ID
     * @param id Product ID to search for
     * @return The product if found, null otherwise
     */
    public Product getProductById(int id) {
        return productDatabase.get(id);
    }

    /**
     * Removes a product from the system
     * @param id ID of product to remove
     * @return Map containing operation status and messages
     */
    public Map<String, String> removeProduct(int id) {
        Map<String, String> result = new HashMap<>();
        
        if (!productDatabase.containsKey(id)) {
            result.put("status", "error");
            result.put("message", "Product with ID " + id + " not found");
            return result;
        }

        Product removed = productDatabase.remove(id);
        loggerService.logToSystem("Product removed: " + removed.getName());
        
        result.put("status", "success");
        result.put("message", "Product removed successfully");
        return result;
    }
}

class Product {
    private int id;
    private String name;
    private int categoryId;

    public Product(int id, String name, int categoryId) {
        this.id = id;
        this.name = name;
        this.categoryId = categoryId;
    }

    // Getters and setters
    public int getId() { return id; }
    public String getName() { return name; }
    public int getCategoryId() { return categoryId; }
    public void setId(int id) { this.id = id; }
    public void setName(String name) { this.name = name; }
    public void setCategoryId(int categoryId) { this.categoryId = categoryId; }
}

interface LoggerService {
    void logToSystem(String message);
}
