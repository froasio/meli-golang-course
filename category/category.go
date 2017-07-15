package category

import (
	"meli-golang-course/meliclient"
)

type Data interface {
	Map() map[string]interface{}
}

type CategoryPriceData struct {
	id          string
	min         float64
	suggested   float64
	max         float64
	total       uint
	pages       uint
	cummulative float64
}

func (cd *CategoryPriceData) Map() map[string]interface{} {
	return map[string]interface{}{
		"Id":        cd.id,
		"Min":       cd.min,
		"Suggested": cd.suggested,
		"Max":       cd.max,
		"Total":     cd.total,
		"Pages":     cd.pages,
	}
}

type CategoryService interface {
	Price(categoryId string) (data Data, err error)
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

func (c *categoryMeli) getCategoryPricingByPage(categoryId string, page uint) (CategoryPriceData, error) {

	categoryItems, err := c.client.GetCategoryItems(categoryId, page, c.pageSize)

	if err != nil {
		return CategoryPriceData{}, err
	}

	totalResults := uint(len(categoryItems.Results))
	if totalResults == 0 {
		return CategoryPriceData{}, nil
	}

	categoryPriceData := CategoryPriceData{
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

	return categoryPriceData, nil
}

func (c *categoryMeli) Price(categoryId string) (data Data, err error) {

	categoryPriceData := &CategoryPriceData{}
	categoryData, err := c.client.GetCategory(categoryId)

	if err != nil {
		return categoryPriceData, err
	}

	totalPages := c.getTotalPages(categoryData.TotalItems)

	categoryPriceData.id = categoryData.Id
	categoryPriceData.total = categoryData.TotalItems
	categoryPriceData.pages = totalPages
	return categoryPriceData, nil

}
