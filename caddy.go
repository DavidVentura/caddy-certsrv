// Copyright 2015 Matthew Holt and The Caddy Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certsrv

import (
	"context"
	"crypto/x509"
	"errors"
	"strings"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
	"gopkg.in/jcmturner/gokrb5.v7/spnego"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func init() {
	caddy.RegisterModule(CertSrvIssuer{})
}

// CertSrvIssuer can request certificates from a
// Microsoft Active Directory Certificate Services instance
type CertSrvIssuer struct {
	CertSrvUrl string `json:"certsrv_url"`
	Realm      string `json:"realm"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty_empty"`
	KeytabFile string `json:"keytab_path,omitempty_empty"`

	keytab *keytab.Keytab
	cl     *spnego.Client
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (CertSrvIssuer) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "tls.issuance.certsrv",
		New: func() caddy.Module { return new(CertSrvIssuer) },
	}
}

// Provision sets up the issuer.
func (iss *CertSrvIssuer) Provision(ctx caddy.Context) error {
	iss.logger = ctx.Logger()
	var cl *spnego.Client
	var err error
	if iss.KeytabFile != "" {
		kt, err := keytab.Load(iss.KeytabFile)
		if err != nil {
			return err
		}
		iss.keytab = kt
		cl, err = MakeClientWithKeytab(iss.CertSrvUrl, iss.Username, iss.keytab, iss.Realm)
	} else {
		cl, err = MakeClientWithPassword(iss.CertSrvUrl, iss.Username, iss.Password, iss.Realm)
	}

	if err != nil {
		return err
	}

	iss.cl = cl

	return nil
}

// Validate config
func (iss *CertSrvIssuer) Validate() error {
	iss.logger.Info("My config is\n", zap.Any("config", iss))
	if !strings.Contains(iss.CertSrvUrl, "//") {
		return errors.New("certsrv_url must be a valid URL")
	}
	if iss.Password == "" && iss.KeytabFile == "" {
		return errors.New("Provide Password or KeytabFile")
	}
	return nil
}

// IssuerKey returns the unique issuer key for the
// confgured CA endpoint.
func (iss CertSrvIssuer) IssuerKey() string {
	return "totally a unique key"
}

// Issue issues a certificate to satisfy the CSR.
func (iss CertSrvIssuer) Issue(ctx context.Context, csr *x509.CertificateRequest) (*certmagic.IssuedCertificate, error) {
	iss.logger.Info("Getting asked to pass a CSR for %s\n", zap.Stringer("Subject", csr.Subject))
	// TODO: honor cancellation in MakeCert etc
	cert := MakeCert(iss.cl, iss.CertSrvUrl, csr)

	return &certmagic.IssuedCertificate{
		// The PEM-encoding of DER-encoded ASN.1 data.
		Certificate: []byte(cert),
	}, nil
}
func (iss *CertSrvIssuer) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "certsrv_url":
				if !d.NextArg() {
					return d.ArgErr()
				}
				iss.CertSrvUrl = d.Val()
			case "realm":
				if !d.NextArg() {
					return d.ArgErr()
				}
				iss.Realm = d.Val()
			case "username":
				if !d.NextArg() {
					return d.ArgErr()
				}
				iss.Username = d.Val()
			case "password":
				if !d.NextArg() {
					return d.ArgErr()
				}
				iss.Password = d.Val()
			case "keytab_path":
				if !d.NextArg() {
					return d.ArgErr()
				}
				iss.KeytabFile = d.Val()
			}
		}
	}
	return nil
}

// Interface guards
var (
	_ caddy.Validator   = (*CertSrvIssuer)(nil)
	_ caddy.Provisioner = (*CertSrvIssuer)(nil)
	_ certmagic.Issuer  = (*CertSrvIssuer)(nil)

// _ provisioner.CertificateModifier = (*customCertLifetime)(nil)
)
