package api

import (
	"backend/db"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {

	ctx := context.Background()
	client := db.Init(ctx)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to iTUYEN",
		})
	})

	SetupCronjob(router)
	SetupExpiryAPI(router, client, ctx)
	SetupProductAPI(router, client, ctx)
	SetupLineCallback(router)
	// connectLineBot(client, ctx)
}
