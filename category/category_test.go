package category

import (
	"fmt"
	"meli-golang-course/meliclient"
	"testing"
)

var totalPagesTest = []struct {
	total    uint
	expected uint
}{
	{0, 0},
	{199, 1},
	{200, 1},
	{250, 2},
}

type Client interface {
	GetCategory(categoryId string) (meliclient.CategoryResponse, error)
	GetCategoryItems(cat string, page uint, pageSize uint) (meliclient.CategoryItemsResponse, error)
}

type meliClientMock struct {
	categoryResponse           meliclient.CategoryResponse
	categoryResponseError      error
	categoryItemsResponse      meliclient.CategoryItemsResponse
	categoryItemsResponseError error
}

func (m *meliClientMock) GetCategory(categoryId string) (meliclient.CategoryResponse, error) {
	return m.categoryResponse, m.categoryResponseError
}

func (m *meliClientMock) GetCategoryItems(cat string, page uint, pageSize uint) (meliclient.CategoryItemsResponse, error) {
	return m.categoryItemsResponse, m.categoryItemsResponseError
}

func TestTotalPages(t *testing.T) {

	cat := New()

	for _, tt := range totalPagesTest {
		actual := cat.getTotalPages(tt.total)
		if actual != tt.expected {
			t.Fail()
		}
	}

}

func TestGetCategoryItems(t *testing.T) {

	max := 10.0
	min := 1.0

	response := meliclient.CategoryItemsResponse{
		Results: []meliclient.Item{
			meliclient.Item{Price: 1.0},
			meliclient.Item{Price: 2.0},
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
	cummulative := 55.0

	meliclient := &meliClientMock{
		categoryResponse:           meliclient.CategoryResponse{},
		categoryResponseError:      nil,
		categoryItemsResponse:      response,
		categoryItemsResponseError: nil,
	}

	cat := &categoryMeli{client: meliclient, pageSize: 200}
	data, _ := cat.getCategoryPricingByPage("MLA1234", 0)

	if data.max != max {
		t.Fail()
	}

	if data.min != min {
		t.Fail()
	}

	if data.cummulative != cummulative {
		t.Fail()
	}

	if data.total != total {
		t.Fail()
	}

}
