package main

import (
	"errors"
	"github.com/froasio/meli-golang-course/category"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"sync"
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

func (c *categoryMeliMock) PriceEstimation(categoryId string) (data category.Data, err error) {

	if categoryId == "MLA1234" {
		return &categoryDataMock{1.0, 5.0, 10.0}, nil
	}

	return &categoryDataMock{}, errors.New("Invalid Category")
}

func getRouter(categoryService category.CategoryService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	setCategoryPricerRoute(router, categoryService)
	return router
}

func TestPricingAPIResponses(t *testing.T) {
	Convey("Given a router", t, func(c C) {
		router := getRouter(&categoryMeliMock{})
		Convey("When the category is valid it should return Ok", func() {
			r, _ := http.NewRequest("GET", "/categories/MLA1234/prices", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			c.So(w.Code, ShouldEqual, http.StatusOK)
		})
		Convey("When the category is invalid it should return Bad Request", func() {
			r, _ := http.NewRequest("GET", "/categories/MLA1231/prices", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			c.So(w.Code, ShouldEqual, http.StatusBadRequest)
		})
		Convey("When router receive many concurrent calls it should respond to all of them", func() {

			var wg sync.WaitGroup

			requestBatch := [5]struct {
				request                *http.Request
				expectedStatusResponse int
			}{}

			requestBatch[0].request, _ = http.NewRequest("GET", "/categories/MLA1234/prices", nil)
			requestBatch[0].expectedStatusResponse = http.StatusOK
			requestBatch[1].request, _ = http.NewRequest("GET", "/categories/MLA1234/prices", nil)
			requestBatch[1].expectedStatusResponse = http.StatusOK
			requestBatch[2].request, _ = http.NewRequest("GET", "/categories/MLA1234/prices", nil)
			requestBatch[2].expectedStatusResponse = http.StatusOK
			requestBatch[3].request, _ = http.NewRequest("GET", "/categories/MLA1231/prices", nil)
			requestBatch[3].expectedStatusResponse = http.StatusBadRequest
			requestBatch[4].request, _ = http.NewRequest("GET", "/categories/MLA1231/prices", nil)
			requestBatch[4].expectedStatusResponse = http.StatusBadRequest

			for i := 0; i < len(requestBatch); i++ {

				wg.Add(1)
				go func(i int) {
					writer := httptest.NewRecorder()
					defer wg.Done()
					router.ServeHTTP(writer, requestBatch[i].request)
					c.So(writer.Code, ShouldEqual, requestBatch[i].expectedStatusResponse)
				}(i)

			}
			wg.Wait()
		})
	})
}

func BenchmarkPricingCalculation(b *testing.B) {

	router := getRouter(category.New())
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/categories/MLA1377/prices", nil)
	router.ServeHTTP(w, r)

}
