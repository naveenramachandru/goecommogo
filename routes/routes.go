package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/naveenramachandru/goecommogo/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	incomingRoutes.POST("users/register", controllers.Register())
	incomingRoutes.POST("users/login", controllers.Login())

	// incomingRoutes.POST("users/address", controller.addAddress())
	// incomingRoutes.GET("users/address", controller.getAddress())

}
