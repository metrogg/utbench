// Converted Java method
import java.awt.*;
import java.util.List;
import javax.swing.*;

/**
 * EnhancedLayoutManager provides advanced layout capabilities beyond standard Swing layouts.
 * It supports dynamic grid-based layouts with flexible constraints and responsive behavior.
 */
class EnhancedLayoutManager {
    private List<Component> components;
    private List<LayoutConstraint> constraints;
    private boolean respectAspectRatio;
    private double aspectRatio;

    /**
     * Constructs an EnhancedLayoutManager with specified components and constraints.
     *
     * @param components List of components to layout
     * @param constraints List of layout constraints for each component
     * @param respectAspectRatio Whether to maintain aspect ratio of components
     * @param aspectRatio The target aspect ratio (width/height)
     */
    public EnhancedLayoutManager(List<Component> components, List<LayoutConstraint> constraints,
                               boolean respectAspectRatio, double aspectRatio) {
        this.components = components;
        this.constraints = constraints;
        this.respectAspectRatio = respectAspectRatio;
        this.aspectRatio = aspectRatio;
    }

    /**
     * Calculates optimal component positions and sizes based on constraints.
     *
     * @param containerWidth Available width for layout
     * @param containerHeight Available height for layout
     * @return Map of components to their calculated bounds
     */
    public java.util.Map<Component, Rectangle> calculateLayout(int containerWidth, int containerHeight) {
        java.util.Map<Component, Rectangle> layout = new java.util.HashMap<>();

        if (components.isEmpty()) {
            return layout;
        }

        // Calculate total weight for proportional sizing
        double totalWeight = constraints.stream()
                .mapToDouble(LayoutConstraint::getWeight)
                .sum();

        // Calculate available space after accounting for fixed-size components
        int remainingWidth = containerWidth;
        int remainingHeight = containerHeight;
        
        for (int i = 0; i < components.size(); i++) {
            Component comp = components.get(i);
            LayoutConstraint constraint = constraints.get(i);

            if (constraint.isFixedSize()) {
                Rectangle bounds = new Rectangle(
                    constraint.getX(),
                    constraint.getY(),
                    constraint.getWidth(),
                    constraint.getHeight()
                );
                layout.put(comp, bounds);
                
                remainingWidth -= bounds.width;
                remainingHeight -= bounds.height;
            }
        }

        // Distribute remaining space proportionally
        for (int i = 0; i < components.size(); i++) {
            Component comp = components.get(i);
            LayoutConstraint constraint = constraints.get(i);

            if (!constraint.isFixedSize()) {
                int width = (int) (remainingWidth * (constraint.getWeight() / totalWeight));
                int height = (int) (remainingHeight * (constraint.getWeight() / totalWeight));
                
                if (respectAspectRatio) {
                    double currentRatio = (double) width / height;
                    if (currentRatio > aspectRatio) {
                        width = (int) (height * aspectRatio);
                    } else {
                        height = (int) (width / aspectRatio);
                    }
                }

                Rectangle bounds = new Rectangle(
                    constraint.getX(),
                    constraint.getY(),
                    width,
                    height
                );
                layout.put(comp, bounds);
            }
        }

        return layout;
    }

    /**
     * Inner class representing layout constraints for components.
     */
    public static class LayoutConstraint {
        private int x;
        private int y;
        private int width;
        private int height;
        private double weight;
        private boolean fixedSize;

        public LayoutConstraint(int x, int y, int width, int height, double weight, boolean fixedSize) {
            this.x = x;
            this.y = y;
            this.width = width;
            this.height = height;
            this.weight = weight;
            this.fixedSize = fixedSize;
        }

        // Getters and setters omitted for brevity
        public int getX() { return x; }
        public int getY() { return y; }
        public int getWidth() { return width; }
        public int getHeight() { return height; }
        public double getWeight() { return weight; }
        public boolean isFixedSize() { return fixedSize; }
    }
}
