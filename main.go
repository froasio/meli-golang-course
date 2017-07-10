package main

import (
	"github.com/gin-gonic/gin"
	"meli-golang-course/category"
	"net/http"
)

func CategoryPricerHandler(categoryService category.CategoryService) gin.HandlerFunc {

	fn := func(c *gin.Context) {

		categoryId := c.Param("category")

		data, err := categoryService.Price(categoryId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, data.Map())
		}
	}

	return gin.HandlerFunc(fn)

}

func setCategoryPricerRoute(router *gin.Engine, categoryService category.CategoryService) {

	router.GET("/categories/:category/prices", CategoryPricerHandler(categoryService))

}

func main() {

	router := gin.Default()

	setCategoryPricerRoute(router, &category.CategoryMeli{})

	router.Run(":8080")

}
