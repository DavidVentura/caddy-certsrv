package certsrv

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

// Returns a PEM-encoded ceritifcate
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
	defer resp.Body.Close()
	cert, err := io.ReadAll(resp.Body)
	return string(cert)
}

func MakeCert(cl *spnego.Client, certSrvUrl string, csr *x509.CertificateRequest) string {
	// FIXME raw?
	req := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csr.Raw,
	})

	data := url.Values{
		"Mode":        {"newreq"},
		"CertRequest": {string(req)},
		"CertAttrib":  {"CertificateTemplate:WebServer(PrivateKeyExportable)"},
	}

	resp, err := cl.PostForm(fmt.Sprintf("%s/certfnsh.asp", certSrvUrl), data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	link, err := parseHTMLResponse(resp.Body)
	if err != nil {
		panic(err)
	}
	return fetchCert(cl, certSrvUrl, link)
}
