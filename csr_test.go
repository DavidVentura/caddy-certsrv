package certsrv

import (
	"strings"
	"testing"
)

func TestCSRGeneration(t *testing.T) {
	_, pem, err := MakeCSR("example.com")
	if err != nil {
		panic(err)
	}
	if !strings.Contains(string(pem), "BEGIN CERTIFICATE REQUEST") {
		t.Errorf("Expected 'BEGIN CERTIFICATE REQUEST' in PEM CSR, got:\n%s", pem)
	}
}
