-- add "on delete cascade" from work_session_active on user_id
ALTER TABLE work_session_active
DROP CONSTRAINT work_session_active_user_id_fkey;

ALTER TABLE work_session_active
ADD CONSTRAINT work_session_active_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;