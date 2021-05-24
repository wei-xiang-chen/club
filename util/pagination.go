package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PAGE_DEFAULT  int = 1
	LIMIT_DEFAULT int = 10
)

func GetPagination(c *gin.Context) (*int, *int, error) {
	var page, limit int

	if value, ok := c.GetQuery("page"); ok {
		p, err := strconv.Atoi(value)
		if err != nil {
			return nil, nil, err
		}
		page = p
	} else {
		page = PAGE_DEFAULT
	}

	if value, ok := c.GetQuery("limit"); ok {
		l, err := strconv.Atoi(value)
		if err != nil {
			return nil, nil, err
		}
		limit = l
	} else {
		limit = LIMIT_DEFAULT
	}

	offset := (page - 1) * limit

	return &offset, &limit, nil
}
