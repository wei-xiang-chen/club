package login

import (
	"club/models"
	"club/models/rest"
	"club/services/user_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) error {
	var restResult rest.RestResult
	var login models.User

	c.ShouldBindJSON(&login)

	user, err := user_service.Insert(&login)
	if err != nil {
		return nil
	}

	restResult.Data = user
	c.JSON(http.StatusOK, restResult)
	return nil
}
