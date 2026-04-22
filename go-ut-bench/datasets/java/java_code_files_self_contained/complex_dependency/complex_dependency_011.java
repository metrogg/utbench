// Converted Java method
import java.util.*;
import java.text.SimpleDateFormat;

class TaskDatabaseManager {
    private final String databaseName;
    private int databaseVersion;
    private List<Map<String, Object>> taskCache;
    private SimpleDateFormat dateFormat;

    public TaskDatabaseManager(String dbName, int version) {
        this.databaseName = dbName;
        this.databaseVersion = version;
        this.taskCache = new ArrayList<>();
        this.dateFormat = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
    }

    /**
     * Simulates inserting a task into the database with comprehensive task data
     * @param taskType The type of task
     * @param taskName The name of the task
     * @param taskContent Detailed content of the task
     * @param studentId Associated student ID
     * @param studentName Associated student name
     * @return Map containing insertion status and generated task ID
     */
    public Map<String, Object> insertTask(int taskType, String taskName, String taskContent, 
                                        int studentId, String studentName) {
        Map<String, Object> task = new HashMap<>();
        String taskId = UUID.randomUUID().toString();
        String currentTime = dateFormat.format(new Date());
        
        task.put("taskId", taskId);
        task.put("taskType", taskType);
        task.put("taskName", taskName);
        task.put("taskContent", taskContent);
        task.put("studentId", studentId);
        task.put("studentName", studentName);
        task.put("creationTime", currentTime);
        task.put("successCount", 0);
        task.put("failCount", 0);
        task.put("score", 0);
        task.put("status", "pending");
        
        taskCache.add(task);
        
        Map<String, Object> result = new HashMap<>();
        result.put("status", "success");
        result.put("taskId", taskId);
        result.put("affectedRows", 1);
        return result;
    }

    /**
     * Queries tasks based on various criteria with pagination support
     * @param filters Map of filter criteria (can include taskType, studentId, status, etc.)
     * @param page Page number for pagination (1-based)
     * @param pageSize Number of items per page
     * @return Map containing query results and pagination info
     */
    public Map<String, Object> queryTasks(Map<String, Object> filters, int page, int pageSize) {
        List<Map<String, Object>> results = new ArrayList<>();
        
        // Apply filters
        for (Map<String, Object> task : taskCache) {
            boolean matches = true;
            for (Map.Entry<String, Object> filter : filters.entrySet()) {
                if (!task.containsKey(filter.getKey()) || 
                    !task.get(filter.getKey()).equals(filter.getValue())) {
                    matches = false;
                    break;
                }
            }
            if (matches) {
                results.add(new HashMap<>(task));
            }
        }
        
        // Apply pagination
        int totalItems = results.size();
        int totalPages = (int) Math.ceil((double) totalItems / pageSize);
        page = Math.max(1, Math.min(page, totalPages));
        
        int fromIndex = (page - 1) * pageSize;
        int toIndex = Math.min(fromIndex + pageSize, totalItems);
        List<Map<String, Object>> paginatedResults = results.subList(fromIndex, toIndex);
        
        Map<String, Object> response = new HashMap<>();
        response.put("tasks", paginatedResults);
        response.put("currentPage", page);
        response.put("totalPages", totalPages);
        response.put("totalItems", totalItems);
        
        return response;
    }

    /**
     * Updates task statistics (success count, fail count, score)
     * @param taskId ID of the task to update
     * @param successIncrement How much to add to success count
     * @param failIncrement How much to add to fail count
     * @param scoreIncrement How much to add to score
     * @return Map containing update status and affected fields
     */
    public Map<String, Object> updateTaskStats(String taskId, int successIncrement, 
                                             int failIncrement, int scoreIncrement) {
        for (Map<String, Object> task : taskCache) {
            if (task.get("taskId").equals(taskId)) {
                int currentSuccess = (int) task.getOrDefault("successCount", 0);
                int currentFail = (int) task.getOrDefault("failCount", 0);
                int currentScore = (int) task.getOrDefault("score", 0);
                
                task.put("successCount", currentSuccess + successIncrement);
                task.put("failCount", currentFail + failIncrement);
                task.put("score", currentScore + scoreIncrement);
                
                // Update status based on new stats
                if ((currentSuccess + successIncrement) > 0) {
                    task.put("status", "completed");
                }
                
                Map<String, Object> result = new HashMap<>();
                result.put("status", "success");
                result.put("taskId", taskId);
                result.put("updatedFields", Arrays.asList(
                    "successCount", "failCount", "score", "status"));
                return result;
            }
        }
        
        Map<String, Object> result = new HashMap<>();
        result.put("status", "error");
        result.put("message", "Task not found");
        return result;
    }
}
