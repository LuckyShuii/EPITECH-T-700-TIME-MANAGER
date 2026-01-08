-- Seed Timemanager (idempotent-ish)

-- This fixture file is launched automatically in dev mode to populate the database
-- with sample data for easier testing.
-- every password for each users here is 'lboillot'

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ------------------------------------------------------------
-- 0) Weekly rates
-- ------------------------------------------------------------
INSERT INTO
    weekly_rate (uuid, rate_name, amount)
SELECT '11111111-1111-1111-1111-111111111111', 'Temps pleins', 35
WHERE
    NOT EXISTS (
        SELECT 1
        FROM weekly_rate
        WHERE
            rate_name = 'Temps pleins'
    );

INSERT INTO
    weekly_rate (uuid, rate_name, amount)
SELECT '22222222-2222-2222-2222-222222222222', 'Temps pleins + RTT', 39
WHERE
    NOT EXISTS (
        SELECT 1
        FROM weekly_rate
        WHERE
            rate_name = 'Temps pleins + RTT'
    );

INSERT INTO
    weekly_rate (uuid, rate_name, amount)
SELECT '33333333-3333-3333-3333-333333333333', 'Temps partiel', 20
WHERE
    NOT EXISTS (
        SELECT 1
        FROM weekly_rate
        WHERE
            rate_name = 'Temps partiel'
    );

-- ------------------------------------------------------------
-- 1) Users (10) : 4 employees, 4 managers, 2 admins
-- ------------------------------------------------------------


WITH rates AS (
  SELECT
    (SELECT id FROM weekly_rate WHERE rate_name = 'Temps pleins' LIMIT 1)        AS full_time_id,
    (SELECT id FROM weekly_rate WHERE rate_name = 'Temps pleins + RTT' LIMIT 1) AS full_time_rtt_id,
    (SELECT id FROM weekly_rate WHERE rate_name = 'Temps partiel' LIMIT 1)      AS part_time_id
),
data AS (
  SELECT
    'a0000000-0000-0000-0000-000000000001'::varchar AS uuid,
    'employee1'::varchar AS username,
    'employee1@timemanager.local'::varchar AS email,
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK'::varchar AS password_hash,
    'Alice'::varchar AS first_name,
    'Martin'::varchar AS last_name,
    '0600000001'::varchar AS phone_number,
    ARRAY['employee']::text[] AS roles,
    'active'::varchar AS status,
    (SELECT full_time_id FROM rates) AS weekly_rate_id,
    1::int AS first_day_of_week

  UNION ALL SELECT
    'a0000000-0000-0000-0000-000000000002','employee2','employee2@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Bruno','Petit','0600000002',
    ARRAY['employee'],'active',(SELECT full_time_id FROM rates),1

  UNION ALL SELECT
    'a0000000-0000-0000-0000-000000000003','employee3','employee3@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Chloé','Durand','0600000003',
    ARRAY['employee'],'pending',(SELECT part_time_id FROM rates),1

  UNION ALL SELECT
    'a0000000-0000-0000-0000-000000000004','employee4','employee4@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'David','Moreau','0600000004',
    ARRAY['employee'],'active',(SELECT full_time_rtt_id FROM rates),1

  UNION ALL SELECT
    'b0000000-0000-0000-0000-000000000001','manager1','manager1@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Emma','Leroy','0600000011',
    ARRAY['employee','manager'],'active',(SELECT full_time_id FROM rates),1

  UNION ALL SELECT
    'b0000000-0000-0000-0000-000000000002','manager2','manager2@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Félix','Roux','0600000012',
    ARRAY['employee','manager'],'active',(SELECT full_time_rtt_id FROM rates),1

  UNION ALL SELECT
    'b0000000-0000-0000-0000-000000000003','manager3','manager3@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Gaëlle','Fontaine','0600000013',
    ARRAY['employee','manager'],'disabled',(SELECT full_time_id FROM rates),1

  UNION ALL SELECT
    'b0000000-0000-0000-0000-000000000004','manager4','manager4@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Hugo','Lambert','0600000014',
    ARRAY['employee','manager'],'active',(SELECT part_time_id FROM rates),1

  UNION ALL SELECT
    'c0000000-0000-0000-0000-000000000001','admin1','admin1@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Inès','Faure','0600000021',
    ARRAY['employee','admin'],'active',(SELECT full_time_id FROM rates),1

  UNION ALL SELECT
    'c0000000-0000-0000-0000-000000000002','admin2','admin2@timemanager.local',
    '$2a$10$0OpjORV/S0JSyM71GKAcyufumggyocR52rBmM5QgTrUc0EwEGIVKK',
    'Jules','Garnier','0600000022',
    ARRAY['employee','admin'],'active',(SELECT full_time_rtt_id FROM rates),1
)
INSERT INTO users (
  uuid, username, email, password_hash, first_name, last_name, phone_number,
  roles, status, weekly_rate_id, first_day_of_week
)
SELECT
  d.uuid, d.username, d.email, d.password_hash, d.first_name, d.last_name, d.phone_number,
  d.roles, d.status, d.weekly_rate_id, d.first_day_of_week
FROM data d
WHERE NOT EXISTS (
  SELECT 1 FROM users u
  WHERE u.email = d.email OR u.username = d.username
);

-- ------------------------------------------------------------
-- 2) Teams + members
-- ------------------------------------------------------------
INSERT INTO
    teams (uuid, name, description)
SELECT 'd0000000-0000-0000-0000-000000000001', 'Engineering', 'Core engineering team'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM teams
        WHERE
            name = 'Engineering'
    );

INSERT INTO
    teams (uuid, name, description)
SELECT 'd0000000-0000-0000-0000-000000000002', 'Support', 'Customer support team'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM teams
        WHERE
            name = 'Support'
    );

INSERT INTO
    teams (uuid, name, description)
SELECT 'd0000000-0000-0000-0000-000000000003', 'Product', 'Product & design team'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM teams
        WHERE
            name = 'Product'
    );

-- Managers (is_manager = true) - uuid généré
INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, TRUE
FROM users u JOIN teams t ON t.name = 'Engineering'
WHERE u.username = 'manager1'
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, TRUE
FROM users u JOIN teams t ON t.name = 'Support'
WHERE u.username = 'manager2'
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, TRUE
FROM users u JOIN teams t ON t.name = 'Product'
WHERE u.username = 'manager4'
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

-- Employees
INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, FALSE
FROM users u JOIN teams t ON t.name = 'Engineering'
WHERE u.username IN ('employee1','employee2')
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, FALSE
FROM users u JOIN teams t ON t.name = 'Support'
WHERE u.username IN ('employee3')
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, FALSE
FROM users u JOIN teams t ON t.name = 'Product'
WHERE u.username IN ('employee4')
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

-- Admins
INSERT INTO teams_members (uuid, user_id, team_id, is_manager)
SELECT gen_random_uuid()::text, u.id, t.id, FALSE
FROM users u JOIN teams t ON t.name = 'Engineering'
WHERE u.username IN ('admin1','admin2')
AND NOT EXISTS (SELECT 1 FROM teams_members tm WHERE tm.user_id = u.id AND tm.team_id = t.id);

-- ------------------------------------------------------------
-- 3) Work sessions (active / archived / history)
-- ------------------------------------------------------------
-- IMPORTANT: uuid généré => on utilise une condition NOT EXISTS basée sur user_id + clock_in

-- employee1: session complétée
INSERT INTO work_session_active (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '3 hours',
  NOW() - INTERVAL '10 minutes',
  170,
  'completed'::work_session_status,
  15
FROM users u
WHERE u.username = 'employee1'
AND NOT EXISTS (
  SELECT 1 FROM work_session_active w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '3 hours')
);

-- employee2: session active
INSERT INTO work_session_active (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '2 hours',
  NULL,
  NULL,
  'active'::work_session_status,
  0
FROM users u
WHERE u.username = 'employee2'
AND NOT EXISTS (
  SELECT 1 FROM work_session_active w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '2 hours')
);

-- manager1: session complétée
INSERT INTO work_session_active (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '5 hours',
  NOW() - INTERVAL '1 hour',
  240,
  'completed'::work_session_status,
  20
FROM users u
WHERE u.username = 'manager1'
AND NOT EXISTS (
  SELECT 1 FROM work_session_active w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '5 hours')
);

-- manager2: session paused
INSERT INTO work_session_active (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '1 hour',
  NULL,
  NULL,
  'paused'::work_session_status,
  10
FROM users u
WHERE u.username = 'manager2'
AND NOT EXISTS (
  SELECT 1 FROM work_session_active w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '1 hour')
);

-- archived
INSERT INTO work_session_archived (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '45 days',
  NOW() - INTERVAL '45 days' + INTERVAL '8 hours',
  480,
  'completed'::work_session_status,
  NOW() - INTERVAL '44 days',
  30
FROM users u
WHERE u.username = 'employee1'
AND NOT EXISTS (
  SELECT 1 FROM work_session_archived w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '45 days')
);

INSERT INTO work_session_archived (
  uuid, user_id, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - INTERVAL '90 days',
  NOW() - INTERVAL '90 days' + INTERVAL '7 hours 30 minutes',
  450,
  'completed'::work_session_status,
  NOW() - INTERVAL '89 days',
  25
FROM users u
WHERE u.username = 'employee4'
AND NOT EXISTS (
  SELECT 1 FROM work_session_archived w
  WHERE w.user_id = u.id
    AND w.clock_in = (NOW() - INTERVAL '90 days')
);

-- history (pas de user_id)
INSERT INTO work_session_history (
  uuid, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  NOW() - INTERVAL '3 years',
  NOW() - INTERVAL '3 years' + INTERVAL '8 hours',
  480,
  'completed'::work_session_status,
  NOW() - INTERVAL '3 years' + INTERVAL '1 day',
  35
WHERE NOT EXISTS (
  SELECT 1 FROM work_session_history h
  WHERE h.clock_in = (NOW() - INTERVAL '3 years')
);

INSERT INTO work_session_history (
  uuid, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes
)
SELECT
  gen_random_uuid()::text,
  NOW() - INTERVAL '2 years 6 months',
  NOW() - INTERVAL '2 years 6 months' + INTERVAL '7 hours',
  420,
  'completed'::work_session_status,
  NOW() - INTERVAL '2 years 6 months' + INTERVAL '2 days',
  20
WHERE NOT EXISTS (
  SELECT 1 FROM work_session_history h
  WHERE h.clock_in = (NOW() - INTERVAL '2 years 6 months')
);

-- ------------------------------------------------------------
-- 4) Breaks (liés à work_session_active)
-- ------------------------------------------------------------
-- On s'accroche à la session active via (user + clock_in) au lieu d'un uuid fixe

-- Breaks pour employee1 sur sa session de NOW() - 3 hours
INSERT INTO breaks (uuid, start_time, end_time, duration_minutes, status, work_session_active_id)
SELECT
  gen_random_uuid(),
  (w.clock_in + INTERVAL '1 hour'),
  (w.clock_in + INTERVAL '1 hour 10 minutes'),
  10,
  'completed'::break_status,
  w.id
FROM work_session_active w
JOIN users u ON u.id = w.user_id
WHERE u.username = 'employee1'
  AND w.clock_in = (NOW() - INTERVAL '3 hours')
AND NOT EXISTS (
  SELECT 1 FROM breaks b
  WHERE b.work_session_active_id = w.id
    AND b.start_time = (w.clock_in + INTERVAL '1 hour')
);

INSERT INTO breaks (uuid, start_time, end_time, duration_minutes, status, work_session_active_id)
SELECT
  gen_random_uuid(),
  (w.clock_in + INTERVAL '2 hours'),
  (w.clock_in + INTERVAL '2 hours 5 minutes'),
  5,
  'completed'::break_status,
  w.id
FROM work_session_active w
JOIN users u ON u.id = w.user_id
WHERE u.username = 'employee1'
  AND w.clock_in = (NOW() - INTERVAL '3 hours')
AND NOT EXISTS (
  SELECT 1 FROM breaks b
  WHERE b.work_session_active_id = w.id
    AND b.start_time = (w.clock_in + INTERVAL '2 hours')
);

-- Break actif pour employee2 sur sa session de NOW() - 2 hours
INSERT INTO breaks (uuid, start_time, end_time, duration_minutes, status, work_session_active_id)
SELECT
  gen_random_uuid(),
  NOW() - INTERVAL '20 minutes',
  NULL,
  NULL,
  'active'::break_status,
  w.id
FROM work_session_active w
JOIN users u ON u.id = w.user_id
WHERE u.username = 'employee2'
  AND w.clock_in = (NOW() - INTERVAL '2 hours')
AND NOT EXISTS (
  SELECT 1 FROM breaks b
  WHERE b.work_session_active_id = w.id
    AND b.status = 'active'::break_status
);
-- ------------------------------------------------------------
-- 5) Massive work sessions generation (50 per user per table)
-- ------------------------------------------------------------

-- Generate 50 work_session_active per user (last 30 days)
INSERT INTO work_session_active (uuid, user_id, clock_in, clock_out, duration_minutes, status, breaks_duration_minutes)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - (i || ' days')::interval - ((8 + (random() * 2)::int) || ' hours')::interval,
  NOW() - (i || ' days')::interval - ((random() * 2)::int || ' hours')::interval,
  420 + (random() * 120)::int,
  'completed'::work_session_status,
  15 + (random() * 30)::int
FROM users u
CROSS JOIN generate_series(1, 50) AS i
WHERE u.status = 'active';

-- Generate 50 work_session_archived per user (30 days to 2 years ago)
INSERT INTO work_session_archived (uuid, user_id, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes)
SELECT
  gen_random_uuid()::text,
  u.id,
  NOW() - ((31 + i * 5) || ' days')::interval - ((8 + (random() * 2)::int) || ' hours')::interval,
  NOW() - ((31 + i * 5) || ' days')::interval - ((random() * 2)::int || ' hours')::interval,
  400 + (random() * 140)::int,
  'completed'::work_session_status,
  NOW() - ((30 + i * 5) || ' days')::interval,
  10 + (random() * 35)::int
FROM users u
CROSS JOIN generate_series(1, 50) AS i
WHERE u.status = 'active';

-- Generate 50 work_session_history entries (more than 2 years ago, no user_id for RGPD)
INSERT INTO work_session_history (uuid, clock_in, clock_out, duration_minutes, status, archived_at, breaks_duration_minutes)
SELECT
  gen_random_uuid()::text,
  NOW() - ((730 + i * 7) || ' days')::interval - ((8 + (random() * 2)::int) || ' hours')::interval,
  NOW() - ((730 + i * 7) || ' days')::interval - ((random() * 2)::int || ' hours')::interval,
  390 + (random() * 150)::int,
  'completed'::work_session_status,
  NOW() - ((729 + i * 7) || ' days')::interval,
  10 + (random() * 40)::int
FROM generate_series(1, 500) AS i;