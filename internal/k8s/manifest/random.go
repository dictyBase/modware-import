package manifest

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
)

func RandContainerName(name, frag string, n int) (string, error) {
	cname, err := RandomAlphaName(n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s-%s-%s",
		strings.ToLower(cname), name, frag,
	), nil
}

func RandomAlphaName(n int) (string, error) {
	const alphabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	cname := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabets))))
		if err != nil {
			return "", errors.Errorf("error in generating secure random number %s", err)
		}
		cname[i] = alphabets[num.Int64()]
	}
	return string(cname), nil
}

func FullName(name, frag string, nameLen int) (string, error) {
	appNameLen := nameLen - (len(name) + len(frag))
	if appNameLen < 0 {
		appNameLen *= -1
	}
	rndName, err := RandomAppName(appNameLen)
	if err != nil {
		return rndName, fmt.Errorf("error in generating random app name %s", err)
	}
	return Trunc(fmt.Sprintf(
		"%s-%s-%s",
		rndName,
		strings.TrimSuffix(name, "-"),
		frag,
	), nameLen), nil
}

// RandomAppName generates a random lowercase alphabetical name of fixed length
// by combining various adjective, noun, verb and adverb.
func RandomAppName(nameLen int) (string, error) {
	var rstr *strings.Builder
	rmapper := RandMapper()
	for _, g := range []string{"adj", "plural", "verb", "adv"} {
		bidx, err := rand.Int(rand.Reader, big.NewInt(int64(len(rmapper[g]()))))
		if err != nil {
			return "", errors.Errorf("error in generating secure random number %s", err)
		}
		rstr.WriteString(rmapper[g]()[bidx.Int64()])
	}
	return strings.ToLower(Trunc(rstr.String(), -1*nameLen)), nil
}

func RandMapper() map[string]func() []string {
	return map[string]func() []string{
		"adj":    Adjective,
		"plural": PluralNoun,
		"verb":   Verb,
		"adv":    Adverb,
	}
}

// Trunc truncates a string to the given length either
// from the start(positive length) or from the
// end(negative length)
func Trunc(in string, ln int) string {
	var out string
	slen := len(in)
	switch {
	case slen == 0:
		out = in
	case slen <= ln:
		out = in
	case ln == 0:
		out = in
	case ln < 0:
		out = in[-ln:]
	default:
		out = in[0:ln]
	}
	return out
}