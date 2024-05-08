package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/login", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World!")
		})
	}

}
