package order

import (
	"api-gateway-svc/config"
	"log"
	"net/url"

	"api-gateway-svc/helpers"

	"github.com/gorilla/mux"
)




func RegisterService(r *mux.Router, c config.Config) {
	url, err := url.Parse(c.OrderSvcUrl)
	if err != nil {
		log.Fatal("Unable to pass order service url")
	}
	r.HandleFunc("/order", 	helpers.RegisterProxyService(url))
}
