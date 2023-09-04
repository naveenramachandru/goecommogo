package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/naveenramachandru/goecommogo/models"

	"github.com/naveenramachandru/goecommogo/db"

	"github.com/gin-gonic/gin"
	"github.com/naveenramachandru/goecommogo/controllers/data"
	razorpay "github.com/razorpay/razorpay-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var PaymentCoolection *mongo.Collection = db.PaymentDetails(db.Client, "Payment")

func DataSmaple2() string {
	return ""
}

func CreatePayment(c *gin.Context, ctx context.Context, paymentData models.Payment) (dataInserted bool, payment_id string, inserterr gin.H) {

	// var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// defer cancel()
	_ = data.DataSmaple()

	paymentData.PaymentID = primitive.NewObjectID()

	if paymentData.UserID == "" {

		data := gin.H{"s": false, "message": "user id is required", "data": emptyarr}

		return false, "", data

	}
	if paymentData.TransactionID == "" {

		data := gin.H{"s": false, "message": "transaction id is required", "data": emptyarr}

		return false, "", data

	}

	checktransCount, checkerr := PaymentCoolection.CountDocuments(ctx, bson.M{"tansaction_id": paymentData.TransactionID})

	if checkerr != nil {
		data := gin.H{"s": false, "message": checkerr, "data": emptyarr}

		return false, "", data

	}

	if checktransCount > 0 {

		data := gin.H{"s": false, "message": "tansaction_id should be unique", "data": emptyarr}

		return false, "", data

	}

	paymentData.PaymentID = primitive.NewObjectID()

	paymentData.PaymentAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	log.Println(paymentData)

	createpaydata, err := PaymentCoolection.InsertOne(ctx, paymentData)

	if err != nil {

		return false, "", gin.H{"s": false, "message": err, "data": emptyarr}

	}

	log.Println(createpaydata)

	return true, paymentData.PaymentID.String(), gin.H{}

}

func GetPaymentDetailsFormRazorapy() gin.HandlerFunc {

	return func(c *gin.Context) {

		client := razorpay.NewClient("rzp_test_GGL6e9Zyd8xRNy", "HUTwntR8kVn9nNnJ6AOzBWEg")

		// queryParams := map[string]interface{}{} // You can provide query parameters here
		// headers := map[string]string{}

		// body, err := client.Payment.Fetch("pay_Lcx1QSbjHMkSa9", queryParams, headers)
		data := map[string]interface{}{
			"speed":   "optimum",
			"receipt": "Receipt No. 31",
		}
		body, err := client.Payment.Refund("pay_Lbq5jOc8g0JA77", 10*100, data, nil)

		if err != nil {

			log.Println(err)

		}

		c.JSON(http.StatusOK, gin.H{"success": false, "message": err, "data": body})

	}

}

func UpdatePayment() {

}
