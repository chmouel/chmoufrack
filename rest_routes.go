package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"ProgramsIndex",
		"GET",
		"/programs",
		RESTProgramsIndex,
	},
	Route{
		"ProgramCreate",
		"POST",
		"/program",
		RESTProgramCreate,
	},
	Route{
		"ProgramCleanup",
		"DELETE",
		"/program/purge/{name}",
		RESTProgramCleanup,
	},
	Route{
		"WorkoutCreate",
		"POST",
		"/workouts",
		RESTMultipleWorkoutsCreate,
	},
}
