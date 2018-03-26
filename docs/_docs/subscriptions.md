---
title: "Subscriptions"
permalink: /docs/subscriptions/
---

Accounts should be profiled and have access to different features based on their own entitlements.

## Rules

- Each account is associated to a subscription
- A subscription is valid for a limited time and it's associated to a subscription plan
- Each subscription plan has a cost of renewal and a duration
- When a subscription is renewed, the validity will be extended of X days, where X is the duration of the associated subscription plan
- Each subscription plan enforces a set of rules to the associated account. Each rule give access to specific feature
- Subscription plan rules are applied to every sub-account (customer, desk) associated with the reseller account

## Plans

Default plans:
  - Every new account will have `premium` subscription plan valid for `365` day (1 year)
  - Every subscription plan has the `price` set to 0 and the `duration` set to 365 days.


|                | Free                   | Basic           | Standard  | Premium |
| --------------------- |:----------------:| -------------:| ------------:|------------:|
| Included SMS                  | 0 | 500 | 1000 | 2000 |
| No. Units/Hotspot             | 1 | 1 | 10 |100 |
| Advanced reports              | No | Yes | Yes | Yes |
| Captive portal customization  | No | No | Yes | Yes |
| Social marketing              | No | No | No |  Yes |

Rules are saved inside `subscription_plans` table as extra fields.



### Included SMS

Number of SMS included for each plan, when the threshold is reached, the system will no longer send SMS. SMS can be "bought" as packs of X elements.

Field: `included_sms` , should be evaluated along with `sms_max_count` inside `account_sms_counts` table.

### Number of units for Hotspot

The number of Units which can be connected to an Hotspot instance.

Field: `max_units`

### Advanced reports

Basic reports include:

- list of active sessions
- list of inactive sessions

Advanced reports include:

- Sessions/ time graph
- New registrations/time graph
- Bandwidth/time graph total
- Bandwidth/time graph on single user
- Average duration of sessions for timeframes

Field: `advanced_report`

### Captive portal customization

Allow changing captive portal aspect like logos and colors.

Field: `wings_customization`

### Social marketing and analytics

**NOT implemented**

Reports created on social collected data, like  customer age and country distribution.
- login type graph (fb, insta, link, sms, email)
- users genre
- day based histogram male/female
- day based histogram age
- day based nationality
- filters (genre/age/nationality/time period) on users (export and/or integration with mailchimp like services)

Field: `social_analytics`


## Enforced rules

- Authentication is denied for all clients connecting to an Hotspot belonging to an expired account
- Resellers and customers wit "Free" and "Basic" subscription can't change Captive Portal view configuration
- One a unit is registered, the system checks the number of allowed units per hotspot
- During reseller account creation the administrator must choose a subscription plan


### Not implemented yet

- An expired (reseller, desk or customer) account has access to the Hotspot Manager but is locked in a page until the expire date is extended
- New accounts have a `free` subscription plan by defeulat
- Each reseller account can upgrade his own subscription plan paying with a credit card
- Admin user can change subscription plan for existing accounts
