package server

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
DROP TABLE If Exists Excercise;
DROP TABLE If Exists Warmup;
DROP TABLE If Exists Warmdown;
DROP TABLE If Exists Interval;
DROP TABLE If Exists Repeat;

CREATE TABLE IF NOT EXISTS Excercise (
	id integer PRIMARY KEY,
	name varchar(255),
    comment text DEFAULT "",
	CONSTRAINT uc_ExcerciseID UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS Warmup (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	excerciseID integer
);

CREATE TABLE IF NOT EXISTS Warmdown (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	excerciseID integer
);

CREATE TABLE IF NOT EXISTS Interval (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	laps tinyint NOT NULL,
    length INTEGER NOT NULL,
	percentage tinyint NOT NULL,
	rest text,
	effort_type varchar(32) DEFAULT "distance",
	repeatID integer,
	excerciseID integer
);

Create Table IF NOT EXISTS Repeat  (
	id integer PRIMARY KEY,
	repeat tinyint,
	position tinyint DEFAULT 0,
	excerciseID integer
);
`

//TODO: remove
var aSample = `
	DELETE FROM Excercise;
	DELETE FROM Warmup;
	DELETE FROM Warmdown;
	DELETE FROM Interval;
	DELETE FROM Repeat;

	INSERT INTO Excercise(name) VALUES("Pyramids Short");

	INSERT INTO Warmup(effort_type, effort, position, excerciseID) VALUES("distance", "5km very easy around", 1, 1);
	INSERT INTO Repeat(Repeat, position, excerciseID) VALUES(5, 2, 1);
	INSERT INTO Warmdown(effort_type, effort, position, excerciseID) VALUES("time", "15 mn footing", 3, 1);

	INSERT INTO Interval(laps, length, percentage, rest, effort_type, repeatID) VALUES(6, 1000, 90, "400m active", "distance", 1);
`

func DBConnect() (err error) {
	// TODO: proper sqlite location
	DB, err = sql.Open("sqlite3", "/tmp/a.db")
	if err != nil {
		return
	}

	_, err = DB.Exec(sqlTable)
	if err != nil {
		return
	}

	_, err = DB.Exec(aSample)
	return
}

func sqlTX(query string, args ...interface{}) (res sql.Result, err error) {
	tx, err := DB.Begin()
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

func getSteps(t string, id int, steps *[]Step) (err error) {
	var getWarmupSQL = fmt.Sprintf(`SELECT position, effort, effort_type FROM Warmup WHERE %s=?`, t)
	var getWarmdownSQL = fmt.Sprintf(`SELECT position, effort, effort_type FROM Warmdown WHERE %s=?`, t)
	var getIntervalSQL = fmt.Sprintf(`SELECT position, laps, length, percentage, rest, effort_type FROM Interval WHERE %s=?`, t)

	rows, err := DB.Query(getWarmupSQL, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
	}
	for rows.Next() {
		var step = Step{
			Type: "warmup",
		}
		err = rows.Scan(&step.Position, &step.Effort, &step.EffortType)
		*steps = append(*steps, step)
	}

	rows, err = DB.Query(getWarmdownSQL, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
	}
	for rows.Next() {
		var step = Step{
			Type: "warmdown",
		}
		err = rows.Scan(&step.Position, &step.Effort, &step.EffortType)
		*steps = append(*steps, step)
	}

	rows, err = DB.Query(getIntervalSQL, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
	}
	for rows.Next() {
		var step = Step{
			Type: "interval",
		}
		err = rows.Scan(&step.Position, &step.Laps, &step.Length, &step.Percentage, &step.Rest, &step.EffortType)
		*steps = append(*steps, step)
	}
	return
}

func getProgram(excerciseName string) (excercise Excercise, err error) {
	var getExcerciseSQL = `SELECT id, comment FROM Excercise WHERE name=?`
	var getRepeatSQL = `SELECT id, position, repeat FROM Repeat WHERE excerciseID=?`
	var steps []Step
	var repeatStep []Step

	excercise = Excercise{
		Name: excerciseName,
	}

	err = DB.QueryRow(getExcerciseSQL, excerciseName).Scan(
		&excercise.ID,
		&excercise.Comment)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = &error404{"ProgramNotFound"}
		return
	}

	err = getSteps("excerciseID", excercise.ID, &steps)
	if err != nil {
		fmt.Println("e")
		return
	}

	step := Step{
		Type: "repeat",
	}
	repeat := Repeat{}
	err = DB.QueryRow(getRepeatSQL, excercise.ID).Scan(
		&repeat.ID,
		&step.Position,
		&repeat.Repeat)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Printf("repeat error2: %s\n", err.Error())
			return
		} else {
			err = nil
		}
	} else {
		err = getSteps("repeatID", excercise.ID, &repeatStep)
		if err != nil {
			if err != sql.ErrNoRows {
				return
			} else {
				err = nil
			}
		}
	}
	repeat.Steps = repeatStep
	step.Repeat = repeat
	if len(repeat.Steps) != 0 {
		steps = append(steps, step)
	}
	excercise.Steps = steps

	sort.Sort(excercise)

	return
}

func getAllPrograms() (excercises []Excercise, err error) {
	var getAllExcercises = `SELECT name from Excercise`
	rows, err := DB.Query(getAllExcercises)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}

	for rows.Next() {
		e := Excercise{}
		err = rows.Scan(&e.Name)
		if err != nil {
			return
		}
		e, err = getProgram(e.Name)
		if err != nil {
			return
		}
		excercises = append(excercises, e)
	}
	return
}

func addProgram(excercise Excercise) (err error) {
	fmt.Println(excercise.ID)
	return
}
