package routes

import (
	"go_ecom_mongo/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("users/register", controllers.Register())
	incomingRoutes.POST("users/login", controllers.Login())

	// incomingRoutes.POST("users/address", controller.addAddress())
	// incomingRoutes.GET("users/address", controller.getAddress())

}
