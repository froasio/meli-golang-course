package category

import (
	"github.com/froasio/meli-golang-course/meliclient"
	"math"
	"sync"
)

type Data interface {
	Map() map[string]interface{}
}

type CategoryPriceData struct {
	min         float64
	suggested   float64
	max         float64
	total       uint
	cummulative float64
}

func (cd *CategoryPriceData) Map() map[string]interface{} {
	return map[string]interface{}{
		"Min":       cd.min,
		"Suggested": cd.suggested,
		"Max":       cd.max,
	}
}

type CategoryService interface {
	PriceEstimation(categoryId string) (data Data, err error)
}

type categoryMeli struct {
	client   meliclient.Client
	pageSize uint
}

func New() *categoryMeli {

	return &categoryMeli{client: meliclient.New(), pageSize: 200}

}

func (c *categoryMeli) getTotalPages(totalItems uint) uint {

	return (totalItems + c.pageSize - 1) / c.pageSize

}

func (c *categoryMeli) getCategoryPricingByPage(categoryId string, page uint) *CategoryPriceData {

	categoryItems, err := c.client.GetCategoryItems(categoryId, page, c.pageSize)

	if err != nil {
		return nil
	}

	totalResults := uint(len(categoryItems.Results))
	if totalResults == 0 {
		return nil
	}

	categoryPriceData := &CategoryPriceData{
		min:         categoryItems.Results[0].Price,
		max:         categoryItems.Results[0].Price,
		total:       totalResults,
		cummulative: categoryItems.Results[0].Price,
	}

	for i := uint(1); i < totalResults; i++ {
		if categoryItems.Results[i].Price > categoryPriceData.max {
			categoryPriceData.max = categoryItems.Results[i].Price
		}
		if categoryItems.Results[i].Price < categoryPriceData.min {
			categoryPriceData.min = categoryItems.Results[i].Price
		}
		categoryPriceData.cummulative += categoryItems.Results[i].Price
	}

	return categoryPriceData
}

func (c *categoryMeli) reduceCategoryPricingPages(pagesData []*CategoryPriceData) *CategoryPriceData {

	categoryPriceData := &CategoryPriceData{
		min:         math.Inf(1),
		max:         math.Inf(-1),
		total:       0,
		cummulative: 0,
	}

	for _, data := range pagesData {
		if data != nil {
			if data.min < categoryPriceData.min {
				categoryPriceData.min = data.min
			}
			if data.max > categoryPriceData.max {
				categoryPriceData.max = data.max
			}
			categoryPriceData.total += data.total
			categoryPriceData.cummulative += data.cummulative
		}
	}

	if categoryPriceData.total != 0 {
		categoryPriceData.suggested = categoryPriceData.cummulative / float64(categoryPriceData.total)
	}

	return categoryPriceData

}

func (c *categoryMeli) PriceEstimation(categoryId string) (data Data, err error) {

	categoryData, err := c.client.GetCategory(categoryId)
	var wg sync.WaitGroup

	if err != nil {
		return &CategoryPriceData{}, err
	}

	totalPages := c.getTotalPages(categoryData.TotalItems)
	pagesData := make([]*CategoryPriceData, totalPages, totalPages)

	for i := uint(0); i < totalPages; i++ {

		wg.Add(1)
		go func(page uint) {

			defer wg.Done()
			pagesData[page] = c.getCategoryPricingByPage(categoryId, page)

		}(i)

	}
	wg.Wait()

	categoryPriceData := c.reduceCategoryPricingPages(pagesData)

	return categoryPriceData, nil

}
