package meliclient

import (
	"encoding/json"
	"net/http"
	"time"
)

type CategoryResponse struct {
	Id                           string
	Total_items_in_this_category int
}

type Client interface {
	GetCategory(categoryId string) (CategoryResponse, error)
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

func (m *meliClient) GetCategory(categoryId string) (CategoryResponse, error) {

	categoryResponse := CategoryResponse{}
	response, err := m.httpClient.Get("https://api.mercadolibre.com/categories/" + categoryId)
	if err != nil {
		return categoryResponse, err
	}
	defer response.Body.Close()

	errDecode := json.NewDecoder(response.Body).Decode(&categoryResponse)
	return categoryResponse, errDecode
}
