package cmd

import (
	"net/http"

	"github.com/fakeyanss/usg-go/internal/api"
	"github.com/gin-gonic/gin"
)

func Main() error {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	group := router.Group("")
	api.RegisterS3Api(group)

	router.GET("/favicon.ico", favicon)

	return router.Run(":8000")
}

// favicon return a empty 200 reponse
func favicon(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.Header("Link", "rel=\"shortcut icon\" href=\"#\"")
	c.Status(http.StatusOK)
}
