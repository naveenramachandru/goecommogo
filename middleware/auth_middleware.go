package middleware

import (
	"fmt"
	// "go_ecom_mongo/controllers"
	"log"
	"net/http"

	generate "github.com/naveenramachandru/goecommogo/utils"

	"github.com/gin-gonic/gin"
	"github.com/naveenramachandru/goecommogo/controllers"
	"go.mongodb.org/mongo-driver/bson"
)

var emptyarr = make([]int, 0)

// Authz validates token and authorizes users
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": fmt.Sprintf("No Authorization header provided"), "data": emptyarr})
			c.Abort()
			return
		}

		claims, err := generate.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": false, "error": err, "data": emptyarr})
			c.Abort()
			return
		}

		tokendata, err1 := controllers.GetUserID(c, claims.Email)

		if err1 != "" {
			log.Println(err)

			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "user not found", "data": emptyarr})
			c.Abort()

			return

		}

		log.Println(bson.M{"email": claims.Email})
		log.Println(bson.M{"userid": tokendata})

		c.Set("email", claims.Email)
		c.Set("userid", tokendata)

		c.Next()

	}
}
