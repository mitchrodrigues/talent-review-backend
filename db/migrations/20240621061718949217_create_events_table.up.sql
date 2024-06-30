-- Up Migration 20240621061718949217 create_events_table

-- beginStatement
 CREATE TABLE events (
  id uuid NOT NULL,

  aggregate_id uuid,
  aggregate_type VARCHAR(64),

  version INTEGER,
  type    VARCHAR(64),

  data JSON,
  metadata JSON,

  user_id UUID,
  organization_id UUID,

  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP,

  PRIMARY KEY (id)

);
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX ON events (aggregate_id, aggregate_type, id, deleted_at);
-- endStatement

-- beginStatement
CREATE INDEX ON events (user_id, organization_id, deleted_at);
-- endStatement
