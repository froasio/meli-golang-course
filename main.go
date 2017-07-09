package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.GET("/categories/:category/prices", func(c *gin.Context) {

		category := c.Param("category")
		c.JSON(http.StatusOK, gin.H{"category": category})

	})
	router.Run(":8080")
}
