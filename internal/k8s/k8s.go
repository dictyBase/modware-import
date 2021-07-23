package k8s

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/cockroachdb/errors"
	v1 "k8s.io/api/core/v1"
)

type AppParams struct {
	Name, Description, Namespace string
	Fragment                     string
}

type ImageSpec struct {
	repo, tag  string
	PullPolicy string
}

func NewImageSpec(repo, tag, policy string) *ImageSpec {
	return &ImageSpec{
		repo:       repo,
		tag:        tag,
		PullPolicy: policy,
	}
}

func (s *ImageSpec) ImageManifest() string {
	return fmt.Sprintf("%s:%s", s.repo, s.tag)
}

func RandContainerName(name, frag string, n int) (string, error) {
	cname, err := RandomAlphaName(n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s-%s-%s",
		string(cname), name, frag,
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

func FullName(name, frag string) string {
	return fmt.Sprintf(
		"%s-%s",
		strings.TrimSuffix(Trunc(name, 63), "-"),
		frag,
	)
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

func ArangoSecManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ARANGODB_PASS",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "arangodb.password",
				},
			},
		},
	}
}

func ArangoConfigManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ARANGODB_USER",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "arangodb.user",
				},
			},
		},
		{
			Name: "ARANGODB_DATABASE",
			ValueFrom: &v1.EnvVarSource{
				ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "arangodb.modware-annotation",
				},
			},
		},
	}
}

func MinioSecManifest() []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name: "ACCESS_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "minio.accesskey",
				},
			},
		},
		{
			Name: "SECRET_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: "dictybase-configuration"},
					Key:                  "minio.secretkey",
				},
			},
		},
	}

}

func LogEnvManifest(level string) []v1.EnvVar {
	return []v1.EnvVar{
		{
			Name:  "LOG_LEVEL",
			Value: level,
		},
	}
}
