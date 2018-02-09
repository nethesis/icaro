---
title: "Client installation"
permalinks: /docs/client_installation/
---

The client is composed by two parts:

- CoovaChilli: you will need a [patched version of CoovaChilli](https://github.com/NethServer/coova-chilli)
- Dedalo: a wrapper around to CoovaChilli

Currently, installation is supported only on [CentOS 7](https://www.centos.org) and [NethServer 7](https://www.nethserver.org).

We are also working to support in on [OpenWrt](https://openwrt.org/).

## RPMs

RPMs for CentOS are automatically built on each release.

- Patched [CoovaChilli]( https://github.com/NethServer/rpm-coova-chilli), available at [NethServer repositories](http://mirror.nethserver.org/nethserver/7/testing/x86_64/Packages/)
- [Dedalo](https://github.com/nethesis/icaro/tree/master/dedalo/dist), available on [Github](https://github.com/nethesis/icaro/releases)

## Install on CentOS

Make sure your system has Internet access and it's connected to a working access point, then:

1. Download and install `dedalo` and `coova-chilli` RPMs from the above links

2. Configure your Dedalo instance using instructions reported inside the [README](https://github.com/nethesis/icaro/tree/master/dedalo)

3. Make sure to enable Dedalo at boot:
   ```
   systemctl enable dedalo
   ```


## Install on NethServer

The installation requires at least 3 ethernet interfaces: 

- 1 for normal LAN clients, marked with ``green`` role
- 1 (or more) for Internet connection, marked with ``red`` role
- 1 one for the Dedalo, marked with ``hotspot`` role

Proceed with the configuration of green and red interfaces, then:

1. Install the ``nethserver-dedalo`` package

   Access a shell and execute:
   ```
   yum --enablerepo=nethserver-testing install nethserver-dedalo
   ```

2. Access the ``Hotspot unit`` page inside the Server Manager (``http://<SERVER_FQDN>:980``)

3. On the first wizard page, insert Sun host name and enter valid reseller credentials.
   Then follow the istructions on the screen.
