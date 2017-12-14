# Hotspot Manager

The **Hotspot Manager** is  web interface used to configure and manage the Hotspot.

## Roles

There are 3 different profiles that can access the Hotspot Manager, each profile different authorizations
and can access only sections related to its own role.

The 3 profiles are the following:

1. Reseller
2. Customer
3. Desk

Desk can see very few things, Customer can see all that Desk see plus more panels, Reseller can see everything.

### Reseller

The vendor which installs the physical firewall and creates an Hotspot entity; he can create users for his own Hotspot instances. 
This is the only profile that can access creation API (or Nethesis portal)
The Reseller can see everything on the hotspot manager, he makes the first configuration of the Hotspot, choosing the best options for the specific customer.

**Every technical or potentially dangerous configuration must be performed only by the reseller.**  
The reseller creates its own account for the specific Hotspot using the APIs (or Nethesis protal),
after that it will login on the Hotppot Manager and ceate login credentials for other profiles.
Every Hotspot can have only 3 different profiles/accounts.
 
### Customer

The owner of the hotel or the local person responsible for the service, he doesn't have technical skills but he needs to do some tasks:

* customize the captive portal
* customize Ads in the splash page
* check how many users are connected
* read users statistics
* export reports
* make marketing activities

### Desk

The is the person at the desk of the hotel, Desk can:

* create accounts
* see the list of accounts
* see how many accounts are connected and the bandwidth consumption (e.g. slow traffic)

### Hotspot manager views

List of the views of the control panel with the associated profile (**R**eseller, **C**ustomer, **D**esk)

* DashBoard (**C,D**)
* Account (creation and management)(**C,D**)
* Report and Log (**C**)
* Marketing (**C**)
* Preferences (include captive portal customization) (**C**)
* Admin Prefs (**R**eseller only)


## Special Hotspot Manager access

Every Nethesis reseller that can access Nethesis portal should make access with the same reseller credentials to the Hostpot Manager,
to see and manage all the reseller's Hotspot instances.

Thus Hotspot Manager should support OAuth authentication and Nethesis portal should become a OAuth provider.
