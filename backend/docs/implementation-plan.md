# JiraFlow — To'liq Implementatsiya Plani
# Jira + Confluence 100% o'rnini bosish

**Holat:** 2026-05-23  
**Maqsad:** Jira Software + Confluence ni to'liq almashtirish  
**Hozirgi holat:** Jira ~55% | Confluence ~40%

---

## Umumiy ko'rinish

```
FAZA 1 — Kritik asoslar         (2–3 hafta)
FAZA 2 — Jira to'ldirish        (3–4 hafta)
FAZA 3 — Confluence to'ldirish  (3–4 hafta)
FAZA 4 — Jira ilg'or            (3–4 hafta)
FAZA 5 — Confluence ilg'or      (2–3 hafta)
FAZA 6 — Integratsiya           (2–3 hafta)
FAZA 7 — Enterprise             (3–4 hafta)
```

**Jami taxminiy vaqt:** 18–25 hafta (1 developer, to'liq vaqt)

---

## FAZA 1 — Kritik Asoslar (Hozir buzilgan yoki bo'sh)

> Bular bo'lmasa tizim yarim ishlaydi. Birinchi qilish shart.

---

### F1.1 — WebSocket Route

**Muammo:** `websocket.Hub` yozilgan, lekin router'da hech qanday `/ws` yo'q.  
**Natija:** Real-time notification umuman ishlamaydi.

**Vazifalar:**
- [ ] `api/router.go` — `/ws` endpoint qo'shish, Auth middleware bilan
- [ ] `api/handlers/v1/ws.go` — handler yozish (userID JWT'dan olinadi)
- [ ] `internal/infrastructura/websocket/notify.go` — `SendToUser(userID, payload)` metodi
- [ ] Notification dispatcher'ni WebSocket bilan ulash

**Migration:** yo'q  
**API:** `GET /api/v1/ws` (WebSocket upgrade)

---

### F1.2 — Issue Ordering (Backlog va Sprint tartib)

**Muammo:** Issues'da `position` maydoni yo'q — backlog drag-and-drop imkonsiz.

**Vazifalar:**
- [ ] Migration: `issues` jadvaliga `position FLOAT8 DEFAULT 0` qo'shish (LexoRank uchun float yaxshi)
- [ ] `entity/issue.go` — `Position float64` maydon qo'shish
- [ ] Repository: `UpdatePosition(ctx, issueID, position)` metod
- [ ] Usecase: `ReorderIssues(ctx, projectID, items []ReorderItem)`
- [ ] Handler: `PUT /api/v1/issues/reorder`
- [ ] Backlog list'da `ORDER BY position ASC` qo'shish

**Migration:** `000014_issue_position.up.sql`

---

### F1.3 — Page Watchers API

**Muammo:** `page_watchers` DB jadvali bor, lekin hech qanday API endpoint yo'q.

**Vazifalar:**
- [ ] `entity/page_watcher.go` — entity struct
- [ ] Repository interface + postgres implementatsiya
- [ ] Usecase: `WatchPage`, `UnwatchPage`, `ListPageWatchers`, `IsWatching`
- [ ] Handler: `POST /api/v1/pages/:id/watchers`
- [ ] Handler: `DELETE /api/v1/pages/:id/watchers`
- [ ] Handler: `GET /api/v1/pages/:id/watchers`
- [ ] Notification dispatcher'ga `PageUpdated` va `PageCommented` event qo'shish

**API:**
```
POST   /api/v1/pages/:id/watchers
DELETE /api/v1/pages/:id/watchers
GET    /api/v1/pages/:id/watchers
```

---

### F1.4 — Page Delete Ruxsati To'g'rilash

**Muammo:** `page.go:133` — faqat muallif o'chira oladi. Space admin ham o'chira olishi kerak.

**Vazifalar:**
- [ ] `usecase/page/page.go` — Space membership tekshiruvi qo'shish
- [ ] Space admin role'ni tekshirish helper

---

### F1.5 — Email Notification Yuborish

**Muammo:** Email template'lar bor, lekin notifikatsiya dispatcher email chaqirmaydi.

**Vazifalar:**
- [ ] `usecase/notification/notification.go` — `Notify()` da email yuborish logikasi
- [ ] User notification preference'ni tekshirish (email opt-in/out)
- [ ] `IssueAssigned` → assigned template
- [ ] `Mentioned` → mentioned template
- [ ] Async yuborish (goroutine, xato bo'lsa log, panic bo'lmasin)

---

## FAZA 2 — Jira To'ldirish

> Jira'ning asosiy funksiyalari — bular bo'lmasa "Jira alternative" bo'lmaydi.

---

### F2.1 — Time Tracking (Vaqt Hisobi)

**Jira analog:** Log Work, Time Spent, Original Estimate

**Vazifalar:**
- [ ] Migration: `issue_worklogs` jadval

```sql
CREATE TABLE issue_worklogs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id         UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    author_id        UUID NOT NULL REFERENCES users(id),
    time_spent_sec   INT NOT NULL CHECK (time_spent_sec > 0),
    started_at       TIMESTAMPTZ NOT NULL,
    comment          TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- [ ] Migration: `issues` jadvaliga qo'shimcha maydonlar:

```sql
ALTER TABLE issues ADD COLUMN original_estimate_sec INT;
ALTER TABLE issues ADD COLUMN remaining_estimate_sec INT;
ALTER TABLE issues ADD COLUMN time_spent_sec INT NOT NULL DEFAULT 0;
```

- [ ] `entity/worklog.go` — struct, CreateReq, UpdateReq, Filter
- [ ] Repository interface + postgres
- [ ] Usecase: `LogWork`, `UpdateWorklog`, `DeleteWorklog`, `ListWorklogs`
- [ ] Issue update'da `time_spent_sec` avtomatik yig'ish
- [ ] Handlers:
  - `POST /api/v1/issues/:id/worklogs`
  - `GET /api/v1/issues/:id/worklogs`
  - `PUT /api/v1/worklogs/:id`
  - `DELETE /api/v1/worklogs/:id`

**Migration:** `000015_time_tracking.up.sql`

---

### F2.2 — Components (Loyiha Komponentlari)

**Jira analog:** Components (frontend, backend, mobile, API...)

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE components (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT,
    lead_id      UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, name)
);

CREATE TABLE issue_components (
    issue_id      UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    component_id  UUID NOT NULL REFERENCES components(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, component_id)
);
```

- [ ] Entity, repository, usecase, handlers
- [ ] Issue CRUD'ga `component_ids` qo'shish
- [ ] Issue filter'ga `component_id` qo'shish

**API:**
```
POST   /api/v1/projects/:id/components
GET    /api/v1/projects/:id/components
PUT    /api/v1/components/:id
DELETE /api/v1/components/:id
```

**Migration:** `000016_components.up.sql`

---

### F2.3 — Versions / Releases

**Jira analog:** Fix Version, Affects Version

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE versions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    description  TEXT,
    status       VARCHAR(16) NOT NULL DEFAULT 'unreleased'
                 CHECK (status IN ('unreleased', 'released', 'archived')),
    start_date   DATE,
    release_date DATE,
    released_at  TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, name)
);

CREATE TABLE issue_fix_versions (
    issue_id    UUID NOT NULL REFERENCES issues(id)   ON DELETE CASCADE,
    version_id  UUID NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, version_id)
);

CREATE TABLE issue_affects_versions (
    issue_id    UUID NOT NULL REFERENCES issues(id)   ON DELETE CASCADE,
    version_id  UUID NOT NULL REFERENCES versions(id) ON DELETE CASCADE,
    PRIMARY KEY (issue_id, version_id)
);
```

- [ ] Entity, repository, usecase, handlers
- [ ] `POST /versions/:id/release` — versiyani release qilish
- [ ] Issue filter'ga `fix_version_id` qo'shish

**API:**
```
POST   /api/v1/projects/:id/versions
GET    /api/v1/projects/:id/versions
PUT    /api/v1/versions/:id
DELETE /api/v1/versions/:id
POST   /api/v1/versions/:id/release
GET    /api/v1/versions/:id/issues
```

**Migration:** `000017_versions.up.sql`

---

### F2.4 — Epic Progress Tracking

**Muammo:** Epic ichidagi child issuelar asosida `%` completeness hisoblanmaydi.

**Vazifalar:**
- [ ] Repository: `GetEpicProgress(ctx, epicID) (total, done int, err error)`
- [ ] Usecase: `GetIssue` da epic uchun `progress` qo'shish
- [ ] `entity/issue.go` — `Progress *EpicProgress` maydon

```go
type EpicProgress struct {
    Total     int     `json:"total"`
    Done      int     `json:"done"`
    Percent   float64 `json:"percent"`
}
```

- [ ] "Done" statusni workflow'dan aniqlash (category='done')
- [ ] Migration: `workflow_statuses` ga `category VARCHAR(16) CHECK (category IN ('todo','in_progress','done'))` qo'shish

---

### F2.5 — Bulk Issue Operations

**Vazifalar:**
- [ ] `entity/issue.go` — `BulkUpdateReq` struct

```go
type BulkUpdateReq struct {
    IssueIDs   []string `json:"issue_ids"  validate:"required,min=1,max=100"`
    AssigneeID *string  `json:"assignee_id"`
    StatusID   *string  `json:"status_id"`
    SprintID   *string  `json:"sprint_id"`
    Priority   *string  `json:"priority"`
    LabelIDs   []string `json:"label_ids"`
}
```

- [ ] Repository: `BulkUpdate(ctx, req)`
- [ ] Usecase + handler
- [ ] `PUT /api/v1/issues/bulk`
- [ ] `DELETE /api/v1/issues/bulk` (bulk delete)

---

### F2.6 — Advanced Issue Filtering (JQL-like)

**Jira analog:** `project = PROJ AND status IN (Done, Review) AND assignee = currentUser()`

**Vazifalar:**
- [ ] `entity/issue.go` — `IssueFilter` kengaytirish:
  - `StatusIDs []string` (bir nechta status)
  - `Types []string`
  - `Priorities []string`
  - `AssigneeIDs []string`
  - `HasDueDate *bool`
  - `DueDateFrom`, `DueDateTo *time.Time`
  - `CreatedFrom`, `CreatedTo *time.Time`
  - `TextSearch string` (full-text)
  - `ComponentIDs []string`
  - `VersionID string`
  - `Unassigned *bool`

- [ ] Repository: SQL query builder (WHERE clause dinamik quriladi)
- [ ] `GET /api/v1/issues` — query param sifatida yuqoridagi filterlar

---

### F2.7 — Saved Filters / Views

**Jira analog:** "My open issues", "All issues in PROJ", saved filters

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE saved_filters (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id  UUID REFERENCES projects(id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    filter_json JSONB NOT NULL,
    is_shared   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- [ ] Entity, repository, usecase, handlers

**API:**
```
POST   /api/v1/saved-filters
GET    /api/v1/saved-filters
PUT    /api/v1/saved-filters/:id
DELETE /api/v1/saved-filters/:id
```

---

### F2.8 — Sprint Reports & Charts (API)

**Jira analog:** Burndown chart, Velocity chart, Sprint report

**Vazifalar:**
- [ ] `entity/report.go` — Report response struct'lari
- [ ] Repository: `GetBurndownData(ctx, sprintID) ([]BurndownPoint, error)`
- [ ] Repository: `GetVelocityData(ctx, projectID, lastN int) ([]VelocityPoint, error)`
- [ ] Repository: `GetSprintReport(ctx, sprintID) (*SprintReport, error)` — completed/incomplete/removed issues

**BurndownPoint:**
```go
type BurndownPoint struct {
    Date          time.Time `json:"date"`
    RemainingWork int       `json:"remaining_work"` // story points yoki issue count
    IdealWork     int       `json:"ideal_work"`
}
```

**API:**
```
GET /api/v1/sprints/:id/reports/burndown
GET /api/v1/sprints/:id/reports/summary
GET /api/v1/projects/:id/reports/velocity
GET /api/v1/projects/:id/reports/cumulative-flow
GET /api/v1/projects/:id/reports/issue-statistics
```

---

### F2.9 — Roadmap / Timeline

**Jira analog:** Roadmap view (Epic'lar timeline'da)

**Vazifalar:**
- [ ] `issues` ga `start_date DATE` qo'shish (migration)
- [ ] `entity/issue.go` — `StartDate *time.Time`
- [ ] Repository: `GetRoadmap(ctx, projectID, from, to time.Time) ([]*Issue, error)` — epic'lar + child count
- [ ] `GET /api/v1/projects/:id/roadmap?from=2026-01-01&to=2026-12-31`

---

### F2.10 — Dashboard / Gadgets

**Jira analog:** Project dashboard, configurable gadgets

**Vazifalar:**
- [ ] `GET /api/v1/projects/:id/dashboard` — yig'ilgan statistika:
  - Ochiq/yopiq issue soni (type bo'yicha)
  - Aktiv sprint holati
  - Assignee bo'yicha distribution
  - Priority bo'yicha distribution
  - Son 7 kun faoliyat

---

### F2.11 — Board Yaxshilanishi

**Vazifalar:**
- [ ] Swimlanes: `GET /api/v1/boards/:id/issues?group_by=assignee|epic|type`
- [ ] Board quick filters: `board_filters` jadval yoki board `filter` JSONB kengaytirish
- [ ] WIP limit hisoblash: column uchun `current_count` qaytarish
- [ ] Issue drag board'da: `PUT /api/v1/issues/:id/move` (column + position)
- [ ] Backlog endpoint: `GET /api/v1/projects/:id/backlog` (sprint yo'q issuelar)
- [ ] Sprint'ga issue qo'shish: `POST /api/v1/sprints/:id/issues`
- [ ] Sprint'dan issue chiqarish: `DELETE /api/v1/sprints/:id/issues/:issue_id`

---

## FAZA 3 — Confluence To'ldirish

> Confluence'ning asosiy funksiyalari.

---

### F3.1 — Page Templates

**Confluence analog:** "Meeting notes", "Project plan", "How-to article" shablonlar

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE page_templates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id    UUID REFERENCES spaces(id) ON DELETE CASCADE,  -- NULL = global
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    content     JSONB NOT NULL DEFAULT '{}'::jsonb,
    icon        VARCHAR(10),  -- emoji
    category    VARCHAR(64),  -- 'meeting', 'documentation', 'planning', ...
    is_global   BOOLEAN NOT NULL DEFAULT FALSE,
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- [ ] Default global template'larni seed qilish (meeting notes, retrospective, decision log)
- [ ] Entity, repository, usecase, handlers
- [ ] Page yaratishda `template_id` qabul qilish

**API:**
```
POST   /api/v1/page-templates
GET    /api/v1/page-templates?space_id=...
GET    /api/v1/page-templates/:id
PUT    /api/v1/page-templates/:id
DELETE /api/v1/page-templates/:id
```

---

### F3.2 — Page Tags / Labels

**Confluence analog:** Page labels (tagging)

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE page_tags (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id   UUID NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
    name       VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (space_id, name)
);

CREATE TABLE page_tag_links (
    page_id UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES page_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (page_id, tag_id)
);
```

- [ ] Entity, repository, usecase, handlers
- [ ] `GET /api/v1/pages?tag=...` filtering

---

### F3.3 — Page Restrictions (Per-page permissions)

**Confluence analog:** Page restrictions — view/edit per user/group

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE page_restrictions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    page_id     UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id     UUID REFERENCES users(id) ON DELETE CASCADE,
    permission  VARCHAR(16) NOT NULL CHECK (permission IN ('view', 'edit')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (page_id, user_id, permission)
);
```

- [ ] Middleware/usecase: page olishda restriction tekshirish
- [ ] Handlers:
  - `POST /api/v1/pages/:id/restrictions`
  - `GET /api/v1/pages/:id/restrictions`
  - `DELETE /api/v1/pages/:id/restrictions/:restriction_id`

---

### F3.4 — Favorites / Starred Pages & Spaces

**Confluence analog:** Starred pages, My spaces

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE user_favorites (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type  VARCHAR(16) NOT NULL CHECK (entity_type IN ('page', 'space')),
    entity_id    UUID NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, entity_type, entity_id)
);
```

- [ ] Entity, repository, usecase, handlers

**API:**
```
POST   /api/v1/favorites        { entity_type, entity_id }
DELETE /api/v1/favorites/:id
GET    /api/v1/favorites?type=page|space
```

---

### F3.5 — Recently Visited

**Confluence analog:** "Recently viewed" sidebar

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE user_recents (
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type  VARCHAR(16) NOT NULL CHECK (entity_type IN ('page', 'space', 'issue', 'project')),
    entity_id    UUID NOT NULL,
    visited_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, entity_type, entity_id)
);
CREATE INDEX idx_user_recents_user_id ON user_recents (user_id, visited_at DESC);
```

- [ ] Page/Space GET handler'larida `upsert` recent record
- [ ] `GET /api/v1/recents?limit=10`

---

### F3.6 — Inline Comments (Page specific text)

**Confluence analog:** Specific text ustida comment qoldirish

**Vazifalar:**
- [ ] Migration: `comments` jadvaliga qo'shimcha maydonlar:

```sql
ALTER TABLE comments ADD COLUMN inline_start INT;      -- character offset
ALTER TABLE comments ADD COLUMN inline_end   INT;
ALTER TABLE comments ADD COLUMN inline_text  TEXT;     -- highlighted text
ALTER TABLE comments ADD COLUMN is_resolved  BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE comments ADD COLUMN resolved_by  UUID REFERENCES users(id);
ALTER TABLE comments ADD COLUMN resolved_at  TIMESTAMPTZ;
```

- [ ] `entity/comment.go` — yangi maydonlar
- [ ] Handler: `POST /api/v1/pages/:id/inline-comments`
- [ ] Handler: `POST /api/v1/comments/:id/resolve`
- [ ] Handler: `GET /api/v1/pages/:id/inline-comments`

---

### F3.7 — Page Analytics (View count)

**Confluence analog:** Page views, unique visitors

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE page_views (
    page_id    UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id    UUID REFERENCES users(id) ON DELETE SET NULL,
    viewed_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_page_views_page_id ON page_views (page_id, viewed_at DESC);
```

- [ ] Page GET handler'ida async view log yozish
- [ ] `GET /api/v1/pages/:id/analytics` — total views, unique viewers, views per day

---

### F3.8 — Page Export

**Confluence analog:** Export to PDF / Word

**Vazifalar:**
- [ ] `internal/infrastructura/export/` — PDF va HTML export
- [ ] TipTap JSON → HTML converter (tiptap/renderer.go mavjud — kengaytirish)
- [ ] HTML → PDF: `chromedp` yoki `wkhtmltopdf` library
- [ ] Handler: `GET /api/v1/pages/:id/export?format=pdf|html`

---

### F3.9 — Space Statistics

**Confluence analog:** Space overview, page count, contributors

**Vazifalar:**
- [ ] `GET /api/v1/spaces/:id/statistics`:
  - Jami page soni
  - Draft / published nisbati
  - So'nggi 30 kun faoliyat
  - Top contributors (kim ko'p yozgan)
  - Ko'p ko'rilgan sahifalar

---

### F3.10 — Space Archive / Restore

**Muammo:** `is_archived` maydoni bor, lekin archive/restore API endpoint yo'q.

**Vazifalar:**
- [ ] `POST /api/v1/spaces/:id/archive`
- [ ] `POST /api/v1/spaces/:id/restore`
- [ ] Arxivlangan space'da yangi page yaratish bloklash

---

### F3.11 — Page Mention Notification

**Muammo:** Issue'da mention bor, lekin page'da `@user` mention notification yuborilmaydi.

**Vazifalar:**
- [ ] `usecase/comment/mention.go` — page uchun ham ishlash (hozir faqat issue)
- [ ] `dispatcher.go` — `PageMentioned` event qo'shish
- [ ] Page kommentida mention parse qilish

---

## FAZA 4 — Jira Ilg'or

---

### F4.1 — Webhook Support

**Jira analog:** Webhooks — tashqi tizimga event yuborish

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE webhooks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id  UUID REFERENCES projects(id) ON DELETE CASCADE,  -- NULL = global
    name        VARCHAR(100) NOT NULL,
    url         TEXT NOT NULL,
    secret      TEXT,  -- HMAC signing secret
    events      TEXT[] NOT NULL,  -- ['issue.created', 'issue.updated', ...]
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

- [ ] `internal/infrastructura/webhook/` — HTTP POST yuborish, HMAC signature
- [ ] Retry logikasi (3 marta, exponential backoff)
- [ ] `webhook_deliveries` jadval — delivery log
- [ ] Events: `issue.*`, `sprint.*`, `page.*`, `project.*`

**API:**
```
POST   /api/v1/webhooks
GET    /api/v1/webhooks
PUT    /api/v1/webhooks/:id
DELETE /api/v1/webhooks/:id
POST   /api/v1/webhooks/:id/test
GET    /api/v1/webhooks/:id/deliveries
```

---

### F4.2 — Issue Dependency (Gantt blocking)

**Jira analog:** Dependency visualization, blocking chains

**Vazifalar:**
- [ ] `GET /api/v1/issues/:id/dependency-chain` — blocks/blocked_by zanjirini qaytarish (recursive)
- [ ] Circular dependency tekshiruvi issue link yaratishda

---

### F4.3 — Daily Digest Email (Cron)

**Muammo:** `daily_digest.html` template bor, lekin cron job yo'q.

**Vazifalar:**
- [ ] `internal/usecase/notification/digest.go` — digest logikasi
  - Foydalanuvchiga tayinlangan ochiq issuelar
  - Kechikkan (overdue) issuelar
  - Bugun due date bo'lgan issuelar
  - Yangi mentionlar (o'qilmagan)
- [ ] Cron scheduler: `internal/pkg/cron/` — har kuni 08:00 da ishga tushadi
- [ ] User preference: digest opt-in/out

---

### F4.4 — Issue Subscribers (CC)

**Jira analog:** Share issue with users, CC on updates

**Vazifalar:**
- [ ] `issue_watchers` ni kengaytirish yoki issue share API
- [ ] `POST /api/v1/issues/:id/share` — email yuborish + watcher qo'shish

---

### F4.5 — Multiple Assignees

**Jira analog:** Service Management'da bir issueга bir nechta assignee

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE issue_assignees (
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    PRIMARY KEY (issue_id, user_id)
);
```

- [ ] `issues.assignee_id` — optional qoldirish (primary assignee)
- [ ] `entity/issue.go` — `Assignees []UserShort`

---

### F4.6 — Issue Rank (Priority ordering)

**Jira analog:** LexoRank — drag-and-drop backlog ordering

**Vazifalar:**
- [ ] `internal/pkg/lexorank/` — LexoRank algoritmini implementatsiya
- [ ] `issues` jadvaliga `rank VARCHAR(255)` ustun
- [ ] `PUT /api/v1/issues/:id/rank` — `{ after_id, before_id }` asosida rank hisoblash

---

### F4.7 — Issue Clone

**Jira analog:** "Clone issue" — barcha maydonlar bilan nusxa

**Vazifalar:**
- [ ] Usecase: `CloneIssue(ctx, issueID, actorID) (*Issue, error)`
  - Yangi issue_number
  - Labels, custom fields, component, assignee ko'chirish
  - Links, worklogs, attachments — ixtiyoriy
- [ ] `POST /api/v1/issues/:id/clone`

---

## FAZA 5 — Confluence Ilg'or

---

### F5.1 — Real-time Collaborative Editing

**Confluence analog:** Bir vaqtda bir nechta odam sahifani tahrirlaydi

**Vazifalar:**
- [ ] `internal/infrastructura/collab/` — OT (Operational Transformation) yoki CRDT
- [ ] WebSocket orqali document sync: `room:{page_id}` channel
- [ ] "X is editing" presence indicator
- [ ] Conflict resolution strategy (last-write-wins minimal version uchun)

> **Izoh:** Bu eng murakkab feature. Birinchi bosqichda "lock-based" editing qilsa ham bo'ladi: bir vaqtda faqat bitta odam tahrirlaydi.

**Minimal versiya:**
- [ ] `page_locks` jadval: `page_id`, `user_id`, `locked_at`, `expires_at`
- [ ] `POST /api/v1/pages/:id/lock`
- [ ] `DELETE /api/v1/pages/:id/lock`
- [ ] Lock avtomatik 30 daqiqadan so'ng tugaydi

---

### F5.2 — Page Reactions

**Confluence analog:** Emoji reactions (👍 ❤️ 💡 ...)

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE page_reactions (
    page_id    UUID NOT NULL REFERENCES pages(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    emoji      VARCHAR(10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (page_id, user_id, emoji)
);
```

- [ ] `POST /api/v1/pages/:id/reactions` — `{ emoji: "👍" }`
- [ ] `DELETE /api/v1/pages/:id/reactions/:emoji`
- [ ] `GET /api/v1/pages/:id/reactions` — grouped by emoji + count

---

### F5.3 — Page Macros / Embeds

**Confluence analog:** Jira issues ro'yxatini page'da ko'rsatish

**Vazifalar:**
- [ ] `GET /api/v1/pages/:id/macros/issue-list?project_id=...&status_id=...` — page embed uchun issue list
- [ ] `GET /api/v1/pages/:id/macros/sprint-status?sprint_id=...` — sprint progress
- [ ] Frontend TipTap'da custom node sifatida render qiladi

---

### F5.4 — Content Import/Export

**Vazifalar:**
- [ ] Confluence XML export import qilish
- [ ] Markdown → TipTap JSON converter
- [ ] Notion export → import
- [ ] Space'ni ZIP sifatida export (barcha page'lar HTML/JSON)

---

### F5.5 — Global Search Yaxshilanishi

**Hozirgi holat:** Search bor, lekin chekli.

**Vazifalar:**
- [ ] Relevance ranking (ts_rank bilan)
- [ ] Highlight/excerpt (ts_headline bilan)
- [ ] Search suggestions / autocomplete (`pg_trgm` bilan)
- [ ] Search history (foydalanuvchi so'nggi qidiruvlari)
- [ ] Filter: `updated_after`, `author_id`, `space_id`, `project_id`
- [ ] `GET /api/v1/search/suggestions?q=...`

---

## FAZA 6 — Jira + Confluence Integratsiya

> Ikki tizimni birlashtiradigan asosiy ko'priklar

---

### F6.1 — Issue ↔ Page Link

**Jira+Confluence analog:** Issue'da "Confluence pages" bo'limi; Page'da "Jira issues" bo'limi

**Vazifalar:**
- [ ] Migration:

```sql
CREATE TABLE issue_page_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    issue_id   UUID NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    page_id    UUID NOT NULL REFERENCES pages(id)  ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (issue_id, page_id)
);
```

- [ ] Entity, repository, usecase, handlers

**API:**
```
POST   /api/v1/issues/:id/page-links     { page_id }
GET    /api/v1/issues/:id/page-links
DELETE /api/v1/issues/:id/page-links/:page_id

GET    /api/v1/pages/:id/issue-links
```

---

### F6.2 — Project ↔ Space Avtomatik Bog'lash

**Jira+Confluence analog:** Har bir Jira project'ga Confluence space avtomatik yaratiladi

**Vazifalar:**
- [ ] Usecase: `CreateProject` da `type='project'` space avtomatik yaratish (ixtiyoriy, default OFF)
- [ ] `POST /api/v1/projects/:id/link-space` — mavjud space'ni bog'lash
- [ ] `DELETE /api/v1/projects/:id/link-space`
- [ ] `GET /api/v1/projects/:id/space` — bog'liq space'ni olish

---

### F6.3 — Sprint Retrospective Page Auto-creation

**Jira+Confluence analog:** Sprint tugaganda retrospective page avtomatik yaratiladi

**Vazifalar:**
- [ ] Sprint complete usecase'ida: project'ga bog'liq space topilsa, retrospective page yaratish
- [ ] Template: "Sprint X Retrospective" — What went well, What to improve, Action items
- [ ] `POST /api/v1/sprints/:id/complete` — `create_retrospective: true` parametr

---

### F6.4 — Unified Activity Feed

**Jira+Confluence analog:** "What's happening" — barcha o'zgarishlar bitta joyda

**Vazifalar:**
- [ ] `GET /api/v1/activity?project_id=...&space_id=...` — issues + pages birlashgan feed
- [ ] Pagination, `since` parametr

---

### F6.5 — Cross-system Mention

**Vazifalar:**
- [ ] Page'da `[PROJ-42]` yozilsa — issue'ga link render qilish
- [ ] Issue description'da `[SPACE:page-title]` — page'ga link
- [ ] Backend: mention parse va `entity_id` aniqlash

---

## FAZA 7 — Enterprise Features

---

### F7.1 — SSO / SAML / OAuth2

**Vazifalar:**
- [ ] Google OAuth2 login
- [ ] Microsoft Azure AD SAML
- [ ] Generic OIDC provider support
- [ ] `POST /api/v1/auth/oauth/google`

---

### F7.2 — Advanced RBAC / Permission Scheme

**Hozirgi holat:** Global role (admin/member/viewer) — loyiha darajasida yetarli emas.

**Vazifalar:**
- [ ] Permission scheme: `can_create_issue`, `can_transition_issue`, `can_manage_sprints`, etc.
- [ ] Project-level role customization
- [ ] Space-level fine-grained permissions
- [ ] Casbin policy'ni kengaytirish

---

### F7.3 — Audit Log Yaxshilanishi

**Hozirgi holat:** Audit log bor, lekin API chekli.

**Vazifalar:**
- [ ] Filter: `action_type`, `user_id`, `entity_type`, `date_from`, `date_to`
- [ ] Export: CSV format
- [ ] Retention policy (eski loglarni tozalash)
- [ ] `GET /api/v1/audit-logs/export?format=csv`

---

### F7.4 — Rate Limiting Yaxshilanishi

**Hozirgi holat:** Global rate limit bor.

**Vazifalar:**
- [ ] Per-user rate limiting (Redis counter)
- [ ] Per-endpoint rate limiting (upload, search uchun alohida)
- [ ] Rate limit headers: `X-RateLimit-Limit`, `X-RateLimit-Remaining`

---

### F7.5 — API Keys (Service Accounts)

**Jira analog:** API tokens for CI/CD, external tools

**Vazifalar:**
- [ ] Migration: `api_keys` jadval
- [ ] `POST /api/v1/api-keys` — key generatsiya
- [ ] `GET /api/v1/api-keys`
- [ ] `DELETE /api/v1/api-keys/:id`
- [ ] Auth middleware'da Bearer token va API key ikkalasini qabul qilish

---

### F7.6 — Data Import

**Vazifalar:**
- [ ] Jira XML/JSON export'ni import qilish
- [ ] Trello board import (JSON)
- [ ] Linear import (CSV)
- [ ] `POST /api/v1/import/jira`

---

### F7.7 — Global Site Settings

**Vazifalar:**
- [ ] `site_settings` jadval: `key-value` store
- [ ] `GET /api/v1/admin/settings`
- [ ] `PUT /api/v1/admin/settings`
- [ ] Sozlamalar: `max_file_size`, `allowed_domains`, `default_workflow`, `smtp_config`

---

## Texnik Qarz (Har doim)

---

### T1 — Test Coverage

**Vazifalar:**
- [ ] Unit tests: har bir usecase uchun (mock repository bilan)
- [ ] Integration tests: postgres container bilan (testcontainers-go)
- [ ] HTTP tests: handler'lar uchun (httptest)
- [ ] Maqsad: >80% coverage

---

### T2 — Migration Naming Convention

**Hozirgi holat:** `000014_issue_position.up.sql` — nomerlash ketayapti.  
**Taklif:** Shu tartibni davom ettirish.

---

### T3 — Error Handling Yaxshilanishi

**Vazifalar:**
- [ ] `apperr` package kengaytirish: `Conflict`, `TooManyRequests` error type'lari
- [ ] Barcha repository error'larini wrap qilish (context yo'qolmasin)
- [ ] Structured error response: `{ error, code, details, request_id }`

---

### T4 — Swagger Yangilash

**Vazifalar:**
- [ ] Har yangi endpoint uchun swagger annotation
- [ ] Request/Response example'lar qo'shish
- [ ] Error response'larni hujjatlash

---

## Yakuniy Holat Maqsad

| Soha | Hozir | Faza 1-2 | Faza 1-4 | Faza 1-7 |
|------|-------|----------|----------|----------|
| **Jira** | 55% | 70% | 90% | 100% |
| **Confluence** | 40% | 55% | 80% | 100% |

---

## Qaysi Fazani Birinchi Boshlash Kerak?

```
HOZIR BOSHLASH ↓

1. F1.1 WebSocket route      ← 1 kun
2. F1.2 Issue ordering       ← 1 kun
3. F1.3 Page watchers API    ← 1 kun
4. F2.1 Time tracking        ← 2-3 kun
5. F2.2 Components           ← 1-2 kun
6. F2.3 Versions/Releases    ← 2-3 kun
7. F2.8 Sprint reports       ← 2-3 kun
8. F3.1 Page templates       ← 2-3 kun
9. F3.4 Favorites            ← 1 kun
10. F6.1 Issue↔Page link     ← 1-2 kun
```

**Bu 10 ta feature ~15-20 ish kuni.  
Jira 85% + Confluence 70% darajasiga chiqadi.**
