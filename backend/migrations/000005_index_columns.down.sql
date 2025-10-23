-- TEAMS
DROP INDEX IF EXISTS idx_teams_id;

DROP INDEX IF EXISTS idx_teams_uuid;

-- TEAM MEMBERS
DROP INDEX IF EXISTS idx_teams_members_is_manager;

DROP INDEX IF EXISTS idx_teams_members_team_id;

DROP INDEX IF EXISTS idx_teams_members_user_id;

-- WORK SESSION HISTORY
DROP INDEX IF EXISTS idx_work_session_history_clock_in;

DROP INDEX IF EXISTS idx_work_session_history_archived_at;

-- WORK SESSION ARCHIVED
DROP INDEX IF EXISTS idx_work_session_archived_archived_at;

DROP INDEX IF EXISTS idx_work_session_archived_clock_in;

DROP INDEX IF EXISTS idx_work_session_archived_user_id;

-- WORK SESSION ACTIVE
DROP INDEX IF EXISTS idx_work_session_active_clock_in;

DROP INDEX IF EXISTS idx_work_session_active_status;

DROP INDEX IF EXISTS idx_work_session_active_user_id;

-- USERS
DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_uuid;

DROP INDEX IF EXISTS idx_users_id;