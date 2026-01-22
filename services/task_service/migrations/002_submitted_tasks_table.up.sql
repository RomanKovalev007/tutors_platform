CREATE TYPE submitted_task_status AS ENUM(
    'PENDING',
    'VERIFIED'
);

CREATE TABLE IF NOT EXISTS submitted_tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID REFERENCES assigned_tasks(id),
    student_id VARCHAR(255) NOT NULL,
    content TEXT,
    status submitted_task_status NOT NULL,
    score SMALLINT,
    tutor_feedback TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    overdue_by INTERVAL
);