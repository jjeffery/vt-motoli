package story

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
