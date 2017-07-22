package meliclient

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWhenCountryCodeValidation(t *testing.T) {

	Convey("Given a client", t, func() {
		client := New()
		Convey("When category country is valid it should return country code with no errors", func() {
			countryCode, err := client.getCountryCode("MLA1234")
			So(err, ShouldBeNil)
			So(countryCode, ShouldEqual, "MLA")
		})
		Convey("When country code is invalid it should return an error", func() {
			Convey("When length is less than 3", func() {
				_, err := client.getCountryCode("ML")
				So(err, ShouldNotBeNil)
			})
			Convey("When country code doesn't exists", func() {
				_, err := client.getCountryCode("MNN1234")
				So(err, ShouldNotBeNil)
			})
		})
	})

}

func TestCategoryItemsRequestPaging(t *testing.T) {

	Convey("Given a client", t, func() {
		client := New()
		Convey("When category is valid, page is 0 and page size is 200 it should return the request path with valid offset and limit", func() {
			req, _ := client.getCategoryItemsRequest("MLA1234", 0, 200)
			query := req.URL.Query()
			So(query.Get("category"), ShouldEqual, "MLA1234")
			So(query.Get("offset"), ShouldEqual, "0")
			So(query.Get("limit"), ShouldEqual, "200")
			So(req.URL.Path, ShouldEqual, "/sites/MLA/search")
		})
	})

}

func getMockedCategoryResponse() string {
	return `{"id":"MLA10272","name":"Registradoras Antiguas","picture":"http://resources.mlstatic.com/category/images/c8ce50fb-c995-4fb5-9c3c-61dd143e8ff4.png","permalink":null,"total_items_in_this_category":581,"path_from_root":[{"id":"MLA1367","name":"Antigüedades"},{"id":"MLA10272","name":"Registradoras Antiguas"}],"children_categories":[],"attribute_types":"none","settings":{"adult_content":false,"buying_allowed":true,"buying_modes":["auction","buy_it_now"],"catalog_domain":null,"coverage_areas":"not_allowed","currencies":["ARS"],"fragile":false,"immediate_payment":"required","item_conditions":["used","not_specified","new"],"items_reviews_allowed":false,"listing_allowed":true,"max_description_length":50000,"max_pictures_per_item":12,"max_pictures_per_item_var":6,"max_sub_title_length":70,"max_title_length":60,"maximum_price":null,"minimum_price":null,"mirror_category":null,"mirror_master_category":null,"mirror_slave_categories":[],"price":"required","reservation_allowed":"not_allowed","restrictions":[],"rounded_address":false,"seller_contact":"not_allowed","shipping_modes":["not_specified","custom"],"shipping_options":["custom","carrier"],"shipping_profile":"optional","show_contact_information":false,"simple_shipping":"optional","stock":"required","sub_vertical":"other","subscribable":false,"tags":[],"vertical":"other","vip_subdomain":"articulo"},"meta_categ_id":null,"attributable":false}`
}

func TestCategoryRequestResponseParsing(t *testing.T) {

	Convey("Given a client", t, func() {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		client := &meliClient{httpClient: &http.Client{}, baseUrl: ts.URL}
		Convey("When category request response is valid it should be parsed into struct", func() {
			getCategoryHandler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(getMockedCategoryResponse()))
			}
			mux.HandleFunc("/categories/", getCategoryHandler)
			categoryResponse, _ := client.GetCategory("MLA10272")
			expectedCategorResponse := &CategoryResponse{Id: "MLA10272", TotalItems: 581}
			So(categoryResponse, ShouldResemble, expectedCategorResponse)
		})
		Convey("When category request response is invalid it should return an error", func() {
			getCategoryHandler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(""))
			}
			mux.HandleFunc("/categories/", getCategoryHandler)
			_, err := client.GetCategory("MLA10272")
			So(err, ShouldNotBeNil)
		})
	})

}

func getMockedCategoryItemsResponse() string {
	return `{"site_id":"MLA","paging":{"total":403187,"offset":0,"limit":2},"results":[{"id":"MLA645327713","site_id":"MLA","title":"Silla Royal Cristal  Tiffany Almohadon Capitone Con Cristal","seller":{"id":231727494,"power_seller_status":null,"car_dealer":false,"real_estate_agency":false,"tags":[]},"price":1650,"currency_id":"ARS","available_quantity":3,"sold_quantity":92,"buying_mode":"buy_it_now","listing_type_id":"gold_special","stop_time":"2036-11-30T01:10:32.000Z","condition":"new","permalink":"http://articulo.mercadolibre.com.ar/MLA-645327713-silla-royal-cristal-tiffany-almohadon-capitone-con-cristal-_JM","thumbnail":"http://mla-s2-p.mlstatic.com/131915-MLA25342980142_022017-I.jpg","accepts_mercadopago":true,"installments":null,"address":{"state_id":"AR-C","state_name":"Capital Federal","city_id":"TUxBQlZJTDQyMjBa","city_name":"Villa Crespo"},"shipping":{"free_shipping":false,"mode":"not_specified","tags":[]},"seller_address":{"id":199863531,"comment":"","address_line":"","zip_code":"","country":{"id":"AR","name":"Argentina"},"state":{"id":"AR-C","name":"Capital Federal"},"city":{"id":"TUxBQlZJTDQyMjBa","name":"Villa Crespo"},"latitude":-34.5984082,"longitude":-58.4303291},"attributes":[],"original_price":null,"category_id":"MLA11794","official_store_id":null,"catalog_product_id":null,"reviews":{"rating_average":5,"total":1}},{"id":"MLA612515596","site_id":"MLA","title":"La Tapera - Estilo Campo / Asador Parrilla Con Cruz Y Disco","seller":{"id":7767036,"power_seller_status":null,"car_dealer":false,"real_estate_agency":false,"tags":[]},"price":11000,"currency_id":"ARS","available_quantity":1,"sold_quantity":12,"buying_mode":"buy_it_now","listing_type_id":"gold_special","stop_time":"2036-03-14T18:39:18.000Z","condition":"new","permalink":"http://articulo.mercadolibre.com.ar/MLA-612515596-la-tapera-estilo-campo-asador-parrilla-con-cruz-y-disco-_JM","thumbnail":"http://mla-s2-p.mlstatic.com/1095-MLA4732518522_072013-I.jpg","accepts_mercadopago":true,"installments":null,"address":{"state_id":"AR-B","state_name":"Buenos Aires","city_id":"","city_name":"Tortuguitas"},"shipping":{"free_shipping":false,"mode":"not_specified","tags":[]},"seller_address":{"id":196531114,"comment":"","address_line":"","zip_code":"","country":{"id":"AR","name":"Argentina"},"state":{"id":"AR-B","name":"Buenos Aires"},"city":{"id":"","name":"Tortuguitas"},"latitude":-34.4902115,"longitude":-58.7766647},"attributes":[],"original_price":null,"category_id":"MLA35923","official_store_id":null,"catalog_product_id":null,"reviews":{"rating_average":4,"total":1}}],"secondary_results":[],"related_results":[],"sort":{"id":"relevance","name":"More relevant"},"available_sorts":[{"id":"price_asc","name":"Lower price"},{"id":"price_desc","name":"Higher price"}],"filters":[{"id":"category","name":"Categories","type":"text","values":[{"id":"MLA1367","name":"Antigüedades","path_from_root":[{"id":"MLA1367","name":"Antigüedades"}]}]}],"available_filters":[{"id":"category","name":"Categories","type":"text","values":[{"id":"MLA5467","name":"Decoración Antigua","results":2403},{"id":"MLA3257","name":"Muebles Antiguos","results":4783},{"id":"MLA7184","name":"Juguetes Antiguos","results":1814},{"id":"MLA3635","name":"Libros Antiguos","results":1545},{"id":"MLA4688","name":"Iluminación Antigua","results":1312},{"id":"MLA4631","name":"Vajilla Antigua","results":1080},{"id":"MLA1375","name":"Joyas y Relojes Antiguos","results":843},{"id":"MLA6652","name":"Audio Antiguo","results":722},{"id":"MLA1379","name":"Platería Antigua","results":646},{"id":"MLA6661","name":"Electrodomésticos Antiguos","results":570},{"id":"MLA4630","name":"Cristalería Antigua","results":496},{"id":"MLA6650","name":"Carteles Antiguos","results":466},{"id":"MLA1381","name":"Indumentaria Antigua","results":420},{"id":"MLA11224","name":"Herramientas Antiguas","results":374},{"id":"MLA6662","name":"Teléfonos Antiguos","results":182},{"id":"MLA34367","name":"Rejas y Portones Antiguos","results":175},{"id":"MLA34242","name":"Sulkys y Carros Antiguos","results":150},{"id":"MLA10081","name":"Balanzas Antiguas","results":148},{"id":"MLA10236","name":"Sifones Antiguos","results":144},{"id":"MLA1378","name":"Cámaras Fotográficas Antiguas","results":141},{"id":"MLA10232","name":"Máquinas de Escribir Antiguas","results":106},{"id":"MLA388312","name":"Máquinas de Coser","results":104},{"id":"MLA1373","name":"Artículos Marítimos Antiguos","results":97},{"id":"MLA10250","name":"Bolsos y Valijas Antiguas","results":92},{"id":"MLA1377","name":"Equipos Científicos Antiguos","results":69},{"id":"MLA1376","name":"Instrumentos Musicales Antig.","results":58},{"id":"MLA10272","name":"Registradoras Antiguas","results":53},{"id":"MLA10249","name":"Llaves y Candados Antiguos","results":45},{"id":"MLA34203","name":"Ropa de Cama Antigua","results":8},{"id":"MLA1383","name":"Otras Antigüedades","results":954}]},{"id":"official_store","name":"Tiendas oficiales","type":"text","values":[{"id":"all","name":"Todas las tiendas oficiales","results":94},{"id":"781","name":"Cebra","results":1}]},{"id":"state","name":"Location","type":"text","values":[{"id":"TUxBUENBUGw3M2E1","name":"Capital Federal","results":8581},{"id":"TUxBUEdSQXJlMDNm","name":"Bs.As. G.B.A. Sur","results":2466},{"id":"TUxBUEdSQWU4ZDkz","name":"Bs.As. G.B.A. Norte","results":2453},{"id":"TUxBUEdSQWVmNTVm","name":"Bs.As. G.B.A. Oeste","results":1817},{"id":"TUxBUFpPTmFpbnRl","name":"Buenos Aires Interior","results":1422},{"id":"TUxBUFNBTmU5Nzk2","name":"Santa Fe","results":799},{"id":"TUxBUENPU2ExMmFkMw","name":"Bs.As. Costa Atlántica","results":585},{"id":"TUxBUENPUmFkZGIw","name":"Córdoba","results":415},{"id":"TUxBUEVOVHMzNTdm","name":"Entre Ríos","results":189},{"id":"TUxBUExBWmE1OWMy","name":"La Pampa","results":88},{"id":"TUxBUE1FTmE5OWQ4","name":"Mendoza","results":82},{"id":"TUxBUE5FVW4xMzMzNQ","name":"Neuquén","results":41},{"id":"TUxBUFNBTGFjMTJi","name":"Salta","results":31},{"id":"TUxBUE1JU3MzNjIx","name":"Misiones","results":26},{"id":"TUxBUFLNT29iZmZm","name":"Río Negro","results":26},{"id":"TUxBUFNBTno3ZmY5","name":"Santa Cruz","results":19},{"id":"TUxBUENPUnM5MjI0","name":"Corrientes","results":14},{"id":"TUxBUFRVQ244NmM3","name":"Tucumán","results":12},{"id":"TUxBUENIVXQxNDM1MQ","name":"Chubut","results":11},{"id":"TUxBUENBVGFiY2Fm","name":"Catamarca","results":10},{"id":"TUxBUENIQW8xMTNhOA","name":"Chaco","results":10},{"id":"TUxBUFNBTnM0ZTcz","name":"San Luis","results":10},{"id":"TUxBUEpVSnk3YmUz","name":"Jujuy","results":7},{"id":"TUxBUFNBTm5lYjU4","name":"San Juan","results":7},{"id":"TUxBUFNBTm9lOTlk","name":"Santiago del Estero","results":6},{"id":"TUxBUEZPUmE1OTk5","name":"Formosa","results":5},{"id":"TUxBUExBWmEyNzY0","name":"La Rioja","results":2},{"id":"TUxBUFRJRVoxM2M5YQ","name":"Tierra del Fuego","results":1}]},{"id":"price","name":"Precio","type":"range","values":[{"id":"*-500.0","name":"Up to $500","results":6524},{"id":"500.0-2500.0","name":"$500 to $2.500","results":6757},{"id":"2500.0-*","name":"More than $2.500","results":6719}]},{"id":"accepts_mercadopago","name":"MercadoPago filter","type":"boolean","values":[{"id":"yes","name":"With MercadoPago","results":20000}]},{"id":"installments","name":"Pago","type":"text","values":[{"id":"yes","name":"Installments","results":19711}]},{"id":"condition","name":"Condición","type":"text","values":[{"id":"new","name":"New","results":4567},{"id":"used","name":"Used","results":14850}]},{"id":"shipping","name":"Shipping","type":"text","values":[{"id":"mercadoenvios","name":"Mercado Envíos","results":2979},{"id":"free","name":"Free shipping","results":780},{"id":"sameday","name":"Moto Express","results":86}]},{"id":"power_seller","name":"Seller quality filter","type":"boolean","values":[{"id":"yes","name":"Best sellers","results":1263}]},{"id":"buying_mode","name":"Buying mode filter","type":"text","values":[{"id":"buy_it_now","name":"Buy it now","results":19731},{"id":"auction","name":"Auction","results":269}]},{"id":"since","name":"Auction start date filter","type":"text","values":[{"id":"today","name":"Publicados hoy","results":27}]},{"id":"until","name":"Auction stop filter","type":"text","values":[{"id":"today","name":"Ending today","results":40}]},{"id":"has_video","name":"Video publications filter","type":"boolean","values":[{"id":"yes","name":"Publications with video","results":251}]},{"id":"has_pictures","name":"Items with images filter","type":"boolean","values":[{"id":"yes","name":"With pictures","results":20000}]}]}`
}

func TestCategoryItemsRequestResponseParsing(t *testing.T) {

	Convey("Given a client", t, func() {
		mux := http.NewServeMux()
		ts := httptest.NewServer(mux)
		defer ts.Close()
		client := &meliClient{httpClient: &http.Client{}, baseUrl: ts.URL}
		Convey("When category items request response is valid it should be parsed into struct", func() {
			getCategoryItemsHandler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(getMockedCategoryItemsResponse()))
			}
			mux.HandleFunc("/sites/MLA/search/", getCategoryItemsHandler)
			categoryItemsResponse, _ := client.GetCategoryItems("MLA1367", 0, 2)
			expectedCategoryItems := &CategoryItemsResponse{
				Results: []Item{
					Item{Price: 1650},
					Item{Price: 11000},
				},
			}

			So(categoryItemsResponse, ShouldResemble, expectedCategoryItems)
		})
		Convey("When category request response is invalid it should return an error", func() {
			getCategoryItemsHandler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(""))
			}
			mux.HandleFunc("/sites/MLA/search/", getCategoryItemsHandler)
			_, err := client.GetCategoryItems("MLA1367", 0, 2)
			So(err, ShouldNotBeNil)
		})
		Convey("When category is invalid it should return an error", func() {
			getCategoryItemsHandler := func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusBadRequest)
				writer.Write([]byte(""))
			}
			mux.HandleFunc("/sites/MLA/search/", getCategoryItemsHandler)
			_, err := client.GetCategoryItems("MKK1367", 0, 2)
			So(err, ShouldNotBeNil)
		})

	})

}
