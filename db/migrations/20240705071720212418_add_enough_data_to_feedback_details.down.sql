-- Down Migration 20240705071720212418 add_enough_data_to_feedback_details

-- beginStatement
ALTER TABLE
    feedback_details DROP COLUMN enough_data;
-- endStatement