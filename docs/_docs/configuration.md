---
title: "Configuration"
permalink: /docs/configuration/
---

# Wings & Wax
Wings component using Wax as backend supports many type of authentication, social or not. To make this fully functional you must configure some params for each type of authentication.

During the provisioning phase of your instance, or in a second step, you must edit the Wax configuration file to get all logins available.

The configuration file is `/opt/icaro/wax/conf.json` and the parameters to configure are:

## Social
For the social logins, each platforms as its application registration process, here the links to create new apps for:
- **Facebook**: [create new app](https://developers.facebook.com/apps)
- **Google+**: [create new app](https://console.cloud.google.com/home/dashboard)
- **LinkedIn**: [create new app](https://www.linkedin.com/developer/apps)
- **Instagram**: [create new app](https://www.instagram.com/developer/)

After the registration insert your variables inside the file `/opt/icaro/wax/conf.json` in the `auth_social` part.
```json
"auth_social": {
	"facebook": {
		"client_id": "fb_client_id",
		"client_secret": "fb_client_secret",
		"redirect_uri": "http://<icaro_instance>/wings/login/facebook"
	},
	"google": {
		"client_id": "gl_client_id",
		"client_secret": "gl_client_secret",
		"redirect_uri": "http://<icaro_instance>/wings/login/google"
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
		"number": "send_number"
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