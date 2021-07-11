package main

import (
	crud "backend/controller"
	"backend/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	// "google.golang.org/api/option"

	"log"
	"net/http"
	"os"
	// "cloud.google.com/go/firestore"
)

var connectBot *linebot.Client

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

func lineBot(c *gin.Context) {
	bot := connectBot
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
				} else if message.Text == "carousel" {
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
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func main() {

	e := models.Product{
		Name:    "Sam",
		Barcode: "12312312312321312",
	}
	e.LeavesRemaining()
	// -------- gin router
	// app := App{}
	// app.Init()
	// app.Run()
	// app.addProduct()
	router := gin.Default()

	var port = os.Getenv("PORT")
	if port == "" {
		// is not heroku
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	connectBot = connectLineBot()

	router.POST("/callback", lineBot)

	productAPI := router.Group("/product")
	{
		productAPI.POST("/create", func(c *gin.Context) {
			var form models.ProductHasExpire
			if c.ShouldBind(&form) == nil {
				result := crud.AddData(form)
				if result != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"message": result.Error(),
					})
					return
				}
				c.JSON(http.StatusCreated, gin.H{
					"message": "Add " + form.Barcode + " OK!",
				})
			} else {
				c.JSON(401, gin.H{"status": "unable to bind data"})
			}
		})
	}

	expireAPI := router.Group("/expire")
	{
		expireAPI.POST("", func(c *gin.Context) {
			var form models.ProductHasExpire
			if c.ShouldBind(&form) == nil {
				result := crud.AddExpire(form)
				if result != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"message": result.Error(),
					})
					return
				}
				// c.JSON(http.StatusCreated, gin.H{
				// 	"message": "Add " + form.Barcode + " OK!",
				// })
			} else {
				c.JSON(401, gin.H{"status": "unable to bind data"})
			}
			c.JSON(http.StatusCreated, nil)
		})
		expireAPI.DELETE("", func(c *gin.Context) {
			var form models.ProductHasExpire
			if c.ShouldBind(&form) == nil {
				fmt.Printf("Delete barcode %v %v \n", form.Barcode, form.Quantity)
				result := crud.RemoveExpire(form)
				if result != nil {
					c.JSON(http.StatusUnprocessableEntity, gin.H{
						"message": result.Error(),
					})
					return
				}
			} else {
				c.JSON(401, gin.H{"status": "unable to bind data"})
			}
			c.JSON(http.StatusNoContent, nil)
		})
	}

	// router.Run(":" + os.Getenv("PORT"))

	if port == "" {
		fmt.Println("Running on Heroku using random PORT")
		router.Run()
	} else {
		fmt.Println("Environment Port : " + port)
		router.Run(":" + port)
	}
}
