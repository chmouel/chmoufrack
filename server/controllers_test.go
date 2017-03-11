// handlers_test.go
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"gopkg.in/gin-gonic/gin.v1"
)

type fakeFBCheck struct {
	ID   string
	Name string
	Link string
}

func (f *fakeFBCheck) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		if f.ID == "" {
			f.ID = "1234"
		}
		if f.Name == "" {
			f.Name = "Chmou EL"
		}
		if f.Link == "" {
			f.Link = "http://facebook.com/testtest"
		}
		c.Set("FBInfo", FBinfo{
			ID:   f.ID,
			Name: f.Name,
			Link: f.Link,
		})
		c.Next()
	}
}

func test_check_http_expected(resp *http.Response, expected int) (err error) {
	if status := resp.StatusCode; status != expected {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		err = errors.New(
			fmt.Sprintf("handler returned wrong status code: expected %v received %v, error: %s",
				expected, status, buf))
	}
	return
}

func TestGETExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	i, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}
	ai := strconv.Itoa(i)

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := http.Get(server.URL + "/v1/exercise/" + ai)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var exercise Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}
}

func TestGETExerciseByName(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := http.Get(server.URL + "/v1/exercise/Test1")
	if err != nil {
		t.Fatal(err)
	}
	var exercise Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}
	if exercise.Name != "Test1" {
		t.Fatalf("TestGetexercisebyname failed %s != Test1.", exercise.Name)
	}
}

func TestGETExerciseNotFound(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := http.Get(server.URL + "/v1/exercise/1200")
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDELETExerciseNotFound(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("DELETE", server.URL+"/v1/exercise/1200", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusNotFound); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDeleteExercise(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	// Delete by Name
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	req, err := http.NewRequest("DELETE", server.URL+"/v1/exercise/Test1", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusNoContent {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Delete by ID
	e = createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	lastid, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}
	url := fmt.Sprintf("%s/v1/exercise/%d", server.URL, lastid)
	req, err = http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusNoContent {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

}

func TestGETExercises(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("GET", server.URL+"/v1/exercises", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var exercises []Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercises); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}

	if len(exercises) != 1 {
		t.Fatal("Did not get enough exercises")
	}

}

func TestPostExcercise(t *testing.T) {
	exercise1 := `{"name": "Test1",
	"comment": "NoComment",
	"steps": [{
	    "effort": "easy warmup todoo",
	    "effort_type": "distance",
	    "type": "warmup"
	},{
	    "effort_type": "distance",
	    "laps": 3,
	    "length": 1234,
	    "percentage": 90,
	    "type": "interval"
	},{
	    "effort": "finish strong",
	    "effort_type": "distance",
	    "type": "warmdown"
	}]}`

	exercise_updated := `{"name": "Test1",
 "comment": "Updaated",
 "steps": []}`

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise1))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusCreated); err != nil {
		t.Fatal(err.Error())
	}

	//check in DB if it was really made
	exercises, err := getAllExercises()
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 1 {
		t.Fatal("did not have a new exercise created")
	}

	req, err = http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise_updated))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusCreated); err != nil {
		t.Fatalf(err.Error())
	}

	exercises, err = getAllExercises()
	if err != nil {
		t.Fatal(err)
	}

	if len(exercises) != 1 {
		t.Fatal("We should have only one excercise since it was an update")
	}
}

func TestPostBadJSON(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString("HALLLLO!!!"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPostNoting(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPostBadContent(t *testing.T) {
	exercise := `{"hello": "moto"}`
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

type emptyFBCheck struct{}

func (f *emptyFBCheck) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func TestWithoutFBInfo(t *testing.T) {
	fbcheck := &emptyFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString("{}"))
	if err != nil {
		t.Fatal(err.Error())
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}

	req, err = http.NewRequest("DELETE", server.URL+"/v1/exercise/blah", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRedirect(t *testing.T) {
	fbcheck := &emptyFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)
	req, err := http.NewRequest("GET", server.URL+"/", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	norclient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := norclient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != 301 {
		t.Fatal("No redirect found on /")
	}
}

func TestDeleteForSomeoneElse(t *testing.T) {
	originalFacebookId := "1234"
	otherFacebookID := "4567"
	name := "Test1"
	e := createSampleExercise(name, "easy warmup todoo", "finish strong", 1000, originalFacebookId)

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	fbcheck := &fakeFBCheck{ID: otherFacebookID}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)
	req, err := http.NewRequest("DELETE", server.URL+"/v1/exercise/"+name, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}

}
