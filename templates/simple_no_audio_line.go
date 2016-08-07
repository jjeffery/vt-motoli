package templates
import (
"text/template"
)
func AddSimpleNoAudioLine(t *template.Template){
	t.New("simple_no_audio_line.html").Parse(`
		<!-- Text line {{.Number}} --><tr><td><span>
		{{ if .OnlyOneSegment}}
		    {{ index .Segments 0 }}
		</span><br />
		    {{ else }}
		    {{range .Segments}}
		    {{ template "segment.html" . }}
		    {{ end }}
		    </span></td></tr>
		{{ end }}
`)
}
