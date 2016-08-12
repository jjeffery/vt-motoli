package templates

import (
	"text/template"
)

func AddSingleSegment(t *template.Template) {
	t.New("single_segment.html").Parse(`
	<span class="line"><span>{{ . }}</span><progress value="0" max="100"></progress></td>
`)
}
