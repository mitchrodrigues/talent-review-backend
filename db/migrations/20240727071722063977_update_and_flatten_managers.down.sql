-- Down Migration 20240727071722063977 update_and_flatten_managers

-- beginStatement
ALTER TABLE employees DROP COLUMN manager_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employees_manager_id;
-- endStatement


-- beginStatement
ALTER TABLE teams DROP COLUMN lead_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_teams_lead_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_histories_employee_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_histories_user_id;
-- endStatement

-- beginStatement
DROP INDEX IF EXISTS idx_employee_histories_created_at;
-- endStatement

-- beginStatement
DROP TABLE IF EXISTS employee_histories;
-- endStatement

-- beginStatement
-- Revert event types and rename fields in team.Created events
UPDATE events
SET type = 'teams.TeamCreated',
    data = jsonb_set(
        data - 'leadID',
        '{ManagerID}',
        data->'leadID'
    ) || jsonb_build_object(
        'ManagerID', data->'leadID',
        'OrganizationID', data->'organizationID',
        'Name', data->'name'
    ) - 'leadID' - 'organizationID' - 'name'
WHERE type = 'team.Created';
--endStatement

-- beginStatement
-- Revert event types and rename fields in team.Updated events
UPDATE events
SET type = 'teams.TeamUpdated',
    data = jsonb_set(
        data - 'leadID',
        '{ManagerID}',
        data->'leadID'
    ) || jsonb_build_object(
        'ManagerID', data->'leadID',
        'Name', data->'name'
    ) - 'leadID' - 'name'
WHERE type = 'team.Updated';
-- endStatement