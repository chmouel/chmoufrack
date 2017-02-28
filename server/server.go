package server

import (
	"fmt"
	"log"
	"net/http"

	fb "github.com/huandu/facebook"
	"gopkg.in/gin-gonic/gin.v1"
)

func FBCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		fbid := c.Query("FBid")
		token := c.Query("FBtoken")

		if fbid == "" || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to specify a fbid or token for this query"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err := fb.Get("/"+fbid, fb.Params{
			"access_token": token,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

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
		v1.GET("/test", FBCheck())
		v1.POST("/exercise", FBCheck(), POSTExercise)
		v1.DELETE("/exercise/:id", FBCheck(), DeleteExercise)
		v1.GET("/exercise/:id", GETExercise)
		v1.GET("/exercises", GETExercises)
	}
	return router
}

func Serve(staticDir string, port int, debug bool) {
	sPort := fmt.Sprintf(":%d", port)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := setupRoutes(staticDir)

	fmt.Printf("Serving on %s with static dir %s\n", sPort, staticDir)
	if err := router.Run(sPort); err != nil {
		log.Fatal(err)
	}
}
