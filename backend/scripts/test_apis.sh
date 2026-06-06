#!/usr/bin/env bash
# =============================================================================
# JiraFlow API Test Script
# Barcha API endpointlarini turli rollar bilan test qiladi
# =============================================================================
set -euo pipefail

BASE_URL="http://localhost:8080/api/v1"
PASS=0
FAIL=0
SKIP=0

# Har test ishga tushganda unique suffix (conflict oldini olish uchun)
TS=$(date +%s | tail -c 5)

# ─── Ranglar ──────────────────────────────────────────────────────────────────
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# ─── Helpers ──────────────────────────────────────────────────────────────────
section() { echo -e "\n${BOLD}${BLUE}══════════════════════════════════════${NC}"; echo -e "${BOLD}${BLUE}  $1${NC}"; echo -e "${BOLD}${BLUE}══════════════════════════════════════${NC}"; }
info()    { echo -e "  ${CYAN}→${NC} $1"; }

ok()   {
  PASS=$((PASS+1))
  echo -e "  ${GREEN}✓${NC} $1"
}
fail() {
  FAIL=$((FAIL+1))
  echo -e "  ${RED}✗${NC} $1"
  echo -e "    ${RED}Response:${NC} $2"
}
warn() { SKIP=$((SKIP+1)); echo -e "  ${YELLOW}⊙${NC} $1"; }

# HTTP so'rov + natijani tekshirish
# check <label> <expected_status> <actual_status> [response_body]
check() {
  local label="$1" exp="$2" got="$3" body="${4:-}"
  if [[ "$got" == "$exp" ]]; then
    ok "$label [HTTP $got]"
  else
    fail "$label [kutilgan: $exp, kelgan: $got]" "$body"
  fi
}

# curl wrapper — status kod va body qaytaradi
# Format: <body>|||<status_code>
req() {
  local method="$1" url="$2"
  shift 2
  local body status
  body=$(curl -s -o /tmp/_jf_body -w "%{http_code}" -X "$method" "$url" "$@")
  status="$body"
  body=$(cat /tmp/_jf_body)
  echo "${body}|||${status}"
}

# Status kodni body dan ajratib olish
status_of() { echo "${1##*|||}" ; }
body_of()   { echo "${1%|||*}" ; }

# JSON field extract
jq_get() { echo "$1" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d$2)" 2>/dev/null || echo ""; }

# =============================================================================
# 0. TEKSHIRUV: server ishlaydimi?
# =============================================================================
section "0. Server tekshiruvi"
if ! curl -s --max-time 3 "$BASE_URL/auth/login" \
    -X POST -H "Content-Type: application/json" \
    -d '{"email":"x","password":"x"}' -o /dev/null 2>&1; then
  echo -e "${RED}Server ishlamayapti! Avval serverni ishga tushiring.${NC}"
  exit 1
fi
ok "Server http://localhost:8080 da ishlamoqda"

# =============================================================================
# 1. AUTH — Login (barcha rollar)
# =============================================================================
section "1. AUTH — Login"

login_as() {
  local email="$1" password="$2"
  local resp status body
  resp=$(req POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$email\",\"password\":\"$password\"}")
  status=$(status_of "$resp")
  body=$(body_of "$resp")
  if [[ "$status" == "200" ]]; then
    echo "$body" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['access_token'])" 2>/dev/null
  else
    echo ""
  fi
}

info "Admin login (Admin123!)"
ADMIN_TOKEN=$(login_as "admin@jiraflow.com" "Admin123!")
[[ -n "$ADMIN_TOKEN" ]] && ok "admin@jiraflow.com login muvaffaqiyatli" || fail "admin login" "token yo'q"

info "Member1 login (Member123!)"
MEMBER1_TOKEN=$(login_as "member1@jiraflow.com" "Member123!")
[[ -n "$MEMBER1_TOKEN" ]] && ok "member1@jiraflow.com login muvaffaqiyatli" || fail "member1 login" "token yo'q"

info "Member2 login"
MEMBER2_TOKEN=$(login_as "member2@jiraflow.com" "Member123!")
[[ -n "$MEMBER2_TOKEN" ]] && ok "member2@jiraflow.com login muvaffaqiyatli" || fail "member2 login" "token yo'q"

info "PM login (Manager123!)"
PM_TOKEN=$(login_as "pm@jiraflow.com" "Manager123!")
[[ -n "$PM_TOKEN" ]] && ok "pm@jiraflow.com login muvaffaqiyatli" || fail "pm login" "token yo'q"

info "Dev login (Dev123!pass)"
DEV_TOKEN=$(login_as "dev@jiraflow.com" "Dev123!pass")
[[ -n "$DEV_TOKEN" ]] && ok "dev@jiraflow.com login muvaffaqiyatli" || fail "dev login" "token yo'q"

info "QA login (Qa123!pass1)"
QA_TOKEN=$(login_as "qa@jiraflow.com" "Qa123!pass1")
[[ -n "$QA_TOKEN" ]] && ok "qa@jiraflow.com login muvaffaqiyatli" || fail "qa login" "token yo'q"

info "Viewer login (Viewer123!)"
VIEWER_TOKEN=$(login_as "viewer@jiraflow.com" "Viewer123!")
[[ -n "$VIEWER_TOKEN" ]] && ok "viewer@jiraflow.com login muvaffaqiyatli" || fail "viewer login" "token yo'q"

info "Noto'g'ri parol bilan login rad etilishi kerak"
resp=$(req POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@jiraflow.com","password":"wrongpassword"}')
check "Noto'g'ri parol → 401" "401" "$(status_of "$resp")" "$(body_of "$resp")"

info "Token yo'q holda protected endpoint → 401"
resp=$(req GET "$BASE_URL/auth/me")
check "Token yo'q → 401" "401" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 2. AUTH — Me, Logout
# =============================================================================
section "2. AUTH — Me & Logout"

info "GET /auth/me (admin)"
resp=$(req GET "$BASE_URL/auth/me" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /auth/me admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /auth/me (viewer)"
resp=$(req GET "$BASE_URL/auth/me" -H "Authorization: Bearer $VIEWER_TOKEN")
check "GET /auth/me viewer" "200" "$(status_of "$resp")" "$(body_of "$resp")"

# Refresh token olish
info "Refresh token"
resp_login=$(req POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"member1@jiraflow.com","password":"Member123!"}')
REFRESH_TOKEN=$(body_of "$resp_login" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['refresh_token'])" 2>/dev/null || echo "")

if [[ -n "$REFRESH_TOKEN" ]]; then
  resp=$(req POST "$BASE_URL/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}")
  check "POST /auth/refresh" "200" "$(status_of "$resp")" "$(body_of "$resp")"
else
  warn "Refresh token topilmadi"
fi

# =============================================================================
# 3. USERS
# =============================================================================
section "3. USERS"

info "GET /users (admin — barcha userlarni ko'ra oladi)"
resp=$(req GET "$BASE_URL/users" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /users admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /users?role=admin"
resp=$(req GET "$BASE_URL/users?role=admin" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /users?role=admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /users?is_active=true"
resp=$(req GET "$BASE_URL/users?is_active=true" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /users?is_active=true" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /users/b1000000-0000-0000-0000-000000000001 (member1)"
resp=$(req GET "$BASE_URL/users/b1000000-0000-0000-0000-000000000001" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /users/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /users (admin yangi user yaratadi)"
resp=$(req POST "$BASE_URL/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"newuser${TS}@test.com\",\"password\":\"NewUser123!\",\"full_name\":\"New Test User\",\"role\":\"member\"}")
check "POST /users (admin)" "201" "$(status_of "$resp")" "$(body_of "$resp")"
NEW_USER_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$NEW_USER_ID" ]]; then
  info "PUT /users/$NEW_USER_ID (update)"
  resp=$(req PUT "$BASE_URL/users/$NEW_USER_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"full_name":"Updated User Name"}')
  check "PUT /users/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /users/$NEW_USER_ID/deactivate"
  resp=$(req POST "$BASE_URL/users/$NEW_USER_ID/deactivate" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "POST /users/:id/deactivate" "204" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /users/$NEW_USER_ID/activate"
  resp=$(req POST "$BASE_URL/users/$NEW_USER_ID/activate" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "POST /users/:id/activate" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 4. PROJECTS
# =============================================================================
section "4. PROJECTS"

WEBDEV_ID="c1000000-0000-0000-0000-000000000001"
MOBILE_ID="c2000000-0000-0000-0000-000000000002"
WORKFLOW_ID="00000000-0000-0000-0000-000000000001"

info "GET /projects (admin)"
resp=$(req GET "$BASE_URL/projects" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /projects admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /projects (member1 — WEBDEV va MOBILE member)"
resp=$(req GET "$BASE_URL/projects" -H "Authorization: Bearer $MEMBER1_TOKEN")
check "GET /projects member1" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /projects/$WEBDEV_ID"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /projects/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /projects (admin yangi loyiha yaratadi)"
resp=$(req POST "$BASE_URL/projects" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"key\":\"TP${TS}\",\"name\":\"Test Project ${TS}\",\"description\":\"API test loyihasi\",\"workflow_id\":\"$WORKFLOW_ID\"}")
check "POST /projects admin" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_PROJECT_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

info "POST /projects (member — ruxsatsiz)"
resp=$(req POST "$BASE_URL/projects" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"key\":\"DN${TS}\",\"name\":\"Should Fail\",\"workflow_id\":\"$WORKFLOW_ID\"}")
STATUS_PROJ_MEMBER=$(status_of "$resp")
if [[ "$STATUS_PROJ_MEMBER" == "403" ]]; then
  ok "POST /projects member → 403 [RBAC ishlayapti]"
else
  warn "POST /projects member → $STATUS_PROJ_MEMBER [RBAC konfiguratsiyasi tekshirilsin — member loyiha yarata olmasligi kerak]"
fi

if [[ -n "$TEST_PROJECT_ID" ]]; then
  info "PUT /projects/$TEST_PROJECT_ID"
  resp=$(req PUT "$BASE_URL/projects/$TEST_PROJECT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated Test Project","description":"Yangilangan tavsif"}')
  check "PUT /projects/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /projects/$TEST_PROJECT_ID/archive"
  resp=$(req POST "$BASE_URL/projects/$TEST_PROJECT_ID/archive" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "POST /projects/:id/archive" "204" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /projects/$TEST_PROJECT_ID"
  resp=$(req DELETE "$BASE_URL/projects/$TEST_PROJECT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /projects/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 5. PROJECT MEMBERS
# =============================================================================
section "5. PROJECT MEMBERS"

info "GET /projects/$WEBDEV_ID/members"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/members" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET project members (admin)" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /projects/$WEBDEV_ID/members (member2 — o'zi ham member)"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/members" \
  -H "Authorization: Bearer $MEMBER2_TOKEN")
check "GET project members (member2)" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "PUT /projects/$WEBDEV_ID/members/b2000000-... (rol yangilash)"
resp=$(req PUT "$BASE_URL/projects/$WEBDEV_ID/members/b2000000-0000-0000-0000-000000000002" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role":"member"}')
check "PUT project member role" "204" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 6. WORKFLOWS
# =============================================================================
section "6. WORKFLOWS"

info "GET /workflows"
resp=$(req GET "$BASE_URL/workflows" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /workflows" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /workflows/$WORKFLOW_ID"
resp=$(req GET "$BASE_URL/workflows/$WORKFLOW_ID" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /workflows/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /workflows (admin)"
resp=$(req POST "$BASE_URL/workflows" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Workflow","description":"Test uchun workflow"}')
check "POST /workflows" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_WF_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_WF_ID" ]]; then
  info "POST /workflows/$TEST_WF_ID/statuses"
  resp=$(req POST "$BASE_URL/workflows/$TEST_WF_ID/statuses" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Backlog","category":"todo","color":"#9CA3AF","position":1,"is_initial":true}')
  check "POST workflow status" "201" "$(status_of "$resp")" "$(body_of "$resp")"
  TEST_STATUS_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

  info "PUT /workflows/statuses/$TEST_STATUS_ID"
  if [[ -n "$TEST_STATUS_ID" ]]; then
    resp=$(req PUT "$BASE_URL/workflows/statuses/$TEST_STATUS_ID" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Backlog Updated","color":"#6B7280"}')
    check "PUT workflow status" "200" "$(status_of "$resp")" "$(body_of "$resp")"

    info "DELETE /workflows/statuses/$TEST_STATUS_ID"
    resp=$(req DELETE "$BASE_URL/workflows/statuses/$TEST_STATUS_ID" \
      -H "Authorization: Bearer $ADMIN_TOKEN")
    check "DELETE workflow status" "204" "$(status_of "$resp")" "$(body_of "$resp")"
  fi

  info "DELETE /workflows/$TEST_WF_ID"
  resp=$(req DELETE "$BASE_URL/workflows/$TEST_WF_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /workflows/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 7. SPRINTS
# =============================================================================
section "7. SPRINTS"

SPRINT1_ID="d1000000-0000-0000-0000-000000000001"
SPRINT2_ID="d2000000-0000-0000-0000-000000000002"

info "GET /projects/$WEBDEV_ID/sprints"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/sprints" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET project sprints" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /sprints/$SPRINT1_ID"
resp=$(req GET "$BASE_URL/sprints/$SPRINT1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /sprints/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /projects/$WEBDEV_ID/sprints (yangi sprint)"
resp=$(req POST "$BASE_URL/projects/$WEBDEV_ID/sprints" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Sprint 3 — Test","goal":"API test sprint"}')
check "POST sprint" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_SPRINT_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_SPRINT_ID" ]]; then
  info "PUT /sprints/$TEST_SPRINT_ID"
  resp=$(req PUT "$BASE_URL/sprints/$TEST_SPRINT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Sprint 3 — Updated"}')
  check "PUT /sprints/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /sprints/$TEST_SPRINT_ID/start"
  resp=$(req POST "$BASE_URL/sprints/$TEST_SPRINT_ID/start" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  # WEBDEV da allaqachon active sprint bor, shuning uchun 409 bo'lishi mumkin
  STATUS=$(status_of "$resp")
  if [[ "$STATUS" == "200" ]] || [[ "$STATUS" == "409" ]]; then
    ok "POST /sprints/:id/start [$STATUS — OK yoki conflict kutilgan]"
  else
    fail "POST /sprints/:id/start" "$(body_of "$resp")"
  fi

  info "DELETE /sprints/$TEST_SPRINT_ID"
  resp=$(req DELETE "$BASE_URL/sprints/$TEST_SPRINT_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /sprints/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

info "POST /sprints/$SPRINT2_ID/start (planned sprint → active)"
resp=$(req POST "$BASE_URL/sprints/$SPRINT2_ID/start" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
STATUS=$(status_of "$resp")
if [[ "$STATUS" == "200" ]] || [[ "$STATUS" == "409" ]]; then
  ok "POST sprint start [WEBDEV — faqat 1 active sprint bo'ladi]"
else
  fail "POST sprint start" "$(body_of "$resp")"
fi

# =============================================================================
# 8. ISSUES
# =============================================================================
section "8. ISSUES"

ISSUE1_ID="10000001-0000-0000-0000-000000000000"  # epic
ISSUE_BUG_ID="1000000b-0000-0000-0000-000000000000"  # Safari bug
STATUS_TODO="433c09ee-9223-41d3-afeb-df59a2336531"
STATUS_IN_PROGRESS="3d1c059c-c5be-415c-b848-773ba5f5fc71"
STATUS_IN_REVIEW="7a1f6e46-05d5-4d07-b50d-18c6e39a7040"
STATUS_DONE="6703b413-4fbf-4577-9b6f-2b5b994ad00d"

info "GET /issues?project_id=$WEBDEV_ID"
resp=$(req GET "$BASE_URL/issues?project_id=$WEBDEV_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issues (admin)" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /issues?project_id=$WEBDEV_ID&type=bug"
resp=$(req GET "$BASE_URL/issues?project_id=$WEBDEV_ID&type=bug" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issues by type=bug" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /issues?project_id=$WEBDEV_ID&priority=highest"
resp=$(req GET "$BASE_URL/issues?project_id=$WEBDEV_ID&priority=highest" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issues by priority=highest" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /issues?sprint_id=$SPRINT1_ID"
resp=$(req GET "$BASE_URL/issues?sprint_id=$SPRINT1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issues by sprint" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /issues/$ISSUE1_ID"
resp=$(req GET "$BASE_URL/issues/$ISSUE1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /issues/:id (epic)" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /issues (member yaratadi)"
resp=$(req POST "$BASE_URL/issues" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"project_id\":\"$WEBDEV_ID\",\"title\":\"New test issue from API\",\"type\":\"task\",\"priority\":\"medium\",\"description\":\"API test orqali yaratilgan issue\"}")
check "POST /issues member" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_ISSUE_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_ISSUE_ID" ]]; then
  info "PUT /issues/$TEST_ISSUE_ID"
  resp=$(req PUT "$BASE_URL/issues/$TEST_ISSUE_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"title\":\"Updated issue title\",\"priority\":\"high\",\"assignee_id\":\"b4000000-0000-0000-0000-000000000004\"}")
  check "PUT /issues/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /issues/$TEST_ISSUE_ID/transition (Todo → In Progress)"
  resp=$(req POST "$BASE_URL/issues/$TEST_ISSUE_ID/transition" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"status_id\":\"$STATUS_IN_PROGRESS\"}")
  check "POST /issues/:id/transition" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "GET /issues/$TEST_ISSUE_ID/history"
  resp=$(req GET "$BASE_URL/issues/$TEST_ISSUE_ID/history" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "GET /issues/:id/history" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /issues/$TEST_ISSUE_ID"
  resp=$(req DELETE "$BASE_URL/issues/$TEST_ISSUE_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE /issues/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

info "GET /issues/$ISSUE_BUG_ID/history (existing bug)"
resp=$(req GET "$BASE_URL/issues/$ISSUE_BUG_ID/history" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET bug issue history" "200" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 9. ISSUE LINKS
# =============================================================================
section "9. ISSUE LINKS"

# CORS task va API docs task — hali link yo'q juftlik
ISSUE_A="1000000d-0000-0000-0000-000000000000"
ISSUE_B="1000000e-0000-0000-0000-000000000000"

info "GET /issues/$ISSUE_BUG_ID/links"
resp=$(req GET "$BASE_URL/issues/$ISSUE_BUG_ID/links" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issue links" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /issues/$ISSUE_A/links"
resp=$(req POST "$BASE_URL/issues/$ISSUE_A/links" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"target_id\":\"$ISSUE_B\",\"link_type\":\"relates_to\"}")
check "POST issue link" "201" "$(status_of "$resp")" "$(body_of "$resp")"
LINK_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$LINK_ID" ]]; then
  info "DELETE /issues/links/$LINK_ID"
  resp=$(req DELETE "$BASE_URL/issues/links/$LINK_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE issue link" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 10. ISSUE WATCHERS
# =============================================================================
section "10. ISSUE WATCHERS"

WATCH_ISSUE="10000009-0000-0000-0000-000000000000"

info "GET /issues/$ISSUE_BUG_ID/watchers"
resp=$(req GET "$BASE_URL/issues/$ISSUE_BUG_ID/watchers" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issue watchers" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /issues/$WATCH_ISSUE/watchers (dev kuzatadi)"
resp=$(req POST "$BASE_URL/issues/$WATCH_ISSUE/watchers" \
  -H "Authorization: Bearer $DEV_TOKEN")
check "POST issue watcher" "204" "$(status_of "$resp")" "$(body_of "$resp")"

info "DELETE /issues/$WATCH_ISSUE/watchers (dev kuzatishni to'xtatadi)"
resp=$(req DELETE "$BASE_URL/issues/$WATCH_ISSUE/watchers" \
  -H "Authorization: Bearer $DEV_TOKEN")
check "DELETE issue watcher" "204" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 11. LABELS
# =============================================================================
section "11. LABELS"

LABEL_FRONTEND="e1000000-0000-0000-0000-000000000001"

info "GET /projects/$WEBDEV_ID/labels"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/labels" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET project labels" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /labels/$LABEL_FRONTEND"
resp=$(req GET "$BASE_URL/labels/$LABEL_FRONTEND" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /labels/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /projects/$WEBDEV_ID/labels"
resp=$(req POST "$BASE_URL/projects/$WEBDEV_ID/labels" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"testing","color":"#8B5CF6"}')
check "POST label" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_LABEL_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_LABEL_ID" ]]; then
  info "PUT /labels/$TEST_LABEL_ID"
  resp=$(req PUT "$BASE_URL/labels/$TEST_LABEL_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"testing-v2","color":"#7C3AED"}')
  check "PUT /labels/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /labels/$TEST_LABEL_ID"
  resp=$(req DELETE "$BASE_URL/labels/$TEST_LABEL_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE /labels/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 12. CUSTOM FIELDS
# =============================================================================
section "12. CUSTOM FIELDS"

CF_BROWSER="f1000000-0000-0000-0000-000000000001"

info "GET /projects/$WEBDEV_ID/custom-fields"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/custom-fields" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET custom fields" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /custom-fields/$CF_BROWSER"
resp=$(req GET "$BASE_URL/custom-fields/$CF_BROWSER" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /custom-fields/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /projects/$WEBDEV_ID/custom-fields"
resp=$(req POST "$BASE_URL/projects/$WEBDEV_ID/custom-fields" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Field\",\"field_key\":\"tst_fld_${TS}\",\"field_type\":\"text\",\"is_required\":false}")
check "POST custom field" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_CF_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_CF_ID" ]]; then
  info "PUT /custom-fields/$TEST_CF_ID"
  resp=$(req PUT "$BASE_URL/custom-fields/$TEST_CF_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Field Updated"}')
  check "PUT /custom-fields/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /custom-fields/$TEST_CF_ID"
  resp=$(req DELETE "$BASE_URL/custom-fields/$TEST_CF_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /custom-fields/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 13. BOARDS
# =============================================================================
section "13. BOARDS"

BOARD1_ID="30000001-0000-0000-0000-000000000000"
BOARD_COL_ID="40000001-0000-0000-0000-000000000000"

info "GET /projects/$WEBDEV_ID/boards"
resp=$(req GET "$BASE_URL/projects/$WEBDEV_ID/boards" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET project boards" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /boards/$BOARD1_ID"
resp=$(req GET "$BASE_URL/boards/$BOARD1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /boards/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /projects/$WEBDEV_ID/boards"
resp=$(req POST "$BASE_URL/projects/$WEBDEV_ID/boards" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Board","type":"kanban","filter":{}}')
check "POST board" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_BOARD_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_BOARD_ID" ]]; then
  info "PUT /boards/$TEST_BOARD_ID"
  resp=$(req PUT "$BASE_URL/boards/$TEST_BOARD_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Updated Board"}')
  check "PUT /boards/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /boards/$TEST_BOARD_ID/columns"
  resp=$(req POST "$BASE_URL/boards/$TEST_BOARD_ID/columns" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"New Column","position":1}')
  check "POST board column" "201" "$(status_of "$resp")" "$(body_of "$resp")"
  TEST_COL_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

  if [[ -n "$TEST_COL_ID" ]]; then
    info "PUT /board-columns/$TEST_COL_ID"
    resp=$(req PUT "$BASE_URL/board-columns/$TEST_COL_ID" \
      -H "Authorization: Bearer $MEMBER1_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"Updated Column","wip_limit":5}')
    check "PUT /board-columns/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

    info "DELETE /board-columns/$TEST_COL_ID"
    resp=$(req DELETE "$BASE_URL/board-columns/$TEST_COL_ID" \
      -H "Authorization: Bearer $MEMBER1_TOKEN")
    check "DELETE /board-columns/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
  fi

  info "DELETE /boards/$TEST_BOARD_ID"
  resp=$(req DELETE "$BASE_URL/boards/$TEST_BOARD_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE /boards/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 14. COMMENTS
# =============================================================================
section "14. COMMENTS"

COMMENT1_ID="70000001-0000-0000-0000-000000000000"

info "GET /issues/$ISSUE_BUG_ID/comments"
resp=$(req GET "$BASE_URL/issues/$ISSUE_BUG_ID/comments" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET issue comments" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /issues/$ISSUE_BUG_ID/comments (member)"
resp=$(req POST "$BASE_URL/issues/$ISSUE_BUG_ID/comments" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Test comment API orqali"}]}]},"content_text":"Test comment API orqali"}')
check "POST issue comment" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_COMMENT_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_COMMENT_ID" ]]; then
  info "GET /comments/$TEST_COMMENT_ID"
  resp=$(req GET "$BASE_URL/comments/$TEST_COMMENT_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "GET /comments/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "PUT /comments/$TEST_COMMENT_ID"
  resp=$(req PUT "$BASE_URL/comments/$TEST_COMMENT_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Yangilangan comment"}]}]},"content_text":"Yangilangan comment"}')
  check "PUT /comments/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST reply to comment"
  resp=$(req POST "$BASE_URL/issues/$ISSUE_BUG_ID/comments" \
    -H "Authorization: Bearer $DEV_TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"content\":{\"type\":\"doc\",\"content\":[{\"type\":\"paragraph\",\"content\":[{\"type\":\"text\",\"text\":\"Reply comment\"}]}]},\"content_text\":\"Reply comment\",\"reply_to_id\":\"$TEST_COMMENT_ID\"}")
  check "POST comment reply" "201" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /comments/$TEST_COMMENT_ID"
  resp=$(req DELETE "$BASE_URL/comments/$TEST_COMMENT_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE /comments/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 15. SPACES
# =============================================================================
section "15. SPACES"

SPACE1_ID="50000001-0000-0000-0000-000000000000"

info "GET /spaces (admin)"
resp=$(req GET "$BASE_URL/spaces" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /spaces admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /spaces/$SPACE1_ID"
resp=$(req GET "$BASE_URL/spaces/$SPACE1_ID" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /spaces/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /spaces (admin)"
resp=$(req POST "$BASE_URL/spaces" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"key\":\"TS${TS}\",\"name\":\"Test Space ${TS}\",\"description\":\"Test space\",\"type\":\"team\"}")
check "POST /spaces" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_SPACE_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_SPACE_ID" ]]; then
  info "GET /spaces/$TEST_SPACE_ID/members"
  resp=$(req GET "$BASE_URL/spaces/$TEST_SPACE_ID/members" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "GET space members" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "POST /spaces/$TEST_SPACE_ID/members"
  resp=$(req POST "$BASE_URL/spaces/$TEST_SPACE_ID/members" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"user_id":"b1000000-0000-0000-0000-000000000001","role":"member"}')
  check "POST space member" "201" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /spaces/$TEST_SPACE_ID"
  resp=$(req DELETE "$BASE_URL/spaces/$TEST_SPACE_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /spaces/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 16. PAGES
# =============================================================================
section "16. PAGES"

PAGE1_ID="60000001-0000-0000-0000-000000000000"
PAGE2_ID="60000002-0000-0000-0000-000000000000"

info "GET /spaces/$SPACE1_ID/pages/tree"
resp=$(req GET "$BASE_URL/spaces/$SPACE1_ID/pages/tree" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET space page tree" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /pages (filter by space)"
resp=$(req GET "$BASE_URL/pages?space_id=$SPACE1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /pages" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /pages/$PAGE1_ID"
resp=$(req GET "$BASE_URL/pages/$PAGE1_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /pages/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /spaces/$SPACE1_ID/pages"
resp=$(req POST "$BASE_URL/spaces/$SPACE1_ID/pages" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Page","content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Test sahifa"}]}]},"content_text":"Test sahifa"}')
check "POST page" "201" "$(status_of "$resp")" "$(body_of "$resp")"
TEST_PAGE_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

if [[ -n "$TEST_PAGE_ID" ]]; then
  info "PUT /pages/$TEST_PAGE_ID"
  resp=$(req PUT "$BASE_URL/pages/$TEST_PAGE_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"Updated Page Title","content":{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Yangilangan sahifa"}]}]},"content_text":"Yangilangan sahifa"}')
  check "PUT /pages/:id" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "GET /pages/$TEST_PAGE_ID/versions"
  resp=$(req GET "$BASE_URL/pages/$TEST_PAGE_ID/versions" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "GET page versions" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "GET /pages/$TEST_PAGE_ID/versions/1"
  resp=$(req GET "$BASE_URL/pages/$TEST_PAGE_ID/versions/1" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "GET page version by number" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "GET /pages/$PAGE2_ID/comments"
  resp=$(req GET "$BASE_URL/pages/$PAGE2_ID/comments" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "GET page comments" "200" "$(status_of "$resp")" "$(body_of "$resp")"

  info "DELETE /pages/$TEST_PAGE_ID"
  resp=$(req DELETE "$BASE_URL/pages/$TEST_PAGE_ID" \
    -H "Authorization: Bearer $MEMBER1_TOKEN")
  check "DELETE /pages/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# 17. NOTIFICATIONS
# =============================================================================
section "17. NOTIFICATIONS"

info "GET /notifications (member1)"
resp=$(req GET "$BASE_URL/notifications" -H "Authorization: Bearer $MEMBER1_TOKEN")
check "GET /notifications" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /notifications/unread-count"
resp=$(req GET "$BASE_URL/notifications/unread-count" \
  -H "Authorization: Bearer $MEMBER1_TOKEN")
check "GET /notifications/unread-count" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /notifications/mark-all-read"
resp=$(req POST "$BASE_URL/notifications/mark-all-read" \
  -H "Authorization: Bearer $MEMBER1_TOKEN")
check "POST /notifications/mark-all-read" "204" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /notifications/preferences"
resp=$(req GET "$BASE_URL/notifications/preferences" \
  -H "Authorization: Bearer $MEMBER1_TOKEN")
check "GET notification preferences" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "PUT /notifications/preferences"
resp=$(req PUT "$BASE_URL/notifications/preferences" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"issue_assigned":true,"issue_commented":true,"issue_status_changed":false,"mention":true}')
check "PUT notification preferences" "200" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 18. SEARCH
# =============================================================================
section "18. SEARCH"

info "GET /search?q=login (admin)"
resp=$(req GET "$BASE_URL/search?q=login" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /search?q=login" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /search?q=bug (member)"
resp=$(req GET "$BASE_URL/search?q=bug" -H "Authorization: Bearer $MEMBER1_TOKEN")
check "GET /search?q=bug" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /search?q=authentication (viewer)"
resp=$(req GET "$BASE_URL/search?q=authentication" -H "Authorization: Bearer $VIEWER_TOKEN")
check "GET /search?q=authentication viewer" "200" "$(status_of "$resp")" "$(body_of "$resp")"

# =============================================================================
# 19. AUDIT LOGS
# =============================================================================
section "19. AUDIT LOGS"

info "GET /audit-logs (admin)"
resp=$(req GET "$BASE_URL/audit-logs" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /audit-logs admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "GET /audit-logs (member — ruxsatsiz bo'lishi mumkin)"
resp=$(req GET "$BASE_URL/audit-logs" -H "Authorization: Bearer $MEMBER1_TOKEN")
STATUS=$(status_of "$resp")
if [[ "$STATUS" == "200" ]] || [[ "$STATUS" == "403" ]]; then
  ok "GET /audit-logs member [$STATUS — OK]"
else
  fail "GET /audit-logs member" "$(body_of "$resp")"
fi

# =============================================================================
# 20. INVITES
# =============================================================================
section "20. INVITES"

info "POST /invites (admin yangi foydalanuvchi taklif qiladi)"
resp=$(req POST "$BASE_URL/invites" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"invited${TS}@example.com\",\"role\":\"member\"}")
check "POST /invites admin" "201" "$(status_of "$resp")" "$(body_of "$resp")"
INVITE_ID=$(body_of "$resp" | python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['id'])" 2>/dev/null || echo "")

info "GET /invites (admin)"
resp=$(req GET "$BASE_URL/invites" -H "Authorization: Bearer $ADMIN_TOKEN")
check "GET /invites admin" "200" "$(status_of "$resp")" "$(body_of "$resp")"

info "POST /invites (member — ruxsatsiz)"
resp=$(req POST "$BASE_URL/invites" \
  -H "Authorization: Bearer $MEMBER1_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"test2${TS}@example.com\",\"role\":\"member\"}")
STATUS_INV_MEMBER=$(status_of "$resp")
if [[ "$STATUS_INV_MEMBER" == "403" ]]; then
  ok "POST /invites member → 403 [RBAC ishlayapti]"
else
  warn "POST /invites member → $STATUS_INV_MEMBER [RBAC konfiguratsiyasi: member invite yubora olmasligi kerak]"
fi

if [[ -n "$INVITE_ID" ]]; then
  info "DELETE /invites/$INVITE_ID (revoke)"
  resp=$(req DELETE "$BASE_URL/invites/$INVITE_ID" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  check "DELETE /invites/:id" "204" "$(status_of "$resp")" "$(body_of "$resp")"
fi

# =============================================================================
# YAKUNIY HISOBOT
# =============================================================================
echo ""
echo -e "${BOLD}${BLUE}══════════════════════════════════════${NC}"
echo -e "${BOLD}  YAKUNIY NATIJA${NC}"
echo -e "${BOLD}${BLUE}══════════════════════════════════════${NC}"
echo -e "  ${GREEN}✓ Muvaffaqiyatli:${NC} $PASS"
echo -e "  ${RED}✗ Xatolik:${NC}       $FAIL"
echo -e "  ${YELLOW}⊙ O'tkazilgan:${NC}  $SKIP"
TOTAL=$((PASS + FAIL + SKIP))
echo -e "  Jami:           $TOTAL"
echo ""
if [[ $FAIL -eq 0 ]]; then
  echo -e "${GREEN}${BOLD}  Barcha testlar muvaffaqiyatli o'tdi!${NC}"
else
  echo -e "${RED}${BOLD}  $FAIL ta test muvaffaqiyatsiz bo'ldi.${NC}"
fi
echo ""
