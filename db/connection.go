package db

import (
	"database/sql"

	"github.com/chmouel/chmoufrack"
	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS Workout (
	id integer PRIMARY KEY,
	repetition float,
    meters int,
	percentage integer,
	repos text,
	programID integer NOT NULL,
	CONSTRAINT uc_comicID UNIQUE (repetition, percentage, meters, repos, programID));

CREATE TABLE IF NOT EXISTS Program (
	id integer PRIMARY KEY,
	date datetime DEFAULT CURRENT_TIMESTAMP,
    name varchar(255),
    comment text DEFAULT "",
	CONSTRAINT uc_ProgramID UNIQUE (name));`

func CreateSchema() (err error) {
	// TODO: proper sqlite location
	chmoufrack.DB, err = sql.Open("sqlite3",
		chmoufrack.ConfigDir+"/chmoufrack.db")
	if err != nil {
		return
	}

	_, err = chmoufrack.DB.Exec(sqlTable)
	return
}

func CreateSample() (err error) {
	var res sql.Result
	res, err = CreateProgram("Pyramidal", "Pyramidial Workout going string and stronger by the strongess")
	if err != nil {
		return
	}

	pyramidID, err := res.LastInsertId()
	if err != nil {
		return
	}

	res, err = CreateWorkout(5, 400, 100, "1.5 minutes", int(pyramidID))
	if err != nil {
		return
	}

	res, err = CreateWorkout(3, 800, 95, "1.5 minutes", int(pyramidID))
	if err != nil {
		return
	}

	res, err = CreateWorkout(2, 1000, 90, "5 minutes", int(pyramidID))
	if err != nil {
		return
	}

	res, err = CreateProgram("8x400", "8x400 is the best!")
	if err != nil {
		return
	}

	troiscentID, err := res.LastInsertId()
	if err != nil {
		return
	}

	res, err = CreateWorkout(8, 400, 95, "Time it takes to complete", int(troiscentID))
	if err != nil {
		return
	}

	return
}

func sqlTX(query string, args ...interface{}) (res sql.Result, err error) {
	tx, err := chmoufrack.DB.Begin()
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
