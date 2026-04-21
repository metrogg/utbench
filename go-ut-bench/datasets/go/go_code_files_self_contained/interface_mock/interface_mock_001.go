package main

type Logger interface {
	Log(message string)
	GetLogs() []string
}

type ConsoleLogger struct {
	logs []string
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{logs: []string{}}
}

func (cl *ConsoleLogger) Log(message string) {
	cl.logs = append(cl.logs, message)
}

func (cl *ConsoleLogger) GetLogs() []string {
	result := make([]string, len(cl.logs))
	for i, log := range cl.logs {
		result[i] = log
	}
	return result
}

type FileLogger struct {
	filename string
	logs     []string
}

func NewFileLogger(filename string) *FileLogger {
	return &FileLogger{filename: filename, logs: []string{}}
}

func (fl *FileLogger) Log(message string) {
	fl.logs = append(fl.logs, "["+fl.filename+"] "+message)
}

func (fl *FileLogger) GetLogs() []string {
	result := make([]string, len(fl.logs))
	for i, log := range fl.logs {
		result[i] = log
	}
	return result
}

type Application struct {
	logger Logger
}

func NewApplication(logger Logger) *Application {
	return &Application{logger: logger}
}

func (app *Application) Run(taskName string) string {
	app.logger.Log("Starting task: " + taskName)
	result := "Task " + taskName + " completed"
	app.logger.Log(result)
	return result
}

func (app *Application) GetHistory() []string {
	return app.logger.GetLogs()
}