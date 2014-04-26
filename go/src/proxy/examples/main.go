package main

import "fmt"
import "strings"
import "strconv"
import "flag"
import "os"
import "time"
import "math/rand"
import "net/http"
import "proxy"


func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


func main() {
	port := flag.Int("port", -1, "http port.")
	sslport := flag.Int("sslport", -1, "https port.")
	devport := flag.Int("devport", -1, "dev http port.")
	prodports := flag.String("prodports", "", "prod http ports.")
	ssldir := flag.String("ssldir", "", "directory with ssl key and certificate.")
	devIds := flag.String("dev-uids", "", "developer uids.")

	flag.Parse()

	if port == nil || *port == -1 || sslport == nil || *sslport == -1 ||
	   devport == nil || *devport == -1 ||
	   prodports == nil || *prodports == "" ||
	   ssldir == nil || *ssldir == "" ||
	   devIds == nil {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	
	rand.Seed(time.Now().UTC().UnixNano())
	prodps := []proxy.RewriteRule{}
	for _, p := range strings.Split(*prodports, ",") {
		port, err := strconv.Atoi(p)
		if err != nil {
			fmt.Printf("strconv.Atoi error: %v.\n", err)
			continue
		}
		prodps = append(prodps, proxy.NewRewritePortRule(port))
	}

	uids := strings.Split(*devIds, ",")
	userTrigger := func(req *http.Request) bool {
		return stringInSlice(req.Form.Get("uid"), uids)
	}

	rules := []proxy.RewriteRule{
		proxy.NewRewriteIfRule(userTrigger, proxy.NewRewritePortRule(*devport)),
		proxy.NewRewriteRandomRule(prodps),
	}

	proxy.ListenAndServe(*port, rules)
	proxy.ListenAndServeTLS(*sslport, *ssldir + "ssl.crt", *ssldir + "ssl.key", rules)

	select{}
}
