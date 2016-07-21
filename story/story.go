package story

import (
	"fmt"
	"regexp"
	"strconv"
)

type Story struct {
	Name      string
	Format    string // "side", ?
	MaxPages  int
	MaxLines  int
	MaxCont   int
	ScaleSide int
	ScaleTop  int
	Pause     int
	Pages     []*Page
}

type Page struct {
	Number int
	Lines  []*Line
	Texts  []string
	Pic    string
}

func newPage(num int) *Page {
	return &Page{
		Number: num,
	}
}

func (p *Page) GetLine(lineNum int) (*Line, error) {
	lineNum = lineNum - 1
	if len(p.Lines) == lineNum {
		p.Lines = append(p.Lines, newLine())
	}
	if len(p.Lines) > lineNum {
		return p.Lines[lineNum], nil
	}
	return nil, fmt.Errorf("line %d missing", len(p.Lines)+1)
}

func (page *Page) SetText(num int, text string) error {
	if len(page.Texts) == num {
		page.Texts = append(page.Texts, text)
		return nil
	}
	if len(page.Texts) > num {
		page.Texts[num] = text
		return nil
	}
	return fmt.Errorf("invalid index: %d", num+1)
}

func (page *Page) Process(command, arg string) error {
	switch command {
	case "Pic":
		page.Pic = arg
	default:
		return fmt.Errorf("invalid command")
	}
	return nil
}

type Line struct {
	Texts []string
	Time  float64
}

func (line *Line) SetText(num int, text string) error {
	if len(line.Texts) == num {
		line.Texts = append(line.Texts, text)
		return nil
	}
	if len(line.Texts) > num {
		line.Texts[num] = text
		return nil
	}
	return fmt.Errorf("invalid index: %d", num+1)
}

func newLine() *Line {
	return &Line{}
}

func New() *Story {
	// TODO(jpj): include default values here
	return &Story{}
}

var (
	storyCmdMap = map[string]func(*Story, string, string) error{
		"StoryName": handleStoryName,
		"Format":    handleFormat,
		"MaxPages":  handleMaxPages,
		"MaxLines":  handleMaxLines,
	}

	pageCmdRE = regexp.MustCompile(`^Page([0-9]+)(.*)$`)
)

func pageCommand(command string) (pageNum int, pageCmd string, ok bool) {
	var err error
	if submatches := pageCmdRE.FindStringSubmatch(command); len(submatches) > 0 {
		pageNum, err = strconv.Atoi(submatches[1])
		if err != nil {
			// this will only happen if the number is too big
			// eg "Page999999999999999999999999"
			return 0, "", false
		}
		pageCmd = submatches[2]
		return pageNum, pageCmd, true
	}
	return 0, "", false
}

func (s *Story) GetPage(pageNum int) (*Page, error) {
	if len(s.Pages) == pageNum {
		s.Pages = append(s.Pages, newPage(pageNum))
	}
	if len(s.Pages) > pageNum {
		return s.Pages[pageNum], nil
	}
	return nil, fmt.Errorf("page %d missing", len(s.Pages))
}

func (s *Story) Process(command string, arg string) error {
	if fn := storyCmdMap[command]; fn != nil {
		return fn(s, command, arg)
	}
	return nil
}

func handleStoryName(s *Story, command, arg string) error {
	s.Name = arg
	return nil
}

func handleFormat(s *Story, command, arg string) error {
	switch arg {
	case "side", "top":
		s.Format = arg
		return nil
	}
	return fmt.Errorf("unknown format: %q", arg)
}

func handleMaxPages(s *Story, command, arg string) error {
	n, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	s.MaxPages = n
	return nil
}

func handleMaxLines(s *Story, command, arg string) error {
	n, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	s.MaxLines = n
	return nil
}
