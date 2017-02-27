// handlers_test.go
package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGETExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	i, err := AddExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}
	ai := strconv.Itoa(i)

	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
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
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)
	_, err := AddExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}
	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
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

	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	resp, err := http.Get(server.URL + "/v1/exercise/1200")
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDeleteExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	_, err := AddExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("DELETE", server.URL+"/v1/exercise/Test1", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if status := resp.StatusCode; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestGETExercises(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	_, err := AddExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("GET", server.URL+"/v1/exercises", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
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
	// TODO(chmou): need to be moved in constructor
	_, err := sqlTX("DELETE FROM Exercise")
	if err != nil {
		t.Fatal("Could not cleanup all exercises")
	}

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

	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise1))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusCreated {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusCreated, buf)
		return
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

	if status := resp.StatusCode; status != http.StatusCreated {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusCreated, buf)
		return
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
	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString("HALLLLO!!!"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusBadRequest {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusBadRequest, buf)
		return
	}
}

func TestPostNoting(t *testing.T) {
	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusBadRequest {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusBadRequest, buf)
		return
	}

}

func TestPostBadContent(t *testing.T) {
	exercise := `{"hello": "moto"}`
	server := httptest.NewServer(setupRoutes("./")) //Creating new server with the user handlers
	req, err := http.NewRequest("POST", server.URL+"/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusBadRequest {
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusBadRequest, buf)
		return
	}
}
