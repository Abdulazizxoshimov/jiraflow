DELETE FROM workflow_transitions
    WHERE workflow_id = '00000000-0000-0000-0000-000000000001'::uuid;
DELETE FROM workflow_statuses
    WHERE workflow_id = '00000000-0000-0000-0000-000000000001'::uuid;
DELETE FROM workflows
    WHERE id = '00000000-0000-0000-0000-000000000001'::uuid;
