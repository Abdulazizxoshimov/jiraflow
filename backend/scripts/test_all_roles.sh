#!/usr/bin/env bash
# =============================================================================
# Barcha rollar uchun API endpoint test script
# Backend: http://localhost:8080
# =============================================================================
set -uo pipefail

BASE="http://localhost:8080/api/v1"
PASS=0
FAIL=0
SKIP=0

# Unique suffix to avoid 409 conflicts on repeated runs
TS=$(date +%s | tail -c 6)

# ── Ranglar ──────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# ── Ma'lumotlar ───────────────────────────────────────────────────────────────
WEBDEV_PROJECT="c1000000-0000-0000-0000-000000000001"
MOBILE_PROJECT="c2000000-0000-0000-0000-000000000002"
SPRINT1="d1000000-0000-0000-0000-000000000001"
SPRINT4="d4000000-0000-0000-0000-000000000004"
ISSUE_BUG="1000000b-0000-0000-0000-000000000000"
ISSUE_TASK="10000005-0000-0000-0000-000000000000"
BOARD1="30000001-0000-0000-0000-000000000000"
SPACE1="50000001-0000-0000-0000-000000000000"
PAGE1="60000001-0000-0000-0000-000000000000"
COMMENT1="70000001-0000-0000-0000-000000000000"
STATUS_TODO="433c09ee-9223-41d3-afeb-df59a2336531"
STATUS_INPROG="3d1c059c-c5be-415c-b848-773ba5f5fc71"
STATUS_INREVIEW="7a1f6e46-05d5-4d07-b50d-18c6e39a7040"
STATUS_DONE="6703b413-4fbf-4577-9b6f-2b5b994ad00d"
MEMBER1_ID="b1000000-0000-0000-0000-000000000001"
MEMBER2_ID="b2000000-0000-0000-0000-000000000002"

# ── Yordamchi funksiyalar ─────────────────────────────────────────────────────
check() {
  local label="$1"
  local expected="$2"
  local actual="$3"

  if [[ "$actual" == "$expected" ]]; then
    echo -e "  ${GREEN}✓${NC} $label (HTTP $actual)"
    ((PASS++))
  elif [[ "$actual" == "SKIP" ]]; then
    echo -e "  ${YELLOW}~${NC} $label (skipped)"
    ((SKIP++))
  else
    echo -e "  ${RED}✗${NC} $label  expected=${expected}  got=${actual}"
    ((FAIL++))
  fi
}

get_code() {
  curl -s -o /dev/null -w "%{http_code}" "$@"
}

login() {
  local email="$1" pass="$2"
  local resp
  resp=$(curl -s -X POST "$BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$pass\"}")
  echo "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d.get('data',{}).get('access_token',''))" 2>/dev/null || true
}

auth_get()  { get_code -H "Authorization: Bearer $1" "${@:2}"; }
auth_post() { get_code -s -X POST -H "Authorization: Bearer $1" -H "Content-Type: application/json" "${@:2}"; }
auth_put()  { get_code -s -X PUT  -H "Authorization: Bearer $1" -H "Content-Type: application/json" "${@:2}"; }
auth_del()  { get_code -s -X DELETE -H "Authorization: Bearer $1" "${@:2}"; }

section() {
  echo ""
  echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════${NC}"
  echo -e "${BOLD}${CYAN}  $1${NC}"
  echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════${NC}"
}

role_header() {
  echo ""
  echo -e "${BOLD}  ▶  Role: $1${NC}  ($2)"
}

# =============================================================================
# 0. Health check
# =============================================================================
section "0. HEALTH CHECK"
code=$(get_code "http://localhost:8080/health")
check "GET /health" "200" "$code"

# =============================================================================
# 1. AUTH — barcha rollar login
# =============================================================================
section "1. AUTH — LOGIN"

declare -A TOKENS

echo "  Logging in all users..."
TOKENS[admin]=$(login "admin@jiraflow.com" "Admin123!")
TOKENS[member1]=$(login "member1@jiraflow.com" "Member123!")
TOKENS[member2]=$(login "member2@jiraflow.com" "Member123!")
TOKENS[pm]=$(login "pm@jiraflow.com" "Manager123!")
TOKENS[dev]=$(login "dev@jiraflow.com" "Dev123!pass")
TOKENS[qa]=$(login "qa@jiraflow.com" "Qa123!pass1")
TOKENS[viewer]=$(login "viewer@jiraflow.com" "Viewer123!")

for role in admin member1 member2 pm dev qa viewer; do
  tok="${TOKENS[$role]}"
  if [[ -n "$tok" ]]; then
    echo -e "  ${GREEN}✓${NC} $role logged in"
    ((PASS++))
  else
    echo -e "  ${RED}✗${NC} $role login FAILED"
    ((FAIL++))
  fi
done

# Auth me — barcha rollar
section "2. GET /auth/me — barcha rollar"
for role in admin member1 member2 pm dev qa viewer; do
  tok="${TOKENS[$role]}"
  code=$(auth_get "$tok" "$BASE/auth/me")
  check "GET /auth/me [$role]" "200" "$code"
done

# Refresh token
section "3. POST /auth/refresh"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@jiraflow.com","password":"Admin123!"}')
check "POST /auth/login (admin)" "200" "$code"

# Forgot password (public, har kim)
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/auth/forgot-password" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@jiraflow.com"}')
check "POST /auth/forgot-password (public)" "204" "$code"

# =============================================================================
# 4. USERS
# =============================================================================
section "4. USERS"

role_header "admin" "global admin"
TOK="${TOKENS[admin]}"
check "GET /users [admin]" "200" "$(auth_get "$TOK" "$BASE/users")"
check "GET /users/:id [admin]" "200" "$(auth_get "$TOK" "$BASE/users/$MEMBER1_ID")"

role_header "member1" "global member"
TOK="${TOKENS[member1]}"
check "GET /users [member1]" "200" "$(auth_get "$TOK" "$BASE/users")"
check "GET /users/:id [member1]" "200" "$(auth_get "$TOK" "$BASE/users/$MEMBER2_ID")"

role_header "viewer" "global viewer"
TOK="${TOKENS[viewer]}"
check "GET /users [viewer]" "200" "$(auth_get "$TOK" "$BASE/users")"

# Admin faqat user yarata oladi
TOK="${TOKENS[admin]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/users" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"testcreate${TS}@jiraflow.com\",\"password\":\"Test123!pass\",\"full_name\":\"Test Create\",\"role\":\"member\"}")
check "POST /users [admin]" "201" "$code"

# Non-admin user yaratishga urinish (403 yoki 201 bo'lishi mumkin)
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/users" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d '{"email":"shouldfail@test.com","password":"Test123!pass","full_name":"Fail","role":"member"}')
check "POST /users [viewer → 403]" "403" "$code"

# =============================================================================
# 5. PROJECTS
# =============================================================================
section "5. PROJECTS"

for role in admin member1 member2 pm dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/projects")
  check "GET /projects [$role]" "200" "$code"
done

for role in admin member1 member2 pm dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT")
  check "GET /projects/:id WEBDEV [$role]" "200" "$code"
done

# Project yaratish — faqat admin va member
TOK="${TOKENS[admin]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/projects" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"key\":\"TP${TS}\",\"name\":\"Test Project ${TS}\",\"description\":\"test\"}")
check "POST /projects [admin]" "201" "$code"

# viewer project yarata olmaydi
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/projects" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d '{"key":"FAILVWR","name":"Fail Viewer","description":"test"}')
check "POST /projects [viewer → 403]" "403" "$code"

# Project members
for role in admin member1 member2 dev; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/members")
  check "GET /projects/:id/members WEBDEV [$role]" "200" "$code"
done

# Project sprints
for role in admin member1 dev qa; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/sprints")
  check "GET /projects/:id/sprints WEBDEV [$role]" "200" "$code"
done

# Project labels, components, versions
TOK="${TOKENS[admin]}"
check "GET /projects/:id/labels [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/labels")"
check "GET /projects/:id/components [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/components")"
check "GET /projects/:id/versions [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/versions")"
check "GET /projects/:id/custom-fields [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/custom-fields")"
check "GET /projects/:id/boards [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/boards")"
check "GET /projects/:id/backlog [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/backlog")"
check "GET /projects/:id/roadmap [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/roadmap")"
check "GET /projects/:id/dashboard [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/dashboard")"
check "GET /projects/:id/velocity [admin]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/velocity")"

# viewer uchun ham
TOK="${TOKENS[viewer]}"
check "GET /projects/:id/labels [viewer]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/labels")"
check "GET /projects/:id/backlog [viewer]" "200" "$(auth_get "$TOK" "$BASE/projects/$WEBDEV_PROJECT/backlog")"

# =============================================================================
# 6. SPRINTS
# =============================================================================
section "6. SPRINTS"

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/sprints/$SPRINT1")
  check "GET /sprints/:id [$role]" "200" "$code"
done

TOK="${TOKENS[admin]}"
check "GET /sprints/:id/report [admin]" "200" "$(auth_get "$TOK" "$BASE/sprints/$SPRINT1/report")"
check "GET /sprints/:id/burndown [admin]" "200" "$(auth_get "$TOK" "$BASE/sprints/$SPRINT1/burndown")"
check "GET /sprints/:id/burnup [admin]" "200" "$(auth_get "$TOK" "$BASE/sprints/$SPRINT1/burnup")"
check "GET /sprints/:id/capacity [admin]" "200" "$(auth_get "$TOK" "$BASE/sprints/$SPRINT1/capacity")"
check "GET /sprints/:id/impediments [admin]" "200" "$(auth_get "$TOK" "$BASE/sprints/$SPRINT1/impediments")"

# Sprint yaratish
TOK="${TOKENS[member1]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/projects/$WEBDEV_PROJECT/sprints" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Sprint ${TS}\",\"goal\":\"test goal\",\"start_date\":\"2026-06-15T00:00:00Z\",\"end_date\":\"2026-06-29T00:00:00Z\"}")
check "POST /projects/:id/sprints [member1=project admin]" "201" "$code"

# qa yarata olmaydi (WEBDEV viewer, project-level restriction)
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/projects/$WEBDEV_PROJECT/sprints" \
  -H "Authorization: Bearer ${TOKENS[qa]}" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"QA Sprint Fail ${TS}\",\"goal\":\"should fail\",\"start_date\":\"2026-06-15T00:00:00Z\",\"end_date\":\"2026-06-29T00:00:00Z\"}")
check "POST /projects/:id/sprints [qa=viewer → 403]" "403" "$code"

# =============================================================================
# 7. ISSUES
# =============================================================================
section "7. ISSUES"

for role in admin member1 member2 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/issues?project_id=$WEBDEV_PROJECT")
  check "GET /issues?project_id=WEBDEV [$role]" "200" "$code"
done

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG")
  check "GET /issues/:id [$role]" "200" "$code"
done

# Issue key by key
code=$(auth_get "${TOKENS[admin]}" "$BASE/issues/key/WEBDEV-1")
check "GET /issues/key/:key [admin]" "200" "$code"

# Issue yaratish — member va admin
for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues" \
    -H "Authorization: Bearer $TOK" \
    -H "Content-Type: application/json" \
    -d "{\"project_id\":\"$WEBDEV_PROJECT\",\"title\":\"Test Issue by $role\",\"type\":\"task\",\"status_id\":\"$STATUS_TODO\",\"priority\":\"medium\"}")
  check "POST /issues [$role]" "201" "$code"
done

# viewer issue yarata olmaydi
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d "{\"project_id\":\"$WEBDEV_PROJECT\",\"title\":\"Viewer Issue Fail\",\"type\":\"task\",\"status_id\":\"$STATUS_TODO\",\"priority\":\"low\"}")
check "POST /issues [viewer → 403]" "403" "$code"

# Issue history, watchers, links
TOK="${TOKENS[admin]}"
check "GET /issues/:id/history [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/history")"
check "GET /issues/:id/watchers [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/watchers")"
check "GET /issues/:id/links [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/links")"
check "GET /issues/:id/worklogs [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/worklogs")"
check "GET /issues/:id/assignees [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/assignees")"
check "GET /issues/:id/votes [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/votes")"
check "GET /issues/:id/time-summary [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/time-summary")"
check "GET /issues/:id/page-links [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/page-links")"
check "GET /issues/:id/epic-progress [admin]" "200" "$(auth_get "$TOK" "$BASE/issues/10000001-0000-0000-0000-000000000000/epic-progress")"

# Issue transition — reset to In Progress first (idempotent via DB), then transition to In Review
psql postgres://postgres:4444@localhost:5432/jiraflow -c \
  "UPDATE issues SET status_id='$STATUS_INPROG' WHERE id='$ISSUE_TASK'" > /dev/null 2>&1

TOK="${TOKENS[dev]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues/$ISSUE_TASK/transition" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"status_id\":\"$STATUS_INREVIEW\"}")
check "POST /issues/:id/transition [dev=member]" "200" "$code"

# global viewer transition qila olmaydi — Casbin blocks it at 403
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues/$ISSUE_TASK/transition" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d "{\"status_id\":\"$STATUS_INREVIEW\"}")
check "POST /issues/:id/transition [viewer → 403]" "403" "$code"

# Worklog yaratish
TOK="${TOKENS[dev]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues/$ISSUE_BUG/worklogs" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"time_spent":3600,"description":"Investigating the bug","started_at":"2026-06-10T10:00:00Z"}')
check "POST /issues/:id/worklogs [dev]" "201" "$code"

# =============================================================================
# 8. COMMENTS
# =============================================================================
section "8. COMMENTS"

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/issues/$ISSUE_BUG/comments")
  check "GET /issues/:id/comments [$role]" "200" "$code"
done

# Comment yaratish
for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues/$ISSUE_BUG/comments" \
    -H "Authorization: Bearer $TOK" \
    -H "Content-Type: application/json" \
    -d "{\"content\":{\"type\":\"doc\",\"content\":[{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Comment from $role\"}]}]}}")
  check "POST /issues/:id/comments [$role]" "201" "$code"
done

# viewer comment yoza olmaydi (WEBDEV viewer)
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/issues/$ISSUE_BUG/comments" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d '{"content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Viewer fail"}]}]}}')
check "POST /issues/:id/comments [viewer → 403]" "403" "$code"

# GET comment by id
TOK="${TOKENS[admin]}"
check "GET /comments/:id [admin]" "200" "$(auth_get "$TOK" "$BASE/comments/$COMMENT1")"

# Comment reactions
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/comments/$COMMENT1/reactions" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"emoji":"👍"}')
check "POST /comments/:id/reactions [admin]" "204" "$code"
check "GET /comments/:id/reactions [admin]" "200" "$(auth_get "$TOK" "$BASE/comments/$COMMENT1/reactions")"

# =============================================================================
# 9. BOARDS
# =============================================================================
section "9. BOARDS"

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/boards/$BOARD1")
  check "GET /boards/:id [$role]" "200" "$code"
done

TOK="${TOKENS[admin]}"
check "GET /boards/:id/swimlanes [admin]" "200" "$(auth_get "$TOK" "$BASE/boards/$BOARD1/swimlanes")"

# Board yaratish
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/projects/$WEBDEV_PROJECT/boards" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Board 99","type":"kanban"}')
check "POST /projects/:id/boards [admin]" "201" "$code"

# =============================================================================
# 10. SPACES & PAGES
# =============================================================================
section "10. SPACES"

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/spaces")
  check "GET /spaces [$role]" "200" "$code"
done

for role in admin member1 dev qa; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/spaces/$SPACE1")
  check "GET /spaces/:id [$role]" "200" "$code"
done

TOK="${TOKENS[admin]}"
check "GET /spaces/:id/members [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/members")"
check "GET /spaces/:id/statistics [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/statistics")"
check "GET /spaces/:id/pages/tree [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/pages/tree")"
check "GET /spaces/:id/page-tags [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/page-tags")"
check "GET /spaces/:id/blog-posts [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/blog-posts")"
check "GET /spaces/:id/webhooks [admin]" "200" "$(auth_get "$TOK" "$BASE/spaces/$SPACE1/webhooks")"

# Space yaratish
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/spaces" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"key\":\"TS${TS}\",\"name\":\"Test Space ${TS}\",\"type\":\"team\"}")
check "POST /spaces [admin]" "201" "$code"

# viewer space yarata olmaydi
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/spaces" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d '{"key":"VIEWERFAIL","name":"Fail","type":"team"}')
check "POST /spaces [viewer → 403]" "403" "$code"

section "11. PAGES"

for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/pages/$PAGE1")
  check "GET /pages/:id [$role]" "200" "$code"
done

TOK="${TOKENS[admin]}"
check "GET /pages [admin]" "200" "$(auth_get "$TOK" "$BASE/pages?space_id=$SPACE1")"
check "GET /pages/:id/versions [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/versions")"
check "GET /pages/:id/watchers [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/watchers")"
check "GET /pages/:id/tags [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/tags")"
check "GET /pages/:id/analytics [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/analytics")"
check "GET /pages/:id/inline-comments [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/inline-comments")"
check "GET /pages/:id/restrictions [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/restrictions")"
check "GET /pages/:id/access [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/access")"
# Acquire a lock first so GET returns 200
curl -s -X POST "$BASE/pages/$PAGE1/lock" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"ttl_seconds":300}' > /dev/null 2>&1
check "GET /pages/:id/lock [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/lock")"
check "GET /pages/:id/reactions [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/reactions")"
check "GET /pages/:id/issue-links [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/issue-links")"
check "GET /pages/:id/macros [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/macros")"
check "GET /pages/:id/comments [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/comments")"

# Page yaratish
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/spaces/$SPACE1/pages" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Page 99","content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Test"}]}]}}')
check "POST /spaces/:id/pages [admin]" "201" "$code"

# Page export
check "GET /pages/:id/export/html [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/export/html")"
check "GET /pages/:id/export/md [admin]" "200" "$(auth_get "$TOK" "$BASE/pages/$PAGE1/export/md")"

# =============================================================================
# 12. NOTIFICATIONS
# =============================================================================
section "12. NOTIFICATIONS"

for role in admin member1 dev viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/notifications")
  check "GET /notifications [$role]" "200" "$code"
done

TOK="${TOKENS[member1]}"
check "GET /notifications/unread-count [member1]" "200" "$(auth_get "$TOK" "$BASE/notifications/unread-count")"
check "GET /notifications/preferences [member1]" "200" "$(auth_get "$TOK" "$BASE/notifications/preferences")"

# Mark all read
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/notifications/mark-all-read" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json")
check "POST /notifications/mark-all-read [member1]" "204" "$code"

# =============================================================================
# 13. SEARCH & ACTIVITY
# =============================================================================
section "13. SEARCH & ACTIVITY"

for role in admin member1 dev qa viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/search?q=auth")
  check "GET /search?q=auth [$role]" "200" "$code"
done

TOK="${TOKENS[admin]}"
check "GET /search/suggestions [admin]" "200" "$(auth_get "$TOK" "$BASE/search/suggestions?q=web")"
check "GET /activity [admin]" "200" "$(auth_get "$TOK" "$BASE/activity")"

# =============================================================================
# 14. AUDIT LOGS
# =============================================================================
section "14. AUDIT LOGS"

TOK="${TOKENS[admin]}"
check "GET /audit-logs [admin]" "200" "$(auth_get "$TOK" "$BASE/audit-logs")"
check "GET /audit-logs/export [admin]" "200" "$(auth_get "$TOK" "$BASE/audit-logs/export")"

# member audit loga kira olmaydi
code=$(auth_get "${TOKENS[member1]}" "$BASE/audit-logs")
check "GET /audit-logs [member1 → 403]" "403" "$code"

# viewer audit loga kira olmaydi
code=$(auth_get "${TOKENS[viewer]}" "$BASE/audit-logs")
check "GET /audit-logs [viewer → 403]" "403" "$code"

# =============================================================================
# 15. WORKFLOWS
# =============================================================================
section "15. WORKFLOWS"

for role in admin member1 viewer; do
  TOK="${TOKENS[$role]}"
  code=$(auth_get "$TOK" "$BASE/workflows")
  check "GET /workflows [$role]" "200" "$code"
done

WORKFLOW_ID="00000000-0000-0000-0000-000000000001"
TOK="${TOKENS[admin]}"
check "GET /workflows/:id [admin]" "200" "$(auth_get "$TOK" "$BASE/workflows/$WORKFLOW_ID")"

# Workflow yaratish
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/workflows" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Workflow 99","description":"test"}')
check "POST /workflows [admin]" "201" "$code"

# viewer yarata olmaydi
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/workflows" \
  -H "Authorization: Bearer ${TOKENS[viewer]}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Viewer Fail Workflow"}')
check "POST /workflows [viewer → 403]" "403" "$code"

# =============================================================================
# 16. FAVORITES
# =============================================================================
section "16. FAVORITES"

for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  check "GET /favorites [$role]" "200" "$(auth_get "$TOK" "$BASE/favorites")"
done

TOK="${TOKENS[member1]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/favorites" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"entity_type\":\"issue\",\"entity_id\":\"$ISSUE_BUG\"}")
check "POST /favorites (issue) [member1]" "201" "$code"

code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/favorites" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"entity_type\":\"page\",\"entity_id\":\"$PAGE1\"}")
check "POST /favorites (page) [member1]" "201" "$code"

check "GET /favorites/check?entity_type=issue&entity_id=$ISSUE_BUG [member1]" "200" \
  "$(auth_get "$TOK" "$BASE/favorites/check?entity_type=issue&entity_id=$ISSUE_BUG")"

check "GET /recently-visited [member1]" "200" "$(auth_get "$TOK" "$BASE/recently-visited")"

# =============================================================================
# 17. API KEYS
# =============================================================================
section "17. API KEYS"

for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  check "GET /api-keys [$role]" "200" "$(auth_get "$TOK" "$BASE/api-keys")"
done

TOK="${TOKENS[member1]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api-keys" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test API Key 99","description":"for testing"}')
check "POST /api-keys [member1]" "201" "$code"

# =============================================================================
# 18. PERMISSION SCHEMES
# =============================================================================
section "18. PERMISSION SCHEMES"

TOK="${TOKENS[admin]}"
check "GET /permission-schemes [admin]" "200" "$(auth_get "$TOK" "$BASE/permission-schemes")"

code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/permission-schemes" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Scheme 99","description":"test"}')
check "POST /permission-schemes [admin]" "201" "$code"

# non-admin
code=$(auth_get "${TOKENS[member1]}" "$BASE/permission-schemes")
check "GET /permission-schemes [member1 → 403]" "403" "$code"

# =============================================================================
# 19. SAVED FILTERS
# =============================================================================
section "19. SAVED FILTERS"

for role in admin member1 dev; do
  TOK="${TOKENS[$role]}"
  check "GET /saved-filters [$role]" "200" "$(auth_get "$TOK" "$BASE/saved-filters")"
done

TOK="${TOKENS[dev]}"
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/saved-filters" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"My Bugs\",\"project_id\":\"$WEBDEV_PROJECT\",\"filter\":{\"type\":\"bug\"}}")
check "POST /saved-filters [dev]" "201" "$code"

# =============================================================================
# 20. ISSUE TYPES, TYPE SCHEMES, FIELD CONFIGS
# =============================================================================
section "20. ISSUE TYPES & SCHEMES"

TOK="${TOKENS[admin]}"
check "GET /issue-types [admin]" "200" "$(auth_get "$TOK" "$BASE/issue-types")"
check "GET /issue-type-schemes [admin]" "200" "$(auth_get "$TOK" "$BASE/issue-type-schemes")"
check "GET /field-configurations [admin]" "200" "$(auth_get "$TOK" "$BASE/field-configurations")"
check "GET /notification-schemes [admin]" "200" "$(auth_get "$TOK" "$BASE/notification-schemes")"
check "GET /security-schemes [admin]" "200" "$(auth_get "$TOK" "$BASE/security-schemes")"
check "GET /project-templates [admin]" "200" "$(auth_get "$TOK" "$BASE/project-templates")"
check "GET /space-categories [admin]" "200" "$(auth_get "$TOK" "$BASE/space-categories")"
check "GET /blueprints [admin]" "200" "$(auth_get "$TOK" "$BASE/blueprints")"
check "GET /page-templates [admin]" "200" "$(auth_get "$TOK" "$BASE/page-templates")"

# member ham o'qiy oladi
check "GET /issue-types [member1]" "200" "$(auth_get "${TOKENS[member1]}" "$BASE/issue-types")"

# viewer
check "GET /project-templates [viewer]" "200" "$(auth_get "${TOKENS[viewer]}" "$BASE/project-templates")"

# =============================================================================
# 21. INVITES
# =============================================================================
section "21. INVITES"

TOK="${TOKENS[admin]}"
check "GET /invites [admin]" "200" "$(auth_get "$TOK" "$BASE/invites")"

# Invite yaratish
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/invites" \
  -H "Authorization: Bearer $TOK" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"invitee${TS}@test.com\",\"role\":\"member\"}")
check "POST /invites [admin]" "201" "$code"

# member invite yarata olmaydi
code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/invites" \
  -H "Authorization: Bearer ${TOKENS[member1]}" \
  -H "Content-Type: application/json" \
  -d '{"email":"shouldfail@test.com","role":"member"}')
check "POST /invites [member1 → 403]" "403" "$code"

# =============================================================================
# 22. SWAGGER & PUBLIC ENDPOINTS
# =============================================================================
section "22. SWAGGER & PUBLIC"

code=$(get_code "http://localhost:8080/swagger/index.html")
check "GET /swagger/index.html (public)" "200" "$code"

# Unauthenticated — protected route
code=$(get_code "$BASE/projects")
check "GET /projects (no auth → 401)" "401" "$code"

code=$(get_code "$BASE/issues")
check "GET /issues (no auth → 401)" "401" "$code"

code=$(get_code "$BASE/spaces")
check "GET /spaces (no auth → 401)" "401" "$code"

# =============================================================================
# YAKUNIY HISOBOT
# =============================================================================
echo ""
echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════${NC}"
echo -e "${BOLD}  YAKUNIY NATIJA${NC}"
echo -e "${BOLD}${CYAN}══════════════════════════════════════════════════════${NC}"
echo -e "  ${GREEN}✓ PASSED : $PASS${NC}"
echo -e "  ${RED}✗ FAILED : $FAIL${NC}"
echo -e "  ${YELLOW}~ SKIPPED: $SKIP${NC}"
echo ""
TOTAL=$((PASS + FAIL))
if [[ $FAIL -eq 0 ]]; then
  echo -e "  ${GREEN}${BOLD}Barcha testlar muvaffaqiyatli o'tdi!${NC}"
else
  PCT=$(( PASS * 100 / TOTAL ))
  echo -e "  ${YELLOW}${BOLD}Muvaffaqiyat: ${PCT}% (${PASS}/${TOTAL})${NC}"
fi
echo ""
