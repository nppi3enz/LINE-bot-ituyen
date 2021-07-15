package controller

import (
	"backend/models"
	"context"
	"errors"
	"fmt"
	"time"

	// initial "firebase/initFirebase"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	// "github.com/mitchellh/mapstructure"
)

var ctx = context.Background()
var client = Init(ctx)

func Init(ctx context.Context) *firestore.Client {
	sa := option.WithCredentialsFile("google-credentials.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	fmt.Println("run client : ")
	fmt.Println(client)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

// func Home(c *gin.Context) {
// 	ProductsData := []models.Product{}
// 	// client := initial.Init(ctx)
// 	iter := client.Collection("products").Documents(ctx)
// 	for {
// 		ProductData := models.Product{}
// 		_, err := iter.Next()
// 		if err == iterator.Done {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}

// 		// mapstructure.Decode(doc.Data(), &IncomeData)
// 		ProductsData = append(ProductsData, ProductData)
// 	}
// 	return c.JSON(http.StatusOK, ProductsData)
// }
func List(p map[string]string) []models.Product {
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
		// if err != nil {
		// 	return []
		// }
		var b models.Product
		if err := doc.DataTo(&b); err != nil {
			// Handle error, possibly by returning the error
		}
		fmt.Println(b)
		products = append(products, b)
		// fmt.Printf("Value = %s: %s", doc.Ref.ID, doc.Data())
	}
	return products
}
func AddData(p models.ProductHasExpiry) error {

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

// func Destroy(c *gin.Context) {
// 	client.Collection("products")
// 	// .Where(c.Param("_id")).Delete(ctx)
// 	return c.JSON(http.StatusNoContent, nil)
// 	// _, _, err := client.Collection("income-v2").Add(ctx, IncomesData)
// }

func AddExpiry(p models.ProductHasExpiry) error {
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

func UpdateExpiry(p models.ProductHasExpiry) error {
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

func RemoveExpiry(p models.ProductHasExpiry) error {
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
