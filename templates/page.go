package templates
import (
    "text/template"
)
func AddPage(t *template.Template) {
    t.New("page.html").Parse(`
    <!-- PAGE {{.Number}} -->
<div class="page" id="p{{.Number}}">
    <div class="pict">
        <div class="underpic"> <!-- Format for a vertically oriented page -->
            <table><tr>
                <td><img src="../../common/back.png" class="btnback" onclick="pageBack()"></td>
                <td><img src="{{.Image}}" class="mainpic"></td>
                <td><img src="../../common/fwd.png" class="btnfwd" onclick="pageFwd()"></td>
            </tr></table>
        </div>

        <div class="sidepic"> <!-- Format for a horizontally oriented page -->
            <img src="{{.Image}}" class="mainpic">
            <table><tr>
                <td><img src="../../common/back.png" class="btnback" onclick="pageBack()"></td>
                <td><span class="title">{{.Story.Name}}</span><br />
                    <span class="pgnum">page #</span></td>
                <td><img src="../../common/fwd.png" class="btnfwd" onclick="pageFwd()"></td>
            </tr></table>
        </div>
    </div>

<div class="text">
    <table class="texttab">
        {{range .Lines}}
            {{ if .Time }}
                {{ template "line.html" . }}
            {{ else }}
                {{ if .IsLineType }}
                    {{ template "simple_no_audio_line.html" . }}
                {{ else }}
                    {{ template "text.html" . }}
                {{ end }}
            {{end}}
        {{ end }}
    </table>
</div>
</div>
<!-- END PAGE {{.Number}} -->
`)
}