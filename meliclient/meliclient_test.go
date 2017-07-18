package meliclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhenCountryIsValidCountryCodeGetterReturnsCountryCode(t *testing.T) {

	client := New()
	countryCode, err := client.getCountryCode("MLA1234")

	if err != nil || countryCode != "MLA" {
		t.Fail()
	}

}

func TestWhenCountryIsInvalidCountryCodeGetterReturnsError(t *testing.T) {

	client := New()
	_, err := client.getCountryCode("ML")

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

	if query.Get("limit") != "200" {
		t.Fail()
	}

	if req.URL.Path != "/sites/MLA/search" {
		t.Fail()
	}

}

func getMockedCategoryResponse() string {
	return `{"id":"MLA10272","name":"Registradoras Antiguas","picture":"http://resources.mlstatic.com/category/images/c8ce50fb-c995-4fb5-9c3c-61dd143e8ff4.png","permalink":null,"total_items_in_this_category":581,"path_from_root":[{"id":"MLA1367","name":"Antigüedades"},{"id":"MLA10272","name":"Registradoras Antiguas"}],"children_categories":[],"attribute_types":"none","settings":{"adult_content":false,"buying_allowed":true,"buying_modes":["auction","buy_it_now"],"catalog_domain":null,"coverage_areas":"not_allowed","currencies":["ARS"],"fragile":false,"immediate_payment":"required","item_conditions":["used","not_specified","new"],"items_reviews_allowed":false,"listing_allowed":true,"max_description_length":50000,"max_pictures_per_item":12,"max_pictures_per_item_var":6,"max_sub_title_length":70,"max_title_length":60,"maximum_price":null,"minimum_price":null,"mirror_category":null,"mirror_master_category":null,"mirror_slave_categories":[],"price":"required","reservation_allowed":"not_allowed","restrictions":[],"rounded_address":false,"seller_contact":"not_allowed","shipping_modes":["not_specified","custom"],"shipping_options":["custom","carrier"],"shipping_profile":"optional","show_contact_information":false,"simple_shipping":"optional","stock":"required","sub_vertical":"other","subscribable":false,"tags":[],"vertical":"other","vip_subdomain":"articulo"},"meta_categ_id":null,"attributable":false}`
}

func TestWhenGettingTheCategoryItGetsParsedIntoStruct(t *testing.T) {

	getCategoryHandler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(getMockedCategoryResponse()))
	}

	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/categories/", getCategoryHandler)
	defer ts.Close()

	client := &meliClient{httpClient: &http.Client{}, baseUrl: ts.URL}
	categoryResponse, _ := client.GetCategory("MLA10272")

	if categoryResponse.Id != "MLA10272" {
		t.Fail()
	}

	if categoryResponse.TotalItems != 581 {
		t.Fail()
	}
}

func TestWhenGettingAnInvalidCategoryItReturnsAnError(t *testing.T) {

	getCategoryHandler := func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(""))
	}

	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/categories/", getCategoryHandler)
	defer ts.Close()

	client := &meliClient{httpClient: &http.Client{}, baseUrl: ts.URL}
	_, err := client.GetCategory("MLA10272")

	if err == nil {
		t.Fail()
	}

}
