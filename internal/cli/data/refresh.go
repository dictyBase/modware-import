package data

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dictyBase/modware-import/internal/registry"
	"github.com/minio/minio-go/v6"

	"github.com/dictyBase/modware-import/internal/git"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

// RefreshCmd updates dictybase data files in S3(minio) storage
var RefreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "updates data files in S3 storage",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, _ := cmd.Flags().GetString("repository")
		branch, _ := cmd.Flags().GetString("branch")
		subf, _ := cmd.Flags().GetString("subfolder")
		bucketPath, _ := cmd.Flags().GetString("s3-bucket-path")
		dir, err := git.CloneRepo(repo, branch)
		if err != nil {
			return err
		}
		fmap, err := newdataFileManager(subf).allFileReaders(dir)
		if err != nil {
			return err
		}
		logger := registry.GetLogger()
		for path, rd := range fmap {
			_, err := registry.GetS3Client().PutObject(
				viper.GetString("s3-bucket"),
				fmt.Sprintf("%s/%s", bucketPath, path),
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

type dataFileManager struct {
	subFolder string
	rdmap     map[string]io.Reader
}

func newdataFileManager(folder string) *dataFileManager {
	return &dataFileManager{
		subFolder: folder,
		rdmap:     make(map[string]io.Reader),
	}
}

func (d *dataFileManager) allFileReaders(
	dir string,
) (map[string]io.Reader, error) {
	return d.rdmap, filepath.Walk(
		filepath.Join(dir, d.subFolder),
		d.pathWalker,
	)
}

func (d *dataFileManager) pathWalker(
	path string,
	info fs.FileInfo,
	err error,
) error {
	if err != nil {
		return fmt.Errorf("error in handling path %s %s", path, err)
	}
	if info.IsDir() {
		return nil
	}
	fh, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error in opening file %s %s", path, err)
	}
	subPath, err := d.pathWithoutFolder(path)
	d.rdmap[subPath] = bufio.NewReader(fh)
	return err
}

func (d *dataFileManager) pathWithoutFolder(path string) (string, error) {
	pathParts := strings.Split(
		strings.TrimLeft(path, "/"),
		string(os.PathSeparator),
	)
	idx := slices.Index(pathParts, d.subFolder)
	if idx == -1 {
		return "", fmt.Errorf(
			"error in finding subfolder in path parts %s",
			path,
		)
	}
	return strings.Join(pathParts[idx+1:], "/"), nil
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
