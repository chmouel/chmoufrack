package server

import (
	"database/sql"
	"log"
	"testing"
)

func removeFromArray(slice []Step, s int) []Step {
	return append(slice[:s], slice[s+1:]...)
}

func setUp() {
	var err error

	DB, err = sql.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(sqlTable)
	if err != nil {
		log.Fatal(err)
	}
}

func newExercise(
	exerciceName, warmupEffort, warmdownEffort string,
	length2 int) (e Exercise) {
	var steps Steps

	step1 := Step{
		Type:       "warmup",
		Effort:     warmupEffort,
		EffortType: "distance",
	}
	steps = append(steps, step1)

	step2 := Step{
		Laps:       3,
		Length:     length2,
		Percentage: 90,
		Type:       "interval",
	}
	steps = append(steps, step2)

	step3 := Step{
		Effort:     warmdownEffort,
		Type:       "warmdown",
		EffortType: "distance",
	}
	steps = append(steps, step3)

	e = Exercise{
		Name:    exerciceName,
		Comment: "NoComment",
		Steps:   steps,
	}
	return
}

func TestAddExercise(t *testing.T) {
	setUp()
	e := newExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	res, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	i, err := res.LastInsertId()

	e, err = getExercise(i)
	if err != nil {
		t.Fatalf("getExercise() failed: %s", err)
	}
	if e.Name != "Test1" {
		t.Fatalf("%s != Test1", e.Name)
	}
	if e.Steps[0].Effort != "easy warmup todoo" {
		t.Fatalf("%s != easy warmup todoo", e.Steps[0].Effort)
	}
	if e.Steps[1].Length != 1234 {
		t.Fatalf("%s != 1234", int(e.Steps[1].Length))
	}
	if e.Steps[2].Effort != "finish strong" {
		t.Fatalf("%s != finish strong", e.Steps[2].Effort)
	}
}

func TestUpdateExercise(t *testing.T) {
	setUp()
	e := newExercise("TestAddUpdate", "easy warmup todoo",
		"finish strong", 4567)

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addProgram() failed: %s", err)
	}

	e.Name = "TestUpdated"
	e.Steps[0].Effort = "New1"
	e.Steps[1].Length = 999
	e.Steps[2].Effort = "New2"

	res, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when updating failed: %s", err)
	}

	i, err := res.LastInsertId()
	e, err = getExercise(i)
	if err != nil {
		t.Fatalf("getExercise() failed: %s", err)
	}
	if e.Name != "TestUpdated" {
		t.Fatalf("%s != TestUpdated", e.Name)
	}

	if e.Steps[0].Effort != "New1" {
		t.Fatalf("fatal updating warmup")
	}
	if e.Steps[1].Length != 999 {
		t.Fatalf("fatal updating interval")
	}
	if e.Steps[2].Effort != "New2" {
		t.Fatalf("fatal updating warmdown")
	}

	e.Steps = removeFromArray(e.Steps, 2)
	res, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when removing failed: %s", err)
	}
	i, err = res.LastInsertId()
	e, err = getExercise(i)
	if len(e.Steps) != 2 {
		t.Fatalf("failing to remove a step %d != 2", e.Steps.Len())
	}

	s := Step{
		Effort:     "OlaMoto",
		Type:       "warmdown",
		EffortType: "distance",
	}
	e.Steps = append(e.Steps, s)
	res, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when add a new element: %s", err)
	}
	i, err = res.LastInsertId()
	e, err = getExercise(i)
	if len(e.Steps) != 3 {
		t.Fatalf("failing to remove a step %d != 3", e.Steps.Len())
	}

}
