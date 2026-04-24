// obs 包提供结构化日志功能
// 基于 Go 标准库的 slog 封装，提供统一格式的日志输出
package obs

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// Logger 结构化日志记录器
// 封装 slog.Logger，提供 Info、Debug、Warn、Error 级别日志
type Logger struct {
	base *slog.Logger // 底层 slog 日志实例
}

// NewLogger 创建新的日志记录器
// 参数:
//   - verbose: 是否启用 Debug 级别日志
//
// 返回值:
//   - *Logger: 日志记录器实例
//
// 功能说明:
//   - 默认使用 Info 级别
//   - verbose=true 时启用 Debug 级别
//   - 输出到标准错误流
func NewLogger(verbose bool) *Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return &Logger{base: slog.New(handler)}
}

// NewLoggerWithWriter creates a logger that writes to both os.Stderr and the given writer.
func NewLoggerWithWriter(verbose bool, w io.Writer) *Logger {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	mw := io.MultiWriter(os.Stderr, w)
	handler := slog.NewTextHandler(mw, &slog.HandlerOptions{Level: level})
	return &Logger{base: slog.New(handler)}
}

// Info 输出信息级别日志
// 参数:
//   - msg: 日志消息
//   - attrs: 键值对属性
func (l *Logger) Info(msg string, attrs ...any) {
	l.base.Info(msg, attrs...)
}

// Debug 输出调试级别日志
// 参数:
//   - msg: 日志消息
//   - attrs: 键值对属性
func (l *Logger) Debug(msg string, attrs ...any) {
	l.base.Debug(msg, attrs...)
}

// Warn 输出警告级别日志
// 参数:
//   - msg: 日志消息
//   - attrs: 键值对属性
func (l *Logger) Warn(msg string, attrs ...any) {
	l.base.Warn(msg, attrs...)
}

// Error 输出错误级别日志
// 参数:
//   - msg: 日志消息
//   - attrs: 键值对属性
func (l *Logger) Error(msg string, attrs ...any) {
	l.base.Error(msg, attrs...)
}

// With 创建带有固定属性的日志记录器
// 参数:
//   - attrs: 键值对属性
//
// 返回值:
//   - *Logger: 新的日志记录器实例，包含指定的属性
func (l *Logger) With(attrs ...any) *Logger {
	return &Logger{base: l.base.With(attrs...)}
}

// WithContext 创建带有上下文的日志记录器（保留方法签名兼容性）
// 参数:
//   - ctx: 上下文（当前未使用）
//
// 返回值:
//   - *Logger: 日志记录器实例
func (l *Logger) WithContext(_ context.Context) *Logger {
	return l
}
