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
		"GetPrograms",
		"GET",
		"/rest/programs",
		GETPrograms,
	},
}
