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
	staticDirName   = "static"
	defaultFileName = ".template"
)

func CompileCatalog(inputDirPath, outputDirPath string) error {
	mdDir := filepath.Join(inputDirPath, mdDirName)

	if err := RecreateDir(outputDirPath); err != nil {
		return err
	}

	if err := filepath.WalkDir(mdDir, func(path string, d fs.DirEntry, err error) error {
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
			mdRelPath, err := filepath.Rel(mdDir, path)
			if err != nil {
				return err
			}

			log.Debugf("relPath=%s, parentDir=%s", mdRelPath, filepath.Dir(mdRelPath))
			htmlFilePath, err := resolveInputResoucePath(mdRelPath, htmlDirName, inputDirPath)
			if err != nil {
				return errors.New("Default HTML file required: " + htmlFilePath)
			}

			cssFilePath, err := resolveInputResoucePath(mdRelPath, filepath.Join(staticDirName, cssDirName), inputDirPath)
			if err != nil {
				cssFilePath = ""
			}

			outputFilePath := filepath.Join(outputDirPath, strings.TrimSuffix(mdRelPath, ".md")+".html")
			log.Debugf("htmlFilePath=%s, path=%s, cssFilePath=%s, outputFilePath=%s", htmlFilePath, path, cssFilePath, outputFilePath)
			if err := RenderFileHTML(htmlFilePath, path, cssFilePath, outputFilePath); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := CopyDir( // copy static files (css + images etc.)
		filepath.Join(inputDirPath, staticDirName),
		filepath.Join(outputDirPath, staticDirName),
	); err != nil {
		return err
	}

	return nil
}

func resolveInputResoucePath(mdRelPath, dirName, inputDirPath string) (string, error) { // get according
	resDir := filepath.Join(inputDirPath, dirName)
	resFilePath := filepath.Join(resDir, strings.TrimSuffix(mdRelPath, ".md")+"."+filepath.Base(dirName))
	if _, err := os.Stat(resFilePath); os.IsNotExist(err) {
		resFilePath = filepath.Join(resDir, defaultFileName+"."+filepath.Base(dirName))
		if _, err := os.Stat(resFilePath); os.IsNotExist(err) {
			return "", err
		}
	}
	return resFilePath, nil
}
