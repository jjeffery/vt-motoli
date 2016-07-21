package main

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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

// Regular expressions for parsing the line contents
var (
	// regular expressions indicating a comment line
	commentREs = []*regexp.Regexp{
		regexp.MustCompile(`^#(#|[^#]+\s)`), // eg "### comment"
		regexp.MustCompile(`^#[^#]\s`),      // eg "#This is a comment"
		regexp.MustCompile(`^\s*<`),         // eg "<Identifier> <-- start of line etc"
	}

	commandRE = regexp.MustCompile(`^#([A-Za-z0-9-]+)#\s+(.*)$`)
)

func scanStory(r io.Reader) *story.Story {
	s := story.New()
	scan := scanner.New(r)
	for scan.Scan() {
		if scan.Err != nil {
			log.Fatal(scan.Err)
		}
		if scan.Prefixes.Match("Page", "Line") {
			pageNum := scan.Prefixes[0].Index
			lineNum := scan.Prefixes[1].Index
			line := s.Page(pageNum).Line(lineNum)
			line.SetText(scan.Prefixes[1].Cont, scan.Arg)
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
	return s
}

func scanStory_1(r io.Reader) *story.Story {
	s := story.New()
	scanner := bufio.NewScanner(r)

	for line := 1; scanner.Scan(); line++ {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		if text == "" {
			debug.Printf("line %d: blank\n", line)
			continue
		}

		if isComment(text) {
			debug.Printf("line %d: comment\n", line)
			continue
		}

		if cmd := commandRE.FindStringSubmatch(text); len(cmd) > 0 {
			debug.Printf("line %d: %q %q\n", line, cmd[1], cmd[2])
			if err := s.Process(cmd[1], cmd[2]); err != nil {
				log.Fatal("line %d: %v", line, err)
			}
			continue
		}

		// TODO(jpj):
		log.Fatal("line %d: unknown: %q\n", line, text)
	}

	return s
}

func isComment(text string) bool {
	for _, re := range commentREs {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}

func printStory(s *story.Story) {
	t := template.Must(template.New("tmpl").Parse(tmpl))
	t.Execute(os.Stdout, s)
}

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
