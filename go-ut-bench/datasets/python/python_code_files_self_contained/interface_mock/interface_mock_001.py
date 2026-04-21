from abc import ABC, abstractmethod


class Logger(ABC):
    @abstractmethod
    def log(self, message: str) -> None:
        pass
    
    @abstractmethod
    def get_logs(self) -> list:
        pass


class ConsoleLogger(Logger):
    def __init__(self):
        self._logs = []
    
    def log(self, message: str) -> None:
        self._logs.append(message)
    
    def get_logs(self) -> list:
        return self._logs.copy()


class FileLogger(Logger):
    def __init__(self, filename: str):
        self._filename = filename
        self._logs = []
    
    def log(self, message: str) -> None:
        self._logs.append(f"[{self._filename}] {message}")
    
    def get_logs(self) -> list:
        return self._logs.copy()


class Application:
    def __init__(self, logger: Logger):
        self._logger = logger
    
    def run(self, task_name: str) -> str:
        self._logger.log(f"Starting task: {task_name}")
        result = f"Task {task_name} completed"
        self._logger.log(result)
        return result
    
    def get_history(self) -> list:
        return self._logger.get_logs()