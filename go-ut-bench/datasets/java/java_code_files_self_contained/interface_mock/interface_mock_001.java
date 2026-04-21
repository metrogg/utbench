interface Logger {
    void log(String message);
    String[] getLogs();
}

class ConsoleLogger implements Logger {
    private String[] logs;
    private int count;
    
    public ConsoleLogger() {
        this.logs = new String[100];
        this.count = 0;
    }
    
    public void log(String message) {
        if (count < logs.length) {
            logs[count++] = message;
        }
    }
    
    public String[] getLogs() {
        String[] result = new String[count];
        for (int i = 0; i < count; i++) {
            result[i] = logs[i];
        }
        return result;
    }
}

class FileLogger implements Logger {
    private String filename;
    private String[] logs;
    private int count;
    
    public FileLogger(String filename) {
        this.filename = filename;
        this.logs = new String[100];
        this.count = 0;
    }
    
    public void log(String message) {
        if (count < logs.length) {
            logs[count++] = "[" + filename + "] " + message;
        }
    }
    
    public String[] getLogs() {
        String[] result = new String[count];
        for (int i = 0; i < count; i++) {
            result[i] = logs[i];
        }
        return result;
    }
}

class Application {
    private Logger logger;
    
    public Application(Logger logger) {
        this.logger = logger;
    }
    
    public String run(String taskName) {
        logger.log("Starting task: " + taskName);
        String result = "Task " + taskName + " completed";
        logger.log(result);
        return result;
    }
    
    public String[] getHistory() {
        return logger.getLogs();
    }
}