package integration_test

import (
	"fmt"
	"meli-golang-course/category"
	"testing"
)

func TestWhenCountryIsValidCountryCodeGetterReturnsCountryCode(t *testing.T) {

	cat := category.New()
	data, _ := cat.Price("MLA1377")
	fmt.Println(data.Map())
}
