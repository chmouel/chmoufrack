package server

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"sort"
)

func cleanupSteps(exerciseType string, id int) (err error) {
	_, err = sqlTX(fmt.Sprintf(`delete from Warmup where %s=?`, exerciseType), id)
	if err != nil {
		return
	}

	_, err = sqlTX(fmt.Sprintf(`delete from Warmdown where %s=?`, exerciseType), id)
	if err != nil {
		return
	}
	_, err = sqlTX(fmt.Sprintf(`delete from Intervals where %s=?`, exerciseType), id)
	if err != nil {
		return
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

	_, err = sqlTX(fmt.Sprintf(`delete from Repeats where %s=?`, exerciseType), id)
	if err != nil {
		return
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
			return
		}
		*steps = append(*steps, step)
	}

	if exerciseType == "repeatID" {
		// We don't need to do the repeat stuff if we are in repeat loop
		return
	}

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

func deleteExercise(e Exercise) (err error) {
	sqlT := `DELETE From Exercise where id=?`
	_, err = sqlTX(sqlT, e.ID)
	return
}

func getExercise(ID int, FBId string) (exercise Exercise, err error) {
	var steps []Step
	var nComment sql.NullString
	var strPublic string

	// Note: Do not get email here so to avoid spam, maybe we'll use email another time
	sqlT := `SELECT e.id,e.name,e.comment,e.public,f.fbid,f.name as fbname,f.link FROM Exercise e LEFT JOIN FBinfo f on f.FBId=e.FBId WHERE e.id=? and (public="1" or f.fbid=?)`

	err = DB.QueryRow(sqlT, ID, FBId).Scan(
		&exercise.ID,
		&exercise.Name,
		&nComment,
		&strPublic,
		&exercise.FB.ID,
		&exercise.FB.Name,
		&exercise.FB.Link,
	)

	if strPublic == "0" {
		exercise.Public = false
	} else {
		exercise.Public = true
	}

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
		return
	}
	exercise.Steps = steps

	sort.Sort(exercise)
	return
}

func checkBadCharacters(s string) (err error) {
	var re = regexp.MustCompile(`[/?&].*`)
	if re.MatchString(s) {
		return &errorUnauthorized{"You have invalid characters"}
	}
	return
}

func convertBoolToZeroUn(c bool) (ret string) {
	ret = "0"
	if c {
		ret = "1"
	}
	return
}

func addFBInfo(fb FBinfo) (lastid int, err error) {
	am := ArgsMap{
		"FBid":  fb.ID,
		"name":  fb.Name,
		"link":  fb.Link,
		"email": fb.Email,
	}
	// TODO(chmou): figure out a better way than just a 1
	lastid, err = SQLInsertOrUpdate("FBinfo", 1, am)

	return
}

func addExercise(exercise Exercise) (lastid int, err error) {
	var oldFbID string
	var oldId int

	if exercise.Name == "" {
		return -1, errors.New("You need to specify an exercise Name")
	}
	sqlT := `SELECT id,fbID from Exercise where name=?`
	err = DB.QueryRow(sqlT, exercise.Name).Scan(
		&oldId,
		&oldFbID,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if oldFbID != "" && exercise.FB.ID != oldFbID {
		return -1, &errorUnauthorized{"You are not allowed to update other user exercise"}
	}

	err = checkBadCharacters(exercise.Name)
	if err != nil {
		return
	}

	am := ArgsMap{
		"name":    exercise.Name,
		"comment": exercise.Comment,
		"fbID":    exercise.FB.ID,
		"public":  convertBoolToZeroUn(exercise.Public),
	}

	if oldId != 0 {
		am["id"] = oldId
	}

	lastid, err = SQLInsertOrUpdate("Exercise", oldId, am)
	if err != nil {
		return
	}

	am = ArgsMap{
		"FBid":  exercise.FB.ID,
		"name":  exercise.FB.Name,
		"link":  exercise.FB.Link,
		"email": exercise.FB.Email,
	}
	// TODO(chmou): figure out a better way than just a 1
	_, err = SQLInsertOrUpdate("FBinfo", 1, am)
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

func getAllExercises(fbID string) (exercises []Exercise, err error) {
	var getAllExercises = `SELECT ID from Exercise where (public="1" or fbId=?);`
	rows, err := DB.Query(getAllExercises, fbID)
	for rows.Next() {
		e := Exercise{}
		err = rows.Scan(&e.ID)
		if err != nil {
			return
		}

		e, err = getExercise(e.ID, fbID)
		if err != nil {
			return
		}
		exercises = append(exercises, e)
	}
	return
}
