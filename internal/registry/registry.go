package registry

import (
	"io"

	"github.com/dictyBase/arangomanager"
	r "github.com/go-redis/redis/v7"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	LOGRUS_KEY           = "logrus"
	MINIO_KEY            = "minio"
	LOG_FILE_KEY         = "log_file"
	REDIS_KEY            = "redis"
	ARANGODB_SESSION_KEY = "arangodb_session"
	ARANGODB             = "arangodb"
)

var v = viper.New()

func SetValue(key, value string) {
	v.Set(key, value)
}

func SetArangoSession(s *arangomanager.Session) {
	v.Set(ARANGODB_SESSION_KEY, s)
}

func SetArangodbConnection(c *arangomanager.Database) {
	v.Set(ARANGODB, c)
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

func SetRedisClient(redis *r.Client) {
	v.Set(REDIS_KEY, redis)
}

func GetArangoSession() *arangomanager.Session {
	s, _ := v.Get(ARANGODB_SESSION_KEY).(*arangomanager.Session)
	return s
}

func GetArangodbConnection() *arangomanager.Database {
	c, _ := v.Get(ARANGODB).(*arangomanager.Database)
	return c
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

func GetRedisClient() *r.Client {
	redis, _ := v.Get(REDIS_KEY).(*r.Client)
	return redis
}
