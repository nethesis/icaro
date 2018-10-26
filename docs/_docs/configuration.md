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
Each Captive portal is served from the same cloud server, but each captive portal related to a particular hotspot inherits its configuration.

You must specify the business name during Hotspot creation that will be replaced in the disclaimers for each hotspots.

Terms of use and marketing disclaimers are visualized before the user chooses the login method. To add your disclaimers modify the `/opt/icaro/wax/conf.json` in the `disclaimers` and use

- `$$COMPANY_NAME$$` for company name.
- `$$COMPANY_VAT$$` for company VAT.
- `$$COMPANY_ADDRESS$$` for company address.
- `$$COMPANY_EMAIL$$` for company email.

inside the disclaimers JSON object:

```json
"disclaimers": {
	"terms_of_use": "This is a disclaimer test\n\n - chapter 1\n - chapter 2 provided by $$BUSINESS_NAME$$ located in $$COMPANY_ADDRESS$$",
	"marketing_use": "This is marketing informationt\n\n - chapter 1\n - chapter 2  provided by $$BUSINESS_NAME$$  located in $$COMPANY_ADDRESS$$"
}
```

Became (if the hotspot's business name is `Great Hotel Ltd` located in `Street of Null` for example)
```
This is a disclaimer test\n\n - chapter 1\n - chapter 2 provided by Great Hotel Ltd located in Street of Null
```

## Social
For the social logins, each platforms as its application registration process, here the links to create new apps for:
- **Facebook**: [create new app](https://developers.facebook.com/apps)
- **LinkedIn**: [create new app](https://www.linkedin.com/developer/apps)
- **Instagram**: [create new app](https://www.instagram.com/developer/)

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
	},
	"instagram": {
		"client_id": "in_client_id",
		"client_secret": "in_client_secret",
		"redirect_uri": "http://<icaro_instance>/wings/login/instagram"
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

##### Instagram

* `INSTAGRAM_CLIENT_ID` Instagram Oauth2 Client ID
* `INSTAGRAM_CLIENT_SECRET` Instagram Oauth2 Client Secret
* `INSTAGRAM_REDIRECT_URL` Instagram Oauth2 Redirect URL

#### Email & SMS

##### Email
* `EMAIL_FORM` Email from address
* `EMAIL_SMTP_HOST` SMTP server address/hostanme
* `EMAIL_SMTP_PORT` Port where the SMTP server is listening
* `EMAIL_SMTP_USER` User to use for connection to SMTP server
* `EMAIL_SMTP_PASSWORD` Password of the user used for SMTP server connection

##### SMS
* `SMS_ACCOUNT_SID` Twilio account SID
* `SMS_AUTH_TOKEN` Twilio auth token
* `SMS_SERVICE_SID` Twilio messaging service sid
* `SMS_SEND_QUOTA_ALERT` `true | false` Enable/Disable SMS quota limit alert.
