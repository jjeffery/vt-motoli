package templates
import (
	"text/template"
)
func AddSegment(t *template.Template) {
	t.New("segment.html").Parse(`
{{.}}<br />
`)
}
