package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/gomdhtml/internal/utils"
	log "github.com/gomdhtml/internal/utils/log"
)

const (
	defaultInputDir  string = "./"
	defaultOutputDir string = "./output"
)

func main() {
	outputDir := flag.String("out", defaultOutputDir, "Path to the output directory")
	debugLogger := flag.Bool("debug", false, "Enable debug log messages")
	flag.Parse()
	inputPaths := flag.Args()

	fmt.Println("DEBUG", *debugLogger)
	log.Setup(*debugLogger)

	inputDir := ""
	switch len(inputPaths) {
	case 0:
		inputDir = defaultInputDir
	case 1:
		inputDir = inputPaths[0]
	default:
		log.Err(errors.New("MultipleInputPaths"), "single input path must be specified")
		os.Exit(1)
	}

	log.Infof("render from %s to %s", inputDir, *outputDir)
	err := utils.CompileCatalog(inputDir, *outputDir)
	if err != nil {
		log.Err(err, "error while compiling catalog")
		os.Exit(1)
	}
	log.Info("successfully rendered")
}
