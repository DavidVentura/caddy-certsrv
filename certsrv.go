package certsrv

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

// Returns a PEM-encoded ceritifcate
func fetchCert(cl *spnego.Client, certSrvUrl string, certUrl string, ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/%s", certSrvUrl, certUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := cl.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	cert, err := io.ReadAll(resp.Body)
	return string(cert), nil
}

// Returns a PEM-encoded ceritifcate
func MakeCert(cl *spnego.Client, certSrvUrl string, csr *x509.CertificateRequest, ctx context.Context) (string, error) {
	req := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csr.Raw,
	})

	data := url.Values{
		"Mode":        {"newreq"},
		"CertRequest": {string(req)},
		"CertAttrib":  {"CertificateTemplate:WebServer(PrivateKeyExportable)"},
	}

	// spnego.Client does not support cancellation
	body := strings.NewReader(data.Encode())
	url := fmt.Sprintf("%s/certfnsh.asp", certSrvUrl)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cl.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	link, err := parseHTMLResponse(resp.Body)
	if err != nil {
		return "", err
	}
	return fetchCert(cl, certSrvUrl, link, ctx)
}
