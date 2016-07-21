package story

import (
	"fmt"
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
	Pages     []Page
}

type Page struct {
	Lines []Line
	Texts []string
}

type Line struct {
	Texts []string
	Time  float64
}

func New() *Story {
	// TODO(jpj): include default values here
	return &Story{}
}

var cmdMap = map[string]func(*Story, string, string) error{
	"StoryName": handleStoryName,
	"Format":    handleFormat,
	"MaxPages":  handleMaxPages,
}

func (s *Story) Process(command string, arg string) error {
	if fn := cmdMap[command]; fn != nil {
		return fn(s, command, arg)
	}
	return nil
}

func stringArg(fn func(*Story, string)) func(*Story, string, string) error {
	return func(s *Story, command, arg string) error {
		fn(s, arg)
		return nil
	}
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
