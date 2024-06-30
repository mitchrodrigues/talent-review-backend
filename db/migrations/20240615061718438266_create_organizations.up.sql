-- Up Migration 20240615061718438266 create_organizations_and_users

-- beginStatement
CREATE TABLE organizations (
    id                  UUID NOT NULL,
    name                 VARCHAR(255),
    idp_id               VARCHAR(255),

    merchant_customer_id VARCHAR(255),
    merchant_plan_id     VARCHAR(255),
    merchant_plan_name   VARCHAR(32),

    activated_at        TIMESTAMP,
    deactivated_at      TIMESTAMP,

    created_at          TIMESTAMP,
    updated_at          TIMESTAMP,
    deleted_at          TIMESTAMP,

    version INT,

    PRIMARY KEY (id)
)
-- endStatement

-- beginStatement
CREATE INDEX org_deleted_idx ON organizations (deleted_at);
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX org_idp_id_idx ON organizations (idp_id, deleted_at)
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX org_merchant_id ON organizations (merchant_customer_id, deleted_at)
-- endStatement