package story

type Line struct {
	Number int
	Texts  map[int]string
	Time   float64
}

func newLine(num int) *Line {
	return &Line{
		Number: num,
		Texts:  make(map[int]string),
	}
}

func (line *Line) SetText(num int, text string) {
	line.Texts[num] = text
}
