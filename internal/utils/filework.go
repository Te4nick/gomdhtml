package utils

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := CreateWithDirs(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	sourceFileInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	err = os.Chmod(dst, sourceFileInfo.Mode())
	if err != nil {
		return err
	}
	return nil
}

func CopyDir(src string, dst string) error {
	return os.CopyFS(dst, os.DirFS(src))
}

func CreateWithDirs(filePath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(filePath)
}

func RecreateDir(dirPath string) error {
	if err := os.RemoveAll(dirPath); err != nil {
		return err
	}

	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}
