-- +migrate Up
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
DROP TABLE triggers;
DROP TABLE trigger_histories;
