package main

import (
	"github.com/froasio/meli-golang-course/category"
	"github.com/gin-gonic/gin"
	"net/http"
)

func setCategoryPricerRoute(router *gin.Engine, categoryService category.CategoryService) {

	router.GET("/categories/:category/prices", func(c *gin.Context) {

		categoryId := c.Param("category")

		data, err := categoryService.Price(categoryId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, data.Map())
		}
	})

}

func main() {

	router := gin.Default()
	setCategoryPricerRoute(router, category.New())
	router.Run(":8080")

}
