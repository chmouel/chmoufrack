package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static/")))
	router.PathPrefix("/static/").Handler(s)

	return router
}

func servRest() {
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
