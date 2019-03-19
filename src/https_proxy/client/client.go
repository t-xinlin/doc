package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	u, err := url.Parse("https://localhost:8888")
	if err != nil {
		//panic(err)
		fmt.Printf("error1 : %+v", err)
	}
	tr := &http.Transport{
		Proxy:        http.ProxyURL(u),
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:26326")
	if err != nil {
		//panic(err)
		fmt.Printf("error2 : %+v", err)
	}
	defer resp.Body.Close()
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		//panic(err)
		fmt.Printf("error3 : %+v", err)
	}
	fmt.Printf(">>>>>>>>>>>>  %+v", dump)
}
