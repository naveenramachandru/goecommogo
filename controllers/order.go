package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/naveenramachandru/goecommogo/db"
	"github.com/naveenramachandru/goecommogo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var OrderCollection *mongo.Collection = db.OrderDeatils(db.Client, "Orders")

type OrderData struct {
	models.Payment `json:"payment"`
	UserOrder      struct {
		OrderList   []models.Order `json:"order_list"`
		OrderStatus string         `json:"order_status"`
		TotalPrice  int            `json:"total_price"`
		Discount    int            `json:"discount"`
	} `json:"order"`
}

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		userid, exists := c.Get("userid")

		if !exists {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email claim", "data": emptyarr})
			return

		}

		userId, ok := userid.(string)
		if !ok {

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email claim", "data": emptyarr})
			return
		}

		userCartCount, carterr := CartCollection.CountDocuments(ctx, bson.M{"userid": userId})

		if carterr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": carterr, "data": emptyarr})

			return

		}

		if userCartCount < 1 {

			c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "item not found in cart", "data": emptyarr})

			return
		}
		byteArr, err := io.ReadAll(c.Request.Body)

		if err != nil {
			log.Println(userId)

			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var orderdata OrderData
		if err := json.Unmarshal(byteArr, &orderdata); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		if orderdata.UserOrder.Discount < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "discount value is improper", "data": emptyarr})

			return

		}

		ispayemnt, payment_id, inserterr := CreatePayment(c, ctx, orderdata.Payment)

		if !ispayemnt {
			c.JSON(http.StatusBadRequest, inserterr)

			return

		}
		if payment_id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to update payment details", "data": emptyarr})

			return

		}

		var i int

		for i = 0; i < len(orderdata.UserOrder.OrderList); i++ {

			var Ordermodel models.Order

			Ordermodel = orderdata.UserOrder.OrderList[i]

			Ordermodel.OrderID = primitive.NewObjectID()
			Ordermodel.PaymentID = payment_id
			Ordermodel.OrdereredAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			Ordermodel.Discount = Ordermodel.Discount
			Ordermodel.OrderStatus = orderdata.UserOrder.OrderStatus

			log.Println(Ordermodel)

			insterOrder, orderinserterr := OrderCollection.InsertOne(ctx, Ordermodel)

			if orderinserterr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "unable to update  details", "data": emptyarr})
				continue

				return
			}

			removeData, removeerr := CartCollection.DeleteOne(ctx, bson.M{"userid": userId, "product_id": Ordermodel.ProductID})

			if removeerr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": removeerr, "data": emptyarr})

				continue

				return

			}

			log.Println(insterOrder)
			log.Println(removeData)

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "data found", "data": emptyarr})

	}

}

func UpdateOrder() {

}
