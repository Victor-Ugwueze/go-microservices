package helpers

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

	func RegisterProxyService(url * url.URL) func(w http.ResponseWriter, r *http.Request) {
		
		proxy := httputil.NewSingleHostReverseProxy(url)

		return func(w http.ResponseWriter, r *http.Request) {
			// fmt.Println(mux.Vars(r)["path"])
			path := strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
			path = fmt.Sprintf("%s%s", "/", path)
			r.URL.Path = path
			fmt.Println(path, strings.Split(r.URL.Path, "/"))
			proxy.ServeHTTP(w, r)
		}
	}