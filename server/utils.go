package server

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func DebugVargs(args ...interface{}) {
	fmt.Printf("(")
	cnt := 1
	for _, v := range args {
		if v == nil {
			fmt.Printf("'', ")
			continue
		}

		switch value := v.(type) {
		case string:
			fmt.Printf(`"` + value + `", `)
		case int:
			fmt.Printf("%d, ", value)
		}
		if cnt != len(args) {
			fmt.Printf(", ")
		}
		cnt += 1
	}
	fmt.Println(")")
}

func Debug(a ...interface{}) {
	f, err := os.Create("/tmp/debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	spew.Fdump(f, a...)
}
