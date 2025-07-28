package filework

import (
	"io"
	"os"
	"path/filepath"

	"github.com/gomdhtml/internal/config"
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

func CreateWithDirs(file string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(file)
}

func RecreateDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func GetInputRelPath(file, dir string) (string, error) {
	relDir := filepath.Join(config.Get().InputDir, dir)
	fileRel, err := filepath.Rel(relDir, file)
	if err != nil {
		return "", err
	}

	return fileRel, nil
}

func AddNameSuffix(file, suffix string) string {
	ext := filepath.Ext(file)
	noExt := file[:len(file)-len(ext)]
	return noExt + "-" + suffix + ext
}

func ReplaceExt(file, newExtName string) string {
	return file[:len(file)-len(filepath.Ext(file))+1] + newExtName
}
