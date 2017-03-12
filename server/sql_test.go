package server

import (
	"fmt"
	"log"
	"testing"
)

func removeFromArray(slice []Step, s int) []Step {
	return append(slice[:s], slice[s+1:]...)
}

func TestAddExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	i, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

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
	if e.Steps[1].Length != 1000 {
		t.Fatalf("%s != 1000", int(e.Steps[1].Length))
	}
	if e.Steps[2].Effort != "finish strong" {
		t.Fatalf("%s != finish strong", e.Steps[2].Effort)
	}
}

func TestAddExerciseAndID(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	oldid, err := addExercise(e)
	if err != nil {
		log.Fatal()
	}

	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	e = createSampleExercise("Test2", "easy warmup todoo", "finish strong", 1000, "1234")
	newid, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	e, err = getExercise(newid)
	if err != nil {
		t.Fatalf("getExercise() failed: %s", err)
	}

	if newid == oldid || oldid > newid {
		t.Fatalf("the new exercices id should not have been the old one: NEW:%d, OLD:%d", newid, oldid)
	}
}

func TestUpdateExercise(t *testing.T) {
	e := createSampleExercise("TestAddUpdate", "easy warmup todoo", "finish strong", 4567, "1234")

	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addProgram() failed: %s", err)
	}

	e.Name = "TestUpdated"
	e.Steps[0].Effort = "New1"
	e.Steps[1].Length = 999
	e.Steps[2].Effort = "New2"

	i, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when updating failed: %s", err)
	}

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
	i, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when removing failed: %s", err)
	}
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
	i, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when add a new element: %s", err)
	}

	e, err = getExercise(i)
	if len(e.Steps) != 3 {
		t.Fatalf("failing to remove a step %d != 3", e.Steps.Len())
	}

}

func TestAddGetRepeat(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
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

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)

	i, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when adding a repeat: %s", err)
	}

	e, err = getExercise(i)
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
	i, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() updating a repeat step: %s", err)
	}

	e, err = getExercise(i)
	if len(e.Steps[3].Repeat.Steps) != 1 {
		t.Fatalf("failing updating a repeat step %d != 1",
			len(e.Steps[3].Repeat.Steps))
	}
	e, err = getExercise(i)
	if e.Steps[3].Repeat.Steps[0].Laps != 99 {
		t.Fatalf("failing updating a repeat laps %d != 99",
			e.Steps[3].Repeat.Steps[0].Laps)
	}

	e.Steps[3].Repeat.Steps = removeFromArray(e.Steps[3].Repeat.Steps, 0)
	i, err = addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() when removing from field failed: %s", err)
	}
	e, err = getExercise(i)
	if len(e.Steps[3].Repeat.Steps) != 0 {
		t.Fatalf("addExercise() failing to remove field: %s", err)
	}
}

func TestAddGetRepeatDoublonMixedUP(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	var repeatSteps Steps
	repeatStep1 := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep1)

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	i, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding first repeat: %s", err)
	}
	e, err = getExercise(i)

	var repeatSteps2 Steps
	repeatStep2 := Step{
		Laps:       10,
		Length:     1000,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps2 = append(repeatSteps2, repeatStep2)
	repeat = Repeats{
		Steps:   repeatSteps2,
		Repeats: 5,
	}
	exerciseStep = Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	i, err = addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding second repeat: %s", err)
	}

	e, err = getExercise(i)
	if len(e.Steps) != 5 {
		t.Fatalf("addExercise() when adding second repeat %s!=5", len(e.Steps))
	}

	e.Steps[3].Repeat.Steps[0].Length = 999
	i, err = addExercise(e)
	if err != nil {
		t.Fatal(err)
	}
	e, err = getExercise(i)
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

func TestGetByName(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}
	i, err := getIdOfExerciseName("Test1")
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	e, err = getExercise(i)
	if err != nil {
		t.Fatalf("getExercise() failed: %s", err)
	}
	if e.Name != "Test1" {
		t.Fatal("Failed to TestGetByName")
	}
}

func TestDBDeleteExercise(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	i, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	err = deleteExercise(e)
	if err != nil {
		t.Fatal(err)
	}

	e, err = getExercise(i)
	if _, ok := err.(*error404); !ok {
		fmt.Println(err)
	}
}

func TestGetAllExercices(t *testing.T) {
	_, _ = DB.Exec("DELETE from Exercise")

	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	_, err := addExercise(e)
	if err != nil {
		t.Fatalf("addExercise() failed: %s", err)
	}

	exercises, err := getAllExercises()
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 1 {
		t.Fatal("did not get all exercises")
	}

}

func TestGetAllExercicesNotFound(t *testing.T) {
	_, err := DB.Exec("DELETE FROM Exercise")
	if err != nil {
		t.Fatal(err)
	}

	exercises, err := getAllExercises()
	if err != nil {
		t.Fatal(err)
	}
	if len(exercises) != 0 {
		t.Fatal("did not get all exercises")
	}
}

func TestAddGetRepeatDoublon(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")

	var repeatSteps Steps
	repeatStep1 := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep1)

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	i, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding first repeat: %s", err)
	}
	e, err = getExercise(i)

	var repeatSteps2 Steps
	repeatStep2 := Step{
		Laps:       10,
		Length:     1000,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps2 = append(repeatSteps2, repeatStep2)
	repeat = Repeats{
		Steps:   repeatSteps2,
		Repeats: 5,
	}
	exerciseStep = Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	i, err = addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding second repeat: %s", err)
	}

	e, err = getExercise(i)
	if len(e.Steps) != 5 {
		t.Fatalf("addExercise() when adding second repeat %s!=5", len(e.Steps))
	}

	e.Steps = removeFromArray(e.Steps, 4)
	i, err = addExercise(e)
	e, err = getExercise(i)
	if len(e.Steps) != 4 {
		t.Fatalf("removing repeat is not working steps %d != 4", len(e.Steps))
	}
}

func TestNotHere(t *testing.T) {
	var err error
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	_, err = addExercise(e)
	if err != nil {
		t.FailNow()
	}
	_, err = getExercise(50)

	if _, ok := err.(*error404); !ok {
		t.FailNow()
	}

	_, err = getIdOfExerciseName("blahblah")
	if _, ok := err.(*error404); !ok {
		t.FailNow()
	}
}

func TestUPAndDown(t *testing.T) {

	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	var repeatSteps Steps
	repeatStep1 := Step{
		Laps:       5,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep1)

	repeatStep2 := Step{
		Laps:       10,
		Length:     1000,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep2)

	repeatStep3 := Step{
		Laps:       15,
		Length:     1000,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep3)

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)
	i, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding first repeat: %s", err)
	}
	e, err = getExercise(i)

	first := e.Steps[3].Repeat.Steps[0]
	second := e.Steps[3].Repeat.Steps[1]
	third := e.Steps[3].Repeat.Steps[2]
	var inversedSteps Steps
	inversedSteps = append(inversedSteps, third)
	inversedSteps = append(inversedSteps, first)
	inversedSteps = append(inversedSteps, second)

	e.Steps[3].Repeat.Steps = inversedSteps
	i, err = addExercise(e)
	e, err = getExercise(i)

	newfirst := e.Steps[3].Repeat.Steps[0]
	newsecond := e.Steps[3].Repeat.Steps[1]
	newthird := e.Steps[3].Repeat.Steps[2]

	if third.Laps != newfirst.Laps {
		t.Fatal("Old third should be in position first")
	}

	if first.Laps != newsecond.Laps {
		t.Fatal("Old first should be in second position")
	}
	if second.Laps != newthird.Laps {
		t.Fatal("Old firssecond should be in third position")

	}

}

func TestUpdateForSomeoneElse(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	_, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding first repeat: %s", err)
	}

	e.FB.ID = "5678"
	_, err = addExercise(e)
	if _, ok := err.(*errorUnauthorized); !ok {
		t.Fatal(err.Error())
	}
}

func TestNotExerciseName(t *testing.T) {
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, "1234")
	e.Name = ""
	_, err := addExercise(e)

	if err == nil {
		t.Fatalf("We should have an error adding new exercise without an exercise name")
	}

}

func TestBadCharacters(t *testing.T) {
	s := "foo/bar"
	err := checkBadCharacters(s)
	if err == nil {
		t.Fatalf("We should have detected that we have a bad characters")
	}
}

func TestAddEmojis(t *testing.T) {
	e := createSampleExercise("Test1 💜", "easy warmup todoo 💜", "finish strong 💜", 1000, "1234")
	_, err := addExercise(e)

	if err != nil {
		t.Fatalf("addExercise() when adding second repeat: %s", err)
	}
}
