package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gomdhtml/internal/utils"
)

const (
	defaultInputDir  string = "./"
	defaultOutputDir string = "./output"
)

func main() {
	outputDir := flag.String("out", defaultOutputDir, "Path to the output directory")
	flag.Parse()
	inputPaths := flag.Args()

	inputDir := ""
	switch len(inputPaths) {
	case 0:
		inputDir = defaultInputDir
	case 1:
		inputDir = inputPaths[0]
	default:
		fmt.Println("Error: single input path must be specified")
		os.Exit(1)
	}

	fmt.Printf("Render result from \"%s\" into \"%s\" directory", inputDir, *outputDir)
	utils.CompileCatalog(inputDir, *outputDir)
}
