package certsrv

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
)

func MakeCSR(cn string) (*rsa.PrivateKey, []byte, error) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 2048)

	subj := pkix.Name{
		CommonName: cn,
		/*
			Country:            []string{"AU"},
			Province:           []string{"Some-State"},
			Locality:           []string{"MyCity"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		*/
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	var pemBuf bytes.Buffer
	pem.Encode(&pemBuf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

	return keyBytes, pemBuf.Bytes(), nil
}
