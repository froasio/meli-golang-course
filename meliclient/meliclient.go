package meliclient

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type CategoryResponse struct {
	Id         string `json:"id"`
	TotalItems uint   `json:"total_items_in_this_category"`
}

type Item struct {
	Price float64
}

type CategoryItemsResponse struct {
	Results []Item
}

type Client interface {
	GetCategory(categoryId string) (CategoryResponse, error)
	GetCategoryItems(cat string, page uint, pageSize uint) (CategoryItemsResponse, error)
}

type meliClient struct {
	httpClient *http.Client
}

func New() *meliClient {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	return &meliClient{httpClient: client}
}

func (m *meliClient) getCountryCode(cat string) (string, error) {

	if len(cat) < 3 {
		return cat, errors.New("Invalid category")
	}
	validCountryCode := map[string]bool{
		"MLA": true,
		"MBO": true,
		"MLB": true,
		"MLC": true,
		"MCO": true,
		"MCR": true,
		"MRD": true,
		"MEC": true,
		"MHN": true,
		"MGT": true,
		"MLM": true,
		"MNI": true,
		"MPY": true,
		"MPA": true,
		"MPE": true,
		"MPT": true,
		"MSV": true,
		"MLU": true,
		"MLV": true,
	}
	countryCode := cat[:3]
	if _, ok := validCountryCode[countryCode]; ok {
		return countryCode, nil
	}
	return countryCode, errors.New("Invalid category")

}

func (m *meliClient) GetCategory(categoryId string) (CategoryResponse, error) {

	categoryResponse := CategoryResponse{}
	response, err := m.httpClient.Get("https://api.mercadolibre.com/categories/" + categoryId)
	if err != nil {
		return categoryResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return categoryResponse, errors.New("Invalid category")
	}

	errDecode := json.NewDecoder(response.Body).Decode(&categoryResponse)
	return categoryResponse, errDecode
}

func (m *meliClient) getCategoryItemsRequest(cat string, page uint, pageSize uint) (*http.Request, error) {

	countryCode, err := m.getCountryCode(cat)

	if err != nil {
		return nil, err
	}

	offset := strconv.Itoa(int(page * pageSize))
	limit := strconv.Itoa(int((page+1)*pageSize - 1))

	req, err := http.NewRequest("GET", "https://api.mercadolibre.com/sites/"+countryCode+"/search", nil)
	if err != nil {
		return req, err
	}

	q := req.URL.Query()
	q.Add("category", cat)
	q.Add("limit", limit)
	q.Add("offset", offset)
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func (m *meliClient) GetCategoryItems(cat string, page uint, pageSize uint) (CategoryItemsResponse, error) {

	categoryItemsResponse := CategoryItemsResponse{}

	request, err := m.getCategoryItemsRequest(cat, page, pageSize)
	if err != nil {
		return categoryItemsResponse, err
	}

	response, err := m.httpClient.Do(request)
	if err != nil {
		return categoryItemsResponse, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return categoryItemsResponse, errors.New("Invalid category items")
	}

	errDecoding := json.NewDecoder(response.Body).Decode(&categoryItemsResponse)
	return categoryItemsResponse, errDecoding
}
