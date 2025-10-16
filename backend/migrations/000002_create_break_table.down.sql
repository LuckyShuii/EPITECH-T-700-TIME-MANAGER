DROP TABLE IF EXISTS breaks;

DROP TYPE IF EXISTS break_status;

ALTER TABLE work_session_active
DROP COLUMN IF EXISTS breaks_duration_minutes;

ALTER TABLE work_session_archived
DROP COLUMN IF EXISTS breaks_duration_minutes;

ALTER TABLE work_session_history
DROP COLUMN IF EXISTS breaks_duration_minutes;