package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

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

	fmt.Printf("%+v\n", s)
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
