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

Survey
------
Icaro can send surveys to the users about feedbacks and reviews. Normally feedbacks are sent during the stay after some hours (customizable) and reviews after the stay, after hours or days (also this parameters are configurable).

The authorization to send surveys is handled during the login, by checking a special field that marks the field `survey_auth` to true.

Each hour a cron checks for all users if the `survey_auth` field is true and sends surveys, if not already sent, based on configurations in the `Marketing` section on Sun-UI (hotspot manager).

The survey is sent by `email` or by `sms` at the user endpoint.

Each survey is one-shot, and the token is handled by the `ade` components using the `ade_tokens` table.

```
mysql> describe ade_tokens;
+--------------------+---------------------+------+-----+---------+----------------+
| Field              | Type                | Null | Key | Default | Extra          |
+--------------------+---------------------+------+-----+---------+----------------+
| id                 | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |
| token              | varchar(250)        | NO   | MUL | NULL    |                |
| feedback_sent_time | datetime            | NO   |     | NULL    |                |
| feedback_left_time | datetime            | YES  |     | NULL    |                |
| review_sent_time   | datetime            | NO   |     | NULL    |                |
| review_left_time   | datetime            | YES  |     | NULL    |                |
| user_id            | bigint(20) unsigned | NO   |     | NULL    |                |
| hotspot_id         | bigint(20) unsigned | NO   | MUL | NULL    |                |
+--------------------+---------------------+------+-----+---------+----------------+
8 rows in set (0.00 sec)
```

When a feedback or review is left, then the token for that user is invalidated and an email is sent to the owner of the hotspot with the feedback or review text.