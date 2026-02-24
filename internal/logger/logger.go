package logger

import (
	"fmt"
	"log/slog"
	"os"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	F "github.com/IBM/fp-go/function"
	T "github.com/IBM/fp-go/tuple"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/urfave/cli/v2"
)

func SetupCliLogger(cltx *cli.Context) error {
	l, err := NewCliLogger(cltx)
	if err != nil {
		return fmt.Errorf("error in getting a new logger %s", err)
	}
	registry.SetLogger(l)
	return nil
}

func ExtractCliLogKey(cltx *cli.Context) T.Tuple2[string, string] {
	return T.MakeTuple2(cltx.String("log-level"), cltx.String("log-format"))
}

func MakeDebugTextHandler() slog.Handler {
	return slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

func MakeDebugJSONHandler() slog.Handler {
	return slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

func MakeInfoTextHandler() slog.Handler {
	return slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
}

func MakeInfoJSONHandler() slog.Handler {
	return slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo},
	)
}

var SlogCliHandlers = map[T.Tuple2[string, string]]func(*cli.Context) slog.Handler{
	T.MakeTuple2("debug", "text"): F.Constant1[*cli.Context](
		MakeDebugTextHandler(),
	),
	T.MakeTuple2("debug", "json"): F.Constant1[*cli.Context](
		MakeDebugJSONHandler(),
	),
	T.MakeTuple2("info", "text"): F.Constant1[*cli.Context](
		MakeInfoTextHandler(),
	),
	T.MakeTuple2("info", "json"): F.Constant1[*cli.Context](
		MakeInfoJSONHandler(),
	),
}

var DefaultCliSlogHandler = F.Constant1[*cli.Context](MakeInfoJSONHandler())

func GetCliSlogHandler(cltx *cli.Context) slog.Handler {
	return F.Pipe1(
		cltx,
		F.Switch(
			ExtractCliLogKey,      // Key extractor
			SlogCliHandlers,       // Handler map
			DefaultCliSlogHandler, // Fallback
		),
	)
}

func NewCliLogger(cltx *cli.Context) (*logrus.Entry, error) {
	format := cltx.String("log-format")
	name := cltx.String("log-level")
	fname := cltx.String("log-file")
	lfmt, err := getLogFmt(format)
	if err != nil {
		return &logrus.Entry{}, err
	}
	level, err := getLogLevel(name)
	if err != nil {
		return &logrus.Entry{}, err
	}
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetFormatter(lfmt)
	logger.SetLevel(level)
	// set hook to write to local file
	if len(fname) != 0 {
		logger.Hooks.Add(lfshook.NewHook(fname, lfmt))
		registry.SetValue(registry.LogFileKey, fname)
	}
	logger.Hooks.Add(logrus_stack.StandardHook())
	return logrus.NewEntry(logger), nil
}

func NewLogger(cmd *cobra.Command) (*logrus.Entry, error) {
	format, _ := cmd.Flags().GetString("log-format")
	name, _ := cmd.Flags().GetString("log-level")
	fname, _ := cmd.Flags().GetString("log-file")
	lfmt, err := getLogFmt(format)
	if err != nil {
		return &logrus.Entry{}, err
	}
	level, err := getLogLevel(name)
	if err != nil {
		return &logrus.Entry{}, err
	}
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetFormatter(lfmt)
	logger.SetLevel(level)
	// set hook to write to local file
	if len(fname) == 0 {
		f, err := os.CreateTemp(os.TempDir(), "loader")
		if err != nil {
			return &logrus.Entry{}, fmt.Errorf(
				"error in creating temp file for logging %s",
				err,
			)
		}
		fname = f.Name()
	}
	logger.Hooks.Add(lfshook.NewHook(fname, lfmt))
	logger.Hooks.Add(logrus_stack.StandardHook())
	registry.SetValue(registry.LogFileKey, fname)
	return logrus.NewEntry(logger), nil
}

func getLogLevel(name string) (logrus.Level, error) {
	var level logrus.Level
	switch name {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warn":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		return level, fmt.Errorf(
			"%s log level is not supported",
			level,
		)
	}
	return level, nil
}

func getLogFmt(format string) (logrus.Formatter, error) {
	var lfmt logrus.Formatter
	switch format {
	case "text":
		lfmt = &logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	case "json":
		lfmt = &logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		}
	default:
		return lfmt, fmt.Errorf(
			"only json and text are supported %s log format is not supported",
			format,
		)
	}
	return lfmt, nil
}

func GetSlogHandler(cmd *cobra.Command) slog.Handler {
	return F.Pipe4(
		cmd,
		getLogLevelFromCmd,
		toSlogLevel,
		toHandlerOptions,
		toSlogHandler(cmd),
	)
}

func getLogLevelFromCmd(cmd *cobra.Command) string {
	s, _ := cmd.Flags().GetString("log-level")
	return s
}

func toSlogLevel(s string) slog.Level {
	if s == "debug" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func toHandlerOptions(l slog.Level) *slog.HandlerOptions {
	return &slog.HandlerOptions{Level: l}
}

func toSlogHandler(cmd *cobra.Command) func(*slog.HandlerOptions) slog.Handler {
	return func(opts *slog.HandlerOptions) slog.Handler {
		fmt, _ := cmd.Flags().GetString("log-format")
		if fmt == "text" {
			return slog.NewTextHandler(os.Stderr, opts)
		}
		return slog.NewJSONHandler(os.Stderr, opts)
	}
}
