package domain

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PageinationOptions struct {
	Page int
	Size int
}

func GetOptions(c *gin.Context) PageinationOptions {
	pageString := c.DefaultQuery("page", "1")
	sizeString := c.DefaultQuery("size", "100")

	var pageinationOptions PageinationOptions

	page, err := strconv.Atoi(pageString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "page number has to be a whole number")
	}

	size, err := strconv.Atoi(sizeString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "page size has to be a whole number")
	}
	if size > 250 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "page size has to be between 1 and 250")
	}

	pageinationOptions.Page = page
	pageinationOptions.Size = size

	return pageinationOptions

}
