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
	db, err = sql.Open("sqlite3", CONFIG_DIR+"/test.db")
	if err != nil {
		return
	}

	_, err = db.Exec(sqlTable)
	return
}

func createSample(db *sql.DB) (err error) {
	var res sql.Result
	res, err = createProgram("Pyramidal", "Pyramidial Workout going string and stronger by the strongess", db)
	if err != nil {
		return
	}
	pyramidID, err := res.LastInsertId()
	if err != nil {
		return
	}

	res, err = createProgram("3x100", "3x100 is the best!", db)
	if err != nil {
		return
	}

	troiscentID, err := res.LastInsertId()
	if err != nil {
		return
	}
	res, err = createWorkout(3, 800, 90, "1.5 minutes", db)
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

	res, err = createWorkout(3, 1000, 95, "5 minutes", db)
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
