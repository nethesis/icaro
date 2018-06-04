---
title: Hotspot Manager
permalink: /docs/manager/
---


The **Hotspot Manager** is  web interface used to configure and manage the Hotspot.

## Roles

There are 3 different profiles that can access the Hotspot Manager, each profile has different authorizations
and can access only sections related to its own role.

The 3 profiles are the following:

1. Reseller: can see everything related to its own Hotspots
2. Customer: can see everything related to one Hotspot
3. Desk: can only see users related to one Hotspot

There is also a special `admin` user who has the ability to create/modify/delete resellers and edit manage all Hotspots.

### Reseller

The vendor which installs the physical firewall and creates an Hotspot entity; he can create users for his own Hotspot instances.
This is the only profile that can access creation API.
The Reseller can see everything on the Hotspot Manager, he makes the first configuration of the Hotspot, choosing the best options for the specific customer.

**Every technical or potentially dangerous configuration must be performed only by the reseller.**

The reseller can access the Hotppot Manager and create login credentials for other profiles.
Every Hotspot can have only 2 different profiles/accounts.

The reseller accounts can be created only from admin user.

### Customer

The owner of the hotel or the local person responsible for the service, he doesn't have technical skills but he needs to do some tasks:

* customize the captive portal
* customize ads in the splash page
* check how many users are connected
* read users statistics
* export reports
* perform marketing activities

### Desk

The is the person at the desk of the hotel, Desk can:

* delete accounts
* see the list of accounts
* see how many accounts are connected and the traffic each one has made
* create and delete vouchers 

### Hotspot manager views

List of the views with the associated profile: **R**eseller, **C**ustomer, **D**esk

* Dashboard (**R, C, D**)
* Hotspot 
    * See all hotspot instances (**R** only)
    * See only its hotspot instance (**C, D**)
    * Voucher management (**R, C, D**)
    * Technical Preferences like auth type, bandwidth... (**R** only)
    * Captive Portal customization (**R, C, D**)
    * MAC address allowed (for special devices) (**R, C**)
* Units ((**R, C**)
* Guests ((**R, C, D**)
* Devices ((**R, C**)
* Sessions (**R,C**)
* Report (**R,C,D**)
* Marketing (**R, C**) (**not yet implemented a specific panel** but some data can be exported)



## Single sign-on

(**Not yet implemented**)

To integrate resellers authentication, sharing the same credentials from another platform that already has the profiles credentials, Icaro must integrate an API to include resellers without adding one by one using a brand new API, using a master key auth mode.

Scenario

A company has an internal platform that already has customers or resellers and they login with their credentials in that portal. The company wants to enable the login of that resellers also in their Icaro instance and to make this possibile Icaro must include an API to integrate the resellers easily and programmatically.

## The Dashboard

(**Partially implemented**)

The dashboard shows informations about the use of the service.
It has 3 levels of informations depending on the user authorizations.

1. General data about the use of the service (login and SMS users only)
2. Social network data (when we have login via social network)
3. Data obtained through a questionnaire (mainly hotel use)

## General


General data about the use of the service (login and SMS users only)

1. Graph Traffic in the last days
2. Graph of number of customers connected in the last days
3. Graph Customers activated in the last days
4. Table of Currently connected Users (as current)
   * Login, name, last access, duration, IN / OUT traffic


### Social information

(**Partially implemented: data are stored but not displayed**)

Data obtained from Social profiles

1. Gender
2. Age
3. Nationality
4. Interests

### Questionnaires

(**Not yet implemented**)

Data obtained through questionnaires

1. Reason of the trip (business, fun, other)
2. Nationality
3. Length of stay
4. Interests
5. Request for feedback on the structure (any problems)
6. Communications via WhatsApp Business (future)


### Data export
(**Partially implemented**)

The customer is able to export data in CSV format to reuse them in mass marketing software (eg. MailChimp).

Exportable data should include:

* email 
* phone number
* username 

Other data (like gender, age...) we can receive from social login are still not exportable,the system will export ONLY data from users that allow marketing usage.


