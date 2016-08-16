package main

//go:generate go run assets_generate.go

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/braintree/manners"
	"github.com/fsnotify/fsnotify"
	"github.com/jjeffery/vt-motoli/graceful"
	"github.com/jjeffery/vt-motoli/scanner"
	"github.com/jjeffery/vt-motoli/story"
	"github.com/jjeffery/vt-motoli/templates"
	"github.com/jjeffery/vt-motoli/touch"
)

var developmentMode bool

var (
	// Version gets set during the formal build.
	Version     = "development"
	showVersion bool
)

func init() {
	flag.BoolVar(&showVersion, "version", false, "show version")
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	debug.Printf("command line args: %v", flag.Args())
	switch len(flag.Args()) {
	case 0:
		break
	case 1:
		debug.Printf("changing directory to %q", flag.Arg(0))
		if err := os.Chdir(flag.Arg(0)); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("usage: %s [ directory ]", os.Args[0])
	}

	var wg sync.WaitGroup

	regenerateObsoleteHtmls(".", false)

	wg.Add(1)
	go func() {
		webServer()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		watchForFileChanges("./")
		wg.Done()
	}()

	wg.Wait()
}

func webServer() {
	staticFileServer := http.FileServer(http.Dir("."))
	assetServer := http.FileServer(assets)
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL.Path)
		if _, err := os.Stat("./" + req.URL.Path); err == nil {
			// path/to/whatever exists
			staticFileServer.ServeHTTP(w, req)
		} else {
			assetServer.ServeHTTP(w, req)
		}
	})

	graceful.OnShutdown(func() { manners.Close() })
	log.Println("Listening...")
	manners.ListenAndServe(":3000", http.DefaultServeMux)
	log.Println("web server stopped")
}

func makeStory(sourceFilename string) {
	sourceFile, err := os.Open(sourceFilename)
	if err != nil {
		log.Fatal(err)
	}
	story := scanStory(sourceFile)

	resultFilename := getCorrespondingHtmlFilename(sourceFilename)
	resultFile, err := os.Create(resultFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer resultFile.Close()
	printStory(story, resultFile)
}

func scanStory(r io.Reader) *story.Story {
	s := story.New(developmentMode)
	scan := scanner.New(r)
	for scan.Scan() {
		if scan.Err != nil {
			page := s.CurrentPage()
			page.Errors = append(page.Errors, scan.Err.Error())
			scan.Err = nil
		} else if scan.Command.Matches("Page", "Line") {
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
			if scan.Arg != "" {
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
	if s == "&nil" {
		s = "<br />"
	}
	return s
}

func isMotoLiSourceFile(filename string) bool {
	debug.Printf("start isMotoLiSourceFile(%q)", filename)
	defer debug.Printf("end isMotoLiSourceFile(%q)", filename)

	if strings.ToLower(filepath.Ext(filename)) != ".txt" {
		debug.Printf("%q: does not end with '.txt'", filename)
		return false
	}
	pageRegex := regexp.MustCompile(`^#Page[0-9]+`)
	sourceFile, err := os.Open(filename)
	if err != nil {
		if err == os.ErrNotExist {
			debug.Printf("%q: does not exist", filename)
			return false
		}
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if pageRegex.MatchString(line) {
			debug.Printf("%q: matches", filename)
			return true
		}
	}
	debug.Printf("%q: does not match")
	return false
}

var watchedDirectories map[string]bool = map[string]bool{}

func watchForFileChanges(baseDirectory string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("cannot create FS watcher:", err)
		graceful.Shutdown()
		return
	}
	defer watcher.Close()

	watchSubdirectories(watcher, baseDirectory)

	for {
		select {
		case <-graceful.Done:
			log.Println("file change watch stopped")
			return
		case event := <-watcher.Events:
			log.Println("modified file:", event.Name)
			if isMotoLiSourceFile(event.Name) {
				log.Println("yippee")
				makeStory(event.Name)
			} else if path.Ext(event.Name) != ".html" {
				touch.RecursiveTouchHtml(".")
			}
		case err := <-watcher.Errors:
			log.Println("file watcher error:", err)
		}
	}
}

func watchSubdirectories(w *fsnotify.Watcher, directory string) {
	watchedDirectories[directory] = true
	err := w.Add(directory)
	err = filepath.Walk(directory, func(filename string, f os.FileInfo, err error) error {
		if filename != directory && f.IsDir() {
			watchSubdirectories(w, filename)
		}
		return nil
	})
	if err != nil {
		log.Printf("cannot walk directory %s: %v", directory, err)
		graceful.Shutdown()
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
	if err := tmpl.Execute(outputFile, s); err != nil {
		panic(err)
	}
}

func regenerateObsoleteHtmls(parentPath string, isFinalReleaseMode bool) {
	filepath.Walk(parentPath, func(path string, textInfo os.FileInfo, err error) error {
		if path != parentPath {

			if textInfo.IsDir() {
				regenerateObsoleteHtmls(path, isFinalReleaseMode)
			} else if isMotoLiSourceFile(path) {
				htmlFilename := getCorrespondingHtmlFilename(path)
				htmlInfo, err := os.Stat(htmlFilename)
				if err != nil || (textInfo.ModTime().After(htmlInfo.ModTime()) || isFinalReleaseMode != getIsFinalDevelopmentMode(htmlFilename)) {
					makeStory(path)
				}
			}
		}
		return nil
	})
}

func getCorrespondingHtmlFilename(sourceFilename string) string {
	dir, filename := path.Split(sourceFilename)
	return path.Join(dir, filename[:len(filename)-len(path.Ext(filename))]+".html")
}

func getIsFinalDevelopmentMode(htmlFilename string) bool {
	debug.Printf("start getIsFinalReleaseMode(%q)", htmlFilename)
	defer debug.Printf("end getIsFinalReleaseMode(%q)", htmlFilename)

	pageRegex := regexp.MustCompile(`^<!-- vt-motoli development mode --!>$`)
	sourceFile, err := os.Open(htmlFilename)
	if err != nil {
		if err == os.ErrNotExist {
			debug.Printf("%q: does not exist", htmlFilename)
			return false
		}
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(sourceFile)
	for scanner.Scan() {
		line := scanner.Text()
		if pageRegex.MatchString(line) {
			debug.Printf("%q: matches", htmlFilename)
			return true
		}
	}
	debug.Printf("%q: does not match")
	return false
}
