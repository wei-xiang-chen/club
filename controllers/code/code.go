package code

import (
	"club/models/rest"
	"club/services/code_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Code(c *gin.Context) error {
	var restResult rest.RestResult

	types := c.QueryArray("type")

	codes, err := code_service.Code(types)
	if err != nil {
		return err
	}

	restResult.Data = codes
	c.JSON(http.StatusOK, restResult)
	return nil
}
