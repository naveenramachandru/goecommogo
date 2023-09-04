package controllers

import (
	"context"
	"log"
	"time"

	"github.com/naveenramachandru/goecommogo/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserID(c *gin.Context, email string) (userId string, err string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user1 models.User

	errval := UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user1)
	if errval != nil {
		log.Println("Error fetching user data:", errval)
		return "", "Error fetching user data"
	}

	return user1.ID.Hex(), ""
}
