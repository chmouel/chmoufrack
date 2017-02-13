package server

import "net/http"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var allRoutes = routes{
	route{
		"GETExercise",
		"GET",
		"/v1/exercise/{name}",
		GETExercise,
	},
	route{
		"PostExercise",
		"POST",
		"/v1/exercise/{name}",
		POSTExercise,
	},
	// TODO: remove
	route{
		"GETExerciseS",
		"GET",
		"/v1/exercises",
		GETExerciseS,
	},
}
