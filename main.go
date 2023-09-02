package main

import (
	"go_ecom_mongo/controllers"
	"go_ecom_mongo/db"
	"go_ecom_mongo/middleware"
	"go_ecom_mongo/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("Port")

	if port == "" {
		port = "3000"
	}

	controllers.NewApplication(db.ProductDetails(db.Client, "Products"), db.UserDeatils(db.Client, "Users"), db.PaymentDetails(db.Client, "Payments"))
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/getProfile", controllers.GetProfile())
	router.POST("/addProduct", controllers.AddProduct())
	router.POST("/editProduct", controllers.UpdateProduct())

	router.GET("/getProduct", controllers.GetAllProduct())
	router.DELETE("/deleteproduct", controllers.DeleteProductbyId())

	router.POST("/addCartProduct", controllers.AddCart())
	router.DELETE("/deleteCartProduct", controllers.DeleteCart())
	router.DELETE("/deleteAllCartProduct", controllers.DeleteAllInCart())
	router.GET("/getCartProduct", controllers.GetCart())

	router.POST("/createOrder", controllers.CreateOrder())

	router.POST("/refundPayment", controllers.GetPaymentDetailsFormRazorapy())

	// router.GET("/getProductbyid", controllers.GetProductbyId())

	// router.POST("/users/signup", controllers.Register())

	// log.Fatal(app)

	log.Fatal(router.Run(":" + port))

	// log.Fatalln(app)

	// app:=controllers.NewApplication()
	// app := contr.NewApplication(database.ProductDetailsData(database.Client, "Products"), database.UserData(database.Client, "Users"))

}
