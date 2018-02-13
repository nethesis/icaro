Name: dedalo
Version: 0.0.3
Release: 1%{?dist}
Summary: Network Access Controller, runs on the firewall and intercepts all guest connections

License: GPLv3
URL: https://github.com/nethesis/icaro
Source0: https://github.com/nethesis/icaro/archive/master.tar.gz

%{?systemd_requires}
BuildRequires: systemd
Requires: coova-chilli

%description
Dedalo is the Network Access Controller, runs on the firewall and intercepts all guest connections.
It's a configuration helper for Coova-Chilli

%post
%systemd_user_post dedalo.service

%preun
%systemd_user_preun dedalo.service


%prep
%setup -n icaro-master


%build


%install
mkdir -p %{buildroot}/opt/icaro/dedalo
install -D -m644 dedalo/dedalo.service %{buildroot}/%{_unitdir}/dedalo.service
install -D -m644 dedalo/config %{buildroot}/opt/icaro/dedalo/config
install -D -m755 dedalo/dedalo %{buildroot}/%{_bindir}/dedalo
mkdir -p %{buildroot}/opt/icaro/dedalo/template
install -D -m644 dedalo/template/chilli.conf.tpl %{buildroot}/opt/icaro/dedalo/template/chilli.conf.tpl
install -D -m775 dedalo/template/engine %{buildroot}/opt/icaro/dedalo/template/engine
mkdir -p %{buildroot}/opt/icaro/dedalo/walled_gardens
install -D -m644 dedalo/walled_gardens/facebook.conf %{buildroot}/opt/icaro/dedalo/walled_gardens/facebook.conf
install -D -m644 dedalo/walled_gardens/linkedin.conf %{buildroot}/opt/icaro/dedalo/walled_gardens/linkedin.conf
install -D -m644 dedalo/walled_gardens/instagram.conf %{buildroot}/opt/icaro/dedalo/walled_gardens/instagram.conf
mkdir -p %{buildroot}/opt/icaro/dedalo/wwww
touch %{buildroot}/opt/icaro/dedalo/local.conf
touch %{buildroot}/opt/icaro/dedalo/walled_gardens/local.conf

%files
/usr/lib/systemd/system/dedalo.service
%{_bindir}/dedalo
%dir /opt/icaro/
/opt/icaro/dedalo/template/engine
/opt/icaro/dedalo/wwww
%config(noreplace) /opt/icaro/dedalo/config
%config /opt/icaro/dedalo/template/chilli.conf.tpl
%config /opt/icaro/dedalo/walled_gardens/facebook.conf
%config /opt/icaro/dedalo/walled_gardens/linkedin.conf
%config /opt/icaro/dedalo/walled_gardens/instagram.conf
%config(noreplace) /opt/icaro/dedalo/walled_gardens/local.conf
%config(noreplace) /opt/icaro/dedalo/local.conf
%doc dedalo/README.md


%changelog
* Wed Feb 07 2018 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.0.3
- Add local config files
- Get rid of jq dependency

* Wed Feb 07 2018 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.0.2
- Use http codes in dedalo register command

* Tue Jan 16 2018 Edoardo Spadoni <edoardo.spadoni@nethesis.it> - 0.0.1
- First public release

