package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	pageDefault  int = 1
	limitDefault int = 10
)

func GetPagination(c *gin.Context) (*int, *int, error) {
	var page, limit *int

	if value, ok := c.GetQuery("page"); ok {
		p, err := strconv.Atoi(value)
		if err != nil {
			return nil, nil, err
		}
		page = &p
	} else {
		page = &pageDefault
	}

	if value, ok := c.GetQuery("limit"); ok {
		l, err := strconv.Atoi(value)
		if err != nil {
			return nil, nil, err
		}
		limit = &l
	} else {
		limit = &limitDefault
	}

	offset := (*page - 1) * *limit

	return &offset, limit, nil
}
