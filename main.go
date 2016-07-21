package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jjeffery/vt-motoli/scanner"
	"github.com/jjeffery/vt-motoli/story"
)

func main() {
	log.SetFlags(0)

	// TODO(jpj): start with very simple command line, can expand later
	if len(os.Args) != 2 {
		log.Fatal("usage: vt-motoli <file>")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	s := scanStory(file)
	printStory(s)
}

func scanStory(r io.Reader) *story.Story {
	s := story.New()
	scan := scanner.New(r)
	for scan.Scan() {
		if scan.Prefixes.Match("Page", "Line") {
			pageNum := scan.Prefixes[0].Index
			lineNum := scan.Prefixes[1].Index
			continuationNum := scan.Prefixes[1].Cont
			line := s.Page(pageNum).Line(lineNum)
			line.SetText(continuationNum, scan.Arg)
		} else if scan.Prefixes.Match("Page", "Text") {
			pageNum := scan.Prefixes[0].Index
			textNum := scan.Prefixes[1].Index
			page := s.Page(pageNum)
			page.SetText(textNum, scan.Arg)
		} else if scan.Prefixes.Match("Page", "Time") {
			pageNum := scan.Prefixes[0].Index
			lineNum := scan.Prefixes[1].Index
			line := s.Page(pageNum).Line(lineNum)
			value, err := strconv.ParseFloat(scan.Arg, 64)
			if err != nil {
				log.Fatalf("line %d: %v", scan.Line, err)
			}
			line.Time = value
		} else if scan.Prefixes.Match("Page") {
			pageNum := scan.Prefixes[0].Index
			page := s.Page(pageNum)
			if err := page.Process(scan.Command, scan.Arg); err != nil {
				log.Fatalf("line %d: %v", scan.Line, err)
			}
		} else if scan.Prefixes.Match() {
			if err := s.Process(scan.Command, scan.Arg); err != nil {
				log.Fatalf("line %d: %v", scan.Line, err)
			}
		}
	}
	if scan.Err != nil {
		log.Fatal(scan.Err)
	}

	return s
}

func printStory(s *story.Story) {
	t := template.Must(template.New("tmpl").Parse(tmpl))
	t.Execute(os.Stdout, s)
}

// TODO(jpj): This is just an example of formatting the HTML based
// on the story data model. Need to examine the existing script and
// update tmpl accordingly.
const tmpl = `<!doctype html>
<html>
<head>
</head>
<body>
<div>Name: {{.Name}}</div>
{{range .Pages -}}
<div>
	<div>Page {{.Number}}</div>
	{{range .Lines}}
		<div>
			{{range .Texts -}}
			<div>{{.}}</div>
			{{- end -}}
		</div>
	{{end}}
</div>
{{- end -}}
</body>
</html>
`
