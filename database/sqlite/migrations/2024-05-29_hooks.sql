-- +migrate Up
CREATE TABLE hooks (
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  schedule TEXT NOT NULL,
  script TEXT NOT NULL,

  created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ')),

  PRIMARY KEY (id),
  CONSTRAINT unique_hooks_name UNIQUE (name)
) STRICT;
-- +migrate StatementBegin
CREATE TRIGGER hooks_updated_timestamp AFTER UPDATE ON hooks BEGIN
  UPDATE hooks SET updated_at = strftime('%Y-%m-%dT%H:%M:%fZ') WHERE id = NEW.id;
END;
-- +migrate StatementEnd

CREATE TABLE hook_worker_groups (
  worker_group_id TEXT NOT NULL,
  hook_id      TEXT NOT NULL,

  PRIMARY KEY (worker_group_id,hook_id),
  CONSTRAINT fk_hook_worker_groups_worker_group FOREIGN KEY (worker_group_id) REFERENCES worker_groups(id) ON DELETE CASCADE,
  CONSTRAINT fk_hook_worker_groups_hook FOREIGN KEY (hook_id) REFERENCES hooks(id) ON DELETE CASCADE
) STRICT;
-- +migrate Down
DROP TABLE hook_worker_groups;
DROP TABLE hooks;
