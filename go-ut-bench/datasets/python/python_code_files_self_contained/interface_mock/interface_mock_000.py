from abc import ABC, abstractmethod


class DataSource(ABC):
    @abstractmethod
    def read(self) -> str:
        pass
    
    @abstractmethod
    def write(self, data: str) -> bool:
        pass


class FileDataSource(DataSource):
    def __init__(self, filepath: str):
        self._filepath = filepath
        self._data = ""
    
    def read(self) -> str:
        return self._data
    
    def write(self, data: str) -> bool:
        self._data = data
        return True


class MockDataSource(DataSource):
    def __init__(self):
        self._data = "mock_data"
    
    def read(self) -> str:
        return self._data
    
    def write(self, data: str) -> bool:
        self._data = data
        return True


class DataProcessor:
    def __init__(self, source: DataSource):
        self._source = source
    
    def process(self) -> str:
        data = self._source.read()
        return data.upper()
    
    def save(self, data: str) -> bool:
        return self._source.write(data)