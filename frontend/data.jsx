// data.jsx — Forge dummy data. DevOps/infrastructure flavor.

const COLORS = ["#6366F1","#06B6D4","#10B981","#F59E0B","#EF4444","#8B5CF6","#EC4899","#14B8A6","#F97316","#3B82F6"];

const PEOPLE = [
  { id: "u1",  name: "Maya Chen",        initials: "MC", color: "#6366F1", email: "maya@forge.dev",       role: "Admin",     tg: "@mayaops", tgId: "847291043", status: "Active",  joined: "Mar 2024" },
  { id: "u2",  name: "Diego Alvarez",    initials: "DA", color: "#06B6D4", email: "diego@forge.dev",      role: "Manager",   tg: "@diego_a", tgId: "192847301", status: "Active",  joined: "Jan 2024" },
  { id: "u3",  name: "Priya Raman",      initials: "PR", color: "#10B981", email: "priya@forge.dev",      role: "Developer", tg: "@priyacodes", tgId: "405829471", status: "Active",  joined: "Apr 2024" },
  { id: "u4",  name: "Jonas Weber",      initials: "JW", color: "#F59E0B", email: "jonas@forge.dev",      role: "Developer", tg: null, tgId: null, status: "Active", joined: "Feb 2024" },
  { id: "u5",  name: "Aisha Okonkwo",    initials: "AO", color: "#EC4899", email: "aisha@forge.dev",      role: "Developer", tg: "@aisha_o", tgId: "938174625", status: "Active",  joined: "May 2024" },
  { id: "u6",  name: "Tomás Silva",      initials: "TS", color: "#8B5CF6", email: "tomas@forge.dev",      role: "Developer", tg: "@tomas_s", tgId: "729384651", status: "Active",  joined: "Jun 2024" },
  { id: "u7",  name: "Hana Suzuki",      initials: "HS", color: "#14B8A6", email: "hana@forge.dev",       role: "Developer", tg: null, tgId: null, status: "Pending", joined: "Nov 2024" },
  { id: "u8",  name: "Leo Marchetti",    initials: "LM", color: "#F97316", email: "leo@forge.dev",        role: "Viewer",    tg: "@leom",    tgId: "564738291", status: "Active",  joined: "Sep 2024" },
  { id: "u9",  name: "Yara Haddad",      initials: "YH", color: "#3B82F6", email: "yara@forge.dev",       role: "Developer", tg: "@yara_h",  tgId: "203948572", status: "Active",  joined: "Aug 2024" },
  { id: "u10", name: "Connor O'Brien",   initials: "CO", color: "#EF4444", email: "connor@forge.dev",     role: "Developer", tg: null, tgId: null, status: "Inactive", joined: "Jul 2023" },
];

const ME = PEOPLE[0]; // Maya Chen, Admin

const PROJECTS = [
  { id: "infra",   key: "INFRA", name: "Core Infrastructure",  color: "#6366F1", icon: "server",   desc: "Kubernetes clusters, networking, baseline infra.", members: 8, openIssues: 47, lead: "u1", updated: "2h ago" },
  { id: "deploy",  key: "DEP",   name: "Deploy Pipeline",      color: "#10B981", icon: "rocket",   desc: "CI/CD, release automation, blue-green deploys.",   members: 6, openIssues: 23, lead: "u2", updated: "1d ago" },
  { id: "obs",     key: "OBS",   name: "Observability",        color: "#F59E0B", icon: "chart",    desc: "Metrics, logs, traces. Grafana, Loki, Tempo.",     members: 5, openIssues: 31, lead: "u3", updated: "3h ago" },
  { id: "sec",     key: "SEC",   name: "Security & Compliance",color: "#EF4444", icon: "shield",   desc: "SOC2, IAM, secrets rotation, audit.",              members: 4, openIssues: 18, lead: "u1", updated: "5h ago" },
  { id: "data",    key: "DATA",  name: "Data Platform",        color: "#8B5CF6", icon: "database", desc: "Warehouse, ETL, lakehouse, streaming.",            members: 7, openIssues: 29, lead: "u5", updated: "yesterday" },
  { id: "plat",    key: "PLAT",  name: "Developer Platform",   color: "#06B6D4", icon: "code",     desc: "Internal tooling, scaffolding, golden paths.",     members: 5, openIssues: 14, lead: "u2", updated: "2d ago" },
];

const ACTIVE_PROJECT_ID = "infra";

const ISSUES = [
  // Backlog
  { id: "INFRA-241", title: "Reduce coldstart latency on edge workers below 80ms p99", type: "Story",    pri: "High",     status: "Backlog",    assignee: "u3", reporter: "u1", points: 5, due: "Dec 14", labels: ["edge","performance"], sprint: null,    comments: 4, sub: 2 },
  { id: "INFRA-242", title: "Document multi-region failover runbook for control plane", type: "Task",     pri: "Medium",   status: "Backlog",    assignee: "u4", reporter: "u2", points: 3, due: "Dec 22", labels: ["docs","oncall"], sprint: null,    comments: 1, sub: 0 },
  { id: "INFRA-243", title: "Evaluate Cilium vs Calico for cluster CNI replacement",   type: "Epic",     pri: "Medium",   status: "Backlog",    assignee: "u2", reporter: "u1", points: 13, due: null, labels: ["networking","spike"], sprint: null, comments: 9, sub: 4 },
  { id: "INFRA-244", title: "API gateway: rate limit headers missing on 429s",         type: "Bug",      pri: "Low",      status: "Backlog",    assignee: "u5", reporter: "u5", points: 1, due: null, labels: ["api"], sprint: null, comments: 0, sub: 0 },

  // Todo
  { id: "INFRA-230", title: "Migrate eu-west-2 NAT gateways to AWS PrivateLink",       type: "Story",    pri: "High",     status: "Todo",       assignee: "u4", reporter: "u1", points: 8, due: "Dec 12", labels: ["aws","networking"], sprint: "Sprint 24", comments: 2, sub: 3 },
  { id: "INFRA-231", title: "Bump Terraform to 1.10 across all modules",               type: "Task",     pri: "Medium",   status: "Todo",       assignee: "u3", reporter: "u2", points: 3, due: "Dec 10", labels: ["terraform"], sprint: "Sprint 24", comments: 0, sub: 0 },
  { id: "INFRA-232", title: "Investigate Loki ingester OOM in observability namespace",type: "Bug",      pri: "Critical", status: "Todo",       assignee: "u5", reporter: "u1", points: 5, due: "Dec 9",  labels: ["loki","oncall"], sprint: "Sprint 24", comments: 7, sub: 1 },
  { id: "INFRA-233", title: "Wire pod-identity to all data-platform services",         type: "Story",    pri: "Medium",   status: "Todo",       assignee: "u6", reporter: "u1", points: 5, due: null, labels: ["iam"], sprint: "Sprint 24", comments: 2, sub: 2 },

  // In Progress
  { id: "INFRA-220", title: "Provision shared etcd cluster for staging fleet",         type: "Story",    pri: "High",     status: "In Progress",assignee: "u3", reporter: "u1", points: 5, due: "Dec 8", labels: ["etcd","staging"], sprint: "Sprint 24", comments: 5, sub: 2 },
  { id: "INFRA-221", title: "Fix flaky e2e test: ClusterAutoscaler scale-from-zero",   type: "Bug",      pri: "High",     status: "In Progress",assignee: "u6", reporter: "u3", points: 3, due: "Dec 7", labels: ["ci","flaky"], sprint: "Sprint 24", comments: 3, sub: 0 },
  { id: "INFRA-222", title: "Roll out cgroup v2 to nodepools us-east-1 → us-west-2",   type: "Task",     pri: "Medium",   status: "In Progress",assignee: "u4", reporter: "u1", points: 8, due: "Dec 15", labels: ["kernel","rolllout"], sprint: "Sprint 24", comments: 1, sub: 4 },

  // In Review
  { id: "INFRA-210", title: "Add SLO burn-rate alerts for control-plane API",          type: "Story",    pri: "High",     status: "In Review",  assignee: "u5", reporter: "u2", points: 3, due: "Dec 6", labels: ["slo","alerts"], sprint: "Sprint 24", comments: 6, sub: 1 },
  { id: "INFRA-211", title: "Helm chart: replace deprecated PSP with PodSecurity",     type: "Task",     pri: "Medium",   status: "In Review",  assignee: "u9", reporter: "u1", points: 2, due: "Dec 5", labels: ["helm","security"], sprint: "Sprint 24", comments: 2, sub: 0 },

  // Done
  { id: "INFRA-200", title: "Cut over staging DNS to Route53 hosted zone",             type: "Story",    pri: "Medium",   status: "Done",       assignee: "u3", reporter: "u2", points: 3, due: "Dec 1", labels: ["dns"], sprint: "Sprint 24", comments: 4, sub: 0 },
  { id: "INFRA-201", title: "Patch CVE-2024-9341 in base image",                       type: "Bug",      pri: "Critical", status: "Done",       assignee: "u1", reporter: "u1", points: 2, due: "Nov 29", labels: ["cve","security"], sprint: "Sprint 24", comments: 8, sub: 0 },
  { id: "INFRA-202", title: "Onboard data-platform to centralized secret store",       type: "Story",    pri: "High",     status: "Done",       assignee: "u5", reporter: "u1", points: 5, due: "Nov 27", labels: ["secrets"], sprint: "Sprint 24", comments: 11, sub: 3 },
];

const COLUMNS = [
  { id: "Backlog",     label: "Backlog",     tone: "muted" },
  { id: "Todo",        label: "Todo",        tone: "muted" },
  { id: "In Progress", label: "In progress", tone: "info" },
  { id: "In Review",   label: "In review",   tone: "purple" },
  { id: "Done",        label: "Done",        tone: "success" },
];

const TYPE_META = {
  Bug:   { tone: "danger",  icon: "bug",   color: "#EF4444" },
  Task:  { tone: "info",    icon: "task",  color: "#3B82F6" },
  Story: { tone: "purple",  icon: "story", color: "#8B5CF6" },
  Epic:  { tone: "orange",  icon: "epic",  color: "#F97316" },
};

const PRIORITY_META = {
  Critical: { tone: "danger",  icon: "prHigh", color: "#DC2626" },
  High:     { tone: "warning", icon: "prHigh", color: "#EA580C" },
  Medium:   { tone: "info",    icon: "prMed",  color: "#3B82F6" },
  Low:      { tone: "muted",   icon: "prLow",  color: "#64748B" },
};

const STATUS_META = {
  "Backlog":     { tone: "muted"   },
  "Todo":        { tone: "muted"   },
  "In Progress": { tone: "info"    },
  "In Review":   { tone: "purple"  },
  "Done":        { tone: "success" },
  "Blocked":     { tone: "danger"  },
};

const ROLE_META = {
  Admin:     { tone: "danger"  },
  Manager:   { tone: "purple"  },
  Developer: { tone: "info"    },
  Viewer:    { tone: "muted"   },
};

const ACTIVITY = [
  { id: 1, who: "u3", verb: "moved",     target: "INFRA-220", to: "In Progress", time: "12m ago" },
  { id: 2, who: "u5", verb: "commented", target: "INFRA-232", time: "28m ago", body: "Confirmed it's the new chunk-encoding path. Rolling back ingester to v3.1.2." },
  { id: 3, who: "u1", verb: "merged PR for", target: "INFRA-201", time: "1h ago" },
  { id: 4, who: "u4", verb: "created",   target: "INFRA-242", time: "3h ago" },
  { id: 5, who: "u2", verb: "started",   target: "Sprint 24", time: "yesterday" },
  { id: 6, who: "u6", verb: "assigned",  target: "INFRA-221", to: "Tomás Silva", time: "yesterday" },
  { id: 7, who: "u3", verb: "resolved",  target: "INFRA-200", time: "2d ago" },
];

const NOTIF_EVENTS = [
  { key: "assigned",  label: "Issue assigned to me",                desc: "When someone assigns an issue to you.",                  default: true },
  { key: "status",    label: "Issue status changed",                desc: "Status moves on issues you watch.",                       default: true },
  { key: "comment",   label: "Comment on my issue",                 desc: "New comments on issues you reported or are assigned.",    default: true },
  { key: "mention",   label: "Mentioned in comment",                desc: "Someone @mentions you anywhere.",                         default: true },
  { key: "sprint",    label: "Sprint started or completed",         desc: "Sprint lifecycle events for projects you belong to.",     default: true },
  { key: "newmember", label: "New member joined project",           desc: "Get a heads-up when teammates join.",                     default: false },
  { key: "due",       label: "Due date reminder (1 day before)",    desc: "Nightly reminder for issues due tomorrow.",               default: true },
  { key: "digest",    label: "Daily morning digest",                desc: "Personalized summary delivered at 9:00 local time.",      default: false },
];

const WIKI_TREE = [
  { id: "root",       title: "Engineering",          children: [
    { id: "onboarding", title: "Onboarding",          children: [
      { id: "tools",       title: "Day 1 tooling" },
      { id: "git",         title: "Git conventions" },
      { id: "envs",        title: "Local environments" },
    ]},
    { id: "runbooks",   title: "Runbooks",            children: [
      { id: "rb-loki",     title: "Loki ingester OOM" },
      { id: "rb-failover", title: "Multi-region failover" },
      { id: "rb-rotate",   title: "Secrets rotation" },
    ]},
    { id: "architecture", title: "Architecture",       children: [
      { id: "control",     title: "Control plane" },
      { id: "edge",        title: "Edge fleet" },
      { id: "data",        title: "Data platform" },
    ]},
    { id: "postmortems",title: "Postmortems",          children: [
      { id: "pm-1129",     title: "Nov 29 — CVE patch rollout" },
      { id: "pm-1108",     title: "Nov 8 — etcd leader churn" },
    ]},
  ]},
];

const COMMITS = [
  { sha: "9a3f12c", msg: "fix(loki): bound chunk encoder buffer", who: "u5", when: "12m ago" },
  { sha: "b4c0028", msg: "chore(tf): bump aws provider to 5.74",  who: "u3", when: "1h ago" },
  { sha: "e7d2841", msg: "feat(api): rate-limit headers on 429",  who: "u5", when: "yesterday" },
];

window.FORGE_DATA = { COLORS, PEOPLE, ME, PROJECTS, ACTIVE_PROJECT_ID, ISSUES, COLUMNS, TYPE_META, PRIORITY_META, STATUS_META, ROLE_META, ACTIVITY, NOTIF_EVENTS, WIKI_TREE, COMMITS };
