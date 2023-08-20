package certsrv

import (
	"os"
	"testing"
)

func TestParseHTMLResponse(t *testing.T) {
	data, err := os.Open("testdata/example1.html")
	if err != nil {
		panic(err)
	}
	link, err := ParseHTMLResponse(data)
	if err != nil {
		panic(err)
	}
	want := "certnew.cer?ReqID=61108&Enc=b64"
	if link != want {
		t.Errorf("got %s want %s\n", link, want)
	}
}
