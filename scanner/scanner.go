package scanner

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/spkg/bom"
)

// Regular expressions for parsing the line contents
var (
	// regular expressions indicating a comment line
	commentREs = []*regexp.Regexp{
		regexp.MustCompile(`^##`),       // eg "### comment"
		regexp.MustCompile(`^#[^#]*\s`), // eg "#This is a comment"
		regexp.MustCompile(`^<`),        // eg "<Identifier> <-- start of line etc"
	}

	commandREs = []*regexp.Regexp{
		regexp.MustCompile(`^#([A-Za-z0-9-]+)#\s+(.*)$`),   // #command# arg
		regexp.MustCompile(`^([A-Za-z0-9-]+)\s*:\s+(.*)$`), // command: arg
	}

	segmentREs = []*regexp.Regexp{
		regexp.MustCompile(`^(Page)([0-9]+)`),
		regexp.MustCompile(`^(Line)([0-9]+)-([0-9]+)`),
		regexp.MustCompile(`^(Line)([0-9]+)`),
		regexp.MustCompile(`^(Text)([0-9]+)-([0-9]+)`),
		regexp.MustCompile(`^(Text)([0-9]+)`),
		regexp.MustCompile(`^(Time)([0-9]+)`),
	}
)

// A Segment is a smaller part of a command that has meaning.
// Each segment has a name and an optional index and continuation.
// Segment examples include "Page1" and "Text1-1".
type Segment struct {
	Name  string
	Index int
	Cont  int
}

// A Command is a list of one or more segments extracted
// from a string.
type Command []Segment

func (p Command) Matches(names ...string) bool {
	if len(p) != len(names) {
		return false
	}
	for i, name := range names {
		if p[i].Name != name {
			return false
		}
	}
	return true
}

func newCommand(text string) Command {
	var segments []Segment

	for _, re := range segmentREs {
		subs := re.FindStringSubmatch(text)
		if len(subs) == 0 {
			continue
		}
		text = text[len(subs[0]):]
		segment := Segment{
			Name: subs[1],
		}
		segment.Index, _ = strconv.Atoi(subs[2])
		if len(subs) > 3 {
			segment.Cont, _ = strconv.Atoi(subs[3])
		}
		segments = append(segments, segment)
	}

	if text != "" {
		segments = append(segments, Segment{
			Name: text,
		})
	}

	return Command(segments)
}

type Scanner struct {
	Command Command
	Arg     string
	Err     error
	Line    int

	scanner *bufio.Scanner
}

func New(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(bom.NewReader(r)),
	}
}

func (s *Scanner) Scan() bool {
	s.Command = nil
	s.Arg = ""
	for s.scanner.Scan() {
		text := s.scanner.Text()
		text = strings.TrimSpace(text)
		s.Line++

		if text == "" || isComment(text) {
			continue
		}

		for _, re := range commandREs {
			if cmd := re.FindStringSubmatch(text); len(cmd) > 0 {
				s.Arg = cmd[2]
				s.Command = newCommand(cmd[1])
				return true
			}
		}

		s.Err = fmt.Errorf("line %d: %q", s.Line, text)
		return false
	}
	return false
}

func isComment(text string) bool {
	for _, re := range commentREs {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}
