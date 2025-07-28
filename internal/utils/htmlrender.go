package utils

import (
	"bytes"
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomdhtml/internal/config"
	"github.com/gomdhtml/internal/filework"
	"github.com/gomdhtml/internal/log"
	"github.com/yuin/goldmark"
)

var (
	titleReHTML *regexp.Regexp = regexp.MustCompile(`<h1.*?>(.*?)<\/h1>`)
	linkReHTML  *regexp.Regexp = regexp.MustCompile(`<a href="([^http][^"]*?)\.md">([^<]*)</a>`)
	// imgHTMLRe   *regexp.Regexp = regexp.MustCompile(`<img\s+src="([^"]+)"\s+alt="([^"]*)"\s*\/?>`)
	// ulRe    *regexp.Regexp = regexp.MustCompile(`(?i)(<ul>.*?</ul>)(\s*<ul>.*?</ul>)+`)
	// liRe    *regexp.Regexp = regexp.MustCompile(`(?i)<li>(.*?)</li>`)
)

func newDataHTML(mdFile string) (map[string]template.HTML, error) {
	contentHTML, err := mdToHTML(mdFile)
	if err != nil {
		return nil, err
	}

	customDataHTML, err := parseCustomDataHTML(mdFile)
	if err != nil {
		return nil, err
	}

	customDataHTML["CSS"] = generateCSSTag(mdFile)
	customDataHTML["Title"] = generateTitleTag(mdFile, contentHTML)
	customDataHTML["Content"] = contentHTML

	return customDataHTML, nil
}

func generateTitleTag(mdFile string, html template.HTML) template.HTML {
	matches := titleReHTML.FindStringSubmatch(string(html))
	if matches != nil && len(matches) < 2 {
		return template.HTML("<title>" + matches[1] + "</title>")
	}

	return template.HTML("<title>" + strings.TrimSuffix(filepath.Base(mdFile), ".md") + "</title>")
}

func generateCSSTag(mdFile string) template.HTML {
	cssFile, err := resolveInputResoucePath(mdFile, filepath.Join(staticDirName, cssDirName), defaultFileName)
	if err != nil {
		return ""
	}

	return template.HTML("<link rel=\"stylesheet\" href=\"/" + cssFile + "\">")
}

func RenderFileHTML(mdFile, outputDir string) error {
	templateFile, err := resolveInputResoucePath(mdFile, htmlDirName, defaultFileName)
	if err != nil {
		return errors.New("Default HTML file required: " + templateFile)
	}

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	data, err := newDataHTML(mdFile)
	if err != nil {
		return err
	}

	mdRel, err := filework.GetInputRelPath(mdFile, mdDirName)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(outputDir, filework.ReplaceExt(mdRel, "html"))
	outHTML, err := filework.CreateWithDirs(outputFile)
	if err != nil {
		log.Err(err, "error creating html output file")
		return err
	}
	defer outHTML.Close()

	err = tmpl.Execute(outHTML, data)
	if err != nil {
		panic(err)
	}

	return nil
}

func mdToHTML(mdFile string) (template.HTML, error) {
	mdContent, err := os.ReadFile(mdFile)
	if err != nil {
		return "", err
	}

	var mdBuf bytes.Buffer
	if err := goldmark.Convert(mdContent, &mdBuf); err != nil {
		return "", err
	}

	convertedHTML := linkReHTML.ReplaceAllString(
		strings.ReplaceAll(mdBuf.String(), "</ul>\n<ul>\n", ""), // combine <ul>
		`<a href="$1.html">$2</a>`,                              // replace link paths from md to html
	)

	return template.HTML(convertedHTML), nil
}

func parseCustomDataHTML(mdFile string) (map[string]template.HTML, error) {
	log.Debug("start parsing custom data keys for " + mdFile)

	customsHTML := map[string]template.HTML{}

	for key, suffix := range config.Get().CustomDataKeys {
		mdDefaultFileName := defaultFileName + "-" + suffix
		mdSuffixFile, err := resolveInputResoucePath(
			filework.AddNameSuffix(mdFile, suffix),
			mdDirName,
			mdDefaultFileName,
		)

		if err != nil {
			log.Warnf(
				"expected default file %s for custom data key {{.%s}}",
				filepath.Join(mdDirName, mdDefaultFileName+".md"),
				key,
			)
			customsHTML[key] = ""
			continue
		}

		html, err := mdToHTML(mdSuffixFile)
		if err != nil {
			return nil, err
		}

		customsHTML[key] = html
	}

	return customsHTML, nil
}
