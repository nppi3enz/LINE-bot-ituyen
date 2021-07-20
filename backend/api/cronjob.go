package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func cronJob(c *gin.Context) {
	bot := connectBot
	flexMessage := expiryDashboard()
	var err error

	// Reply Message
	if _, err = bot.BroadcastMessage(flexMessage).Do(); err != nil {
		log.Print(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"result": "complete",
	})

}

func SetupCronjob(router *gin.Engine) {
	router.GET("/cronjob", cronJob)
}
