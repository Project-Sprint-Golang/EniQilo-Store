package routes

import (
	"net/http"

	controller "github.com/Project-Sprint-Golang/EniQilo-Store/app/controllers"
	middleware "github.com/Project-Sprint-Golang/EniQilo-Store/app/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/login", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World!")
		})
		v1.POST("/staff/register", controller.RegisterStaff)
		v1.POST("/customer/register", controller.RegisterCustomer)
		v1.POST("/staff/login", controller.Login)

		v1.Use(middleware.AuthMiddleware())
		v1.POST("/product", controller.AddProduct)
		v1.GET("/product", controller.GetAllProduct)
		v1.PUT("/product/:id", controller.UpdateProduct)
		v1.DELETE("/product/:id", controller.DeleteProduct)
		v1.GET("/product/customer", controller.GetSKUProduct)

	}

}
