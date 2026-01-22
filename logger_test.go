package p2p

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSlogLogger_WithNil(t *testing.T) {
	logger := NewSlogLogger(nil)
	require.NotNil(t, logger)
	assert.NotNil(t, logger.logger)
}

func TestNewSlogLogger_WithCustomLogger(t *testing.T) {
	var buf bytes.Buffer

	customLogger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	logger := NewSlogLogger(customLogger)
	require.NotNil(t, logger)
	assert.Equal(t, customLogger, logger.logger)
}

func TestSlogLogger_Infof(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Infof("test message %s %d", "arg1", 42)

	output := buf.String()
	assert.Contains(t, output, "test message arg1 42")
	assert.Contains(t, output, "INFO")
}

func TestSlogLogger_Errorf(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Errorf("error: %s", "something went wrong")

	output := buf.String()
	assert.Contains(t, output, "error: something went wrong")
	assert.Contains(t, output, "ERROR")
}

func TestSlogLogger_Debugf(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Debugf("debug info: %v", map[string]int{"key": 123})

	output := buf.String()
	assert.Contains(t, output, "debug info:")
	assert.Contains(t, output, "DEBUG")
}

func TestSlogLogger_Warnf(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelWarn})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Warnf("warning: %s", "be careful")

	output := buf.String()
	assert.Contains(t, output, "warning: be careful")
	assert.Contains(t, output, "WARN")
}

func TestSlogLogger_DebugFiltered(t *testing.T) {
	var buf bytes.Buffer
	// Set level to INFO, which should filter out DEBUG
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Debugf("this should not appear")

	output := buf.String()
	assert.Empty(t, strings.TrimSpace(output), "debug message should be filtered when level is INFO")
}

func TestSlogLogger_FormatArgs(t *testing.T) {
	var buf bytes.Buffer

	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	customLogger := slog.New(handler)
	logger := NewSlogLogger(customLogger)

	logger.Infof("values: %s, %d, %f, %t", "string", 123, 45.67, true)

	output := buf.String()
	assert.Contains(t, output, "string")
	assert.Contains(t, output, "123")
	assert.Contains(t, output, "45.67")
	assert.Contains(t, output, "true")
}
