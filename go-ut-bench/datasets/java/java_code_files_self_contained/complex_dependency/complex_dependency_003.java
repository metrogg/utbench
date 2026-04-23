// Converted Java method
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

/**
 * Enhanced OTG Transformer Factory with additional capabilities
 * - Supports dynamic registration of new transformers
 * - Provides transformer validation
 * - Includes transformer metadata
 */
class EnhancedOTGTransformerFactory {
    private final Map<String, TransformerDescriptor> transformers;
    private final Map<String, TransformerMetadata> transformerMetadata;

    public EnhancedOTGTransformerFactory() {
        transformers = new HashMap<>();
        transformerMetadata = new HashMap<>();
        initializeDefaultTransformers();
    }

    private void initializeDefaultTransformers() {
        registerTransformer("contactreportToXml", 
            new TransformerDescriptor("CONTACT_REPORT_FORMAT", TransformerType.TXT2XML),
            new TransformerMetadata("Contact Report", "1.0", "Converts contact reports to XML"));
        
        registerTransformer("opnoteToXml", 
            new TransformerDescriptor("OPNOTE_FORMAT", TransformerType.TXT2XML),
            new TransformerMetadata("Op Note", "1.1", "Converts operational notes to XML"));
        
        registerTransformer("fourwhiskyToXml", 
            new TransformerDescriptor("FOURWHISKY_FORMAT", TransformerType.TXT2XML),
            new TransformerMetadata("Four Whisky", "2.0", "Converts Four Whisky format to XML"));
        
        // Add reverse transformers
        registerTransformer("contactreportToTxt", 
            new TransformerDescriptor("CONTACT_REPORT_FORMAT", TransformerType.XML2TXT),
            new TransformerMetadata("Contact Report Reverse", "1.0", "Converts XML back to contact report format"));
        
        registerTransformer("opnoteToTxt", 
            new TransformerDescriptor("OPNOTE_FORMAT", TransformerType.XML2TXT),
            new TransformerMetadata("Op Note Reverse", "1.1", "Converts XML back to operational notes"));
        
        registerTransformer("fourwhiskyToTxt", 
            new TransformerDescriptor("FOURWHISKY_FORMAT", TransformerType.XML2TXT),
            new TransformerMetadata("Four Whisky Reverse", "2.0", "Converts XML back to Four Whisky format"));
    }

    public void registerTransformer(String format, TransformerDescriptor descriptor, 
                                  TransformerMetadata metadata) {
        if (format == null || format.trim().isEmpty()) {
            throw new IllegalArgumentException("Format cannot be null or empty");
        }
        if (descriptor == null) {
            throw new IllegalArgumentException("Descriptor cannot be null");
        }
        
        transformers.put(format, descriptor);
        if (metadata != null) {
            transformerMetadata.put(format, metadata);
        }
    }

    public TransformerDescriptor getDescriptor(String format) {
        if (!transformers.containsKey(format)) {
            throw new IllegalArgumentException("Unknown transformer format: " + format);
        }
        return transformers.get(format);
    }

    public TransformerMetadata getMetadata(String format) {
        return transformerMetadata.get(format);
    }

    public Collection<String> getSupportedTransformers() {
        return transformers.keySet();
    }

    public boolean isTransformerSupported(String format) {
        return transformers.containsKey(format);
    }

    public enum TransformerType {
        TXT2XML, XML2TXT
    }

    public static class TransformerDescriptor {
        private final String format;
        private final TransformerType type;

        public TransformerDescriptor(String format, TransformerType type) {
            this.format = format;
            this.type = type;
        }

        public String getFormat() {
            return format;
        }

        public TransformerType getType() {
            return type;
        }
    }

    public static class TransformerMetadata {
        private final String name;
        private final String version;
        private final String description;

        public TransformerMetadata(String name, String version, String description) {
            this.name = name;
            this.version = version;
            this.description = description;
        }

        public String getName() {
            return name;
        }

        public String getVersion() {
            return version;
        }

        public String getDescription() {
            return description;
        }
    }
}
