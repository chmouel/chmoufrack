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

func POSTFbinfo(c *gin.Context) {
	var fbinfo FBinfo
	var tokenFBInfo FBinfo

	v, exist := c.Get("FBInfo")
	if !exist {
		handle_error_nf_bad(c,
			&errorUnauthorized{"Why FBinfo do not exist, this should not happen"})
		return
	}

	if err := c.Bind(&fbinfo); err != nil {
		handle_error_nf_bad(c, err)
		return
	}

	tokenFBInfo = v.(FBinfo)
	if tokenFBInfo.ID != fbinfo.ID {
		handle_error_nf_bad(c,
			&errorUnauthorized{"You are not allowed to update other people FB infos"})
		return
	}

	_, err := addFBInfo(fbinfo)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func POSTExercise(c *gin.Context) {
	var exercise Exercise
	var err error

	v, exist := c.Get("FBInfo")
	if !exist {
		handle_error_nf_bad(c,
			&errorUnauthorized{"Why FBinfo do not exist, this should not happen"})
		return
	}

	if err := c.Bind(&exercise); err != nil {
		handle_error_nf_bad(c, err)
		return
	}

	// Needs to be after the bind
	exercise.FB = v.(FBinfo)

	_, err = addExercise(exercise)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}
	c.Status(http.StatusCreated)
}

func DeleteExercise(c *gin.Context) {
	var err error
	var i int
	var fb FBinfo
	exerciseID := c.Param("id")

	v, exist := c.Get("FBInfo")
	if !exist {
		handle_error_nf_bad(c,
			&errorUnauthorized{"Why FBinfo do not exist, this should not happen"})
		return
	}
	fb = v.(FBinfo)

	if i, err = strconv.Atoi(exerciseID); err != nil {
		i, err = getIdOfExerciseName(exerciseID)
		if err != nil {
			handle_error_nf_bad(c, err)
			return
		}
	}

	e, err := getExercise(i)
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}

	if e.FB.ID != fb.ID {
		handle_error_nf_bad(c,
			&errorUnauthorized{"You have no right to delete this exercise"})
		return
	}

	err = deleteExercise(e)
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

	exercises, err = getAllExercises()
	if err != nil {
		handle_error_nf_bad(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, exercises)
}
