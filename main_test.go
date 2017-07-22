package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"meli-golang-course/category"
	"net/http"
	"net/http/httptest"
	"testing"
)

type categoryMeliMock struct {
}

type categoryDataMock struct {
	min       float64
	suggested float64
	max       float64
}

func (cd *categoryDataMock) Map() map[string]interface{} {
	return map[string]interface{}{
		"Min":       cd.min,
		"Suggested": cd.suggested,
		"Max":       cd.max,
	}
}

func (c *categoryMeliMock) Price(categoryId string) (data category.Data, err error) {

	if categoryId == "MLA1234" {
		return &categoryDataMock{1.0, 5.0, 10.0}, nil
	}

	return &categoryDataMock{}, errors.New("Invalid Category")
}

func getRouter(categoryService category.CategoryService) (*gin.Engine, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setCategoryPricerRoute(router, categoryService)
	w := httptest.NewRecorder()
	return router, w
}

func TestPricingAPIResponses(t *testing.T) {
	Convey("Given a router", t, func() {
		router, w := getRouter(&categoryMeliMock{})
		Convey("When the category is valid it should return Ok", func() {
			r, _ := http.NewRequest("GET", "/categories/MLA1234/prices", nil)
			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("When the category is invalid it should return Bad Request", func() {
			r, _ := http.NewRequest("GET", "/categories/MLA1231/prices", nil)
			router.ServeHTTP(w, r)
			So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

func BenchmarkPricingCalculation(b *testing.B) {

	router, w := getRouter(category.New())
	r, _ := http.NewRequest("GET", "/categories/MLA1377/prices", nil)
	router.ServeHTTP(w, r)

}
