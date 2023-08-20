package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"os"

	certsrv "github.com/davidventura/caddy-certsrv"

	"gopkg.in/jcmturner/gokrb5.v7/keytab"
)

var a keytab.Keytab

func main() {
	if len(os.Args) != 6 {
		fmt.Printf("Expected 5 args: cert-domain cert-srv keytab username realm\n")
		os.Exit(1)
	}

	certFor := os.Args[1]
	certSrv := os.Args[2]
	//keytabPath := os.Args[3]
	password := os.Args[3]
	username := os.Args[4]
	realm := os.Args[5]

	/*keytab, err := keytab.Load(keytabPath)
	if err != nil {
		panic(err)
	}
	*/
	fmt.Printf("Obtaining a cert for %s by asking %s\n", certFor, certSrv)
	spnegoCl, err := certsrv.MakeClientWithPassword(certSrv, username, password, realm)

	keyBytes, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	subj := pkix.Name{
		CommonName: certFor,
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	if err != nil {
		panic(err)
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		panic(err)
	}

	cert := certsrv.MakeCert(spnegoCl, certSrv, csr)
	fmt.Printf("%s\n", cert)
}
