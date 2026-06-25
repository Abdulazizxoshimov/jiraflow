-- =============================================================================
-- YLL Project Mock Data Seed
-- Manba: Jira export — project "Yalla App Development" (YLL)
-- Ko'ringan issues: YLL-712 .. YLL-745 (33 ta)
--
-- Ishlatish:
--   psql $DATABASE_URL -f scripts/seed_yll_data.sql
--
-- Talablar:
--   - Migration 000012 (default workflow) va 000013 (admin user) ishga tushgan bo'lsin
--   - Barcha foydalanuvchi paroli: Member123!
-- =============================================================================

BEGIN;

-- =============================================================================
-- 1. FOYDALANUVCHILAR
-- =============================================================================
-- Parol xeshi: Member123!  ($2a$12$4Co.VkCl...)
INSERT INTO users (id, email, password_hash, full_name, role, timezone, language, is_active)
VALUES
  ('c0000001-0000-0000-0000-000000000000',
   'khikmatullo.khakimov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Khikmatullo Khakimov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000002-0000-0000-0000-000000000000',
   'baxodir@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Baxodir', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000003-0000-0000-0000-000000000000',
   'r.zokirov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Rustam Zokirov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000004-0000-0000-0000-000000000000',
   'muhtorxon@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Muhtorxon', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000005-0000-0000-0000-000000000000',
   'ozodbek.kamolov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Ozodbek Kamolov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000006-0000-0000-0000-000000000000',
   'azizbek@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Azizbek', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000007-0000-0000-0000-000000000000',
   'muhammadjon.tokhirov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Muhammadjon Tokhirov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000008-0000-0000-0000-000000000000',
   'mukhammadali.yolbarsbekov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Mukhammadali Yolbarsbekov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c0000009-0000-0000-0000-000000000000',
   'islom@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Islom', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c000000a-0000-0000-0000-000000000000',
   'farrux.ismoilov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Farrux Ismoilov', 'member', 'Asia/Tashkent', 'uz', TRUE),

  ('c000000b-0000-0000-0000-000000000000',
   'abdulvosid.khasanov@yalla.uz',
   '$2a$12$4Co.VkClYuSZlI9temFtDO7Wr2im6sAAaqLZYdhu.I3E2GpudWeMy',
   'Abdulvosid Khasanov', 'member', 'Asia/Tashkent', 'uz', TRUE)

ON CONFLICT (email) DO NOTHING;

-- =============================================================================
-- 2. DEFAULT WORKFLOW'GA QA / VERIFICATION STATUSI QO'SHISH
--    (Jira-da mavjud: To Do, In Progress, Review, QA/Verification, Done)
--    Mavjud default statuslar: To Do(1), In Progress(2), In Review(3), Done(4)
-- =============================================================================
INSERT INTO workflow_statuses (workflow_id, name, category, color, position)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'QA / Verification',
  'in_progress',
  '#8B5CF6',
  5
)
ON CONFLICT (workflow_id, name) DO NOTHING;

-- =============================================================================
-- 3. YLL LOYIHASI
-- =============================================================================
INSERT INTO projects (id, key, name, description, lead_id, workflow_id, issue_counter)
VALUES (
  'd0000001-0000-0000-0000-000000000000',
  'YLL',
  'Yalla App Development',
  'Yalla taxi ilovasi — mijoz, haydovchi va panel ilovalarini qamrab oluvchi asosiy loyiha.',
  'a0000000-0000-0000-0000-000000000001',          -- admin (migration 000013)
  '00000000-0000-0000-0000-000000000001',          -- default workflow
  746                                              -- YLL-745 dan keyingi raqam
)
ON CONFLICT (key) DO UPDATE
  SET issue_counter = GREATEST(projects.issue_counter, EXCLUDED.issue_counter);

-- =============================================================================
-- 4. LOYIHA A'ZOLARI
-- =============================================================================
INSERT INTO project_members (project_id, user_id, role)
VALUES
  ('d0000001-0000-0000-0000-000000000000', 'a0000000-0000-0000-0000-000000000001', 'admin'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000001-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000002-0000-0000-0000-000000000000', 'admin'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000003-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000004-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000005-0000-0000-0000-000000000000', 'admin'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000006-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000007-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000008-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c0000009-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c000000a-0000-0000-0000-000000000000', 'member'),
  ('d0000001-0000-0000-0000-000000000000', 'c000000b-0000-0000-0000-000000000000', 'member')
ON CONFLICT DO NOTHING;

-- =============================================================================
-- 5. ISSUES
--    Jira status → local workflow_status mapping:
--      'To Do' / 'TODO'        → 'To Do'           (category: todo)
--      'In Progress'           → 'In Progress'      (category: in_progress)
--      'Review'                → 'In Review'        (category: in_progress)
--      'QA / Verification'     → 'QA / Verification'(category: in_progress)
--      'Done'                  → 'Done'             (category: done)
--    Jira type → local type:
--      Task → task, Bug → bug, Story → story, Subtask → subtask
-- =============================================================================
DO $$
DECLARE
  wf_id   UUID := '00000000-0000-0000-0000-000000000001';
  proj_id UUID := 'd0000001-0000-0000-0000-000000000000';

  s_todo   UUID;
  s_inprog UUID;
  s_review UUID;
  s_qa     UUID;
  s_done   UUID;

  u_khikmat      CONSTANT UUID := 'c0000001-0000-0000-0000-000000000000';
  u_baxodir      CONSTANT UUID := 'c0000002-0000-0000-0000-000000000000';
  u_zokirov      CONSTANT UUID := 'c0000003-0000-0000-0000-000000000000';
  u_muhtorxon    CONSTANT UUID := 'c0000004-0000-0000-0000-000000000000';
  u_ozodbek      CONSTANT UUID := 'c0000005-0000-0000-0000-000000000000';
  u_azizbek      CONSTANT UUID := 'c0000006-0000-0000-0000-000000000000';
  u_muhammadjon  CONSTANT UUID := 'c0000007-0000-0000-0000-000000000000';
  u_mukhammadali CONSTANT UUID := 'c0000008-0000-0000-0000-000000000000';
  u_islom        CONSTANT UUID := 'c0000009-0000-0000-0000-000000000000';
  u_farrux       CONSTANT UUID := 'c000000a-0000-0000-0000-000000000000';
  u_abdulvosid   CONSTANT UUID := 'c000000b-0000-0000-0000-000000000000';
BEGIN
  SELECT id INTO s_todo   FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'To Do';
  SELECT id INTO s_inprog FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'In Progress';
  SELECT id INTO s_review FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'In Review';
  SELECT id INTO s_qa     FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'QA / Verification';
  SELECT id INTO s_done   FROM workflow_statuses WHERE workflow_id = wf_id AND name = 'Done';

  IF s_todo IS NULL OR s_inprog IS NULL OR s_review IS NULL OR s_qa IS NULL OR s_done IS NULL THEN
    RAISE EXCEPTION 'Workflow statuslari topilmadi. Migration 000012 ishga tushganligini tekshiring.';
  END IF;

  INSERT INTO issues (
    id, project_id, issue_number,
    title, description,
    type, status_id, priority,
    assignee_id, reporter_id,
    created_at, updated_at
  )
  VALUES

  -- YLL-745 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 745,
   'Promokod berayotgan va supportda otirgan qizlar uchun sms yuborishni avtomat qilish kerak',
   E'## Maqsad\nCall-markaz operatorlari mijozlarga qo''ng''iroq qilib, taksi xizmatidan foydalanmaslik sabablarini aniqlashadi. Agar mijoz muvaffaqiyatli oprosdan o''tsa va operator promokod berishga qaror qilsa, bir tugma orqali SMS yuborsin.\n\n## User Story\nCall-markaz operatori sifatida, men qo''ng''iroq yakunlangandan keyin mijozga bir tugma orqali promokodli SMS yubormoqchiman.\n\n## Acceptance Criteria\n- Operator bitta tugma orqali SMS yubora oladi\n- SMS matni shablon asosida shakllanadi\n- Har bir yuborilgan SMS audit logda saqlanadi\n- Brandlar uchun alohida SMS shablonlari qo''llab-quvvatlanadi',
   'task', s_todo, 'medium',
   NULL, u_khikmat,
   '2026-06-11T15:53:51.887+0500'::timestamptz,
   '2026-06-11T15:58:28.587+0500'::timestamptz),

  -- YLL-744 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 744,
   'socket: extend safe.Recover to remaining worker goroutines and confirm Sentry capture',
   E'**Service:** `go-socket-service`\n\n## Problem\nThe critical per-connection goroutines already run under `pkg/safe` — the read/write pumps are covered. Remaining goroutines have no recover:\n- `internal/ws/hub.go:132,138` — register/unregister worker loops\n- `internal/ws/sharded_clients.go:145,179,212,251` — parallel shard-scan fan-out\n- `pkg/nats_client/client.go:33` — consumer worker-pool worker\n\nAlso: confirm `safe.Recover` reports recovered panics to Sentry.\n\n## Proposed fix\n- Wrap remaining worker/fan-out goroutines with `safe.Go` / `defer safe.Recover(...)`\n- Ensure `pkg/safe.Recover` captures to Sentry',
   'subtask', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T15:11:25.306+0500'::timestamptz,
   '2026-06-11T15:11:25.386+0500'::timestamptz),

  -- YLL-743 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 743,
   'panel: add recover helper, cover worker goroutines, add Sentry to existing recovers',
   E'**Service:** `go-panel-service`\n\n## Problem\nPer-request fire-and-forget goroutines already recover inline but **log only, with no Sentry capture**. Other goroutines have no recover:\n- `internal/queue/worker.go:59` — queue server goroutine\n- `internal/adapters/repository/postgres/pool.go` — pool-stats monitor\n\n## Proposed fix\n- Add `pkg/goroutine/goroutine.Go` per the YLL-253 spec (recover + stack + Sentry)\n- Route fire-and-forget goroutines through it\n- Add Sentry capture to existing inline recovers',
   'subtask', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T15:11:12.582+0500'::timestamptz,
   '2026-06-11T15:11:12.728+0500'::timestamptz),

  -- YLL-742 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 742,
   'audit: add recover() to the NATS consumer write path',
   E'**Service:** `go_audit_services`\n\n## Problem\nThe entire consumer write path runs with no `recover()` anywhere. A panic crashes the whole audit sink and drops in-flight buffered events. Uncovered goroutines:\n- `internal/transport/consumer/nats/audit.go:188` — `flushWorker`\n- `internal/transport/consumer/nats/audit.go:193` — `bufferDepthMonitor`\n- `internal/transport/consumer/nats/client.go:334` — per-message handler\n\n## Proposed fix\n- Add a recover helper per the YLL-253 spec\n- Wrap per-message handler and add defer recover to flushWorker and bufferDepthMonitor\n- A recovered flush worker must relaunch, not silently die',
   'subtask', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T15:11:01.401+0500'::timestamptz,
   '2026-06-11T15:11:01.545+0500'::timestamptz),

  -- YLL-741 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 741,
   'executor: add pkg/goroutine recover helper and wrap fire-and-forget DB writes',
   E'**Service:** `go-executor-service`\n\n## Problem\nThree per-request fire-and-forget goroutines do DB writes with no `recover()`:\n- `app/services/ExecutorEXServices.go:200` — app-version mismatch update\n- `app/services/ExecutorPlanEXServices.go:296` — plan_history insert (buy)\n- `app/services/ExecutorPlanEXServices.go:377` — plan_history insert (remove)\n\n## Proposed fix\n- Add `pkg/goroutine/goroutine.Go` per the YLL-253 spec\n- Wrap the three sites with `goroutine.Go(c.log, func(){...})`',
   'subtask', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T15:10:49.524+0500'::timestamptz,
   '2026-06-11T15:10:49.665+0500'::timestamptz),

  -- YLL-740 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 740,
   'logic: wrap fire-and-forget goroutines with goroutine.Go (recover + Sentry)',
   E'**Service:** `go_ildam_logic_services`\n\n## Problem\nThis service has `pkg/goroutine/goroutine.Go` and uses it in ~28 places. But many fire-and-forget goroutines still launch raw `go func()` with no recover. Uncovered sites:\n- Audit writes: `intercity_flight_storage.go:269`, `intercity_schedule_storage.go:95`\n- Notification broadcasts: `executor_publisher.go:90,247`\n- NATS consumer fan-out: `executor_became_free.go:124`\n- Dashboard per-request fan-out: `service/grpc/service/dashboard.go` (7 goroutines)\n- Pool workers: `service/worker/pool.go` worker() runs task() with no recover\n\n## Proposed fix\n- Route fire-and-forget goroutines through `goroutine.Go`\n- For WaitGroup-tracked fan-outs add `defer` recover inside the goroutine body\n- Add recover wrapper around `task()` in pool.go',
   'subtask', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T15:10:28.436+0500'::timestamptz,
   '2026-06-11T15:10:28.593+0500'::timestamptz),

  -- YLL-739 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 739,
   'Panel: harden the async worker pool (unbounded overflow, shutdown panics, no panic recovery)',
   E'**Service:** `go-panel-service`. File: `internal/worker/pool.go`\n\n## Problems\n1. **(Medium) Unbounded overflow.** Submit on full buffer spawns `go task()` directly — unbounded growth.\n2. **(Low-Medium) Submit/Stop race** → send on a closed channel → panic.\n3. **(Low-Medium) Stop is not idempotent.** Calling Stop twice panics on double close.\n4. **(Medium) Workers have no panic recovery.** A single panicking task crashes the whole process.\n5. **(Low) Overflow goroutines aren''t tracked** in wg — abandoned on shutdown.\n\n## Proposed fix\n- Replace `close(p.tasks)` with a `stopCh` + `sync.Once`\n- On overflow: drop + log + metric instead of `go task()`\n- Wrap `task()` with recover + structured log\n- Reuse `go_ildam_logic_services` pool design for monorepo consistency',
   'story', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T14:30:48.565+0500'::timestamptz,
   '2026-06-11T14:30:48.736+0500'::timestamptz),

  -- YLL-738 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 738,
   'Integration: initialize the Firebase app once instead of on every FCM push',
   E'**Service:** `go_ildam_integration_services`\n\n## Problem\n`firebaseStorage.SendMessage` and `SendMessageMultiple` call `connectFirebase` on every push, which runs `firebase.NewApp` — reading the service-account file from disk and creating a fresh OAuth2 token source each time.\n\n## Proposed fix\n- Initialize the Firebase `App` and `messaging.Client` once at wiring time\n- Reuse it for every send\n\n## Affected files\n- `storage/helpers/firebase_storage.go:23` (SendMessageMultiple)\n- `storage/helpers/firebase_storage.go:64` (SendMessage)\n- `storage/helpers/firebase_storage.go:102` (connectFirebase)',
   'story', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T12:54:43.797+0500'::timestamptz,
   '2026-06-11T12:56:57.098+0500'::timestamptz),

  -- YLL-737 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 737,
   'Integration: reuse a shared HTTP client instead of building a new transport per outbound call',
   E'**Service:** `go_ildam_integration_services`\n\n## Problem\nSeveral integration helpers build a brand-new `http.Client` with its own `http.Transport` on every outbound request:\n- `paymentStorage.Post` — PayNet\n- `twoGisMapStorage.TwoGisRouting` — 2GIS\n- SMS sender\n- HM Bank (2 sites)\n\nEach transport is its own connection pool — connections never reused, transports never closed. On hot paths this leads to "too many open files" / memory growth.\n\nThe service already has the correct pattern in `pkg/httpclient`.\n\n## Proposed fix\n- Build the HTTP client(s) once in storage constructors\n- Keep one shared client per TLS config\n- Set `IdleConnTimeout` (~90s) and sane `MaxIdleConnsPerHost` (10-50)',
   'story', s_todo, 'medium',
   NULL, u_baxodir,
   '2026-06-11T12:54:30.047+0500'::timestamptz,
   '2026-06-11T12:56:40.898+0500'::timestamptz),

  -- YLL-736 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 736,
   'Driver app sign-up: Android hardware back button on final review screen jumps to first step',
   E'**Environment:** Stage / Driver app / Android\n\n## Steps to reproduce\n1. Go through the sign-up flow to the last stage (driver rechecks all entered info)\n2. Press the Android hardware back button\n\n**Actual:** App navigates to the first page (choose city), discarding the flow.\n**Expected:** Hardware back button should go back one step, not reset to the first step.\n\n**Reported by:** Abdumumin (QA)',
   'bug', s_review, 'medium',
   u_abdulvosid, u_zokirov,
   '2026-06-11T06:15:03.603+0500'::timestamptz,
   '2026-06-11T11:14:18.450+0500'::timestamptz),

  -- YLL-734 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 734,
   'Driver app sign-up: "Yuk" service option should be removed from "Xizmatni tanlang"',
   E'**Environment:** Stage / Driver app / Android\n\n## Steps to reproduce\n1. Open the driver app (Stage env) on Android\n2. Start the sign-up flow\n3. Reach the "Xizmatni tanlang" screen\n\n**Actual:** Three options shown — Taksi, Shaharlararo, **Yuk**.\n**Expected:** "Yuk" (cargo) option should not be available during driver registration. Only Taksi and Shaharlararo.\n\n**Reported by:** Abdumumin (QA)',
   'bug', s_inprog, 'medium',
   u_muhtorxon, u_zokirov,
   '2026-06-11T06:03:53.545+0500'::timestamptz,
   '2026-06-11T13:52:29.012+0500'::timestamptz),

  -- YLL-733 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 733,
   'Image URL is wrong in Dev environment',
   NULL,
   'bug', s_qa, 'high',
   u_ozodbek, u_ozodbek,
   '2026-06-10T15:30:52.523+0500'::timestamptz,
   '2026-06-10T17:46:59.858+0500'::timestamptz),

  -- YLL-732 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 732,
   'Sync Prod, Staging and Dev Minio for shared images',
   NULL,
   'task', s_qa, 'highest',
   u_ozodbek, u_ozodbek,
   '2026-06-10T15:30:09.085+0500'::timestamptz,
   '2026-06-10T17:48:12.495+0500'::timestamptz),

  -- YLL-731 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 731,
   'Harden executor (driver) SMS-OTP to YLL-524 standard',
   E'## Overview\nExecutor (driver) login and registration now use the same hardened SMS-OTP mechanism as the client via a single Redis-backed channel (`executor_otp:*`).\n\n## Rules & limits\n| Rule | Value |\n|---|---|\n| Code length | 5 digits |\n| Code lifetime (TTL) | 180s |\n| Resend cooldown | 60s |\n| Wrong attempts → block | 5 → blocked 300s |\n| Daily limit (per phone) | 10/day |\n| Per-IP rate limit | 5/60s → 1h block |\n\n## Status\nImplemented on branch `YLL-731-executor-otp-security`. Unit-tested, Swagger updated. Pending MR to `dev`.',
   'story', s_done, 'high',
   u_baxodir, u_baxodir,
   '2026-06-10T12:21:14.305+0500'::timestamptz,
   '2026-06-10T17:03:33.492+0500'::timestamptz),

  -- YLL-730 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 730,
   'foto-control PUT fails — shape column varchar(50) too small for full CDN image URL (SQLSTATE 22001)',
   E'**Symptom**\n`PUT /foto-control/{id}` returns 500 when `shape` holds a full CDN image URL.\n```\nSQLSTATE[22001]: String data, right truncated: value too long for type character varying(50)\n```\n\n**Root cause:** `foto_controls.shape` column is `varchar(50)`. The URL `https://test-cdn-api.ildam.uz/images/181777971932.png` is 53 chars.\n\n**Fix options:**\n1. Quick: widen `foto_controls.shape` to `varchar(255)` or `text`\n2. Proper (preferred): store only path/filename, build full URL at read time\n\n**Reported by:** Hojiakbar (Frontend Lead)',
   'bug', s_inprog, 'medium',
   u_muhtorxon, u_zokirov,
   '2026-06-09T19:21:35.059+0500'::timestamptz,
   '2026-06-11T12:20:50.728+0500'::timestamptz),

  -- YLL-729 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 729,
   'Expose audit read API via gateway for admin panel (ListAudit/GetAudit)',
   E'## Problem\n`go_audit_services` implements the read API but it is not reachable from the admin panel. The gateway has no audit route/handler/client.\n\n## Scope (gateway repo)\n- Regenerate ONLY the audit_services genproto from `ildam_protos` main\n- Add audit gRPC client + register in ServiceManager\n- Add HTTP handlers + routes (employee auth): `GET /panel/v2/audits` and `GET /panel/v2/audit/:id`\n- Config/Vault: `AUDIT_SERVICE_BASE_GRPC_URI`\n- Frontend: consume the two endpoints to render audit history\n\n## Out of scope\nAudit write/publish parity (YLL-620), old-log migration (YLL-273)',
   'subtask', s_review, 'medium',
   u_muhtorxon, u_muhtorxon,
   '2026-06-09T17:06:02.385+0500'::timestamptz,
   '2026-06-10T11:44:58.953+0500'::timestamptz),

  -- YLL-728 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 728,
   '/cli/me returns malformed client image URL — missing "/" before default.png',
   E'**Symptom**\n`/cli/me` returns a malformed `image` value:\n```\n"image": "https://test-cdn-api.ildam.uzdefault.png"\n```\nExpected: `https://test-cdn-api.ildam.uz/default.png`\n\n**Likely root cause:** Backend URL builder doing `base_url + filename` without the `/` separator for the default-avatar case.\n\n**Related:** YLL-592 (full `image_url` from file_upload), YLL-476 (double `/images/` slash bug)',
   'bug', s_qa, 'medium',
   u_ozodbek, u_zokirov,
   '2026-06-09T13:24:43.883+0500'::timestamptz,
   '2026-06-10T18:03:07.047+0500'::timestamptz),

  -- YLL-727 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 727,
   'Incorrect value (0) displayed in the top right metric badge on iOS',
   E'**Environment:** iOS Staging Driver\n\nBosh ekranda yuqori o''ng burchakdagi binafsharang dumaloq ikonka ichida doimiy ravishda **"0"** qiymati ko''rsatilmoqda.\n\n## Steps to Reproduce\n1. iOS qurilmada ilovani oching\n2. Xarita ekraniga o''ting\n3. Statusni "Faol" rejimiga o''tkazing\n4. Yuqori o''ng burchakdagi ko''rsatkichga e''tibor bering\n\n**Expected:** Real ko''rsatkich (prioritet bali, buyurtmalar soni)\n**Actual:** Doimiy "0"',
   'task', s_qa, 'highest',
   u_mukhammadali, u_azizbek,
   '2026-06-09T12:02:28.823+0500'::timestamptz,
   '2026-06-11T10:24:35.700+0500'::timestamptz),

  -- YLL-726 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 726,
   'Auth otp confirmation requires phone to be sent for driver',
   NULL,
   'task', s_done, 'medium',
   u_mukhammadali, u_muhammadjon,
   '2026-06-09T10:24:07.200+0500'::timestamptz,
   '2026-06-09T10:59:39.980+0500'::timestamptz),

  -- YLL-725 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 725,
   'make phone number format max length required in login Android Driver',
   NULL,
   'task', s_review, 'medium',
   u_abdulvosid, u_abdulvosid,
   '2026-06-09T09:48:57.279+0500'::timestamptz,
   '2026-06-11T06:25:39.135+0500'::timestamptz),

  -- YLL-724 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 724,
   'Logging enhancement',
   NULL,
   'task', s_review, 'highest',
   u_ozodbek, u_ozodbek,
   '2026-06-09T03:00:01.005+0500'::timestamptz,
   '2026-06-10T12:55:45.587+0500'::timestamptz),

  -- YLL-723 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 723,
   '[Android Driver] Location Accuracy permission dialog re-prompts infinitely after "No thanks"',
   E'**Platform:** Android driver app.\n\n**Summary:** The Google Play Services "Location Accuracy" dialog keeps reappearing in a loop. Tapping "No thanks" dismisses it, then it immediately pops up again.\n\n**Where:** Driver "Faol/Active" screen (waiting-for-orders state).\n\n**Expected:** Respect the decline — don''t immediately re-prompt. Re-check only on a meaningful trigger.\n\n**Likely cause:** `SettingsClient.checkLocationSettings` is being re-invoked on a tight loop and re-launching the RESOLUTION_REQUIRED resolution every cycle without tracking the user''s decline.',
   'bug', s_done, 'medium',
   u_abdulvosid, u_zokirov,
   '2026-06-09T00:12:47.316+0500'::timestamptz,
   '2026-06-11T10:35:51.518+0500'::timestamptz),

  -- YLL-722 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 722,
   '[Android Driver] Internet connect/disconnect voice notification plays 5+ times per event',
   E'**Platform:** Android driver app.\n\n**Summary:** The voice notification tied to internet connectivity events fires repeatedly instead of once.\n\n**Steps / observed:**\n- On internet disconnect, the voice notification plays 5+ times.\n- On internet reconnect, same — plays 5+ times.\n\n**Expected:** Exactly one voice notification per connectivity transition.\n\n**Likely cause:** Connectivity listener fires multiple callbacks per transition with no debounce/de-duplication against last known state.',
   'bug', s_review, 'high',
   u_abdulvosid, u_zokirov,
   '2026-06-09T00:01:52.398+0500'::timestamptz,
   '2026-06-11T06:23:18.563+0500'::timestamptz),

  -- YLL-721 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 721,
   'Pricing calculated incorrectly for intercity routes (order-of-magnitude error)',
   E'**Severity:** P0 — pricing core flow broken.\n\n**Summary:** Fare calculation returns wildly incorrect, internally inconsistent prices for intercity routes.\n\n**Observed (staging):**\n- Route: Fergana ↔ Tashkent (~300 km, ~257 min ETA)\n- Client app quoted: Standart ~700 so''m, Komfort 8,000 so''m, Yetkazish 734,400 so''m\n- Driver app showed the same order at ~8,500 so''m (Komfort)\n\n**Expected:** Intercity fare must reflect distance/time per the tariff; values must be within sane bounds.\n\n**Impact:** Customers can be quoted near-zero fares (revenue loss) or absurd fares (lost orders). Blocks release.',
   'bug', s_done, 'highest',
   u_farrux, u_zokirov,
   '2026-06-08T23:42:36.586+0500'::timestamptz,
   '2026-06-11T16:14:25.596+0500'::timestamptz),

  -- YLL-720 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 720,
   'Add Swagger API docs to php-taxi-service (SwaggerLume)',
   E'Same story as the gateway — `php-taxi-service` ships no API docs.\n\n**Plan:**\n- Add `darkaonline/swagger-lume` (Lumen 8 compatible) and publish the config\n- Register the provider in `bootstrap/app.php`\n- Add the base `@OA\Info` / server block and annotate controllers\n- Generate via `php artisan swagger-lume:generate`\n\n**Done when:** `/api/documentation` serves the taxi-service docs and the main client + executor endpoints are covered.',
   'task', s_done, 'highest',
   u_baxodir, u_baxodir,
   '2026-06-08T22:01:23.736+0500'::timestamptz,
   '2026-06-09T02:41:44.506+0500'::timestamptz),

  -- YLL-719 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 719,
   'Add Swagger API docs to php-gateway-service (SwaggerLume)',
   E'Right now `php-gateway-service` has no API documentation at all. The `api.php` route file alone is ~64KB.\n\n**Plan:**\n- Pull in `darkaonline/swagger-lume` and publish its config\n- Register the service provider in `bootstrap/app.php`\n- Add base `@OA\Info` / server annotations, then annotate controllers incrementally\n- Make sure the Swagger UI route works through `public/index.php` and inside Docker\n\n**Done when:**\n- `/api/documentation` loads the Swagger UI\n- Core client + executor endpoints are documented\n- Generation is repeatable (make/composer target in README)',
   'task', s_done, 'highest',
   u_baxodir, u_baxodir,
   '2026-06-08T22:01:11.779+0500'::timestamptz,
   '2026-06-09T02:41:45.889+0500'::timestamptz),

  -- YLL-718 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 718,
   'mobile OTP/SMS tasdiqlash klaviaturasida alfanumerik kiritish. ANDROID DRIVER',
   E'Card qo''shish va payment''ga aloqador OTP/SMS tasdiqlash klaviaturasida harf ham kiritsa bo''ladigan qilish (faqat raqam emas, alfanumerik).\n\nPlatforma: Android Driver',
   'subtask', s_review, 'medium',
   u_abdulvosid, u_khikmat,
   '2026-06-08T16:54:15.308+0500'::timestamptz,
   '2026-06-11T10:34:27.092+0500'::timestamptz),

  -- YLL-717 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 717,
   'mobile OTP/SMS tasdiqlash klaviaturasida alfanumerik kiritish. ANDROID CLIENT',
   E'Card qo''shish va payment''ga aloqador OTP/SMS tasdiqlash klaviaturasida harf ham kiritsa bo''ladigan qilish (faqat raqam emas, alfanumerik).\n\nPlatforma: Android Client',
   'subtask', s_todo, 'medium',
   u_islom, u_khikmat,
   '2026-06-08T16:53:59.711+0500'::timestamptz,
   '2026-06-10T12:25:52.665+0500'::timestamptz),

  -- YLL-716 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 716,
   'mobile OTP/SMS tasdiqlash klaviaturasida alfanumerik kiritish. IOS CLIENT',
   E'Card qo''shish va payment''ga aloqador OTP/SMS tasdiqlash klaviaturasida harf ham kiritsa bo''ladigan qilish (faqat raqam emas, alfanumerik).\n\nPlatforma: iOS Client',
   'subtask', s_review, 'medium',
   u_muhammadjon, u_khikmat,
   '2026-06-08T16:53:48.184+0500'::timestamptz,
   '2026-06-08T23:57:17.369+0500'::timestamptz),

  -- YLL-715 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 715,
   'mobile OTP/SMS tasdiqlash klaviaturasida alfanumerik kiritish. IOS DRIVER',
   E'Card qo''shish va payment''ga aloqador OTP/SMS tasdiqlash klaviaturasida harf ham kiritsa bo''ladigan qilish (faqat raqam emas, alfanumerik).\n\nPlatforma: iOS Driver',
   'task', s_done, 'medium',
   u_mukhammadali, u_khikmat,
   '2026-06-08T16:53:12.713+0500'::timestamptz,
   '2026-06-09T11:01:26.798+0500'::timestamptz),

  -- YLL-714 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 714,
   'mobile OTP/SMS tasdiqlash klaviaturasida alfanumerik kiritish',
   E'Card qo''shish va payment''ga aloqador OTP/SMS tasdiqlash klaviaturasida harf ham kiritsa bo''ladigan qilish (faqat raqam emas, alfanumerik).\n\nAsosiy task — subtasklar: iOS Driver (YLL-715), iOS Client (YLL-716), Android Client (YLL-717), Android Driver (YLL-718).',
   'task', s_qa, 'medium',
   u_abdulvosid, u_khikmat,
   '2026-06-08T16:53:06.581+0500'::timestamptz,
   '2026-06-10T11:40:22.031+0500'::timestamptz),

  -- YLL-713 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 713,
   'Client app: hide client''s own location marker on map, show only driver location',
   E'**Platform:** Yalla — Client app\n\n**Current behavior:** On the map the client sees both his own location marker and the driver''s location. When client and driver are in the same car, the driver''s marker can appear behind the client marker.\n\n**Requested change:** Hide the icon indicating the client''s current location — show only the driver''s location.\n\n**Reference:** Same behavior in Yandex, Wildberries (WB) and Uklon Taxi.\n\n**Reported by:** Abdumumin (Manual QA)',
   'task', s_todo, 'medium',
   u_islom, u_zokirov,
   '2026-06-08T03:43:09.542+0500'::timestamptz,
   '2026-06-10T04:21:35.640+0500'::timestamptz),

  -- YLL-712 ─────────────────────────────────────────────────────────────────
  (gen_random_uuid(), proj_id, 712,
   'error text dan api ni ko''rsatmaslik (asosan internet errorlarida "api ga ulanib bo''lmadi" ko''rinishida)',
   E'Internet ulanish xatolarida foydalanuvchiga API URL ko''rinmasligi kerak. Masalan: "Failed to connect to api.example.com" o''rniga "API ga ulanib bo''lmadi" yoki "Internet aloqasini tekshiring" kabi tushunarli xabar ko''rsatilsin.',
   'task', s_review, 'medium',
   u_abdulvosid, u_abdulvosid,
   '2026-06-06T14:29:20.639+0500'::timestamptz,
   '2026-06-11T07:15:52.997+0500'::timestamptz)

  ON CONFLICT (project_id, issue_number) DO NOTHING;

END $$;

-- =============================================================================
-- Natija tekshiruvi
-- =============================================================================
DO $$
DECLARE
  user_count    INT;
  issue_count   INT;
BEGIN
  SELECT COUNT(*) INTO user_count  FROM users  WHERE email LIKE '%@yalla.uz';
  SELECT COUNT(*) INTO issue_count FROM issues WHERE project_id = 'd0000001-0000-0000-0000-000000000000';
  RAISE NOTICE 'Seed yakunlandi: % ta YLL foydalanuvchi, % ta YLL issue', user_count, issue_count;
END $$;

COMMIT;
