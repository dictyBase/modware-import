package git

import (
	"fmt"
	"os"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var rgxp = regexp.MustCompile(`^v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?$`)

func CloneRepo(repo, branch string) (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "*-github")
	if err != nil {
		return dir, fmt.Errorf("error in making temp folder %s", err)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		ReferenceName: gitRefName(branch),
		URL:           repo,
		SingleBranch:  true,
	})
	return dir, err
}

func gitRefName(ref string) plumbing.ReferenceName {
	if rgxp.MatchString(ref) {
		return plumbing.NewTagReferenceName(ref)
	}
	return plumbing.NewBranchReferenceName(ref)
}
