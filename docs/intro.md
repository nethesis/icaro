# Introduction

Icaro is a hotspot built for small and medium public activities like hotels and pubs.
It handles 2 different business models.

## Free HotSpot

* Users are attracted to your location because of the possibility of connecting to Internet for free
* Users login via a splash page in a very easy way, at the same time we can use the splash page to show ADs and collect users data (the Wi-Fi network become an effective lead generation system)
* Total control over data and access via Dashboard
* Export data for marketing or use via specific APIs

## Paid HotSpot

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

Example: https://nethservice.nethesis.it/nextcloud/index.php/apps/gallery/s/SVvOOQcYMEfvFko

### Unity

A real server (Physical or Virtual) with NethSecurity associated to an Hotspot instance 

### Nethesis portal

Web portal located at [my.nethesis.it](https://my.nethesis.it) accessible only by Nethesis resellers.


## Abstract workflow

A reseller who wants to sell the Hotspot service to a customer, usually does the following steps.

1. Installs a new firewall inside the hotel and configures the local CoovaChilli instance connected to the wifi network.
2. Creates a new Hotspot instance using the API.
   Nethesis resellers can do it by accessing Nethesis portal and enter the "Applications" -> "Hotspot" page.
3. Associates one or more server (unities) to the newly created Hotspot instance using the API.
   Nethesis resellers can use the above portal, the list of selectable unities
   displays all servers not associated to any other Hotspot (the list could be filtered by customer)
4. Creates administrator access for the Hotspot instance using the API (or the portal above)
6. Enter the Hotspot Manager with the specific administrator user/pass and generate credentials for other profiles.
