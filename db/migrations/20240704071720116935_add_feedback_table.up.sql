-- Up Migration 20240704071720116935 add_feedback_table

-- beginStatement
CREATE TABLE feedbacks (
    id UUID NOT NULL PRIMARY KEY,
    
    email VARCHAR(255) NOT NULL,
    code VARCHAR(16) NOT NULL,

    version INT,
    
    collection_end_at  TIMESTAMP WITH TIME ZONE,
    submitted_at TIMESTAMP WITH TIME ZONE,

    employee_id UUID NOT NULL,
    organization_id UUID NOT NULL,
    owner_id UUID NOT NULL,

    updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- endStatement

-- beginStatement
CREATE INDEX ON feedbacks (employee_id, email);
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX ON feedbacks (code, deleted_at);
-- endStatement

-- beginStatement
CREATE INDEX ON feedbacks (organization_id, deleted_at);
-- endStatement

-- beginStatement
CREATE TABLE feedback_details (
    id UUID NOT NULL PRIMARY KEY,
    
    feedback_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    organization_id UUID NOT NULL,

    enough_data BOOLEAN NOT NULL,

    strenghts TEXT,
    opportunities TEXT,
    additional TEXT,
    rating INT,

    updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
-- endStatement

-- beginStatement
CREATE UNIQUE INDEX ON feedback_details (feedback_id, organization_id, deleted_at);
-- endStatement
