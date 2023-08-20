package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	certsrv "github.com/davidventura/caddy-certsrv"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func fetchCert(cl *spnego.Client, certSrvUrl string, certUrl string) string {
	url := fmt.Sprintf("%s/%s", certSrvUrl, certUrl)
	fmt.Println(url)
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := cl.Do(r)
	if err != nil {
		panic(err)
	}
	log.Printf("%#v\n", resp)
	defer resp.Body.Close()
	cert, err := io.ReadAll(resp.Body)
	return string(cert)
}
func reqCert(cl *spnego.Client, csr string, certSrvUrl string) string {
	data := url.Values{
		"Mode":        {"newreq"},
		"CertRequest": {csr},
		"CertAttrib":  {"CertificateTemplate:WebServer(PrivateKeyExportable)"},
	}

	r, err := http.NewRequest("PostForm", fmt.Sprintf("%s/certfnsh.asp", certSrvUrl), strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	resp, err := cl.Do(r)
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
	if len(os.Args) != 6 {
		fmt.Printf("Expected 5 args: cert-domain cert-srv keytab username realm\n")
		os.Exit(1)
	}

	certFor := os.Args[1]
	certSrv := os.Args[2]
	keytabPath := os.Args[3]
	username := os.Args[4]
	realm := os.Args[5]

	keytab, err := keytab.Load(keytabPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Obtaining a cert for %s by asking %s\n", certFor, certSrv)
	krb5Str := `
`
	cfg, err := config.NewConfigFromString(krb5Str)
	if err != nil {
		panic(err)
	}
	cl := client.NewClientWithKeytab(username, realm, keytab, cfg)
	spnegoCl := spnego.NewClient(cl, nil, "")

	_, pem, err := certsrv.MakeCSR(certFor)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(pem))
	link := reqCert(spnegoCl, string(pem), certSrv)
	fmt.Printf("%s\n", link)
	cert := fetchCert(spnegoCl, certSrv, link)
	fmt.Printf("%s\n", cert)
}
