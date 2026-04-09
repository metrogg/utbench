// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class EmojiManager {
    private List<EmojiCategory> categories;
    private EmojiCategory currentCategory;
    private int currentCategoryIndex;
    private Map<String, List<Emoji>> customEmojiResources;

    public EmojiManager(boolean enableRecent, boolean enableCustom, boolean enableSystem) {
        this.categories = new ArrayList<>();
        this.customEmojiResources = new HashMap<>();
        
        if (enableRecent) {
            addRecentEmojiCategory();
        }
        if (enableCustom) {
            addCustomEmojiCategory();
        }
        if (enableSystem) {
            addSystemEmojiCategory();
        }
        
        setCurrentCategory(0);
    }

    public void addRecentEmojiCategory() {
        EmojiCategory recent = new EmojiCategory("Recent", new ArrayList<>());
        categories.add(recent);
    }

    public void addCustomEmojiCategory() {
        EmojiCategory custom = new EmojiCategory("Custom", new ArrayList<>());
        categories.add(custom);
    }

    public void addSystemEmojiCategory() {
        EmojiCategory system = new EmojiCategory("System", generateDefaultSystemEmojis());
        categories.add(system);
    }

    private List<Emoji> generateDefaultSystemEmojis() {
        List<Emoji> systemEmojis = new ArrayList<>();
        // Add some common emojis
        systemEmojis.add(new Emoji("😀", "Grinning Face"));
        systemEmojis.add(new Emoji("😂", "Face with Tears of Joy"));
        systemEmojis.add(new Emoji("👍", "Thumbs Up"));
        systemEmojis.add(new Emoji("❤️", "Red Heart"));
        systemEmojis.add(new Emoji("😊", "Smiling Face with Smiling Eyes"));
        return systemEmojis;
    }

    public void addCustomEmojiResource(String resourceName, List<Emoji> emojis) {
        customEmojiResources.put(resourceName, emojis);
        // Update custom category if it exists
        for (EmojiCategory category : categories) {
            if (category.getName().equals("Custom")) {
                category.setEmojis(new ArrayList<>(emojis));
                break;
            }
        }
    }

    public int getCategoryCount() {
        return categories.size();
    }

    public int getCurrentCategorySize() {
        return currentCategory != null ? currentCategory.getEmojiCount() : 0;
    }

    public int getTotalEmojiCount() {
        int total = 0;
        for (EmojiCategory category : categories) {
            total += category.getEmojiCount();
        }
        return total;
    }

    public void setCurrentCategory(int index) {
        if (index >= 0 && index < categories.size()) {
            currentCategoryIndex = index;
            currentCategory = categories.get(index);
        }
    }

    public EmojiCategory getCategory(int index) {
        if (index >= 0 && index < categories.size()) {
            return categories.get(index);
        }
        return null;
    }

    public List<Emoji> getEmojisInRange(int startIndex, int count) {
        List<Emoji> result = new ArrayList<>();
        if (currentCategory == null) return result;
        
        List<Emoji> emojis = currentCategory.getEmojis();
        int endIndex = Math.min(startIndex + count, emojis.size());
        
        for (int i = startIndex; i < endIndex; i++) {
            result.add(emojis.get(i));
        }
        return result;
    }

    public void addToRecent(Emoji emoji) {
        for (EmojiCategory category : categories) {
            if (category.getName().equals("Recent")) {
                // Prevent duplicates
                if (!category.getEmojis().contains(emoji)) {
                    category.getEmojis().add(0, emoji);
                    // Limit recent emojis to 20
                    if (category.getEmojis().size() > 20) {
                        category.getEmojis().remove(category.getEmojis().size() - 1);
                    }
                }
                break;
            }
        }
    }

    public static class EmojiCategory {
        private String name;
        private List<Emoji> emojis;

        public EmojiCategory(String name, List<Emoji> emojis) {
            this.name = name;
            this.emojis = emojis;
        }

        public String getName() {
            return name;
        }

        public List<Emoji> getEmojis() {
            return emojis;
        }

        public void setEmojis(List<Emoji> emojis) {
            this.emojis = emojis;
        }

        public int getEmojiCount() {
            return emojis.size();
        }
    }

    public static class Emoji {
        private String character;
        private String description;

        public Emoji(String character, String description) {
            this.character = character;
            this.description = description;
        }

        public String getCharacter() {
            return character;
        }

        public String getDescription() {
            return description;
        }

        @Override
        public boolean equals(Object obj) {
            if (this == obj) return true;
            if (obj == null || getClass() != obj.getClass()) return false;
            Emoji emoji = (Emoji) obj;
            return character.equals(emoji.character);
        }

        @Override
        public int hashCode() {
            return character.hashCode();
        }
    }
}
