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
	"time"

	"github.com/caddyserver/certmagic"
	"go.uber.org/zap"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddypki"
)

func init() {
	caddy.RegisterModule(CertSrvIssuer{})
}

// CertSrvIssuer is a certificate issuer that generates
// certificates internally using a locally-configured
// CA which can be customized using the `pki` app.
type CertSrvIssuer struct {
	// The ID of the CA to use for signing. The default
	// CA ID is "local". The CA can be configured with the
	// `pki` app.
	CA string `json:"ca,omitempty"`

	// The validity period of certificates.
	Lifetime caddy.Duration `json:"lifetime,omitempty"`

	// If true, the root will be the issuer instead of
	// the intermediate. This is NOT recommended and should
	// only be used when devices/clients do not properly
	// validate certificate chains.
	SignWithRoot bool `json:"sign_with_root,omitempty"`

	ca     *caddypki.CA
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

	// set some defaults
	if iss.CA == "" {
		iss.CA = caddypki.DefaultCAID
	}

	// get a reference to the configured CA
	appModule, err := ctx.App("pki")
	if err != nil {
		return err
	}
	pkiApp := appModule.(*caddypki.PKI)
	ca, err := pkiApp.GetCA(ctx, iss.CA)
	if err != nil {
		return err
	}
	iss.ca = ca

	// set any other default values
	if iss.Lifetime == 0 {
		iss.Lifetime = caddy.Duration(defaultInternalCertLifetime)
	}

	return nil
}

// IssuerKey returns the unique issuer key for the
// confgured CA endpoint.
func (iss CertSrvIssuer) IssuerKey() string {
	return iss.ca.ID
}

// Issue issues a certificate to satisfy the CSR.
func (iss CertSrvIssuer) Issue(ctx context.Context, csr *x509.CertificateRequest) (*certmagic.IssuedCertificate, error) {
	return nil, nil
	/*
		return &certmagic.IssuedCertificate{
			Certificate: buf.Bytes(),
		}, nil
	*/
}

// UnmarshalCaddyfile deserializes Caddyfile tokens into iss.
//
//	... internal {
//	    ca       <name>
//	    lifetime <duration>
//	    sign_with_root
//	}
func (iss *CertSrvIssuer) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		for d.NextBlock(0) {
			switch d.Val() {
			case "ca":
				if !d.AllArgs(&iss.CA) {
					return d.ArgErr()
				}

			case "lifetime":
				if !d.NextArg() {
					return d.ArgErr()
				}
				dur, err := caddy.ParseDuration(d.Val())
				if err != nil {
					return err
				}
				iss.Lifetime = caddy.Duration(dur)

			case "sign_with_root":
				if d.NextArg() {
					return d.ArgErr()
				}
				iss.SignWithRoot = true

			}
		}
	}
	return nil
}

const defaultInternalCertLifetime = 12 * time.Hour

// Interface guards
var (
	_ caddy.Provisioner = (*CertSrvIssuer)(nil)
	_ certmagic.Issuer  = (*CertSrvIssuer)(nil)

// _ provisioner.CertificateModifier = (*customCertLifetime)(nil)
)
