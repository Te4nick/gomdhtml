package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"

	log "github.com/gomdhtml/internal/utils/log"
	"github.com/yuin/goldmark"
)

var (
	titleHTMLRe *regexp.Regexp = regexp.MustCompile(`<h1.*?>(.*?)<\/h1>`)
	linkHTMLRe  *regexp.Regexp = regexp.MustCompile(`<a href="([^http][^"]*?)\.md">([^<]*)</a>`)
	// imgHTMLRe   *regexp.Regexp = regexp.MustCompile(`<img\s+src="([^"]+)"\s+alt="([^"]*)"\s*\/?>`)
	// ulRe    *regexp.Regexp = regexp.MustCompile(`(?i)(<ul>.*?</ul>)(\s*<ul>.*?</ul>)+`)
	// liRe    *regexp.Regexp = regexp.MustCompile(`(?i)<li>(.*?)</li>`)
)

type dataHTML struct {
	CSS     template.HTML
	Title   template.HTML
	Content template.HTML
}

func newDataHTML(mdBytes []byte, cssFilePath string) (*dataHTML, error) {
	var mdBuf bytes.Buffer
	if err := goldmark.Convert(mdBytes, &mdBuf); err != nil {
		return nil, err
	}

	convertedHTML := linkHTMLRe.ReplaceAllString(
		strings.ReplaceAll(mdBuf.String(), "</ul>\n<ul>\n", ""), // combine <ul>
		`<a href="$1.html">$2</a>`,                              // replace link paths from md to html
	)
	title, err := generateTitleTag(convertedHTML)
	if err != nil {
		return nil, err
	}

	return &dataHTML{
		CSS:     generateCSSTag(cssFilePath),
		Title:   template.HTML(title),
		Content: template.HTML(convertedHTML),
	}, nil
}

func generateTitleTag(html string) (template.HTML, error) {
	matches := titleHTMLRe.FindStringSubmatch(html)
	if len(matches) < 2 {
		return "", fmt.Errorf("no <h1> tag found")
	}

	return template.HTML("<title>" + matches[1] + "</title>"), nil
}

func generateCSSTag(cssFilePath string) template.HTML {
	if cssFilePath == "" {
		return ""
	}

	return template.HTML("<link rel=\"stylesheet\" href=\"/" + cssFilePath + "\">")
}

func RenderFileHTML(templateFilePath, mdFilePath, cssFilePath, outputFilePath string) error {
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	mdContent, err := os.ReadFile(mdFilePath)
	if err != nil {
		return err
	}

	data, err := newDataHTML(mdContent, cssFilePath)
	if err != nil {
		return err
	}

	outHTML, err := CreateWithDirs(outputFilePath)
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
