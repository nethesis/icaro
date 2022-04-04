cmdsocket       /run/dedalo/chilli.sock
unixipc         /run/dedalo/chilli.ipc
pidfile         /run/dedalo/chilli.pid

wwwdir          /opt/icaro/dedalo/www
wwwbin          /etc/chilli/wwwsh

net             ${HS_NETWORK}
dhcpif          ${HS_INTERFACE}
tundev          tun-dedalo
uamanydns
dns1            "${HS_DNS1:-8.8.8.8}"
dns2            "${HS_DNS2:-1.1.1.1}"
uamport         3990
locationname    "${HS_UNIT_NAME}"
uamaaaurl       "${HS_AAA_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}&timezone=${HS_TIMEZONE}"
radiusserver1   "localhost"
uamserver       "${HS_SPLASH_PAGE_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}&timezone=${HS_TIMEZONE}"
radiusnasid     "${HS_ID}"
alloworigin     "${HS_ALLOW_ORIGINS}"
${HS_DHCPSTART:+"dhcpstart  $HS_DHCPSTART"}
${HS_DHCPEND:+"dhcpend  $HS_DHCPEND"}
macauth

uamdomain       "${HS_AAA_HOST}"
${HS_MAXCLIENTS:+"maxclients       "$HS_MAXCLIENTS""}

include /opt/icaro/dedalo/walled_gardens/local.conf
include /opt/icaro/dedalo/walled_gardens/facebook.conf
include /opt/icaro/dedalo/walled_gardens/google.conf
include /opt/icaro/dedalo/walled_gardens/linkedin.conf
include /opt/icaro/dedalo/walled_gardens/instagram.conf
include /opt/icaro/dedalo/walled_gardens/wifi4eu.conf

# list of integrations walled garden
${HS_WALLED_GARDEN_INTEGRATIONS}

${USE_UPDOWN_SCRIPTS:+"ipup=/opt/icaro/dedalo/up.sh"}
${USE_UPDOWN_SCRIPTS:+"ipdown=/opt/icaro/dedalo/down.sh"}

include /opt/icaro/dedalo/local.conf
