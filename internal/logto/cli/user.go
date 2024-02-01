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

// FixedLenRandomInt generates a random string of fixed length using digits 1-9.
// It takes an integer length as input and returns the generated string.
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

// FixedLenRandomString generates a random string of fixed length.
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

// ImportUser is a function that imports user data from a CSV file and creates users in a logto system.
// It takes a cli.Context parameter and returns an error if an error occurs during execution.
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
		err = processCSVRecord(record, lclient, logger, cltx, tcache)
		if err != nil {
			return cli.Exit(err.Error(), 2)
		}
	}

	return nil
}

func addCustomUserInformation(
	lclient *logto.Client,
	token string,
	userId string,
	record []string,
) error {
	isSubscribed := false
	if record[15] == "Y" {
		isSubscribed = true
	}
	// Call lclient.AddCustomUserInformation with provided token, userId, and custom data
	err := lclient.AddCustomUserInformation(
		token,
		userId,
		&logto.APIUsersPatchCustomData{
			CustomData: logto.AdditionalUserInformation{
				Profession:       record[5],
				JobTitle:         record[6],
				Institution:      record[7],
				Address:          record[8],
				SecondaryAddress: record[9],
				City:             record[10],
				State:            record[11],
				Region:           record[12],
				Country:          record[13],
				Zipcode:          record[14],
				Subscribed:       isSubscribed,
				Phone:            record[16],
				ResearchInterest: record[17],
			},
		},
	)

	return err
}

func createUser(
	lclient *logto.Client,
	token string,
	record []string,
	normUser string,
) (string, error) {
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

	return userId, err
}

func processCSVRecord(
	record []string,
	lclient *logto.Client,
	logger *logrus.Entry,
	cltx *cli.Context,
	tcache *ttlcache.Cache[string, string],
) error {
	token, err := retrieveToken(&retrieveTokenProperties{
		Lclient: lclient,
		Logger:  logger,
		Cltx:    cltx,
		Tcache:  tcache,
	})
	if err != nil {
		return err
	}
	ok, _, err := lclient.CheckUser(
		token,
		record[0],
	)
	if err != nil {
		return err
	}
	if ok {
		logger.Infof(
			"user with email %s exist, skipping creation",
			record[0],
		)
		return nil
	}
	normUser := normalizeUserName(record[2], record[3])
	ok, _, err = lclient.CheckUserWithUserName(
		token,
		normUser,
	)
	if err != nil {
		return err
	}
	if ok {
		logger.Infof(
			"user with username %s exist, skipping creation",
			normUser,
		)
		return nil
	}
	logger.Debugf("username %s does not exist, going to create", record[0])
	userId, err := createUser(lclient, token, record, normUser)
	if err != nil {
		return fmt.Errorf("error in creating user %s %s", record[0], err)
	}
	logger.Infof("created user with email %s id %s\n", record[0], userId)
	err = addCustomUserInformation(lclient, token, userId, record)
	if err != nil {
		return err
	}
	logger.Debugf(
		"created custom data for user with email %s id %s\n",
		record[0],
		userId,
	)

	return nil
}
