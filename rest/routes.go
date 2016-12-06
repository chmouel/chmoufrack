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
		"/rest/programs",
		GETPrograms,
	},
	Route{
		"ProgramWorkoutGet",
		"GET",
		"/rest/program/{name}/workouts",
		GetWorkoutsForProgram,
	},
	Route{
		"ProgramWorkoutCreate",
		"POST",
		"/rest/program/{name}/workouts",
		CreateMultipleWorkouts,
	},
	Route{
		"ProgramWorkoutsGetFull",
		"GET",
		"/rest/program/{name}/full",
		GetProgramFull,
	},
	Route{
		"ProgramCreate",
		"POST",
		"/rest/program/{name}",
		CreateProgram,
	},
	Route{
		"ProgramCleanup",
		"DELETE",
		"/rest/program/{name}/purge",
		CleanupProgram,
	},
	Route{
		"HTMLProgramShow",
		"GET",
		"/html/program/{name}",
		HTMLProgramShow,
	},
}
