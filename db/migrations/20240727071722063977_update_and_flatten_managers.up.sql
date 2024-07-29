-- Up Migration 20240727071722063977 update_and_flatten_managers

-- beginStatement
ALTER TABLE events
ALTER COLUMN data SET DATA TYPE jsonb USING data::jsonb;
-- endStatement

-- beginStatement
ALTER TABLE
    employees
ADD
    COLUMN manager_id uuid REFERENCES employees(id);
-- endStatement


-- beginStatement
CREATE TABLE employee_histories (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    employee_id uuid NOT NULL REFERENCES employees(id),
    user_id uuid REFERENCES users(id),
    change jsonb NOT NULL,
    
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
-- endStatement

-- beginStatement
CREATE INDEX idx_employee_histories_employee_id ON employee_histories(employee_id, deleted_at);
-- endStatement

-- beginStatement
CREATE INDEX idx_employee_histories_user_id ON employee_histories(user_id, deleted_at);
-- endStatement

-- beginStatement
CREATE INDEX idx_employee_histories_created_at ON employee_histories(created_at, deleted_at);
-- endStatement


-- beginStatement
UPDATE employees
SET manager_id = sub.manager_id
FROM (
    SELECT e.id, t.manager_id
    FROM employees e
    JOIN teams t ON e.team_id = t.id
    WHERE t.id IS NOT NULL AND t.manager_id IS NOT NULL
    AND t.manager_id <> '00000000-0000-0000-0000-000000000000'
) sub
WHERE employees.id = sub.id;
-- endStatement


-- beginStatement
INSERT INTO employee_histories (employee_id, user_id, change, created_at, updated_at)
SELECT 
    e.id AS employee_id, 
    NULL AS user_id, 
    jsonb_build_object('previous', NULL, 'current', t.manager_id, 'field', 'manager_id') AS change,
    e.created_at AS created_at,
    e.updated_at AS updated_at
FROM employees e
JOIN teams t ON e.team_id = t.id
WHERE t.id IS NOT NULL AND t.manager_id IS NOT NULL
AND t.manager_id <> '00000000-0000-0000-0000-000000000000';
-- endStatement

-- beginStatement
ALTER TABLE teams ADD COLUMN lead_id uuid REFERENCES employees(id);
-- endStatement

-- beginStatement
UPDATE teams SET lead_id = manager_id
WHERE manager_id <> '00000000-0000-0000-0000-000000000000'
-- endStatement

-- beginStatement
CREATE INDEX idx_teams_lead_id ON teams(lead_id, deleted_at);
-- endStatement


-- beginStatement
-- Update event types and rename fields in TeamCreated events
UPDATE events
SET type = 'team.Created',
    data = jsonb_set(
        jsonb_set(
            data,
            '{LeadID}',
            data->'ManagerID'
        ),
        '{ManagerID}',
        'null'::jsonb
    ) || jsonb_build_object(
        'leadID', data->'ManagerID',
        'organizationID', data->'OrganizationID',
        'name', data->'Name'
    ) - 'ManagerID' - 'OrganizationID' - 'Name'
WHERE type = 'teams.TeamCreated';
-- endStatement

-- beginStatement
-- Update event types and rename fields in TeamUpdated events
UPDATE events
SET type = 'team.Updated',
    data = jsonb_set(
        jsonb_set(
            data,
            '{LeadID}',
            data->'ManagerID'
        ),
        '{ManagerID}',
        'null'::jsonb
    ) || jsonb_build_object(
        'leadID', data->'ManagerID',
        'name', data->'Name'
    ) - 'ManagerID' - 'Name'
WHERE type = 'teams.TeamUpdated';
-- endStatement