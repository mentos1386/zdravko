-- +migrate Up
CREATE TABLE o_auth2_states  (
  state TEXT,
  expiry DATETIME,
  PRIMARY KEY (state)
);

CREATE TABLE monitors (
  slug TEXT,
  name TEXT,
  schedule TEXT,
  script TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME,

  PRIMARY KEY (slug),
  CONSTRAINT unique_monitors_name UNIQUE (name)
);

CREATE TABLE worker_groups  (
  slug TEXT,
  name TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME,

  PRIMARY KEY (slug),
  CONSTRAINT unique_worker_groups_name UNIQUE (name)
);

CREATE TABLE monitor_worker_groups (
  worker_group_slug TEXT,
  monitor_slug TEXT,
  PRIMARY KEY (worker_group_slug,monitor_slug),
  CONSTRAINT fk_monitor_worker_groups_worker_group FOREIGN KEY (worker_group_slug) REFERENCES worker_groups(slug),
  CONSTRAINT fk_monitor_worker_groups_monitor FOREIGN KEY (monitor_slug) REFERENCES monitors(slug)
);

CREATE TABLE "monitor_histories"  (
  monitor_slug TEXT,
  status TEXT,
  note TEXT,

  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (monitor_slug, created_at),
  CONSTRAINT fk_monitors_history FOREIGN KEY (monitor_slug) REFERENCES monitors(slug)
);

-- +migrate Down
DROP TABLE o_auth2_states;
DROP TABLE monitor_worker_groups;
DROP TABLE worker_groups;
DROP TABLE monitor_histories;
DROP TABLE monitors;
