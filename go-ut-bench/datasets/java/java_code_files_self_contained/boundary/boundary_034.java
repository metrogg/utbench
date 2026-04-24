// Converted Java method
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.zip.*;

class FileHandler {

    /**
     * Handles file operations including upload, download, and zip creation
     * @param operationType Type of operation ("upload", "download", "zip")
     * @param files List of files to process (for upload/zip)
     * @param targetDir Target directory for operations
     * @param fileName Specific file name for download
     * @return Operation result message or byte array of downloaded file
     * @throws IOException If file operations fail
     */
    public static Object handleFileOperation(String operationType, 
                                           List<File> files, 
                                           String targetDir,
                                           String fileName) throws IOException {
        switch (operationType.toLowerCase()) {
            case "upload":
                return handleUpload(files, targetDir);
            case "download":
                return handleDownload(targetDir, fileName);
            case "zip":
                return handleZipCreation(files, targetDir);
            default:
                throw new IllegalArgumentException("Invalid operation type");
        }
    }

    private static String handleUpload(List<File> files, String targetDir) throws IOException {
        Path targetPath = Paths.get(targetDir);
        if (!Files.exists(targetPath)) {
            Files.createDirectories(targetPath);
        }

        for (File file : files) {
            Path dest = targetPath.resolve(file.getName());
            Files.copy(file.toPath(), dest, StandardCopyOption.REPLACE_EXISTING);
        }
        return "Upload successful: " + files.size() + " files processed";
    }

    private static byte[] handleDownload(String targetDir, String fileName) throws IOException {
        Path filePath = Paths.get(targetDir, fileName);
        return Files.readAllBytes(filePath);
    }

    private static byte[] handleZipCreation(List<File> files, String targetDir) throws IOException {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        try (ZipOutputStream zos = new ZipOutputStream(baos)) {
            for (File file : files) {
                ZipEntry entry = new ZipEntry(file.getName());
                zos.putNextEntry(entry);
                byte[] bytes = Files.readAllBytes(file.toPath());
                zos.write(bytes, 0, bytes.length);
                zos.closeEntry();
            }
        }
        return baos.toByteArray();
    }
}
