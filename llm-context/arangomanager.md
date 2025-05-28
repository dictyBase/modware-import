This file is a merged representation of a subset of the codebase, containing specifically included files, combined into a single document by Repomix.
The content has been processed where comments have been removed.

# File Summary

## Purpose
This file contains a packed representation of the entire repository's contents.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Only files matching these patterns are included: **/*.go
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Code comments have been removed from supported file types
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
collection/
  collection.go
command/
  flag/
    flag.go
query/
  query_maps.go
  query_test.go
  query.go
testarango/
  doc.go
  testarango.go
database_test.go
database.go
datasource.go
result.go
resultset_test.go
resultset.go
session.go
statement.go
test_common.go
transaction_test.go
transaction.go
```

# Files

## File: testarango/doc.go
```go
package testarango
```

## File: testarango/testarango.go
```go
package testarango

import (
	"fmt"
	"os"
	"strconv"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/arangomanager"
)

const (
	aPort  = 8529
	minLen = 6
	maxLen = 8
)







func CheckArangoEnv() error {
	envs := []string{
		"ARANGO_USER",
		"ARANGO_HOST",
		"ARANGO_PASS",
	}
	for _, e := range envs {
		if len(os.Getenv(e)) == 0 {
			return fmt.Errorf("env %s is not set", e)
		}
	}

	return nil
}



type TestArango struct {
	*arangomanager.ConnectParams
	*arangomanager.Session
}











func NewTestArangoFromEnv(isCreate bool) (*TestArango, error) {
	tra := new(TestArango)
	if err := CheckArangoEnv(); err != nil {
		return tra, err
	}
	tra.ConnectParams = &arangomanager.ConnectParams{
		User: os.Getenv("ARANGO_USER"),
		Pass: os.Getenv("ARANGO_PASS"),
		Host: os.Getenv("ARANGO_HOST"),
		Port: aPort,
	}
	if len(os.Getenv("ARANGO_PORT")) > 0 {
		aport, _ := strconv.Atoi(os.Getenv("ARANGO_PORT"))
		tra.ConnectParams.Port = aport
	}
	sess, err := arangomanager.Connect(
		tra.ConnectParams.Host,
		tra.ConnectParams.User,
		tra.ConnectParams.Pass,
		tra.ConnectParams.Port,
		false,
	)
	if err != nil {
		return tra, fmt.Errorf("error in connecting %s", err)
	}
	tra.Session = sess
	if isCreate {
		if err := tra.CreateTestDb(arangomanager.RandomString(minLen, maxLen), &driver.CreateDatabaseOptions{}); err != nil {
			return tra, err
		}
	}

	return tra, nil
}





func NewTestArango(
	user, pass, host string,
	port int,
	isCreate bool,
) (*TestArango, error) {
	tra := new(TestArango)
	tra.ConnectParams = &arangomanager.ConnectParams{
		User: user,
		Pass: pass,
		Host: host,
		Port: port,
	}
	sess, err := arangomanager.Connect(
		tra.ConnectParams.Host,
		tra.ConnectParams.User,
		tra.ConnectParams.Pass,
		tra.ConnectParams.Port,
		false,
	)
	if err != nil {
		return tra, fmt.Errorf("error in connecting %s", err)
	}
	tra.Session = sess
	if isCreate {
		err = tra.CreateTestDb(
			arangomanager.RandomString(minLen, maxLen),
			&driver.CreateDatabaseOptions{},
		)
		if err != nil {
			return tra, err
		}
	}

	return tra, nil
}


func (ta *TestArango) CreateTestDb(
	name string,
	opt *driver.CreateDatabaseOptions,
) error {
	if err := ta.CreateDB(name, opt); err != nil {
		return fmt.Errorf("error in creating database %s", err)
	}
	ta.Database = name

	return nil
}
```

## File: datasource.go
```go
package arangomanager


type ConnectParams struct {
	User     string `validate:"required"`
	Pass     string `validate:"required"`
	Database string `validate:"required"`
	Host     string `validate:"required"`
	Port     int    `validate:"required"`
	Istls    bool
}
```

## File: result.go
```go
package arangomanager

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/fatih/structs"
)


type Result struct {
	cursor driver.Cursor
	empty  bool
}


func (r *Result) IsEmpty() bool {
	return r.empty
}


func (r *Result) Read(iface interface{}) error {
	meta, err := r.cursor.ReadDocument(context.TODO(), iface)
	if err != nil {
		return fmt.Errorf("error in reading document %s", err)
	}
	if !structs.IsStruct(iface) {
		return nil
	}
	s := structs.New(iface)
	if f, ok := s.FieldOk("DocumentMeta"); ok {
		if f.IsEmbedded() {
			if err := f.Set(meta); err != nil {
				return fmt.Errorf("error in assigning DocumentMeta to the structure %s", err)
			}
		}
	}

	return nil
}
```

## File: session.go
```go
package arangomanager

import (
	"context"
	"crypto/tls"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	validator "github.com/go-playground/validator/v10"
)


type Session struct {
	client driver.Client
}






func NewSessionFromClient(client driver.Client) *Session {
	return &Session{client}
}


func Connect(
	host, user, password string,
	port int,
	istls bool,
) (*Session, error) {
	connConf := http.ConnectionConfig{
		Endpoints: []string{
			fmt.Sprintf("http://%s:%d", host, port),
		},
	}
	if istls {
		connConf.Endpoints = []string{
			fmt.Sprintf("https://%s:%d", host, port),
		}
		connConf.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	conn, err := http.NewConnection(connConf)
	if err != nil {
		return &Session{}, fmt.Errorf("could not connect %s", err)
	}
	client, err := driver.NewClient(
		driver.ClientConfig{
			Connection: conn,
			Authentication: driver.BasicAuthentication(
				user,
				password,
			),
		})
	if err != nil {
		return &Session{}, fmt.Errorf("could not get a client instance %s", err)
	}

	return &Session{client}, nil
}



func NewSessionDb(connP *ConnectParams) (*Session, *Database, error) {
	var sess *Session
	var dbr *Database
	validate := validator.New()
	if err := validate.Struct(connP); err != nil {
		return sess, dbr, fmt.Errorf("error in validation %s", err)
	}
	sess, err := Connect(
		connP.Host,
		connP.User,
		connP.Pass,
		connP.Port,
		connP.Istls,
	)
	if err != nil {
		return sess, dbr, err
	}
	dbr, err = sess.DB(connP.Database)
	if err != nil {
		return sess, dbr, err
	}

	return sess, dbr, nil
}


func (s *Session) CurrentDB() (*Database, error) {
	return s.getDatabase("_system")
}


func (s *Session) CreateDB(
	name string,
	opt *driver.CreateDatabaseOptions,
) error {
	isOk, err := s.client.DatabaseExists(context.Background(), name)
	if err != nil {
		return fmt.Errorf(
			"error in checking existence of database %s %s",
			name,
			err,
		)
	}
	if !isOk {
		_, err = s.client.CreateDatabase(context.Background(), name, opt)
		if err != nil {
			return fmt.Errorf("error in creating database %s %s", name, err)
		}
	}

	return nil
}


func (s *Session) DB(name string) (*Database, error) {
	return s.getDatabase(name)
}


func (s *Session) CreateUser(user, pass string) error {
	ok, err := s.client.UserExists(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error in finding user %s", err)
	}
	if !ok {
		isActive := true
		_, err := s.client.CreateUser(
			context.Background(),
			user,
			&driver.UserOptions{Password: pass, Active: &isActive},
		)
		if err != nil {
			return fmt.Errorf("error in creating user %s", err)
		}
	}

	return nil
}


func (s *Session) GrantDB(database, user, grant string) error {
	ok, err := s.client.UserExists(context.Background(), user)
	if err != nil {
		return fmt.Errorf("error in finding user %s", err)
	}
	if !ok {
		return fmt.Errorf("user %s does not exist", user)
	}
	dbuser, err := s.client.User(context.Background(), user)
	if err != nil {
		return fmt.Errorf(
			"error in getting user %s from database %s",
			user,
			err,
		)
	}
	dbh, err := s.client.Database(context.Background(), database)
	if err != nil {
		return fmt.Errorf("cannot get a database instance %s", err)
	}
	err = dbuser.SetDatabaseAccess(context.Background(), dbh, getGrant(grant))
	if err != nil {
		return fmt.Errorf("error in setting database access %s", err)
	}

	return nil
}

func getGrant(g string) driver.Grant {
	var grnt driver.Grant
	switch g {
	case "rw":
		grnt = driver.GrantReadWrite
	case "ro":
		grnt = driver.GrantReadOnly
	default:
		grnt = driver.GrantNone
	}

	return grnt
}
func (s *Session) getDatabase(name string) (*Database, error) {
	isOk, err := s.client.DatabaseExists(context.Background(), name)
	if err != nil {
		return &Database{}, fmt.Errorf(
			"error in checking existing of database %s",
			err,
		)
	}
	if !isOk {
		return &Database{}, fmt.Errorf(
			"error in finding database named %s: %s",
			name,
			err,
		)
	}
	dbh, err := s.client.Database(context.Background(), name)
	if err != nil {
		return &Database{}, fmt.Errorf(
			"unable to get database instance %s",
			err,
		)
	}

	return &Database{dbh}, nil
}
```

## File: statement.go
```go
package arangomanager

const (
	truncateFn = `
		function(params) {
			const db = require('@arangodb').db
			params[0]
				.map(name => db._collection(name))
				.map(collection => collection.truncate())
		}
`
)
```

## File: collection/collection.go
```go
package collection

import (
	"cmp"
	"iter"
	"slices"
)



func Map[T1, T2 any](slc []T1, fnc func(T1) T2) []T2 {
	ret := make([]T2, 0)
	for _, elem := range slc {
		ret = append(ret, fnc(elem))
	}

	return ret
}




func CurriedMap[T1, T2 any](fnc func(T1) T2) func([]T1) []T2 {
	return func(slc []T1) []T2 {
		return Map(slc, fnc)
	}
}



func Include[T cmp.Ordered](slice []T, element T) bool {
	if !slices.IsSorted(slice) {
		slices.Sort(slice)
	}
	_, found := slices.BinarySearch(slice, element)

	return found
}



func RemoveStringItems(slice []string, items ...string) []string {
	str := make([]string, 0)
	for _, val := range slice {
		if !Include(items, val) {
			str = append(str, val)
		}
	}

	return str
}



func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}



func CurriedFilter[T any](predicate func(T) bool) func([]T) []T {
	return func(slice []T) []T {
		return Filter(slice, predicate)
	}
}



func MapSeq[T1, T2 any](seq iter.Seq[T1], fn func(T1) T2) iter.Seq[T2] {
	return func(yield func(T2) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}





func PartitionTuple2[T any](
	slice []T,
	predicate func(T) bool,
) Tuple2[[]T, []T] {
	trueSlice := make([]T, 0)
	falseSlice := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return NewTuple2(trueSlice, falseSlice)
}




func CurriedPartitionTuple2[T any](
	predicate func(T) bool,
) func([]T) Tuple2[[]T, []T] {
	return func(slice []T) Tuple2[[]T, []T] {
		return PartitionTuple2(slice, predicate)
	}
}





func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	trueSlice := make([]T, 0)
	falseSlice := make([]T, 0)
	for _, item := range slice {
		if predicate(item) {
			trueSlice = append(trueSlice, item)
		} else {
			falseSlice = append(falseSlice, item)
		}
	}

	return trueSlice, falseSlice
}




func CurriedPartition[T any](predicate func(T) bool) func([]T) ([]T, []T) {
	return func(slice []T) ([]T, []T) {
		return Partition(slice, predicate)
	}
}





func Pipe2[T1, T2, T3 any](tup T1, f1 func(T1) T2, f2 func(T2) T3) T3 {
	return f2(f1(tup))
}





func Pipe3[T1, T2, T3, T4 any](
	initial T1,
	f1 func(T1) T2,
	f2 func(T2) T3,
	f3 func(T3) T4,
) T4 {
	return f3(f2(f1(initial)))
}





func Pipe4[T1, T2, T3, T4, T5 any](
	initial T1,
	f1 func(T1) T2,
	f2 func(T2) T3,
	f3 func(T3) T4,
	f4 func(T4) T5,
) T5 {
	return f4(f3(f2(f1(initial))))
}



type Tuple2[T1, T2 any] struct {
	First  T1
	Second T2
}


func NewTuple2[T1, T2 any](first T1, second T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{
		First:  first,
		Second: second,
	}
}




func SliceToTuple2[T1, T2 any](slice []any) Tuple2[T1, T2] {
	var first T1
	var second T2
	if len(slice) > 0 {
		if val, ok := slice[0].(T1); ok {
			first = val
		}
	}
	if len(slice) > 1 {
		if val, ok := slice[1].(T2); ok {
			second = val
		}
	}

	return NewTuple2(first, second)
}




func TFold[A, B, R any](
	tup Tuple2[A, B],
	folder func(Tuple2[A, B]) R,
) R {
	return folder(tup)
}



func CurriedTFold[A, B, R any](
	folder func(Tuple2[A, B]) R,
) func(Tuple2[A, B]) R {
	return func(tup Tuple2[A, B]) R {
		return TFold(tup, folder)
	}
}



func IsEmpty[T any](slice []T) bool {
	return len(slice) == 0
}
```

## File: command/flag/flag.go
```go
package flag

import (
	"github.com/urfave/cli"
)

































func ArangoFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "arangodb-pass, pass",
			EnvVar:   "ARANGODB_PASS",
			Usage:    "arangodb database password",
			Required: true,
		},
		cli.StringFlag{
			Name:     "arangodb-user, user",
			EnvVar:   "ARANGODB_USER",
			Usage:    "arangodb database user",
			Required: true,
		},
		cli.StringFlag{
			Name:     "arangodb-host, host",
			Value:    "arangodb",
			EnvVar:   "ARANGODB_SERVICE_HOST",
			Usage:    "arangodb database host",
			Required: true,
		},
		cli.StringFlag{
			Name:   "arangodb-port",
			EnvVar: "ARANGODB_SERVICE_PORT",
			Usage:  "arangodb database port",
			Value:  "8529",
		},
		cli.BoolFlag{
			Name:  "is-secure",
			Usage: "flag for secured or unsecured arangodb endpoint",
		},
	}
}












func ArangodbFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "arangodb-pass, pass",
			EnvVar:   "ARANGODB_PASS",
			Usage:    "arangodb database password",
			Required: true,
		},
		cli.StringFlag{
			Name:     "arangodb-database, db",
			EnvVar:   "ARANGODB_DATABASE",
			Usage:    "arangodb database name",
			Required: true,
		},
		cli.StringFlag{
			Name:     "arangodb-user, user",
			EnvVar:   "ARANGODB_USER",
			Usage:    "arangodb database user",
			Required: true,
		},
		cli.StringFlag{
			Name:     "arangodb-host, host",
			Value:    "arangodb",
			EnvVar:   "ARANGODB_SERVICE_HOST",
			Usage:    "arangodb database host",
			Required: true,
		},
		cli.StringFlag{
			Name:   "arangodb-port",
			EnvVar: "ARANGODB_SERVICE_PORT",
			Usage:  "arangodb database port",
			Value:  "8529",
		},
		cli.BoolTFlag{
			Name:  "is-secure",
			Usage: "flag for secured or unsecured arangodb endpoint",
		},
	}
}
```

## File: query/query_maps.go
```go
package query

func getLogic(input string) string {
	lmap := map[string]string{",": "OR", ";": "AND"}

	return lmap[input]
}




func getOperatorMap() map[string]string {
	return map[string]string{
		"==":  "==",
		"===": "==",
		"!=":  "!=",
		">":   ">",
		"<":   "<",
		">=":  ">=",
		"<=":  "<=",
		"=~":  "=~",
		"!~":  "!~",
		"$==": "==",
		"$>":  ">",
		"$<":  "<",
		"$>=": ">=",
		"$<=": "<=",
		"@==": "==",
		"@=~": "=~",
		"@!~": "!~",
		"@!=": "!=",
	}
}

func getOperator(opt string) string {
	omap := getOperatorMap()

	return omap[opt]
}

func getArrayOpertaor(opt string) string {
	amap := getArrayOperatorMap()

	return amap[opt]
}

func hasOperator(opt string) bool {
	omap := getOperatorMap()
	_, isok := omap[opt]

	return isok
}

func hasDateOperator(opt string) bool {
	dmap := getDateOperatorMap()
	_, isok := dmap[opt]

	return isok
}

func hasArrayOperator(opt string) bool {
	amap := getArrayOperatorMap()
	_, isok := amap[opt]

	return isok
}




func getDateOperatorMap() map[string]string {
	return map[string]string{
		"$==": "==",
		"$>":  ">",
		"$<":  "<",
		"$>=": ">=",
		"$<=": "<=",
	}
}





func getArrayOperatorMap() map[string]string {
	return map[string]string{
		"@==": "==",
		"@=~": "=~",
		"@!~": "!~",
		"@!=": "!=",
	}
}
```

## File: query/query_test.go
```go
package query

import (
	"fmt"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/arangomanager/testarango"
	"github.com/stretchr/testify/require"
)

const (
	minLen = 20
	maxLen = 30
)


var fmap = map[string]string{
	"created_at": "created_at",
	"sport":      "sports",
	"email":      "email",
	"label":      "label",
	"tag":        "tag",
	"ontology":   "ontology",
	"summary":    "summary",
}

var qmap = map[string]string{
	"created_at": "foo.created_at",
	"sport":      "bar.game",
	"email":      "fizz.identifier",
	"label":      "v.label",
	"tag":        "s.tag",
	"ontology":   "cvterm.ontology",
	"summary":    "v.summary",
}

func setupTestArango(
	assert *require.Assertions,
) (*arangomanager.Database, string) {
	ta, err := testarango.NewTestArangoFromEnv(true)
	assert.NoError(
		err,
		"should not produce any error from testarango constructor",
	)
	dbh, err := ta.DB(ta.Database)
	assert.NoError(err, "should not produce any database error")
	crnd := arangomanager.RandomString(minLen, maxLen)
	_, err = dbh.CreateCollection(crnd, &driver.CreateCollectionOptions{})
	if err != nil {
		errDbh := dbh.Drop()
		assert.NoError(
			errDbh,
			"should not produce any error from database removal",
		)
	}

	return dbh, crnd
}

func cleanupAfterEach(assert *require.Assertions, dbh *arangomanager.Database) {
	err := dbh.Drop()
	assert.NoError(err, "should not produce any error from database removal")
}

func TestInvalidFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	filters := []*Filter{
		{Field: "sport", Operator: "+++", Value: "football", Logic: "AND"},
		{
			Field:    "email",
			Operator: "^^^",
			Value:    "gmail@gmail.com",
			Logic:    "AND",
		},
		{Field: "tag", Operator: "^^^", Value: "bozama"},
	}
	_, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: filters, Doc: "doc"},
	)
	assert.Error(err, "expect to have error with filter operator")
}

func TestParseFilterString(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	fls, err := ParseFilterString("sport===football;email===mahomes@gmail.com")
	assert.NoError(err, "should not return any parse error")
	assert.Len(fls, 2, "should match length of two items in filter array")
	assert.Equal(fls[0].Value, "football", "should match the sport query")
	assert.Equal(
		fls[1].Value,
		"mahomes@gmail.com",
		"should match the email query",
	)
	assert.Equal(fls[0].Field, "sport", "should match field sport")
	assert.Equal(fls[1].Field, "email", "should match fieldi email")
	assert.Equal(fls[0].Operator, "===", "should match equal operator")
	assert.Equal(fls[1].Operator, "===", "should match equal operator")
	assert.Equal(fls[0].Logic, ";", "should have parsed colon logic")
	assert.Empty(fls[1].Logic, "should have empty logic value")

	fls2, err := ParseFilterString("ontology!~dicty annotation;tag=~logicx")
	assert.NoError(err, "should not return any parse error")
	assert.Len(fls2, 2, "should have two items in filter array")
	assert.Equal(
		fls2[0].Value,
		"dicty annotation",
		"should match ontology query",
	)
	assert.Equal(fls2[1].Value, "logicx", "should match tag query")
	assert.Equal(fls2[0].Field, "ontology", "should match field ontology")
	assert.Equal(fls2[1].Field, "tag", "should match field annotation")
	assert.Equal(fls2[0].Operator, "!~", "should match regexp match operator")
	assert.Equal(
		fls2[1].Operator,
		"=~",
		"should match regexp negation operator",
	)
	assert.Equal(fls2[0].Logic, ";", "should have parsed colon logic")
	assert.Empty(fls2[1].Logic, "should have empty logic value")

	b, err := ParseFilterString("xyz")
	assert.NoError(err, "should not return any parse error")
	assert.Len(b, 0, "should have empty slice since regex doesn't match string")
}

func TestQualifiedMixedLogicStatement(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	fstr, err := ParseFilterString(
		"summary===bhokchoi;ontology===dicty_strain_property;tag===general strain,tag===REMI-seq",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt, err := GenQualifiedAQLFilterStatement(qmap, fstr)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		stmt,
		"( s.tag == 'general strain'\n OR s.tag == 'REMI-seq' )",
		"should have the expected substring",
	)
	assert.Contains(
		stmt,
		"v.summary == 'bhokchoi'\n AND cvterm.ontology == 'dicty_strain_property'",
		"should have the expected substring",
	)
	err = dbh.ValidateQ(genFullStmt(stmt, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	fstr2, err := ParseFilterString(
		"ontology===dicty_strain_property;tag===general strain,tag===REMI-seq;summary===bhokchoi",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt2, err := GenQualifiedAQLFilterStatement(qmap, fstr2)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		stmt2,
		"s.tag == 'REMI-seq' ) \n AND v.summary == 'bhokchoi'",
		"should have the expected substring",
	)
	err = dbh.ValidateQ(genFullStmt(stmt2, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	fstr3, err := ParseFilterString(
		"ontology===dicty_strain_property;tag===general strain,tag===REMI-seq,tag===bacterial strain;summary===bhokchoi",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt3, err := GenQualifiedAQLFilterStatement(qmap, fstr3)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		stmt3,
		"( s.tag == 'general strain'\n OR s.tag == 'REMI-seq'\n OR s.tag == 'bacterial strain' )",
		"should have the expected substring",
	)
	assert.Contains(
		stmt3,
		"cvterm.ontology == 'dicty_strain_property'\n AND  ( s.tag == 'general strain'",
		"should have the expected substring",
	)
	err = dbh.ValidateQ(genFullStmt(stmt3, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
	for _, stm := range []string{stmt, stmt2, stmt3} {
		assert.Contains(stm, "(", "should have starting parenthesis")
		assert.Contains(stm, ")", "should have ending parenthesis")
	}
}

func TestQualifiedEqualFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)

	f, err := ParseFilterString(
		"email===mahomes@gmail.com,email===brees@gmail.com",
	)
	assert.NoError(err, "should not return any parsing error")
	nqa, err := GenQualifiedAQLFilterStatement(qmap, f)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Equal(
		nqa,
		"FILTER  ( fizz.identifier == 'mahomes@gmail.com'\n OR fizz.identifier == 'brees@gmail.com' ) ",
		"should match filter statement",
	)
	err = dbh.ValidateQ(genFullQualifiedStmt(nqa, "fizz", cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestQualifiedSubstringFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)

	qf, err := ParseFilterString("label=~GWDI")
	assert.NoError(err, "should not return any parsing error")
	qsa, err := GenQualifiedAQLFilterStatement(qmap, qf)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Equal(
		qsa,
		"FILTER v.label =~ 'GWDI'",
		"should contain GWDI substring",
	)
	err = dbh.ValidateQ(genFullQualifiedStmt(qsa, "v", cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestQualifiedDateFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)

	df, err := ParseFilterString("created_at$==2019,created_at$==2018")
	assert.NoError(err, "should not return any parsing error")
	dfl, err := GenQualifiedAQLFilterStatement(qmap, df)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Equal(
		dfl,
		"FILTER  ( foo.created_at == DATE_ISO8601('2019')\n OR foo.created_at == DATE_ISO8601('2018') ) ",
	)
	err = dbh.ValidateQ(genFullQualifiedStmt(dfl, "foo", cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestQualifiedArrayFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)

	af, err := ParseFilterString("sport@==basketball")
	assert.NoError(err, "should not return any parsing error")
	afn, err := GenQualifiedAQLFilterStatement(qmap, af)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Contains(
		afn,
		"LET",
		"should contain LET term, indicating array item",
	)
	assert.Contains(
		afn,
		"FILTER 'basketball' IN bar.game[*]",
		"should contain an array containing statement",
	)
	err = dbh.ValidateQ(genFullQualifiedStmt(afn, "bar", cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	af2, err := ParseFilterString("sport@=~basket")
	assert.NoError(err, "should not return any parsing error")
	an2, err := GenQualifiedAQLFilterStatement(qmap, af2)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Contains(
		an2,
		"FILTER CONTAINS(x, LOWER('basket'))",
		"should contain FILTER CONTAINS statement, indicating array item substring",
	)
	assert.Contains(an2, "LIMIT 1", "should match limit of one")
	err = dbh.ValidateQ(genFullQualifiedStmt(an2, "bar", cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	bf, err := ParseFilterString("sport@!=banana,sport@!=apple")
	assert.NoError(err, "should not return any parsing error")
	bns, err := GenQualifiedAQLFilterStatement(qmap, bf)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Contains(
		bns,
		"NOT IN",
		"should contain NOT IN statement, indicating item is not in array",
	)
	assert.Contains(bns, "OR", "should contain OR term")
	assert.Contains(
		bns,
		"FILTER 'banana' NOT IN bar.game[*]",
		"should contain filter with NOT IN operator with collection and column name",
	)
	assert.Contains(
		bns,
		"FILTER 'apple' NOT IN bar.game[*]",
		"should contain filter with NOT IN operator with collection and column name",
	)
	err = dbh.ValidateQ(genFullQualifiedStmt(bns, "bar", cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestMixedLogicStatement(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	fstr, err := ParseFilterString(
		"summary===bhokchoi;ontology===dicty_strain_property;tag===general strain,tag===REMI-seq",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: fstr, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	err = dbh.ValidateQ(genFullStmt(stmt, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	fstr2, err := ParseFilterString(
		"ontology===dicty_strain_property;tag===general strain,tag===REMI-seq;summary===bhokchoi",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt2, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: fstr2, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	err = dbh.ValidateQ(genFullStmt(stmt2, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	fstr3, err := ParseFilterString(
		"ontology===dicty_strain_property;tag===general strain,tag===REMI-seq,tag===bacterial strain;summary===bhokchoi",
	)
	assert.NoError(err, "should not have any error from parsing string")
	stmt3, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: fstr3, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	err = dbh.ValidateQ(genFullStmt(stmt3, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
	for _, stm := range []string{stmt, stmt2, stmt3} {
		assert.Contains(stm, "(", "should have starting parenthesis")
		assert.Contains(stm, ")", "should have ending parenthesis")
	}
}

func TestAQLArrayFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	as, err := ParseFilterString("sport@==basketball")
	assert.NoError(err, "should not have any error from parsing string")
	afn, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: as, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		afn,
		"LET",
		"should contain LET term, indicating array item",
	)
	assert.Contains(
		afn,
		"FILTER 'basketball' IN doc.sports[*]",
		"should contain FILTER and IN term",
	)
	err = dbh.ValidateQ(genFullStmt(afn, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	a, err := ParseFilterString("sport@=~basket")
	assert.NoError(err, "should not have any error from parsing string")
	afl, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: a, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		afn,
		"LET",
		"should contain LET term, indicating array item",
	)
	assert.Contains(
		afl,
		"FILTER CONTAINS(x, LOWER('basket')) ",
		"should contain FILTER CONTAINS statement, indicating array item substring",
	)
	err = dbh.ValidateQ(genFullStmt(afl, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	b, err := ParseFilterString("sport@!=banana,sport@==apple")
	assert.NoError(err, "should not have any error from parsing string")
	bfa, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: b, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	moreFilterTests2(bfa, assert)
	err = dbh.ValidateQ(genFullStmt(bfa, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	b2, err := ParseFilterString("sport@=~banana;sport@==apple")
	assert.NoError(err, "should not have any error from parsing string")
	bf2, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: b2, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	moreFilterTests(bf2, assert)
	err = dbh.ValidateQ(genFullStmt(bfa, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func moreFilterTests2(bfa string, assert *require.Assertions) {
	assert.Contains(
		bfa,
		"FILTER 'apple' IN doc.sports[*]",
		"should contain IN statement",
	)
	assert.Contains(
		bfa,
		"FILTER 'banana' NOT IN doc.sports[*]",
		"should contain NOT IN statement",
	)
	assert.Contains(bfa, "OR", "should contain OR term")
}

func moreFilterTests(bf2 string, assert *require.Assertions) {
	assert.Contains(
		bf2,
		"FILTER 'apple' IN doc.sports[*]",
		"should contain IN statement",
	)
	assert.Contains(
		bf2,
		"FILTER CONTAINS(x, LOWER('banana'))",
		"should contain CONTAINS statement",
	)
	assert.Contains(bf2, "AND", "should contain AND logic")
}

func TestAQLDateFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	ds, err := ParseFilterString("created_at$==2019,created_at$>2018")
	assert.NoError(err, "should not have any error from parsing string")
	dfl, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: ds, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		dfl,
		"doc.created_at == DATE_ISO8601('2019')",
		"should contain DATE_ISO8601 term",
	)
	assert.Contains(
		dfl,
		"doc.created_at > DATE_ISO8601('2018')",
		"should contain DATE_ISO8601 term",
	)
	assert.Contains(dfl, "OR", "should contain OR term")
	err = dbh.ValidateQ(genFullStmt(dfl, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
	ds2, err := ParseFilterString(
		"created_at$<2019;created_at$<=2018;created_at$>=2020",
	)
	assert.NoError(err, "should not have any error from parsing string")
	dn2, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: ds2, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(
		dn2,
		"FILTER doc.created_at < DATE_ISO8601('2019')",
		"should contain DATE_ISO8601 term",
	)
	assert.Contains(
		dn2,
		"doc.created_at <= DATE_ISO8601('2018')",
		"should contain DATE_ISO8601 term",
	)
	assert.Contains(
		dn2,
		"doc.created_at >= DATE_ISO8601('2020')",
		"should contain DATE_ISO8601 term",
	)
	assert.Contains(dn2, "AND", "should contain AND term")
	err = dbh.ValidateQ(genFullStmt(dfl, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestAQLSubstringFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	qf, err := ParseFilterString("label=~GWDI")
	assert.NoError(err, "should not return any parsing error")
	qsa, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: qf, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Contains(qsa, "FILTER", "should contain FILTER term")
	assert.Contains(qsa, "doc.label =~ 'GWDI'", "should contain GWDI substring")
	err = dbh.ValidateQ(genFullStmt(qsa, cstr))
	assert.NoError(err, "should not have any invalid AQL query")

	qf2, err := ParseFilterString("label=~GWDI;email===brady@gmail.com")
	assert.NoError(err, "should not return any parsing error")
	qs2, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: qf2, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not return any error when generating AQL filter statement",
	)
	assert.Contains(qs2, "FILTER", "should contain FILTER term")
	assert.Contains(qs2, "doc.label =~ 'GWDI'", "should contain GWDI substring")
	assert.Contains(
		qs2,
		"doc.email == 'brady@gmail.com'",
		"should contain proper == statement",
	)
	err = dbh.ValidateQ(genFullStmt(qs2, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestAQLOperatorFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)

	s2, err := ParseFilterString(
		"email===mahomes@gmail.com;email===brees@gmail.com",
	)
	assert.NoError(err, "should not have any error from parsing string")
	na2, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: s2, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(na2, "FILTER", "should contain FILTER term")
	assert.Contains(
		na2,
		"doc.email == 'mahomes@gmail.com'",
		"should contain proper == statement",
	)
	assert.Contains(
		na2,
		"doc.email == 'brees@gmail.com'",
		"should contain proper == statement",
	)
	assert.Contains(na2, "AND", "should contain AND term")
	err = dbh.ValidateQ(
		genFullStmt(na2, cstr),
	)
	assert.NoError(err, "should not have any invalid AQL query")
}

func TestAQLEqualFilter(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	dbh, cstr := setupTestArango(assert)
	defer cleanupAfterEach(assert, dbh)
	s, err := ParseFilterString(
		"email===mahomes@gmail.com,email===brees@gmail.com",
	)
	assert.NoError(err, "should not have any error from parsing string")
	naf, err := GenAQLFilterStatement(
		&StatementParameters{Fmap: fmap, Filters: s, Doc: "doc"},
	)
	assert.NoError(
		err,
		"should not have any error from generating AQL filter statement",
	)
	assert.Contains(naf, "FILTER", "should contain FILTER term")
	assert.Contains(
		naf,
		"doc.email == 'mahomes@gmail.com'",
		"should contain proper == statement",
	)
	assert.Contains(
		naf,
		"doc.email == 'brees@gmail.com'",
		"should contain proper == statement",
	)
	assert.Contains(naf, "OR", "should contain OR term")
	err = dbh.ValidateQ(genFullStmt(naf, cstr))
	assert.NoError(err, "should not have any invalid AQL query")
}

func genFullQualifiedStmt(filter, name, coll string) string {
	return fmt.Sprintf(
		`
		FOR %s in %s
			%s
			RETURN %s
		`, name, coll, filter, name,
	)
}

func genFullStmt(filter, coll string) string {
	return fmt.Sprintf(
		`
		FOR doc in %s
			%s
			RETURN doc
		`,
		coll, filter,
	)
}

func TestGenQualifiedAQLFilterStatementFieldValidation(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	invalidFilters := []*Filter{
		{Field: "missing_field", Operator: "==", Value: "value"},
	}
	_, err := GenQualifiedAQLFilterStatement(qmap, invalidFilters)
	assert.Error(err, "should return error for missing field")
	assert.Contains(
		err.Error(),
		"missing field mappings in filter map",
		"error should mention missing field mappings",
	)
	assert.Contains(
		err.Error(),
		"missing_field",
		"error should contain the name of the missing field",
	)


	multipleInvalidFilters := []*Filter{
		{Field: "missing_field1", Operator: "==", Value: "value1"},
		{Field: "sport", Operator: "==", Value: "valid_field"},
		{Field: "missing_field2", Operator: "==", Value: "value2"},
	}
	_, err = GenQualifiedAQLFilterStatement(qmap, multipleInvalidFilters)
	assert.Error(err, "should return error for missing fields")
	assert.Contains(
		err.Error(),
		"missing field mappings in filter map",
		"error should mention missing field mappings",
	)
	assert.Contains(
		err.Error(),
		"missing_field1",
		"error should contain the first missing field",
	)
	assert.Contains(
		err.Error(),
		"missing_field2",
		"error should contain the second missing field",
	)


	validFilters := []*Filter{
		{Field: "sport", Operator: "==", Value: "basketball"},
		{Field: "email", Operator: "==", Value: "test@example.com"},
	}
	stmt, err := GenQualifiedAQLFilterStatement(qmap, validFilters)
	assert.NoError(err, "should not return error when all fields are valid")
	assert.NotEmpty(stmt, "should return a non-empty statement")
}
```

## File: resultset_test.go
```go
package arangomanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyResultsetScan(t *testing.T) {

	rs := &Resultset{empty: true}

	assert := require.New(t)
	assert.False(rs.Scan(), "Scan on empty resultset should return false")


	assert.False(
		rs.Scan(),
		"Multiple Scan calls on empty resultset should return false",
	)


	err := rs.Close()
	assert.NoError(err, "Close on empty resultset should not return error")
}

func TestEmptyResultsetClose(t *testing.T) {

	rs := &Resultset{empty: true}


	assert := require.New(t)
	err := rs.Close()
	assert.NoError(err, "Close on empty resultset should not return error")


	err = rs.Close()
	assert.NoError(
		err,
		"Multiple Close calls on empty resultset should not return error",
	)
}

func TestEmptyResultsetIsEmpty(t *testing.T) {

	rs := &Resultset{empty: true}


	assert := assert.New(t)
	assert.True(rs.IsEmpty(), "IsEmpty on empty resultset should return true")
}



func TestResultsetWorkflow(t *testing.T) {
	assert := assert.New(t)


	rs1 := &Resultset{empty: true}
	assert.True(rs1.IsEmpty(), "IsEmpty should return true for empty resultset")
	assert.False(rs1.Scan(), "Scan should return false for empty resultset")


	var data struct{}
	err := rs1.Read(&data)
	assert.Error(err, "Read on empty resultset should return error")


	err = rs1.Close()
	assert.NoError(err, "Close on empty resultset should not return error")



	rs2 := &Resultset{empty: true}
	assert.False(rs2.Scan(), "First scan should return false")
	assert.False(rs2.Scan(), "Second scan should also return false")
	assert.NoError(rs2.Close(), "First close should not error")
	assert.NoError(rs2.Close(), "Second close should not error")
}
```

## File: test_common.go
```go
package arangomanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/require"
)

const (
	genderQ = `
		FOR d IN @@collection
			FILTER d.gender == @gender
			RETURN d
	`
	genderQNoParam = `
		FOR d IN %s
			FILTER d.gender == '%s'
			RETURN d
	`
	userQ = `
		FOR d in @@collection
			FILTER d.name.first == @first
			FILTER d.name.last == @last
			RETURN d
	`
	userIns = `
		INSERT {
			name: {
				first: @first,
				last: @last
			},
			gender: @gender,
			contact: {
				region: @region,
				address: {
					city: @city,
					state: @state,
					zip: @zip
				}
			}
		} INTO %s
	`
	aPort  = 8529
	minLen = 10
	maxLen = 15
)


type DocParams struct {
	T         *testing.T
	TX        *TransactionHandler
	Coll      driver.Collection
	FirstName string
	LastName  string
}


type TxParams struct {
	T        *testing.T
	DB       *Database
	Coll     driver.Collection
	ReadOnly bool
}


type DocExistsParams struct {
	T           *testing.T
	DB          *Database
	Coll        driver.Collection
	FirstName   string
	LastName    string
	ShouldExist bool
}

func randomIntInRange(min, max int) (int, error) {
	if min >= max {
		return 0, fmt.Errorf("Invalid range")
	}

	possibleValues := big.NewInt(int64(max - min))

	randomValue, err := rand.Int(rand.Reader, possibleValues)
	if err != nil {
		return 0, err
	}

	return min + int(randomValue.Int64()), nil
}


func RandomInt(num int) (int, error) {
	randomValue, err := rand.Int(rand.Reader, big.NewInt(int64(num)))
	if err != nil {
		return 0, err
	}
	return int(randomValue.Int64()), nil
}

func FixedLenRandomString(length int) string {
	alphanum := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	byt := make([]byte, 0)
	alen := len(alphanum)
	for i := 0; i < length; i++ {
		pos, _ := RandomInt(alen)
		byt = append(byt, alphanum[pos])
	}

	return string(byt)
}


func RandomString(min, max int) string {
	alphanum := []byte("abcdefghijklmnopqrstuvwxyz")
	size, _ := randomIntInRange(min, max)
	byt := make([]byte, size)
	alen := len(alphanum)
	for i := 0; i < size; i++ {
		pos, _ := RandomInt(alen)
		byt[i] = alphanum[pos]
	}

	return string(byt)
}

type testArango struct {
	*ConnectParams
	*Session
}

type testUserDb struct {
	driver.DocumentMeta
	Birthday *time.Time `json:"birthday"`
	Contact  struct {
		Address struct {
			City   string `json:"city"`
			State  string `json:"state"`
			Street string `json:"street"`
			Zip    string `json:"zip"`
		} `json:"address"`
		Email  []string `json:"email"`
		Phone  []string `json:"phone"`
		Region string   `json:"region"`
	} `json:"contact"`
	Gender      string     `json:"gender"`
	Likes       []string   `json:"likes"`
	MemberSince *time.Time `json:"memberSince"`
	Name        struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
}

type testUser struct {
	driver.DocumentMeta
	Birthday *userDate `json:"birthday"`
	Contact  struct {
		Address struct {
			City   string `json:"city"`
			State  string `json:"state"`
			Street string `json:"street"`
			Zip    string `json:"zip"`
		} `json:"address"`
		Email  []string `json:"email"`
		Phone  []string `json:"phone"`
		Region string   `json:"region"`
	} `json:"contact"`
	Gender      string    `json:"gender"`
	Likes       []string  `json:"likes"`
	MemberSince *userDate `json:"memberSince"`
	Name        struct {
		First string `json:"first"`
		Last  string `json:"last"`
	} `json:"name"`
}

type userDate struct {
	time.Time
}

func (ud *userDate) UnmarshalJSON(in []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(in), `"`))
	if err != nil {
		return fmt.Errorf("error in parsing time %s", err)
	}
	ud.Time = t

	return nil
}

func checkArangoEnv() error {
	envs := []string{
		"ARANGO_USER",
		"ARANGO_HOST",
		"ARANGO_PASS",
	}
	for _, e := range envs {
		if len(os.Getenv(e)) == 0 {
			return fmt.Errorf("env %s is not set", e)
		}
	}

	return nil
}

func teardown(t *testing.T, c driver.Collection) {
	t.Helper()
	if err := c.Remove(context.Background()); err != nil {
		t.Fatalf("unable to truncate collection %s %s", c.Name(), err)
	}
}

func setup(t *testing.T, db *Database) driver.Collection {
	t.Helper()
	coll, err := db.FindOrCreateCollection(
		RandomString(minLen, maxLen),
		&driver.CreateCollectionOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err = loadTestData(coll); err != nil {
		t.Fatal(err)
	}

	return coll
}

func newTestArangoFromEnv(isCreate bool) (*testArango, error) {
	tra := new(testArango)
	if err := checkArangoEnv(); err != nil {
		return tra, err
	}
	tra.ConnectParams = &ConnectParams{
		User: os.Getenv("ARANGO_USER"),
		Pass: os.Getenv("ARANGO_PASS"),
		Host: os.Getenv("ARANGO_HOST"),
		Port: aPort,
	}
	if len(os.Getenv("ARANGO_PORT")) > 0 {
		aport, _ := strconv.Atoi(os.Getenv("ARANGO_PORT"))
		tra.ConnectParams.Port = aport
	}
	sess, err := Connect(
		tra.ConnectParams.Host,
		tra.ConnectParams.User,
		tra.ConnectParams.Pass,
		tra.ConnectParams.Port,
		false,
	)
	if err != nil {
		return tra, err
	}
	tra.Session = sess
	tra.Database = RandomString(minLen, maxLen)
	if isCreate {
		if err := sess.CreateDB(tra.Database, &driver.CreateDatabaseOptions{}); err != nil {
			return tra, err
		}
	}

	return tra, nil
}

func getReader() (io.Reader, error) {
	buff := bytes.NewBuffer(make([]byte, 0))
	dir, err := os.Getwd()
	if err != nil {
		return buff, fmt.Errorf("unable to get current dir %s", err)
	}
	fhr, err := os.Open(
		filepath.Join(
			dir, "testdata", "names.json",
		),
	)
	if err != nil {
		return fhr, fmt.Errorf("error in opening file %s", err)
	}

	return fhr, nil
}

func loadTestData(coll driver.Collection) error {
	reader, err := getReader()
	if err != nil {
		return err
	}
	dec := json.NewDecoder(reader)
	var ausr []*testUser
	for {
		var usr *testUser
		if err := dec.Decode(&usr); err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("error in decoding %s", err)
		}
		ausr = append(ausr, usr)
	}
	_, err = coll.ImportDocuments(
		context.Background(),
		ausr,
		&driver.ImportDocumentOptions{Complete: true},
	)
	if err != nil {
		return fmt.Errorf("error in importing document %s", err)
	}

	return nil
}


func setupTestTx(t *testing.T) (*Database, driver.Collection, func()) {
	t.Helper()

	ta, err := newTestArangoFromEnv(true)
	if err != nil {
		t.Fatalf("failed to create test database: %s", err)
	}


	cleanup := func() {

		dbh, _ := ta.Session.client.Database(context.Background(), ta.Database)
		if dbh != nil {
			if err := dbh.Remove(context.Background()); err != nil {
				t.Logf("failed to drop test database: %s", err)
			}
		}
	}

	db, err := ta.Session.DB(ta.Database)
	if err != nil {
		cleanup()
		t.Fatalf("failed to get database: %s", err)
	}


	coll := setup(t, db)

	return db, coll, cleanup
}


func beginTestTransaction(params TxParams) *TransactionHandler {
	params.T.Helper()

	opts := &TransactionOptions{}
	if params.ReadOnly {
		opts.ReadCollections = []string{params.Coll.Name()}
	} else {
		opts.WriteCollections = []string{params.Coll.Name()}
	}

	tx, err := params.DB.BeginTransaction(context.Background(), opts)
	if err != nil {
		params.T.Fatalf("failed to begin transaction: %s", err)
	}

	return tx
}


func assertTxCanceled(
	t *testing.T,
	tx *TransactionHandler,
	expectedCanceled bool,
) {
	t.Helper()
	assert := require.New(t)
	assert.Equal(expectedCanceled, tx.canceled,
		"Transaction canceled state mismatch, expected: %v, got: %v",
		expectedCanceled, tx.canceled)
}


func insertTestDocument(params DocParams) {
	params.T.Helper()
	assert := require.New(params.T)

	query := fmt.Sprintf(userIns, params.Coll.Name())
	bindVars := map[string]interface{}{
		"first":  params.FirstName,
		"last":   params.LastName,
		"gender": "male",
		"region": "test",
		"city":   "TestCity",
		"state":  "TestState",
		"zip":    "12345",
	}

	err := params.TX.Do(query, bindVars)
	assert.NoError(err)
}


func assertDocumentExists(params DocExistsParams) {
	params.T.Helper()
	assert := require.New(params.T)

	result, err := params.DB.GetRow(
		fmt.Sprintf(
			"FOR d IN %s FILTER d.name.first == @first AND d.name.last == @last RETURN d",
			params.Coll.Name(),
		),
		map[string]interface{}{
			"first": params.FirstName,
			"last":  params.LastName,
		},
	)
	assert.NoError(err)

	if params.ShouldExist {
		assert.False(
			result.IsEmpty(),
			"Document should exist but was not found: %s %s",
			params.FirstName,
			params.LastName,
		)
	} else {
		assert.True(result.IsEmpty(),
			"Document should not exist but was found: %s %s", params.FirstName, params.LastName)
	}
}
```

## File: transaction_test.go
```go
package arangomanager

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestTransactionHandlerMethods(t *testing.T) {

	db, coll, cleanup := setupTestTx(t)
	defer cleanup()
	defer teardown(t, coll)


	tx, err := db.BeginTransaction(context.Background(), &TransactionOptions{
		ReadCollections:  []string{coll.Name()},
		WriteCollections: []string{coll.Name()},
	})
	if err != nil {
		t.Fatalf("failed to begin transaction: %s", err)
	}

	t.Run("Context", func(t *testing.T) {
		assert := require.New(t)
		ctx := tx.Context()
		assert.NotNil(ctx)
	})

	t.Run("ID", func(t *testing.T) {
		assert := require.New(t)
		id := tx.ID()
		assert.NotEmpty(id)
		assert.Equal(tx.id, id)
	})

	t.Run("Status", func(t *testing.T) {
		assert := require.New(t)
		status, err := tx.Status()
		assert.NoError(err)
		assert.NotNil(status)
	})


	if err := tx.Abort(); err != nil {
		t.Fatalf("failed to abort transaction: %s", err)
	}
}


func TestTransactionLifecycle(t *testing.T) {

	db, coll, cleanup := setupTestTx(t)
	defer cleanup()
	defer teardown(t, coll)

	t.Run("Commit", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: false,
		})
		assert.NotNil(tx)
		assertTxCanceled(t, tx, false)


		err := tx.Commit()
		assert.NoError(err)
		assertTxCanceled(t, tx, true)


		err = tx.Commit()
		assert.Error(err)
		assert.Contains(err.Error(), "cannot commit a canceled transaction")
	})

	t.Run("Abort", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: false,
		})
		assert.NotNil(tx)
		assertTxCanceled(t, tx, false)


		err := tx.Abort()
		assert.NoError(err)
		assertTxCanceled(t, tx, true)


		err = tx.Abort()
		assert.Error(err)
		assert.Contains(err.Error(), "transaction already canceled")
	})
}


func TestTransactionDo(t *testing.T) {

	db, coll, cleanup := setupTestTx(t)
	defer cleanup()
	defer teardown(t, coll)

	t.Run("DoSuccessful", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: false,
		})


		insertTestDocument(DocParams{
			T:         t,
			TX:        tx,
			Coll:      coll,
			FirstName: "TestUser",
			LastName:  "DoMethod",
		})


		err := tx.Commit()
		assert.NoError(err)


		assertDocumentExists(DocExistsParams{
			T:           t,
			DB:          db,
			Coll:        coll,
			FirstName:   "TestUser",
			LastName:    "DoMethod",
			ShouldExist: true,
		})
	})

	t.Run("DoWithInvalidQuery", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: false,
		})


		err := tx.Do("INVALID QUERY", nil)
		assert.Error(err)


		err = tx.Abort()
		assert.NoError(err)
	})
}


func TestTransactionDoRun(t *testing.T) {

	db, coll, cleanup := setupTestTx(t)
	defer cleanup()
	defer teardown(t, coll)

	t.Run("DoRunSuccessful", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: true,
		})


		query := fmt.Sprintf(
			"FOR d IN %s FILTER d.gender == @gender RETURN d",
			coll.Name(),
		)
		bindVars := map[string]interface{}{"gender": "male"}

		result, err := tx.DoRun(query, bindVars)
		assert.NoError(err)
		assert.NotNil(result)
		assert.False(result.IsEmpty())


		err = tx.Abort()
		assert.NoError(err)
	})

	t.Run("DoRunWithInvalidQuery", func(t *testing.T) {
		assert := require.New(t)

		tx := beginTestTransaction(TxParams{
			T:        t,
			DB:       db,
			Coll:     coll,
			ReadOnly: true,
		})


		result, err := tx.DoRun("INVALID QUERY", nil)
		assert.Error(err)
		assert.True(result.IsEmpty())


		err = tx.Abort()
		assert.NoError(err)
	})
}


func TestTransactionIsolation(t *testing.T) {
	assert := require.New(t)

	db, coll, cleanup := setupTestTx(t)
	defer cleanup()
	defer teardown(t, coll)


	tx := beginTestTransaction(TxParams{
		T:        t,
		DB:       db,
		Coll:     coll,
		ReadOnly: false,
	})


	insertTestDocument(DocParams{
		T:         t,
		TX:        tx,
		Coll:      coll,
		FirstName: "Isolation",
		LastName:  "Test",
	})


	assertDocumentExists(DocExistsParams{
		T:           t,
		DB:          db,
		Coll:        coll,
		FirstName:   "Isolation",
		LastName:    "Test",
		ShouldExist: false,
	})


	err := tx.Commit()
	assert.NoError(err)
	assertTxCanceled(t, tx, true)


	assertDocumentExists(DocExistsParams{
		T:           t,
		DB:          db,
		Coll:        coll,
		FirstName:   "Isolation",
		LastName:    "Test",
		ShouldExist: true,
	})
}
```

## File: transaction.go
```go
package arangomanager

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
)


type TransactionHandler struct {
	db       *Database
	id       driver.TransactionID
	ctx      context.Context
	canceled bool
}


func (t *TransactionHandler) Context() context.Context {
	return t.ctx
}


func (t *TransactionHandler) ID() driver.TransactionID {
	return t.id
}


func (t *TransactionHandler) Commit() error {
	if t.canceled {
		return fmt.Errorf("cannot commit a canceled transaction")
	}

	if err := t.db.dbh.CommitTransaction(context.Background(), t.id, nil); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	t.canceled = true
	return nil
}


func (t *TransactionHandler) Abort() error {
	if t.canceled {
		return fmt.Errorf("transaction already canceled")
	}

	if err := t.db.dbh.AbortTransaction(context.Background(), t.id, nil); err != nil {
		return fmt.Errorf("failed to abort transaction: %w", err)
	}

	t.canceled = true
	return nil
}


func (t *TransactionHandler) Status() (driver.TransactionStatusRecord, error) {
	status, err := t.db.dbh.TransactionStatus(context.Background(), t.id)
	if err != nil {
		return driver.TransactionStatusRecord{}, fmt.Errorf(
			"failed to get transaction status: %w",
			err,
		)
	}

	return status, nil
}


func (t *TransactionHandler) Do(
	query string,
	bindVars map[string]interface{},
) error {
	ctx := driver.WithSilent(t.ctx)
	_, err := t.db.dbh.Query(ctx, query, bindVars)
	if err != nil {
		return fmt.Errorf("error in data modification query %w", err)
	}

	return nil
}



func (t *TransactionHandler) DoRun(
	query string,
	bindVars map[string]interface{},
) (*Result, error) {
	if err := t.db.dbh.ValidateQuery(t.ctx, query); err != nil {
		return &Result{
				empty: true,
			}, fmt.Errorf(
				"error in validating the query %s",
				err,
			)
	}
	cqr, err := t.db.dbh.Query(t.ctx, query, bindVars)
	return t.db.getResult(cqr, err)
}
```

## File: database_test.go
```go
package arangomanager

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	driver "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/require"
)

var (
	ahost, aport, auser, apass, adb string
	adbh                            *Database
)

type genderCountParams struct {
	require    *require.Assertions
	collection driver.Collection
	gender     string
	count      int64
}

func TestMain(m *testing.M) {
	tra, err := newTestArangoFromEnv(true)
	if err != nil {
		log.Fatalf("unable to construct new TestArango instance %s", err)
	}
	dbh, err := tra.DB(tra.Database)
	if err != nil {
		log.Fatalf("unable to get database %s", err)
	}
	auser = tra.User
	apass = tra.Pass
	ahost = tra.Host
	aport = strconv.Itoa(tra.Port)
	adb = tra.Database
	adbh = dbh
	code := m.Run()
	if err := dbh.Drop(); err != nil {
		log.Fatalf("error in dropping database %s", err)
	}
	os.Exit(code)
}

func TestCount(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	fc, err := adbh.Count(fmt.Sprintf(genderQNoParam, conn.Name(), "female"))
	require := require.New(t)
	require.NoErrorf(
		err,
		"expect no error from counting query, received error %s",
		err,
	)
	require.Equalf(fc, int64(15), "expect %d received %d", 15, fc)
	mc, err := adbh.Count(fmt.Sprintf(genderQNoParam, conn.Name(), "male"))
	require.NoErrorf(
		err,
		"expect no error from counting query, received error %s",
		err,
	)
	require.Equalf(mc, int64(15), "expect %d received %d", 15, mc)
}

func TestCountWithParams(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	require := require.New(t)
	for _, g := range []string{"male", "female"} {
		testGenderCount(&genderCountParams{
			require:    require,
			collection: conn,
			gender:     g,
			count:      int64(15),
		})
	}
}

func TestCollection(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	_, err := adbh.Collection(RandomString(6, 8))
	require := require.New(t)
	require.Error(
		err,
		"expect to return an error for an non-existent collection",
	)
	nc, err := adbh.Collection(conn.Name())
	require.NoError(
		err,
		"not expect to return an error for existent collection",
	)
	require.Equalf(
		conn.Name(),
		nc.Name(),
		"expect %s, received %s",
		conn.Name(),
		nc.Name(),
	)
}

func TestCreateCollection(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	_, err := adbh.CreateCollection(conn.Name(), nil)
	require := require.New(t)
	require.Error(err, "expect to return existing collection error")
	ncoll := RandomString(9, 11)
	nc, err := adbh.CreateCollection(ncoll, nil)
	require.NoError(
		err,
		"not expect to return an error for non-existent collection",
	)
	require.Equalf(
		ncoll,
		nc.Name(),
		"expect %s, received %s",
		"bogus",
		nc.Name(),
	)
}

func TestFindOrCreateCollection(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	ec, err := adbh.FindOrCreateCollection(c.Name(), nil)
	require := require.New(t)
	require.NoError(
		err,
		"not expect to return an error for existent collection",
	)
	require.Equalf(
		c.Name(),
		ec.Name(),
		"expect %s, received %s",
		c.Name(),
		ec.Name(),
	)
	ncoll := RandomString(12, 15)
	nc, err := adbh.FindOrCreateCollection(ncoll, nil)
	require.NoError(
		err,
		"not expect to return an error for existent collection",
	)
	require.Equalf(
		ncoll,
		nc.Name(),
		"expect %s, received %s",
		"bogus",
		nc.Name(),
	)
}

func TestEnsureGeoIndex(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	require := require.New(t)
	name := "value"
	index, b, err := adbh.EnsureGeoIndex(
		c.Name(),
		[]string{name},
		&driver.EnsureGeoIndexOptions{
			Name: name,
		},
	)
	require.NoError(err, "should not return error for geo index method")
	require.True(b, "should create geo index")
	require.Exactly(
		index.Type(),
		driver.GeoIndex,
		"should return geo index type",
	)
	require.Exactly(index.UserName(), name, "should match provided name option")
	_, _, err = adbh.EnsureGeoIndex(
		"wrong name",
		[]string{name},
		&driver.EnsureGeoIndexOptions{},
	)
	require.Error(err, "should return error for wrong collection name")
}

func TestEnsureHashIndex(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	require := require.New(t)
	name := "entry_id"
	index, b, err := adbh.EnsureHashIndex(
		c.Name(),
		[]string{name},
		&driver.EnsureHashIndexOptions{
			Name: name,
		},
	)
	require.NoError(err, "should not return error for hash index method")
	require.True(b, "should create hash index")
	require.Exactly(
		index.Type(),
		driver.HashIndex,
		"should return hash index type",
	)
	require.Exactly(index.UserName(), name, "should match provided name option")
	_, _, err = adbh.EnsureHashIndex(
		"wrong name",
		[]string{name},
		&driver.EnsureHashIndexOptions{},
	)
	require.Error(err, "should return error for wrong collection name")
}

func TestEnsurePersistentIndex(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	require := require.New(t)
	name := "entry_id"
	index, b, err := adbh.EnsurePersistentIndex(
		c.Name(),
		[]string{name},
		&driver.EnsurePersistentIndexOptions{
			Name: name,
		},
	)
	require.NoError(err, "should not return error for index method")
	require.True(b, "should create index")
	require.Exactly(
		index.Type(),
		driver.PersistentIndex,
		"should return persistent index type",
	)
	require.Exactly(index.UserName(), name, "should match provided name option")
	_, _, err = adbh.EnsurePersistentIndex(
		"wrong name",
		[]string{name},
		&driver.EnsurePersistentIndexOptions{},
	)
	require.Error(err, "should return error for wrong collection name")
}

func TestEnsureSkipListIndex(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	require := require.New(t)
	name := "created_at"
	index, b, err := adbh.EnsureSkipListIndex(
		c.Name(),
		[]string{name},
		&driver.EnsureSkipListIndexOptions{
			Name: name,
		},
	)
	require.NoError(err, "should not return error for skip list index method")
	require.True(b, "should create skip list index")
	require.Exactly(
		index.Type(),
		driver.SkipListIndex,
		"should return skip list index type",
	)
	require.Exactly(index.UserName(), name, "should match provided name option")
	_, _, err = adbh.EnsureSkipListIndex(
		"wrong name",
		[]string{name},
		&driver.EnsureSkipListIndexOptions{},
	)
	require.Error(err, "should return error for wrong collection name")
}

func TestSearchRowsWithParams(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	frs, err := adbh.SearchRows(
		genderQ,
		map[string]interface{}{
			"@collection": conn.Name(),
			"gender":      "female",
		},
	)
	testSearchRs(t, frs, err)
	wrs, err := adbh.SearchRows(
		genderQ,
		map[string]interface{}{
			"@collection": conn.Name(),
			"gender":      "wakanda",
		},
	)
	testSearchRsNoRow(t, wrs, err)
}

func TestSearchRows(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	frs, err := adbh.Search(fmt.Sprintf(genderQNoParam, c.Name(), "female"))
	testSearchRs(t, frs, err)
	wrs, err := adbh.Search(fmt.Sprintf(genderQNoParam, c.Name(), "wakanda"))
	testSearchRsNoRow(t, wrs, err)
}

func TestDo(t *testing.T) {
	t.Parallel()
	c := setup(t, adbh)
	defer teardown(t, c)
	err := adbh.Do(
		fmt.Sprintf(userIns, c.Name()),
		map[string]interface{}{
			"first":  "Chitkini",
			"last":   "Dey",
			"gender": "male",
			"region": "gram",
			"city":   "porgona",
			"state":  "wb",
			"zip":    "48943",
		},
	)
	require := require.New(t)
	require.NoErrorf(
		err,
		"expect no error from insert query, received error %s",
		err,
	)
}

func TestGetRow(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	row, err := adbh.GetRow(
		userQ,
		map[string]interface{}{
			"@collection": conn.Name(),
			"first":       "Mickie",
			"last":        "Menchaca",
		},
	)
	require := require.New(t)
	require.NoErrorf(
		err,
		"expect no error from search query, received error %s",
		err,
	)
	require.False(row.IsEmpty(), "expect result to be not empty")
	var u testUserDb
	err = row.Read(&u)
	require.NoError(err, "expect no error from reading the data")
	require.Equal(u.Gender, "female", "expect gender to be female")
	require.Equal(
		u.Contact.Address.City,
		"Beachwood",
		"should match city Beachwood",
	)
	require.Equal(u.Contact.Region, "732", "should match region 732")
	erow, err := adbh.GetRow(
		userQ,
		map[string]interface{}{
			"@collection": conn.Name(),
			"first":       "Pantu",
			"last":        "Boka",
		},
	)
	require.NoErrorf(
		err,
		"expect no error from row query, received error %s",
		err,
	)
	require.True(erow.IsEmpty(), "expect empty resultset")
}

func TestTruncate(t *testing.T) {
	t.Parallel()
	conn := setup(t, adbh)
	defer teardown(t, conn)
	require := require.New(t)
	for _, g := range []string{"male", "female"} {
		testGenderCount(&genderCountParams{
			require:    require,
			collection: conn,
			gender:     g,
			count:      int64(15),
		})
	}
	err := adbh.Truncate(conn.Name())
	require.NoErrorf(
		err,
		"expect no error from truncation, received error %s",
		err,
	)
	for _, g := range []string{"male", "female"} {
		testGenderCount(&genderCountParams{
			require:    require,
			collection: conn,
			gender:     g,
			count:      int64(0),
		})
	}
}

func testGenderCount(args *genderCountParams) {
	gcp, err := adbh.CountWithParams(
		genderQ,
		map[string]interface{}{
			"@collection": args.collection.Name(),
			"gender":      args.gender,
		},
	)
	args.require.NoErrorf(
		err,
		"expect no error from counting query, received error %s",
		err,
	)
	args.require.Equalf(
		gcp,
		args.count,
		"expect %d received %d",
		args.count,
		gcp,
	)
}

func testAllRows(rs *Resultset, require *require.Assertions, count int) {
	for i := 0; i < count; i++ {
		require.True(rs.Scan(), "expect scanning of record")
		var u testUserDb
		err := rs.Read(&u)
		require.NoError(err, "expect no error from reading the data")
		require.Equal(u.Gender, "female", "expect gender to be female")
	}
}

func testSearchRs(t *testing.T, rs *Resultset, err error) {
	t.Helper()
	require := require.New(t)
	require.NoErrorf(
		err,
		"expect no error from search query, received error %s",
		err,
	)
	require.False(rs.IsEmpty(), "expect resultset to be not empty")
	testAllRows(rs, require, 15)
	require.False(rs.Scan(), "should be false")
	require.NoError(rs.Close(), "should not return error")
}

func testSearchRsNoRow(t *testing.T, rs *Resultset, err error) {
	t.Helper()
	require := require.New(t)
	require.NoErrorf(
		err,
		"expect no error from search query, received error %s",
		err,
	)
	require.True(rs.IsEmpty(), "expect empty resultset")

	require.False(rs.Scan(), "scan on empty resultset should return false")
	require.NoError(rs.Close(), "close on empty resultset should not error")
}
```

## File: query/query.go
```go
package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dictyBase/arangomanager"
	"github.com/dictyBase/arangomanager/collection"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/now"
)

const (
	logicIdx         = 2
	charSet          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	filterStrLen     = 5
	strSeedLen       = 10
	arrQualMatchTmpl = `
		LET %s = (
			FOR x IN %s[*]
				FILTER CONTAINS(x, LOWER('%s'))
				LIMIT 1
				RETURN 1
		)
	`
	arrMatchTmpl = `
	      LET %s = (
	 		FOR x IN %s.%s[*]
				FILTER CONTAINS(x, LOWER('%s'))
				LIMIT 1
				RETURN 1
		)
	`
	arrQualEqualTmpl = `
		LET %s = (
			FILTER '%s' IN %s[*]
			RETURN 1
		)
	`
	arrEqualTmpl = `
		LET %s = (
				FILTER '%s' IN %s.%s[*]
				RETURN 1
		)
	`
	arrNotEqualTmpl = `
		LET %s = (
				FILTER '%s' NOT IN %s.%s[*]
				RETURN 1
		)
	`
	arrQualNotEqualTmpl = `
		LET %s = (
				FILTER '%s' NOT IN %s[*]
				RETURN 1
		)
	`
	dateTmpl = "%s.%s %s DATE_ISO8601('%s')"
)

var (
	startPrefixRegxp = regexp.MustCompile(`\(`)
	endPrefixRegxp   = regexp.MustCompile(`\)`)
	validate         = validator.New()
)

func init() {

	_ = validate.RegisterValidation("operator_validation", validateOperator)
}


func validateOperator(fl validator.FieldLevel) bool {
	validOperators := map[string]bool{
		"==": true, "!=": true, "===": true, "!==": true,
		"=~": true, "!~": true, ">": true, "<": true,
		">=": true, "<=": true, "$==": true, "$>": true,
		"$>=": true, "$<": true, "$<=": true, "@==": true,
		"@!=": true, "@!~": true, "@=~": true,
	}

	return validOperators[fl.Field().String()]
}


type AQLFilterParams struct {

	Fmap map[string]string `validate:"required,min=1"`

	Filters []*Filter `validate:"required,min=1,dive"`
}


type Filter struct {

	Field string `validate:"required"`

	Operator string `validate:"required,operator_validation"`

	Value string `validate:"required"`

	Logic string `validate:"omitempty"`
}


type StatementParameters struct {

	Fmap map[string]string `validate:"required"`

	Filters []*Filter `validate:"required,dive"`

	Doc string `validate:"required"`

	Vert string
}


func buildFilter() (*regexp.Regexp, error) {
	var bldr strings.Builder
	bldr.WriteString(`(\w+)`)
	bldr.WriteString(`(\=\=|\!\=|\=\=\=|\!\=\=|`)
	bldr.WriteString(`\=\~|\!\~|>|<|>\=|`)
	bldr.WriteString(`\=<|\$\=\=|\$\>|`)
	bldr.WriteString(`\$\>\=|\$\<|\$\<\=|`)
	bldr.WriteString(`\@\=\=|\@\!\=|`)
	bldr.WriteString(`\@\!\~|\@\=\~)`)
	bldr.WriteString(`([\w-@.\s]+)(\,|\;)?`)
	rgxp, err := regexp.Compile(bldr.String())
	if err != nil {
		return rgxp, fmt.Errorf("error in compiling regexp %s", err)
	}

	return rgxp, nil
}



func buildDate() (*regexp.Regexp, error) {
	var bldr strings.Builder
	bldr.WriteString(`^\d{4}\-(0[1-9]|1[012])$|`)
	bldr.WriteString(`^\d{4}$|^\d{4}\-(0[1-9]|`)
	bldr.WriteString(`1[012])\-(0[1-9]|[12][0-9]|3[01])$`)
	rgxp, err := regexp.Compile(bldr.String())
	if err != nil {
		return rgxp, fmt.Errorf("error in compiling regexp %s", err)
	}

	return rgxp, nil
}










func ParseFilterString(fstr string) ([]*Filter, error) {
	filters := make([]*Filter, 0)
	qre, err := buildFilter()
	if err != nil {
		return filters, err
	}
	m := qre.FindAllStringSubmatch(fstr, -1)
	if len(m) == 0 {
		return filters, nil
	}
	omap := getOperatorMap()
	for _, mtc := range m {

		if _, ok := omap[mtc[2]]; !ok {
			return filters, fmt.Errorf("filter operator %s not allowed", mtc[2])
		}
		flt := &Filter{
			Field:    mtc[1],
			Operator: mtc[2],
			Value:    mtc[3],
		}
		if len(mtc) == filterStrLen {
			flt.Logic = mtc[4]
		}
		filters = append(filters, flt)
	}

	return filters, nil
}



func validateFilterFields(fmap map[string]string, filters []*Filter) error {

	missingFields := collection.Filter(filters, func(f *Filter) bool {
		_, exists := fmap[f.Field]
		return !exists
	})

	if len(missingFields) > 0 {
		missingFieldNames := collection.Map(
			missingFields,
			func(f *Filter) string {
				return f.Field
			},
		)
		return fmt.Errorf(
			"missing field mappings in filter map: %v",
			missingFieldNames,
		)
	}

	return nil
}


func handleQualifiedArrayFilter(
	stmts map[string]*arraylist.List,
	flt *Filter,
	fmap map[string]string,
) {
	randStr := arangomanager.FixedLenRandomString(strSeedLen)
	switch getArrayOpertaor(flt.Operator) {
	case "=~":
		stmts["let"].Insert(
			0,
			fmt.Sprintf(
				arrQualMatchTmpl,
				randStr,
				fmap[flt.Field],
				flt.Value,
			),
		)
	case "==":
		stmts["let"].Insert(
			0,
			fmt.Sprintf(
				arrQualEqualTmpl,
				randStr,
				flt.Value,
				fmap[flt.Field],
			),
		)
	case "!=":
		stmts["let"].Insert(
			0,
			fmt.Sprintf(
				arrQualNotEqualTmpl,
				randStr,
				flt.Value,
				fmap[flt.Field],
			))
	}
	stmts["nonlet"].Add(fmt.Sprintf("LENGTH(%s) > 0", randStr))
}

















func GenQualifiedAQLFilterStatement(
	fmap map[string]string,
	filters []*Filter,
) (string, error) {
	if err := validate.Struct(&AQLFilterParams{
		Fmap:    fmap,
		Filters: filters,
	}); err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	if err := validateFilterFields(fmap, filters); err != nil {
		return "", err
	}

	stmts := map[string]*arraylist.List{
		"let":    arraylist.New(),
		"nonlet": arraylist.New(),
	}


	for _, flt := range filters {
		switch {
		case hasArrayOperator(flt.Operator):
			handleQualifiedArrayFilter(stmts, flt, fmap)
		case hasDateOperator(flt.Operator):
			if err := dateValidator(flt.Value); err != nil {
				return "", err
			}
			// write time conversion into AQL query
			stmts["nonlet"].Add(fmt.Sprintf("%s %s DATE_ISO8601('%s')",
				fmap[flt.Field], getOperator(flt.Operator), flt.Value,
			))
		case hasOperator(flt.Operator):

			stmts["nonlet"].Add(fmt.Sprintf(
				"%s %s %s",
				fmap[flt.Field], getOperator(flt.Operator),
				addQuoteToStrings(flt.Operator, flt.Value),
			))
		default:
			return "", fmt.Errorf(
				"unknown opertaor for parsing %s",
				flt.Operator,
			)
		}

		addLogic(stmts["nonlet"], flt)
	}

	return toFullStatement(stmts), nil
}

func handleArrayOpertaor(
	prms *StatementParameters,
	flt *Filter,
	randStr string,
) string {
	inner := prms.Doc
	var stmt string
	switch getArrayOpertaor(flt.Operator) {
	case "=~":
		stmt = fmt.Sprintf(
			arrMatchTmpl,
			randStr,
			inner,
			prms.Fmap[flt.Field],
			flt.Value,
		)
	case "==":
		stmt = fmt.Sprintf(
			arrEqualTmpl,
			randStr,
			flt.Value,
			inner,
			prms.Fmap[flt.Field],
		)
	case "!=":
		stmt = fmt.Sprintf(
			arrNotEqualTmpl,
			randStr,
			flt.Value,
			inner,
			prms.Fmap[flt.Field],
		)
	}
	return stmt
}

















func GenAQLFilterStatement(prms *StatementParameters) (string, error) {
	if err := validate.Struct(prms); err != nil {
		return "", fmt.Errorf(
			"validation error in StatementParameters: %w",
			err,
		)
	}

	inner := prms.Doc
	stmts := arraylist.New()
	if len(prms.Vert) > 0 {
		inner = prms.Vert
	}
	for _, flt := range prms.Filters {
		switch {
		case hasArrayOperator(flt.Operator):
			randStr := arangomanager.FixedLenRandomString(strSeedLen)
			stmts.Insert(0, handleArrayOpertaor(prms, flt, randStr))
			stmts.Add(fmt.Sprintf("LENGTH(%s) > 0", randStr))
		case hasDateOperator(flt.Operator):
			if err := dateValidator(flt.Value); err != nil {
				return "", err
			}
			// write time conversion into AQL query
			stmts.Add(
				fmt.Sprintf(
					dateTmpl, inner, prms.Fmap[flt.Field],
					getOperator(flt.Operator), flt.Value,
				),
			)
		case hasOperator(flt.Operator):
			// write the rest of AQL statement based on regular string data
			stmts.Add(fmt.Sprintf(
				"%s.%s %s %s",
				inner,
				prms.Fmap[flt.Field],
				getOperator(
					flt.Operator,
				),
				addQuoteToStrings(flt.Operator, flt.Value),
			))
		default:
			return "", fmt.Errorf(
				"unknown opertaor for parsing %s",
				flt.Operator,
			)
		}
		addLogic(stmts, flt)
	}

	return toString(stmts), nil
}

func addLogic(stmts *arraylist.List, flt *Filter) {
	currSize := stmts.Size()
	if len(flt.Logic) == 0 {
		addClosingParen(stmts, currSize)

		return
	}
	logic := getLogic(flt.Logic)
	switch logic {
	case "OR":
		addStartingParen(stmts, currSize)
	case "AND":
		addClosingParen(stmts, currSize)
	}
	stmts.Add(fmt.Sprintf("\n %s ", logic))
}

func addStartingParen(stmts *arraylist.List, currSize int) {
	if !isBalancedParens(stmts) {
		return
	}
	stmts.Insert(currSize-1, " ( ")
}

func addClosingParen(stmts *arraylist.List, currSize int) {
	if isBalancedParens(stmts) {
		return
	}
	elem, _ := stmts.Get(currSize - logicIdx)
	if val, ok := elem.(string); ok {
		if strings.TrimSpace(val) == "OR" {
			stmts.Add(" ) ")
		}
	}
}

func isBalancedParens(stmts *arraylist.List) bool {
	strStmt := stmts.String()
	startLen := len(startPrefixRegxp.FindAllString(strStmt, -1))
	endLen := len(endPrefixRegxp.FindAllString(strStmt, -1))

	return startLen == endLen
}

func toFullStatement(mst map[string]*arraylist.List) string {
	var clause strings.Builder

	if letList, ok := mst["let"]; ok {
		itr := letList.Iterator()
		for itr.Next() {
			clause.WriteString(itr.Value().(string))
		}
	}
	clause.WriteString("FILTER ")
	if nonletList, ok := mst["nonlet"]; ok {
		itr := nonletList.Iterator()
		for itr.Next() {
			clause.WriteString(itr.Value().(string))
		}
	}

	return clause.String()
}

func toString(l *arraylist.List) string {
	var clause strings.Builder
	itr := l.Iterator()
	for itr.Next() {

		if strings.Contains(itr.Value().(string), "LET ") {
			clause.WriteString(itr.Value().(string))
		}
	}

	clause.WriteString("FILTER ")
	itr.Begin()
	for itr.Next() {

		if !strings.Contains(itr.Value().(string), "LET ") {
			clause.WriteString(itr.Value().(string))
		}
	}

	return clause.String()
}


func addQuoteToStrings(ops, value string) string {
	stringOperators := map[string]int{
		"==":  1,
		"===": 1,
		"!=":  1,
		"=~":  1,
		"!~":  1,
	}
	if _, ok := stringOperators[ops]; ok {
		return fmt.Sprintf("'%s'", value)
	}

	return value
}

func dateValidator(str string) error {

	dre, err := buildDate()
	if err != nil {
		return err
	}
	mtch := dre.FindString(str)
	if len(mtch) == 0 {
		return fmt.Errorf("error in validating date %s", str)
	}

	if _, err := now.Parse(mtch); err != nil {
		return fmt.Errorf("could not parse date %s %s", str, err)
	}

	return nil
}
```

## File: database.go
```go
package arangomanager

import (
	"context"
	"fmt"
	"math"
	"time"

	driver "github.com/arangodb/go-driver"
)

const tranSize = 12


type TransactionOptions struct {

	ReadCollections []string

	WriteCollections []string

	ExclusiveCollections []string

	WaitForSync bool

	AllowImplicit bool

	LockTimeout int

	MaxTransactionSize int
}


type Database struct {
	dbh driver.Database
}


func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		MaxTransactionSize: int(math.Pow10(tranSize)),
	}
}

func (d *Database) BeginTransaction(
	ctx context.Context,
	opts *TransactionOptions,
) (*TransactionHandler, error) {
	if opts == nil {
		opts = DefaultTransactionOptions()
	}


	beginOpts := &driver.BeginTransactionOptions{

		MaxTransactionSize: uint64(0),
	}
	beginOpts.WaitForSync = opts.WaitForSync
	if opts.LockTimeout > 0 {
		beginOpts.LockTimeout = time.Duration(opts.LockTimeout) * time.Second
	}

	txID, err := d.dbh.BeginTransaction(
		ctx, driver.TransactionCollections{
			Read:      opts.ReadCollections,
			Write:     opts.WriteCollections,
			Exclusive: opts.ExclusiveCollections,
		}, beginOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := driver.WithTransactionID(ctx, txID)
	return &TransactionHandler{
		db:       d,
		id:       txID,
		ctx:      txCtx,
		canceled: false,
	}, nil
}


func (d *Database) Handler() driver.Database {
	return d.dbh
}



func (d *Database) SearchRows(
	query string,
	bindVars map[string]interface{},
) (*Resultset, error) {

	if err := d.dbh.ValidateQuery(context.Background(), query); err != nil {
		return &Resultset{
				empty: true,
			}, fmt.Errorf(
				"error in validating the query %s",
				err,
			)
	}
	ctx := context.Background()
	cqr, err := d.dbh.Query(ctx, query, bindVars)
	if err != nil {
		return &Resultset{
				empty: true,
			}, fmt.Errorf(
				"error in running search %s",
				err,
			)
	}
	if !cqr.HasMore() {
		return &Resultset{empty: true}, nil
	}

	return &Resultset{cursor: cqr, ctx: ctx}, nil
}


func (d *Database) Search(query string) (*Resultset, error) {
	return d.SearchRows(query, nil)
}



func (d *Database) CountWithParams(
	query string,
	bindVars map[string]interface{},
) (int64, error) {

	if err := d.dbh.ValidateQuery(context.Background(), query); err != nil {
		return 0, fmt.Errorf("error in validating the query %s", err)
	}
	cobj, err := d.dbh.Query(
		driver.WithQueryCount(context.Background(), true),
		query,
		bindVars,
	)
	if err != nil {
		return 0, fmt.Errorf("error with query %s", err)
	}

	return cobj.Count(), nil
}


func (d *Database) Count(query string) (int64, error) {
	return d.CountWithParams(query, nil)
}



func (d *Database) Exec(query string) error {
	return d.Do(query, nil)
}



func (d *Database) Do(query string, bindVars map[string]interface{}) error {
	ctx := driver.WithSilent(context.Background())
	_, err := d.dbh.Query(ctx, query, bindVars)
	if err != nil {
		return fmt.Errorf("error in data modification query %s", err)
	}

	return nil
}



func (d *Database) GetRow(
	query string,
	bindVars map[string]interface{},
) (*Result, error) {
	if err := d.dbh.ValidateQuery(context.Background(), query); err != nil {
		return &Result{
				empty: true,
			}, fmt.Errorf(
				"error in validating the query %s",
				err,
			)
	}
	cqr, err := d.dbh.Query(context.Background(), query, bindVars)

	return d.getResult(cqr, err)
}



func (d *Database) DoRun(
	query string,
	bindVars map[string]interface{},
) (*Result, error) {
	return d.GetRow(query, bindVars)
}


func (d *Database) Get(query string) (*Result, error) {
	return d.GetRow(query, nil)
}



func (d *Database) Run(query string) (*Result, error) {
	return d.GetRow(query, nil)
}


func (d *Database) Collection(name string) (driver.Collection, error) {
	var coll driver.Collection
	ok, err := d.dbh.CollectionExists(context.Background(), name)
	if err != nil {
		return coll, fmt.Errorf("unable to check for collection %s", name)
	}
	if !ok {
		return coll, fmt.Errorf("collection %s has to be created", name)
	}
	coll, err = d.dbh.Collection(context.Background(), name)
	if err != nil {
		return coll, fmt.Errorf("error in getting collection %s", err)
	}

	return coll, nil
}


func (d *Database) CreateCollection(
	name string,
	opt *driver.CreateCollectionOptions,
) (driver.Collection, error) {
	var coll driver.Collection
	ok, err := d.dbh.CollectionExists(context.Background(), name)
	if err != nil {
		return coll, fmt.Errorf("error in collection lookup %s", err)
	}
	if ok {
		return coll, fmt.Errorf("collection %s exists", name)
	}
	coll, err = d.dbh.CreateCollection(context.TODO(), name, opt)
	if err != nil {
		return coll, fmt.Errorf("error in creating collection %s", err)
	}

	return coll, nil
}




func (d *Database) FindOrCreateCollection(
	name string,
	opt *driver.CreateCollectionOptions,
) (driver.Collection, error) {
	var coll driver.Collection
	ok, err := d.dbh.CollectionExists(context.Background(), name)
	if err != nil {
		return coll, fmt.Errorf("unable to check for collection %s", name)
	}
	if ok {
		coll, err = d.dbh.Collection(context.Background(), name)
		if err != nil {
			return coll, fmt.Errorf("error in fetching collection %s", err)
		}

		return coll, nil
	}
	coll, err = d.dbh.CreateCollection(context.TODO(), name, opt)
	if err != nil {
		return coll, fmt.Errorf("error in creating collection %s", err)
	}

	return coll, nil
}


func (d *Database) FindOrCreateGraph(
	name string,
	defs []driver.EdgeDefinition,
) (driver.Graph, error) {
	var grph driver.Graph
	ok, err := d.dbh.GraphExists(context.Background(), name)
	if err != nil {
		return grph, fmt.Errorf("error in graph %s lookup %s", name, err)
	}
	if ok {
		grph, err = d.dbh.Graph(context.Background(), name)
		if err != nil {
			return grph, fmt.Errorf("error in fetching graph %s", err)
		}

		return grph, nil
	}
	grph, err = d.dbh.CreateGraphV2(
		context.Background(),
		name,
		&driver.CreateGraphOptions{EdgeDefinitions: defs},
	)
	if err != nil {
		return grph, fmt.Errorf("error in creating graph %s", err)
	}

	return grph, nil
}


func (d *Database) EnsureGeoIndex(
	coll string, fields []string,
	opts *driver.EnsureGeoIndexOptions,
) (driver.Index, bool, error) {
	var idx driver.Index
	cobj, err := d.Collection(coll)
	if err != nil {
		return idx, false, fmt.Errorf("unable to check for collection %s", coll)
	}
	idx, isOk, err := cobj.EnsureGeoIndex(context.Background(), fields, opts)
	if err != nil {
		return idx, isOk, fmt.Errorf("error in handling index %s", err)
	}

	return idx, isOk, nil
}


func (d *Database) EnsureHashIndex(
	coll string, fields []string,
	opts *driver.EnsureHashIndexOptions,
) (driver.Index, bool, error) {
	var idx driver.Index
	cobj, err := d.Collection(coll)
	if err != nil {
		return idx, false, fmt.Errorf("unable to check for collection %s", coll)
	}
	idx, isOk, err := cobj.EnsureHashIndex(context.Background(), fields, opts)
	if err != nil {
		return idx, isOk, fmt.Errorf("error in handling index %s", err)
	}

	return idx, isOk, nil
}


func (d *Database) EnsurePersistentIndex(
	coll string, fields []string,
	opts *driver.EnsurePersistentIndexOptions,
) (driver.Index, bool, error) {
	var idx driver.Index
	cobj, err := d.Collection(coll)
	if err != nil {
		return idx, false, fmt.Errorf("unable to check for collection %s", coll)
	}
	idx, isOk, err := cobj.EnsurePersistentIndex(
		context.Background(),
		fields,
		opts,
	)
	if err != nil {
		return idx, isOk, fmt.Errorf("error in handling index %s", err)
	}

	return idx, isOk, nil
}


func (d *Database) EnsureSkipListIndex(
	coll string, fields []string,
	opts *driver.EnsureSkipListIndexOptions,
) (driver.Index, bool, error) {
	var idx driver.Index
	cobj, err := d.Collection(coll)
	if err != nil {
		return idx, false, fmt.Errorf("unable to check for collection %s", coll)
	}
	idx, isOk, err := cobj.EnsureSkipListIndex(
		context.Background(),
		fields,
		opts,
	)
	if err != nil {
		return idx, isOk, fmt.Errorf("error in handling index %s", err)
	}

	return idx, isOk, nil
}


func (d *Database) Drop() error {
	if err := d.dbh.Remove(context.Background()); err != nil {
		return fmt.Errorf("error in removing database %s", err)
	}

	return nil
}


func (d *Database) ValidateQ(q string) error {
	if err := d.dbh.ValidateQuery(context.Background(), q); err != nil {
		return fmt.Errorf("error in validating the query %s", err)
	}

	return nil
}


func (d *Database) Truncate(names ...string) error {
	for _, n := range names {
		if _, err := d.Collection(n); err != nil {
			return err
		}
	}
	_, err := d.dbh.Transaction(
		context.Background(),
		truncateFn,
		&driver.TransactionOptions{
			WriteCollections: names,
			ReadCollections:  names,
			Params:           []interface{}{names},
			MaxTransactionSize: func() int {
				size := math.Pow10(tranSize)
				if size > float64(math.MaxInt) {
					return math.MaxInt
				}
				return int(size)
			}(),
		})
	if err != nil {
		return fmt.Errorf("error in truncating collections %s", err)
	}

	return nil
}

func (d *Database) getResult(cdr driver.Cursor, err error) (*Result, error) {
	if err != nil {
		return &Result{empty: true}, fmt.Errorf("error in query %s", err)
	}
	if !cdr.HasMore() {
		return &Result{empty: true}, nil
	}

	return &Result{cursor: cdr}, nil
}
```

## File: resultset.go
```go
package arangomanager

import (
	"context"
	"fmt"

	driver "github.com/arangodb/go-driver"
	"github.com/fatih/structs"
)


type Resultset struct {
	cursor driver.Cursor
	ctx    context.Context
	empty  bool
}


func (r *Resultset) IsEmpty() bool {
	return r.empty
}


func (r *Resultset) Scan() bool {
	if r.empty {
		return false
	}
	if r.cursor.HasMore() {
		return true
	}


	_ = r.cursor.Close()

	return false
}



func (r *Resultset) Read(iface interface{}) error {
	if r.empty {
		return fmt.Errorf("cannot read from empty resultset")
	}

	meta, err := r.cursor.ReadDocument(r.ctx, iface)
	if err != nil {
		return fmt.Errorf("error in reading document %s", err)
	}
	if !structs.IsStruct(iface) {
		return nil
	}
	s := structs.New(iface)
	if f, ok := s.FieldOk("DocumentMeta"); ok {
		if f.IsEmbedded() {
			if err := f.Set(meta); err != nil {
				return fmt.Errorf(
					"error in assigning DocumentMeta to the structure %s",
					err,
				)
			}
		}
	}

	return nil
}




func (r *Resultset) Close() error {
	if r.empty {
		return nil
	}
	if err := r.cursor.Close(); err != nil {
		return fmt.Errorf("error in closing cursor %s", err)
	}

	return nil
}
```
