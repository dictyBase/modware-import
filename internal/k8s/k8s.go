package k8s

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppParams struct {
	Name, Description, Namespace string
	fragment                     string
}

type Application struct {
	metav1.ObjectMeta
	*ImageSpec
	description string
	randomizer  map[string]func() []string
}

func NewApp(args *AppParams, ispec *ImageSpec) (*Application, error) {
	qname, err := RandomFullName(args.Name, args.fragment, 10)
	if err != nil {
		return &Application{}, err
	}
	return &Application{
		description: args.Description,
		ImageSpec:   ispec,
		ObjectMeta: metav1.ObjectMeta{
			Name:            qname,
			Namespace:       args.Namespace,
			ResourceVersion: "v1.0.0",
			Labels: map[string]string{
				"heritage": "naml",
			},
		},
	}, nil
}

// Meta returns the Kubernetes native ObjectMeta which is used to manage applications with naml.
func (a *Application) Meta() *metav1.ObjectMeta {
	return &a.ObjectMeta
}

// Description returns the application description
func (a *Application) Description() string {
	return a.description
}

func (a *Application) RandContainerName(n int, suffix string) (string, error) {
	cname, err := RandomAlphaName(n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s-%s-%s",
		string(cname), a.Meta().Name, suffix,
	), nil
}

type ImageSpec struct {
	repo, tag string
}

func NewImageSpec(repo, tag string) *ImageSpec {
	return &ImageSpec{repo: repo, tag: tag}
}

func (s *ImageSpec) ImageManifest() string {
	return fmt.Sprintf("%s:%s", s.repo, s.tag)
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

func RandomFullName(name, frag string, n int) (string, error) {
	qname := name
	if len(qname) == 0 {
		n, err := RandomAppName()
		if err != nil {
			return n, err
		}
		qname = n
	}
	cname, err := RandomAlphaName(n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s-%s-%s",
		strings.TrimSuffix(Trunc(qname, 63), "-"),
		frag, string(cname),
	), nil
}

// RandomAppName generates a random lowercase alphabetical name
// by combining various adjective, noun, verb and adverb.
func RandomAppName() (string, error) {
	var rstr *strings.Builder
	rmapper := RandMapper()
	for _, g := range []string{"adj", "plural", "verb", "adv"} {
		bidx, err := rand.Int(rand.Reader, big.NewInt(int64(len(rmapper[g]()))))
		if err != nil {
			return "", errors.Errorf("error in generating secure random number %s", err)
		}
		rstr.WriteString(rmapper[g]()[bidx.Int64()])
	}
	return strings.ToLower(rstr.String()), nil
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
		break
	case slen <= ln:
		out = in
		break
	case ln == 0:
		out = in
		break
	case ln < 0:
		out = in[-ln:]
		break
	default:
		out = in[0:ln]
	}
	return out
}
