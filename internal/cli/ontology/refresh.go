package ontology

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/minio/minio-go/v6"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RefreshCmd updates obojson formatted ontologies in S3(minio) storage
var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "updates ontologies in S3 storage",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cloneRepo()
		if err != nil {
			return err
		}
		files, err := obojsonFiles(dir)
		if err != nil {
			return err
		}
		logger := registry.GetLogger()
		for _, e := range files {
			_, err := registry.GetS3Client().FPutObject(
				viper.GetString("s3-bucket"),
				fmt.Sprintf("%s/%s", viper.GetString("s3-bucket-path"), filepath.Base(e)),
				e, minio.PutObjectOptions{
					UserMetadata: map[string]string{
						GroupTag: viper.GetString("group"),
					},
				},
			)
			if err != nil {
				return err
			}
			logger.Debugf("uploaded %s file", filepath.Base(e))
		}
		logger.Infof("refreshed %d obo files", len(files))
		return nil
	},
}

func init() {
	refreshFlags()
	viper.BindPFlags(RefreshCmd.Flags())
}

func refreshFlags() {
	RefreshCmd.Flags().String(
		"repository",
		"https://github.com/dictyBase/migration-data",
		"github repository source of obojson files",
	)
	RefreshCmd.Flags().String(
		"branch",
		"master",
		"branch of github repository",
	)
	RefreshCmd.Flags().String(
		"subfolder",
		"obograph-json",
		"sub folder of clone repository for the source of obojson files",
	)
	RefreshCmd.Flags().String(
		"group",
		"",
		"ontology group name",
	)
}

func obojsonFiles(dir string) ([]string, error) {
	fullPath := filepath.Join(dir, viper.GetString("subfolder"))
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return []string{""}, err
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), "json") {
			continue
		}
		files = append(files, filepath.Join(fullPath, e.Name()))
	}
	return files, nil
}

func cloneRepo() (string, error) {
	dir, err := os.MkdirTemp(os.TempDir(), "*-github")
	if err != nil {
		return dir, fmt.Errorf("error in making temp folder %s", err)
	}
	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:          viper.GetString("repository"),
		SingleBranch: true,
		ReferenceName: plumbing.NewBranchReferenceName(
			viper.GetString("branch"),
		)},
	)
	return dir, err
}
