package server

import (
	"database/sql"
	"fmt"
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
	DELETE FROM Repeat;

	INSERT INTO Exercise(name) VALUES("Pyramids Short");

	INSERT INTO Warmup(effort_type, effort, position, exerciseID) VALUES("distance", "5km very easy around", 0, 1);
	INSERT INTO Warmdown(effort_type, effort, position, exerciseID) VALUES("time", "15 mn footing", 2, 1);

	INSERT INTO Repeat(Repeat, position, exerciseID) VALUES(5, 1, 1);
	INSERT INTO Interval(laps, length, percentage, rest, effort_type, repeatID) VALUES(6, 1000, 90, "400m active", "distance", 1);

`

type ArgsMap map[string]interface{}

func SQLInsertOrUpdate(table string, id int, am ArgsMap) (lastid int, err error) {
	var c int
	var newID int64
	var res sql.Result

	var keys []interface{} = make([]interface{}, 0)
	var values []interface{} = make([]interface{}, 0)
	for k, v := range am {
		keys = append(keys, k)
		values = append(values, v)
	}

	query := "SELECT 1 FROM " + table + " WHERE id=? "
	if am["repeatID"] != nil {
		query += fmt.Sprintf("and repeatID=%d", am["repeatID"].(int))
	} else if am["exerciseID"] != nil {
		query += fmt.Sprintf("and exerciseID=%d", am["exerciseID"].(int))
	}

	var existing int
	err = DB.QueryRow(query, id).Scan(
		&existing,
	)

	if existing == 0 {
		query = "INSERT INTO " + table + "("
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
		newID, _ = res.LastInsertId()
		lastid = int(newID)
		return
	}

	c = 1
	query = "UPDATE " + table + " SET "
	for _, k := range keys {
		query += k.(string) + "=?"
		if c != len(am) {
			query += ", "
		}
		c += 1
	}
	query += fmt.Sprintf(" WHERE ID=%d", id)
	res, err = sqlTX(query, values...)
	lastid = id
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
	return

}
