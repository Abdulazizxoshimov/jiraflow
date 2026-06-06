-- =============================================================================
-- TEST SEED DATA — turli rollar bilan API test uchun mock ma'lumotlar
-- =============================================================================
-- Foydalanuvchilar va parollar:
--   admin@jiraflow.com   / Admin123!    → global admin
--   member1@jiraflow.com / Member123!   → global member, WEBDEV admin
--   member2@jiraflow.com / Member123!   → global member, WEBDEV member
--   pm@jiraflow.com      / Manager123!  → global member, MOBILE admin
--   dev@jiraflow.com     / Dev123!pass  → global member, WEBDEV+MOBILE member
--   qa@jiraflow.com      / Qa123!pass1  → global member, WEBDEV viewer
--   viewer@jiraflow.com  / Viewer123!   → global viewer
-- =============================================================================

BEGIN;

-- =============================================================================
-- 1. FOYDALANUVCHILAR
-- =============================================================================
INSERT INTO users (id, email, password_hash, full_name, role, timezone, language, is_active) VALUES
  ('b1000000-0000-0000-0000-000000000001', 'member1@jiraflow.com',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Bobur Toshmatov', 'member', 'Asia/Tashkent', 'uz', TRUE),
  ('b2000000-0000-0000-0000-000000000002', 'member2@jiraflow.com',
   '$2a$12$PB4hGip9Eq1PT.GyEalnXeOWlFX8HS8NhQ13rUjOJ4C1UWc2wGdhi',
   'Zulfiya Rahimova', 'member', 'Asia/Tashkent', 'uz', TRUE),
  ('b3000000-0000-0000-0000-000000000003', 'pm@jiraflow.com',
   '$2a$12$tgjh.L2O2O5QaA1EeONl.OQwmHDRqePTX0yKSXQss9RZHTQNN97g.',
   'Jasur Mirzayev', 'member', 'Asia/Tashkent', 'en', TRUE),
  ('b4000000-0000-0000-0000-000000000004', 'dev@jiraflow.com',
   '$2a$12$mpAmNILdA/IIxjn2WQQvz..Na9BzJqnaHKJEp67XOcKgCqNb3YJsa',
   'Dilnoza Yusupova', 'member', 'UTC', 'en', TRUE),
  ('b5000000-0000-0000-0000-000000000005', 'qa@jiraflow.com',
   '$2a$12$W7KppJwtCiWdkMtmV1WdU.N3Jpdv6bDUnD/yx6wmA0/QikKq6MsUi',
   'Sherzod Karimov', 'member', 'UTC', 'en', TRUE),
  ('b6000000-0000-0000-0000-000000000006', 'viewer@jiraflow.com',
   '$2a$12$oALJC4cu7.wVh3SBILpEdOaakZZco9X/.MQKzpX8cBAu3AvABICkW',
   'Nodira Hasanova', 'viewer', 'Asia/Tashkent', 'uz', TRUE)
ON CONFLICT (email) DO NOTHING;

-- =============================================================================
-- 2. LOYIHALAR
-- =============================================================================
INSERT INTO projects (id, key, name, description, lead_id, workflow_id, issue_counter, is_archived) VALUES
  ('c1000000-0000-0000-0000-000000000001',
   'WEBDEV', 'Web Development Platform',
   'Kompaniyaning asosiy veb platformasini yaratish va boshqarish',
   'b1000000-0000-0000-0000-000000000001',
   '00000000-0000-0000-0000-000000000001',
   0, FALSE),
  ('c2000000-0000-0000-0000-000000000002',
   'MOBILE', 'Mobile Application',
   'iOS va Android uchun cross-platform mobil ilova',
   'b3000000-0000-0000-0000-000000000003',
   '00000000-0000-0000-0000-000000000001',
   0, FALSE)
ON CONFLICT (key) DO NOTHING;

-- =============================================================================
-- 3. PROJECT MEMBERS
-- =============================================================================
-- WEBDEV: admin→admin, member1→admin, member2→member, dev→member, qa→viewer, viewer→viewer
-- MOBILE: admin→admin, pm→admin, dev→member, qa→member, member1→viewer
INSERT INTO project_members (project_id, user_id, role) VALUES
  ('c1000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('c1000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', 'admin'),
  ('c1000000-0000-0000-0000-000000000001', 'b2000000-0000-0000-0000-000000000002', 'member'),
  ('c1000000-0000-0000-0000-000000000001', 'b4000000-0000-0000-0000-000000000004', 'member'),
  ('c1000000-0000-0000-0000-000000000001', 'b5000000-0000-0000-0000-000000000005', 'viewer'),
  ('c1000000-0000-0000-0000-000000000001', 'b6000000-0000-0000-0000-000000000006', 'viewer'),
  ('c2000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('c2000000-0000-0000-0000-000000000002', 'b3000000-0000-0000-0000-000000000003', 'admin'),
  ('c2000000-0000-0000-0000-000000000002', 'b4000000-0000-0000-0000-000000000004', 'member'),
  ('c2000000-0000-0000-0000-000000000002', 'b5000000-0000-0000-0000-000000000005', 'member'),
  ('c2000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001', 'viewer')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 4. SPRINTLAR
-- =============================================================================
INSERT INTO sprints (id, project_id, name, goal, status, start_date, end_date, started_at, created_by) VALUES
  ('d1000000-0000-0000-0000-000000000001',
   'c1000000-0000-0000-0000-000000000001',
   'Sprint 1 — Authentication & Users',
   'Foydalanuvchi autentifikatsiyasi va profil sahifasini yaratish',
   'active', '2026-05-01', '2026-05-14', '2026-05-01 09:00:00+05',
   'b1000000-0000-0000-0000-000000000001'),
  ('d2000000-0000-0000-0000-000000000002',
   'c1000000-0000-0000-0000-000000000001',
   'Sprint 2 — Dashboard & Projects',
   'Dashboard va loyiha boshqarish modulini yaratish',
   'planned', '2026-05-15', '2026-05-28', NULL,
   'b1000000-0000-0000-0000-000000000001'),
  ('d3000000-0000-0000-0000-000000000003',
   'c2000000-0000-0000-0000-000000000002',
   'Sprint 1 — App Setup & Onboarding',
   'Asosiy ilova tuzilmasi va onboarding ekranlari',
   'completed', '2026-04-15', '2026-04-28', '2026-04-15 09:00:00+05',
   'b3000000-0000-0000-0000-000000000003'),
  ('d4000000-0000-0000-0000-000000000004',
   'c2000000-0000-0000-0000-000000000002',
   'Sprint 2 — Core Features',
   'Asosiy funksiyalar: login, home, profile',
   'active', '2026-05-01', '2026-05-14', '2026-05-01 09:00:00+05',
   'b3000000-0000-0000-0000-000000000003')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 5. LABELLAR
-- =============================================================================
INSERT INTO labels (id, project_id, name, color) VALUES
  ('e1000000-0000-0000-0000-000000000001', 'c1000000-0000-0000-0000-000000000001', 'frontend',    '#3B82F6'),
  ('e2000000-0000-0000-0000-000000000002', 'c1000000-0000-0000-0000-000000000001', 'backend',     '#10B981'),
  ('e3000000-0000-0000-0000-000000000003', 'c1000000-0000-0000-0000-000000000001', 'urgent',      '#EF4444'),
  ('e4000000-0000-0000-0000-000000000004', 'c1000000-0000-0000-0000-000000000001', 'database',    '#8B5CF6'),
  ('e5000000-0000-0000-0000-000000000005', 'c1000000-0000-0000-0000-000000000001', 'design',      '#F59E0B'),
  ('e6000000-0000-0000-0000-000000000006', 'c2000000-0000-0000-0000-000000000002', 'ios',         '#6B7280'),
  ('e7000000-0000-0000-0000-000000000007', 'c2000000-0000-0000-0000-000000000002', 'android',     '#22C55E'),
  ('e8000000-0000-0000-0000-000000000008', 'c2000000-0000-0000-0000-000000000002', 'crash',       '#EF4444'),
  ('e9000000-0000-0000-0000-000000000009', 'c2000000-0000-0000-0000-000000000002', 'performance', '#F97316')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 6. CUSTOM FIELDS
-- =============================================================================
INSERT INTO custom_fields (id, project_id, name, field_key, field_type, is_required, options, position) VALUES
  ('f1000000-0000-0000-0000-000000000001',
   'c1000000-0000-0000-0000-000000000001',
   'Browser', 'browser', 'select', FALSE,
   '["Chrome","Firefox","Safari","Edge","Opera"]'::jsonb, 1),
  ('f2000000-0000-0000-0000-000000000002',
   'c1000000-0000-0000-0000-000000000001',
   'Environment', 'environment', 'select', TRUE,
   '["development","staging","production"]'::jsonb, 2),
  ('f3000000-0000-0000-0000-000000000003',
   'c1000000-0000-0000-0000-000000000001',
   'Estimated Hours', 'estimated_hours', 'number', FALSE, NULL, 3),
  ('f4000000-0000-0000-0000-000000000004',
   'c2000000-0000-0000-0000-000000000002',
   'Platform', 'platform', 'multi_select', FALSE,
   '["iOS","Android","Both"]'::jsonb, 1),
  ('f5000000-0000-0000-0000-000000000005',
   'c2000000-0000-0000-0000-000000000002',
   'App Version', 'app_version', 'text', FALSE, NULL, 2),
  ('f6000000-0000-0000-0000-000000000006',
   'c2000000-0000-0000-0000-000000000002',
   'Device', 'device', 'text', FALSE, NULL, 3)
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 7. ISSUES — WEBDEV (15 ta)
-- Workflow status UUID lar:
--   433c09ee-9223-41d3-afeb-df59a2336531 = To Do
--   3d1c059c-c5be-415c-b848-773ba5f5fc71 = In Progress
--   7a1f6e46-05d5-4d07-b50d-18c6e39a7040 = In Review
--   6703b413-4fbf-4577-9b6f-2b5b994ad00d = Done
-- =============================================================================
INSERT INTO issues (id, project_id, issue_number, title, description, type, status_id, priority,
                    assignee_id, reporter_id, sprint_id, story_points, due_date, custom_fields) VALUES
  -- EPIClar
  ('10000001-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 1,
   'Authentication System',
   'Foydalanuvchi autentifikatsiya tizimini toliq yaratish: login, register, JWT, refresh',
   'epic', '3d1c059c-c5be-415c-b848-773ba5f5fc71', 'high',
   'b1000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 21, '2026-05-14',
   '{"environment": "development"}'::jsonb),

  ('10000002-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 2,
   'Dashboard MVP',
   'Asosiy dashboard sahifasi: statistikalar, grafik, songi faoliyat',
   'epic', '433c09ee-9223-41d3-afeb-df59a2336531', 'medium',
   'b2000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
   'd2000000-0000-0000-0000-000000000002', 34, '2026-05-28',
   '{"environment": "development"}'::jsonb),

  -- STORYlar
  ('10000003-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 3,
   'As a user, I can login with email and password',
   'Login sahifasi UI va backend integratsiyasi',
   'story', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'highest',
   'b2000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 8, '2026-05-07',
   '{"environment": "staging", "browser": "Chrome"}'::jsonb),

  ('10000004-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 4,
   'As a user, I can reset my forgotten password',
   'Email orqali parol tiklash: token yuborish va yangi parol ortnatish',
   'story', '7a1f6e46-05d5-4d07-b50d-18c6e39a7040', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 5, '2026-05-10',
   '{"environment": "development"}'::jsonb),

  -- TASKlar
  ('10000005-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 5,
   'Create JWT middleware for authentication',
   'Bearer token validatsiyasi uchun Gin middleware yozish',
   'task', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 3, '2026-05-05',
   '{"environment": "development", "estimated_hours": 8}'::jsonb),

  ('10000006-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 6,
   'Design login page UI mockup',
   'Figmada login sahifasi dizayni: light/dark mode, responsive',
   'task', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'medium',
   'b2000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 2, '2026-05-04',
   '{"environment": "development", "browser": "Chrome", "estimated_hours": 4}'::jsonb),

  ('10000007-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 7,
   'Set up PostgreSQL database schema',
   'Migrations yozish va test DB sozlash',
   'task', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'highest',
   'b4000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 5, '2026-05-03',
   '{"environment": "development", "estimated_hours": 12}'::jsonb),

  ('10000008-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 8,
   'Implement refresh token rotation',
   'Refresh tokenni yangilash va eski tokenni invalid qilish',
   'task', '7a1f6e46-05d5-4d07-b50d-18c6e39a7040', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 3, '2026-05-08',
   '{"environment": "staging", "estimated_hours": 6}'::jsonb),

  ('10000009-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 9,
   'Add CORS configuration',
   'Frontend origin larni CORS orqali ruxsat berish',
   'task', '433c09ee-9223-41d3-afeb-df59a2336531', 'low',
   NULL, 'b1000000-0000-0000-0000-000000000001',
   'd2000000-0000-0000-0000-000000000002', 1, NULL,
   '{"environment": "development"}'::jsonb),

  ('1000000a-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 10,
   'Write API documentation for auth endpoints',
   'Swagger annotatsiyalari va README qoshish',
   'task', '433c09ee-9223-41d3-afeb-df59a2336531', 'low',
   'b2000000-0000-0000-0000-000000000002', 'b1000000-0000-0000-0000-000000000001',
   'd2000000-0000-0000-0000-000000000002', 2, NULL,
   '{"environment": "development", "estimated_hours": 3}'::jsonb),

  -- BUGlar
  ('1000000b-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 11,
   'Login button not working on mobile Safari',
   'iPhone 15 Safari da login tugmasi bosilganda hech narsa bolmaydi',
   'bug', '3d1c059c-c5be-415c-b848-773ba5f5fc71', 'highest',
   'b2000000-0000-0000-0000-000000000002', 'b5000000-0000-0000-0000-000000000005',
   'd1000000-0000-0000-0000-000000000001', NULL, '2026-05-06',
   '{"browser": "Safari", "environment": "staging"}'::jsonb),

  ('1000000c-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 12,
   'Token expiry not handled gracefully',
   'Access token tugaganda 401 orniga app crash bolyapti',
   'bug', '3d1c059c-c5be-415c-b848-773ba5f5fc71', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b5000000-0000-0000-0000-000000000005',
   'd1000000-0000-0000-0000-000000000001', NULL, '2026-05-07',
   '{"browser": "Chrome", "environment": "production"}'::jsonb),

  ('1000000d-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 13,
   'Password reset email not sent in production',
   'SMTP konfiguratsiyasi productionda ishlamayapti',
   'bug', '433c09ee-9223-41d3-afeb-df59a2336531', 'high',
   NULL, 'b5000000-0000-0000-0000-000000000005',
   'd1000000-0000-0000-0000-000000000001', NULL, '2026-05-09',
   '{"environment": "production"}'::jsonb),

  -- SUBTASKlar (parent: Authentication System epic)
  ('1000000e-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 14,
   'Write unit tests for auth service',
   'Auth usecase uchun unit testlar yozish',
   'subtask', '433c09ee-9223-41d3-afeb-df59a2336531', 'medium',
   'b4000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 3, NULL,
   '{"environment": "development", "estimated_hours": 4}'::jsonb),

  ('1000000f-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001', 15,
   'Set up Redis for session management',
   'Redisda refresh token va session malumotlarini saqlash',
   'subtask', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b1000000-0000-0000-0000-000000000001',
   'd1000000-0000-0000-0000-000000000001', 2, '2026-05-04',
   '{"environment": "development"}'::jsonb)
ON CONFLICT DO NOTHING;

-- Parent - child bog'lanishlar (WEBDEV)
UPDATE issues SET parent_id = '10000001-0000-0000-0000-000000000000'
WHERE id IN (
  '10000003-0000-0000-0000-000000000000',
  '10000004-0000-0000-0000-000000000000',
  '10000005-0000-0000-0000-000000000000',
  '10000006-0000-0000-0000-000000000000',
  '10000007-0000-0000-0000-000000000000',
  '10000008-0000-0000-0000-000000000000',
  '1000000e-0000-0000-0000-000000000000',
  '1000000f-0000-0000-0000-000000000000'
);

UPDATE projects SET issue_counter = 15 WHERE id = 'c1000000-0000-0000-0000-000000000001';

-- =============================================================================
-- 8. ISSUES — MOBILE (10 ta)
-- =============================================================================
INSERT INTO issues (id, project_id, issue_number, title, description, type, status_id, priority,
                    assignee_id, reporter_id, sprint_id, story_points, due_date, custom_fields) VALUES
  ('20000001-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 1,
   'App Architecture Setup',
   'Flutter loyihasi tuzilmasi: Clean Architecture, BLoC pattern, DI',
   'epic', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'highest',
   'b3000000-0000-0000-0000-000000000003', 'b3000000-0000-0000-0000-000000000003',
   'd3000000-0000-0000-0000-000000000003', 13, '2026-04-28',
   '{"platform": ["iOS", "Android"]}'::jsonb),

  ('20000002-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 2,
   'Login screen implementation',
   'Email/parol bilan kirish ekrani: UI va API integratsiyasi',
   'story', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'highest',
   'b4000000-0000-0000-0000-000000000004', 'b3000000-0000-0000-0000-000000000003',
   'd4000000-0000-0000-0000-000000000004', 8, '2026-05-05',
   '{"platform": ["Both"], "app_version": "1.0.0"}'::jsonb),

  ('20000003-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 3,
   'Implement biometric authentication',
   'Face ID va fingerprint orqali kirish imkoniyati',
   'task', '3d1c059c-c5be-415c-b848-773ba5f5fc71', 'high',
   'b4000000-0000-0000-0000-000000000004', 'b3000000-0000-0000-0000-000000000003',
   'd4000000-0000-0000-0000-000000000004', 5, '2026-05-10',
   '{"platform": ["iOS", "Android"], "app_version": "1.0.0"}'::jsonb),

  ('20000004-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 4,
   'Home screen — project list',
   'Foydalanuvchi loyihalarini royxatda korsatish: filter, search, sort',
   'task', '433c09ee-9223-41d3-afeb-df59a2336531', 'medium',
   'b5000000-0000-0000-0000-000000000005', 'b3000000-0000-0000-0000-000000000003',
   'd4000000-0000-0000-0000-000000000004', 3, '2026-05-12',
   '{"platform": ["Both"], "app_version": "1.0.0"}'::jsonb),

  ('20000005-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 5,
   'App crashes on Android 12 during startup',
   'Samsung Galaxy S22 da ilova ochilganda crash bolyapti',
   'bug', '3d1c059c-c5be-415c-b848-773ba5f5fc71', 'highest',
   'b4000000-0000-0000-0000-000000000004', 'b5000000-0000-0000-0000-000000000005',
   'd4000000-0000-0000-0000-000000000004', NULL, '2026-05-06',
   '{"platform": ["Android"], "app_version": "0.9.5", "device": "Samsung Galaxy S22"}'::jsonb),

  ('20000006-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 6,
   'Push notifications not working on iOS',
   'APNs konfiguratsiyasi notogri: notification kelmasayapti',
   'bug', '7a1f6e46-05d5-4d07-b50d-18c6e39a7040', 'high',
   'b3000000-0000-0000-0000-000000000003', 'b5000000-0000-0000-0000-000000000005',
   'd4000000-0000-0000-0000-000000000004', NULL, '2026-05-08',
   '{"platform": ["iOS"], "app_version": "1.0.0", "device": "iPhone 15 Pro"}'::jsonb),

  ('20000007-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 7,
   'Implement offline mode for issue list',
   'Internetsiz holatda cache dan issue list korsatish',
   'task', '433c09ee-9223-41d3-afeb-df59a2336531', 'medium',
   NULL, 'b3000000-0000-0000-0000-000000000003',
   NULL, 8, '2026-05-20',
   '{"platform": ["Both"]}'::jsonb),

  ('20000008-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 8,
   'Dark mode support',
   'Ilova boylab dark/light mode toggle qoshish',
   'task', '433c09ee-9223-41d3-afeb-df59a2336531', 'low',
   NULL, 'b3000000-0000-0000-0000-000000000003',
   NULL, 5, NULL,
   '{"platform": ["Both"]}'::jsonb),

  ('20000009-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 9,
   'Performance: slow image loading',
   'Profil rasmlari yuklanishi 3-4 soniya olayapti, optimizatsiya kerak',
   'bug', '433c09ee-9223-41d3-afeb-df59a2336531', 'medium',
   'b4000000-0000-0000-0000-000000000004', 'b5000000-0000-0000-0000-000000000005',
   'd4000000-0000-0000-0000-000000000004', NULL, NULL,
   '{"platform": ["Both"], "app_version": "1.0.0"}'::jsonb),

  ('2000000a-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002', 10,
   'Onboarding flow — 3 steps',
   'Yangi foydalanuvchilar uchun 3 bosqichli onboarding ekranlar',
   'story', '6703b413-4fbf-4577-9b6f-2b5b994ad00d', 'medium',
   'b4000000-0000-0000-0000-000000000004', 'b3000000-0000-0000-0000-000000000003',
   'd3000000-0000-0000-0000-000000000003', 5, '2026-04-25',
   '{"platform": ["Both"], "app_version": "0.9.0"}'::jsonb)
ON CONFLICT DO NOTHING;

UPDATE projects SET issue_counter = 10 WHERE id = 'c2000000-0000-0000-0000-000000000002';

-- =============================================================================
-- 9. ISSUE_LABELS
-- =============================================================================
INSERT INTO issue_labels (issue_id, label_id) VALUES
  ('10000005-0000-0000-0000-000000000000', 'e2000000-0000-0000-0000-000000000002'),
  ('10000006-0000-0000-0000-000000000000', 'e1000000-0000-0000-0000-000000000001'),
  ('10000006-0000-0000-0000-000000000000', 'e5000000-0000-0000-0000-000000000005'),
  ('10000007-0000-0000-0000-000000000000', 'e2000000-0000-0000-0000-000000000002'),
  ('10000007-0000-0000-0000-000000000000', 'e4000000-0000-0000-0000-000000000004'),
  ('1000000b-0000-0000-0000-000000000000', 'e3000000-0000-0000-0000-000000000003'),
  ('1000000b-0000-0000-0000-000000000000', 'e1000000-0000-0000-0000-000000000001'),
  ('1000000c-0000-0000-0000-000000000000', 'e3000000-0000-0000-0000-000000000003'),
  ('1000000c-0000-0000-0000-000000000000', 'e2000000-0000-0000-0000-000000000002'),
  ('20000005-0000-0000-0000-000000000000', 'e8000000-0000-0000-0000-000000000008'),
  ('20000005-0000-0000-0000-000000000000', 'e7000000-0000-0000-0000-000000000007'),
  ('20000006-0000-0000-0000-000000000000', 'e6000000-0000-0000-0000-000000000006'),
  ('20000009-0000-0000-0000-000000000000', 'e9000000-0000-0000-0000-000000000009')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 10. ISSUE_WATCHERS
-- =============================================================================
INSERT INTO issue_watchers (issue_id, user_id) VALUES
  ('1000000b-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001'),
  ('1000000b-0000-0000-0000-000000000000', 'b1000000-0000-0000-0000-000000000001'),
  ('1000000c-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001'),
  ('1000000c-0000-0000-0000-000000000000', 'b5000000-0000-0000-0000-000000000005'),
  ('20000005-0000-0000-0000-000000000000', 'b3000000-0000-0000-0000-000000000003'),
  ('20000005-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 11. ISSUE_LINKS
-- =============================================================================
INSERT INTO issue_links (source_id, target_id, link_type, created_by) VALUES
  ('1000000b-0000-0000-0000-000000000000', '1000000c-0000-0000-0000-000000000000',
   'relates_to', 'b5000000-0000-0000-0000-000000000005'),
  ('20000005-0000-0000-0000-000000000000', '20000009-0000-0000-0000-000000000000',
   'relates_to', 'b3000000-0000-0000-0000-000000000003'),
  ('10000003-0000-0000-0000-000000000000', '1000000b-0000-0000-0000-000000000000',
   'blocks', 'b1000000-0000-0000-0000-000000000001')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 12. BOARDLAR
-- =============================================================================
INSERT INTO boards (id, project_id, name, type, filter, created_by) VALUES
  ('30000001-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001',
   'WEBDEV Kanban Board', 'kanban',
   '{"sprint_id": null}'::jsonb,
   'b1000000-0000-0000-0000-000000000001'),
  ('30000002-0000-0000-0000-000000000000',
   'c1000000-0000-0000-0000-000000000001',
   'Sprint 1 Scrum Board', 'scrum',
   '{"sprint_id": "d1000000-0000-0000-0000-000000000001"}'::jsonb,
   'b1000000-0000-0000-0000-000000000001'),
  ('30000003-0000-0000-0000-000000000000',
   'c2000000-0000-0000-0000-000000000002',
   'MOBILE Kanban Board', 'kanban',
   '{"sprint_id": null}'::jsonb,
   'b3000000-0000-0000-0000-000000000003')
ON CONFLICT DO NOTHING;

INSERT INTO board_columns (id, board_id, name, position, wip_limit) VALUES
  ('40000001-0000-0000-0000-000000000000', '30000001-0000-0000-0000-000000000000', 'To Do',       1, NULL),
  ('40000002-0000-0000-0000-000000000000', '30000001-0000-0000-0000-000000000000', 'In Progress',  2, 3),
  ('40000003-0000-0000-0000-000000000000', '30000001-0000-0000-0000-000000000000', 'In Review',    3, 2),
  ('40000004-0000-0000-0000-000000000000', '30000001-0000-0000-0000-000000000000', 'Done',         4, NULL),
  ('40000005-0000-0000-0000-000000000000', '30000003-0000-0000-0000-000000000000', 'To Do',        1, NULL),
  ('40000006-0000-0000-0000-000000000000', '30000003-0000-0000-0000-000000000000', 'In Progress',  2, 5),
  ('40000007-0000-0000-0000-000000000000', '30000003-0000-0000-0000-000000000000', 'In Review',    3, 2),
  ('40000008-0000-0000-0000-000000000000', '30000003-0000-0000-0000-000000000000', 'Done',         4, NULL)
ON CONFLICT DO NOTHING;

-- Board column → status mapping (board_column_statuses)
INSERT INTO board_column_statuses (column_id, status_id) VALUES
  ('40000001-0000-0000-0000-000000000000', '433c09ee-9223-41d3-afeb-df59a2336531'),
  ('40000002-0000-0000-0000-000000000000', '3d1c059c-c5be-415c-b848-773ba5f5fc71'),
  ('40000003-0000-0000-0000-000000000000', '7a1f6e46-05d5-4d07-b50d-18c6e39a7040'),
  ('40000004-0000-0000-0000-000000000000', '6703b413-4fbf-4577-9b6f-2b5b994ad00d'),
  ('40000005-0000-0000-0000-000000000000', '433c09ee-9223-41d3-afeb-df59a2336531'),
  ('40000006-0000-0000-0000-000000000000', '3d1c059c-c5be-415c-b848-773ba5f5fc71'),
  ('40000007-0000-0000-0000-000000000000', '7a1f6e46-05d5-4d07-b50d-18c6e39a7040'),
  ('40000008-0000-0000-0000-000000000000', '6703b413-4fbf-4577-9b6f-2b5b994ad00d')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 13. SPACES
-- =============================================================================
INSERT INTO spaces (id, key, name, description, type, lead_id, project_id, is_archived) VALUES
  ('50000001-0000-0000-0000-000000000000',
   'WEBDOC', 'WEBDEV Documentation',
   'Web Development loyihasining texnik dokumentatsiyasi',
   'project', 'b1000000-0000-0000-0000-000000000001',
   'c1000000-0000-0000-0000-000000000001', FALSE),
  ('50000002-0000-0000-0000-000000000000',
   'MOBDOC', 'Mobile App Documentation',
   'Mobile ilova arxitekturasi va API dokumentatsiyasi',
   'project', 'b3000000-0000-0000-0000-000000000003',
   'c2000000-0000-0000-0000-000000000002', FALSE),
  ('50000003-0000-0000-0000-000000000000',
   'TEAMWIKI', 'Team Space',
   'Umumiy jamoa uchun wiki: onboarding, qoidalar, jarayonlar',
   'team', 'a0000000-0000-0000-0000-000000000001',
   NULL, FALSE)
ON CONFLICT (key) DO NOTHING;

INSERT INTO space_members (space_id, user_id, role) VALUES
  ('50000001-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('50000001-0000-0000-0000-000000000000', 'b1000000-0000-0000-0000-000000000001', 'admin'),
  ('50000001-0000-0000-0000-000000000000', 'b2000000-0000-0000-0000-000000000002', 'member'),
  ('50000001-0000-0000-0000-000000000000', 'b4000000-0000-0000-0000-000000000004', 'member'),
  ('50000001-0000-0000-0000-000000000000', 'b5000000-0000-0000-0000-000000000005', 'viewer'),
  ('50000002-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('50000002-0000-0000-0000-000000000000', 'b3000000-0000-0000-0000-000000000003', 'admin'),
  ('50000002-0000-0000-0000-000000000000', 'b4000000-0000-0000-0000-000000000004', 'member'),
  ('50000003-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('50000003-0000-0000-0000-000000000000', 'b1000000-0000-0000-0000-000000000001', 'member'),
  ('50000003-0000-0000-0000-000000000000', 'b3000000-0000-0000-0000-000000000003', 'member')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 14. PAGES
-- =============================================================================
INSERT INTO pages (id, space_id, title, content, content_text, author_id, last_editor_id, parent_id, position, status) VALUES
  ('60000001-0000-0000-0000-000000000000',
   '50000001-0000-0000-0000-000000000000',
   'Getting Started',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"WEBDEV loyihasiga xush kelibsiz. Bu sahifada boshlash uchun zarur malumotlar bor."}]}]}'::jsonb,
   'WEBDEV loyihasiga xush kelibsiz.',
   'b1000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', NULL, 1, 'published'),

  ('60000002-0000-0000-0000-000000000000',
   '50000001-0000-0000-0000-000000000000',
   'API Documentation',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Backend API endpointlari royxati va ularning tavsifi."}]}]}'::jsonb,
   'Backend API endpointlari royxati va ularning tavsifi.',
   'b1000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001',
   '60000001-0000-0000-0000-000000000000', 1, 'published'),

  ('60000003-0000-0000-0000-000000000000',
   '50000001-0000-0000-0000-000000000000',
   'Database Schema',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"PostgreSQL jadvallar tuzilmasi va munosabatlar diagrammasi."}]}]}'::jsonb,
   'PostgreSQL jadvallar tuzilmasi.',
   'b2000000-0000-0000-0000-000000000002', 'b2000000-0000-0000-0000-000000000002',
   '60000001-0000-0000-0000-000000000000', 2, 'published'),

  ('60000004-0000-0000-0000-000000000000',
   '50000002-0000-0000-0000-000000000000',
   'Mobile Architecture',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Flutter Clean Architecture, BLoC pattern va folder tuzilmasi."}]}]}'::jsonb,
   'Flutter Clean Architecture va BLoC pattern.',
   'b3000000-0000-0000-0000-000000000003', 'b3000000-0000-0000-0000-000000000003', NULL, 1, 'published'),

  ('60000005-0000-0000-0000-000000000000',
   '50000003-0000-0000-0000-000000000000',
   'Team Onboarding',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Yangi jamoa azosi uchun boshlash qollanmasi."}]}]}'::jsonb,
   'Yangi jamoa azosi uchun boshlash qollanmasi.',
   'a0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', NULL, 1, 'published')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 15. COMMENTS
-- =============================================================================
INSERT INTO comments (id, parent_type, parent_id, author_id, content, content_text) VALUES
  ('70000001-0000-0000-0000-000000000000',
   'issue', '1000000b-0000-0000-0000-000000000000',
   'b5000000-0000-0000-0000-000000000005',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"iPhone 15 Safari 17.3 da tasdiqlangan. Login tugmasi click eventini ushlamayapti."}]}]}'::jsonb,
   'iPhone 15 Safari 17.3 da tasdiqlangan. Login tugmasi click eventini ushlamayapti.'),

  ('70000002-0000-0000-0000-000000000000',
   'issue', '1000000b-0000-0000-0000-000000000000',
   'b2000000-0000-0000-0000-000000000002',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Muammoni topdim: form submit da preventDefault() chaqirilmayapti. Fix qilyapman."}]}]}'::jsonb,
   'Muammoni topdim: form submit da preventDefault() chaqirilmayapti. Fix qilyapman.'),

  ('70000003-0000-0000-0000-000000000000',
   'issue', '20000005-0000-0000-0000-000000000000',
   'b4000000-0000-0000-0000-000000000004',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Stack trace: NullPointerException in MainActivity.onCreate. Android 12 permission model ozgardi."}]}]}'::jsonb,
   'Stack trace: NullPointerException in MainActivity.onCreate. Android 12 permission model ozgardi.'),

  ('70000004-0000-0000-0000-000000000000',
   'page', '60000001-0000-0000-0000-000000000000',
   'b2000000-0000-0000-0000-000000000002',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Environment setup qismi yangilash kerak, Docker compose versiyasi ozgardi."}]}]}'::jsonb,
   'Environment setup qismi yangilash kerak, Docker compose versiyasi ozgardi.')
ON CONFLICT DO NOTHING;

INSERT INTO comments (id, parent_type, parent_id, author_id, content, content_text, reply_to_id) VALUES
  ('70000005-0000-0000-0000-000000000000',
   'issue', '1000000b-0000-0000-0000-000000000000',
   'b1000000-0000-0000-0000-000000000001',
   '{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Zulfiya rahmat! Pull request ochasanmi?"}]}]}'::jsonb,
   'Zulfiya rahmat! Pull request ochasanmi?',
   '70000002-0000-0000-0000-000000000000')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 16. ISSUE HISTORY
-- =============================================================================
INSERT INTO issue_history (issue_id, user_id, field, old_value, new_value) VALUES
  ('1000000b-0000-0000-0000-000000000000',
   'a0000000-0000-0000-0000-000000000001', 'status',
   '"433c09ee-9223-41d3-afeb-df59a2336531"'::jsonb,
   '"3d1c059c-c5be-415c-b848-773ba5f5fc71"'::jsonb),
  ('1000000b-0000-0000-0000-000000000000',
   'b1000000-0000-0000-0000-000000000001', 'assignee',
   'null'::jsonb,
   '"b2000000-0000-0000-0000-000000000002"'::jsonb),
  ('10000005-0000-0000-0000-000000000000',
   'b4000000-0000-0000-0000-000000000004', 'status',
   '"3d1c059c-c5be-415c-b848-773ba5f5fc71"'::jsonb,
   '"6703b413-4fbf-4577-9b6f-2b5b994ad00d"'::jsonb),
  ('20000005-0000-0000-0000-000000000000',
   'b3000000-0000-0000-0000-000000000003', 'priority',
   '"high"'::jsonb, '"highest"'::jsonb)
ON CONFLICT DO NOTHING;

COMMIT;

-- =============================================================================
-- YAKUNIY HISOBOT
-- =============================================================================
SELECT '=== SEED MUVAFFAQIYATLI ===' as status;
SELECT 'USERS' as entity, count(*) as total FROM users;
SELECT 'PROJECTS' as entity, count(*) as total FROM projects;
SELECT 'SPRINTS' as entity, count(*) as total FROM sprints;
SELECT 'ISSUES' as entity, count(*) as total FROM issues;
SELECT 'LABELS' as entity, count(*) as total FROM labels;
SELECT 'BOARDS' as entity, count(*) as total FROM boards;
SELECT 'SPACES' as entity, count(*) as total FROM spaces;
SELECT 'PAGES' as entity, count(*) as total FROM pages;
SELECT 'COMMENTS' as entity, count(*) as total FROM comments;

SELECT email, role as global_role FROM users ORDER BY role, email;
