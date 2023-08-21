package certsrv

import (
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"
)

func MakeClientWithPassword(certSrv string, username string, password string, realm string) (*spnego.Client, error) {
	cfg, err := config.Load("/etc/krb5.conf")
	if err != nil {
		return nil, err
	}
	cl := client.NewClientWithPassword(username, realm, password, cfg, client.DisablePAFXFAST(true))
	return spnego.NewClient(cl, nil, ""), nil
}

func MakeClientWithKeytab(certSrv string, username string, kt *keytab.Keytab, realm string) (*spnego.Client, error) {
	cfg, err := config.Load("/etc/krb5.conf")
	if err != nil {
		return nil, err
	}
	cl := client.NewClientWithKeytab(username, realm, kt, cfg, client.DisablePAFXFAST(true))
	return spnego.NewClient(cl, nil, ""), nil
}
