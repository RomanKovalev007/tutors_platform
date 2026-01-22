CREATE OR REPLACE FUNCTION mark_expired_tasks()
RETURNS TABLE(
    task_id UUID,
    old_status VARCHAR(50),
    expired_at TIMESTAMPTZ
) AS $$
BEGIN
    RETURN QUERY
    WITH updated AS (
        UPDATE assigned_tasks
        SET task_status = 'EXPIRED',
            updated_at = NOW()
        WHERE task_status = 'ACTIVE'
          AND deadline IS NOT NULL
          AND deadline < NOW()
        RETURNING 
            id, 
            task_status AS old_status,
            NOW() AS expired_at
    )
    SELECT 
        u.id AS task_id,
        u.old_status,
        u.expired_at
    FROM updated u;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION mark_expired_tasks() IS 'Автоматически помечает просроченные активные задачи как expired';