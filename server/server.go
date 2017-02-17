package server

import (
	"fmt"
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

func router(staticDir string) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range allRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).
			Name(route.Name).Handler(handler)
	}

	s := http.StripPrefix("/", http.FileServer(http.Dir(staticDir)))
	router.PathPrefix("/").Handler(s)

	return router
}

func Serve(staticDir string, port int) {
	router := router(staticDir)
	// TODO(chmou): setting
	fmt.Printf("Serving on port %d with static dir %s\n", port, staticDir)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
