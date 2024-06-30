-- Up Migration 20240619061718827910 create_users_tables

-- beginStatement
CREATE TABLE users (
    id UUID NOT NULL,

    idp_id        VARCHAR(255),
    idp_invite_id VARCHAR(255),

    first_name  VARCHAR(255),
    last_name   VARCHAR(255),
    email       VARCHAR(255),

    status            VARCHAR(32),
    status_updated_at TIMESTAMP,
 
    organization_id UUID NOT NULL, 
    inviter_id      UUID,

    version INT,

    invited_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    PRIMARY KEY(id)
)
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX users_idp_id_idx ON users (idp_id, deleted_at)
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX users_idp_invite_idx ON users (idp_invite_id, deleted_at)
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX users_email_idx ON users (email, deleted_at)
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX users_organization_idx ON users (organization_id, deleted_at, id)
-- endStatement
