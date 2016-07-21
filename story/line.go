package story

import (
	"fmt"
)

type Line struct {
	Number int
	Texts  []string
	Time   float64
}

func newLine(num int) *Line {
	return &Line{
		Number: num,
	}
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
