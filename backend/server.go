package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"net/http"
	"os"
)

func connectLineBot() *linebot.Client {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	return bot
}

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	bot := connectLineBot()

	r.POST("/callback", func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.JSON(400, gin.H{
					"message": "Error",
				})
			} else {
				c.JSON(500, gin.H{
					"message": "Error",
				})
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if message.Text == "hello" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("สวัสดีค่ะ")).Do(); err != nil {
							log.Print(err)
						}
					} else if message.Text == "เช็ควันหมดอายุ" {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("รออีกแปปน๊า")).Do(); err != nil {
							log.Print(err)
						}
					}
					// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					// 	log.Print(err)
					// }
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	"message": "callback",
		// })
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
	r.Run(":" + os.Getenv("PORT"))
}
