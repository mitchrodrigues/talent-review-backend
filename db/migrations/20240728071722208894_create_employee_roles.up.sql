-- Up Migration 20240728071722208894 create_employee_roles
-- beginStatement
CREATE TABLE employee_roles (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    title VARCHAR(255) NOT NULL,
    track VARCHAR(32) NOT NULL,
    level INT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX idx_employee_roles_org_title_deleted ON employee_roles(organization_id, title, deleted_at);
-- endStatement

-- beginStatement
CREATE INDEX idx_employee_roles_deleted_at ON employee_roles(deleted_at);
-- endStatement

-- beginStatement
ALTER TABLE
    employees
ADD
    COLUMN employee_role_id UUID REFERENCES employee_roles(id),
ALTER 
    COLUMN type DROP NOT NULL,
ALTER 
    COLUMN level DROP NOT NULL,      
ALTER 
    COLUMN title DROP NOT NULL;
-- endStatement

-- beginStatement
CREATE INDEX idx_employees_org_role_deleted ON employees(employee_role_id, organization_id, deleted_at);
-- endStatement