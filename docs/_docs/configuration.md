---
title: "Configuration"
permalink: /docs/configuration/
---

# Sun UI
To configure the front-end part aka Sun-UI, you can edit the file in `static/config/config.js` and choose your app name. Copy also the logos and background in the `static` folder.

The configuration file it's easy to edit
```js
const CONFIG = {
  APP_NAME: '<YOUR_APP_NAME>',
}

```

# Wings & Wax
Wings component using Wax as backend supports many type of authentication, social or not. To make this fully functional you must configure some params for each type of authentication.

During the provisioning phase of your instance, or in a second step, you must edit the Wax configuration file to get all logins available.

The configuration file is `/opt/icaro/wax/conf.json` and the parameters to configure are the following.

## Database
Wax and Sun database configuration is saved inside the file `/opt/icaro/wax/conf.json` in the `database` part.
```json
"database": {
    "host":"localhost",
    "port":"3306",
    "name":"icaro",
    "user": "sun-api",
    "password": "Sun-ApiMariaDBPassWordHere"
}
```

## Disclaimers
### Default
Each Captive portal is served from the same cloud server, but each captive portal related to a particular hotspot inherits its configuration.

You must specify the business name during Hotspot creation that will be replaced in the disclaimers for each hotspots.

Terms of use and marketing disclaimers are visualized before the user chooses the login method. To add your disclaimers modify the `/opt/icaro/wax/conf.json` in the `disclaimers` and use

- `{% raw %}{{ .BusinessName }}{% endraw %}` for company name.
- `{% raw %}{{ .BusinessVAT }}{% endraw %}` for company VAT.
- `{% raw %}{{ .BusinessAddress }}{% endraw %}` for company address.
- `{% raw %}{{ .BusinessEmail }}{% endraw %}` for company email.
- `{% raw %}{{ .BusinessDPO }}{% endraw %}` for DPO name.
- `{% raw %}{{ .BusinessDPOMail }}{% endraw %}` for DPO mail.
- `{% raw %}{{ .IntegrationTerms }}{% endraw %}` for external integration privacies.

inside the disclaimers JSON object:

```json
"disclaimers": {
	"terms_of_use": "This is a disclaimer test\n\n - chapter 1\n - chapter 2 provided by {% raw %}{{ .BusinessName }}{% endraw %} located in {% raw %}{{ .BusinessAddress }}{% endraw %}",
	"marketing_use": "This is marketing informationt\n\n - chapter 1\n - chapter 2  provided by {% raw %}{{ .BusinessName }}{% endraw %} located in {% raw %}{{ .BusinessAddress }}{% endraw %}"
}
```

Became (if the hotspot's business name is `Great Hotel Ltd` located in `Street of Null` for example)
```
This is a disclaimer test\n\n - chapter 1\n - chapter 2 provided by Great Hotel Ltd located in Street of Null
```

### Custom
For each reseller it is possible to upload customized disclaimers. The `admin` user can select a particular reseller and upload disclaimers to his profile both for privacy and for terms of use with a specific attribute.

Once loaded into the profile, the reseller will be able to choose which disclaimer to use by default for all his installations and hotspots, by going to his profile and selecting the relative attribute.

For each hotspot it is possible to override the default disclaimer, choosing in detail the hotspot which one to use.

If the reseller has not uploaded any disclaimers, then the Icaro installation default is used for all hotspots.

**Tables**
Used to load all custom disclaimers
```sql
CREATE TABLE `disclaimers` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(250) NOT NULL,
  `title` varchar(250) NOT NULL,
  `body` text NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
)
```

Used to map custom disclaimers to hotspot, to handle override
```sql
CREATE TABLE `disclaimers_hotspots` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `disclaimer_id` bigint(20) unsigned NOT NULL,
  `hotspot_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `hotspot_id` (`hotspot_id`),
  KEY `disclaimer_id` (`disclaimer_id`,`hotspot_id`),
  CONSTRAINT `disclaimers_hotspots_ibfk_1` FOREIGN KEY (`disclaimer_id`) REFERENCES `disclaimers` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `disclaimers_hotspots_ibfk_2` FOREIGN KEY (`hotspot_id`) REFERENCES `hotspots` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
)
```

Used to map custom loaded disclaimers to a particular reseller
```sql
CREATE TABLE `disclaimers_accounts` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `disclaimer_id` bigint(20) unsigned NOT NULL,
  `account_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  KEY `account_id` (`account_id`),
  KEY `disclaimer_id` (`disclaimer_id`,`account_id`),
  CONSTRAINT `disclaimers_accounts_ibfk_1` FOREIGN KEY (`disclaimer_id`) REFERENCES `disclaimers` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `disclaimers_accounts_ibfk_2` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
)
```

## Social
For the social logins, each platforms as its application registration process, here the links to create new apps for:
- **Facebook**: [create new app](https://developers.facebook.com/apps)
- **LinkedIn**: [create new app](https://www.linkedin.com/developer/apps)

After the registration insert your variables inside the file `/opt/icaro/wax/conf.json` in the `auth_social` part.
```json
"auth_social": {
	"facebook": {
		"client_id": "fb_client_id",
		"client_secret": "fb_client_secret",
		"redirect_uri": "https://<icaro_instance>/wings/login/facebook"
	},
	"linkedin": {
		"client_id": "li_client_id",
		"client_secret": "li_client_secret",
		"redirect_uri": "http://<icaro_instance>/wings/login/linkedin"
	}
}
```

## Email & SMS
For the Email login process the only requested part is a valid `smtp` host to send verification emails, configure the `/opt/icaro/wax/conf.json` file in the `endpoints` and `email` part.

For the SMS login process Icaro use Twilio as SMS sender, register your app here: [Twilio Apps](https://www.twilio.com/console), and edit all required params in `/opt/icaro/wax/conf.json` file in the `endpoints` and `sms` part.

```json
"endpoints": {
	"sms": {
		"account_sid": "twilio_account_sid",
		"auth_token": "twilio_account_token",
		"service_sid": "twilio_messaging_service_sid",
		"send_quota_alert": true
	},
	"email": {
		"from": "email_from",
		"smtp_host": "your_smtp_host",
		"smtp_port": 25,
		"smtp_user": "your_smtp_user",
		"smtp_password": "your_smtp_password"
	}
}

```

If `send_quota_alert` is set to true, resellers will be notify when theirs SMS quota limit is reached. (default: false)

## Captive portal defaults
The captive portal can be customized for each Hotspot, but if the user doesn't have the rights, you can set captive portal default options under the the `captive_portal` section:
```JSON
"captive_portal": {
    "redirect": "https://nethesis.github.io/icaro",
    "title": "Icaro",
    "subtitle": "The Open Source Hotspot",
    "description": "Free as in freedom",
    "background": "#2a87be",
    "logo": "logo.png",
    "banner": "banner.png"
}
```

## Envirioment Variables

You can also configure some properties via environment variables, if present **the values in the environment variable take the precedence** over the ones declared in files.

### Common variables

This variables are common to all backends:

* `DB_HOST` Mysql or Mariadb hostname/address
* `DB_PORT` Port where database instance is listening
* `DB_USER` User to use for connection to database instance
* `DB_PASSWORD` Password of the user used for database instance connection
* `DB_NAME` Name of database used by Icaro's backends (default: icaro)

### Wax & Sun-Api variables

* `CORS_ORIGINS` List of space separated origins allowed to perform cross-site requests (see more [here](https://www.w3.org/TR/cors/#access-control-allow-origin-response-header))

### Wax specific variables

#### Social

##### Facebook

* `FACEBOOK_CLIENT_ID` Facebook Oauth2 Client ID
* `FACEBOOK_CLIENT_SECRET` Facebook Oauth2 Client Secret
* `FACEBOOK_REDIRECT_URL` Facebook Redirect URL

##### LinkedIn

* `LINKEDIN_CLIENT_ID` LinkedIn Oauth2 Client ID
* `LINKEDIN_CLIENT_SECRET` LinkedIn Oauth2 Client Secret
* `LINKEDIN_REDIRECT_URL` LinkedIn Oauth2 Redirect URL

#### Email & SMS

##### Email
* `EMAIL_FORM` Email from address
* `EMAIL_FORM_NAME` Email from name
* `EMAIL_SMTP_HOST` SMTP server address/hostanme
* `EMAIL_SMTP_PORT` Port where the SMTP server is listening
* `EMAIL_SMTP_USER` User to use for connection to SMTP server
* `EMAIL_SMTP_PASSWORD` Password of the user used for SMTP server connection

##### SMS
* `SMS_ACCOUNT_SID` Twilio account SID
* `SMS_AUTH_TOKEN` Twilio auth token
* `SMS_SERVICE_SID` Twilio messaging service sid
* `SMS_SEND_QUOTA_ALERT` `true | false` Enable/Disable SMS quota limit alert.

### Ade specific variables

##### Survey
* `SURVEY_URL` Url to host the Ade UI component
