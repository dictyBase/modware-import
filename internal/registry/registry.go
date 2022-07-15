package registry

import (
	"io"

	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/go-obograph/storage"
	r "github.com/go-redis/redis/v7"
	"github.com/minio/minio-go/v6"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
)

const (
	LogrusKey          = "logrus"
	MinioKey           = "minio"
	LogFileKey         = "log_file"
	RedisKey           = "redis"
	ArangodbSessionKey = "arangodb_session"
	Arangodb           = "arangodb"
	OboStorageKey      = "obostorage"
	OboReadersKey      = "oboreaders"
	KubeClientKey      = "kubeconfig"
)

var v = viper.New()

func SetValue(key, value string) {
	v.Set(key, value)
}

func SetArangoOboStorage(s storage.DataSource) {
	v.Set(OboStorageKey, s)
}

func SetArangoSession(s *arangomanager.Session) {
	v.Set(ArangodbSessionKey, s)
}

func SetArangodbConnection(c *arangomanager.Database) {
	v.Set(Arangodb, c)
}

func SetLogger(l *logrus.Entry) {
	v.Set(LogrusKey, l)
}

func SetS3Client(s3c *minio.Client) {
	v.Set(MinioKey, s3c)
}

func SetReader(key string, r io.Reader) {
	v.Set(key, r)
}

func SetAllReaders(key string, rds map[string]io.Reader) {
	v.Set(key, rds)
}

func SetWriter(key string, w io.Writer) {
	v.Set(key, w)
}

func SetRedisClient(redis *r.Client) {
	v.Set(RedisKey, redis)
}

func SetKubeClient(key string, client *kubernetes.Clientset) {
	v.Set(key, client)
}

func GetArangoOboStorage() storage.DataSource {
	s, _ := v.Get(OboStorageKey).(storage.DataSource)
	return s
}

func GetArangoSession() *arangomanager.Session {
	s, _ := v.Get(ArangodbSessionKey).(*arangomanager.Session)
	return s
}

func GetArangodbConnection() *arangomanager.Database {
	c, _ := v.Get(Arangodb).(*arangomanager.Database)
	return c
}

func GetLogger() *logrus.Entry {
	l, _ := v.Get(LogrusKey).(*logrus.Entry)
	return l
}

func GetS3Client() *minio.Client {
	s3c, _ := v.Get(MinioKey).(*minio.Client)
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

func GetAllReaders(key string) map[string]io.Reader {
	rds, _ := v.Get(key).(map[string]io.Reader)
	return rds
}

func GetValue(key string) string {
	val, _ := v.Get(key).(string)
	return val
}

func GetRedisClient() *r.Client {
	redis, _ := v.Get(RedisKey).(*r.Client)
	return redis
}

func GetKubeClient(key string) *kubernetes.Clientset {
	client, _ := v.Get(key).(*kubernetes.Clientset)
	return client
}
