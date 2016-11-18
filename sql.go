package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Workout (
	id integer PRIMARY KEY,
	repetition float,
    meters int,
	percentage integer,
	repos text,
	CONSTRAINT uc_comicID UNIQUE (repetition, percentage, meters, repos));

CREATE TABLE IF NOT EXISTS Program (
	id integer PRIMARY KEY,
	date datetime DEFAULT CURRENT_TIMESTAMP,
    name varchar(255),
    comment text DEFAULT "",
	CONSTRAINT uc_ProgramID UNIQUE (name));

CREATE TABLE IF NOT EXISTS ProgramWorkout (
	id integer PRIMARY KEY,
	ProgramID integer,
    WorkoutID integer);`

func createSchema() (db *sql.DB, err error) {
	// TODO: proper sqlite location
	db, err = sql.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		return
	}

	_, err = db.Exec(sqlTable)
	if err != nil {
		return
	}

	err = createSample(db)
	if err != nil {
		return
	}
	return
}

func createSample(db *sql.DB) (err error) {
	var res sql.Result
	res, err = createProgram("Pyramidal", db)
	if err != nil {
		return
	}
	pyramidID, err := res.LastInsertId()
	if err != nil {
		return
	}

	res, err = createProgram("3x100", db)
	if err != nil {
		return
	}

	troiscentID, err := res.LastInsertId()
	if err != nil {
		return
	}
	res, err = createWorkout(3, 90, 800, "1.5 minutes", db)
	if err != nil {
		return
	}

	lastinsertid, err := res.LastInsertId()
	if err != nil {
		return
	}

	_, err = associateWorkoutProgram(pyramidID, lastinsertid, db)
	if err != nil {
		return
	}

	_, err = associateWorkoutProgram(troiscentID, lastinsertid, db)
	if err != nil {
		return
	}

	res, err = createWorkout(3, 95, 1000, "5 minutes", db)
	if err != nil {
		return
	}

	lastinsertid, err = res.LastInsertId()
	if err != nil {
		return
	}

	_, err = associateWorkoutProgram(pyramidID, lastinsertid, db)
	if err != nil {
		return
	}

	res, err = createWorkout(3, 100, 100, "3 minutes", db)
	if err != nil {
		return
	}

	lastinsertid, err = res.LastInsertId()
	if err != nil {
		return
	}

	_, err = associateWorkoutProgram(pyramidID, lastinsertid, db)
	if err != nil {
		return
	}

	return
}

func sqlTX(db *sql.DB, query string, args ...interface{}) (res sql.Result, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	if err != nil {
		return

	}
	defer stmt.Close()

	res, err = stmt.Exec(args...)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

// createPrograms ...
func createProgram(name string, db *sql.DB) (res sql.Result, err error) {
	var createProgramSQL = `INSERT INTO Program(name) VALUES(?);`
	res, err = sqlTX(db, createProgramSQL, name)
	return
}

func createWorkout(repetition int, percentage int, meters int, repos string, db *sql.DB) (res sql.Result, err error) {
	var createWorkoutSQL = `INSERT INTO Workout(repetition, percentage, meters, repos) VALUES(?, ?, ?, ?);`
	res, err = sqlTX(db, createWorkoutSQL, repetition, percentage, meters, repos)
	return
}

func associateWorkoutProgram(programid int64, workoutid int64, db *sql.DB) (res sql.Result, err error) {
	var associateWorkoutProgramSQL = `INSERT INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(?, ?)`
	res, err = sqlTX(db, associateWorkoutProgramSQL, programid, workoutid)
	return
}

// getWorkouts: get a workout for a
func getWorkouts(name string, db *sql.DB) (rounds []Workout, err error) {
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
