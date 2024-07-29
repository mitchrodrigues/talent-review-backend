-- Down Migration 20240729071722230323 add_terminated_at_to_employees

ALTER TABLE employees DROP COLUMN terminated_at;
