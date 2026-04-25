// obs 包提供结构化日志功能
// 基于 Go 标准库的 slog 封装，提供统一格式的日志输出
package obs

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level 定义日志级别
type Level int

const (
	LevelTrace Level = iota - 2 // 比 Debug 更细
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (l Level) toSlogLevel() slog.Level {
	switch l {
	case LevelTrace:
		return slog.Level(-8) // 比 Debug(-4) 更低
	case LevelDebug:
		return slog.LevelDebug
	case LevelInfo:
		return slog.LevelInfo
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Logger 结构化日志记录器
// 支持终端输出（中文、简洁）和文件输出（JSON、完整）
type Logger struct {
	terminalHandler slog.Handler
	fileHandlers    map[string]slog.Handler
	context         map[string]any
	level           Level
	mu              sync.Mutex
}

// NewLogger 创建新的日志记录器
// 参数:
//   - verbose: 是否启用 Debug 级别日志
//   - logDir: 文件日志目录（为空则不输出文件日志）
//
// 返回值:
//   - *Logger: 日志记录器实例
func NewLogger(verbose bool, logDir string) *Logger {
	level := LevelInfo
	if verbose {
		level = LevelDebug
	}

	terminalHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level.toSlogLevel(),
	})

	fileHandlers := make(map[string]slog.Handler)
	if logDir != "" {
		os.MkdirAll(logDir, 0o755)
		categories := []string{"run", "runner", "evaluator", "api", "errors"}
		for _, cat := range categories {
			path := filepath.Join(logDir, cat+".log")
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
			if err == nil {
				fileHandlers[cat] = slog.NewJSONHandler(f, &slog.HandlerOptions{
					Level: slog.Level(-8), // Trace 级别
				})
			}
		}
	}

	return &Logger{
		terminalHandler: terminalHandler,
		fileHandlers:    fileHandlers,
		context:         make(map[string]any),
		level:           level,
	}
}

// NewLoggerWithWriter creates a logger whose terminal stream is the provided writer.
// File category logs are disabled because this is used for per-run Web log capture.
func NewLoggerWithWriter(verbose bool, w io.Writer) *Logger {
	level := LevelInfo
	if verbose {
		level = LevelDebug
	}
	if w == nil {
		w = os.Stderr
	}
	return &Logger{
		terminalHandler: slog.NewTextHandler(w, &slog.HandlerOptions{
			Level: level.toSlogLevel(),
		}),
		fileHandlers: make(map[string]slog.Handler),
		context:      make(map[string]any),
		level:        level,
	}
}

// log 内部日志方法
func (l *Logger) log(lvl Level, msg string, attrs ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	allAttrs := make([]slog.Attr, 0, len(l.context)*2+len(attrs)/2)
	for k, v := range l.context {
		allAttrs = append(allAttrs, slog.Any(k, v))
	}
	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok {
			allAttrs = append(allAttrs, slog.Any(key, attrs[i+1]))
		}
	}

	record := slog.Record{
		Time:    time.Now(),
		Level:   lvl.toSlogLevel(),
		Message: msg,
	}
	if len(allAttrs) > 0 {
		record.AddAttrs(allAttrs...)
	}

	if lvl >= l.level {
		l.terminalHandler.Handle(context.Background(), record)
	}

	for _, handler := range l.fileHandlers {
		handler.Handle(context.Background(), record)
	}
}

// Trace 输出 Trace 级别日志
// 用于 API 请求/响应详情、HTTP headers、Token消耗等
func (l *Logger) Trace(msg string, attrs ...any) {
	l.log(LevelTrace, msg, attrs...)
}

// Debug 输出 Debug 级别日志
// 用于内部状态变化、临时文件路径、命令执行详情
func (l *Logger) Debug(msg string, attrs ...any) {
	l.log(LevelDebug, msg, attrs...)
}

// Info 输出 Info 级别日志
// 用于阶段开始/结束、关键进度、统计摘要
func (l *Logger) Info(msg string, attrs ...any) {
	l.log(LevelInfo, msg, attrs...)
}

// Warn 输出 Warn 级别日志
// 用于非致命错误、降级处理
func (l *Logger) Warn(msg string, attrs ...any) {
	l.log(LevelWarn, msg, attrs...)
}

// Error 输出 Error 级别日志
// 用于致命错误、编译失败、测试失败
func (l *Logger) Error(msg string, attrs ...any) {
	l.log(LevelError, msg, attrs...)
}

// With 创建带有固定属性的日志记录器
func (l *Logger) With(attrs ...any) *Logger {
	newLogger := &Logger{
		terminalHandler: l.terminalHandler,
		fileHandlers:    l.fileHandlers,
		context:         make(map[string]any),
		level:           l.level,
	}

	for k, v := range l.context {
		newLogger.context[k] = v
	}

	for i := 0; i < len(attrs)-1; i += 2 {
		if key, ok := attrs[i].(string); ok {
			newLogger.context[key] = attrs[i+1]
		}
	}

	return newLogger
}

// ToFile 创建指向特定分类文件的 Logger
func (l *Logger) ToFile(category string) *Logger {
	newLogger := &Logger{
		terminalHandler: l.terminalHandler,
		fileHandlers:    make(map[string]slog.Handler),
		context:         make(map[string]any),
		level:           l.level,
	}

	for k, v := range l.context {
		newLogger.context[k] = v
	}

	if handler, ok := l.fileHandlers[category]; ok {
		newLogger.fileHandlers[category] = handler
	}

	return newLogger
}

// WithContext 创建带有上下文的日志记录器（保留方法签名兼容性）
func (l *Logger) WithContext(_ context.Context) *Logger {
	return l
}

// LogError 记录错误到 errors.log 和终端
func (l *Logger) LogError(stage, msg string, err error, extra ...any) {
	attrs := []any{"stage", stage, "error", err.Error()}
	attrs = append(attrs, extra...)
	l.log(LevelError, msg, attrs...)
}

// LogAPIRequest 记录 API 请求详情
func (l *Logger) LogAPIRequest(model, lang, sampleID string, tokens int, latencyMS int) {
	l.ToFile("api").Trace("api_request",
		"model", model,
		"language", lang,
		"sample_id", sampleID,
		"tokens", tokens,
		"latency_ms", latencyMS,
	)
}

// LogAPIResponse 记录 API 响应详情
func (l *Logger) LogAPIResponse(model, lang, sampleID string, success bool, truncated bool, errMsg string) {
	attrs := []any{"model", model, "language", lang, "sample_id", sampleID, "success", success, "truncated", truncated}
	if errMsg != "" {
		attrs = append(attrs, "error", errMsg)
	}
	l.ToFile("api").Trace("api_response", attrs...)
}

// MarshalJSON 用于 JSON 序列化（调试用）
func (l *Logger) MarshalJSON() ([]byte, error) {
	type alias Logger
	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(l),
	})
}
