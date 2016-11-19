package main

import "database/sql"

// createPrograms ...
func createProgram(name string, comment string, db *sql.DB) (res sql.Result, err error) {
	var createProgramSQL = `INSERT OR REPLACE INTO Program(name, comment) VALUES(?, ?);`
	res, err = sqlTX(db, createProgramSQL, name, comment)
	return
}

func createWorkout(repetition int, percentage int, meters int, repos string, db *sql.DB) (res sql.Result, err error) {
	var createWorkoutSQL = `INSERT OR REPLACE INTO Workout(repetition, percentage, meters, repos) VALUES(?, ?, ?, ?);`
	res, err = sqlTX(db, createWorkoutSQL, repetition, percentage, meters, repos)
	return
}

func associateWorkoutProgram(programid int64, workoutid int64, db *sql.DB) (res sql.Result, err error) {
	var associateWorkoutProgramSQL = `INSERT OR REPLACE INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(?, ?)`
	res, err = sqlTX(db, associateWorkoutProgramSQL, programid, workoutid)
	return
}

// getWorkouts: get a workout for a
func getWorkoutsforProgram(name string, db *sql.DB) (rounds []Workout, err error) {
	var getWorkoutSQL = `
    SELECT W.repetition, W.meters, W.percentage, W.repos
       FROM Program P, Workout W, ProgramWorkout PW
       WHERE P.name = $1 AND PW.WorkoutID == W.ID
       AND PW.ProgramID == P.id`

	rows, err := db.Query(getWorkoutSQL, name)
	for rows.Next() {
		var w Workout
		err = rows.Scan(&w.Repetition, &w.Meters, &w.Percentage, &w.Repos)
		if err != nil {
			return
		}
		rounds = append(rounds, w)
	}
	return
}

// getPrograms ...
func getPrograms(db *sql.DB) (programs []Program, err error) {
	var getProgramsSQL = `SELECT name, date, comment from Program`
	rows, err := db.Query(getProgramsSQL)
	if err != nil {
		return
	}
	for rows.Next() {
		var p Program
		err = rows.Scan(&p.Name, &p.Date, &p.Comment)
		if err != nil {
			return
		}
		programs = append(programs, p)
	}
	return
}

// getWorkout ...
func getWorkouts(db *sql.DB) (workouts []Workout, err error) {
	var getProgramsSQL = `SELECT repetition, meters, percentage, repos from Workout`
	rows, err := db.Query(getProgramsSQL)
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
