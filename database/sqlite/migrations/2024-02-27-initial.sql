-- +migrate Up
CREATE TABLE oauth2_states  (
  state      TEXT NOT NULL,
  expires_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (state)
) STRICT;

CREATE TABLE monitors (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  "group" TEXT NOT NULL DEFAULT 'default',
  schedule TEXT NOT NULL,
  script TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_monitors_name UNIQUE (name)
) STRICT;


--CREATE TRIGGER monitors_updated_timestamp AFTER UPDATE ON monitors BEGIN
--  update monitors set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = new.id;
--END;

CREATE TABLE worker_groups  (
  id TEXT NOT NULL,
  name TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_worker_groups_name UNIQUE (name)
) STRICT;

--CREATE TRIGGER worker_groups_updated_timestamp AFTER UPDATE ON worker_groups BEGIN
--  update worker_groups set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = new.id;
--END;

CREATE TABLE monitor_worker_groups (
  worker_group_id TEXT NOT NULL,
  monitor_id      TEXT NOT NULL,

  PRIMARY KEY (worker_group_id,monitor_id),
  CONSTRAINT fk_monitor_worker_groups_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_monitor_worker_groups_monitor FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE monitor_histories  (
  monitor_id TEXT NOT NULL,
  worker_group_id TEXT NOT NULL,

  status       TEXT NOT NULL,
  note         TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (monitor_id, worker_group_id, created_at),
  CONSTRAINT fk_monitor_histories_monitor FOREIGN KEY (monitor_id) REFERENCES monitors(id) ON DELETE CASCADE,
  CONSTRAINT fk_monitor_histories_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE
) STRICT;

-- +migrate Down
DROP TABLE oauth2_states;
DROP TABLE monitor_worker_groups;
DROP TABLE worker_groups;
DROP TABLE monitor_histories;
DROP TABLE monitors;
