// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	log.SetFlags(0)

	fs := http.Dir("assets")

	if err := vfsgen.Generate(fs, vfsgen.Options{}); err != nil {
		log.Fatal(err)
	}
}
