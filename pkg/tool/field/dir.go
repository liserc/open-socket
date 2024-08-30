package field

import (
	"fmt"
	"github.com/open-socket/pkg/tool/errs"
	"os"
	"path/filepath"
)

// OutDir creates the absolute path name from path and checks if the path exists and is a directory.
// Returns absolute path including trailing '/' or error if the path does not exist or is not a directory.
func OutDir(path string) (string, error) {
	outDir, err := filepath.Abs(path)
	if err != nil {
		return "", errs.WrapMsg(err, "failed to resolve absolute path", "path", path)
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errs.WrapMsg(err, "output directory does not exist", "path", outDir)
		}
		return "", errs.WrapMsg(err, "failed to stat output directory", "path", outDir)
	}

	if !stat.IsDir() {
		return "", errs.WrapMsg(fmt.Errorf("specified path %s is not a directory", outDir), "specified path is not a directory", "path", outDir)
	}
	outDir += "/"
	return outDir, nil
}
