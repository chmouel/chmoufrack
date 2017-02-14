package server

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

var sqlTable = `
DROP TABLE If Exists Exercise;
DROP TABLE If Exists Warmup;
DROP TABLE If Exists Warmdown;
DROP TABLE If Exists Interval;
DROP TABLE If Exists Repeat;

CREATE TABLE IF NOT EXISTS Exercise (
	id integer PRIMARY KEY,
	name varchar(255),
    comment text DEFAULT "",
	CONSTRAINT uc_ExerciseID UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS Warmup (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	exerciseID integer
);

CREATE TABLE IF NOT EXISTS Warmdown (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	exerciseID integer
);

CREATE TABLE IF NOT EXISTS Interval (
	id integer PRIMARY KEY,
	position tinyint DEFAULT 0,
	laps tinyint NOT NULL,
    length INTEGER NOT NULL,
	percentage tinyint NOT NULL,
	rest text,
	effort_type varchar(32) DEFAULT "distance",
	effort text DEFAULT "", -- storing time in there
	repeatID integer,
	exerciseID integer);

Create Table IF NOT EXISTS Repeat  (
	id integer PRIMARY KEY,
	repeat tinyint,
	position tinyint DEFAULT 0,
	exerciseID integer
);
`

//TODO: remove
var aSample = `
	DELETE FROM Exercise;
	DELETE FROM Warmup;
	DELETE FROM Warmdown;
	DELETE FROM Interval;
	DELETE FROM Repeat;

	INSERT INTO Exercise(name) VALUES("Pyramids Short");

	INSERT INTO Warmup(effort_type, effort, position, exerciseID) VALUES("distance", "5km very easy around", 0, 1);
	INSERT INTO Warmdown(effort_type, effort, position, exerciseID) VALUES("time", "15 mn footing", 2, 1);

	INSERT INTO Repeat(Repeat, position, exerciseID) VALUES(5, 1, 1);
	INSERT INTO Interval(laps, length, percentage, rest, effort_type, repeatID) VALUES(6, 1000, 90, "400m active", "distance", 1);

`

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

func cleanupSteps(exerciseType string, id int) (res sql.Result, err error) {
	res, err = sqlTX(fmt.Sprintf(`delete from Warmup where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}

	res, err = sqlTX(fmt.Sprintf(`delete from Warmdown where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}
	res, err = sqlTX(fmt.Sprintf(`delete from Interval where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}

func getSteps(exerciseType string, targetID int, steps *[]Step) (err error) {
	getWarmupSQL := fmt.Sprintf(`SELECT id, position, effort,
					 effort_type FROM Warmup
					 WHERE %s=?`, exerciseType)
	rows, err := DB.Query(getWarmupSQL, targetID)
	for rows.Next() {
		var step = Step{
			Type: "warmup",
		}
		err = rows.Scan(
			&step.ID, &step.Position,
			&step.Effort, &step.EffortType)
		if err != nil {
			return
		}
		*steps = append(*steps, step)
	}

	getWarmdownSQL := fmt.Sprintf(`SELECT id, position, effort,
					 effort_type FROM Warmdown
					 WHERE %s=?`, exerciseType)
	rows, err = DB.Query(getWarmdownSQL, targetID)
	for rows.Next() {
		var step = Step{
			Type: "warmdown",
		}
		err = rows.Scan(
			&step.ID, &step.Position,
			&step.Effort, &step.EffortType)
		if err != nil {
			return
		}
		*steps = append(*steps, step)
	}

	getIntervalSQL := fmt.Sprintf(`SELECT id, position, laps, length,
					   percentage, rest, effort_type,
					   effort FROM Interval WHERE %s=?`, exerciseType)

	rows, err = DB.Query(getIntervalSQL, targetID)
	for rows.Next() {

		step := Step{
			Type: "interval",
		}
		err = rows.Scan(&step.ID, &step.Position, &step.Laps,
			&step.Length, &step.Percentage, &step.Rest,
			&step.EffortType, &step.Effort)
		if err != nil {
			return
		}
		*steps = append(*steps, step)
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

	//TODO: cleanup
	getRepeatSQL := `SELECT id, repeat, position from Repeat where exerciseID=?`
	rows, err = DB.Query(getRepeatSQL, targetID)
	for rows.Next() {
		//TODO: cleanup
		step := Step{
			Type: "repeat",
		}

		err = rows.Scan(&step.Repeat.ID, &step.Repeat.Repeat,
			&step.Position)
		if err != nil {
			return
		}

		var repeatSteps []Step
		err = getSteps("repeatID", step.Repeat.ID, &repeatSteps)
		if err != nil {
			return
		}
		step.Repeat.Steps = repeatSteps

		*steps = append(*steps, step)
	}

	return
}

func getExercise(ID int64) (exercise Exercise, err error) {
	var steps []Step

	sql := `SELECT id, name, comment from Exercise where id=?`
	err = DB.QueryRow(sql, ID).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Comment,
	)
	if err != nil {
		return
	}
	err = getSteps("exerciseID", exercise.ID, &steps)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	exercise.Steps = steps

	sort.Sort(exercise)
	return
}

func addExercise(exercise Exercise) (res sql.Result, err error) {
	sql := `insert or replace into Exercise (ID, name, comment) values (?, ?, ?);`
	res, err = sqlTX(sql, exercise.ID, exercise.Name, exercise.Comment)
	if err != nil {
		return
	}

	res, err = cleanupSteps("exerciseID", exercise.ID)
	if err != nil {
		return
	}

	for position, value := range exercise.Steps {
		_, err = addStep(value, "exerciseId", position, exercise.ID)
		if err != nil {
			return
		}
	}

	return
}

func addStep(value Step, exerciseType string, position, targetID int) (
	res sql.Result, err error) {
	if value.Type == "warmup" {
		sql := fmt.Sprintf(`INSERT OR REPLACE into Warmup
			(id, effort_type, effort, position, %s)
			VALUES (?, ?, ?, ?, ?);`, exerciseType)
		res, err = sqlTX(sql, value.ID, value.EffortType,
			value.Effort, position, targetID)
		if err != nil {
			return
		}
	} else if value.Type == "warmdown" {
		sql := fmt.Sprintf(`INSERT OR REPLACE into Warmdown
			(id, effort_type, effort, position, %s)
			VALUES (?, ?, ?, ?, ?);`, exerciseType)
		res, err = sqlTX(sql, value.ID, value.EffortType,
			value.Effort, position, targetID)
		if err != nil {
			return
		}
	} else if value.Type == "interval" {
		sql := fmt.Sprintf(`insert or replace into Interval (
			position, laps,
			length, percentage, rest,
			effort_type, effort, %s) values
			(?, ?, ?,
			 ?, ?, ?,
			 ?, ?)`, exerciseType)
		res, err = sqlTX(sql,
			position, value.Laps,
			value.Length, value.Percentage, value.Rest,
			value.EffortType, value.Effort, targetID)
		if err != nil {
			return
		}
	} else if value.Type == "repeat" {
		sql := `insert or replace into Repeat
				(ID, repeat, position, exerciseId) values
				(?, ?, ?, ?);`
		res, err = sqlTX(sql,
			value.Repeat.ID, value.Repeat.Repeat,
			position, targetID)
		if err != nil {
			return
		}
		for position, value := range value.Repeat.Steps {
			_, err = addStep(value, "repeatId", position, targetID)
			if err != nil {
				return
			}
		}
	}
	return
}
