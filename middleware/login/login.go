package login

import (
	"club/pojo"
	"club/pojo/rest"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	restResult rest.RestResult
)

func Login(c *gin.Context) {
	var login pojo.Login

	c.ShouldBindJSON(&login)

	restResult.Data = loginService.Insert(&login)
	c.JSON(http.StatusOK, restResult)
}
