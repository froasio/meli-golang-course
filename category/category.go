package category

import (
	"meli-golang-course/meliclient"
)

type Data interface {
	Map() map[string]interface{}
}

type CategoryPriceData struct {
	id        string
	min       float64
	suggested float64
	max       float64
	total     uint
	pages     uint
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

func (c *categoryMeli) getTotalPages(totalItems uint, pageSize uint) uint {

	return (totalItems + pageSize - 1) / pageSize

}

func (c *categoryMeli) Price(categoryId string) (data Data, err error) {

	categoryPriceData := &CategoryPriceData{}
	categoryData, err := c.client.GetCategory(categoryId)

	if err != nil {
		return categoryPriceData, err
	}

	totalPages := c.getTotalPages(categoryData.TotalItems, c.pageSize)

	categoryPriceData.id = categoryData.Id
	categoryPriceData.total = categoryData.TotalItems
	categoryPriceData.pages = totalPages
	return categoryPriceData, nil

}
