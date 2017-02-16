package server

import (
	"database/sql"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

func cleanupSteps(exerciseType string, id int) (err error) {
	_, err = sqlTX(fmt.Sprintf(`delete from Warmup where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = sqlTX(fmt.Sprintf(`delete from Warmdown where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = sqlTX(fmt.Sprintf(`delete from Interval where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

	_, err = sqlTX(fmt.Sprintf(`delete from Repeat where %s=?`, exerciseType), id)
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

		sort.Sort(step.Repeat.Steps)

		*steps = append(*steps, step)
	}

	return
}

func getExercise(ID int64) (exercise Exercise, err error) {
	var steps []Step

	sqlT := `SELECT id, name, comment from Exercise where id=?`
	err = DB.QueryRow(sqlT, ID).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Comment,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = &error404{"Exercise Not Found"}
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

	err = cleanupSteps("exerciseID", exercise.ID)
	if err != nil {
		return
	}

	for position, value := range exercise.Steps {
		err = addStep(value, "exerciseID", position, exercise.ID)
		if err != nil {
			return
		}
	}

	return
}

func addStep(value Step, exerciseType string, position, targetID int) (err error) {
	if value.Type == "warmup" {
		am := ArgsMap{
			"effort_type": value.EffortType,
			"effort":      value.Effort,
			"position":    position,
		}
		am[exerciseType] = targetID

		_, err = SQLInsertOrUpdate("Warmup", value.ID, am)
		if err != nil {
			return
		}
	} else if value.Type == "warmdown" {
		am := ArgsMap{
			"effort_type": value.EffortType,
			"effort":      value.Effort,
			"position":    position,
		}
		am[exerciseType] = targetID

		_, err = SQLInsertOrUpdate("Warmdown", value.ID, am)
		if err != nil {
			return
		}
	} else if value.Type == "interval" {
		am := ArgsMap{
			"position":    position,
			"laps":        value.Laps,
			"length":      value.Length,
			"percentage":  value.Percentage,
			"rest":        value.Rest,
			"effort_type": value.EffortType,
			"effort":      value.Effort}
		am[exerciseType] = targetID
		_, err = SQLInsertOrUpdate("Interval", value.ID, am)

		if err != nil {
			return
		}

	} else if value.Type == "repeat" {
		var lastid int
		am := ArgsMap{
			"position": position,
			"repeat":   value.Repeat.Repeat}
		am[exerciseType] = targetID
		lastid, err = SQLInsertOrUpdate("Repeat", value.Repeat.ID, am)
		if err != nil {
			return
		}

		err = cleanupSteps("repeatID", value.Repeat.ID)
		if err != nil {
			return
		}

		for position, value := range value.Repeat.Steps {
			err = addStep(value, "repeatID", position, lastid)
			if err != nil {
				return
			}
		}
	}
	return
}

func getAllExercises() (exercises []Exercise, err error) {
	var getAllExercises = `SELECT ID from Exercise`
	rows, err := DB.Query(getAllExercises)

	for rows.Next() {
		e := Exercise{}
		err = rows.Scan(&e.ID)
		if err != nil {
			return
		}
		e, err = getExercise(int64(e.ID))
		if err != nil {
			return
		}
		exercises = append(exercises, e)
	}
	return
}
