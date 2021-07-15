package main

import (
	crud "backend/controller"
	"backend/models"
	"fmt"
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
					// Unmarshal JSON
					flexContainer, err := linebot.UnmarshalFlexMessageJSON([]byte(`{
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
								  "text": "ประจำวันที่ 15/07/2021",
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
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/restaurant_regular_32.png"
								},
								{
								  "type": "text",
								  "text": "รายการสินค้าชิ้นที่ 1",
								  "flex": 0,
								  "weight": "bold",
								  "margin": "sm"
								},
								{
								  "type": "text",
								  "text": "1 Days",
								  "size": "sm",
								  "color": "#CD113B",
								  "align": "end"
								}
							  ],
							  "height": "30px"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/restaurant_regular_32.png"
								},
								{
								  "type": "text",
								  "text": "รายการสินค้าชิ้นที่ 2",
								  "flex": 0,
								  "weight": "bold",
								  "margin": "sm"
								},
								{
								  "type": "text",
								  "text": "3 Days",
								  "size": "sm",
								  "color": "#FF7600",
								  "align": "end"
								}
							  ],
							  "height": "30px"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "icon",
								  "url": "https://scdn.line-apps.com/n/channel_devcenter/img/fx/restaurant_regular_32.png"
								},
								{
								  "type": "text",
								  "text": "รายการสินค้าชิ้นที่ 3",
								  "flex": 0,
								  "weight": "bold",
								  "margin": "sm"
								},
								{
								  "type": "text",
								  "text": "5 Days",
								  "size": "sm",
								  "color": "#52006A",
								  "align": "end"
								}
							  ],
							  "height": "30px"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "text",
								  "text": "รายการสินค้าชิ้นที่ 4",
								  "flex": 0,
								  "weight": "bold",
								  "margin": "sm"
								},
								{
								  "type": "text",
								  "text": "10 Days",
								  "size": "sm",
								  "color": "#aaaaaa",
								  "align": "end"
								}
							  ],
							  "height": "30px"
							},
							{
							  "type": "box",
							  "layout": "baseline",
							  "contents": [
								{
								  "type": "text",
								  "text": "รายการสินค้าชิ้นที่ 5",
								  "flex": 0,
								  "weight": "bold",
								  "margin": "sm"
								},
								{
								  "type": "text",
								  "text": "30 Days",
								  "size": "sm",
								  "color": "#aaaaaa",
								  "align": "end"
								}
							  ],
							  "height": "30px"
							}
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
					  }`))
					if err != nil {
						log.Println(err)
					}

					// New Flex Message
					flexMessage := linebot.NewFlexMessage("เช็ควันหมดอายุของวันนี้", flexContainer)

					// Reply Message
					if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
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
