---
title: Tokens
permalink: /docs/api_token/
---

## API tokens
If you want to grant access to the Icaro platform to external applications or users, you can create an access token with particulary **ACL** and **Role**.

The `access_tokens` table is:
```sql
CREATE TABLE `access_tokens` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `account_id` bigint(20) unsigned NOT NULL,
  `token` varchar(200) NOT NULL,
  `role` varchar(200) NOT NULL,
  `type` varchar(150) DEFAULT 'login',
  `expires` datetime NOT NULL,
  `acls` varchar(150) DEFAULT 'read',
  `description` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`),
  UNIQUE KEY `account_id` (`account_id`,`id`),
  CONSTRAINT `access_tokens_ibfk_1` FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
)
```

- `account_id` is referenced to a particular account: `admin`, `reseller` etc...
- `token` is the value of the token, normally a `SHA-256` string
- `role`: must be of of these: `admin`, `reseller`, `customer`, `desk`
- `type`: for users that can access to hotspot manager the type is `login`, otherwise is `api`
- `expires`: indicated the expiration of the token (only for `login` type)
- `acls`: must be `read` (only for `GET` requests), `write` (only for `GET`, `POST`, `PUT`), `full` (for all requests)
- `description`: describe for what the token is used

```
!!! At the moment to grant external access, the token must be inserted via SQL query from `admin` account in the `access_tokens` table.
```