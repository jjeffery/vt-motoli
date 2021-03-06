package templates

import (
	"text/template"
)

func AddLine(t *template.Template) {
	t.New("line.html").Parse(`
{{ if .OnlyOneSegment }}
<!-- Text line {{.Number}} -->
<tr{{ if .Lang}} class="l-{{.Lang}} {{.Lang}}"{{ end}}><td class="textLine">
    {{ template "single_segment.html" index .Segments 0}}
{{ else }}
<tr{{ if .Lang}} class="l-{{.Lang}} {{.Lang}}"{{ end}}><!-- Text line {{.Number}} -->
    <td class="textLine wrap">
{{range .Segments}}<span class="line"><span>{{.}}</span><br /><progress value="0"></progress><br />
        {{end}}
    </td>
{{ end }}
    <!-- Audio file, time in secs, page ID, line No. -->
    <td class="button"><img src="../../common/play.png" onclick="playAudio('p{{.Page.Number}}sound{{.Number}}.mp3', '{{.Time}}', 'p{{.Page.Number}}', '{{.Number}}' )"></td>

</tr>
`)
}
