---
title: "Software components"
permalink: /docs/components/
---

Icaro is divided in 4 parts:

- [Dedalo](https://github.com/nethesis/icaro/tree/master/dedalo): network access controller,
  it runs on the firewall and intercepts all guest connections, based on [CoovaChilli](http://coova.github.io/CoovaChilli/)
- [Wax](https://github.com/nethesis/icaro/tree/master/wax): accounting server, it speaks RADIUS AAA protocol via HTTP with Dedalo
- [Wings](https://github.com/nethesis/icaro/tree/master/wings): captive portal hosted along with Wax, it's static responsive web page for user login. It talks with Wax and Dedalo
- Sun: Hotspot Manager applications which is dived in 3 parts
  - [Sun UI](https://github.com/nethesis/icaro/tree/master/sun/sun-ui): web interface written in [Vue.js](https://vuejs.org/) respecting [Patternfly](http://www.patternfly.org) design
  - [Sun API](https://github.com/nethesis/icaro/tree/master/sun/sun-api): [APIs]({{site.api_url}}) used by Sun UI, implemented using [Golang ](https://golang.org/)
  - [Sun Tasks](https://github.com/nethesis/icaro/tree/master/sun/sun-tasks): Golang applications to terminated sessions and cleanup invalid tokens
  - [Ade UI](https://github.com/nethesis/icaro/tree/master/ade/ade-ui): web interface written with VueJs that handles the feedbacks and reviews using tokens for each user.
  - [Ade API](https://github.com/nethesis/icaro/tree/master/ade/ade-api): used by Ade UI, implemented in [Golang ](https://golang.org/), handles the requests and send email to hotspot's owner
  - [Ade tasks](https://github.com/nethesis/icaro/tree/master/ade/ade-tasks): a cronjob that check if users need to receive surveys (feedbacks or reviews)

Wax, Sun and Ade use the same MariaDB (or MySQL) database.

Wax, Wings, Ade and Sun are considered server-side componentes, while Dedalo is considered the client one.
All techincal documentation can be found in the links above.

## How it works

This is an high-level overview of the whole architecture:

![schema](../img/schema.png "Schema")


### Interactions

- An user tries to access Internet using its own smartphone via WiFi

- The traffic is intercepted by Dedalo (**1** in figure)

- Dedalo talks to Wax and checks if the user is already authorized to access Internet.
  If not, the smartphone client is redirect to Wings captive portal (**2** in figure)

- The user is forced to login using one of the available methods.

- When the authentication is done, Wings sends an authorization message to Dedalo using [JSON interface](http://coova.github.io/CoovaChilli/JSON/) (**3** in figure)

- Dedalo asks authorization confirmation to Wax, if the authorization is granted (**4** in figure)

### Secondary interactions
- After hours or days, if the hotspot's owners have enabled the surveys, an email or sms are sent to the users
- The users click on the link inside the email or sms and fill the feedbacks or reviews
- When survey has finished a recap email is sent to the hotspot owner