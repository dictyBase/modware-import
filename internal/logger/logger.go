package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewLogger(cmd *cobra.Command) (*logrus.Entry, error) {
	e := new(logrus.Entry)
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	format, _ := cmd.Flags().GetString("log-format")
	switch format {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		})
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "02/Jan/2006:15:04:05",
		})
	default:
		return e, fmt.Errorf(
			"only json and text are supported %s log format is not supported",
			format,
		)
	}
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
	return logrus.NewEntry(logger), nil
}
