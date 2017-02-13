package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func logger(inner http.Handler, name string) http.Handler {
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

func router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range allRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).
			Name(route.Name).Handler(handler)
	}

	s := http.StripPrefix("/", http.FileServer(http.Dir("../client")))
	router.PathPrefix("/").Handler(s)

	return router
}

func Serve() {
	router := router()
	// TODO(chmou): setting
	log.Fatal(http.ListenAndServe(":8080", router))
}