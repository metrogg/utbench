// Converted Java method
import java.util.List;
import java.util.Map;
import java.util.HashMap;
import java.util.Collections;
import java.util.Comparator;
import java.util.stream.Collectors;

class LogAnalyzer {
    
    /**
     * Analyzes a list of system logs and provides comprehensive statistics
     * including frequency analysis, time statistics, and severity distribution.
     * 
     * @param logs List of system logs to analyze
     * @return Map containing various statistics about the logs
     */
    public Map<String, Object> analyzeLogs(List<SysLog> logs) {
        if (logs == null || logs.isEmpty()) {
            return Collections.singletonMap("error", "No logs provided for analysis");
        }

        Map<String, Object> analysis = new HashMap<>();
        
        // Basic statistics
        analysis.put("totalLogs", logs.size());
        
        // User frequency analysis
        Map<String, Long> userFrequency = logs.stream()
            .collect(Collectors.groupingBy(SysLog::getUsername, Collectors.counting()));
        analysis.put("userFrequency", userFrequency);
        
        // Most active user
        String mostActiveUser = userFrequency.entrySet().stream()
            .max(Map.Entry.comparingByValue())
            .map(Map.Entry::getKey)
            .orElse("N/A");
        analysis.put("mostActiveUser", mostActiveUser);
        
        // Time statistics
        List<Long> timestamps = logs.stream()
            .map(SysLog::getTime)
            .sorted()
            .collect(Collectors.toList());
        
        if (!timestamps.isEmpty()) {
            analysis.put("earliestLog", timestamps.get(0));
            analysis.put("latestLog", timestamps.get(timestamps.size() - 1));
            analysis.put("timeRange", timestamps.get(timestamps.size() - 1) - timestamps.get(0));
        }
        
        // Severity distribution
        Map<String, Long> severityDistribution = logs.stream()
            .collect(Collectors.groupingBy(SysLog::getSeverity, Collectors.counting()));
        analysis.put("severityDistribution", severityDistribution);
        
        // Longest message
        String longestMessage = logs.stream()
            .max(Comparator.comparingInt(log -> log.getMessage().length()))
            .map(SysLog::getMessage)
            .orElse("N/A");
        analysis.put("longestMessage", longestMessage);
        
        return analysis;
    }
    
    // Simplified SysLog class to make the code self-contained
    public static class SysLog {
        private String username;
        private String message;
        private long time;
        private String severity;
        
        // Constructors, getters, and setters
        public SysLog(String username, String message, long time, String severity) {
            this.username = username;
            this.message = message;
            this.time = time;
            this.severity = severity;
        }
        
        public String getUsername() { return username; }
        public String getMessage() { return message; }
        public long getTime() { return time; }
        public String getSeverity() { return severity; }
    }
}
