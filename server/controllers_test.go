// handlers_test.go
package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETExercise(t *testing.T) {
	e := newExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	server := httptest.NewServer(router("./")) //Creating new server with the user handlers
	resp, err := http.Get(server.URL + "/v1/exercise/0")
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var exercise Exercise
	if err = json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}
}

func TestGETExerciseNotFound(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(router("./")) //Creating new server with the user handlers
	resp, err := http.Get(server.URL + "/v1/exercise/0")
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestGETExerciseBadReq(t *testing.T) {
	server := httptest.NewServer(router("./")) //Creating new server with the user handlers
	resp, err := http.Get(server.URL + "/v1/exercise/xxx")
	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestGETExercises(t *testing.T) {
	e := newExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	req, err := http.NewRequest("GET", "/v1/exercises", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GETExercises)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var exercises []Exercise
	if err = json.NewDecoder(rr.Body).Decode(&exercises); err != nil {
		t.Fatal("Could not decode body, not proper json.")
	}

	if len(exercises) != 1 {
		t.Fatal("Did not get enough exercises")
	}

}

func TestPostExcercise(t *testing.T) {
	exercise := `{
  "id": 0,
  "name": "Test1",
  "comment": "NoComment",
  "steps": [
    {
      "id": 1,
      "effort": "easy warmup todoo",
      "effort_type": "distance",
      "type": "warmup",
      "repeat": {
        "id": 0
      }
    },
    {
      "id": 1,
      "effort_type": "distance",
      "laps": 3,
      "length": 1234,
      "percentage": 90,
      "type": "interval",
      "repeat": {
        "id": 0
      }
    },
    {
      "id": 1,
      "effort": "finish strong",
      "effort_type": "distance",
      "type": "warmdown",
      "repeat": {
        "id": 0
      }
    }
  ]
}`

	req, err := http.NewRequest("POST", "/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTExercise)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	//check in DB if it was really made
	exercises, err := getAllExercises()
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 1 {
		t.Fatal("did not get all exercises")
	}
}

func TestPostEmpty(t *testing.T) {
	exercise := ""

	req, err := http.NewRequest("POST", "/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTExercise)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostNoting(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/exercise", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTExercise)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestPostBadContent(t *testing.T) {
	exercise := `{"hello": "moto"}`
	req, err := http.NewRequest("POST", "/v1/exercise", bytes.NewBufferString(exercise))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(POSTExercise)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
}
