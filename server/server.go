package server

import (
	"fmt"
	"log"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

func setupRoutes(staticDir string) *gin.Engine {
	router := gin.Default()

	// Not ideal compared to the old imp on direct net/http since I cannot have
	// a / catchall for StaticFS
	router.Static("/html", staticDir)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/html")
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/exercise", POSTExercise)
		v1.DELETE("/exercise/:id", DeleteExercise)
		v1.GET("/exercise/:id", GETExercise)
		v1.GET("/exercises", GETExercises)
	}
	return router
}

func Serve(staticDir string, port int) {
	sPort := fmt.Sprintf(":%d", port)
	router := setupRoutes(staticDir)

	fmt.Printf("Serving on %s with static dir %s\n", sPort, staticDir)
	if err := router.Run(sPort); err != nil {
		log.Fatal(err)
	}
}
