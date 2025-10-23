-- This is a transaction to move old records from work_session_archived to work_session_history
-- and then delete them from work_session_archived to keep the table size manageable.
-- Records older than 730 days (2 years) will be moved and the user_id link will also be removed to respect RGPD
-- and ensure data privacy.

BEGIN;

INSERT INTO
    work_session_history (
        uuid,
        clock_in,
        clock_out,
        duration_minutes,
        status,
        updated_at,
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
    created_at,
    NOW() AS archived_at
FROM work_session_archived
WHERE
    created_at < NOW() - INTERVAL '730 days';

DELETE FROM work_session_archived
WHERE
    created_at < NOW() - INTERVAL '730 days';

COMMIT;