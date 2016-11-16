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
    comment text);

CREATE TABLE IF NOT EXISTS ProgramWorkout (
	id integer PRIMARY KEY,
	ProgramID integer,
    WorkoutID integer);`

var sqlSamples = `
INSERT INTO Program(name) VALUES("Pyramidal");
INSERT INTO Program(name) VALUES("3x100");

INSERT INTO Workout(repetition, percentage, meters, repos) VALUES(3, 90, 1500, "1.5 minutes");
INSERT INTO Workout(repetition, percentage, meters, repos) VALUES(3, 95, 1000, "1.5 minutes");

INSERT INTO Workout(repetition, percentage, meters, repos) VALUES(3, 100, 100, "3 minutes");

INSERT INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(1, 1);
INSERT INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(1, 2);

INSERT INTO ProgramWorkout(ProgramID, WorkoutID) VALUES(2, 3);

SELECT W.id FROM Program P, Workout W, ProgramWorkout PW WHERE P.name = 'Pyramidal' AND PW.WorkoutID == W.ID AND PW.ProgramID == P.id;
SELECT W.id FROM Program P, Workout W, ProgramWorkout PW WHERE P.name = '3x100' AND PW.WorkoutID == W.ID AND PW.ProgramID == P.id;
`

func createTable() (db *sql.DB, err error) {
	// TODO: proper sqlite location
	db, err = sql.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		return
	}

	_, err = db.Exec(sqlTable)
	if err != nil {
		return
	}

	// TODO: remove samples
	_, err = db.Exec(sqlSamples)
	if err != nil {
		return
	}

	return
}
