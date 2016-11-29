package rest

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
		GETPrograms,
	},
	Route{
		"ProgramShow",
		"GET",
		"/program/{name}",
		GETProgram,
	},
	Route{
		"ProgramCreate",
		"POST",
		"/program",
		CreateProgram,
	},
	Route{
		"ProgramCleanup",
		"DELETE",
		"/program/purge/{name}",
		CleanupProgram,
	},
	Route{
		"WorkoutCreate",
		"POST",
		"/workouts",
		CreateMultipleWorkouts,
	},
}
