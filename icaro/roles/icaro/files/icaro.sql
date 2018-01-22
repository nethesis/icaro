/* HOTSPOT MANAGER */
CREATE TABLE `accounts` (
  `id` serial,
  `creator_id` bigint unsigned NOT NULL,
  `uuid` varchar(200),
  `type` varchar(200) NOT NULL,
  `name` varchar(200) NOT NULL,
  `username` varchar (200) NOT NULL,
  `password` varchar (200) NOT NULL,
  `email` varchar(250),
  `created` datetime,
  KEY(`username`),
  KEY(`uuid`),
  PRIMARY KEY(`id`)
);

/* CREATE DEFAULT ADMIN USER */
INSERT INTO `accounts` VALUES (1, 0, "", "admin", "Admin", "admin", MD5("admin"), "", NOW());

CREATE TABLE `account_preferences` (
  `id` serial,
  `account_id` bigint unsigned NOT NULL,
  `key` varchar(250) NOT NULL,
  `value` varchar(250) NOT NULL,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`account_id`, `key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `access_tokens` (
  `id` serial,
  `account_id` bigint unsigned NOT NULL,
  `token` varchar(200) NOT NULL,
  `role` varchar(200) NOT NULL,
  `expires` datetime NOT NULL,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`account_id`, `id`),
  PRIMARY KEY(`id`)
);
/* --------------- */

/* ACCOUNTING AAA */
CREATE TABLE `hotspots` (
  `id` serial,
  `account_id` bigint unsigned NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` varchar(250),
  `created` datetime,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  KEY(`name`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `accounts_hotspots` (
  `id` serial,
  `account_id` bigint unsigned NOT NULL,
  `hotspot_id` bigint unsigned NOT NULL,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  KEY(`account_id`, `hotspot_id`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `hotspot_preferences` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `key` varchar(250) NOT NULL,
  `value` varchar(16384) NOT NULL,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`hotspot_id`, `key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `hotspot_vouchers` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `code` varchar(250) NOT NULL,
  `expires` datetime NOT NULL,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`hotspot_id`, `code`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `users` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `name` varchar(200) NOT NULL,
  `username` varchar(200) NOT NULL,
  `password` varchar(200) NOT NULL,
  `email` varchar(200),
  `account_type` varchar(200) NOT NULL,
  `kbps_down` integer unsigned,
  `kbps_up` integer unsigned,
  `valid_from` datetime,
  `valid_until` datetime,
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`hotspot_id`, `username`),
  KEY(`created`),
  KEY(`username`),
  KEY(`email`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `user_sessions` (
  `id` serial,
  `user_id` bigint unsigned NOT NULL,
  `session_key` varchar(200),
  `created` datetime,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  PRIMARY KEY(`id`)
);

CREATE TABLE `devices` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `mac_address` varchar(200) NOT NULL,
  `ip_address` varchar(200),
  `description` varchar(200),
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  KEY(`mac_address`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `units` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `mac_address` varchar(200) NOT NULL,
  `description` varchar(200),
  `uuid` varchar(200) NOT NULL,
  `secret` varchar(200) NOT NULL,
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  KEY(`mac_address`),
  KEY(`uuid`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `sessions` (
  `id` serial,
  `unit_id` bigint unsigned NOT NULL,
  `hotspot_id` bigint unsigned NOT NULL,
  `device_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `bytes_up` bigint unsigned,
  `bytes_down` bigint unsigned,
  `duration` bigint unsigned,
  `auth_time` datetime,
  `start_time` datetime,
  `update_time` datetime,
  `stop_time` datetime,
  `session_key` varchar(200),
  FOREIGN KEY (`unit_id`) REFERENCES units(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (`device_id`) REFERENCES devices(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  PRIMARY KEY(`id`)
);
/* -------------------- */
