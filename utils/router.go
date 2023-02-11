package utils

import (
	"github.com/gin-gonic/gin"
)

var MainRouter *gin.Engine

func InitMainRouter() {
	if PRODUCTION {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		SendErrorResponse(c, NOT_FOUND_ERROR)
	})
	router.NoMethod(func(c *gin.Context) {
		SendErrorResponse(c, INVALID_REQUEST_METHOD_ERROR)
	})

	MainRouter = router
}
