package server

import (
	"database/sql"
	"fmt"
	"sort"
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
	_, err = sqlTX(fmt.Sprintf(`delete from Intervals where %s=?`, exerciseType), id)
	if err != nil {
		fmt.Println(err.Error())
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

	_, err = sqlTX(fmt.Sprintf(`delete from Repeats where %s=?`, exerciseType), id)
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

	query := fmt.Sprintf(`SELECT id, position, laps, length,
					   percentage, rest, effort_type,
					   effort FROM Intervals WHERE %s=?`, exerciseType)
	getIntervalSQL := query
	rows, err = DB.Query(getIntervalSQL, targetID)
	for rows.Next() {
		var nEffort sql.NullString
		step := Step{
			Type: "interval",
		}
		err = rows.Scan(&step.ID, &step.Position, &step.Laps,
			&step.Length, &step.Percentage, &step.Rest,
			&step.EffortType, &nEffort)

		if nEffort.Valid {
			step.Effort = nEffort.String
		}

		if err != nil {
			fmt.Printf("%+v\n", step.Effort)
			return
		}
		*steps = append(*steps, step)
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

	//TODO: cleanup
	getRepeatSQL := `SELECT id, repeats, position from Repeats where exerciseID=?`
	rows, err = DB.Query(getRepeatSQL, targetID)
	for rows.Next() {
		//TODO: cleanup
		step := Step{
			Type: "repeat",
		}

		err = rows.Scan(&step.Repeat.ID, &step.Repeat.Repeats,
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

func getIdOfExerciseName(name string) (id int, err error) {
	sqlT := `SELECT id from Exercise where name=?`
	err = DB.QueryRow(sqlT, name).Scan(
		&id,
	)

	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = &error404{"Exercise " + name + " Not Found"}
	}
	return
}

func deleteExercise(ID int) (err error) {
	sqlT := `DELETE From Exercise where id=?`
	_, err = sqlTX(sqlT, ID)
	return
}

func getExercise(ID int) (exercise Exercise, err error) {
	var steps []Step
	var nComment sql.NullString

	sqlT := `SELECT id, name, comment from Exercise where id=?`
	err = DB.QueryRow(sqlT, ID).Scan(
		&exercise.ID,
		&exercise.Name,
		&nComment,
	)

	if nComment.Valid {
		exercise.Comment = nComment.String
	}

	if err != nil {
		if err != sql.ErrNoRows {
			return
		}
		err = &error404{"Exercise " + exercise.Name + " Not Found"}
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

func AddExercise(exercise Exercise) (lastid int, err error) {
	if exercise.Name != "" {
		sqlT := `SELECT id from Exercise where name=?`
		err = DB.QueryRow(sqlT, exercise.Name).Scan(
			&exercise.ID,
		)
		if err != nil && err != sql.ErrNoRows {
			return
		}
	}

	am := ArgsMap{
		"name":    exercise.Name,
		"comment": exercise.Comment,
	}
	lastid, err = SQLInsertOrUpdate("Exercise", exercise.ID, am)
	if err != nil {
		return
	}

	err = cleanupSteps("exerciseID", lastid)
	if err != nil {
		return
	}

	for position, value := range exercise.Steps {
		err = addStep(value, "exerciseID", position, lastid)
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
		_, err = SQLInsertOrUpdate("Intervals", value.ID, am)

		if err != nil {
			return
		}

	} else if value.Type == "repeat" {
		var lastid int
		am := ArgsMap{
			"position": position,
			"repeats":  value.Repeat.Repeats}
		am[exerciseType] = targetID
		lastid, err = SQLInsertOrUpdate("Repeats", value.Repeat.ID, am)
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
		e, err = getExercise(e.ID)
		if err != nil {
			return
		}
		exercises = append(exercises, e)
	}
	return
}
