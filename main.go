package main

import (
	"text/template"
	"io"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"github.com/jjeffery/vt-motoli/scanner"
	"github.com/jjeffery/vt-motoli/story"
)

func main() {
	log.SetFlags(0)

	// TODO(jpj): start with very simple command line, can expand later
	if len(os.Args) != 2 {
		log.Fatal("usage: vt-motoli <file>")
	}

	makeStory(os.Args[1])

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)




	go func() {
		log.Println("Listening...")
		http.ListenAndServe(":3000", nil)
	}()
	watchForFileChanges(os.Args[1])

}

func makeStory(sourceFilename string){
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		log.Fatal(err)
	}
	story := scanStory(sourceFile)

	dir, filename  := path.Split(sourceFilename)
	resultFilename := path.Join(dir,filename[:len(filename)-len(path.Ext(filename))] + ".html")
	resultFile, err := os.Create(resultFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()
	printStory(story, resultFile)
}

func scanStory(r io.Reader) *story.Story {
	s := story.New()
	scan := scanner.New(r)
	for scan.Scan() {
		if scan.Command.Matches("Page", "Line") {
			pageNum := scan.Command[0].Index
			lineNum := scan.Command[1].Index
			continuationNum := scan.Command[1].Cont
			s.Page(pageNum).Line(lineNum, true).Segments[continuationNum] = scan.Arg
		} else if scan.Command.Matches("Page", "Text") {
			pageNum := scan.Command[0].Index
			textNum := scan.Command[1].Index
			continuationNum := scan.Command[1].Cont
			s.Page(pageNum).Line(textNum, false).Segments[continuationNum] = scan.Arg
		} else if scan.Command.Matches("Page", "Time") {
			pageNum := scan.Command[0].Index
			lineNum := scan.Command[1].Index
			s.Page(pageNum).Line(lineNum, true).Time = scan.Arg
		} else if scan.Command.Matches("Page", "Pic") {
			pageNum := scan.Command[0].Index
			if (scan.Arg != ""){
				s.Page(pageNum).Image = fmt.Sprintf("../../common/%s.jpg", scan.Arg)
			}
		} else if scan.Command.Matches("StoryName") {
			s.Name = scan.Arg
		} else if scan.Command.Matches("Format") {
			s.Format = scan.Arg
		} else if scan.Command.Matches("MaxPages") ||
			scan.Command.Matches("MaxLines") ||
			scan.Command.Matches("MaxCont") {
			// do nothing: not needed anymore
		} else if scan.Command.Matches("ScaleSide") {
			s.ScaleSide = scan.Arg
		} else if scan.Command.Matches("ScaleTop") {
			s.ScaleTop = scan.Arg
		} else if scan.Command.Matches("Pause") {
			s.Pause = scan.Arg
		} else {
			log.Fatalf("line %d: unknown command", scan.Line)
		}
	}
	if scan.Err != nil {
		log.Fatal(scan.Err)
	}

	return s
}

func floatArg(scan *scanner.Scanner) float64 {
	v, err := strconv.ParseFloat(scan.Arg, 64)
	if err != nil {
		log.Fatalf("line %d: %v", scan.Line, err)
	}
	return v
}

func intArg(scan *scanner.Scanner) int {
	v, err := strconv.Atoi(scan.Arg)
	if err != nil {
		log.Fatalf("line %d: %v", scan.Line, err)
	}
	return v
}

func substitute(s string) string {
	s = strings.Replace(s, "|", "</span><span class=\"pause\">|</span><span>", -1)
	if (s == "&nil") {
		s = "<br />"
	}
	return s
}

func watchForFileChanges(sourceFile string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					makeStory(sourceFile)

				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func printStory(s *story.Story, outputFile *os.File) {
	tmpl, err := template.New("").Parse(`{{ template "story.html" . }}`)
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.ParseFiles(
		"../../../templates/story.html",
		"../../../templates/page.html",
		"../../../templates/text.html",
		"../../../templates/line.html",
		"../../../templates/single_segment.html",
		"../../../templates/simple_no_audio_line.html",
		"../../../templates/segment.html",
		"../../../templates/text.html")
	if err != nil {
		panic(err)
	}

	for k, v := range s.Pages {
		  for k1, v1 := range v.Lines {
			  for k3, v3 := range v1.Segments {
				  //if(v3=="&nil" && len(v1.Segments)==1){
					//  s.Pages[k].Lines[k1].IsLineType = true
				  //}
				  s.Pages[k].Lines[k1].Segments[k3] = substitute(v3)
			  }

		  }
		//for k2, v2 := range v.Texts {
		//	for k3, v3 := range v2.Segments {
		//		s.Pages[k].Texts[k2].Segments[k3] = substitute(v3)
		//	}
		//}
	}
	err = tmpl.Execute(outputFile, s)
	if err != nil {
		panic(err)
	}

}
