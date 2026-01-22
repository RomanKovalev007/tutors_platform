-- Для быстрой фильтрации по статусу (частый запрос)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_assigned_tasks_status 
ON assigned_tasks(task_status);

-- Оптимизация для функции mark_expired_tasks: активные задачи с дедлайном
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_assigned_tasks_active_deadline 
ON assigned_tasks(deadline) 
WHERE task_status = 'ACTIVE' AND deadline IS NOT NULL;

-- Для поиска задач по группе (если часто фильтруете по group_id)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_assigned_tasks_group_id 
ON assigned_tasks(group_id);

-- Для поиска задач преподавателя (если часто по tutor_id)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_assigned_tasks_tutor_id 
ON assigned_tasks(tutor_id);

-- Для поиска всех работ по конкретной задаче (самый частый запрос)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submitted_tasks_task_id 
ON submitted_tasks(task_id);

-- Для поиска работ конкретного студента
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submitted_tasks_student_id 
ON submitted_tasks(student_id);

-- Составной индекс: студент + задача (проверка сдачи)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_submitted_tasks_student_task 
ON submitted_tasks(student_id, task_id);
