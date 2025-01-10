package logging

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	t.Parallel()

	logger := NewLogger("debug", true)

	if logger == nil {
		t.Error("Expected logger to be created")
	}

	if logger.GetLevel() != zerolog.DebugLevel {
		t.Errorf("Expected logger level to be debug, got %s", logger.GetLevel())
	}

	logger = NewLogger("info", false)

	if logger == nil {
		t.Error("Expected logger to be created")
	}

	if logger.GetLevel() != zerolog.InfoLevel {
		t.Errorf("Expected logger level to be info, got %s", logger.GetLevel())
	}
}

func TestDefaultLogger(t *testing.T) {
	t.Parallel()

	logger := DefaultLogger()

	if logger == nil {
		t.Error("Expected logger to be created")
	}

	logger2 := DefaultLogger()

	if logger2 == nil {
		t.Error("Expected logger to be created")
	}

	if logger != logger2 {
		t.Error("Expected logger to be a singleton")
	}
}

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger1 := FromContext(ctx)
	if logger1 == nil {
		t.Fatal("expected logger to never be nil")
	}

	ctx = WithLogger(ctx, logger1)

	logger2 := FromContext(ctx)
	if logger1 != logger2 {
		t.Errorf("expected %#v to be %#v", logger1, logger2)
	}
}

func TestNewLoggerFromEnv(t *testing.T) {
	t.Parallel()

	logger := NewLoggerFromEnv()

	if logger == nil {
		t.Error("Expected logger to be created")
	}
}

func TestWithLogger(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := NewLogger("debug", true)

	ctx = WithLogger(ctx, logger)

	if FromContext(ctx) != logger {
		t.Error("Expected logger to be attached to context")
	}
}