-- This is a transaction to move active records from work_session_active to work_session_archived
-- and then delete them from work_session_active to keep the table size manageable.
-- Records older than 30 days will be moved

BEGIN;

INSERT INTO
    work_session_archived (
        uuid,
        clock_in,
        clock_out,
        duration_minutes,
        status,
        updated_at,
        user_id,
        created_at,
        archived_at
    )
SELECT
    uuid,
    clock_in,
    clock_out,
    duration_minutes,
    status,
    updated_at,
    user_id,
    created_at,
    NOW() AS archived_at
FROM work_session_active
WHERE
    created_at < NOW() - INTERVAL '30 days';

DELETE FROM work_session_active
WHERE
    created_at < NOW() - INTERVAL '30 days';

COMMIT;