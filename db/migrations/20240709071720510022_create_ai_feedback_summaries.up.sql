-- Up Migration 20240709071720510022 create_ai_feedback_summaries
CREATE TABLE feedback_summaries (
    id uuid PRIMARY KEY,
    summary text,
    action_items jsonb,
    employee_id uuid NOT NULL,
    organization_id uuid NOT NULL,
    feedback_id UUID not null,
    updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
)
