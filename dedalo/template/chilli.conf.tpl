cmdsocket       /run/dedalo/chilli.sock
unixipc         /run/dedalo/chilli.ipc
pidfile         /run/dedalo/chilli.pid

wwwdir          /opt/icaro/dedalo/www
wwwbin          /etc/chilli/wwwsh

net             ${HS_NETWORK}
dhcpif          ${HS_INTERFACE}
tundev          tun-dedalo
uamanydns
dns1            "${HS_DNS1:-208.67.222.222}"
dns2            "${HS_DNS2:-208.67.220.220}"
uamport         3990
locationname    "${HS_UNIT_NAME}"
uamaaaurl       "${HS_AAA_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}&timezone=${HS_TIMEZONE}"
radiusserver1   "localhost"
uamserver       "${HS_SPLASH_PAGE_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}&timezone=${HS_TIMEZONE}"
radiusnasid     "${HS_ID}"
alloworigin     "${HS_ALLOW_ORIGINS}"
macauth

uamdomain       "${HS_AAA_HOST}"

include /opt/icaro/dedalo/walled_gardens/local.conf
include /opt/icaro/dedalo/walled_gardens/facebook.conf
include /opt/icaro/dedalo/walled_gardens/linkedin.conf
include /opt/icaro/dedalo/walled_gardens/instagram.conf

${USE_UPDOWN_SCRIPTS:+"ipup=/opt/icaro/dedalo/up.sh"}
${USE_UPDOWN_SCRIPTS:+"ipdown=/opt/icaro/dedalo/down.sh"}

include /opt/icaro/dedalo/local.conf
