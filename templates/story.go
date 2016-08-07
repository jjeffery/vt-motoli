package templates
import (
    "text/template"
)
func AddStory(t *template.Template) {
    t.New("story.html").Parse(`
    <!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>
        {{.Name}}
    </title>
    <link rel="stylesheet" href="../../common/storyStyles.css">
    <script src="../../common/jquery-3.1.0.min.js"></script>
    <script src="../../common/storyScripts.js"></script>
</head>
<body onload="initialise()" onresize="limitPictureSize()">
<p id="loadmsg">Loading the story...</p>
<div class="{{.Format}}" id="story">

    <!-- Data for this story -->
    <p class="data">
        <span id="platform"></span>
        <span id="onloads">0</span>
        <span id="pScaleSide">{{.ScaleSide}}</span>
        <span id="pScaleTop">{{.ScaleTop}}</span>
        <span id="barPause">{{.Pause}}</span>
        <span id="backColor">white</span>
        <span id="diagnostic">off</span>
    </p>
    {{range .Pages}}
        {{ template "page.html" . }}
    {{end}}
    <!-- PAGE IDENTIFICATION SECTION - DO NOT CHANGE -->
    <div class="footer underpic">
        <strong><span class="title">{{.Name}}</span></strong>
        <span class="pgnum">page #</span>
    </div> <!-- END PAGE IDENTIFICATION SECTION -->

    <div>
        <audio id="AudioPlayer"></audio>
    </div>
</div>
<!-- END OF STORY -->
</body>
</html>
`)
}