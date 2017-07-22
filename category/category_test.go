package category

import (
	"github.com/froasio/meli-golang-course/meliclient"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
)

type Client interface {
	GetCategory(categoryId string) (*meliclient.CategoryResponse, error)
	GetCategoryItems(cat string, page uint, pageSize uint) (*meliclient.CategoryItemsResponse, error)
}

type meliClientMock struct {
	categoryResponse           meliclient.CategoryResponse
	categoryResponseError      error
	categoryItemsResponse      meliclient.CategoryItemsResponse
	categoryItemsResponseError error
}

func (m *meliClientMock) GetCategory(categoryId string) (*meliclient.CategoryResponse, error) {
	return &(m.categoryResponse), m.categoryResponseError
}

func (m *meliClientMock) GetCategoryItems(cat string, page uint, pageSize uint) (*meliclient.CategoryItemsResponse, error) {
	return &(m.categoryItemsResponse), m.categoryItemsResponseError
}

func TestCategoryPriceDataMapping(t *testing.T) {

	Convey("Given a category price data struct it should transform to map", t, func() {
		cat := &CategoryPriceData{min: 1.0, max: 10.0, suggested: 5.0}
		catMap := cat.Map()
		expectedMap := map[string]interface{}{
			"Min":       1.0,
			"Max":       10.0,
			"Suggested": 5.0,
		}
		So(catMap, ShouldResemble, expectedMap)
	})

}

func TestTotalPages(t *testing.T) {
	Convey("Given an amount of items", t, func() {

		cat := New()
		totalPagesTest := []struct {
			total    uint
			expected uint
		}{
			{0, 0},
			{199, 1},
			{200, 1},
			{250, 2},
		}

		for _, tt := range totalPagesTest {
			Convey("When total items is "+strconv.Itoa(int(tt.total)), func() {
				So(cat.getTotalPages(tt.total), ShouldEqual, tt.expected)
			})
		}

	})
}

func TestGetCategoryItemsPricingCalculation(t *testing.T) {

	Convey("Given a category service", t, func() {

		cat := &categoryMeli{client: nil, pageSize: 200}

		Convey("When response is not empty is should return pricint parameters", func() {
			max := 10.0
			min := 1.0
			cummulative := 55.0

			response := meliclient.CategoryItemsResponse{
				Results: []meliclient.Item{
					meliclient.Item{Price: 2.0},
					meliclient.Item{Price: 1.0},
					meliclient.Item{Price: 3.0},
					meliclient.Item{Price: 4.0},
					meliclient.Item{Price: 5.0},
					meliclient.Item{Price: 6.0},
					meliclient.Item{Price: 7.0},
					meliclient.Item{Price: 8.0},
					meliclient.Item{Price: 9.0},
					meliclient.Item{Price: 10.0},
				},
			}
			total := uint(len(response.Results))

			cat.client = &meliClientMock{
				categoryResponse:           meliclient.CategoryResponse{},
				categoryResponseError:      nil,
				categoryItemsResponse:      response,
				categoryItemsResponseError: nil,
			}

			data := cat.getCategoryPricingByPage("MLA1234", 0)
			expectedData := &CategoryPriceData{
				min:         min,
				max:         max,
				cummulative: cummulative,
				total:       total,
			}
			So(data, ShouldResemble, expectedData)
		})
		Convey("When response is empty pricing total should be 0", func() {
			response := meliclient.CategoryItemsResponse{
				Results: []meliclient.Item{},
			}

			cat.client = &meliClientMock{
				categoryResponse:           meliclient.CategoryResponse{},
				categoryResponseError:      nil,
				categoryItemsResponse:      response,
				categoryItemsResponseError: nil,
			}

			data := cat.getCategoryPricingByPage("MLA1234", 0)
			So(data, ShouldBeNil)
		})
	})

}

func TestReducingPagesResults(t *testing.T) {

	Convey("Given an array of pages results", t, func() {
		cat := New()
		Convey("When array is not empty it should reduce the pricing result", func() {
			pricingPages := []*CategoryPriceData{
				&CategoryPriceData{
					min:         2.0,
					max:         10.0,
					cummulative: 100.0,
					total:       10,
				},
				nil,
				&CategoryPriceData{
					min:         1.0,
					max:         11.0,
					cummulative: 100.0,
					total:       10,
				},
			}

			data := cat.reduceCategoryPricingPages(pricingPages)
			expectedData := &CategoryPriceData{
				max:         11.0,
				min:         1.0,
				suggested:   10.0,
				total:       20,
				cummulative: 200,
			}
			So(data, ShouldResemble, expectedData)
		})
	})
}

func TestCalculatingItemsPricing(t *testing.T) {

	max := 10.0
	min := 1.0
	Convey("Given a category", t, func() {
		categoryId := "MLA1234"
		Convey("When the response is no empty it should return the pricing estimation", func() {
			response := meliclient.CategoryItemsResponse{
				Results: []meliclient.Item{
					meliclient.Item{Price: 2.0},
					meliclient.Item{Price: 1.0},
					meliclient.Item{Price: 3.0},
					meliclient.Item{Price: 4.0},
					meliclient.Item{Price: 5.0},
					meliclient.Item{Price: 6.0},
					meliclient.Item{Price: 7.0},
					meliclient.Item{Price: 8.0},
					meliclient.Item{Price: 9.0},
					meliclient.Item{Price: 10.0},
				},
			}

			meliclient := &meliClientMock{
				categoryResponse:           meliclient.CategoryResponse{TotalItems: 10},
				categoryResponseError:      nil,
				categoryItemsResponse:      response,
				categoryItemsResponseError: nil,
			}

			cat := &categoryMeli{client: meliclient, pageSize: 200}
			data, _ := cat.Price(categoryId)
			dataMapping := data.Map()
			expectedMap := map[string]interface{}{
				"Min":       min,
				"Max":       max,
				"Suggested": 5.5,
			}
			So(dataMapping, ShouldResemble, expectedMap)
		})
	})

}
