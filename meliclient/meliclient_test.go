package meliclient

import (
	"testing"
)

func TestWhenCountryIsValidCountryCodeGetterReturnsCountryCode(t *testing.T) {

	client := New()
	countryCode, err:= client.getCountryCode("MLA1234")

	if err != nil || countryCode != "MLA" {
		t.Fail()
	}

}

func TestWhenCountryIsInvalidCountryCodeGetterReturnsError(t *testing.T) {

	client := New()
	_, err:= client.getCountryCode("ML")

	if err == nil {
		t.Fail()
	}

	_, err = client.getCountryCode("MNN1234")

	if err == nil {
		t.Fail()
	}

}

func TestWhenGivenCategoryAndFirstPageReturnsARequestWithCategoryLimitAndOffset(t *testing.T) {

	client := New()
	req, _ := client.getCategoryItemsRequest("MLA1234", 0, 200)

	query := req.URL.Query()

	if query.Get("category") != "MLA1234" {
		t.Fail()
	}

	if query.Get("offset") != "0" {
		t.Fail()
	}

	if query.Get("limit") != "199" {
		t.Fail()
	}

	if req.URL.Path != "/sites/MLA/search" {
		t.Fail()
	}

}