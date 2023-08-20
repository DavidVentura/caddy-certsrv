package certsrv

import (
	"fmt"

	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func MakeClient(certSrv string, username string, password string, realm string) *spnego.Client {
	cfg, err := config.Load("/etc/krb5.conf")
	if err != nil {
		panic(err)
	}
	//cl := client.NewClientWithKeytab(username, realm, keytab, cfg)
	cl := client.NewClientWithPassword(username, realm, password, cfg, client.DisablePAFXFAST(true))
	fmt.Printf("user %s realm %s kt %s\n", username, realm, "keytabPath")
	return spnego.NewClient(cl, nil, "")
}
