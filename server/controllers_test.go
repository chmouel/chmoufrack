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
)

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

// test_do_request ...
func test_do_request(method, url string, data io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer bearer-token")
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	return
}

func TestRestGETExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, true, "Test User", "1234")

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

func TestRestGETExerciseByName(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, true, "Test User", "1234")
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

func TestRestGETExerciseNotFound(t *testing.T) {
	_, err := DB.Exec(SQLresetDB)
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

func TestRestDELETExerciseNotFound(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("DELETE", server.URL+"/v1/exercise", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusNotFound); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestDeleteExercise(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	// Delete by Name
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, true, "Test User", "1234")

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	resp, err := test_do_request("DELETE", server.URL+"/v1/exercise/Test1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusNoContent {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	// Delete by ID
	e = createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, true, "Test User", "1234")
	lastid, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	url := fmt.Sprintf("%s/v1/exercise/%d", server.URL, lastid)
	resp, err = test_do_request("DELETE", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusNoContent {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

}

func TestRestGETExercises(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, true, "Test User", "1234")

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

func TestRestCreateExcercice(t *testing.T) {
	_, err := DB.Exec(SQLresetDB)
	if err != nil {
		t.Fatal(err)
	}

	exercise_public := `{"name": "TestPublic",
	"comment": "NoComment",
	"public": true,
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

	exercise_private := `{
	  "name": "TestPrivate",
	  "comment": "NoComment",
      "public": false,
	  "steps": []}`

	exercise_updated := `{"name": "Test1",
 "comment": "Updaated",
 "steps": []}`

	fbcheck := &fakeFBCheck{ID: "012345678933"}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise_public))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusCreated); err != nil {
		t.Fatal(err.Error())
	}

	resp, err = test_do_request("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise_private))
	if err != nil {
		t.Fatal(err)
	}

	//check in DB if it was really made, we should only have one exercise cause we have not set the fbid
	exercises, err := getAllExercises("")
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 1 {
		t.Fatalf("we did not or had too many public exercises: %d != 1", len(exercises))
	}

	//check in DB, now we are getting the private too
	exercises, err = getAllExercises("012345678933")
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 2 {
		t.Fatalf("we did not get our private exercise: %d != 2", len(exercises))
	}

	resp, err = test_do_request("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise_updated))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusCreated); err != nil {
		t.Fatalf(err.Error())
	}

	exercises, err = getAllExercises("")
	if err != nil {
		t.Fatal(err)
	}

	if len(exercises) != 1 {
		t.Fatal("We should have only one excercice since it was an update")
	}
}

func TestRestGetPrivateExercise(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	fbid := "123456789"
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, false, "Test User", fbid)
	_, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	fbcheck := &fakeFBCheck{ID: fbid}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	// Get exercise
	resp, err := test_do_request("GET", server.URL+"/v1/exercise/Test1", nil)
	if err != nil {
		t.Fatal(err)
	}
	var exercise Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}
	if exercise.Name != "Test1" {
		t.Fatalf("TestRestGetPrivateExercise failed %s != Test1.", exercise.Name)
	}

	// Get all excercices
	resp, err = test_do_request("GET", server.URL+"/v1/exercises", nil)
	if err != nil {
		t.Fatal(err)
	}
	var exercises []Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercises); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}

	if exercises[0].Name != "Test1" {
		t.Fatalf("TestRestGetPrivateExercise failed %s != Test1.", exercise.Name)
	}

	// No authorization should be public
	resp, err = http.Get(server.URL + "/v1/Test1")
	if err != nil {
		t.Fatal(err)
	}

	// Since it's a public request and ours is private, then fail
	if status := resp.StatusCode; status != http.StatusNotFound {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestRestPostNoFB(t *testing.T) {
	fbcheck := &emptyFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestPostBadJSON(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("POST", server.URL+"/v1/exercise", bytes.NewBufferString("HALLO"))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestPostNoting(t *testing.T) {
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("POST", server.URL+"/v1/exercise", nil)
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestPostBadContent(t *testing.T) {
	exercise := `{"hello": "moto"}`
	fbcheck := &fakeFBCheck{}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestRestWithoutFBInfo(t *testing.T) {
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

func TestRestRedirect(t *testing.T) {
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

func TestRestDeleteForSomeoneElse(t *testing.T) {
	_, err := DB.Exec(SQLresetDB)
	if err != nil {
		t.Fatal(err)
	}

	originalFacebookId := "1234"
	otherFacebookID := "4567"
	name := "Test1"
	e := createSampleExercise(name, "easy warmup todoo", "finish strong", 1000, true, "Test User", originalFacebookId)

	_, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	fbcheck := &fakeFBCheck{ID: otherFacebookID}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)
	resp, err := test_do_request("DELETE", server.URL+"/v1/exercise/"+name, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}

}

func TestRestCreateFBInfo(t *testing.T) {
	fbid := "4557"
	fbinfo_rest := `{
	"id": "%s",
    "name": "Ola Chica",
    "link": "https://www.facebook.com/app_scoped_user_id/10157827176205251/"
  }`

	fbcheck := &fakeFBCheck{ID: fbid}
	server := httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err := test_do_request("POST", server.URL+"/v1/fbinfo", bytes.NewBufferString(
		fmt.Sprintf(fbinfo_rest, fbid)))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusCreated); err != nil {
		t.Fatal(err.Error())
	}

	fbcheck = &fakeFBCheck{ID: "FAKENEWS!"}
	server = httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	resp, err = test_do_request("POST", server.URL+"/v1/fbinfo", bytes.NewBufferString(fbinfo_rest))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatal(err.Error())
	}

	fbcheck = &fakeFBCheck{ID: fbid}
	server = httptest.NewServer(
		setupRoutes("./", fbcheck),
	)

	// Bad content no string
	resp, err = test_do_request("POST", server.URL+"/v1/fbinfo", bytes.NewBufferString("FAKENEWS"))
	if err != nil {
		t.Fatal(err)
	}

	if err = test_check_http_expected(resp, http.StatusBadRequest); err != nil {
		t.Fatal(err.Error())
	}

	emptyFBCheck := &emptyFBCheck{}
	server = httptest.NewServer(
		setupRoutes("./", emptyFBCheck),
	)

	req, err := http.NewRequest("POST", server.URL+"/v1/fbinfo", bytes.NewBufferString(fbinfo_rest))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err = test_check_http_expected(resp, http.StatusUnauthorized); err != nil {
		t.Fatalf(err.Error())
	}
}
