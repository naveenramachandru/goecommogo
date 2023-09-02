package controllers

import (
	"context"

	"go_ecom_mongo/db"
	"go_ecom_mongo/models"
	generate "go_ecom_mongo/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

var emptyarr = make([]int, 0)

var UserCollection *mongo.Collection = db.UserDeatils(db.Client, "Users")

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection

	paymentCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection, paymentCollection *mongo.Collection) *Application {

	return &Application{
		prodCollection:    prodCollection,
		userCollection:    userCollection,
		paymentCollection: paymentCollection,
	}

}

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

			return
		}

		//check the fields (valiadator)

		validatorErr := Validate.Struct(user)
		if validatorErr != nil {

			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "Validation Error", "data": emptyarr})

			return

		}

		//check the user is present

		countemail, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

			return

		}

		if countemail > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "User email already exits", "data": emptyarr})
			return

		}

		countPhoneNo, err := UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "somethiong went worng", "data": emptyarr})

			return

		}

		if countPhoneNo > 0 {

			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "User phone number already exits", "data": emptyarr})
			return

		}

		//bycrypt the password
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()

		//generate token and refresh token

		token, refreshToken, err := generate.TokenGenrator(*user.Email, *user.Password, user.ID.String())

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "unable to generate token", "data": err})

			return
		}
		if token == "" || refreshToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "unable to generate token", "data": err})

			return

		}

		user.Token = &token
		user.Refresh_Token = &refreshToken

		//insert the details in db

		_, insertErr := UserCollection.InsertOne(ctx, user)

		if insertErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"s": false, "message": "Error to add this data", "data": insertErr})

			return

		}
		defer cancel()

		//final output

		c.JSON(http.StatusOK, gin.H{"data": user, "s": true, "message": "data found"})

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var user models.User
		var founduser models.User

		if err := c.BindJSON(&user); err != nil {

			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"s": false, "message": err, "data": emptyarr})
			return

		}

		//check the user is present

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)

		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"s": false, "message": "user not found", "data": emptyarr})
			return

		}

		//verify password

		CheckPasswordIsValid, msg := VerifyPassword(*founduser.Password, *user.Password)

		if !CheckPasswordIsValid {

			c.JSON(http.StatusInternalServerError, gin.H{"s": false, "message": msg, "data": emptyarr})

			return

		}

		token, refresh, err := generate.TokenGenrator(*user.Email, *user.Password, user.ID.String())

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"s": false, "message": err, "data": emptyarr})

			return

		}

		if token == "" || refresh == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"s": false, "message": "some thing went wrong", "data": emptyarr})

			return

		}

		UpdateAllTokens(token, refresh, founduser.User_ID)

		//get user details

		c.JSON(http.StatusOK, gin.H{"s": false, "message": "data found", "data": founduser})

	}
}

func HashPassword(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)

}

func VerifyPassword(dbpassword string, enteredPassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(enteredPassword))

	valid := true
	msg := ""

	if err != nil {
		msg = "login/password is incorrect "
		valid = false
	}

	return valid, msg

}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateobj primitive.D
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserCollection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	},
		&opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
