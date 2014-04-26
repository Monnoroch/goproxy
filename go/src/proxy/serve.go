package proxy

import "fmt"
import "strconv"
import "net/http"
import "net/http/httputil"


func getServer(port int, scheme string, rules []RewriteRule) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.ParseForm()
			for _, rule := range rules {
			    if rule.Active(req) {
			        rule.Rewrite(req)
			        break
			    }
			}
			req.URL.Scheme = scheme
			req.URL.Host = req.Host
		},
	})

	return &http.Server{
		Addr: ":" + strconv.Itoa(port),
		Handler: mux,
	}
}

func ListenAndServe(port int, rules []RewriteRule) {
	s := getServer(port, "http", rules)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println("http.Server.ListenAndServe error: %v.", err)
		}
	}()
}

func ListenAndServeTLS(port int, cert string, key string, rules []RewriteRule) {
	s := getServer(port, "http", rules)
	go func() {
		err := s.ListenAndServeTLS(cert, key)
		if err != nil {
			fmt.Println("http.Server.ListenAndServeTLS error: %v.", err)
		}
	}()
}
