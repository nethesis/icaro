---
title: Marketing
permalink: /docs/marketing/
---

Icaro collects some data about the authenticated users to make the logins option available, for example each login options (Facebook, LinkedIn ...) collects informations that are stored in plain mode as a JSON object in the database. All of this data are usefull to create marketing campains, inspect the users likes, user working positions or other social activities.

All of this information are accepted by the users when the social login prompts ask them to authorize the app to get free navigation.

The table involved is `user_marketings`:
```
MariaDB [icaro]> describe user_marketings;
+--------------+---------------------+------+-----+---------+----------------+
| Field        | Type                | Null | Key | Default | Extra          |
+--------------+---------------------+------+-----+---------+----------------+
| id           | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |
| user_id      | bigint(20) unsigned | NO   | UNI | NULL    |                |
| account_type | varchar(200)        | NO   |     | NULL    |                |
| data         | longtext            | NO   |     | NULL    |                |
| created      | datetime            | YES  |     | NULL    |                |
+--------------+---------------------+------+-----+---------+----------------+
5 rows in set (0.00 sec)
```

The `data` field contains all login information provided by the social auth mode.

The main goal of Icaro is to provide an easy and simple Hotspot software, the marketing part is secondary, but Icaro provides information that can be mined with third party softwares.