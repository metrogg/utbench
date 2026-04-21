// Converted Java method
import java.util.ArrayList;
import java.util.List;

class ItemCatService {
    private List<Category> categories;

    public ItemCatService() {
        this.categories = new ArrayList<>();
    }

    /**
     * Adds a new category to the service
     * @param id The category ID
     * @param name The category name
     * @param parentId The parent category ID (0 for root categories)
     */
    public void addCategory(long id, String name, long parentId) {
        categories.add(new Category(id, name, parentId));
    }

    /**
     * Gets the complete category tree structure
     * @return CatResult containing the hierarchical category structure
     */
    public CatResult getItemCatList() {
        CatResult result = new CatResult();
        List<CategoryNode> rootNodes = new ArrayList<>();
        
        // Find all root categories (parentId = 0)
        for (Category category : categories) {
            if (category.getParentId() == 0) {
                rootNodes.add(buildCategoryTree(category));
            }
        }
        
        result.setData(rootNodes);
        return result;
    }

    /**
     * Recursively builds the category tree structure
     * @param current The current category node
     * @return CategoryNode with all its children
     */
    private CategoryNode buildCategoryTree(Category current) {
        CategoryNode node = new CategoryNode(current.getId(), current.getName());
        List<CategoryNode> children = new ArrayList<>();
        
        // Find all children of the current category
        for (Category category : categories) {
            if (category.getParentId() == current.getId()) {
                children.add(buildCategoryTree(category));
            }
        }
        
        if (!children.isEmpty()) {
            node.setChildren(children);
        }
        return node;
    }

    /**
     * Inner class representing a category
     */
    private static class Category {
        private final long id;
        private final String name;
        private final long parentId;

        public Category(long id, String name, long parentId) {
            this.id = id;
            this.name = name;
            this.parentId = parentId;
        }

        public long getId() { return id; }
        public String getName() { return name; }
        public long getParentId() { return parentId; }
    }
}

class CatResult {
    private List<CategoryNode> data;

    public List<CategoryNode> getData() { return data; }
    public void setData(List<CategoryNode> data) { this.data = data; }
}

class CategoryNode {
    private long id;
    private String name;
    private List<CategoryNode> children;

    public CategoryNode(long id, String name) {
        this.id = id;
        this.name = name;
    }

    public long getId() { return id; }
    public String getName() { return name; }
    public List<CategoryNode> getChildren() { return children; }
    public void setChildren(List<CategoryNode> children) { this.children = children; }
}
