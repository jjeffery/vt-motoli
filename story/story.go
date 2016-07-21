package story

import (
	"fmt"
	"regexp"
	"strconv"
)

type Story struct {
	Name      string
	Format    string // "side", ?
	ScaleSide int
	ScaleTop  int
	Pause     int
	Pages     map[int]*Page
}

func New() *Story {
	// TODO(jpj): include default values here
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

func (s *Story) Process(command string, arg string) error {
	if fn := storyCmdMap[command]; fn != nil {
		return fn(s, command, arg)
	}
	return nil
}

var (
	storyCmdMap = map[string]func(*Story, string, string) error{
		"StoryName": handleStoryName,
		"Format":    handleFormat,
		"MaxPages":  handleIgnore,
		"MaxLines":  handleIgnore,
		"ScaleSide": handleScaleSide,
		"ScaleTop":  handleScaleTop,
		"Pause":     handlePause,
	}

	pageCmdRE = regexp.MustCompile(`^Page([0-9]+)(.*)$`)
)

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

func handleIgnore(s *Story, command, arg string) error {
	// ignore, not needed
	return nil
}

func handleScaleSide(s *Story, command, arg string) error {
	v, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	s.ScaleSide = v
	return nil
}

func handleScaleTop(s *Story, command, arg string) error {
	v, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	s.ScaleTop = v
	return nil
}

func handlePause(s *Story, command, arg string) error {
	v, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}
	s.Pause = v
	return nil
}
