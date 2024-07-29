-- Down Migration 20240728071722208894 create_employee_roles

-- beginStatement
DROP INDEX IF EXISTS idx_employees_org_role_deleted;
-- endStatement

-- beginStatement
ALTER TABLE employees DROP COLUMN IF EXISTS employee_role_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_roles_org_title_deleted;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_roles_org_level_deleted;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_roles_version_deleted;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_roles_deleted_at;
-- endStatement

-- beginStatement
DROP TABLE IF EXISTS employee_roles;
-- endStatement