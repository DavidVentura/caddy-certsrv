Name:           caddy
Version:        1.0
Release:        1%{?dist}
Summary:        Caddy server RPM package

Conflicts:      caddy
License:        MIT
Source0:        %{name}-%{version}.tar.gz

%description
Caddy is a web server written in Go. This package provides the binary and service files to run it as a systemd service.
It's been built with the following plugins:
- certsrv

%prep
%setup
mkdir -p asd

%install
install -m 755 -d %{buildroot}/usr/bin
install -m 755 caddy %{buildroot}/usr/bin

install -m 755 -d %{buildroot}/usr/lib/systemd/system
install -m 644 caddy.service %{buildroot}/usr/lib/systemd/system

%pre
/usr/bin/getent group caddy  || /usr/sbin/groupadd -r caddy
/usr/bin/getent passwd caddy || /usr/sbin/useradd -r -d /var/caddy -s /sbin/nologin -g caddy caddy
mkdir -p /var/caddy
chown caddy:caddy /var/caddy

%files
/usr/bin/caddy
/usr/lib/systemd/system/caddy.service
