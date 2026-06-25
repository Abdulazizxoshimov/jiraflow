// adapters.jsx — Normalize backend responses to the shape views expect.
// Backend uses snake_case + UUID IDs + different field names.
// Views expect: name, pri, points, status-as-string, etc.

// ─── Helpers ──────────────────────────────────────────────────────────────
export function userInitials(name) {
  if (!name) return "?";
  return name.split(" ").filter(Boolean).slice(0, 2).map((w) => w[0].toUpperCase()).join("");
}

export function fmtDate(iso) {
  if (!iso) return null;
  const d = new Date(iso);
  return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
}

// Priority: backend → UI
const PRI_MAP = {
  highest: "Critical",
  high:    "High",
  medium:  "Medium",
  low:     "Low",
  lowest:  "Low",
};
// Priority: UI → backend
const PRI_REVERSE = {
  Critical: "highest",
  High:     "high",
  Medium:   "medium",
  Low:      "low",
};

function adaptType(t) {
  if (!t) return "Task";
  return t.charAt(0).toUpperCase() + t.slice(1);
}

// ─── User ─────────────────────────────────────────────────────────────────
// Backend: { id, email, full_name, avatar_url, color, role, is_active, ... }
// UI:      { id, name, initials, color, email, role, status, avatar, tg }
export function adaptUser(u) {
  if (!u) return null;
  return {
    id:       u.id,
    name:     u.full_name,
    initials: userInitials(u.full_name),
    color:    u.color || "#6366F1",
    email:    u.email,
    role:     u.role ? u.role.toLowerCase() : "member",
    status:   u.is_active ? "Active" : "Inactive",
    avatar:   u.avatar_url || null,
    tg:       null,
    joined:   fmtDate(u.created_at),
    _raw:     u,
  };
}

// ─── Issue ────────────────────────────────────────────────────────────────
export function adaptIssue(issue, projectKey) {
  if (!issue) return null;
  const key = projectKey && issue.issue_number
    ? projectKey + "-" + issue.issue_number
    : issue.id;
  return {
    id:         key,
    _id:        issue.id,
    project_id: issue.project_id,
    title:      issue.title,
    type:       adaptType(issue.type),
    pri:        PRI_MAP[issue.priority] || "Medium",
    points:     issue.story_points || 0,
    status:     issue.status ? issue.status.name : null,
    status_id:  issue.status_id,
    assignee:   issue.assignee_id,
    reporter:   issue.reporter_id,
    labels:     (issue.labels || []).map((l) => l.name || l),
    sprint:     issue.sprint_id || null,
    due:        fmtDate(issue.due_date),
    comments:   0,
    sub:        0,
    assigneeUser: issue.assignee ? adaptUser(issue.assignee) : null,
    reporterUser: issue.reporter ? adaptUser(issue.reporter) : null,
    _raw: issue,
  };
}

// ─── Project ──────────────────────────────────────────────────────────────
const PROJECT_COLORS = ["#6366F1","#06B6D4","#10B981","#F59E0B","#EF4444","#8B5CF6","#EC4899","#14B8A6","#F97316","#3B82F6"];
const PROJECT_ICONS  = ["server","rocket","chart","shield","database","code","kanban","briefcase"];

function projectColor(id) {
  let h = 0;
  for (let i = 0; i < id.length; i++) h = (h * 31 + id.charCodeAt(i)) >>> 0;
  return PROJECT_COLORS[h % PROJECT_COLORS.length];
}
function projectIcon(key) {
  if (!key) return "briefcase";
  let h = 0;
  for (let i = 0; i < key.length; i++) h = (h * 31 + key.charCodeAt(i)) >>> 0;
  return PROJECT_ICONS[h % PROJECT_ICONS.length];
}

export function adaptProject(p) {
  if (!p) return null;
  return {
    id:         p.id,
    key:        p.key,
    name:       p.name,
    desc:       p.description || "",
    color:      projectColor(p.id),
    icon:       projectIcon(p.key),
    lead:       p.lead_id,
    leadUser:   p.lead ? adaptUser(p.lead) : null,
    members:    p.member_count || 0,
    openIssues: p.open_issues  || 0,
    updated:    fmtDate(p.updated_at),
    isArchived: p.is_archived,
    _raw:       p,
  };
}

// ─── WorkflowStatus → COLUMNS shape ──────────────────────────────────────
export function adaptStatus(s) {
  const TONE = { todo: "muted", in_progress: "info", done: "success" };
  return {
    id:    s.name,
    label: s.name,
    tone:  TONE[s.category] || "muted",
    _id:   s.id,
  };
}

// ─── Sprint ───────────────────────────────────────────────────────────────
export function adaptSprint(s) {
  if (!s) return null;
  return {
    id:        s.id,
    name:      s.name,
    status:    s.status,
    startDate: s.start_date,
    endDate:   s.end_date,
    goal:      s.goal,
    _raw:      s,
  };
}

// ─── Comment ─────────────────────────────────────────────────────────────
export function adaptComment(c) {
  if (!c) return null;
  return {
    id:     c.id,
    who:    c.author_id,
    author: c.author ? adaptUser(c.author) : null,
    body:   c.body,
    at:     fmtDate(c.created_at),
    edited: c.updated_at !== c.created_at,
    _raw:   c,
  };
}

// ─── Utility: adapt list ─────────────────────────────────────────────────
export function adaptList(items, fn, ...args) {
  return (items || []).map((i) => fn(i, ...args));
}

// ─── Priority converters (for create/update) ─────────────────────────────
export function toBackendPriority(uiPriority) {
  return PRI_REVERSE[uiPriority] || "medium";
}
export function toBackendType(uiType) {
  return (uiType || "task").toLowerCase();
}
