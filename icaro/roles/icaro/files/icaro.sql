/* ACCOUNTING AAA */
CREATE TABLE hotspots (
  id serial,
  name varchar(200),
  description varchar(250),
  created datetime,
  KEY(name),
  PRIMARY KEY(id)
);

CREATE TABLE hotspot_preferences (
  id serial,
  hotspot_id bigint,
  key varchar(250),
  value varchar(250),
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  UNIQUE KEY (hotspot_id, key),
  PRIMARY KEY(id)
);

CREATE TABLE users (
  id serial,
  hotspot_id bigint unsigned,
  username varchar(200),
  password varchar(200),
  email varchar(200),
  account_type varchar(200),
  kbps_down integer unsigned,
  kbps_up integer unsigned,
  valid_from datetime,
  valid_until datetime,
  created datetime,
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  UNIQUE KEY (hotspot_id, username),
  KEY(created),
  KEY(username),
  KEY(email),
  PRIMARY KEY(id)
);

CREATE TABLE devices (
  id serial,
  hotspot_id bigint unsigned,
  user_id bigint unsigned,
  mac_address varchar(200),
  description varchar(200),
  created datetime,
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  KEY(mac_address),
  PRIMARY KEY(id)
);

CREATE TABLE units (
  id serial,
  hotspot_id bigint unsigned,
  mac_address varchar(200),
  description varchar(200),
  uuid varchar(200),
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  KEY(mac_address),
  KEY(uuid),
  PRIMARY KEY(id)
);

CREATE TABLE sessions (
  id serial,
  unit_id bigint unsigned,
  hotspot_id bigint unsigned,
  device_id bigint unsigned,
  user_id bigint unsigned,
  bytes_up bigint unsigned,      -- bytes uploaded by user
  bytes_down bigint unsigned,    -- bytes downloaded by user
  duration bigint unsigned,      -- duration in seconds
  auth_time datetime,            -- set to now() at authentication
  start_time datetime,           -- set to now() on accounting start
  update_time datetime,          -- set to now() on accounting start,update,stop
  stop_time datetime,            -- set to now() on accounting stop
  session_key varchar(200),      -- a unique key generated from session data
  FOREIGN KEY (unit_id) REFERENCES units(id),
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  FOREIGN KEY (device_id) REFERENCES devices(id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  KEY(session_key),
  KEY(bytes_up),
  KEY(bytes_down),
  KEY(duration),
  KEY(auth_time),
  KEY(start_time),
  KEY(update_time),
  KEY(stop_time),
  PRIMARY KEY(id)
);
/* -------------------- */

/* HOTSPOT MANAGER */
CREATE TABLE resellers (
  id serial,
  name varchar(200),
  email varchar(250),
  username varchar (200),
  password varchar (200),
  created datetime,
  uuid varchar(200),
  KEY(username),
  KEY(uuid),
  PRIMARY KEY(id)
);

CREATE TABLE reseller_preferences (
  id serial,
  reseller_id bigint,
  key varchar(250),
  value varchar(250),
  FOREIGN KEY (reseller_id) REFERENCES resellers(id),
  UNIQUE KEY (reseller_id, key),
  PRIMARY KEY(id)
);

CREATE TABLE hotspot_accounts (
  id serial,
  email varchar(250),
  username varchar (200),
  password varchar (200),
  created datetime,
  type varchar(200),
  KEY(username),
  PRIMARY KEY(id)
)
/* --------------- */