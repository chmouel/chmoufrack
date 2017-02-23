package server

import "database/sql"

var sqlTable = `
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
	exerciseID integer,
    CHECK(repeatID is not NULL or exerciseID is not NULL)
)
;

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
	DELETE FROM Repeat;`

type ArgsMap map[string]interface{}

func createSampleExercise(
	exerciceName, warmupEffort, warmdownEffort string,
	length2 int) (e Exercise) {
	var steps Steps

	step1 := Step{
		Type:       "warmup",
		Effort:     warmupEffort,
		EffortType: "distance",
	}
	steps = append(steps, step1)

	step2 := Step{
		Laps:       3,
		Length:     length2,
		Percentage: 90,
		Type:       "interval",
		EffortType: "distance",
	}
	steps = append(steps, step2)

	step3 := Step{
		Effort:     warmdownEffort,
		Type:       "warmdown",
		EffortType: "distance",
	}
	steps = append(steps, step3)

	e = Exercise{
		Name:    exerciceName,
		Comment: "NoComment",
		Steps:   steps,
	}
	return
}

func SQLInsertOrUpdate(table string, id int, am ArgsMap) (lastid int, err error) {
	var c int
	var res sql.Result
	var begin, query string

	var keys []interface{} = make([]interface{}, 0)
	var values []interface{} = make([]interface{}, 0)
	for k, v := range am {
		keys = append(keys, k)
		values = append(values, v)
	}

	if id != 0 {
		begin = "INSERT OR REPLACE INTO "
	} else {
		begin = "INSERT INTO "
	}

	query = begin + table + "("
	c = 1
	for _, k := range keys {
		query += `"` + k.(string) + `"`
		if c != len(am) {
			query += ","
		}
		c += 1
	}
	query += ") VALUES ("
	c = 1
	for range keys {
		query += `?`
		if c != len(am) {
			query += ","
		}
		c += 1
	}
	query += ");"
	res, err = sqlTX(query, values...)
	if err != nil {
		return
	}

	n, _ := res.LastInsertId()
	lastid = int(n)
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

func DBConnect(location string) (err error) {
	// TODO: proper sqlite location
	DB, err = sql.Open("sqlite3", location)
	if err != nil {
		return
	}

	_, err = DB.Exec(sqlTable)
	return
}

func InitFixturesDB() (err error) {
	_, err = DB.Exec(aSample)

	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1234)

	var repeatSteps Steps
	repeatStep := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep)

	repeat := Repeat{
		Steps:  repeatSteps,
		Repeat: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)

	_, err = AddExercise(e)
	if err != nil {
		return
	}
	return

}
