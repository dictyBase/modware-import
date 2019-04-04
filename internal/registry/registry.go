package registry

import (
	"io"

	minio "github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOGRUS_KEY   = "logrus"
	MINIO_KEY    = "minio"
	LOG_FILE_KEY = "log_file"
)

var v = viper.New()

func SetValue(key, value string) {
	v.Set(key, value)
}

func SetLogger(l *logrus.Entry) {
	v.Set(LOGRUS_KEY, l)
}

func SetS3Client(s3c *minio.Client) {
	v.Set(MINIO_KEY, s3c)
}

func SetReader(key string, r io.Reader) {
	v.Set(key, r)
}

func SetWriter(key string, w io.Writer) {
	v.Set(key, w)
}

func GetLogger() *logrus.Entry {
	l, _ := v.Get(LOGRUS_KEY).(*logrus.Entry)
	return l
}

func GetS3Client() *minio.Client {
	s3c, _ := v.Get(MINIO_KEY).(*minio.Client)
	return s3c
}

func GetWriter(key string) io.Writer {
	w, _ := v.Get(key).(io.Writer)
	return w
}

func GetReader(key string) io.Reader {
	r, _ := v.Get(key).(io.Reader)
	return r
}

func GetValue(key string) string {
	val, _ := v.Get(key).(string)
	return val
}
