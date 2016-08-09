package touch
import (
	"path/filepath"
	"os"
	"fmt"
	"path"
	"time"
)

func touch(filename string, f os.FileInfo, err error) error {
	if(path.Ext(filename)==".html"){
		fmt.Printf("Visited: %s\n", filename)
		os.Chtimes(filename, time.Now(), time.Now())
	}
	return nil
}

func RecursiveTouchHtml(root string){
	err := filepath.Walk(root, touch)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
