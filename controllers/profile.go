package controllers

import (
	"context"
	"go_ecom_mongo/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// Assuming 'email' is stored in the claims as you are trying to access it
		emailClaim, exists := c.Get("email")

		if !exists {
			log.Println(emailClaim)
			log.Println(exists)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Something went wrong", "data": emptyarr})
			return
		}

		// Convert the email claim to a string
		email, ok := emailClaim.(string)
		if !ok {
			log.Println(email)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid email claim", "data": emptyarr})
			return
		}

		err1 := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

		if err1 != nil {
			log.Println(err1)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user not found", "data": emptyarr})
			return

		}

		c.JSON(http.StatusOK, gin.H{"success": true, "message": "user found", "data": user})
	}
}

// QEJM2c9l5RNOxdEh
// 122.172.87.33/32
