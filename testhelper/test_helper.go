package testhelper

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func ProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}

func MustArgumentFilePath() string {
	return fmt.Sprintf("%s", "scripts/gen_test_data/arguments")
}

func ToColoredFilenamePathFromRoot(args []string) string {
	return ToFilename(args) + "_colored.txt"
}

func ToPlainFilenamePathFromRoot(args []string) string {
	return ToFilename(args) + "_plain.txt"
}

func ToFilename(args []string) string {
	safeArgs := make([]string, len(args))
	for i, arg := range args {
		safeArgs[i] = strings.Replace(arg, "/", "__", -1)
	}
	return fmt.Sprintf("%s/printer/testdata/%s", ProjectRoot(), strings.Join(safeArgs, "_"))
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
