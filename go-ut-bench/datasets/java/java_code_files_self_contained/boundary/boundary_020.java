// Converted Java method
import java.io.IOException;
import java.nio.file.*;
import java.util.*;
import java.util.stream.Collectors;

class AdvancedFileProcessor {
    private final Path directoryPath;
    
    /**
     * Initializes the file processor with a directory path.
     * @param directoryPath The path to the directory to process
     * @throws IOException If the directory cannot be accessed
     */
    public AdvancedFileProcessor(String directoryPath) throws IOException {
        this.directoryPath = Paths.get(directoryPath);
        if (!Files.isDirectory(this.directoryPath)) {
            throw new IllegalArgumentException("Path must be a directory");
        }
    }
    
    /**
     * Processes all files in the directory and returns statistics.
     * @return Map containing file count, total size, largest file, and word frequency
     * @throws IOException If files cannot be read
     */
    public Map<String, Object> processFiles() throws IOException {
        Map<String, Object> stats = new HashMap<>();
        List<Path> files = new ArrayList<>();
        long totalSize = 0;
        Path largestFile = null;
        long maxSize = 0;
        Map<String, Integer> wordFrequency = new HashMap<>();
        
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(directoryPath)) {
            for (Path file : stream) {
                if (Files.isRegularFile(file)) {
                    files.add(file);
                    long size = Files.size(file);
                    totalSize += size;
                    
                    if (size > maxSize) {
                        maxSize = size;
                        largestFile = file;
                    }
                    
                    // Process file content for word frequency
                    processFileContent(file, wordFrequency);
                }
            }
        }
        
        stats.put("fileCount", files.size());
        stats.put("totalSize", totalSize);
        stats.put("largestFile", largestFile != null ? largestFile.getFileName().toString() : "N/A");
        stats.put("wordFrequency", sortByValue(wordFrequency));
        
        return stats;
    }
    
    private void processFileContent(Path file, Map<String, Integer> wordFrequency) throws IOException {
        try {
            String content = new String(Files.readAllBytes(file));
            String[] words = content.split("\\s+");
            
            for (String word : words) {
                if (!word.isEmpty()) {
                    String cleanedWord = word.toLowerCase().replaceAll("[^a-zA-Z]", "");
                    if (!cleanedWord.isEmpty()) {
                        wordFrequency.merge(cleanedWord, 1, Integer::sum);
                    }
                }
            }
        } catch (IOException e) {
            System.err.println("Error reading file: " + file.getFileName());
            throw e;
        }
    }
    
    private Map<String, Integer> sortByValue(Map<String, Integer> map) {
        return map.entrySet()
                .stream()
                .sorted(Map.Entry.<String, Integer>comparingByValue().reversed())
                .collect(Collectors.toMap(
                        Map.Entry::getKey,
                        Map.Entry::getValue,
                        (e1, e2) -> e1,
                        LinkedHashMap::new));
    }
}
