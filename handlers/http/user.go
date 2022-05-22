package api

import (
	"finance/car-finance/back-end/entities"
	"finance/car-finance/back-end/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAPIHandler struct {
	userRepo repositories.UserRepository
}

func (h *UserAPIHandler) GetUser(c *gin.Context) {

	h.userRepo.GetUser(c)
	c.JSON(401, map[string]interface{}{
		"id":     "go.micro.client",
		"code":   401,
		"detail": "Unauthorized",
		"status": "Unauthorized",
	})
}

func (h *UserAPIHandler) GetUsers(c *gin.Context) {
	var RequestUser *entities.RequestUser
	// var res *entities.UserResponse
	if err := c.ShouldBindJSON(&RequestUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if models, paging, err := h.userRepo.GetUsers(c, RequestUser); err == nil {
		// res.Message = "success"
		Data := make([]*entities.UserDetail, 0)
		for _, v := range models {
			Data = append(Data, v.ToProtobuf())
		}
		meta := paging.ToProtobuf()
		// res.Meta = meta
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success",
			"data":    Data,
			"meta":    meta,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
		// return microError.InternalServerError("500", "%s", "Internal Server Error")
	}

}

func (h *UserAPIHandler) Create(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result, err := h.userRepo.Create(c, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, &result)

}
