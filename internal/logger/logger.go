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
	e := new(logrus.Entry)
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	format, _ := cmd.Flags().GetString("log-format")
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
		return e, fmt.Errorf(
			"only json and text are supported %s log format is not supported",
			format,
		)
	}
	logger.SetFormatter(lfmt)
	level, _ := cmd.Flags().GetString("log-level")
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		return e, fmt.Errorf(
			"%s log level is not supported",
			level,
		)
	}
	// set hook to write to local file
	fname, _ := cmd.Flags().GetString("log-file")
	if len(fname) == 0 {
		f, err := ioutil.TempFile(os.TempDir(), "loader")
		if err != nil {
			return e, fmt.Errorf("error in creating temp file for logging %s", err)
		}
		fname = f.Name()
	}
	logger.Hooks.Add(lfshook.NewHook(fname, lfmt))
	registry.SetValue(registry.LOG_FILE_KEY, fname)
	return logrus.NewEntry(logger), nil
}
