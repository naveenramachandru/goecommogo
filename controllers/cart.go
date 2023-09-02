package controllers

import (
	"context"
	"go_ecom_mongo/db"
	"go_ecom_mongo/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CartCollection *mongo.Collection = db.UserDeatils(db.Client, "Cart")

func AddCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var cartproduct models.ProductUser

		if err := c.BindJSON(&cartproduct); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

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

		productcount, errmsg := CartCollection.CountDocuments(ctx, bson.M{"product_id": *cartproduct.Product_ID, "userid": userId})

		if errmsg != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": errmsg, "data": emptyarr})

			return

		}
		log.Println(productcount)

		if productcount > 0 {

			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "item is already in cart", "data": emptyarr})

			return

		}

		cartproduct.ID = primitive.NewObjectID()

		cartproduct.UserID = &userId

		addproducttocart, err := CartCollection.InsertOne(ctx, cartproduct)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "item added to cart", "data": addproducttocart})

	}
}

func DeleteCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		cartid, err := c.GetQuery("cart_id")

		if !err {
			log.Println(err)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "data not found", "data": emptyarr})
			return

		}

		log.Println(cartid)

		// var pordidhex primitive.ObjectID

		id, err1 := primitive.ObjectIDFromHex(cartid)

		if err1 != nil {

			log.Println(err1)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "data not found", "data": emptyarr})
			return

		}

		data, msg := CartCollection.DeleteOne(ctx, bson.M{"_id": id})
		if msg != nil {
			log.Println(msg)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": msg, "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "data  found", "data": data})

	}
}

func DeleteAllInCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

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

		deleteallcart, err := CartCollection.DeleteMany(ctx, bson.M{"userid": userId})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "deleted all data", "data": deleteallcart})

	}
}

func GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var cartdata []models.ProductUser

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

		userCartData, err := CartCollection.Find(ctx, bson.M{"userid": userId})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})
			return

		}

		defer userCartData.Close(ctx)

		for userCartData.Next(ctx) {

			var cartmodel models.ProductUser

			if err := userCartData.Decode(&cartmodel); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err, "data": emptyarr})

				return
			}

			cartdata = append(cartdata, cartmodel)

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "data found", "data": cartdata})

	}
}
