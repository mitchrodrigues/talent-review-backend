-- Up Migration 20240621061718944212 create_employees_tables

-- beginStatement
CREATE TABLE employees (
    id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    worker_type char(4) NOT NULL,
    type VARCHAR(12) NOT NULL,
    level INTEGER NOT NULL,
    
    level_start_at TIMESTAMP,
    employement_start_at TIMESTAMP,

    team_id UUID,
    user_id UUID,    
    organization_id UUID NOT NULL,

    version INTEGER DEFAULT 1,

    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    PRIMARY KEY(id)
);
-- endStatement

-- beginStatement
CREATE INDEX employee_org_type_level_idx ON employees (organization_id, level, type, deleted_at);
-- endStatement

-- beginStatement
CREATE TABLE IF NOT EXISTS teams (
    id UUID NOT NULL,
    manager_id UUID,

    name VARCHAR(64) NOT NULL,
    organization_id UUID NOT NULL,

    version INTEGER DEFAULT 1,


    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    PRIMARY KEY(id)
);
-- endStatement

-- beginStatement
CREATE INDEX team_org_idx ON teams (organization_id, deleted_at);
-- endStatement
