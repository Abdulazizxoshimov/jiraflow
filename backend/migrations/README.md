# Database Migrations — Jira & Confluence Backend

## Tuzilma (12 ta migration, 31 ta jadval)

| # | Migration | Nima qo'shadi |
|---|-----------|---------------|
| 001 | `init_extensions` | pgcrypto, citext, pg_trgm, btree_gin + `set_updated_at()` trigger fn |
| 002 | `users_and_auth` | users, refresh_tokens, password_resets, invites |
| 003 | `workflows` | workflows, workflow_statuses, workflow_transitions |
| 004 | `projects` | projects, project_members, custom_fields |
| 005 | `sprints` | sprints (bir loyihada faqat 1 ta active) |
| 006 | `issues` | issues, labels, issue_labels, issue_watchers, issue_links, issue_history |
| 007 | `spaces_and_pages` | spaces, space_members, pages, page_versions, page_watchers |
| 008 | `comments` | comments (polymorphic), comment_mentions |
| 009 | `attachments` | attachments (issue/page/comment uchun) |
| 010 | `notifications` | notifications, notification_preferences |
| 011 | `boards_and_audit` | boards, board_columns, board_column_statuses, audit_logs |
| 012 | `seed_default_workflow` | Standart To Do → In Progress → In Review → Done |

## Ishga tushirish

### golang-migrate bilan
```bash
# O'rnatish
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Up
migrate -path ./migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Down (bitta)
migrate -path ./migrations \
  -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1
```

### Yoki psql bilan
```bash
for f in migrations/*.up.sql; do
  psql -d mydb -v ON_ERROR_STOP=1 -f "$f"
done
```

## Asosiy dizayn qarorlari

**UUID PK** — barcha jadvallarda `gen_random_uuid()`.

**Soft delete** — `deleted_at TIMESTAMPTZ` ko'pchilik jadvallarda. Indekslar `WHERE deleted_at IS NULL` bilan partial qilingan — performance uchun.

**Audit trail** — `created_at`/`updated_at` + alohida `issue_history` va `page_versions` jadvallar. Trigger `set_updated_at()` har bir UPDATE'da ishlaydi.

**Full-text search** — `pages.search_vector` va `issues.search_vector` — `GENERATED ALWAYS AS ... STORED` — avtomatik to'ldiriladi, alohida trigger kerak emas. GIN index bilan tez qidiriladi.

**Polymorphic associations** — `comments` va `attachments` `parent_type` + `parent_id` orqali issue/page'ga bog'lanadi. FK constraint yo'q (PostgreSQL polymorphic FK qo'llab-quvvatlamaydi), bu kod tomonidan tekshiriladi.

**Custom fields** — `issues.custom_fields JSONB` + GIN index. Loyiha sozlamalari `custom_fields` jadvalida.

**Active sprint constraint** — `uq_sprints_one_active_per_project` — partial unique index `WHERE status='active'`. Bir loyihada faqat 1 ta active sprint.

**Workflow flexibility** — `workflow_transitions.from_status_id` NULL bo'lsa, istalgan holatdan o'tish mumkin (global transition, masalan "Force close").

## Tekshirilgan

Barcha migration'lar PostgreSQL 16'da test qilingan:
- ✅ Up: 12/12 muvaffaqiyatli
- ✅ Down: 12/12 toza rollback
- ✅ Constraints ishlaydi (key regex, unique, check)
- ✅ Full-text search ishlaydi
- ✅ Sprint uniqueness ishlaydi

## Keyingi qadam

Go kodi: `internal/platform/database/` — pgx pool, transaction helpers, migration runner.
sqlc bilan typed query'larni generatsiya qilamiz.
