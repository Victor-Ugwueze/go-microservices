package product

import (
	"api-gateway-svc/config"
	"log"
	"net/url"

	"api-gateway-svc/helpers"

	"github.com/gorilla/mux"
)




func RegisterService(r *mux.Router, c config.Config) {
	url, err := url.Parse(c.ProductSvcUrl)
	if err != nil {
		log.Fatal("Unable to pass product service url")
	}
	r.HandleFunc("/product", 	helpers.RegisterProxyService(url))
}

