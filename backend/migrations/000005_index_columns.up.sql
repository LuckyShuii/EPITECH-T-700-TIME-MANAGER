-- USERS
CREATE INDEX idx_users_id ON users (id);

CREATE INDEX idx_users_uuid ON users (uuid);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);

-- WORK SESSION ACTIVE
CREATE INDEX idx_work_session_active_user_id ON work_session_active (user_id);

CREATE INDEX idx_work_session_active_status ON work_session_active (status);

CREATE INDEX idx_work_session_active_clock_in ON work_session_active (clock_in);

-- WORK SESSION ARCHIVED
CREATE INDEX idx_work_session_archived_user_id ON work_session_archived (user_id);

CREATE INDEX idx_work_session_archived_clock_in ON work_session_archived (clock_in);

CREATE INDEX idx_work_session_archived_archived_at ON work_session_archived (archived_at);

-- WORK SESSION HISTORY
CREATE INDEX idx_work_session_history_archived_at ON work_session_history (archived_at);

CREATE INDEX idx_work_session_history_clock_in ON work_session_history (clock_in);

-- TEAM MEMBERS
CREATE INDEX idx_teams_members_user_id ON teams_members (user_id);

CREATE INDEX idx_teams_members_team_id ON teams_members (team_id);

CREATE INDEX idx_teams_members_is_manager ON teams_members (is_manager);

-- TEAMS
CREATE INDEX idx_teams_uuid ON teams (uuid);

CREATE INDEX idx_teams_id ON teams (id);