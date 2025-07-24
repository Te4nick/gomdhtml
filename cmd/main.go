package main

import (
	"errors"
	"flag"
	"os"

	"github.com/gomdhtml/internal/config"
	"github.com/gomdhtml/internal/log"
	"github.com/gomdhtml/internal/utils"
)

const (
	defaultConfigFile  string = ""
	defaultInputDir    string = "./"
	defaultOutputDir   string = "./output"
	defaultDebug       bool   = false
	configPathFlagName string = "config"
	outputDirFlagName  string = "out"
	debugFlagName      string = "debug"
)

func main() {
	outputDir := flag.String(outputDirFlagName, defaultOutputDir, "Path to the output directory")
	configPath := flag.String(configPathFlagName, defaultConfigFile, "Path to config file")
	debug := flag.Bool(debugFlagName, defaultDebug, "Enable debug log messages")
	flag.Parse()
	inputPaths := flag.Args()

	log.Setup(*debug)
	log.Debug("debug messages are on")

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

	config.ParseConfig(*configPath, config.CLIArgs{
		InputDir: inputDir,
		OutputDir: *outputDir,
	})

	log.Infof("render from %s to %s", inputDir, *outputDir)
	err := utils.CompileCatalog(inputDir, *outputDir)
	if err != nil {
		log.Err(err, "error while compiling catalog")
		os.Exit(1)
	}
	log.Info("successfully rendered")
}
