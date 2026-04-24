import java.util.HashMap;
import java.util.Map;

class WebFragmentSimulator {
    private String currentUrl;
    private boolean isFragmentVisible;
    private Map<String, String> intentExtras;
    private boolean networkAvailable;

    public WebFragmentSimulator() {
        this.intentExtras = new HashMap<>();
        this.isFragmentVisible = false;
        this.networkAvailable = true;
    }

    /**
     * Simulates loading a URL in a web view with various conditions
     * @param url The URL to load
     * @param isFragmentVisible Whether the fragment is currently visible to user
     * @param networkAvailable Whether network connection is available
     * @return Map containing load status and any additional information
     */
    public Map<String, Object> loadWebContent(String url, boolean isFragmentVisible, boolean networkAvailable) {
        Map<String, Object> result = new HashMap<>();
        this.currentUrl = url;
        this.isFragmentVisible = isFragmentVisible;
        this.networkAvailable = networkAvailable;

        if (!isFragmentVisible) {
            result.put("status", "pending");
            result.put("message", "Fragment not visible - loading deferred");
            return result;
        }

        if (!networkAvailable) {
            result.put("status", "error");
            result.put("message", "Network unavailable");
            return result;
        }

        if (url == null || url.isEmpty()) {
            result.put("status", "error");
            result.put("message", "Invalid URL");
            return result;
        }

        // Simulate successful load
        result.put("status", "success");
        result.put("url", url);
        result.put("loadedTime", System.currentTimeMillis());
        
        // Check if URL needs parameters appended
        if (url.contains("ziroomupin")) {
            String modifiedUrl = appendUrlParameters(url);
            result.put("modifiedUrl", modifiedUrl);
        }

        return result;
    }

    private String appendUrlParameters(String originalUrl) {
        if (originalUrl.contains("?")) {
            return originalUrl + "&app_version=1.0&os=android";
        } else {
            return originalUrl + "?app_version=1.0&os=android";
        }
    }

    /**
     * Simulates starting a new activity with intent extras
     * @param targetActivity The target activity type
     * @param title The title for the new activity
     * @param url The URL to pass
     * @param isExternal Whether this is an external page
     */
    public void startNewActivity(int targetActivity, String title, String url, boolean isExternal) {
        intentExtras.clear();
        intentExtras.put("title", title);
        
        if (isExternal) {
            intentExtras.put("url", appendUrlParameters(url));
        } else {
            intentExtras.put("url", url);
        }
        
        intentExtras.put("activityType", String.valueOf(targetActivity));
    }

    public Map<String, String> getIntentExtras() {
        return new HashMap<>(intentExtras);
    }
}
