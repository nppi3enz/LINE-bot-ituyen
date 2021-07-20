package api

import (
	"backend/models"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func SetupExpiryAPI(router *gin.Engine, client *firestore.Client, ctx context.Context) {

	expiryAPI := router.Group("/expiry")
	{
		expiryAPI.GET("", func(c *gin.Context) {
			barcode := c.Query("barcode")
			input := map[string]string{
				"Barcode": barcode,
			}
			result := ListExpiry(input, client, ctx)
			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})
		expiryAPI.POST("", func(c *gin.Context) {
			var form models.ProductHasExpiry
			if c.ShouldBind(&form) == nil {
				result := AddExpiry(form, client, ctx)
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
				result := UpdateExpiry(form, client, ctx)
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
				result := RemoveExpiry(form, client, ctx)
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
}

func ListExpiry(p map[string]string, client *firestore.Client, ctx context.Context) []models.Expiry {
	col := client.Collection("expiry")

	var query firestore.Query
	query = col.Query
	if p["Barcode"] != "" {
		query = col.Where("product.barcode", "==", p["Barcode"])
	} else {
		query = col.OrderBy("expireDate", firestore.Asc)
	}

	docs := query.Documents(ctx)

	var expiries []models.Expiry
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		// if err != nil {
		// 	return []
		// }
		var b models.Expiry
		fmt.Println(b)
		if err := doc.DataTo(&b); err != nil {
			// Handle error, possibly by returning the error
		}
		fmt.Println(b)
		expiries = append(expiries, b)
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
	}
	return expiries
}

func AddExpiry(p models.ProductHasExpiry, client *firestore.Client, ctx context.Context) error {
	result := client.Collection("expiry").Where("product.barcode", "==", p.Barcode).Documents(ctx)
	var docID string
	var docData map[string]interface{}

	for {
		doc, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
		docID = doc.Ref.ID
		docData = doc.Data()
	}
	if docData == nil {
		result = client.Collection("products").Where("barcode", "==", p.Barcode).Documents(ctx)
		for {
			doc, err := result.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
			docID = doc.Ref.ID
			docData = doc.Data()
		}

		fmt.Println(docData)
		expiredTime, err := time.Parse(time.RFC3339, p.ExpireDate+"T00:00:00.000+07:00")
		if err != nil {
			fmt.Println(err)
		}

		InitialData := map[string]interface{}{
			"expireDate": expiredTime,
			"quantity":   p.Quantity,
			"product": map[string]interface{}{
				"ID":      docID,
				"name":    docData["name"],
				"barcode": docData["barcode"],
			},
		}

		_, _, err = client.Collection("expiry").Add(ctx, InitialData)

		if err != nil {
			log.Fatalf("Failed adding expired: %v", err)
		}
	} else {
		return errors.New("Please Delete old expiry before add new")
	}
	return nil
}

func UpdateExpiry(p models.ProductHasExpiry, client *firestore.Client, ctx context.Context) error {
	result := client.Collection("expiry").Where("product.barcode", "==", p.Barcode).Documents(ctx)
	var docID string
	var docData map[string]interface{}

	for {
		doc, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
		docID = doc.Ref.ID
		docData = doc.Data()
	}
	if docData != nil {
		_, err := client.Collection("expiry").Doc(docID).Update(ctx, []firestore.Update{
			{
				Path:  "quantity",
				Value: p.Quantity,
			},
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
	} else {
		return errors.New("Barcode not Found")
	}
	return nil
}

func RemoveExpiry(p models.ProductHasExpiry, client *firestore.Client, ctx context.Context) error {
	result := client.Collection("expiry").Where("product.barcode", "==", p.Barcode).Documents(ctx)
	var docID string
	var docData map[string]interface{}

	for {
		doc, err := result.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
		docID = doc.Ref.ID
		docData = doc.Data()
	}
	if docData != nil {
		if docData["quantity"].(int64) <= 1 {
			// delete
			_, err := client.Collection("expiry").Doc(docID).Delete(ctx)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err)
			}
		} else {
			// minus 1 item
			_, err := client.Collection("expiry").Doc(docID).Update(ctx, []firestore.Update{
				{
					Path:  "quantity",
					Value: firestore.Increment(-1 * p.Quantity),
				},
			})
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err)
			}
			// fmt.Println("NOK")
		}
	} else {
		return errors.New("Barcode not Found")
	}

	return nil
}
