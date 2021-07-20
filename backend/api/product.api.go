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

func SetupProductAPI(router *gin.Engine, client *firestore.Client, ctx context.Context) {
	productAPI := router.Group("/product")
	{
		productAPI.GET("", func(c *gin.Context) {
			barcode := c.Query("barcode")
			input := map[string]string{
				"Barcode": barcode,
			}
			result := ListProduct(input, client, ctx)
			c.JSON(http.StatusOK, gin.H{
				"data": result,
			})
		})
		productAPI.POST("/create", func(c *gin.Context) {
			var form models.ProductHasExpiry
			if c.ShouldBind(&form) == nil {
				result := CreateProduct(form, client, ctx)
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
}
func ListProduct(p map[string]string, client *firestore.Client, ctx context.Context) []models.Product {
	col := client.Collection("products")

	var query firestore.Query
	query = col.Query
	if p["Barcode"] != "" {
		query = col.Where("barcode", "==", p["Barcode"])
	}

	docs := query.Documents(ctx)

	var products []models.Product
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
		// if err != nil {
		// 	return []
		// }
		var b models.Product
		if err := doc.DataTo(&b); err != nil {
			// Handle error, possibly by returning the error
			fmt.Println("Error!")
		}

		products = append(products, b)
	}
	return products
}
func CreateProduct(p models.ProductHasExpiry, client *firestore.Client, ctx context.Context) error {

	result := client.Collection("products").Where("barcode", "==", p.Barcode).Documents(ctx)

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
		// docID = doc.Ref.ID
		docData = doc.Data()
	}
	if docData == nil {
		ProductsData := map[string]interface{}{
			"name":    p.Name,
			"barcode": p.Barcode,
		}
		doc, _, err := client.Collection("products").Add(ctx, ProductsData)

		if err != nil {
			log.Fatalf("Failed adding product: %v", err)
		}
		expiredTime, err := time.Parse(time.RFC3339, p.ExpireDate+"T00:00:00.000+07:00")
		if err != nil {
			fmt.Println(err)
		}

		InitialData := map[string]interface{}{
			"expireDate": expiredTime,
			"quantity":   p.Quantity,
			"product": map[string]interface{}{
				"ID":      doc.ID,
				"name":    p.Name,
				"barcode": p.Barcode,
			},
		}

		doc, _, err = client.Collection("expiry").Add(ctx, InitialData)

		if err != nil {
			log.Fatalf("Failed adding expired: %v", err)
		}
	} else {
		return errors.New("Already Add Barcode")
	}

	// return c.JSON(http.StatusCreated, nil)
	return nil
}
