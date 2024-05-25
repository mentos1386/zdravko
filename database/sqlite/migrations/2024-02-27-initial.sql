-- +migrate Up
CREATE TABLE oauth2_states  (
  state      TEXT NOT NULL,
  expires_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (state)
) STRICT;

CREATE TABLE checks (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  schedule TEXT NOT NULL,
  script TEXT NOT NULL,

  filter TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_checks_name UNIQUE (name)
) STRICT;
-- +migrate StatementBegin
CREATE TRIGGER checks_updated_timestamp AFTER UPDATE ON checks BEGIN
  UPDATE checks SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') WHERE id = NEW.id;
END;
-- +migrate StatementEnd

CREATE TABLE worker_groups  (
  id TEXT NOT NULL,
  name TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_worker_groups_name UNIQUE (name)
) STRICT;
-- +migrate StatementBegin
CREATE TRIGGER worker_groups_updated_timestamp AFTER UPDATE ON worker_groups BEGIN
  UPDATE worker_groups SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') WHERE id = NEW.id;
END;
-- +migrate StatementEnd

CREATE TABLE check_worker_groups (
  worker_group_id TEXT NOT NULL,
  check_id      TEXT NOT NULL,

  PRIMARY KEY (worker_group_id,check_id),
  CONSTRAINT fk_check_worker_groups_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_check_worker_groups_check FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE
) STRICT;


CREATE TABLE triggers (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  script TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_triggers_name UNIQUE (name)
) STRICT;
-- +migrate StatementBegin
CREATE TRIGGER triggers_updated_timestamp AFTER UPDATE ON triggers BEGIN
  UPDATE triggers SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') WHERE id = NEW.id;
END;
-- +migrate StatementEnd

CREATE TABLE targets (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  "group" TEXT NOT NULL,

  visibility TEXT NOT NULL,
  state TEXT NOT NULL,

  metadata TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_targets_name UNIQUE (name)
) STRICT;
-- +migrate StatementBegin
CREATE TRIGGER targets_updated_timestamp AFTER UPDATE ON targets BEGIN
  UPDATE targets SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') WHERE id = NEW.id;
END;
-- +migrate StatementEnd

CREATE TABLE target_histories  (
  target_id       TEXT NOT NULL,
  worker_group_id TEXT NOT NULL,
  check_id        TEXT NOT NULL,

  status       TEXT NOT NULL,
  note         TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (target_id, worker_group_id, check_id, created_at),
  CONSTRAINT fk_target_histories_target FOREIGN KEY (target_id) REFERENCES targets(id) ON DELETE CASCADE,
  CONSTRAINT fk_target_histories_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_target_histories_check FOREIGN KEY (check_id) REFERENCES checks(id) ON DELETE CASCADE
) STRICT;

-- +migrate Down
DROP TABLE oauth2_states;
DROP TABLE check_worker_groups;
DROP TABLE worker_groups;
DROP TRIGGER worker_groups_updated_timestamp;
DROP TABLE checks;
DROP TRIGGER checks_updated_timestamp;
DROP TABLE triggers;
DROP TRIGGER triggers_updated_timestamp;
