---
title: Network
permalink: /docs/network/
---

## Hotspot Configuration 

* You must have at least one free interface that will be used **only for hotspot purposes**, the interface can be **physical or logical (VLAN)**
* Configure the Hotspot service specifying the correct interface
* Define a private network for the Hotspot (the clients will receive IP addresses belonging to this network)
* Connect all the Access Points to the interface 


## Access Points Configuration

* Access Point should only allow the traffic between clients and the hotspot service, they must behave like dumb network switches
* Configure APs with open access to everybody without authentication, disable DHCP server on every APs because the DHCP service must be performed by the hotspot service
* Disable every service (PnP, security and so on) in order to avoid any interference with the hotspot service
* Map the APs in a private network different from the one used for the hotspot clients (the IP addresses can be static or assigned by the AP's controller if present)
* In case of multiple APs without a controller use different SSIDs (es: SCHOOL-1, SCHOOL-2 and so on ) in order to find more easily problems in the AP units
* Use different radio channels on the APs to minimize interferences 
* Enable client isolation if present
* If the hotspot make use of a VLAN configure the switches where APs are linked, so that VLAN tags are not cutted by the switches
* Finding problems in Hotspot environments can be very tricky because there are many variables that you can't control easily (client devices,access points, auth methods and many others...), please avoid low quality AP cause they can generate further problems.
* With the APs provide a good  coverage of the area where the service must be available, lack of signal also is linked to a bad quality service and increase the probability of malfunctioning
