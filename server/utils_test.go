package server

import (
	"log"
	"os"
	"testing"

	"github.com/huandu/facebook"
)

type fakeFBCheck struct {
	ID   string
	Name string
	Link string
}

func (f *fakeFBCheck) FBGet(url, token string) (fbc facebook.Result, err error) {
	if f.ID == "" {
		f.ID = "1234"
	}
	if f.Name == "" {
		f.Name = "Chmou EL"
	}
	if f.Link == "" {
		f.Link = "http://facebook.com/testtest"
	}

	fbc = facebook.Result{
		"id":    f.ID,
		"name":  f.Name,
		"link":  "http://facebook.com/" + f.Link,
		"email": "chmouel@chmouel.com",
	}
	return
}

type emptyFBCheck struct{}

func (f *emptyFBCheck) FBGet(url, token string) (fbc facebook.Result, err error) {
	return
}

func TestMain(m *testing.M) {
	dblocation := os.Getenv("FRACK_TEST_DB")

	if dblocation == "" {
		log.Fatal("You need to specify a FRACK_TEST_DB variable")
	}

	err := DBConnect(dblocation, true)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
