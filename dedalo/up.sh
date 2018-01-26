#!/bin/sh

TUNTAP=$(basename $DEV)
UNDO_FILE=/var/run/dedalo.$TUNTAP.sh

HS_WANIF=${HS_WANIF:-$(route -n|grep '^0.0.0.0'|head -n1|awk '{print $8}')}

[ -e "$UNDO_FILE" ] && sh $UNDO_FILE 2>/dev/null
rm -f $UNDO_FILE 2>/dev/null

ipt() {
    opt=$1; shift
    echo "iptables -D $*" >> $UNDO_FILE
    iptables $opt $*
}

ipt_in() {
    ipt -I INPUT -i $TUNTAP $*
}

run_up() {
    ipt_in --dst $ADDR -j DROP

    ipt_in -p tcp -m tcp --dport $UAMPORT --dst $ADDR -j ACCEPT

    ipt_in -p udp -d 255.255.255.255 --destination-port 67:68 -j ACCEPT
    ipt_in -p udp -d $ADDR --destination-port 67:68 -j ACCEPT
    ipt_in -p udp --dst $ADDR --dport 53 -j ACCEPT
    ipt_in -p icmp --dst $ADDR -j ACCEPT

    ipt -I INPUT -i $DHCPIF -j DROP

    ipt -I FORWARD -i $DHCPIF -j DROP
    ipt -I FORWARD -o $DHCPIF -j DROP

    ipt -I FORWARD -i $TUNTAP -j ACCEPT
    ipt -I FORWARD -o $TUNTAP -j ACCEPT

    # Help out conntrack to not get confused
    # (stops masquerading from working)
    #ipt -I PREROUTING -t raw -j NOTRACK -i $DHCPIF
    #ipt -I OUTPUT -t raw -j NOTRACK -o $DHCPIF

    # Help out MTU issues with PPPoE or Mesh
    ipt -I FORWARD -p tcp -m tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu
    ipt -I FORWARD -t mangle -p tcp -m tcp --tcp-flags SYN,RST SYN -j TCPMSS --clamp-mss-to-pmtu

    ipt -I FORWARD -i $TUNTAP \! -o $HS_WANIF -j DROP

    ipt -I FORWARD -i $TUNTAP -o $HS_WANIF -j ACCEPT

    # site specific stuff optional
    [ -e /opt/icaro/dedalo/ipup.sh ] && . /opt/icaro/dedalo/ipup.sh
}

FLOCK=$(which flock)

if [ -n "$FLOCK" ] && [ -z "$LOCKED_FILE" ]
then
    export LOCKED_FILE=/tmp/.dedalo-flock
    flock -x $LOCKED_FILE -c "$0 $@"
else
    run_up
fi
