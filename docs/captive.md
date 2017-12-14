### The Captive Portal

The Captive Portal is  a web page which is displayed to newly Hotspot connected users before they are granted Internet access.
It is a responsive web page and must require user credentials if the user isn't already authenticated.
The Captive Portal can be used to make login but also for showing Ads and give a better image of the hotel/bar.

The Captive portal has the following features:

* it is installed on a remote public server (in-cloud)
* it has to be customizable by the reseller in a simple and effective way, the result should be good looking (e.g. limit customization possibilities)
* it must permit to insert ADs or informations (always made by the reseller)
* it must allow to show a _landing page_ , usually used just for Ads prior to go to the login page
* it must allow users to report Hotspot problems to the staff (e.g. used in squares without any desk person)       

Landing page:


![landing page](https://nethservice.nethesis.it/nextcloud/index.php/s/PohSiE2E4CuIHRy/download)

Login page:


![login page](https://nethservice.nethesis.it/nextcloud/index.php/s/yFYGzRznxfWFxq1/download)

Aboce images five an idea of the login page.

### Login Methods and type of Accounts

#### Login Methods

Voucher (used for security/payment) can be:

 * `none`: Hotspot access is free, use one of the auths below to navigate

 * `block`: is used to protectd the access from people outside the company structure, then use one of the auths below to navigate

 * `paid`: the code can be sold from office desk for each user, then use one of the auths below to navigate

Login type after the voucher check:

  * Email
  * SMS   
  * Social (only if voucher are not payed)
    * Facebook
      * Like Gate not mandatory with timeout option (avoid to put like if you waiting for N seconds)
      * FB Checkin
    * Google
    * Instagram
    * Linkedin
    * Twitter

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

##### SMS service

The system should be functional only with one or two SMS providers (e.g. smshosting.it) to avoid configuration by the user.
Nethesis reseller could by prepaid packages of SMS from Nethesis shop.

Also WhatsApp Business is coming, it could be an alternative way to make login without pay SMS
Drawback: WA should be permitted without authentication, will they make login to the Hotspot?

#### Special device accounts

Devices in the Hotspot unit network which needs direct access to the Internet (eg. printers, special workstations),
can have special access rights directly in the firewall configuration using CoovaChilli.

