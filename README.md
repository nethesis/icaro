# Icaro

<p align="center">
   <img src ="https://github.com/nethesis/icaro/raw/master/logo/logo.png" />
</p>

**"Easy and simple HotSpot for small and medium hotels"**

## Introduction

Icaro is a hotspot built for small and medium public activities like hotels and pubs.
It handles 2 different business models.

### Free HotSpot

* Users are attracted to your location because of the possibility to connect to internet for free
* Users make login via splash page in a very easy way, at the same time we can use the splash page to show ADs and collect users data (the Wi-Fi network become an effective lead generation system)
* Total control over data and access via DashBoard
* Export data for marketing or use via specific APIs 

### Payed HotSpot

* The administrator allow only payed authentication methods (e.g. COUPON)
* Users can see the splash page but not connect to internet
* Users can buy a Wi-Fi Coupon (or Voucher)and connect to internet for a certain period of time and with limitations
* Users should chose your connection because of its cost or performances (you should have a good connectivity )
* You can earn money selling Internet access


## Hotspot Manager

The **Hotspot Manager** is  web interface used to configure and manage the Hotspot.

There are 3 different profiles that can access the Hotspot Manager, each profile has its own purpose.
The Hotspot Manger must be written with a modular logic so that every profile can access only to sections intended for him.

Each profile has his/her own credentials, a profile give different authorizations to the user.

The 3 profiles are the following:

1. The Reseller
2. The Customer
3. The Desk

Desk can see very few things, Customer can see all that Desk see plus more panels, Reseller can see everything.

**1. The Reseller** 

The vendor (he sells Nethesis solutions)  which installs the physical firewall and creates an Hotspot entity; he can create users for his own Hotspot instances, this is the only one profile that can access also to my.nethesis.it.
The Reseller can see everything on the hotspot manager, he makes the first configuration of the hotspot, choosing the best options for the specific customer.

**Every technical or potentially dangerous configuration must be performed only by the reseller.**  
The reseller create its own access for the specific hotspot on _my.nethesis.it_, after that it will make access on _hotspot.nethesis.it_ and he will create access for other profiles (every hotspot will have only 3 profiles/account).
 
**2. The Customer**

The owner of the hotel or the local person responsible for the service, he doesn't have technical skills but he needs to do some tasks:

* modify the captive portal
* modify ads in the splash page
* check how many users are connectes
* read users statistics
* export reports
* make marketing activities

**3. The Desk**

Is the person at the desk of the hotel, Desk can:
* Create accounts
* See the list of accounts
* See how many accounts are connected and the bandwidth consumption (e.g. slow traffic )


## Software components

Icaro is divided in 3 parts:

- CoovaChilli: runs on the firewal and intercepts all guest connections
- Accounting server
- Captive portal


