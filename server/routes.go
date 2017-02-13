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
		"GETExcercise",
		"GET",
		"/v1/excercise/{name}",
		GETExcercise,
	},
	route{
		"PostExcercise",
		"POST",
		"/v1/excercise/{name}",
		POSTExcercise,
	},
	// TODO: remove
	route{
		"GETExcerciseS",
		"GET",
		"/v1/excercises",
		GETExcerciseS,
	},
}
