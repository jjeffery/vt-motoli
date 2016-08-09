package main

import (
	"bufio"
	"text/template"
	"io"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"regexp"
	"strings"
	_ "github.com/jteeuwen/go-bindata"
	"github.com/jjeffery/vt-motoli/scanner"
	"github.com/jjeffery/vt-motoli/story"
	"github.com/jjeffery/vt-motoli/templates"
	"github.com/jjeffery/vt-motoli/touch"
	"github.com/kardianos/osext"
	//"github.com/spkg/bom"
	"github.com/spkg/zipfs"
)

//func MotoLiHandler(w http.ResponseWriter, r *http.Request) {
//	fs := http.FileServer(http.Dir("."))
//	fs(w, r)
//	//w.Write([]byte("Hello World"))
//}

func main() {
	filename, _ := osext.Executable()
	fmt.Println(filename)

	log.SetFlags(0)

	go func() {

		staticFileServer := http.FileServer(http.Dir("."))
		zipFileSystem, err := zipfs.New("resources.zip")
		if err != nil {
			log.Fatal(err)
		}
		zipFileServer := zipfs.FileServer(zipFileSystem)
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			log.Println(req.URL.Path)
			if _, err := os.Stat("./"+req.URL.Path); err == nil {
				// path/to/whatever exists
				staticFileServer.ServeHTTP(w,req)
			} else {
				zipFileServer.ServeHTTP(w,req)
			}

		})
		log.Println("Listening...")
		http.ListenAndServe(":3000", nil)
	}()
	watchForFileChanges("./")
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
		} else if scan.Command.Matches("Page", "Lang") {
			pageNum := scan.Command[0].Index
			textNum := scan.Command[1].Index
			s.Page(pageNum).Line(textNum, false).Lang = scan.Arg

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

func isMotoLiSourceFile(filename string) bool {
	pageRegex := regexp.MustCompile(`^#Page[0-9]+`)
	sourceFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if(pageRegex.MatchString(line)){
			return true;
		}
	}
	return false;

}

var watchedDirectories map[string]bool = map[string]bool{}
func watchForFileChanges(baseDirectory string) {
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
				log.Println("modified file:", event.Name)
				if(isMotoLiSourceFile(event.Name)){
					log.Println("yippee")
					makeStory(event.Name)
				}else if(path.Ext(event.Name)!=".html"){
					touch.RecursiveTouchHtml(".")
				}

			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	watchSubdirectories(watcher, baseDirectory)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func watchSubdirectories(w *fsnotify.Watcher, directory string){
	watchedDirectories[directory] = true
	err := w.Add(directory)
	err = filepath.Walk(directory, func(filename string, f os.FileInfo, err error) error {
		if(filename!=directory && f.IsDir()){
			watchSubdirectories(w, filename)
		}
		return nil
})
	if err != nil {
		log.Fatal(err)
	}
}

func printStory(s *story.Story, outputFile *os.File) {
	tmpl, err := template.New("").Parse(`{{ template "story.html" . }}`)
	if err != nil {
		panic(err)
	}
	templates.AddStory(tmpl)
	templates.AddPage(tmpl)
	templates.AddText(tmpl)
	templates.AddLine(tmpl)
	templates.AddSingleSegment(tmpl)
	templates.AddSimpleNoAudioLine(tmpl)
	templates.AddSegment(tmpl)
	templates.AddText(tmpl)
	templates.AddSimpleNoAudioLine(tmpl)

	for k, v := range s.Pages {
		  for k1, v1 := range v.Lines {
			  for k3, v3 := range v1.Segments {
				  s.Pages[k].Lines[k1].Segments[k3] = substitute(v3)
			  }

		  }
	}
	err = tmpl.Execute(outputFile, s)
	if err != nil {
		panic(err)
	}

}
