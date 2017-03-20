package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFBCheck(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("GET", server.URL+"/v1/exercises", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Authorization", "XX")
	resp, err := http.DefaultClient.Do(req)
	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}

	emptyFB := &emptyFBCheck{}
	server = httptest.NewServer(
		setupRoutes("./", emptyFB),
	)

	req, err = http.NewRequest("GET", server.URL+"/v1/exercises", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Authorization", "Bearer Hallo")
	resp, err = http.DefaultClient.Do(req)

	fmt.Println(resp.StatusCode)

}
