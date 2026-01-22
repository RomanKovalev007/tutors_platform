CREATE TYPE assigned_task_status AS ENUM(
    'ACTIVE',
    'EXPIRED'
);

CREATE TABLE IF NOT EXISTS assigned_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id VARCHAR(255) NOT NULL,
    tutor_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    max_score SMALLINT,
    deadline TIMESTAMPTZ,
    task_status assigned_task_status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);