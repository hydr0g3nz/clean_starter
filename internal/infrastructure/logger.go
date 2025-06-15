package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hydr0g3nz/wallet_topup_system/internal/domain/infra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// AppLogger defines the interface for application logging
// NOTE: In a real Clean Architecture, this interface should be defined
// in a higher-level package (like "app" or "shared"), not infrastructure.
// It's defined here for simplicity of this single file example.

// Logger implements the AppLogger interface using zap
type Logger struct {
	zap *zap.Logger
}

// NewLogger creates a new logger instance and returns it as the AppLogger interface
func NewLogger(isProduction bool) (*Logger, error) { // Return the interface type
	var config zap.Config
	var err error

	if isProduction {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// config.Level.SetLevel(zapcore.InfoLevel) // Adjust level if needed
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// config.Level.SetLevel(zapcore.DebugLevel) // Keep Debug for development
	}

	// Create log directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("can't create log directory: %w", err)
	}

	// Set up log file path (basic daily file for demonstration)
	// For production, use a robust file rotation library like lumberjack
	logFile := filepath.Join(logDir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))

	// Create a file syncer, handling potential errors gracefully
	fileSyncer := zapcore.AddSync(os.Stderr) // Default to stderr if file opening fails
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to open log file %s: %v. Logging only to stdout/stderr.\n", logFile, err)
	} else {
		fileSyncer = zapcore.AddSync(file)
		// Note: Proper file handle closing on application exit is important
		// but not fully implemented in this basic example's Close/Sync.
	}

	// Create cores for writing to file (JSON) and stdout (Console)
	cores := []zapcore.Core{}

	// Add file core if file was successfully opened
	if file != nil && err == nil {
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(config.EncoderConfig),
			fileSyncer,
			config.Level,
		)
		cores = append(cores, fileCore)
	}

	// Always add stdout core
	stdoutCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		config.Level,
	)
	cores = append(cores, stdoutCore)

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger
	// Add SkipCaller 1 to skip the wrapper method itself in the log output
	zapLogger := zap.New(core, zap.AddCallerSkip(1))

	// Return the concrete struct, cast as the interface type
	return &Logger{zapLogger}, nil
}

// mapToZapFields converts a map[string]interface{} into a slice of zapcore.Field
func mapToZapFields(fields map[string]interface{}) []zapcore.Field {
	if len(fields) == 0 {
		return nil // Return nil or an empty slice if no fields
	}

	// Use make with initial capacity for efficiency
	zapFields := make([]zapcore.Field, 0, len(fields))
	for key, value := range fields {
		// Use zap.Any to handle various types from interface{}
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}

// Implement the AppLogger methods
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.zap.Debug(msg, toZapFields(fields...)...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zap.Sugar().Debugf(format, args...)
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.zap.Info(msg, toZapFields(fields...)...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.zap.Sugar().Infof(format, args...)
}

func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.zap.Warn(msg, toZapFields(fields...)...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zap.Sugar().Warnf(format, args...)
}

func (l *Logger) Error(msg string, fields ...interface{}) {
	l.zap.Error(msg, toZapFields(fields...)...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zap.Sugar().Errorf(format, args...)
}

func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.zap.Fatal(msg, toZapFields(fields...)...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zap.Sugar().Fatalf(format, args...)
}

func (l *Logger) Sync() error {
	// Sync attempts to flush buffered logs.
	// If you were managing the file handle directly, add file.Close() here.
	return l.zap.Sync()
}
func (l *Logger) With(fields ...interface{}) infra.Logger {
	return &Logger{
		zap: l.zap.With(toZapFields(fields...)...),
	}
}

// Optional Close method if you need explicit resource cleanup
func (l *Logger) Close() error {
	syncErr := l.zap.Sync()
	// Add file.Close() logic here if you stored the file handle
	return syncErr
}
func toZapFields(fields ...interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields)/2)

	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue // ข้ามถ้า key ไม่ใช่ string
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}

	return zapFields
}
