package category

import (
	"errors"
)

type Data interface {
	Map() map[string]interface{}
}

type CategoryData struct {
	min       float64
	suggested float64
	max       float64
}

func (cd *CategoryData) Map() map[string]interface{} {
	return map[string]interface{}{
		"Min":       cd.min,
		"Suggested": cd.suggested,
		"Max":       cd.max,
	}
}

type CategoryService interface {
	Price(categoryId string) (data Data, err error)
}

type CategoryMeli struct {
}

func (c *CategoryMeli) Price(categoryId string) (data Data, err error) {

	if categoryId == "MLA1234" {
		return &CategoryData{1.0, 5.0, 10.0}, nil
	}

	return &CategoryData{}, errors.New("Invalid Category")
}
