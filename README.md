# Caddy Plugin for Microsoft Active Directory Certificate Services

🤮

Build with

```bash
xcaddy build --with  github.com/davidventura/caddy-certsrv=.
```

Package with

```bash
mkdir caddy-1.0
cp -t caddy-1.0/ caddy caddy.service
tar czf ~/rpmbuild/SOURCES/caddy-1.0.tar.gz caddy-1.0/
rpmbuild -ba caddy.spec
```

