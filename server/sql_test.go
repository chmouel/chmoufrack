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
		EffortType: "distance",
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

func TestAddGetRepeat(t *testing.T) {
	setUp()
	e := newExercise("Test1", "easy warmup todoo", "finish strong", 1234)
	originSteps := len(e.Steps)

	var repeatSteps Steps
	repeatStep := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep)

	repeat := Repeat{
		Steps:  repeatSteps,
		Repeat: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when adding a repeat: %s", err)
	}

	e, err = getExercise(0)
	if len(e.Steps) != originSteps+1 {
		t.Fatalf("failing to add a repeat %d != %d", e.Steps.Len(),
			originSteps+1)
	}

	// This should test positioning too
	if len(e.Steps[3].Repeat.Steps) != 1 {
		t.Fatalf("failing to getAdd a repeat step %d != 1",
			len(e.Steps[3].Repeat.Steps))
	}

	if e.Steps[3].Repeat.Steps[0].Laps != 6 {
		t.Fatalf("failing to get step field %s != 6",
			e.Steps[3].Repeat.Steps[0].Laps)
	}

	e.Steps[3].Repeat.Steps[0].Laps = 99
	_, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() updating a repeat step: %s", err)
	}

	e, err = getExercise(0)
	if len(e.Steps[3].Repeat.Steps) != 1 {
		t.Fatalf("failing updating a repeat step %d != 1",
			len(e.Steps[3].Repeat.Steps))
	}
	e, err = getExercise(0)
	if e.Steps[3].Repeat.Steps[0].Laps != 99 {
		t.Fatalf("failing updating a repeat laps %d != 99",
			e.Steps[3].Repeat.Steps[0].Laps)
	}

	e.Steps[3].Repeat.Steps = removeFromArray(e.Steps[3].Repeat.Steps, 0)
	_, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when removing from field failed: %s", err)
	}
	e, err = getExercise(0)
	if len(e.Steps[3].Repeat.Steps) != 0 {
		t.Fatalf("addExercise() failing to remove field: %s", err)
	}
}

func TestAddGetRepeatDoublonMixedUP(t *testing.T) {
	setUp()
	e := newExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	var repeatSteps Steps
	repeatStep1 := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep1)

	repeat := Repeat{
		Steps:  repeatSteps,
		Repeat: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	_, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding first repeat: %s", err)
	}
	e, err = getExercise(0)

	var repeatSteps2 Steps
	repeatStep2 := Step{
		Laps:       10,
		Length:     1000,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps2 = append(repeatSteps2, repeatStep2)
	repeat = Repeat{
		Steps:  repeatSteps2,
		Repeat: 5,
	}
	exerciseStep = Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	_, err = addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding second repeat: %s", err)
	}

	e, err = getExercise(0)
	if len(e.Steps) != 5 {
		t.Fatalf("addExercise() when adding second repeat %s!=5", len(e.Steps))
	}

	e.Steps[3].Repeat.Steps[0].Length = 999
	_, err = addExercise(e)
	if err != nil {
		t.Fatal(err)
	}
	e, err = getExercise(0)
	if err != nil {
		t.Fatal(err)
	}
	if e.Steps[3].Repeat.Steps[0].Length != 999 {
		t.Fatalf("failing to update repeat")
	}
	if len(e.Steps[4].Repeat.Steps) != 1 {
		t.Fatalf("failing to update repeat")
	}

	if e.Steps[4].Repeat.Steps[0].Length == 999 {
		t.Fatalf("failing to update repeat")
	}
}
