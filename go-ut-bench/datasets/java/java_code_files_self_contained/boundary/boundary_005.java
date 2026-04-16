// Converted Java method
import java.util.HashMap;
import java.util.Map;

class PostCreator {
    /**
     * Validates and creates a post with the given details and image data.
     * 
     * @param title The post title
     * @param description The post description
     * @param price The item price
     * @param country The country location
     * @param stateProvince The state/province location
     * @param city The city location
     * @param contactEmail The contact email
     * @param imageData Either a file path (String) or byte array representing the image
     * @return A map containing validation results and post data if successful
     */
    public Map<String, Object> createPost(
            String title, String description, String price,
            String country, String stateProvince, String city,
            String contactEmail, Object imageData) {
        
        Map<String, Object> result = new HashMap<>();
        Map<String, String> errors = new HashMap<>();

        // Validate required fields
        if (title == null || title.trim().isEmpty()) {
            errors.put("title", "Title is required");
        } else if (title.length() > 100) {
            errors.put("title", "Title must be less than 100 characters");
        }

        if (description == null || description.trim().isEmpty()) {
            errors.put("description", "Description is required");
        }

        // Validate price format
        try {
            if (price != null && !price.trim().isEmpty()) {
                double priceValue = Double.parseDouble(price);
                if (priceValue < 0) {
                    errors.put("price", "Price cannot be negative");
                }
            }
        } catch (NumberFormatException e) {
            errors.put("price", "Invalid price format");
        }

        // Validate email format
        if (contactEmail != null && !contactEmail.trim().isEmpty() && 
            !contactEmail.matches("^[A-Za-z0-9+_.-]+@(.+)$")) {
            errors.put("contactEmail", "Invalid email format");
        }

        // Validate image data
        if (imageData == null) {
            errors.put("image", "Image is required");
        } else if (!(imageData instanceof String) && !(imageData instanceof byte[])) {
            errors.put("image", "Invalid image data format");
        }

        if (!errors.isEmpty()) {
            result.put("success", false);
            result.put("errors", errors);
            return result;
        }

        // Create post data if validation passes
        Map<String, String> postData = new HashMap<>();
        postData.put("title", title);
        postData.put("description", description);
        postData.put("price", price);
        postData.put("country", country);
        postData.put("stateProvince", stateProvince);
        postData.put("city", city);
        postData.put("contactEmail", contactEmail);

        result.put("success", true);
        result.put("postData", postData);
        result.put("imageType", imageData instanceof String ? "path" : "bytes");

        return result;
    }
}
