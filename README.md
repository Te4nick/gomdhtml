# gomdhtml
Compile HTMX, CSS, Markdown into web static files

## Get and Run
```bash
git clone https://github.com/Te4nick/gomdhtml.git
cd gomdhtml
go run ./cmd
```

## Usage
```bash
go run ./cmd --out ./output ./ 
```

## Documentation
### Files and Directories
This tool accepts input directory path as argument (default `./`) and output directory path with `--out`
optional argument (default `./output`). It searches markdown files in `md` directory of input directory
path and pastes them into html files in `html` directory of input directory path. By default tool uses 
`html/.template.html` as template for html files. If you want to use different template, you can redefine 
it per markdown file by creating html file with the same name as markdown file in `html` directory. 
The default css file is expected to be `css/.template.css` and also redefinable for each markdown file 
by creating css file with the same name as markdown file in `css` directory. The project directory 
structure with default agruments and compiled files in `output` direcory can be seen here:
```
.
├── css
│   ├── .template.css
│   └── articles.css
├── html
│   ├── .template.html
│   └── articles.html
├── md
│   ├── articles
│   │   └── gomdhtml.md
│   ├── articles.md
│   └── index.md
└── output
    ├── articles
    │   └── gomdhtml.html
    ├── css
    │   ├── .template.css
    │   └── articles.css
    ├── articles.html
    └── index.html
```
To create links to other pages just use markdown link with relative path to desired markdown file.
It will be automatically converted into link to corresponding html file in output directory making
it consistent between markdown and html navigation.

### HTML Templates
Tool expects `{{.Content}}` htmx tag in html file to insert markdown content. It also aptionally 
searches for `{{.CSS}}` tag to insert css file and `{{.Title}}` tag is populated with first H1
appearence in markdown file. Check `example` project directory for use case.