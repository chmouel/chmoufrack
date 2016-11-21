package main

import (
	"database/sql"
	"fmt"
	"log"
)

// createPrograms ...
func createProgram(name string, comment string) (res sql.Result, err error) {
	var createProgramSQL = `INSERT OR REPLACE INTO Program(name, comment) VALUES(?, ?);`
	res, err = sqlTX(createProgramSQL, name, comment)
	return
}

func createWorkout(repetition int, meters int, percentage int, repos string) (res sql.Result, err error) {
	var createWorkoutSQL = `INSERT OR REPLACE INTO Workout(repetition, percentage, meters, repos) VALUES(?, ?, ?, ?);`
	res, err = sqlTX(createWorkoutSQL, repetition, percentage, meters, repos)
	return
}

func deleteProgram(name string) (res sql.Result, err error) {
	p, err := getProgram(name)
	if err != nil {
		log.Fatalf("No results for Program %s", name)
	}
	var deletePWSQL = `DELETE FROM ProgramWorkout WHERE ProgramID == ?`
	res, err = sqlTX(deletePWSQL, p.ID)

	var deleteProgramSQL = `DELETE FROM Program WHERE name=?`
	res, err = sqlTX(deleteProgramSQL, name)
	return
}

func deleteWorkout(id int64) (res sql.Result, err error) {
	var deleteWorkoutSQL = `DELETE FROM ProgramWorkout WHERE WorkoutID == ?`
	res, err = sqlTX(deleteWorkoutSQL, id)

	deleteWorkoutSQL = `DELETE FROM Workout WHERE id=?`
	res, err = sqlTX(deleteWorkoutSQL, id)
	return
}

func associateWorkoutProgram(programid int64, workoutid int64) (res sql.Result, err error) {
	var associateWorkoutProgramSQL = `INSERT OR REPLACE INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(?, ?)`
	res, err = sqlTX(associateWorkoutProgramSQL, programid, workoutid)
	return
}

func associateWorkoutProgramByName(workoutName string, programName string) (res sql.Result, err error) {
	w, err := getWorkoutByName(workoutName)
	if err != nil {
		return
	}
	if w.Meters == "" {
		//TODO: Error handling
		log.Fatalf("No results for workout %s", workoutName)
	}

	p, err := getProgram(programName)
	if err != nil {
		log.Fatalf("No results for program %s", programName)
	}
	_, err = associateWorkoutProgram(w.ID, p.ID)
	return
}

// getWorkouts: get a workout for a
func getWorkoutsforProgram(name string) (rounds []Workout, err error) {
	var getWorkoutSQL = `
    SELECT W.id, W.repetition, W.meters, W.percentage, W.repos
       FROM Program P, Workout W, ProgramWorkout PW
       WHERE P.name = $1 AND PW.WorkoutID == W.ID
       AND PW.ProgramID == P.id`

	rows, err := DB.Query(getWorkoutSQL, name)
	for rows.Next() {
		var w Workout
		err = rows.Scan(&w.ID, &w.Repetition, &w.Meters, &w.Percentage, &w.Repos)
		if err != nil {
			return
		}
		rounds = append(rounds, w)
	}
	return
}

// getPrograms ...
func getPrograms() (programs []Program, err error) {
	var getProgramsSQL = `SELECT id, name, date, comment from Program`
	rows, err := DB.Query(getProgramsSQL)
	if err != nil {
		return
	}
	for rows.Next() {
		var p Program
		err = rows.Scan(&p.ID, &p.Name, &p.Date, &p.Comment)
		if err != nil {
			return
		}
		programs = append(programs, p)
	}
	return
}

func getProgram(programName string) (program Program, err error) {
	var getProgramsSQL = `SELECT id, name, date, comment from Program where name = $1`
	err = DB.QueryRow(getProgramsSQL, programName).Scan(
		&program.ID, &program.Name,
		&program.Date, &program.Comment)
	if err != nil {
		return
	}
	return
}

// getWorkout ...
func getWorkouts() (workouts []Workout, err error) {
	var getProgramsSQL = `SELECT repetition, meters, percentage, repos from Workout`
	rows, err := DB.Query(getProgramsSQL)
	if err != nil {
		return
	}
	for rows.Next() {
		var w Workout
		err = rows.Scan(&w.Repetition, &w.Meters, &w.Percentage, &w.Repos)
		if err != nil {
			return
		}
		workouts = append(workouts, w)
	}
	return
}

func getWorkoutByName(workoutName string) (workout Workout, err error) {
	var getProgramsSQL = `SELECT id, repetition, meters, percentage, repos from Workout`
	rows, err := DB.Query(getProgramsSQL)
	if err != nil {
		return
	}
	for rows.Next() {
		var w Workout
		err = rows.Scan(&w.ID, &w.Repetition, &w.Meters, &w.Percentage, &w.Repos)
		if err != nil {
			return
		}
		name := fmt.Sprintf("%sx%s@%s", w.Repetition, w.Meters, w.Percentage)
		if name == workoutName {
			workout = w
		}
	}
	return
}
