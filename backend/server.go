package main

import (
	crud "backend/controller"
	"backend/models"
	"fmt"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	// "google.golang.org/api/option"

	"log"
	"net/http"
	"os"

	// "cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
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
func expiryDashboard() *linebot.FlexMessage {
	result := crud.ListExpiry(make(map[string]string))
	var outputMsg string
	now := time.Now()

	for _, val := range result {
		dateExpiry := val.ExpireDate
		diff := dateExpiry.Sub(now).Hours()
		calculateDay := math.Ceil(diff / 24)

		nameProduct := val.Product["name"]

		var color string
		var iconText string
		if calculateDay > 10 {
			color = "#aaaaaa"
		} else if 7 < calculateDay && calculateDay <= 10 {
			color = "#FFA900"
			iconText = `{ "type": "icon", "url": "https://ituyen.herokuapp.com/alert-yellow.png" },`
		} else if 4 < calculateDay && calculateDay <= 7 {
			color = "#FF7600"
			iconText = `{ "type": "icon", "url": "https://ituyen.herokuapp.com/alert-orange.png" },`
		} else if 0 < calculateDay && calculateDay <= 4 {
			color = "#CD113B"
			iconText = `{ "type": "icon", "url": "https://ituyen.herokuapp.com/alert-red.png" },`
		} else {
			color = "#52006A"
			iconText = `{ "type": "icon", "url": "https://ituyen.herokuapp.com/alert-purple.png" },`
		}

		outputMsg += `{
			  "type": "box",
			  "layout": "baseline",
			  "contents": [
				` + iconText + `
				{
				  "type": "text",
				  "text": "` + nameProduct.(string) + `",
				  "flex": 3,
				  "weight": "bold",
				  "margin": "sm",
				  "wrap": true
				},
				{
				  "type": "text",
				  "text": "` + fmt.Sprint(calculateDay) + ` Days",
				  "flex": 1,
				  "size": "sm",
				  "color": "` + color + `",
				  "align": "end"
				}
			  ],
			  "height": "30px"
			},`
	}
	inputFmt := outputMsg[:len(outputMsg)-1]
	currentTime := fmt.Sprintf("%02d/%02d/%d", now.Day(), now.Month(), now.Year())
	jsonOutput := `{
		"type": "bubble",
		"size": "mega",
		"header": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "box",
			  "layout": "vertical",
			  "contents": [
				{
				  "type": "text",
				  "text": "รายการสินค้าคงเหลือ",
				  "color": "#ffffff",
				  "size": "xl",
				  "flex": 4,
				  "weight": "bold"
				},
				{
				  "type": "text",
				  "text": "ประจำวันที่ ` + currentTime + `",
				  "size": "sm",
				  "color": "#ffffff66"
				}
			  ]
			}
		  ],
		  "paddingAll": "20px",
		  "backgroundColor": "#0367D3",
		  "spacing": "md",
		  "height": "90px",
		  "paddingTop": "22px"
		},
		"body": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			` + inputFmt + `
		  ]
		},
		"footer": {
		  "type": "box",
		  "layout": "vertical",
		  "contents": [
			{
			  "type": "button",
			  "action": {
				"type": "uri",
				"label": "นำสินค้าออก",
				"uri": "https://liff.line.me/1656205141-1QNAezQL"
			  },
			  "style": "primary",
			  "color": "#0367D3",
			  "height": "sm"
			}
		  ]
		}
	  }`

	// Unmarshal JSON
	flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(jsonOutput))
	if err != nil {
		log.Println(err)
	}

	// New Flex Message
	return linebot.NewFlexMessage("เช็ควันหมดอายุของวันนี้", flexContainer)
}
func expiryBroadcast() {

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
					flexMessage := expiryDashboard()

					// Reply Message
					if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
						log.Print(err)
					}
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
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

func main() {

	// e := models.Product{
	// 	Name:    "Sam",
	// 	Barcode: "12312312312321312",
	// }
	// e.LeavesRemaining()

	router := gin.Default()

	// Set up CORS middleware options
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
			// return origin == "https://line-bot-ituyen.web.app"
		},
		MaxAge: 12 * time.Hour,
	}

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

	router.GET("/cronjob", cronJob)

	router.POST("/callback", lineBot)

	productAPI := router.Group("/product")
	{
		productAPI.GET("", func(c *gin.Context) {
			barcode := c.Query("barcode")
			input := map[string]string{
				"Barcode": barcode,
			}
			result := crud.List(input)
			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})
		productAPI.POST("/create", func(c *gin.Context) {
			var form models.ProductHasExpiry
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

	expiryAPI := router.Group("/expiry")
	{
		expiryAPI.GET("", func(c *gin.Context) {
			barcode := c.Query("barcode")
			input := map[string]string{
				"Barcode": barcode,
			}
			result := crud.ListExpiry(input)
			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})
		expiryAPI.POST("", func(c *gin.Context) {
			var form models.ProductHasExpiry
			if c.ShouldBind(&form) == nil {
				result := crud.AddExpiry(form)
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
		expiryAPI.PUT("", func(c *gin.Context) {
			var form models.ProductHasExpiry
			if c.ShouldBind(&form) == nil {
				result := crud.UpdateExpiry(form)
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
		expiryAPI.DELETE("", func(c *gin.Context) {
			var form models.ProductHasExpiry
			if c.ShouldBind(&form) == nil {
				fmt.Printf("Delete barcode %v %v \n", form.Barcode, form.Quantity)
				result := crud.RemoveExpiry(form)
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
	router.Use(cors.New(config))
	if port == "" {
		fmt.Println("Running on Heroku using random PORT")
		router.Run()
	} else {
		fmt.Println("Environment Port : " + port)
		router.Run(":" + port)
	}
}
