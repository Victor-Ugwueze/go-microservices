package helpers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

	func RegisterProxyService(url * url.URL) func(w http.ResponseWriter, r *http.Request) {
		
		proxy := httputil.NewSingleHostReverseProxy(url)

		return func(w http.ResponseWriter, r *http.Request) {
			// path := strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
			// path = fmt.Sprintf("%s%s", "/", path)
			// r.Host = url.Host
			// fmt.Println(path, strings.Split(r.URL.Path, "/"))
			proxy.ServeHTTP(w, r)
		}
	}