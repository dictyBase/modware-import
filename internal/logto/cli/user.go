package cli

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"time"

	logto "github.com/dictyBase/modware-import/internal/logto/client"
	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/jellydator/ttlcache/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const TokenKey = "token"

var unameRgxp = regexp.MustCompile(`\s+|[.\-\'\(\)\,\?\"]`)

func normalizeUserName(first, last string) string {
	return unameRgxp.ReplaceAllString(fmt.Sprintf("%s%s", first, last), "")
}

// Generate a random number using crypto/rand.
func RandomInt(num int) (int, error) {
	randomValue, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		return 0, err
	}

	return int(randomValue.Int64()), nil
}

func FixedLenRandomInt(length int) string {
	num := []byte("123456789")
	byt := make([]byte, 0)
	alen := len(num)
	for i := 0; i < length; i++ {
		pos, _ := RandomInt(alen)
		byt = append(byt, num[pos])
	}

	return string(byt)
}

func FixedLenRandomString(length int) string {
	alphanum := []byte(
		"123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	)
	byt := make([]byte, 0)
	alen := len(alphanum)
	for i := 0; i < length; i++ {
		pos, _ := RandomInt(alen)
		byt = append(byt, alphanum[pos])
	}

	return string(byt)
}

type retrieveTokenProperties struct {
	Lclient *logto.Client
	Cltx    *cli.Context
	Logger  *logrus.Entry
	Tcache  *ttlcache.Cache[string, string]
}

func retrieveToken(args *retrieveTokenProperties) (string, error) {
	logger := args.Logger
	cache := args.Tcache
	cltx := args.Cltx
	item := cache.Get(TokenKey)
	if item != nil {
		logger.Debug("access token not expired getting from cache")
		return item.Value(), nil
	}
	logger.Debug("retrieving a fresh access token")
	aresp, err := args.Lclient.AccessToken(
		cltx.String("app-id"),
		cltx.String("app-secret"),
		cltx.String("api-resource"),
	)
	if err != nil {
		return item.Value(), err
	}
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", aresp.ExpiresIn-1000))
	if err != nil {
		return item.Value(), fmt.Errorf(
			"error in parsing duration %d",
			aresp.ExpiresIn,
		)
	}
	cache.Set(TokenKey, aresp.AccessToken, dur)
	logger.Debug("cached the new access token")
	return aresp.AccessToken, nil
}

func ImportUser(cltx *cli.Context) error {
	logger := registry.GetLogger()
	lclient := logto.NewClient(cltx.String("endpoint"))
	token, err := retrieveToken(&retrieveTokenProperties{
		Lclient: lclient,
		Logger:  logger,
		Cltx:    cltx,
		Tcache:  registry.GetTTLCache(),
	})
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	userId, err := lclient.CreateUser(
		token,
		&logto.APIUsersPostReq{
			PrimaryEmail: "bola@bola.com",
			PrimaryPhone: "19343049303438",
			Username:     "hello",
			Password:     "r93r938493*7043",
			Name:         "bola",
		},
	)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}
	logger.Infof("got user id %s\n", userId)
	return nil
}
