package cmd

import (
	"github.com/fakeyanss/usg-go/internal/api"
	"github.com/gin-gonic/gin"
)

func Main() error {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	group := router.Group("")
	api.RegisterS3Api(group)

	return router.Run(":8000")
}
