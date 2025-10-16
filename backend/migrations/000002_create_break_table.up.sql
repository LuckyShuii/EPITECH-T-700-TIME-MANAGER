CREATE TYPE break_status AS ENUM(
    'active',
    'completed'
);

CREATE TABLE breaks (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    duration_minutes INTEGER,
    status break_status DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    work_session_active_id INTEGER NOT NULL REFERENCES work_session_active (id) ON DELETE CASCADE
);

ALTER TABLE work_session_active
ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;

ALTER TABLE work_session_archived
ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;

ALTER TABLE work_session_history
ADD COLUMN breaks_duration_minutes INTEGER DEFAULT 0;