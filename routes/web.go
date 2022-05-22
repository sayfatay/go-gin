package routes

import (
	httpHandler "finance/car-finance/back-end/handlers/http"

	"github.com/gin-gonic/gin"
)

// WebRouter web router
func SetupRouter(r *gin.Engine) *gin.Engine {
	// r := gin.Default()

	userAPI := new(httpHandler.UserAPIHandler)
	r.Any("/", func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			"code":    200,
			"message": "success",
		})
	})

	r.GET("/user/:id", userAPI.GetUser)
	r.POST("/user", userAPI.GetUsers)
	r.POST("/user/create", userAPI.Create)
	// r.PUT("/customers/:id", h.UpdateCustomer)
	// r.DELETE("/customers/:id", h.DeleteCustomer)

	return r
}
