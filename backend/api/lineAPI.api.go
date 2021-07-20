package api

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
)

var connectBot *linebot.Client
var client *firestore.Client
var ctx context.Context

func SetupLineCallback(router *gin.Engine, c *firestore.Client, cx context.Context) {
	connectBot = connectLineBot()
	client = c
	ctx = cx

	router.POST("/callback", lineBot)
}

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
	result := ListExpiry(make(map[string]string), client, ctx)
	var outputMsg string
	now := time.Now()

	for _, val := range result {
		dateExpiry := val.ExpireDate
		diff := dateExpiry.Sub(now).Hours()
		calculateDay := math.Ceil(diff / 24)

		nameProduct := val.Product["name"].(string)

		if val.Quantity > 1 {
			nameProduct = `x` + fmt.Sprint(val.Quantity) + ` ` + nameProduct
		}

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
				  "text": "` + nameProduct + `",
				  "flex": 3,
				  "weight": "bold",
				  "margin": "sm",
				  "wrap": true
				},
				{
				  "type": "text",
				  "text": "` + fmt.Sprint(calculateDay) + ` วัน",
				  "flex": 1,
				  "size": "sm",
				  "color": "` + color + `",
				  "align": "end"
				}
			  ],
			  "paddingTop": "5px",
			  "paddingBottom": "5px"
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
				if message.Text == "เช็ควันหมดอายุ" {
					flexMessage := expiryDashboard()

					// Reply Message
					if _, err = bot.ReplyMessage(event.ReplyToken, flexMessage).Do(); err != nil {
						log.Print(err)
					}
				}
				break
			}
		}
	}
}
