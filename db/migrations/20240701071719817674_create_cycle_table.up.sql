-- Up Migration 20240701071719817674 create_cycle_table

CREATE TABLE cycles (
    id UUID PRIMARY KEY,
    
    owner_id UUID NOT NULL,
    employee_id UUID NOT NULL,
    organization_id UUID NOT NULL,

    type VARCHAR(32),

    start_at TIMESTAMP WITH TIME ZONE NOT NULL,
    end_at TIMESTAMP WITH TIME ZONE NOT NULL,

    updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);
