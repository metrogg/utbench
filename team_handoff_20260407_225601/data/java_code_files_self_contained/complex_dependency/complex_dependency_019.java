import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

/**
 * Simulates a game character with customizable properties and skins.
 * This is a more generalized version of the original Minecraft boss code.
 */
class GameCharacter {
    private final UUID id;
    private final String name;
    private boolean isInvulnerable;
    private Map<String, String> properties;
    private Map<String, String> textures;

    /**
     * Creates a new GameCharacter with the given name.
     * @param name The character's name
     */
    public GameCharacter(String name) {
        this.id = UUID.randomUUID();
        this.name = name;
        this.isInvulnerable = true; // Default to invulnerable
        this.properties = new HashMap<>();
        this.textures = new HashMap<>();
    }

    /**
     * Simulates setting a skin for the character using a mock skin service.
     * @param skinId The ID of the skin to apply
     * @return true if the skin was applied successfully, false otherwise
     */
    public boolean setSkin(String skinId) {
        // Simulate a skin service response
        Map<String, String> mockResponse = mockSkinService(skinId);
        
        if (mockResponse != null) {
            String textureValue = mockResponse.get("value");
            String signature = mockResponse.get("signature");
            
            if (textureValue != null && signature != null) {
                textures.put("textures", textureValue);
                properties.put("signature", signature);
                return true;
            }
        }
        return false;
    }

    /**
     * Mock skin service that simulates the behavior of a real skin service.
     * @param skinId The skin ID to request
     * @return A map containing texture data or null if the request fails
     */
    private Map<String, String> mockSkinService(String skinId) {
        // Simulate different responses based on input
        if (skinId == null || skinId.isEmpty()) {
            return null;
        }

        Map<String, String> response = new HashMap<>();
        
        // For valid skin IDs, generate mock texture data
        if (skinId.equals("valid_skin")) {
            response.put("value", "texture_data_for_" + skinId);
            response.put("signature", "signature_for_" + skinId);
            return response;
        } 
        else if (skinId.equals("invalid_skin")) {
            return null; // Simulate failed request
        }
        else {
            // Default case for other skin IDs
            response.put("value", "default_texture_data");
            response.put("signature", "default_signature");
            return response;
        }
    }

    // Getters and setters for character properties
    public UUID getId() { return id; }
    public String getName() { return name; }
    public boolean isInvulnerable() { return isInvulnerable; }
    public void setInvulnerable(boolean invulnerable) { isInvulnerable = invulnerable; }
    public Map<String, String> getProperties() { return properties; }
    public Map<String, String> getTextures() { return textures; }
}
