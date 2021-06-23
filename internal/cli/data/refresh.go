package data

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"

	"github.com/dictyBase/modware-import/internal/cli"
	"github.com/dictyBase/modware-import/internal/collection"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RefreshCmd updates dictybase data files in S3(minio) storage
var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "updates data files in S3 storage",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := cli.CloneRepo(
			viper.GetString("repository"),
			viper.GetString("branch"),
		)
		if err != nil {
			return err
		}
		fmap, err := allFileReaders(dir)
		if err != nil {
			return err
		}
		logger := registry.GetLogger()
		for path, rd := range fmap {
			_, err := registry.GetS3Client().PutObject(
				viper.GetString("s3-bucket"),
				fmt.Sprintf("%s/%s", viper.GetString("s3-bucket-path"), path),
				rd, -1, minio.PutObjectOptions{
					UserMetadata: map[string]string{
						GroupTag: viper.GetString("group"),
					},
				},
			)
			if err != nil {
				return err
			}
			logger.Debugf("uploaded file to path %s", path)
		}
		logger.Infof("refreshed %d data files", len(fmap))
		return nil
	},
}

func allFileReaders(dir string) (map[string]io.Reader, error) {
	fullPath := filepath.Join(dir, viper.GetString("subfolder"))
	fmap := make(map[string]io.Reader)
	err := filepath.Walk(fullPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error in handling path %s %s", path, err)
		}
		if info.IsDir() {
			return filepath.SkipDir
		}
		fh, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error in opening file %s %s", path, err)
		}
		// logic to get parent dirs except the root
		pathParts := strings.Split(
			strings.TrimLeft(path, "/"),
			string(os.PathSeparator),
		)
		idx := collection.Index(pathParts, viper.GetString("subfolder"))
		if idx == -1 {
			return fmt.Errorf("error in finding subfolder in path parts %s", path)
		}
		fmap[strings.Join(pathParts[idx+1:], "/")] = fh
		return nil
	})
	return fmap, err
}

func init() {
	refreshFlags()
	viper.BindPFlags(RefreshCmd.Flags())
}

func refreshFlags() {
	RefreshCmd.Flags().String(
		"repository",
		"https://github.com/dictyBase/migration-data",
		"github repository source of data files",
	)
	RefreshCmd.Flags().String(
		"branch",
		"master",
		"branch of github repository",
	)
	RefreshCmd.Flags().String(
		"subfolder",
		"import",
		"sub folder of clone repository for the source of data files",
	)
	RefreshCmd.Flags().String(
		"group",
		"",
		"S3 group name metadata",
	)
}
