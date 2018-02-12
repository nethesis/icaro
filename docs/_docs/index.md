---
title: Introduction
permalink: /docs/home/
redirect_from: /docs/index.html
---


Icaro has been built around two business models:

- free hotspot
- paid hotspot

## Free Hotspot

* Users are attracted to your location because of the possibility of connecting to Internet for free
* Users login via a splash page in a very easy way, at the same time we can use the splash page to show ADs and collect users data (the Wi-Fi network become an effective lead generation system)
* Total control over data and access via Dashboard
* Export data for marketing or use via specific APIs

## Paid Hotspot

* The administrator allow only paid authentication methods: `voucher`
* Users can see the splash page but not connect to Internet
* Users can buy a Wi-Fi voucher and connect to Internet for a certain period of time and with limitations
* Users should chose your connection because of its cost or performances (you should have a good connectivity )
* You can earn money selling Internet access


## Glossary

### Hotspot

It's an abstract entity which describes how connected users can access the Internet.
Each hotspot can have one or more attached firewall. Multiple firewall attached to same Hotspot share the same users and configuration.

### Hotspot manager

It's a web interface used to configure and manage the Hotspot instances and units.

### Captive portal

A responsive web page which is accessed by the customer using a smartphone or a PC, the page must require user credentials if the user isn't already authenticated.

### Unity

A real server (Physical or Virtual) with associated to an Hotspot instance.
It can be a [NethServer](https://www.nethserver.org) or another any supported platform.


## Abstract workflow

A reseller who wants to sell the Hotspot service to a customer, usually does the following steps.

1. Installs a new firewall inside the hotel and configures the local CoovaChilli instance connected to the WiFi network.

2. Creates a new Hotspot instance accessing Hotspot Manager (eg. https://myhotspot.nethesis.it) or using the APIs.

3. Associates one or more server (unities) to the newly created Hotspot instance accessing Hotspot Manager (or using APIs).

4. Enter the Hotspot Manager with personal cerdentials and  creates Customer and Desk users for the Hotspot instance

