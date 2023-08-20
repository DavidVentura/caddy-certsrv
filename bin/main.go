package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	certsrv "github.com/davidventura/caddy-certsrv"

	"github.com/MarshallWace/go-spnego"
)

func reqCert(csr string, certSrvUrl string) string {
	data := url.Values{
		"Mode":        {"newreq"},
		"CertRequest": {csr},
		"CertAttrib":  {"CertificateTemplate:WebServer(PrivateKeyExportable)"},
	}

	c := &http.Client{
		Transport: &spnego.Transport{},
	}

	resp, err := c.PostForm(certSrvUrl, data)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v\n", resp)
	defer resp.Body.Close()
	link, err := certsrv.ParseHTMLResponse(resp.Body)
	if err != nil {
		panic(err)
	}
	return link
}
func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Expected 2 args: cert-domain cert-srv\n")
		os.Exit(1)
	}
	certFor := os.Args[1]
	certSrv := os.Args[2]
	fmt.Printf("Obtaining a cert for %s by asking %s\n", certFor, certSrv)
	_, pem, err := certsrv.MakeCSR(certFor)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(pem))
	reqCert(string(pem), certSrv)
}
