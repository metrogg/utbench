// Converted Java method
import java.util.HashMap;
import java.util.Map;

class FragmentManager {
    private Map<String, Fragment> fragmentMap = new HashMap<>();
    private Fragment currentFragment;
    private String currentFragmentTag;

    /**
     * Adds a fragment to the manager with a specific tag
     * @param tag The identifier for the fragment
     * @param fragment The fragment instance
     */
    public void addFragment(String tag, Fragment fragment) {
        if (tag == null || fragment == null) {
            throw new IllegalArgumentException("Tag and fragment cannot be null");
        }
        fragmentMap.put(tag, fragment);
    }

    /**
     * Switches to the fragment with the specified tag
     * @param tag The identifier of the fragment to switch to
     * @return true if the switch was successful, false otherwise
     */
    public boolean switchFragment(String tag) {
        if (!fragmentMap.containsKey(tag)) {
            return false;
        }

        Fragment newFragment = fragmentMap.get(tag);
        if (currentFragment == newFragment) {
            return true;
        }

        // In a real Android environment, this would handle fragment transactions
        currentFragment = newFragment;
        currentFragmentTag = tag;
        return true;
    }

    /**
     * Gets the currently displayed fragment tag
     * @return The current fragment tag
     */
    public String getCurrentFragmentTag() {
        return currentFragmentTag;
    }

    /**
     * Gets the currently displayed fragment
     * @return The current fragment instance
     */
    public Fragment getCurrentFragment() {
        return currentFragment;
    }

    /**
     * Simple Fragment representation for testing purposes
     */
    static class Fragment {
        private final String name;

        public Fragment(String name) {
            this.name = name;
        }

        public String getName() {
            return name;
        }
    }
}
