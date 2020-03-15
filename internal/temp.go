package worker

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func (w *Worker) generateDirectory() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get current wd: %w", err)
	}

	dir, err := ioutil.TempDir(filepath.Join(pwd, "builds"), "")
	if err != nil {
		return "", fmt.Errorf("could not create temp dir: %w", err)
	}

	return dir, nil
}

func removeExtension(fileName string) string {
	return fileName[0 : len(fileName)-len(filepath.Ext(fileName))]
}
