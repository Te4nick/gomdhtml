package utils

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomdhtml/internal/config"
	"github.com/gomdhtml/internal/filework"
	"github.com/gomdhtml/internal/log"
)

const (
	mdDirName       = "md"
	htmlDirName     = "html"
	cssDirName      = "css"
	staticDirName   = "static"
	defaultFileName = ".template"
)

func CompileCatalog(inputDir, outputDir string) error {
	mdDir := filepath.Join(inputDir, mdDirName)

	if err := filework.RecreateDir(outputDir); err != nil {
		return err
	}

	if err := filepath.WalkDir(mdDir, func(mdFile string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		log.Debug("walking over filename=" + info.Name())

		if !d.IsDir() && strings.HasSuffix(info.Name(), ".md") {
			log.Debug("entering compile for filename=" + mdFile)
			if err := RenderFileHTML(mdFile, outputDir); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	if err := filework.CopyDir( // copy static files (css + images etc.)
		filepath.Join(inputDir, staticDirName),
		filepath.Join(outputDir, staticDirName),
	); err != nil {
		return err
	}

	return nil
}

func resolveInputResoucePath(mdFile, dirName string) (string, error) { // get according
	mdRel, err := filework.GetInputRelPath(mdFile, dirName)
	if err != nil {
		return "", err
	}
	resDir := filepath.Join(config.Get().InputDir, dirName)
	resFile := filepath.Join(resDir, strings.TrimSuffix(mdRel, ".md")+"."+filepath.Base(dirName))
	if _, err := os.Stat(resFile); os.IsNotExist(err) {
		resFile = filepath.Join(resDir, defaultFileName+"."+filepath.Base(dirName))
		if _, err := os.Stat(resFile); os.IsNotExist(err) {
			return "", err
		}
	}
	return resFile, nil
}
