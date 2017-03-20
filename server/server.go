package server

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/huandu/facebook"
	"gopkg.in/gin-gonic/gin.v1"
)

type ACLCheck interface {
	FBGet(url, token string) (fbc facebook.Result, err error)
}

type FBCheck struct{}

func (fbcheck *FBCheck) FBGet(url, token string) (fbc facebook.Result, err error) {
	fbc, err = facebook.Get(url, facebook.Params{
		"access_token": token,
		"fields":       "name,email,link",
	})
	return
}

// Really need to find a way how to test that,
func Check(aclcheck ACLCheck) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Request.Header["Authorization"]) == 0 {
			c.Next()
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

		fb, err := aclcheck.FBGet("/me", token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var fbInfo FBinfo
		err = fb.Decode(&fbInfo)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Error while decoding: " + err.Error()})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if fb.Get("id") == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Could not get fields from facebookInfo"})
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fbInfo.ID = fb.Get("id").(string)

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

	v1 := router.Group("/v1", Check(acl))
	{
		v1.POST("/fbinfo", POSTFbinfo)
		v1.POST("/exercise", POSTExercise)
		v1.DELETE("/exercise/:id", DeleteExercise)
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
