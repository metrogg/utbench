// Converted Java method
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class CategoryManager {
    private Map<String, List<Note>> categories;
    private String currentCategory;

    public CategoryManager() {
        categories = new HashMap<>();
        currentCategory = "生活"; // Default category
    }

    /**
     * Adds a new category to the manager.
     * @param categoryName Name of the category to add
     * @return true if added successfully, false if category already exists
     */
    public boolean addCategory(String categoryName) {
        if (categories.containsKey(categoryName)) {
            return false;
        }
        categories.put(categoryName, new ArrayList<>());
        return true;
    }

    /**
     * Removes a category and all its notes.
     * @param categoryName Name of the category to remove
     * @return List of removed notes, or null if category didn't exist
     */
    public List<Note> removeCategory(String categoryName) {
        if (!categories.containsKey(categoryName)) {
            return null;
        }
        List<Note> removedNotes = categories.remove(categoryName);
        if (currentCategory.equals(categoryName)) {
            currentCategory = categories.isEmpty() ? "" : categories.keySet().iterator().next();
        }
        return removedNotes;
    }

    /**
     * Adds a note to the specified category.
     * @param note The note to add
     * @param categoryName Target category name
     * @return true if added successfully, false if category doesn't exist
     */
    public boolean addNoteToCategory(Note note, String categoryName) {
        if (!categories.containsKey(categoryName)) {
            return false;
        }
        categories.get(categoryName).add(note);
        return true;
    }

    /**
     * Gets all notes from the specified category.
     * @param categoryName Name of the category to query
     * @return List of notes in the category, or null if category doesn't exist
     */
    public List<Note> getNotesByCategory(String categoryName) {
        if (!categories.containsKey(categoryName)) {
            return null;
        }
        return new ArrayList<>(categories.get(categoryName));
    }

    /**
     * Gets the current active category.
     * @return Name of the current category
     */
    public String getCurrentCategory() {
        return currentCategory;
    }

    /**
     * Sets the current active category.
     * @param categoryName Name of the category to set as current
     * @return true if set successfully, false if category doesn't exist
     */
    public boolean setCurrentCategory(String categoryName) {
        if (!categories.containsKey(categoryName)) {
            return false;
        }
        currentCategory = categoryName;
        return true;
    }

    /**
     * Gets all available category names.
     * @return List of all category names
     */
    public List<String> getAllCategories() {
        return new ArrayList<>(categories.keySet());
    }

    /**
     * Note class representing a simple note with title and content.
     */
    public static class Note {
        private String title;
        private String content;
        private long id;

        public Note(String title, String content) {
            this.title = title;
            this.content = content;
            this.id = System.currentTimeMillis();
        }

        public String getTitle() {
            return title;
        }

        public String getContent() {
            return content;
        }

        public long getId() {
            return id;
        }

        @Override
        public String toString() {
            return "Note{" +
                    "title='" + title + '\'' +
                    ", content='" + content + '\'' +
                    ", id=" + id +
                    '}';
        }
    }
}
