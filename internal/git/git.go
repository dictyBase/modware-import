package git

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneRepo(repo, branch string) (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "*-github")
	if err != nil {
		return dir, fmt.Errorf("error in making temp folder %s", err)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		URL:           repo,
		SingleBranch:  true,
	})
	return dir, err
}
