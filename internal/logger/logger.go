package logger

import (
	"fmt"
	"io/ioutil"
	"os"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
		f, err := ioutil.TempFile(os.TempDir(), "loader")
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
