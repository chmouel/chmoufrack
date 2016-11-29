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
		"ProgramWorkoutGet",
		"GET",
		"/program/{name}/workouts",
		GetWorkoutsForProgram,
	},
	Route{
		"ProgramWorkoutCreate",
		"POST",
		"/program/{name}/workouts",
		CreateMultipleWorkouts,
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
		"/program/{name}/purge",
		CleanupProgram,
	},
}
