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
	Texts  map[int]*Text
	Image  string
	Story  *Story
}

func (story *Story)newPage(num int) *Page {
	return &Page{
		Number: num,
		Lines:  make(map[int]*Line),
		Texts:  make(map[int]*Text),
		Image:  fmt.Sprintf("p%d", num), // TODO: find out naming convention for pictures
		Story: story,
	}
}

func (p *Page) Line(lineNum int) *Line {
	line := p.Lines[lineNum]
	if line == nil {
		line = p.newLine(lineNum)
		p.Lines[lineNum] = line
	}
	return line
}

func (p *Page) Text(textNum int) *Text {
	text := p.Texts[textNum]
	if text == nil {
		text = p.newText(textNum)
		p.Texts[textNum] = text
	}
	return text
}

type Line struct {
	Number   int
	Segments map[int]string
	Time     string
	Page *Page
}

func (line *Line)OnlyOneSegment() bool {
	return len(line.Segments)==1
}

func (page *Page)newLine(num int) *Line {
	return &Line{
		Number:   num,
		Segments: make(map[int]string),
		Page: page,
	}
}

type Text struct {
	Number   int
	Segments map[int]string
	Page *Page
}
func (text *Text)OnlyOneSegment() bool {
	return len(text.Segments)==1
}

func (page *Page)newText(num int) *Text {
	return &Text{
		Number:   num,
		Segments: make(map[int]string),
		Page: page,
	}
}
