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
  UNIQUE KEY (`username`),
  UNIQUE KEY (`uuid`),
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
  `value` mediumtext NOT NULL,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`hotspot_id`, `key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `hotspot_vouchers` (
  `id` serial,
  `hotspot_id` bigint unsigned NOT NULL,
  `code` varchar(250) NOT NULL,
  `auto_login` tinyint NOT NULL,
  `bandwidth_up` integer unsigned,
  `bandwidth_down` integer unsigned,
  `duration` integer unsigned,
  `remain_use` integer,
  `expires` datetime DEFAULT NULL,
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
  `auto_login` tinyint NOT NULL,
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

CREATE TABLE `user_marketings` (
  `id` serial,
  `user_id` bigint unsigned NOT NULL,
  `account_type` varchar(200) NOT NULL,
  `data` longtext NOT NULL,
  `created` datetime,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`user_id`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `user_sessions` (
  `id` serial,
  `user_id` bigint unsigned NOT NULL,
  `session_key` varchar(200),
  `created` datetime,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`user_id`, `session_key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `user_temp_sessions` (
  `id` serial,
  `user_id` bigint unsigned NOT NULL,
  `device_mac` varchar(200) NOT NULL,
  `session_key` varchar(200),
  `created` datetime,
  FOREIGN KEY (`user_id`) REFERENCES users(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`user_id`, `session_key`),
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
  `name` varchar(200),
  `description` varchar(200),
  `uuid` varchar(200) NOT NULL,
  `secret` varchar(200) NOT NULL,
  `created` datetime,
  FOREIGN KEY (`hotspot_id`) REFERENCES hotspots(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`uuid`),
  KEY(`mac_address`),
  KEY(`uuid`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `sessions` (
  `id` serial,
  `unit_id` bigint unsigned NOT NULL,
  `unit_mac` varchar(200) DEFAULT NULL,
  `hotspot_id` bigint unsigned NOT NULL,
  `hotspot_desc` varchar(1024) DEFAULT NULL,
  `device_id` bigint unsigned NOT NULL,
  `device_mac` varchar(200) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `username` varchar(500) NOT NULL,
  `bytes_up` bigint unsigned,
  `bytes_down` bigint unsigned,
  `duration` bigint unsigned,
  `auth_time` datetime,
  `start_time` datetime,
  `update_time` datetime,
  `stop_time` datetime,
  `session_key` varchar(200),
  UNIQUE KEY `session_key` (`session_key`),
  PRIMARY KEY(`id`)
);

CREATE TABLE `session_histories` (
  `id` serial,
  `session_id` bigint unsigned NOT NULL,
  `unit_id` bigint unsigned NOT NULL,
  `unit_mac` varchar(200) DEFAULT NULL,
  `hotspot_id` bigint unsigned NOT NULL,
  `hotspot_desc` varchar(1024) DEFAULT NULL,
  `device_id` bigint unsigned NOT NULL,
  `device_mac` varchar(200) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `username` varchar(500) NOT NULL,
  `bytes_up` bigint unsigned,
  `bytes_down` bigint unsigned,
  `duration` bigint unsigned,
  `auth_time` datetime,
  `start_time` datetime,
  `update_time` datetime,
  `stop_time` datetime,
  `session_key` varchar(200),
  PRIMARY KEY(`id`)
);

/* -------------------- */

/* SUBSCRIPTIONS */
CREATE TABLE subscription_plans (
    id serial not null primary key,
    code varchar(200) not null,
    name varchar(200) not null,
    description varchar(200) not null,
    price decimal(5,2),
    period integer default null,
    included_sms integer not null,
    max_units integer not null,
    advanced_report boolean default false,
    wings_customization boolean default false,
    social_analytics boolean default false
);

INSERT INTO subscription_plans VALUES (1, 'free', 'Free', 'Free limited plan', 0.00, 365, 0, 1, false, false, false);
INSERT INTO subscription_plans VALUES (2, 'basic', 'Basic', 'Basic plan', 0.00, 365, 500, 1, true, false, false);
INSERT INTO subscription_plans VALUES (3, 'standard', 'Standard', 'Standard lan', 0.00, 365, 1000, 10, true, true, false);
INSERT INTO subscription_plans VALUES (4, 'premium', 'Premium', 'Premium plan', 0.00, 3650, 2000, 100, true, true, true);

CREATE TABLE subscriptions (
    id serial not null primary key,
    account_id bigint unsigned not null references accounts(id),
    subscription_plan_id bigint not null references subscription_plans(id),
    valid_from timestamp null,
    valid_until timestamp null,
    created timestamp default current_timestamp
);

/* -------------------- */

/* EXTRAS */
CREATE TABLE `account_sms_counts` (
  `id` serial,
  `account_id` bigint unsigned NOT NULL,
  `sms_max_count` bigint unsigned,
  `sms_count` bigint unsigned,
  FOREIGN KEY (`account_id`) REFERENCES accounts(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  UNIQUE KEY (`account_id`),
  PRIMARY KEY(`id`)
);
/* ------ */
