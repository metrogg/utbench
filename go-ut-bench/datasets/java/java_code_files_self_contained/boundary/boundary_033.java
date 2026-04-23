import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

class FileOperationsManager {
    
    /**
     * Writes multiple lines of text to a file with specified encoding.
     * Creates parent directories if they don't exist.
     * 
     * @param filePath The path to the file to write
     * @param lines List of strings to write to the file
     * @param append Whether to append to existing file or overwrite
     * @return true if operation succeeded, false otherwise
     */
    public boolean writeToFile(String filePath, List<String> lines, boolean append) {
        File file = new File(filePath);
        
        try {
            // Create parent directories if they don't exist
            File parent = file.getParentFile();
            if (parent != null && !parent.exists()) {
                if (!parent.mkdirs()) {
                    return false;
                }
            }
            
            try (BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(
                            new FileOutputStream(file, append), StandardCharsets.UTF_8))) {
                
                for (String line : lines) {
                    writer.write(line);
                    writer.newLine();
                }
                return true;
            }
        } catch (IOException e) {
            return false;
        }
    }
    
    /**
     * Reads all lines from a file with specified encoding.
     * 
     * @param filePath The path to the file to read
     * @return List of strings containing file lines, or empty list if error occurs
     */
    public List<String> readFromFile(String filePath) {
        List<String> lines = new ArrayList<>();
        File file = new File(filePath);
        
        if (!file.exists() || !file.isFile()) {
            return lines;
        }
        
        try (BufferedReader reader = new BufferedReader(
                new InputStreamReader(
                        new FileInputStream(file), StandardCharsets.UTF_8))) {
            
            String line;
            while ((line = reader.readLine()) != null) {
                lines.add(line);
            }
        } catch (IOException e) {
            // Return whatever we've read so far
        }
        
        return lines;
    }
    
    /**
     * Checks if a file exists and is readable.
     * 
     * @param filePath The path to check
     * @return true if file exists and is readable, false otherwise
     */
    public boolean fileExists(String filePath) {
        File file = new File(filePath);
        return file.exists() && file.canRead();
    }
}
