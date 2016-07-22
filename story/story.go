package story

import (
	"fmt"
)

type Story struct {
	Name      string
	Format    string
	ScaleSide int
	ScaleTop  int
	Pause     int
	Pages     map[int]*Page
}

func New() *Story {
	return &Story{
		Pages: make(map[int]*Page),
	}
}

func (s *Story) Page(pageNum int) *Page {
	page := s.Pages[pageNum]
	if page == nil {
		page = newPage(pageNum)
		s.Pages[pageNum] = page
	}
	return page
}

type Page struct {
	Number int
	Lines  map[int]*Line
	Texts  map[int]*Text
	Image  string
}

func newPage(num int) *Page {
	return &Page{
		Number: num,
		Lines:  make(map[int]*Line),
		Texts:  make(map[int]*Text),
		Image:  fmt.Sprintf("p%d", num), // TODO: find out naming convention for pictures
	}
}

func (p *Page) Line(lineNum int) *Line {
	line := p.Lines[lineNum]
	if line == nil {
		line = newLine(lineNum)
		p.Lines[lineNum] = line
	}
	return line
}

func (p *Page) Text(textNum int) *Text {
	text := p.Texts[textNum]
	if text == nil {
		text = newText(textNum)
		p.Texts[textNum] = text
	}
	return text
}

type Line struct {
	Number   int
	Segments map[int]string
	Time     float64
}

func newLine(num int) *Line {
	return &Line{
		Number:   num,
		Segments: make(map[int]string),
	}
}

type Text struct {
	Number   int
	Segments map[int]string
}

func newText(num int) *Text {
	return &Text{
		Number:   num,
		Segments: make(map[int]string),
	}
}
