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
		if scan.Command.Matches("Page", "Line") {
			pageNum := scan.Command[0].Index
			lineNum := scan.Command[1].Index
			continuationNum := scan.Command[1].Cont
			s.Page(pageNum).Line(lineNum).Texts[continuationNum] = scan.Arg
		} else if scan.Command.Matches("Page", "Text") {
			pageNum := scan.Command[0].Index
			textNum := scan.Command[1].Index
			s.Page(pageNum).Texts[textNum] = scan.Arg
		} else if scan.Command.Matches("Page", "Time") {
			pageNum := scan.Command[0].Index
			lineNum := scan.Command[1].Index
			s.Page(pageNum).Line(lineNum).Time = floatArg(scan)
		} else if scan.Command.Matches("Page", "Pic") {
			pageNum := scan.Command[0].Index
			s.Page(pageNum).Image = scan.Arg
		} else if scan.Command.Matches("StoryName") {
			s.Name = scan.Arg
		} else if scan.Command.Matches("Format") {
			s.Format = scan.Arg
		} else if scan.Command.Matches("MaxPages") ||
			scan.Command.Matches("MaxLines") ||
			scan.Command.Matches("MaxCont") {
			// do nothing: not needed anymore
		} else if scan.Command.Matches("ScaleSide") {
			s.ScaleSide = intArg(scan)
		} else if scan.Command.Matches("ScaleTop") {
			s.ScaleTop = intArg(scan)
		} else if scan.Command.Matches("Pause") {
			s.Pause = intArg(scan)
		} else {
			log.Fatalf("line %d: unknown command", scan.Line)
		}
	}
	if scan.Err != nil {
		log.Fatal(scan.Err)
	}

	return s
}

func floatArg(scan *scanner.Scanner) float64 {
	v, err := strconv.ParseFloat(scan.Arg, 64)
	if err != nil {
		log.Fatalf("line %d: %v", scan.Line, err)
	}
	return v
}

func intArg(scan *scanner.Scanner) int {
	v, err := strconv.Atoi(scan.Arg)
	if err != nil {
		log.Fatalf("line %d: %v", scan.Line, err)
	}
	return v
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
