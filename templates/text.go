package templates
import (
    "text/template"
)
func AddText(t *template.Template) {
    t.New("text.html").Parse(`
    {{ if .OnlyOneSegment }}
    <!-- Text line {{.Number}} --><tr><td><span>
    {{ index .Segments 0 }}
</span><br />
{{ else }}
    <!-- Text line {{.Number}} -->
    <tr><td><span>
        {{range .Segments}}
            {{ template "segment.html" . }}
        {{ end }}
    </span></td></tr>
{{ end }}
`)
}

