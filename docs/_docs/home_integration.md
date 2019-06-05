---
title: Introduction
permalink: /docs/home_integration/
---

## General
Icaro functionalities can be extended using external integrations, to add new features or demand particulary operations in different platforms.

Each integrations can be added in the Icaro instance and can be activated or deactived for each hotspots. When an integrations has been activated a payload is generated and sended to the external platform.

The integration works in 2 different way:

- The external platform receives the payload and uses that data in different ways, activated by **hotspot manager**. (`Sun-API`)
- During the authentication process a `URL` is called to pass data in `pre` or `post` authentication step, in the **captive portal**. (`Wings`)

# Integration parameter
Each record in the database of an integration has the following fields
```sql
CREATE TABLE `integrations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(250) NOT NULL,
  `description` varchar(250) NOT NULL,
  `site` varchar(250) NOT NULL,
  `logo` mediumtext NOT NULL,
  `webhook_url` varchar(250) NOT NULL,
  `webhook_token` varchar(250) NOT NULL,
  `pre_auth_redirect_url` varchar(250) NOT NULL,
  `post_auth_redirect_url` varchar(250) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
)
```
- `name`, `description`, `site`, `logo`: metadata of integration
- `webhook_url` and `webhook_token`: url and authentication token of the external application call by `Sun-API` during the integration.
- `pre_auth_redirect_url` and `post_auth_redirect_url`: external url called during authentication process in the `Wings` component.

# Hotspot manager (`Sun-API`)
When an integration is enabled a payload like in the example below is created and sended to the `webhook_url` using the `webhook_token`. The `webhook_token` is added in the request with `X-Icaro-WebHook-Token` header.

Example:
```bash
POST -> webhook_url with { Headers: X-Icaro-WebHook-Token: 060a5ad5108fada38d10ba54756}
```
```json
{
  "status": true,
  "token": "5ca6448cf3475717799718e90060a5ad5108fada38d10ba547564fcd348111a9",
  "hotspot": {
    "id": 1,
    "uuid": "8cad0e8b-892a-4d40-9b07-044702ef5a45",
    "name": "test",
    "description": "test",
    "business_name": "test",
    "business_vat": "test",
    "business_address": "test",
    "business_email": "test@test.it"
  },
  "accounts": [
    {
      "id": 3,
      "uuid": "eb3c1fd3-80c0-4177-a7e5-e804f0322edf",
      "type": "customer",
      "name": "c1",
      "username": "c1",
      "password_hash": "3c27daec33252e33039c743617bd39b2",
      "email": "c1@c1.it"
    },
    {
      "id": 4,
      "uuid": "947e4916-391c-4e16-b759-949cca96043b",
      "type": "desk",
      "name": "d1",
      "username": "d1",
      "password_hash": "3c27daec33252e33039c743617bd39b2",
      "email": "d1@d1.it"
    }
  ],
  "units": [
    {
      "id": 1,
      "mac_address": "E4-95-6E-44-1D-5C",
      "name": "test",
      "description": "unit_description test",
      "uuid": "8b758c35-7fb4-48a3-9021-df7a51a82c90"
    }
  ]
}
```
Fields explanation:

- `status`: tells if the service is `activated` or `deactivated` from user. When a `reseller` wants to remove the integration the value will be `false`
- `token`: the token used to authenticate the requests from the external platform to Icaro
- `hotspot`: is the hotspot where the integration is enabled
- `accounts`: list of accounts of the hotspot with relative credentials
- `units`: list of the units of the hotspot

# Captive portal (`Wings`)
In the authentication process there are 2 different step:
- the `pre_auth_` step, the moment before choosing an authentication method
- the `post_auth` step, the moment after the authentication method is choosen and the user is created, but **not** authenticated

When one of this steps is reached, the captive portal check if the relative hotspot has the integration enabled and the `pre_auth_redirect_url` or `post_auth_redirect_url` fields are defined.

If one of this field is defined, the captive portal creates the `GET` request and append it to the `pre` or `post` url, after the redirect a callback is called in the integration page and the user returns in the captive portal.

Use case.
1. The `post_auth_redirect_url` is `https://my.integration.org`
2. The captive portal creates the `GET` request getting all the parameters in the current session after a login authentication method is choosen
3. The captive portal open a new page in `https://my.integration.org/wings/login/sms?digest=...`
4. User navigate in the external page and confirm the process
5. The external integration calls a defined *callback* to redirect user to the captive portal
6. The user is now **authenticated**