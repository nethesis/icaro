cmdsocket       /run/dedalo/chilli.sock
unixipc         /run/dedalo/chilli.ipc
pidfile         /run/dedalo/chilli.pid

net             10.1.0.0/255.255.255.0
uamlisten       10.1.0.1
dhcpif          ${HS_INTERFACE}
uamanydns
dns1            "208.67.222.222"
dns2            "208.67.220.220"
uamport         3990
locationname    "${HS_UNIT_NAME}"
uamaaaurl       "${HS_AAA_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}"
radiusserver1   "localhost"
uamserver       "${HS_SPLASH_PAGE_URL}/?digest=${HS_DIGEST}&uuid=${HS_UUID}"
radiusnasid     "${HS_ID}"
alloworigin     "${HS_ALLOW_ORIGINS}"

include /opt/icaro/dedalo/walled_gardens/local.conf
include /opt/icaro/dedalo/walled_gardens/facebook.conf
include /opt/icaro/dedalo/walled_gardens/google.conf
include /opt/icaro/dedalo/walled_gardens/linkedin.conf

ipup=/etc/chilli/up.sh
ipdown=/etc/chilli/down.sh