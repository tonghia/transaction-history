package args

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func ParsePeriod(period string) (string, error) {
	if period == "" {
		flag.Usage()
		return "", errors.New("-period argument are required")
	}

	return period, nil
}

func ParseFilePath(filePathPtr string) (string, error) {
	if filePathPtr == "" {
		flag.Usage()
		return "", errors.New("-file arguments are required")
	}

	// Verify that the file exists and is a regular file.
	absPath, err := filepath.Abs(filePathPtr)
	if err != nil {
		return "", fmt.Errorf("error resolving absolute path: %v", err)
	}

	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", absPath)
	}
	if err != nil {
		return "", fmt.Errorf("error accessing file: %v", err)
	}
	if fileInfo.IsDir() {
		return "", fmt.Errorf("provided file path is a directory, not a file: %s", absPath)
	}

	return absPath, nil
}
