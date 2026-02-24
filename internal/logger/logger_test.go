package logger

import (
	"context"
	"flag"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestGetCliSlogHandler(t *testing.T) {
	tests := []struct {
		name      string
		logLevel  string
		logFormat string
		wantLevel slog.Level
	}{
		{
			name:      "Debug Level",
			logLevel:  "debug",
			logFormat: "text",
			wantLevel: slog.LevelDebug,
		},
		{
			name:      "Info Level",
			logLevel:  "info",
			logFormat: "json",
			wantLevel: slog.LevelInfo,
		},
		{
			name:      "Default (Empty)",
			logLevel:  "",
			logFormat: "",
			wantLevel: slog.LevelInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := flag.NewFlagSet("test", 0)
			set.String("log-level", tt.logLevel, "")
			set.String("log-format", tt.logFormat, "")
			c := cli.NewContext(nil, set, nil)

			handler := GetCliSlogHandler(c)
			assert.True(t, handler.Enabled(context.Background(), tt.wantLevel))
			if tt.wantLevel == slog.LevelInfo {
				assert.False(t, handler.Enabled(context.Background(), slog.LevelDebug))
			}
		})
	}
}
