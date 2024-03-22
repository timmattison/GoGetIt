package internal

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrCouldNotStatFile       = fmt.Errorf("could not stat file")
	ErrPathIsDirectoryNotFile = fmt.Errorf("path is a directory, not a file")
	ErrPathIsFileNotDirectory = fmt.Errorf("path is a file, not a directory")
)

func FileExists(filename string) error {
	var info os.FileInfo
	var err error

	if info, err = os.Stat(filename); err != nil {
		if !os.IsNotExist(err) {
			return errors.Join(ErrCouldNotStatFile, err)
		}

		// The file does not exist
		return os.ErrNotExist
	}

	if info.IsDir() {
		return ErrPathIsDirectoryNotFile
	}

	return nil
}

func DirectoryExists(directory string) error {
	var err error

	if err = FileExists(directory); err != nil {
		if errors.Is(err, ErrPathIsDirectoryNotFile) {
			// Found it and it is a directory
			return nil
		}

		// Probably not found
		return err
	}

	// Found it and it is a file
	return ErrPathIsFileNotDirectory
}
