interface DataSource {
    String read();
    boolean write(String data);
}

class FileDataSource implements DataSource {
    private String filepath;
    private String data;
    
    public FileDataSource(String filepath) {
        this.filepath = filepath;
        this.data = "";
    }
    
    public String read() {
        return data;
    }
    
    public boolean write(String data) {
        this.data = data;
        return true;
    }
}

class MockDataSource implements DataSource {
    private String data;
    
    public MockDataSource() {
        this.data = "mock_data";
    }
    
    public String read() {
        return data;
    }
    
    public boolean write(String data) {
        this.data = data;
        return true;
    }
}

class DataProcessor {
    private DataSource source;
    
    public DataProcessor(DataSource source) {
        this.source = source;
    }
    
    public String process() {
        String data = source.read();
        return data.toUpperCase();
    }
    
    public boolean save(String data) {
        return source.write(data);
    }
}