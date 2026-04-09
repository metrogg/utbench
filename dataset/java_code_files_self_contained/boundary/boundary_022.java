import java.util.ArrayList;
import java.util.List;

class TaskManager {
    private List<String> tasks;
    private String filename;

    /**
     * Initializes a TaskManager with a specified filename.
     * @param filename The name of the file to store tasks
     */
    public TaskManager(String filename) {
        this.filename = filename;
        this.tasks = new ArrayList<>();
    }

    /**
     * Adds a new task to the task list.
     * @param taskDescription The description of the task to add
     * @return true if task was added successfully, false otherwise
     */
    public boolean addTask(String taskDescription) {
        if (taskDescription == null || taskDescription.trim().isEmpty()) {
            return false;
        }
        tasks.add((tasks.size() + 1) + ". " + taskDescription);
        return true;
    }

    /**
     * Deletes a task by its position or all tasks.
     * @param position The position of the task to delete ("-1" for all tasks)
     * @return true if deletion was successful, false otherwise
     */
    public boolean deleteTask(String position) {
        try {
            if (position.equals("-1")) {
                tasks.clear();
                return true;
            } else {
                int pos = Integer.parseInt(position);
                if (pos > 0 && pos <= tasks.size()) {
                    tasks.remove(pos - 1);
                    // Re-number remaining tasks
                    for (int i = 0; i < tasks.size(); i++) {
                        String task = tasks.get(i);
                        tasks.set(i, (i + 1) + task.substring(task.indexOf('.')));
                    }
                    return true;
                }
            }
        } catch (NumberFormatException e) {
            return false;
        }
        return false;
    }

    /**
     * Gets the current list of tasks.
     * @return List of tasks
     */
    public List<String> getTasks() {
        return new ArrayList<>(tasks);
    }

    /**
     * Gets the number of tasks.
     * @return Count of tasks
     */
    public int getTaskCount() {
        return tasks.size();
    }
}
