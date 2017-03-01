package server

import (
	"net/http"
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func handle_error_nf_bad(c *gin.Context, err error) {
	if _, ok := err.(*error404); ok {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if _, ok := err.(*errorUnauthorized); ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func POSTExercise(c *gin.Context) {
	var exercise Exercise
	var err error
	var fbID int

	if c.Query("fbID") != "" {
		if fbID, err = strconv.Atoi(c.Query("fbID")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	}

	if err := c.Bind(&exercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exercise.FBid = fbID

	_, err = AddExercise(exercise)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func DeleteExercise(c *gin.Context) {
	var err error
	var i, id int
	exerciseID := c.Param("id")

	if i, err = strconv.Atoi(exerciseID); err == nil {
		id = i
	} else {
		id, err = getIdOfExerciseName(exerciseID)

	}

	err = deleteExercise(id)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func GETExercise(c *gin.Context) {
	var exercise Exercise
	var err error
	var i, id int
	exerciseID := c.Param("id")

	if i, err = strconv.Atoi(exerciseID); err == nil {
		id = i
	} else {
		id, err = getIdOfExerciseName(exerciseID)
		if err != nil {
			handle_error_nf_bad(c, err)
			return
		}
	}

	exercise, err = getExercise(id)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, exercise)
}

func GETExercises(c *gin.Context) {
	var exercises []Exercise
	var err error

	exercises, err = GetAllExercises()
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, exercises)
}
