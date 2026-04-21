#include <string>
#include <stdexcept>

class DataSource {
public:
    virtual std::string read() = 0;
    virtual bool write(const std::string& data) = 0;
    virtual ~DataSource() {}
};

class FileDataSource : public DataSource {
private:
    std::string filepath;
    std::string data;
    
public:
    FileDataSource(const std::string& path) : filepath(path), data("") {}
    
    std::string read() override {
        return data;
    }
    
    bool write(const std::string& newData) override {
        data = newData;
        return true;
    }
};

class MockDataSource : public DataSource {
private:
    std::string data;
    
public:
    MockDataSource() : data("mock_data") {}
    
    std::string read() override {
        return data;
    }
    
    bool write(const std::string& newData) override {
        data = newData;
        return true;
    }
};

class DataProcessor {
private:
    DataSource* source;
    
public:
    DataProcessor(DataSource* src) : source(src) {}
    
    std::string process() {
        std::string data = source->read();
        std::string result;
        for (char c : data) {
            if (c >= 'a' && c <= 'z') {
                result += static_cast<char>(c - 32);
            } else {
                result += c;
            }
        }
        return result;
    }
    
    bool save(const std::string& data) {
        return source->write(data);
    }
};