package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

// TempFileName generates a temporary filename for use in testing or whatever
// http://stackoverflow.com/a/28005931
func TempFileName(suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), hex.EncodeToString(randBytes)+suffix)
}
