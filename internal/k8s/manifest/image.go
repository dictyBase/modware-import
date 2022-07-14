package manifest

import (
	"fmt"
)

func Image(repo, tag string) string {
	return fmt.Sprintf("%s:%s", repo, tag)
}
