-- +migrate Up
CREATE TABLE oauth2_states  (
  state      TEXT NOT NULL,
  expires_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (state)
) STRICT;

CREATE TABLE checks (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  "group" TEXT NOT NULL DEFAULT 'default',
  schedule TEXT NOT NULL,
  script TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_checks_name UNIQUE (name)
) STRICT;


--CREATE TRIGGER checks_updated_timestamp AFTER UPDATE ON checks BEGIN
--  update checks set updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') where id = new.id;
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

CREATE TABLE check_worker_groups (
  worker_group_id TEXT NOT NULL,
  check_id      TEXT NOT NULL,

  PRIMARY KEY (worker_group_id,check_id),
  CONSTRAINT fk_check_worker_groups_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_check_worker_groups_check FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE check_histories  (
  check_id TEXT NOT NULL,
  worker_group_id TEXT NOT NULL,

  status       TEXT NOT NULL,
  note         TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (check_id, worker_group_id, created_at),
  CONSTRAINT fk_check_histories_check FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE,
  CONSTRAINT fk_check_histories_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE
) STRICT;

CREATE TABLE triggers (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  script TEXT NOT NULL,
  status TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_triggers_name UNIQUE (name)
) STRICT;

CREATE TABLE trigger_histories  (
  trigger_id TEXT NOT NULL,

  status       TEXT NOT NULL,
  note         TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (trigger_id, created_at),
  CONSTRAINT fk_trigger_histories_trigger FOREIGN KEY (trigger_id) REFERENCES triggers(id) ON DELETE CASCADE
) STRICT;

-- +migrate Down
DROP TABLE oauth2_states;
DROP TABLE check_worker_groups;
DROP TABLE worker_groups;
DROP TABLE check_histories;
DROP TABLE checks;
DROP TABLE triggers;
DROP TABLE trigger_histories;
