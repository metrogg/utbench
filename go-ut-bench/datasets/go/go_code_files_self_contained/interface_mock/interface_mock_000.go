package main

type DataSource interface {
	Read() string
	Write(data string) bool
}

type FileDataSource struct {
	filepath string
	data     string
}

func NewFileDataSource(filepath string) *FileDataSource {
	return &FileDataSource{filepath: filepath, data: ""}
}

func (fds *FileDataSource) Read() string {
	return fds.data
}

func (fds *FileDataSource) Write(data string) bool {
	fds.data = data
	return true
}

type MockDataSource struct {
	data string
}

func NewMockDataSource() *MockDataSource {
	return &MockDataSource{data: "mock_data"}
}

func (mds *MockDataSource) Read() string {
	return mds.data
}

func (mds *MockDataSource) Write(data string) bool {
	mds.data = data
	return true
}

type DataProcessor struct {
	source DataSource
}

func NewDataProcessor(source DataSource) *DataProcessor {
	return &DataProcessor{source: source}
}

func (dp *DataProcessor) Process() string {
	data := dp.source.Read()
	result := ""
	for _, c := range data {
		if c >= 'a' && c <= 'z' {
			result += string(c - 32)
		} else {
			result += string(c)
		}
	}
	return result
}

func (dp *DataProcessor) Save(data string) bool {
	return dp.source.Write(data)
}