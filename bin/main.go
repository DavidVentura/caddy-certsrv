package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	certsrv "github.com/davidventura/caddy-certsrv"
)

func reqCert(csr string) string {
	data := url.Values{
		"Mode":        {"newreq"},
		"CertRequest": {csr},
		"CertAttrib":  {"CertificateTemplate:WebServer(PrivateKeyExportable)"},
	}
	resp, err := http.PostForm("http://example.com/upload", data)
	log.Printf("%#v\n", resp)
	defer resp.Body.Close()
	link, err := certsrv.ParseHTMLResponse(resp.Body)
	if err != nil {
		panic(err)
	}
	return link
}
func main() {
	_, pem, err := certsrv.MakeCSR("example.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(pem))
	reqCert(string(pem))
}
