package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"errors"

	log "github.com/gomdhtml/internal/utils/log"
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

	if err := os.RemoveAll(outputDirPath); err != nil {
		return err
	}

	if err := os.Mkdir(outputDirPath, os.ModePerm); err != nil {
		return err
	}

	err := filepath.WalkDir(mdDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		log.Debug("walking over filename=" + info.Name())

		if !d.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			log.Debug("entering compile for filename=" + info.Name())
			relPath, err := filepath.Rel(mdDir, path)
			if err != nil {
				return err
			}
			log.Debugf("relPath=%s, parentDir=%s", relPath, filepath.Dir(relPath))
			// if err := os.MkdirAll(filepath.Dir(relPath), os.ModePerm); err != nil {

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
			log.Debugf("htmlFilePath=%s, path=%s, cssFilePath=%s, outputFilePath=%s", htmlFilePath, path, cssFilePath, outputFilePath)
			if err := RenderFileHTML(htmlFilePath, path, cssFilePath, outputFilePath); err != nil {
				return err
			}

			if cssFilePath != "" {
				if err := CopyFile(cssFilePath, filepath.Join(outputDirPath, cssDirName, filepath.Base(cssFilePath))); err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}
