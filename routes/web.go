package routes

import (
	httpHandler "finance/car-finance/back-end/handlers/http"

	"github.com/gin-gonic/gin"
)

// WebRouter web router
// func WebRouter(serviceName string, router *gin.Engine) *gin.Engine {
// 	r := router.Group(fmt.Sprintf("/%s", serviceName)) // do not change
// 	{
// 		userAPI := new(httpHandler.UserAPIHandler)
// 		r.Any("/", func(c *gin.Context) {
// 			c.JSON(200, map[string]interface{}{
// 				"code":    200,
// 				"message": "success",
// 			})
// 		})
// 		r.POST("/user/get", userAPI.GetUser)
// 		// r.GET("/cdn/get-file", cdnAPI.GetFile)
// 		// r.POST("/cdn/validate-file-member", cdnAPI.ValidateFileImportFileMember)
// 		// r.POST("/cdn/import-file-member", cdnAPI.ImportFileMember)
// 	}
// 	return router
// }
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
