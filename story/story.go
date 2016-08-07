package story

import (
	"fmt"
)

type Story struct {
	Name      string
	Format    string
	ScaleSide string
	ScaleTop  string
	Pause     string
	BackColor string
	Debug     string
	Pages     map[int]*Page
}

func New() *Story {
	return &Story{
		Pages: make(map[int]*Page),
		Format: "side",
		ScaleSide : "40",
		ScaleTop : "60",
		Pause : "300",
		BackColor : "white",
		Debug : "off",
	}
}

func (s *Story) Page(pageNum int) *Page {
	page := s.Pages[pageNum]
	if page == nil {
		page = s.newPage(pageNum)
		s.Pages[pageNum] = page
	}
	return page
}

type Page struct {
	Number int
	Lines  map[int]*Line
	Image  string
	Story  *Story
}

func (story *Story)newPage(num int) *Page {
	return &Page{
		Number: num,
		Lines:  make(map[int]*Line),
		Image:  fmt.Sprintf("p%dpic.jpg", num), // TODO: find out naming convention for pictures
		Story: story,
	}
}

func (p *Page) Line(lineNum int, isLineType bool) *Line {
	line := p.Lines[lineNum]
	if line == nil {
		line = p.newLine(lineNum, isLineType)
		p.Lines[lineNum] = line
	}
	return line
}

type Line struct {
	Number   int
	Segments map[int]string
	Time     string
	Page *Page
	IsLineType bool // as opposed to a text type
	Lang   string
}

func (line *Line)OnlyOneSegment() bool {
	return len(line.Segments)==1
}

func (page *Page)newLine(num int, isLineType bool) *Line {
	return &Line{
		Number:   num,
		Segments: make(map[int]string),
		Page: page,
		IsLineType: isLineType,
	}
}
