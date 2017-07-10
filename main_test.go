package main

import (
	"errors"
	"github.com/gin-gonic/gin"
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

func TestWhenServiceCanCalculateCategoryPricingTheApiRespondsOk(t *testing.T) {
	router, w := getRouter(&categoryMeliMock{})
	r, _ := http.NewRequest("GET", "/categories/MLA1234/prices", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Fail()
	}
}

func TestWhenServiceCannotCalculateCategoryPricingTheApiRespondsBadRequest(t *testing.T) {

	gin.SetMode(gin.TestMode)
	router := gin.New()
	setCategoryPricerRoute(router, &categoryMeliMock{})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/categories/MLA1231/prices", nil)
	router.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fail()
	}
}
