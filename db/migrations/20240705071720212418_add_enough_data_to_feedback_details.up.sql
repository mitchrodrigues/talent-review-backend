-- Up Migration 20240705071720212418 add_enough_data_to_feedback_details

-- beginStatement
ALTER TABLE feedback_details 
    ADD COLUMN enough_data BOOLEAN NOT NULL,
    ADD COLUMN opportunities TEXT,
    DROP COLUMN opporuntities;
-- endStatement