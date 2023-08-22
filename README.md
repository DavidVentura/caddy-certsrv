# Caddy Plugin for Microsoft Active Directory Certificate Services

🤮

Build with

```bash
xcaddy build --with  github.com/davidventura/caddy-certsrv=.
```

Package with

```bash
mkdir caddy-2.7.4
cp -t caddy-2.7.4/ caddy caddy.service Caddyfile
tar czf ~/rpmbuild/SOURCES/caddy-2.7.4.tar.gz caddy-2.7.4/
rpmbuild -ba caddy.spec
```

