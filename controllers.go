package chmoufrack

import (
	"database/sql"
	"fmt"
	"log"
)

// TODO: Inserting comments is not working, figure this out!
// createPrograms ...
func CreateProgram(name, comment string) (res sql.Result, err error) {
	var createProgramSQL = `INSERT INTO Program(name, comment) VALUES(?, ?)`
	fmt.Println(createProgramSQL)
	res, err = sqlTX(createProgramSQL, name, comment)
	return
}

func CreateWorkout(repetition, meters, percentage int, repos string, programID int) (res sql.Result, err error) {
	var createWorkoutSQL = `INSERT INTO Workout(repetition, percentage, meters, repos, programID) VALUES(?, ?, ?, ?, ?);`
	res, err = sqlTX(createWorkoutSQL, repetition, percentage, meters, repos, programID)
	return
}

func DeleteProgram(name string) (res sql.Result, err error) {
	p, err := GetProgram(name)
	if err != nil {
		log.Fatalf("No results for Program %s", name)
	}
	var deletePWSQL = `DELETE FROM Workout WHERE programID == ?`
	res, err = sqlTX(deletePWSQL, p.ID)

	var deleteProgramSQL = `DELETE FROM Program WHERE name=?`
	res, err = sqlTX(deleteProgramSQL, name)
	return
}

func DeleteWorkout(programID, id int64) (res sql.Result, err error) {
	deleteWorkoutSQL := `DELETE FROM Workout WHERE programID=? and id = ?`
	res, err = sqlTX(deleteWorkoutSQL, programID, id)
	return
}

func DeleteAllWorkoutProgram(programID int64) (res sql.Result, err error) {
	deleteWorkoutSQL := `DELETE FROM Workout WHERE programID=?`
	res, err = sqlTX(deleteWorkoutSQL, programID)
	return
}

// getWorkouts: get a workout for a
func GetWorkoutsforProgram(name string) (rounds []Workout, err error) {
	var getWorkoutSQL = `
    SELECT W.id, W.repetition, W.meters, W.percentage, W.repos
       FROM Program P, Workout W WHERE P.name = $1
       AND W.ProgramID == P.id`

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
func GetPrograms() (programs []Program, err error) {
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

func GetProgram(programName string) (program Program, err error) {
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
func GetWorkouts() (workouts []Workout, err error) {
	var getProgramsSQL = `SELECT W.repetition, W.meters, W.percentage, W.repos
					        FROM Program P, Workout W WHERE W.ProgramID == P.ID`
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

func GetWorkoutByName(workoutName string) (workout Workout, err error) {
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
