- [📄Articles](../articles.md)

# Markdown to static html pages converter

## Description
ok. the idea was that i write something like obsidian database and then render as static pages website.
now let's together see how it turned out.

## Used tech
- [go 1.24.5](https://go.dev/doc/devel/release#go1.24.minor)
- [yuin/goldmark](https://github.com/yuin/goldmark)

## Current state
for now tool only works with manually defining .md to render as .html in main.go.

## Plan
- add support for auto-detecting .md files in provided directory
- add support for .css files including in .html

### Folder hierarchy
will make binary to accept input folder path as argument, and then will render all .md files in that folder to .html files
into provided output folder `./output` as default.
let's say that `./md` will be default folder for markdown. then there will be `./html` for templates to render md into.
by default let's use `./html/template.html` to render into. if in `./html` there will be .html with the same name as md in `./md`
it will render into that template and if in `./css` there will be .css with the same name as md in `./md` it will include that
particular .css