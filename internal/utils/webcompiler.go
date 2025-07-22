package utils

import (
	"os"
	"path/filepath"
	"strings"

	"errors"
)

const (
	mdDirName       = "md"
	htmlDirName     = "html"
	cssDirName      = "css"
	defaultFileName = ".template"
)

func CompileCatalog(inputDirPath, outputDirPath string) error {
	mdDir := filepath.Join(inputDirPath, mdDirName)
	htmlDir := filepath.Join(inputDirPath, htmlDirName)
	cssDir := filepath.Join(inputDirPath, cssDirName)

	if err := os.MkdirAll(filepath.Join(outputDirPath, cssDirName), os.ModePerm); err != nil {
		return err
	}

	err := filepath.Walk(mdDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			relPath, err := filepath.Rel(mdDir, path)
			if err != nil {
				return err
			}

			htmlFilePath := filepath.Join(htmlDir, strings.TrimSuffix(relPath, ".md")+".html")
			if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
				htmlFilePath = filepath.Join(htmlDir, defaultFileName+".html")
				if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
					return errors.New("Default HTML file required: " + htmlFilePath)
				}
			}

			cssFilePath := filepath.Join(cssDir, strings.TrimSuffix(relPath, ".md")+".css")
			if _, err := os.Stat(cssFilePath); os.IsNotExist(err) {
				cssFilePath = filepath.Join(cssDir, defaultFileName+".css")
				if _, err := os.Stat(htmlFilePath); os.IsNotExist(err) {
					cssFilePath = ""
				}
			}

			outputFilePath := filepath.Join(outputDirPath, strings.TrimSuffix(relPath, ".md")+".html")
			if err := RenderFileHTML(htmlFilePath, path, cssFilePath, outputFilePath); err != nil {
				return err
			}

			if err := CopyFile(cssFilePath, filepath.Join(outputDirPath, cssDirName, filepath.Base(cssFilePath))); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
