package story

import (
	"fmt"
)

type Page struct {
	Number int
	Lines  map[int]*Line
	Texts  map[int]string
	Image  string
}

func newPage(num int) *Page {
	return &Page{
		Number: num,
		Lines:  make(map[int]*Line),
		Texts:  make(map[int]string),
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

func (page *Page) SetText(num int, text string) {
	page.Texts[num] = text
}

func (page *Page) Process(command, arg string) error {
	switch command {
	case "Pic":
		page.Image = arg
	default:
		return fmt.Errorf("invalid command")
	}
	return nil
}
