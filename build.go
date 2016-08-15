// +build ignore

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var targets = []struct {
	goos   string
	goarch string
}{
	{"windows", "amd64"},
	{"windows", "386"},
	{"linux", "amd64"},
	{"linux", "386"},
	{"darwin", "amd64"},
}

func main() {
	log.SetFlags(0)
	checkGitNotDirty()
	version := getVersion()
	log.Printf("version=%q", version)

	for _, target := range targets {
		compile(target.goos, target.goarch, version)
	}
}

func checkGitNotDirty() {
	cmd := exec.Cmd{
		Path:   findExe("git"),
		Args:   []string{"git", "diff", "--quiet"},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	if err := cmd.Run(); err != nil {
		log.Fatal("uncommited changes in git repository")
	}
}

func getVersion() string {
	return fmt.Sprintf("%s/%s", time.Now().Format("2006-01-02T15:04:05-0700"), getGitRevision())
}

func getGitRevision() string {
	buf := bytes.Buffer{}
	cmd := exec.Cmd{
		Path:   findExe("git"),
		Args:   []string{"git", "rev-parse", "--short=12", "HEAD"},
		Stdout: &buf,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	if err := cmd.Run(); err != nil {
		log.Fatalf("cannot run git: %v", err)
	}

	return strings.Trim(string(buf.Bytes()), " \r\n\t")
}

func compile(goos string, goarch string, version string) {
	artifactDir := filepath.Join("artifacts", fmt.Sprintf("%s_%s", goos, goarch))
	if err := os.MkdirAll(artifactDir, 0755); err != nil {
		log.Fatal("cannot create %s: %v", artifactDir, err)
	}

	outputFile := filepath.Join(artifactDir, "vt-motoli")
	if runtime.GOOS == "windows" {
		outputFile += ".exe"
	}

	cmd := exec.Cmd{
		Path:   findExe("go"),
		Args:   []string{"go", "build", "-o", outputFile},
		Env:    buildEnv(goos, goarch),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}

	cmd.Args = append(cmd.Args, "-ldflags")
	cmd.Args = append(cmd.Args, "-X main.Version="+version)
	cmd.Args = append(cmd.Args, ".")

	if err := cmd.Run(); err != nil {
		log.Fatal("cannot compile:", err)
	}
}

// findExe finds an executable on the PATH.
func findExe(name string) string {
	path := os.Getenv("PATH")
	if path == "" {
		log.Fatal("no PATH environment defined")
	}

	for _, dir := range filepath.SplitList(path) {
		file := filepath.Join(dir, name)
		if runtime.GOOS == "windows" {
			file += ".exe"
		}
		if fileInfo, err := os.Stat(file); err == nil && !fileInfo.IsDir() {
			return file
		}
	}

	log.Fatal("cannot find %s on PATH")
	return "" // not reached
}

func buildEnv(goos string, goarch string) []string {
	var env []string

	getenv := func(name string) string {
		v := os.Getenv(name)
		return name + "=" + v
	}

	if runtime.GOOS == "windows" {
		env = append(env, getenv("TEMP"))
	}
	env = append(env, getenv("PATH"))
	env = append(env, getenv("GOPATH"))
	env = append(env, "GOOS="+goos)
	env = append(env, "GOARCH="+goarch)
	return env
}
