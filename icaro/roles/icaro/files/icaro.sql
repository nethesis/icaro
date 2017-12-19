/* HOTSPOT MANAGER */
CREATE TABLE `accounts` (
  `id` serial,
  `uuid` varchar(200),
  `type` varchar(200),
  `name` varchar(200),
  `username` varchar (200),
  `password` varchar (200),
  `email` varchar(250),
  `created` datetime,
  KEY(`username`),
  KEY(`uuid`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `account_preferences` (
  `id` serial,
  `account_id` bigint unsigned,
  `key` varchar(250),
  `value` varchar(250),
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`),
  UNIQUE KEY (`account_id`, `key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `access_tokens` (
  `id` serial,
  `account_id` bigint unsigned,
  `token` varchar(200),
  `role` varchar(200),
  `expires` datetime,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`),
  UNIQUE KEY (`account_id`, `id`),
  PRIMARY KEY(`id`)
);
/* --------------- */

/* ACCOUNTING AAA */
CREATE TABLE `hotspots` (
  `id` serial,
  `account_id` bigint unsigned,
  `name` varchar(200),
  `description` varchar(250),
  `created` datetime,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`),
  KEY(`name`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `hotspot_preferences` (
  `id` serial,
  `hotspot_id` bigint unsigned,
  `key` varchar(250),
  `value` varchar(250),
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`),
  UNIQUE KEY (`hotspot_id`, `key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `users` (
  `id` serial,
  `hotspot_id` bigint unsigned,
  `name` varchar(200),
  `username` varchar(200),
  `password` varchar(200),
  `email` varchar(200),
  `account_type` varchar(200),
  `kbps_down` integer unsigned,
  `kbps_up` integer unsigned,
  `valid_from` datetime,
  `valid_until` datetime,
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`),
  UNIQUE KEY (`hotspot_id`, `username`),
  KEY(`created`),
  KEY(`username`),
  KEY(`email`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `devices` (
  `id` serial,
  `hotspot_id` bigint unsigned,
  `user_id` bigint unsigned,
  `mac_address` varchar(200),
  `description` varchar(200),
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  KEY(`mac_address`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `units` (
  `id` serial,
  `hotspot_id` bigint unsigned,
  `mac_address` varchar(200),
  `description` varchar(200),
  `uuid` varchar(200),
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`),
  KEY(`mac_address`),
  KEY(`uuid`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `sessions` (
  `id` serial,
  `unit_id` bigint unsigned,
  `hotspot_id` bigint unsigned,
  `device_id` bigint unsigned,
  `user_id` bigint unsigned,
  `bytes_up` bigint unsigned,
  `bytes_down` bigint unsigned,
  `duration` bigint unsigned,
  `auth_time` datetime,
  `start_time` datetime,
  `update_time` datetime,
  `stop_time` datetime,
  `session_key` varchar(200),
  FOREIGN KEY (`unit_id`) REFERENCES units(`id`),
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`),
  FOREIGN KEY (`device_id`) REFERENCES devices(`id`),
  FOREIGN KEY (`user_id`) REFERENCES users(`id`),
  PRIMARY KEY(`id`)
);
/* -------------------- */