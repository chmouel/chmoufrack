package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	fb "github.com/huandu/facebook"
	"gopkg.in/gin-gonic/gin.v1"
)

type ACLCheck interface {
	Check() gin.HandlerFunc
}

type FBCheck struct{}

// Really need to find a way how to test that,
func (fbcheck *FBCheck) Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Header["Authorization"]) == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to have an authorization header in your request"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := c.Request.Header["Authorization"][0]

		if len(token) > 6 && strings.ToUpper(token[0:6]) == "BEARER" {
			token = token[7:]
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You need to specify a token for this query"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_fb, err := fb.Get("/me", fb.Params{
			"access_token": token,
			"fields":       "name,email,link",
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var fbInfo FBinfo
		err = _fb.Decode(&fbInfo)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error while decoding: " + err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fbInfo.ID = _fb.Get("id").(string)

		c.Set("FBInfo", fbInfo)
		c.Next()
	}
}

func setupRoutes(staticDir string, acl ACLCheck) *gin.Engine {
	router := gin.Default()

	// Not ideal compared to the old imp on direct net/http since I cannot have
	// a / catchall for StaticFS
	router.Static("/html", staticDir)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/html")
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/fbinfo", acl.Check(), POSTFbinfo)
		v1.POST("/exercise", acl.Check(), POSTExercise)
		v1.DELETE("/exercise/:id", acl.Check(), DeleteExercise)
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

	fbcheck := &FBCheck{}
	router := setupRoutes(staticDir, fbcheck)

	fmt.Printf("Serving on %s with static dir %s\n", sPort, staticDir)
	if err := router.Run(sPort); err != nil {
		log.Fatal(err)
	}
}
