package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func router(staticDir string) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range allRoutes {
		var handler http.Handler
		handler = route.HandlerFunc
		router.Methods(route.Method).Path(route.Pattern).
			Name(route.Name).Handler(handler)
	}

	s := http.StripPrefix("/", http.FileServer(http.Dir(staticDir)))
	router.PathPrefix("/").Handler(s)

	logRouter := handlers.LoggingHandler(os.Stdout, router)
	return logRouter
}

func Serve(staticDir string, port int) {
	router := router(staticDir)
	sPort := fmt.Sprintf(":%d", port)
	fmt.Printf("Serving on port %d with static dir %s\n", port, staticDir)
	log.Fatal(http.ListenAndServe(sPort, router))
}
