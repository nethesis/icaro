---
title: "Captive portal"
permalink: /docs/captive/
---

The Captive Portal is a web page which is displayed to newly connected users before they are granted Internet access.
It is a responsive web page and must require user credentials if the user isn't already authenticated.
The Captive Portal is used to authenticate users but also for showing Ads and give a better image of the hotel/bar.

The Captive Portal has the following features:

* it is installed on a remote public server (in-cloud)
* it has to be customizable by the reseller in a simple and effective way, the result should be good looking (e.g. limit customization possibilities)
* it must permit to insert Ads or extra information (function accesible by the reseller)
* it must allow to show a _landing page_ , usually used just for Ads prior to go to the login page
* it must allow users to report Hotspot problems to the staff (e.g. used in squares without any desk person) (**not yet implemented**)

Example of landing page:

![Captive portal](../img/captive.png "Schema")


### Login Methods and type of Accounts
For more details, see [Login configuration](/icaro/docs/configuration/) page, to properly configure each logins option.

#### Login Methods

Voucher (used for security/payment) can be:

 * `none`: Hotspot access is free, use one of the auths below to navigate, no voucher is required

 * `voucher unlimited` : this kind of voucher can be re-used infinite times, allowing you to protect the access from people outside the company structure, just communicate this voucher code to guests inside the company and they will be able to autenticate themselves using one of the authentication methods to navigate

* `voucher limited` : this kind of voucher can be re-used only a finite number of times, its tipical use is with premium (payed) accounts with better conditions (autologin, more bandwidth, more traffic and so on).

* `voucher authentication` : this kind of voucher can be used to authenticate immediately the guest by inserting the code. This voucher can be used combined with the above voucher or alone.

Voucher can be exported in CSV format or printed (individually or globally).

Login type after the voucher check:

  * Email
  * SMS
  * Social (only if voucher are not payed)
    * Facebook
    * Linkedin
  * Code (this is the voucher authentication explained above)

Social logins are valid only for a specific period of time.

The period can be customizable for each Hotspot instance. Available choices:

- 1 week
- 2 weeks
- 1 month
- 3 months
- 6 months
- 1 year


Each login can be used on any units referring to single Hotspot instance.
For example, given an hotel chain with an Hotspot instance, the user can login with same
credentials on any hotel belonging to the chain.

The system should track to which unity a single account is connected.

##### Email auth
When user requests an Email authentication, Dedalo open a temporary session of N minutes (5 minutes default) the make possible receive the email with the authentication code.

##### SMS service

The system should be functional only with one or two SMS providers (eg. [Twilio](https://www.twilio.com/)) to avoid configuration by the user.

Also WhatsApp Business has been released, it could be an alternative way to make login without pay SMS
Drawback: WA should be permitted without authentication, will they make login to the Hotspot?

#### Guests Limits and features

We can limit guests operativity with different parameters:

* Maximum download Bandwidth (kbps)
* Maximum upload Bandwidth (kbps)
* Daily Navigation Time limit (Minutes)
* Daily Navigation Traffic limit (MB)

The **Autologin** feature allows guests to autenticate themselves just once, the next time the hotspot will recognize the guests and will make them browse without requiring any authentication.


#### Special device accounts

Devices in the Hotspot unit network which needs direct access to the Internet (eg. printers, special workstations),
can have special access rights directly in the firewall configuration using CoovaChilli.
