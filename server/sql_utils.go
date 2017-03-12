package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var sqlTable = `
CREATE TABLE IF NOT EXISTS FBinfo (
	ID int NOT NULL AUTO_INCREMENT,
	FBId bigint not null,
	name varchar(255) not null,
	link varchar(255) not null,
	email varchar(255) not null,
	PRIMARY KEY(ID),
	CONSTRAINT uc_U UNIQUE (FBid,name,link)
);

CREATE TABLE IF NOT EXISTS Exercise (
	ID int NOT NULL AUTO_INCREMENT,
	name varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    comment text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
	fbID bigint NOT NULL,
    PRIMARY KEY(ID),
	CONSTRAINT uc_U UNIQUE (ID,name,fbID)
);

CREATE TABLE IF NOT EXISTS Warmup (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
	repeatID integer,
	exerciseID integer,
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Warmdown (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	effort_type varchar(32) DEFAULT "distance",
    effort text NOT NULL,
	repeatID integer,
	exerciseID integer,
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Intervals (
	id int NOT NULL AUTO_INCREMENT,
	position tinyint DEFAULT 0,
	laps tinyint NOT NULL,
    length INTEGER NOT NULL,
	percentage tinyint NOT NULL,
	rest text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
	effort_type varchar(32) DEFAULT "distance",
	effort text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
	repeatID integer,
	exerciseID integer,
    CHECK(repeatID is not NULL or exerciseID is not NULL),
    PRIMARY KEY(ID)
);

CREATE TABLE IF NOT EXISTS Repeats  (
	id int NOT NULL AUTO_INCREMENT,
	repeats tinyint,
	position tinyint DEFAULT 0,
	exerciseID integer,
    PRIMARY KEY(ID)
);
`

//TODO: remove
var SQLDropTable = `
SET FOREIGN_KEY_CHECKS = 0;
SET GROUP_CONCAT_MAX_LEN=32768;
SET @tables = NULL;
SELECT GROUP_CONCAT(table_name) INTO @tables
  FROM information_schema.tables
  WHERE table_schema = (SELECT DATABASE());
SELECT IFNULL(@tables,'dummy') INTO @tables;

SET @tables = CONCAT('DROP TABLE IF EXISTS ', @tables);
PREPARE stmt FROM @tables;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
SET FOREIGN_KEY_CHECKS = 1;`

var SQLresetDB = `
	DELETE FROM Exercise;
	DELETE FROM Warmup;
	DELETE FROM Warmdown;
	DELETE FROM Intervals;
	DELETE FROM Repeats;`

type ArgsMap map[string]interface{}

func createSampleExercise(exerciceName, warmupEffort, warmdownEffort string, length int, facebookid string) (e Exercise) {
	var steps Steps

	fbinfo := FBinfo{
		ID:    facebookid,
		Name:  "Chmou EL",
		Link:  "https://www.facebook.com/app_scoped_user_id/10157827176205251/",
		Email: "email@email.com",
	}

	step1 := Step{
		Type:       "warmup",
		Effort:     warmupEffort,
		EffortType: "distance",
	}
	steps = append(steps, step1)

	step2 := Step{
		Laps:       3,
		Length:     length,
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
		FB:      fbinfo,
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
		begin = "REPLACE INTO "
	} else {
		begin = "INSERT INTO "
	}

	query = begin + table + "("
	c = 1
	for _, k := range keys {
		query += k.(string)
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

func DBConnect(dbconnection string, reset string) (err error) {
	DB, err = sql.Open("mysql", dbconnection+"?multiStatements=true&collation=utf8mb4_bin")

	if err != nil {
		return
	}

	if reset != "" {
		_, err = DB.Exec(SQLDropTable)
		if err != nil {
			return
		}
	}

	_, err = DB.Exec(sqlTable)
	return
}

func InitFixturesDB(facebookid string) (err error) {
	_, err = DB.Exec(SQLresetDB)
	e := createSampleExercise("Test1", "easy warmup todoo", "finish strong", 1000, facebookid)
	var repeatSteps Steps
	repeatStep := Step{
		Laps:       6,
		Length:     400,
		Percentage: 100,
		Type:       "interval",
		EffortType: "distance",
	}
	repeatSteps = append(repeatSteps, repeatStep)

	repeat := Repeats{
		Steps:   repeatSteps,
		Repeats: 5,
	}
	exerciseStep := Step{
		Type:   "repeat",
		Repeat: repeat,
	}
	e.Steps = append(e.Steps, exerciseStep)

	_, err = addExercise(e)
	if err != nil {
		return
	}
	return
}
