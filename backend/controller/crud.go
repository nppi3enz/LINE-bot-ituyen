package controller

import (
	"backend/models"
	"context"

	// initial "firebase/initFirebase"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	// "github.com/mitchellh/mapstructure"
)

var ctx = context.Background()
var client = Init(ctx)

func Init(ctx context.Context) *firestore.Client {
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
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

func AddData(p models.Product) {
	// ProductsData := new(models.Product)
	ProductsData := map[string]interface{}{
		"Name":    p.Name,
		"Barcode": p.Barcode,
	}
	_, _, err := client.Collection("products").Add(ctx, ProductsData)

	if err != nil {
		log.Fatalf("Failed adding product: %v", err)
	}

	// return c.JSON(http.StatusCreated, nil)

}

// func Destroy(c *gin.Context) {
// 	client.Collection("products")
// 	// .Where(c.Param("_id")).Delete(ctx)
// 	return c.JSON(http.StatusNoContent, nil)
// 	// _, _, err := client.Collection("income-v2").Add(ctx, IncomesData)
// }
