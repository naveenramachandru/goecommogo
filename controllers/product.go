package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/naveenramachandru/goecommogo/db"
	"github.com/naveenramachandru/goecommogo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProductCollection *mongo.Collection = db.UserDeatils(db.Client, "Products")

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

			return
		}
		if *product.Product_Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Product name is required", "data": nil})
			return
		}
		if *product.Price == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Product Price is required", "data": nil})
			return
		}

		product.Product_ID = primitive.NewObjectID()

		userid, exists := c.Get("userid")

		if !exists {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong", "data": emptyarr})
			return
		}

		// Convert the email claim to a string
		userId, ok := userid.(string)
		if !ok {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email claim", "data": emptyarr})
			return
		}
		product.VednorID = &userId

		log.Println(product)
		log.Println(ctx)

		insertProduct, err := ProductCollection.InsertOne(ctx, product)

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "added sucesssfully", "data": insertProduct})

		//

	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

			return
		}
		if *product.Product_Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Product name is required", "data": nil})
			return
		}
		if *product.Price == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Product Price is required", "data": nil})
			return
		}

		userid, exists := c.Get("userid")

		if !exists {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong", "data": emptyarr})
			return
		}

		// Convert the email claim to a string
		userId, ok := userid.(string)

		if !ok {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email claim", "data": emptyarr})
			return
		}
		log.Println("user data us")

		log.Println(userId)
		log.Println(product.Product_ID.String())

		// productID, err := primitive.ObjectIDFromHex(product.Product_ID.String())

		// productdata := ProductCollection.FindOne(ctx, bson.M{"vednorid": userId, "_id": product.Product_ID}).Decode(&foundproduct)
		productcount, err := ProductCollection.CountDocuments(ctx, bson.M{"vednorid": userId, "_id": product.Product_ID})

		log.Println(productcount)
		log.Println(err)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error counting product", "data": emptyarr})
			return
		}

		if productcount == 0 {

			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "No Product Found", "data": emptyarr})
			return

		}

		updatepro, err := ProductCollection.UpdateOne(ctx, bson.M{"vednorid": userId, "_id": product.Product_ID}, bson.M{"$set": bson.M{"product_name": *product.Product_Name, "price": *product.Price, "rating": *product.Rating, "image": *product.Image, "quantity": *product.Quantity, "isactive": *product.IsActive, "avaliablestock": *product.AvaliableStock}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Error updating product", "data": emptyarr})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Product updated successfully", "data": updatepro})

	}
}

func GetAllProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var products []models.Product

		ProdutIdparms, ok := c.GetQuery("product_id")

		produtnameparms, ok1 := c.GetQuery("product_name")
		productIsActive, ok2 := c.GetQuery("active")

		searchQuery := c.DefaultQuery("product_name", "")

		log.Println(produtnameparms)

		var parmaId primitive.ObjectID

		if ok == true {

			id, err := primitive.ObjectIDFromHex(ProdutIdparms)
			if err != nil {
				log.Println("Error converting ID:", err)
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid ID format", "data": emptyarr})
				return
			}

			parmaId = id

		}

		page := 1
		perPage := 10 // Number of products per page

		// Parse query parameters for pagination
		if pageStr := c.DefaultQuery("page", c.GetHeader("page")); pageStr != "" {
			page, _ = strconv.Atoi(pageStr)
		}
		if perPageStr := c.DefaultQuery("per_page", c.GetHeader("per_page")); perPageStr != "" {
			perPage, _ = strconv.Atoi(perPageStr)
		}

		// Calculate skip and limit values for pagination
		skip := (page - 1) * perPage
		limit := perPage

		basonMap := make(map[string]interface{})

		if ok == true {

			basonMap = bson.M{"_id": parmaId}

		} else if ok1 == true {

			if searchQuery != "" {

				// basonMap["product_name"] = primitive.Regex{Pattern: searchQuery, Options: "i"}
				basonMap = bson.M{
					"product_name": primitive.Regex{Pattern: searchQuery, Options: "i"},
				}

				// Include the active status if provided
				if productIsActive != "" {
					boolValue, err := strconv.ParseBool(productIsActive)
					if err != nil {

						c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
						return
					}
					basonMap["isactive"] = boolValue
				}

			}

			// basonMap = bson.M{"product_name": produtnameparms}

		} else if ok2 == true {
			boolValue, err := strconv.ParseBool(productIsActive)
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
				return
			}

			basonMap = bson.M{"isactive": boolValue}

		} else {
			basonMap = bson.M{}

		}

		log.Println("map is ")
		log.Println(basonMap)

		getProduct, err := ProductCollection.Find(ctx, basonMap, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
			return

		}
		defer getProduct.Close(ctx)

		for getProduct.Next(ctx) {

			var pro models.Product

			if err := getProduct.Decode(&pro); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})

				return
			}

			products = append(products, pro)

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "data found", "data": products})

	}
}

func DeleteProductbyId() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		proid, err := c.GetQuery("product_id")

		if !err {
			log.Println(err)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "data not found", "data": emptyarr})
			return

		}

		log.Println(proid)

		// var pordidhex primitive.ObjectID

		id, err1 := primitive.ObjectIDFromHex(proid)

		if err1 != nil {

			log.Println(err1)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "data not found", "data": emptyarr})
			return

		}

		data, msg := ProductCollection.DeleteOne(ctx, bson.M{"_id": id})
		if msg != nil {
			log.Println(msg)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": msg, "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "data  found", "data": data})

	}
}

func GetProductForVendor() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
