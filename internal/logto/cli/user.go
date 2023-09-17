package cli

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
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
	tcache := registry.GetTTLCache()
	reader := csv.NewReader(registry.GetReader("USER_INPUT"))
	lclient := logto.NewClient(cltx.String("endpoint"))
	header := false
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return cli.Exit(
				fmt.Sprintf("error in reading csv record %s", err),
				2,
			)
		}
		if !header {
			header = true
			continue
		}
		if record[1] != "Valid" {
			logger.Debugf("user with email %s in not valid", record[0])
			continue
		}
		token, err := retrieveToken(&retrieveTokenProperties{
			Lclient: lclient,
			Logger:  logger,
			Cltx:    cltx,
			Tcache:  tcache,
		})
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		ok, _, err := lclient.CheckUser(
			token,
			record[0],
		)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		if ok {
			logger.Infof(
				"user with email %s exist, skipping creation",
				record[0],
			)
			continue
		}
		normUser := normalizeUserName(record[2], record[3])
		ok, _, err = lclient.CheckUserWithUserName(
			token,
			normUser,
		)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
		if ok {
			logger.Infof(
				"user with username %s exist, skipping creation",
				normUser,
			)
			continue
		}
		logger.Debugf("username %s does not exist, going to create", record[0])
		userId, err := lclient.CreateUser(
			token,
			&logto.APIUsersPostReq{
				PrimaryEmail: record[0],
				Username:     normUser,
				Name:         fmt.Sprintf("%s %s", record[2], record[3]),
				PrimaryPhone: FixedLenRandomInt(10),
				Password:     FixedLenRandomString(80),
			},
		)
		if err != nil {
			return cli.Exit(
				fmt.Sprintf("error in creating user %s %s", record[0], err),
				2,
			)
		}
		logger.Infof("created user with email %s id %s\n", record[0], userId)
	}
	return nil
}
