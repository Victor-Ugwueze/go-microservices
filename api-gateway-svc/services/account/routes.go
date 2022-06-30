package account

import (
	"api-gateway-svc/config"
	"log"
	"net/url"

	"api-gateway-svc/helpers"

	"github.com/gorilla/mux"
)


func RegisterService(r *mux.Router, c config.Config) {
	url, err := url.Parse(c.AccountSvcUrl)
	if err != nil {
		log.Fatal("Unable to pass account service url")
	}
	r.HandleFunc("/account/{rest}", 	helpers.RegisterProxyService(url))
}
