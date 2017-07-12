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
	total     int
}

func (cd *CategoryPriceData) Map() map[string]interface{} {
	return map[string]interface{}{
		"Id":        cd.id,
		"Min":       cd.min,
		"Suggested": cd.suggested,
		"Max":       cd.max,
		"Total":     cd.total,
	}
}

type CategoryService interface {
	Price(categoryId string) (data Data, err error)
}

type categoryMeli struct {
	client meliclient.Client
}

func New() *categoryMeli {
	return &categoryMeli{client: meliclient.New()}
}

func (c *categoryMeli) Price(categoryId string) (data Data, err error) {
	categoryPriceData := &CategoryPriceData{}
	categoryData, err := c.client.GetCategory(categoryId)
	if err != nil {
		return categoryPriceData, err
	}

	categoryPriceData.id = categoryData.Id
	categoryPriceData.total = categoryData.Total_items_in_this_category
	return categoryPriceData, nil
}
