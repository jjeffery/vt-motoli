package scanner

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions for parsing the line contents
var (
	// regular expressions indicating a comment line
	commentREs = []*regexp.Regexp{
		regexp.MustCompile(`^#(#|[^#]+\s)`), // eg "### comment"
		regexp.MustCompile(`^#[^#]\s`),      // eg "#This is a comment"
		regexp.MustCompile(`^\s*<`),         // eg "<Identifier> <-- start of line etc"
	}

	commandRE = regexp.MustCompile(`^#([A-Za-z0-9-]+)#\s+(.*)$`)

	prefixREs = []*regexp.Regexp{
		regexp.MustCompile(`^(Page)([0-9]+)`),
		regexp.MustCompile(`^(Line)([0-9]+)-([0-9]+)`),
		regexp.MustCompile(`^(Line)([0-9]+)`),
		regexp.MustCompile(`^(Text)([0-9]+)-([0-9]+)`),
		regexp.MustCompile(`^(Text)([0-9]+)`),
		regexp.MustCompile(`^(Time)([0-9]+)`),
	}
)

type Prefix struct {
	Name  string
	Index int
	Cont  int
}

type Prefixes []Prefix

func (p Prefixes) Next() Prefixes {
	if len(p) <= 1 {
		return nil
	}
	return p[1:]
}

func (p Prefixes) Match(names ...string) bool {
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

func (p Prefixes) Index() int {
	if len(p) == 0 {
		return 0
	}
	return p[0].Index
}

func (p Prefixes) Cont() int {
	if len(p) == 0 {
		return 0
	}
	return p[0].Cont
}

type Scanner struct {
	Prefixes Prefixes
	Command  string
	Arg      string
	Err      error
	Line     int

	scanner *bufio.Scanner
}

func New(r io.Reader) *Scanner {
	return &Scanner{
		scanner: bufio.NewScanner(r),
	}
}

func (s *Scanner) Scan() bool {
	s.Prefixes = nil
	s.Command = ""
	s.Arg = ""
	for s.scanner.Scan() {
		text := s.scanner.Text()
		text = strings.TrimSpace(text)
		s.Line++
		debug.Println("Line", s.Line)

		if text == "" || isComment(text) {
			continue
		}

		if cmd := commandRE.FindStringSubmatch(text); len(cmd) > 0 {
			s.Arg = cmd[2]
			s.Prefixes, s.Command = stripPrefixes(cmd[1])
			return true
		}

		s.Err = fmt.Errorf("line %d: %q", s.Line, text)
		return false
	}
	return false
}

func stripPrefixes(command string) ([]Prefix, string) {
	debug.Printf("strip prefixes: %q\n", command)
	var prefixes []Prefix

	for _, re := range prefixREs {
		debug.Println("command=", command)
		debug.Printf("re=%v\n", re)
		subs := re.FindStringSubmatch(command)
		if len(subs) == 0 {
			continue
		}
		command = command[len(subs[0]):]
		prefix := Prefix{
			Name: subs[1],
		}
		prefix.Index, _ = strconv.Atoi(subs[2])
		if len(subs) > 3 {
			prefix.Cont, _ = strconv.Atoi(subs[3])
		}
		prefixes = append(prefixes, prefix)
	}

	debug.Printf("strip prefixes: %v, %q\n", prefixes, command)
	return prefixes, command
}

func isComment(text string) bool {
	for _, re := range commentREs {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}
