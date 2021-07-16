package k8s

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// RandomReleaseName geneates a random lowercase alphabetical name
func RandomReleaseName() string {
	aidx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Adjective()))))
	pidx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(PluralNoun()))))
	vidx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Verb()))))
	advidx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Adverb()))))
	return strings.ToLower(
		fmt.Sprintf(
			"%s%s%s%s",
			Adjective()[aidx.Int64()],
			PluralNoun()[pidx.Int64()],
			Verb()[vidx.Int64()],
			Adverb()[advidx.Int64()],
		))
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

type Release struct {
	Namespace string
	Service   string
	Name      string
}

type Application struct {
	Release     *Release
	Version     string
	Description string
	Name        string
}

func (app *Application) QualifiedName() string {
	var b *strings.Builder
	b.WriteString(app.Name)
	rname := app.Release.Name
	if len(rname) == 0 {
		rname = RandomReleaseName()
	}
	b.WriteString(fmt.Sprintf("%s-", rname))
	return strings.TrimSuffix(Trunc(b.String(), 63), "-")
}
