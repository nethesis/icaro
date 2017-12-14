CREATE TABLE hotspots (
  id serial,
  name varchar(200),
  -- defaults for 'classes' of service
  defcode_redirection_url varchar(200),
  defcode_idle_timeout integer unsigned,
  defcode_kbps_down integer unsigned,
  defcode_kbps_up integer unsigned,
  defcode_limit_interval varchar(10), -- values: minute, hour, day, week, month
  defcode_limit_value integer unsigned, -- times limit_interval
  defcode_repeat_interval boolean default false,
  defcode_session_time integer unsigned,
  defcode_kbytes_total integer unsigned,
  defcode_kbytes_down integer unsigned,
  defcode_kbytes_up integer unsigned,
  defuser_redirection_url varchar(200),
  defuser_idle_timeout integer unsigned,
  defuser_kbps_down integer unsigned,
  defuser_kbps_up integer unsigned,
  defuser_limit_interval varchar(10), -- values: minute, hour, day, week, month
  defuser_limit_value integer unsigned, -- times limit_interval
  defuser_repeat_interval boolean default false,
  defuser_session_time integer unsigned,
  defuser_kbytes_total integer unsigned,
  defuser_kbytes_down integer unsigned,
  defuser_kbytes_up integer unsigned,
  defdev_always_allow boolean default false,
  defdev_redirection_url varchar(200),
  defdev_idle_timeout integer unsigned,
  defdev_kbps_down integer unsigned,
  defdev_kbps_up integer unsigned,
  defdev_limit_interval varchar(10), -- values: minute, hour, day, week, month
  defdev_limit_value integer unsigned, -- times limit_interval
  defdev_repeat_interval boolean default false,
  defdev_session_time integer unsigned,
  defdev_kbytes_total integer unsigned,
  defdev_kbytes_down integer unsigned,
  defdev_kbytes_up integer unsigned,
  uamsecret varchar(200),
  KEY(name),
  PRIMARY KEY(id)
);

CREATE TABLE users (
  id serial,
  hotspot_id bigint unsigned,
  username varchar(200),
  password varchar(200),
  email varchar(200),
  -- attributes for acces control
  redirection_url varchar(200),
  idle_timeout integer unsigned,
  kbps_down integer unsigned,
  kbps_up integer unsigned,
  limit_interval varchar(10), -- values: minute, hour, day, week, month
  limit_value integer unsigned, -- times limit_interval
  repeat_interval boolean default false,
  session_time integer unsigned,
  kbytes_total integer unsigned,
  kbytes_down integer unsigned,
  kbytes_up integer unsigned,
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
  -- attributes for acces control
  always_reject boolean default false,
  reply_message varchar(200),
  -- the following default from 'hotspots'
  redirection_url varchar(200),
  always_allow boolean default false,
  idle_timeout integer unsigned,
  kbps_down integer unsigned,
  kbps_up integer unsigned,
  limit_interval varchar(10), -- values: minute, hour, day, week, month
  limit_value integer unsigned, -- times limit_interval
  repeat_interval boolean default false,
  session_time integer unsigned,
  kbytes_total integer unsigned,
  kbytes_down integer unsigned,
  kbytes_up integer unsigned,
  valid_from datetime,
  valid_until datetime,
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
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id),
  KEY(mac_address),
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