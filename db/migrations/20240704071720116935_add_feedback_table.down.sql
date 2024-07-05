-- Down Migration 20240704071720116935 add_feedback_table

-- beginStatement
DROP TABLE IF EXISTS feedbacks;
-- endStatement

-- beginStatement
DROP TABLE IF EXISTS feedback_details;
-- endStatement