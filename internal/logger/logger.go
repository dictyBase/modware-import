package logger

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewLogger(cmd *cobra.Command) (*logrus.Entry, error) {
	lfmt, err := getLogFmt(cmd)
	if err != nil {
		return &logrus.Entry{}, err
	}
	level, err := getLogLevel(cmd)
	if err != nil {
		return &logrus.Entry{}, err
	}
	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: lfmt,
		Level:     level,
	}
	// set hook to write to local file
	fname, _ := cmd.Flags().GetString("log-file")
	if len(fname) == 0 {
		f, err := ioutil.TempFile(os.TempDir(), "loader")
		if err != nil {
			return &logrus.Entry{}, fmt.Errorf("error in creating temp file for logging %s", err)
		}
		fname = f.Name()
	}
	logger.Hooks.Add(lfshook.NewHook(fname, lfmt))
	registry.SetValue(registry.LOG_FILE_KEY, fname)
	return logrus.NewEntry(logger), nil
}

func getLogLevel(cmd *cobra.Command) (logrus.Level, error) {
	var level logrus.Level
	name, _ := cmd.Flags().GetString("log-level")
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

func getLogFmt(cmd *cobra.Command) (logrus.Formatter, error) {
	var lfmt logrus.Formatter
	format, _ := cmd.Flags().GetString("log-format")
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
