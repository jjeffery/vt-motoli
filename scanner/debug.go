package scanner

import "fmt"

type debugT bool

func (d debugT) Printf(format string, a ...interface{}) {
	if d {
		fmt.Printf("debug: "+format, a...)
	}
}

func (d debugT) Println(a ...interface{}) {
	if d {
		a = append([]interface{}{"debug:"}, a...)
		fmt.Println(a...)
	}
}

var debug = debugT(true)
