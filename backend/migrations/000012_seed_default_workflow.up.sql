-- =============================================================================
-- SEED DATA — Standart workflow (yangi loyihalar uchun)
-- =============================================================================

-- 1. Default workflow
INSERT INTO workflows (id, name, description, is_default, created_by)
SELECT
    '00000000-0000-0000-0000-000000000001'::uuid,
    'Default Workflow',
    'Standart To Do → In Progress → In Review → Done workflow',
    TRUE,
    -- Birinchi admin foydalanuvchi yo'q bo'lsa, system user yaratish kerak.
    -- Hozircha NULL'ga ruxsat bermayotgan bo'lsak, seed'ni alohida bosqichda
    -- (admin user yaratilgandan keyin) bajarishimiz mumkin.
    -- Bu yerda biz dummy admin yaratamiz:
    (SELECT id FROM users WHERE role = 'admin' LIMIT 1)
WHERE EXISTS (SELECT 1 FROM users WHERE role = 'admin')
ON CONFLICT (id) DO NOTHING;

-- 2. Default statuses
-- Bu seed faqat workflow yaratilgan bo'lsa ishlaydi
DO $$
DECLARE
    wf_id UUID := '00000000-0000-0000-0000-000000000001'::uuid;
    todo_id UUID;
    in_progress_id UUID;
    in_review_id UUID;
    done_id UUID;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM workflows WHERE id = wf_id) THEN
        RAISE NOTICE 'Default workflow yaratilmagan, seed skip qilinmoqda';
        RETURN;
    END IF;

    -- To Do
    INSERT INTO workflow_statuses (workflow_id, name, category, color, position, is_initial)
    VALUES (wf_id, 'To Do', 'todo', '#6B7280', 1, TRUE)
    ON CONFLICT (workflow_id, name) DO NOTHING
    RETURNING id INTO todo_id;

    -- In Progress
    INSERT INTO workflow_statuses (workflow_id, name, category, color, position)
    VALUES (wf_id, 'In Progress', 'in_progress', '#3B82F6', 2)
    ON CONFLICT (workflow_id, name) DO NOTHING
    RETURNING id INTO in_progress_id;

    -- In Review
    INSERT INTO workflow_statuses (workflow_id, name, category, color, position)
    VALUES (wf_id, 'In Review', 'in_progress', '#F59E0B', 3)
    ON CONFLICT (workflow_id, name) DO NOTHING
    RETURNING id INTO in_review_id;

    -- Done
    INSERT INTO workflow_statuses (workflow_id, name, category, color, position)
    VALUES (wf_id, 'Done', 'done', '#10B981', 4)
    ON CONFLICT (workflow_id, name) DO NOTHING
    RETURNING id INTO done_id;

    -- ID'larni alohida olib chiqamiz (RETURNING DO NOTHING bilan ishlamaydi)
    SELECT id INTO todo_id        FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'To Do';
    SELECT id INTO in_progress_id FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'In Progress';
    SELECT id INTO in_review_id   FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'In Review';
    SELECT id INTO done_id        FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'Done';

    -- Transitions: linear
    INSERT INTO workflow_transitions (workflow_id, from_status_id, to_status_id, name) VALUES
        (wf_id, todo_id,        in_progress_id, 'Start progress'),
        (wf_id, in_progress_id, in_review_id,   'Submit for review'),
        (wf_id, in_review_id,   done_id,        'Mark as done'),
        -- Backward:
        (wf_id, in_progress_id, todo_id,        'Stop progress'),
        (wf_id, in_review_id,   in_progress_id, 'Reject review'),
        (wf_id, done_id,        in_progress_id, 'Reopen'),
        -- Global (istalgan holatdan):
        (wf_id, NULL,           done_id,        'Force close')
    ON CONFLICT (workflow_id, from_status_id, to_status_id) DO NOTHING;
END $$;
