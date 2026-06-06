// api.jsx — Forge mock backend.
// Installs a fetch() interceptor for /api/v1/* so the prototype runs without a
// live Go server. Feature code calls the real REST endpoints via the api()
// helper exactly as it would in production; only this file knows it's mocked.
// Remove this script + load order entry to run against the real backend.

(function () {
  // Seed an auth token so Authorization: Bearer <token> always has a value.
  if (!localStorage.getItem("token")) localStorage.setItem("token", "forge_demo_jwt_token");

  const P = window.FORGE_DATA;
  const ME = P.ME;
  const now = Date.now();
  const ago = (m) => new Date(now - m * 60000).toISOString();
  let seq = 1000;
  const nid = (pfx) => (pfx || "id") + "_" + (++seq);
  const clone = (x) => JSON.parse(JSON.stringify(x));

  // ─── In-memory store ──────────────────────────────────────────────────
  const db = {
    links: {
      "INFRA-232": [
        { id: nid("lk"), link_type: "is blocked by", issue_id: "INFRA-220" },
        { id: nid("lk"), link_type: "relates to", issue_id: "INFRA-210" },
        { id: nid("lk"), link_type: "blocks", issue_id: "INFRA-244" },
      ],
    },
    watchers: { "INFRA-232": ["u1", "u3", "u5", "u2"] },
    votes: { "INFRA-232": { voters: ["u3", "u5", "u9"] } },
    assignees: { "INFRA-232": ["u5", "u3"] },
    worklogs: {
      "INFRA-232": [
        { id: nid("wl"), user_id: "u5", time_spent: "2h 30m", minutes: 150, description: "Reproduced OOM on staging, captured heap profile.", started_at: ago(60 * 26) },
        { id: nid("wl"), user_id: "u3", time_spent: "1h 15m", minutes: 75, description: "Bisected the regression to the chunk-encoder buffer pool.", started_at: ago(60 * 8) },
      ],
    },
    estimates: { "INFRA-232": { original: 480 } }, // minutes (8h)
    history: {
      "INFRA-232": [
        { id: nid("h"), user_id: "u1", field: "Issue", from: null, to: "created", at: ago(60 * 24 * 5) },
        { id: nid("h"), user_id: "u1", field: "Label", from: null, to: "loki", at: ago(60 * 24 * 5) },
        { id: nid("h"), user_id: "u2", field: "Sprint", from: "Backlog", to: "Sprint 24", at: ago(60 * 24 * 4) },
        { id: nid("h"), user_id: "u5", field: "Status", from: "Backlog", to: "Todo", at: ago(60 * 24 * 2) },
        { id: nid("h"), user_id: "u3", field: "Priority", from: "Medium", to: "Critical", at: ago(60 * 22) },
        { id: nid("h"), user_id: "u5", field: "Status", from: "Todo", to: "In Progress", at: ago(60 * 3) },
      ],
    },

    workflows: null, // built below
    components: [
      { id: nid("cp"), name: "Control plane", description: "API server, scheduler, etcd.", lead_id: "u1", issue_count: 14 },
      { id: nid("cp"), name: "Networking", description: "CNI, ingress, service mesh.", lead_id: "u2", issue_count: 9 },
      { id: nid("cp"), name: "Edge fleet", description: "Edge workers and PoP rollout.", lead_id: "u3", issue_count: 6 },
    ],
    versions: [
      { id: nid("vr"), name: "2024.11 — Hardening", start_date: "2024-11-01", release_date: "2024-11-29", status: "released", description: "CVE patches, PSP migration." },
      { id: nid("vr"), name: "2024.12 — Edge & Reliability", start_date: "2024-12-02", release_date: "2024-12-15", status: "unreleased", description: "Coldstart latency, failover runbooks." },
      { id: nid("vr"), name: "2025.01 — Observability", start_date: "2025-01-06", release_date: "2025-01-24", status: "unreleased", description: "Loki, Tempo, burn-rate alerts." },
    ],
    customFields: [
      { id: nid("cf"), name: "Severity", type: "select", description: "Incident severity classification.", required: true, options: ["SEV1", "SEV2", "SEV3"], position: 0 },
      { id: nid("cf"), name: "Root cause", type: "text", description: "Free-text post-incident root cause.", required: false, options: [], position: 1 },
      { id: nid("cf"), name: "Customer impact", type: "number", description: "Estimated affected customers.", required: false, options: [], position: 2 },
      { id: nid("cf"), name: "Runbook URL", type: "url", description: "Link to the relevant runbook.", required: false, options: [], position: 3 },
    ],
    labels: [
      { id: nid("lb"), name: "oncall", color: "#EF4444", issue_count: 18 },
      { id: nid("lb"), name: "networking", color: "#6366F1", issue_count: 12 },
      { id: nid("lb"), name: "performance", color: "#F59E0B", issue_count: 9 },
      { id: nid("lb"), name: "security", color: "#10B981", issue_count: 14 },
      { id: nid("lb"), name: "terraform", color: "#8B5CF6", issue_count: 7 },
      { id: nid("lb"), name: "docs", color: "#06B6D4", issue_count: 5 },
      { id: nid("lb"), name: "flaky", color: "#EC4899", issue_count: 4 },
      { id: nid("lb"), name: "spike", color: "#F97316", issue_count: 3 },
    ],
    webhooks: [
      { id: nid("wh"), url: "https://hooks.forge.dev/ci/deploy-pipeline", secret: "whsec_••••a91f", events: ["issue.transition", "sprint.completed"], is_active: true, last_status: "success", last_code: 200 },
      { id: nid("wh"), url: "https://api.pagerduty.com/integration/forge/enqueue", secret: "", events: ["issue.created", "issue.assigned"], is_active: true, last_status: "success", last_code: 202 },
      { id: nid("wh"), url: "https://staging.internal/forge-webhook", secret: "whsec_••••6b2c", events: ["page.updated"], is_active: false, last_status: "failed", last_code: 503 },
    ],
    apiKeys: [
      { id: nid("ak"), name: "CI pipeline (GitHub Actions)", prefix: "jfk_7Hd2", created_at: ago(60 * 24 * 40), last_used_at: ago(35) },
      { id: nid("ak"), name: "Terraform provider", prefix: "jfk_q9Lm", created_at: ago(60 * 24 * 120), last_used_at: ago(60 * 24 * 2) },
      { id: nid("ak"), name: "Grafana alerting", prefix: "jfk_b3Xs", created_at: ago(60 * 24 * 12), last_used_at: null },
    ],
    savedFilters: [
      { id: nid("sf"), name: "My critical bugs", filter_json: { type: "Bug", priority: "Critical", assignee: ME.id } },
      { id: nid("sf"), name: "Sprint 24 — unassigned", filter_json: { sprint: "Sprint 24", assignee: "none" } },
    ],
    boardSwimlane: "none",
    notifPrefs: { email_assigned: true, email_mentioned: true, email_commented: true, email_status: false, email_watcher: true },
    importJobs: {},
    sprintCapacity: {
      s24: P.PEOPLE.slice(0, 7).map((u, i) => ({ user_id: u.id, available_hours: [40, 40, 32, 40, 24, 40, 16][i], logged_hours: [31, 28, 30, 18, 26, 12, 4][i], story_points: [13, 11, 14, 9, 8, 5, 2][i] })),
    },
    auditLogs: null, // built below
    spaces: null, pages: null, pageTree: null, blogPosts: null, inlineComments: {}, reactions: {},
    activity: null,
  };

  // Workflows ------------------------------------------------------------
  (function buildWorkflows() {
    const wf = (name, statuses, transitions) => ({ id: nid("wf"), name, statuses, transitions });
    const st = (name, category, color) => ({ id: nid("st"), name, category, color });
    const mk = [];
    const def = [
      st("Backlog", "todo", "#94A3B8"),
      st("Todo", "todo", "#64748B"),
      st("In Progress", "in_progress", "#3B82F6"),
      st("In Review", "in_progress", "#8B5CF6"),
      st("Done", "done", "#10B981"),
    ];
    const trs = [];
    for (let i = 0; i < def.length - 1; i++) trs.push({ id: nid("tr"), from_id: def[i].id, to_id: def[i + 1].id });
    trs.push({ id: nid("tr"), from_id: def[2].id, to_id: def[1].id });
    mk.push(wf("Software development", def, trs));
    const bug = [st("Open", "todo", "#EF4444"), st("Triaged", "todo", "#F59E0B"), st("Fixing", "in_progress", "#3B82F6"), st("Verifying", "in_progress", "#8B5CF6"), st("Closed", "done", "#10B981")];
    const btr = [];
    for (let i = 0; i < bug.length - 1; i++) btr.push({ id: nid("tr"), from_id: bug[i].id, to_id: bug[i + 1].id });
    mk.push(wf("Bug triage", bug, btr));
    db.workflows = mk;
  })();

  // Audit logs -----------------------------------------------------------
  (function buildAudit() {
    const actions = [
      ["members.invite", "create", "member", "aisha@forge.dev"],
      ["sprint.complete", "update", "sprint", "Sprint 23"],
      ["integrations.telegram.connect", "auth", "integration", "@forge_team_bot"],
      ["members.role.change", "update", "member", "Hana Suzuki → Developer"],
      ["project.settings.update", "update", "project", "Core Infrastructure"],
      ["issue.delete", "delete", "issue", "INFRA-188"],
      ["apikey.create", "create", "api_key", "CI pipeline"],
      ["webhook.delete", "delete", "webhook", "staging.internal"],
      ["auth.login", "auth", "session", "maya@forge.dev"],
      ["version.release", "update", "version", "2024.11 — Hardening"],
      ["label.create", "create", "label", "performance"],
      ["page.delete", "delete", "page", "Old onboarding"],
    ];
    const logs = [];
    for (let i = 0; i < 124; i++) {
      const a = actions[i % actions.length];
      const u = P.PEOPLE[i % P.PEOPLE.length];
      logs.push({ id: nid("al"), actor_id: u.id, action: a[0], category: a[1], resource_type: a[2], resource_id: a[3], meta: { ip: "10.0.4." + (10 + (i % 40)) }, created_at: ago(i * 47 + 5) });
    }
    db.auditLogs = logs;
  })();

  // Activity -------------------------------------------------------------
  (function buildActivity() {
    const verbs = [
      ["created", "issue", "INFRA-242"], ["moved", "issue", "INFRA-220"], ["commented", "issue", "INFRA-232"],
      ["assigned", "issue", "INFRA-221"], ["closed", "issue", "INFRA-201"], ["started", "sprint", "Sprint 24"],
      ["updated", "page", "Loki ingester OOM"], ["created", "issue", "INFRA-244"], ["resolved", "issue", "INFRA-200"],
      ["logged time on", "issue", "INFRA-232"], ["released", "version", "2024.11"], ["watched", "issue", "INFRA-210"],
    ];
    const meta = { "INFRA-220": "Done", "INFRA-221": "Tomás Silva" };
    const out = [];
    for (let i = 0; i < 24; i++) {
      const v = verbs[i % verbs.length];
      const u = P.PEOPLE[(i * 3) % P.PEOPLE.length];
      out.push({ id: nid("act"), actor_id: u.id, action: v[0], entity_type: v[1], entity_id: v[2], meta: { to: meta[v[2]] }, created_at: ago(i * 23 + 2) });
    }
    db.activity = out;
  })();

  // Wiki: spaces, pages, tree, blog --------------------------------------
  (function buildWiki() {
    db.spaces = [
      { id: "sp_eng", name: "Engineering", key: "ENG", type: "team", description: "Runbooks, architecture, postmortems.", archived: false, page_count: 14 },
      { id: "sp_inf", name: "Infrastructure", key: "INFRA", type: "project", description: "Core infra design docs and SOPs.", archived: false, page_count: 9 },
      { id: "sp_me", name: "Maya's notes", key: "MAYA", type: "personal", description: "Personal scratchpad.", archived: false, page_count: 3 },
    ];

    const loremBody = `<p>This runbook covers the symptoms, diagnostics, and mitigation steps for Loki ingester pods getting OOMKilled under steady-state log volume in production.</p>
<blockquote><strong>TL;DR</strong> — If ingester memory climbs above 4GB and stays there, scale ingesters horizontally first and rotate the chunk store afterwards. Rollback to v3.1.2 is the nuclear option.</blockquote>
<h2>Symptoms</h2>
<ul><li><code>kubectl get pods -n observability</code> shows <code>OOMKilled</code> in restart history</li><li>Log ingestion lag spikes to 8+ minutes (alert: <code>LokiIngestionLagHigh</code>)</li><li>Grafana panel "Loki — ingester memory" ramps past 4GB</li></ul>
<h2>Mitigation</h2>
<p>Scale ingesters from 6 to 9 replicas. This is the safest first move and rarely makes things worse. If horizontal scale alone does not hold, reduce the chunk encoder buffer pool from 1024 to 256.</p>`;

    const pages = {};
    const mkPage = (id, space_id, title, parent_id, body) => {
      pages[id] = {
        id, space_id, title, parent_id, body: body || `<p>${title} — content goes here. Edit this page to add detail.</p>`,
        author_id: "u3", updated_at: ago(28), versions: [], reactions: {},
      };
      // a couple of versions per page
      const base = pages[id];
      base.versions = [
        { version: 1, author_id: "u3", at: ago(60 * 24 * 9), size: (base.body.length * 0.6) | 0, body: "<p>Initial draft.</p>" },
        { version: 2, author_id: "u5", at: ago(60 * 24 * 3), size: (base.body.length * 0.85) | 0, body: base.body.replace("9 replicas", "8 replicas") },
        { version: 3, author_id: "u3", at: ago(28), size: base.body.length, body: base.body },
      ];
      return pages[id];
    };

    mkPage("pg_onb", "sp_eng", "Onboarding", null);
    mkPage("pg_tools", "sp_eng", "Day 1 tooling", "pg_onb");
    mkPage("pg_git", "sp_eng", "Git conventions", "pg_onb");
    mkPage("pg_rb", "sp_eng", "Runbooks", null);
    mkPage("pg_loki", "sp_eng", "Loki ingester OOM", "pg_rb", loremBody);
    mkPage("pg_failover", "sp_eng", "Multi-region failover", "pg_rb");
    mkPage("pg_rotate", "sp_eng", "Secrets rotation", "pg_rb");
    mkPage("pg_arch", "sp_eng", "Architecture", null);
    mkPage("pg_control", "sp_eng", "Control plane", "pg_arch");
    mkPage("pg_edge", "sp_eng", "Edge fleet", "pg_arch");
    // infra space
    mkPage("pg_sop", "sp_inf", "Standard operating procedures", null);
    mkPage("pg_net", "sp_inf", "Network topology", null);
    // personal
    mkPage("pg_scratch", "sp_me", "Scratchpad", null);

    db.pages = pages;

    db.inlineComments = {
      pg_loki: [
        { id: nid("ic"), page_id: "pg_loki", user_id: "u2", text: "Should we mention cgroup v2 memory limits here too?", anchor_text: "OOMKilled", resolved: false, created_at: ago(120) },
        { id: nid("ic"), page_id: "pg_loki", user_id: "u1", text: "Confirmed — 256 is the right buffer cap.", anchor_text: "1024 to 256", resolved: true, created_at: ago(400) },
      ],
    };
    db.reactions = { pg_loki: { "👍": ["u1", "u2", "u4"], "🚀": ["u5"], "👀": ["u3", "u6"] } };

    db.blogPosts = {
      sp_eng: [
        { id: nid("bp"), space_id: "sp_eng", title: "Q4 reliability retro", author_id: "u1", body: "<p>What we learned shipping the edge fleet this quarter.</p>", published: true, published_at: ago(60 * 24 * 6) },
        { id: nid("bp"), space_id: "sp_eng", title: "Migrating to cgroup v2 — field notes", author_id: "u4", body: "<p>Draft — rollout notes from us-east-1.</p>", published: false, published_at: null },
      ],
      sp_inf: [], sp_me: [],
    };
  })();

  // ─── Helpers ──────────────────────────────────────────────────────────
  const user = (id) => P.PEOPLE.find((p) => p.id === id) || null;
  const ok = (data, status) => ({ status: status || 200, data });
  const err = (status, message) => ({ status, data: { message } });
  const J = (res) => ok(res);

  function pageTree(spaceId) {
    const all = Object.values(db.pages).filter((p) => p.space_id === spaceId);
    const byParent = {};
    all.forEach((p) => { (byParent[p.parent_id || "root"] ||= []).push(p); });
    const build = (parent) => (byParent[parent] || []).map((p) => ({ id: p.id, title: p.title, children: build(p.id) }));
    return build("root");
  }

  function diffLines(a, b) {
    const strip = (h) => h.replace(/<[^>]+>/g, "\n").split("\n").map((s) => s.trim()).filter(Boolean);
    const la = strip(a), lb = strip(b);
    const setA = new Set(la), setB = new Set(lb);
    const out = [];
    la.forEach((l) => { if (!setB.has(l)) out.push({ type: "removed", text: l }); else out.push({ type: "same", text: l }); });
    lb.forEach((l) => { if (!setA.has(l)) out.push({ type: "added", text: l }); });
    return out;
  }

  // ─── Route table ──────────────────────────────────────────────────────
  // Each: [METHOD, RegExp, (m, body, query) => {status, data}]
  const R = [];
  const route = (method, pattern, fn) => R.push([method, pattern, fn]);

  // Issue links
  route("GET", /^\/issues\/([^/]+)\/links$/, (m) => {
    const list = (db.links[m[1]] || []).map((l) => {
      const iss = P.ISSUES.find((i) => i.id === l.issue_id) || { id: l.issue_id, title: "Unknown issue", status: "Backlog", type: "Task" };
      return { id: l.id, link_type: l.link_type, issue: { id: iss.id, title: iss.title, status: iss.status, type: iss.type } };
    });
    return ok(list);
  });
  route("POST", /^\/issues\/([^/]+)\/links$/, (m, b) => {
    const link = { id: nid("lk"), link_type: b.link_type, issue_id: b.linked_issue_id };
    (db.links[m[1]] ||= []).push(link);
    return ok(link, 201);
  });
  route("DELETE", /^\/issues\/links\/([^/]+)$/, (m) => {
    Object.keys(db.links).forEach((k) => { db.links[k] = db.links[k].filter((l) => l.id !== m[1]); });
    return ok({ ok: true });
  });

  // Watchers
  route("GET", /^\/issues\/([^/]+)\/watchers$/, (m) => ok((db.watchers[m[1]] || []).map(user).filter(Boolean)));
  route("POST", /^\/issues\/([^/]+)\/watchers$/, (m, b) => {
    const arr = (db.watchers[m[1]] ||= []);
    const uid = b.user_id || ME.id;
    if (!arr.includes(uid)) arr.push(uid);
    return ok(arr.map(user));
  });
  route("DELETE", /^\/issues\/([^/]+)\/watchers$/, (m) => {
    db.watchers[m[1]] = (db.watchers[m[1]] || []).filter((u) => u !== ME.id);
    return ok(db.watchers[m[1]].map(user));
  });

  // Votes
  route("GET", /^\/issues\/([^/]+)\/votes$/, (m) => {
    const v = (db.votes[m[1]] ||= { voters: [] });
    return ok({ count: v.voters.length, voted_by_me: v.voters.includes(ME.id) });
  });
  route("POST", /^\/issues\/([^/]+)\/votes$/, (m) => {
    const v = (db.votes[m[1]] ||= { voters: [] });
    if (v.voters.includes(ME.id)) v.voters = v.voters.filter((u) => u !== ME.id);
    else v.voters.push(ME.id);
    return ok({ count: v.voters.length, voted_by_me: v.voters.includes(ME.id) });
  });

  // Assignees
  route("GET", /^\/issues\/([^/]+)\/assignees$/, (m) => {
    const iss = P.ISSUES.find((i) => i.id === m[1]);
    const seed = db.assignees[m[1]] || (iss && iss.assignee ? [iss.assignee] : []);
    db.assignees[m[1]] = seed;
    return ok(seed.map(user).filter(Boolean));
  });
  route("PUT", /^\/issues\/([^/]+)\/assignees$/, (m, b) => {
    db.assignees[m[1]] = b.user_ids || [];
    return ok(db.assignees[m[1]].map(user).filter(Boolean));
  });
  route("DELETE", /^\/issues\/([^/]+)\/assignees\/([^/]+)$/, (m) => {
    db.assignees[m[1]] = (db.assignees[m[1]] || []).filter((u) => u !== m[2]);
    return ok(db.assignees[m[1]].map(user).filter(Boolean));
  });

  // Worklogs + time summary
  const summarize = (id) => {
    const logs = db.worklogs[id] || [];
    const spent = logs.reduce((s, w) => s + (w.minutes || 0), 0);
    const orig = (db.estimates[id] || { original: 0 }).original;
    return { original_estimate: orig, time_spent: spent, time_remaining: Math.max(0, orig - spent) };
  };
  const parseMin = (s) => {
    let t = 0; const mm = String(s || "").match(/(\d+)\s*h/); const mn = String(s || "").match(/(\d+)\s*m/);
    if (mm) t += parseInt(mm[1]) * 60; if (mn) t += parseInt(mn[1]); if (!mm && !mn) t = parseInt(s) || 0; return t;
  };
  route("GET", /^\/issues\/([^/]+)\/worklogs$/, (m) => ok((db.worklogs[m[1]] || []).map((w) => ({ ...w, user: user(w.user_id) }))));
  route("POST", /^\/issues\/([^/]+)\/worklogs$/, (m, b) => {
    const w = { id: nid("wl"), user_id: ME.id, time_spent: b.time_spent, minutes: parseMin(b.time_spent), description: b.description || "", started_at: b.started_at || new Date().toISOString() };
    (db.worklogs[m[1]] ||= []).push(w);
    return ok({ ...w, user: user(w.user_id) }, 201);
  });
  route("PUT", /^\/issues\/([^/]+)\/worklogs\/([^/]+)$/, (m, b) => {
    const w = (db.worklogs[m[1]] || []).find((x) => x.id === m[2]);
    if (!w) return err(404, "Worklog not found");
    Object.assign(w, { time_spent: b.time_spent, minutes: parseMin(b.time_spent), description: b.description, started_at: b.started_at });
    return ok({ ...w, user: user(w.user_id) });
  });
  route("DELETE", /^\/issues\/([^/]+)\/worklogs\/([^/]+)$/, (m) => {
    db.worklogs[m[1]] = (db.worklogs[m[1]] || []).filter((x) => x.id !== m[2]);
    return ok({ ok: true });
  });
  route("GET", /^\/issues\/([^/]+)\/time-summary$/, (m) => ok(summarize(m[1])));

  // History
  route("GET", /^\/issues\/([^/]+)\/history$/, (m) => ok((db.history[m[1]] || []).slice().reverse().map((h) => ({ ...h, user: user(h.user_id) }))));

  // Workflows
  route("GET", /^\/workflows$/, () => ok(db.workflows.map((w) => ({ id: w.id, name: w.name, status_count: w.statuses.length }))));
  route("POST", /^\/workflows$/, (m, b) => {
    const w = { id: nid("wf"), name: b.name, statuses: [], transitions: [] };
    db.workflows.push(w);
    return ok({ id: w.id, name: w.name, status_count: 0 }, 201);
  });
  route("GET", /^\/workflows\/([^/]+)$/, (m) => {
    const w = db.workflows.find((x) => x.id === m[1]);
    return w ? ok(w) : err(404, "Workflow not found");
  });
  route("PUT", /^\/workflows\/([^/]+)$/, (m, b) => {
    const w = db.workflows.find((x) => x.id === m[1]);
    if (!w) return err(404, "Not found");
    if (b.name != null) w.name = b.name;
    if (b.status_order) w.statuses.sort((a, c) => b.status_order.indexOf(a.id) - b.status_order.indexOf(c.id));
    return ok(w);
  });
  route("DELETE", /^\/workflows\/([^/]+)$/, (m) => { db.workflows = db.workflows.filter((x) => x.id !== m[1]); return ok({ ok: true }); });
  route("POST", /^\/workflows\/([^/]+)\/statuses$/, (m, b) => {
    const w = db.workflows.find((x) => x.id === m[1]);
    if (!w) return err(404, "Not found");
    const s = { id: nid("st"), name: b.name, category: b.category, color: b.color };
    w.statuses.push(s);
    return ok(s, 201);
  });
  route("PUT", /^\/workflows\/statuses\/([^/]+)$/, (m, b) => {
    for (const w of db.workflows) { const s = w.statuses.find((x) => x.id === m[1]); if (s) { Object.assign(s, b); return ok(s); } }
    return err(404, "Status not found");
  });
  route("DELETE", /^\/workflows\/statuses\/([^/]+)$/, (m) => {
    for (const w of db.workflows) { w.statuses = w.statuses.filter((x) => x.id !== m[1]); w.transitions = w.transitions.filter((t) => t.from_id !== m[1] && t.to_id !== m[1]); }
    return ok({ ok: true });
  });
  route("POST", /^\/workflows\/([^/]+)\/transitions$/, (m, b) => {
    const w = db.workflows.find((x) => x.id === m[1]);
    if (!w) return err(404, "Not found");
    const t = { id: nid("tr"), from_id: b.from_id, to_id: b.to_id };
    w.transitions.push(t);
    return ok(t, 201);
  });
  route("DELETE", /^\/workflows\/transitions\/([^/]+)$/, (m) => {
    for (const w of db.workflows) w.transitions = w.transitions.filter((t) => t.id !== m[1]);
    return ok({ ok: true });
  });

  // Components
  route("GET", /^\/projects\/([^/]+)\/components$/, () => ok(db.components.map((c) => ({ ...c, lead: user(c.lead_id) }))));
  route("POST", /^\/projects\/([^/]+)\/components$/, (m, b) => {
    const c = { id: nid("cp"), name: b.name, description: b.description || "", lead_id: b.lead_id || null, issue_count: 0 };
    db.components.push(c);
    return ok({ ...c, lead: user(c.lead_id) }, 201);
  });
  route("PUT", /^\/components\/([^/]+)$/, (m, b) => {
    const c = db.components.find((x) => x.id === m[1]); if (!c) return err(404, "Not found");
    Object.assign(c, { name: b.name ?? c.name, description: b.description ?? c.description, lead_id: b.lead_id ?? c.lead_id });
    return ok({ ...c, lead: user(c.lead_id) });
  });
  route("DELETE", /^\/components\/([^/]+)$/, (m) => { db.components = db.components.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Versions
  route("GET", /^\/projects\/([^/]+)\/versions$/, () => ok(db.versions));
  route("POST", /^\/projects\/([^/]+)\/versions$/, (m, b) => {
    const v = { id: nid("vr"), name: b.name, start_date: b.start_date || null, release_date: b.release_date || null, status: "unreleased", description: b.description || "" };
    db.versions.push(v);
    return ok(v, 201);
  });
  route("PUT", /^\/versions\/([^/]+)$/, (m, b) => {
    const v = db.versions.find((x) => x.id === m[1]); if (!v) return err(404, "Not found");
    Object.assign(v, b); return ok(v);
  });
  route("DELETE", /^\/versions\/([^/]+)$/, (m) => { db.versions = db.versions.filter((x) => x.id !== m[1]); return ok({ ok: true }); });
  route("POST", /^\/versions\/([^/]+)\/release$/, (m) => {
    const v = db.versions.find((x) => x.id === m[1]); if (!v) return err(404, "Not found");
    v.status = "released"; v.release_date = v.release_date || new Date().toISOString().slice(0, 10);
    return ok(v);
  });

  // Custom fields
  route("GET", /^\/projects\/([^/]+)\/custom-fields$/, () => ok(db.customFields.slice().sort((a, b) => a.position - b.position)));
  route("POST", /^\/projects\/([^/]+)\/custom-fields$/, (m, b) => {
    const f = { id: nid("cf"), name: b.name, type: b.type, description: b.description || "", required: !!b.required, options: b.options || [], position: db.customFields.length };
    db.customFields.push(f); return ok(f, 201);
  });
  route("PUT", /^\/custom-fields\/([^/]+)$/, (m, b) => {
    if (b.reorder) { // {reorder:[ids]}
      b.reorder.forEach((id, i) => { const f = db.customFields.find((x) => x.id === id); if (f) f.position = i; });
      return ok(db.customFields.slice().sort((a, c) => a.position - c.position));
    }
    const f = db.customFields.find((x) => x.id === m[1]); if (!f) return err(404, "Not found");
    Object.assign(f, b); return ok(f);
  });
  route("DELETE", /^\/custom-fields\/([^/]+)$/, (m) => { db.customFields = db.customFields.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Labels
  route("GET", /^\/projects\/([^/]+)\/labels$/, () => ok(db.labels));
  route("POST", /^\/projects\/([^/]+)\/labels$/, (m, b) => {
    const l = { id: nid("lb"), name: b.name, color: b.color || "#6366F1", issue_count: 0 };
    db.labels.push(l); return ok(l, 201);
  });
  route("PUT", /^\/labels\/([^/]+)$/, (m, b) => {
    const l = db.labels.find((x) => x.id === m[1]); if (!l) return err(404, "Not found");
    Object.assign(l, { name: b.name ?? l.name, color: b.color ?? l.color }); return ok(l);
  });
  route("DELETE", /^\/labels\/([^/]+)$/, (m) => { db.labels = db.labels.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Webhooks
  route("GET", /^\/projects\/([^/]+)\/webhooks$/, () => ok(db.webhooks));
  route("POST", /^\/webhooks$/, (m, b) => {
    const w = { id: nid("wh"), url: b.url, secret: b.secret || "", events: b.events || [], is_active: b.is_active !== false, last_status: null, last_code: null };
    db.webhooks.push(w); return ok(w, 201);
  });
  route("GET", /^\/webhooks\/([^/]+)\/deliveries$/, (m) => {
    const out = [];
    for (let i = 0; i < 8; i++) out.push({ id: nid("dl"), event: ["issue.created", "issue.transition", "sprint.completed"][i % 3], status: i % 4 === 0 ? "failed" : "success", code: i % 4 === 0 ? 503 : 200, at: ago(i * 90 + 10), duration_ms: 120 + i * 33 });
    return ok(out);
  });
  route("GET", /^\/webhooks\/([^/]+)$/, (m) => { const w = db.webhooks.find((x) => x.id === m[1]); return w ? ok(w) : err(404, "Not found"); });
  route("PUT", /^\/webhooks\/([^/]+)$/, (m, b) => { const w = db.webhooks.find((x) => x.id === m[1]); if (!w) return err(404, "Not found"); Object.assign(w, b); return ok(w); });
  route("DELETE", /^\/webhooks\/([^/]+)$/, (m) => { db.webhooks = db.webhooks.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // API keys
  route("GET", /^\/api-keys$/, () => ok(db.apiKeys.map((k) => ({ id: k.id, name: k.name, prefix: k.prefix, created_at: k.created_at, last_used_at: k.last_used_at }))));
  route("POST", /^\/api-keys$/, (m, b) => {
    const prefix = "jfk_" + Math.random().toString(36).slice(2, 6);
    const plain = prefix + "_" + Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2);
    const k = { id: nid("ak"), name: b.name, prefix, created_at: new Date().toISOString(), last_used_at: null };
    db.apiKeys.unshift(k);
    return ok({ ...k, plain_key: plain }, 201);
  });
  route("DELETE", /^\/api-keys\/([^/]+)$/, (m) => { db.apiKeys = db.apiKeys.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Audit logs
  route("GET", /^\/audit-logs$/, (m, b, q) => {
    let list = db.auditLogs.slice();
    if (q.user_id) list = list.filter((l) => l.actor_id === q.user_id);
    if (q.action) list = list.filter((l) => l.category === q.action || l.action === q.action);
    if (q.from) list = list.filter((l) => l.created_at >= q.from);
    if (q.to) list = list.filter((l) => l.created_at <= q.to);
    const page = parseInt(q.page || "1"), limit = parseInt(q.limit || "50");
    const total = list.length;
    const items = list.slice((page - 1) * limit, page * limit).map((l) => ({ ...l, actor: user(l.actor_id) }));
    return ok({ items, total, page, limit });
  });
  route("GET", /^\/audit-logs\/export$/, (m, b, q) => {
    const rows = [["timestamp", "actor", "action", "resource_type", "resource_id", "ip"]];
    db.auditLogs.forEach((l) => rows.push([l.created_at, (user(l.actor_id) || {}).name || l.actor_id, l.action, l.resource_type, l.resource_id, (l.meta || {}).ip || ""]));
    const csv = rows.map((r) => r.map((c) => '"' + String(c).replace(/"/g, '""') + '"').join(",")).join("\n");
    return { status: 200, blob: new Blob([csv], { type: "text/csv" }), filename: "audit-logs.csv" };
  });

  // Saved filters
  route("GET", /^\/saved-filters$/, () => ok(db.savedFilters));
  route("POST", /^\/saved-filters$/, (m, b) => { const f = { id: nid("sf"), name: b.name, filter_json: b.filter_json || {} }; db.savedFilters.push(f); return ok(f, 201); });
  route("GET", /^\/saved-filters\/([^/]+)$/, (m) => { const f = db.savedFilters.find((x) => x.id === m[1]); return f ? ok(f) : err(404, "Not found"); });
  route("PUT", /^\/saved-filters\/([^/]+)$/, (m, b) => { const f = db.savedFilters.find((x) => x.id === m[1]); if (!f) return err(404, "Not found"); Object.assign(f, b); return ok(f); });
  route("DELETE", /^\/saved-filters\/([^/]+)$/, (m) => { db.savedFilters = db.savedFilters.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Board swimlanes
  route("GET", /^\/boards\/([^/]+)\/swimlanes$/, (m, b, q) => {
    const type = q.type || db.boardSwimlane;
    return ok({ swimlane_type: type });
  });
  route("PUT", /^\/boards\/([^/]+)\/swimlane-type$/, (m, b) => { db.boardSwimlane = b.swimlane_type; return ok({ swimlane_type: db.boardSwimlane }); });

  // Spaces
  route("GET", /^\/spaces$/, () => ok(db.spaces.filter((s) => !s.archived)));
  route("POST", /^\/spaces$/, (m, b) => { const s = { id: nid("sp"), name: b.name, key: b.key, type: b.type || "team", description: b.description || "", archived: false, page_count: 0 }; db.spaces.push(s); db.blogPosts[s.id] = []; return ok(s, 201); });
  route("GET", /^\/spaces\/([^/]+)\/pages\/tree$/, (m) => ok(pageTree(m[1])));
  route("POST", /^\/spaces\/([^/]+)\/pages$/, (m, b) => {
    const id = nid("pg");
    db.pages[id] = { id, space_id: m[1], title: b.title || "Untitled", parent_id: b.parent_id || null, body: "<p>Start writing…</p>", author_id: ME.id, updated_at: new Date().toISOString(), versions: [{ version: 1, author_id: ME.id, at: new Date().toISOString(), size: 20, body: "<p>Start writing…</p>" }], reactions: {} };
    const sp = db.spaces.find((s) => s.id === m[1]); if (sp) sp.page_count++;
    return ok(db.pages[id], 201);
  });
  route("GET", /^\/spaces\/([^/]+)\/blog-posts$/, (m) => ok((db.blogPosts[m[1]] || []).map((p) => ({ ...p, author: user(p.author_id) }))));
  route("POST", /^\/spaces\/([^/]+)\/blog-posts$/, (m, b) => {
    const p = { id: nid("bp"), space_id: m[1], title: b.title || "Untitled post", author_id: ME.id, body: b.body || "<p>Draft…</p>", published: false, published_at: null };
    (db.blogPosts[m[1]] ||= []).push(p); return ok({ ...p, author: user(p.author_id) }, 201);
  });
  route("GET", /^\/spaces\/([^/]+)$/, (m) => { const s = db.spaces.find((x) => x.id === m[1]); return s ? ok(s) : err(404, "Not found"); });
  route("PUT", /^\/spaces\/([^/]+)$/, (m, b) => { const s = db.spaces.find((x) => x.id === m[1]); if (!s) return err(404, "Not found"); Object.assign(s, b); return ok(s); });
  route("POST", /^\/spaces\/([^/]+)\/archive$/, (m) => { const s = db.spaces.find((x) => x.id === m[1]); if (s) s.archived = true; return ok({ ok: true }); });
  route("DELETE", /^\/spaces\/([^/]+)$/, (m) => { db.spaces = db.spaces.filter((x) => x.id !== m[1]); return ok({ ok: true }); });

  // Pages
  route("GET", /^\/pages\/([^/]+)\/versions\/([^/]+)\/diff\/([^/]+)$/, (m) => {
    const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found");
    const v1 = pg.versions.find((v) => String(v.version) === m[2]);
    const v2 = pg.versions.find((v) => String(v.version) === m[3]);
    if (!v1 || !v2) return err(404, "Version not found");
    return ok({ from: v1.version, to: v2.version, lines: diffLines(v1.body, v2.body) });
  });
  route("GET", /^\/pages\/([^/]+)\/versions\/([^/]+)$/, (m) => {
    const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found");
    const v = pg.versions.find((x) => String(x.version) === m[2]); return v ? ok({ ...v, author: user(v.author_id) }) : err(404, "Not found");
  });
  route("GET", /^\/pages\/([^/]+)\/versions$/, (m) => {
    const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found");
    return ok(pg.versions.slice().reverse().map((v) => ({ version: v.version, author: user(v.author_id), at: v.at, size: v.size })));
  });
  route("PUT", /^\/pages\/([^/]+)\/move$/, (m, b) => { const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found"); pg.parent_id = b.parent_id || null; return ok(pg); });
  // inline comments
  route("GET", /^\/pages\/([^/]+)\/inline-comments$/, (m) => ok((db.inlineComments[m[1]] || []).map((c) => ({ ...c, user: user(c.user_id) }))));
  route("POST", /^\/pages\/([^/]+)\/inline-comments$/, (m, b) => {
    const c = { id: nid("ic"), page_id: m[1], user_id: ME.id, text: b.text, anchor_text: b.anchor_text || "", anchor_offset: b.anchor_offset || 0, resolved: false, created_at: new Date().toISOString() };
    (db.inlineComments[m[1]] ||= []).push(c); return ok({ ...c, user: user(c.user_id) }, 201);
  });
  // reactions
  route("GET", /^\/pages\/([^/]+)\/reactions$/, (m) => {
    const r = db.reactions[m[1]] || {};
    const counts = {}, mine = [];
    Object.entries(r).forEach(([emo, arr]) => { counts[emo] = arr.length; if (arr.includes(ME.id)) mine.push(emo); });
    return ok({ counts, my_reactions: mine });
  });
  route("POST", /^\/pages\/([^/]+)\/reactions$/, (m, b) => {
    const r = (db.reactions[m[1]] ||= {});
    const arr = (r[b.emoji] ||= []);
    if (arr.includes(ME.id)) r[b.emoji] = arr.filter((u) => u !== ME.id); else arr.push(ME.id);
    const counts = {}, mine = [];
    Object.entries(r).forEach(([emo, a]) => { if (a.length) counts[emo] = a.length; if (a.includes(ME.id)) mine.push(emo); });
    return ok({ counts, my_reactions: mine });
  });
  // export
  route("GET", /^\/pages\/([^/]+)\/export\/(html|pdf|md|docx)$/, (m) => {
    const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found");
    const fmt = m[2];
    const text = pg.body.replace(/<[^>]+>/g, (t) => (t === "</p>" || t === "</h2>" || t === "</li>" ? "\n" : ""));
    let blob, ext = fmt;
    if (fmt === "html") blob = new Blob(["<!doctype html><meta charset=utf-8><title>" + pg.title + "</title><h1>" + pg.title + "</h1>" + pg.body], { type: "text/html" });
    else if (fmt === "md") blob = new Blob(["# " + pg.title + "\n\n" + text], { type: "text/markdown" });
    else blob = new Blob(["%" + fmt.toUpperCase() + " export\n\n" + pg.title + "\n\n" + text], { type: fmt === "pdf" ? "application/pdf" : "application/vnd.openxmlformats-officedocument.wordprocessingml.document" });
    return { status: 200, blob, filename: pg.title.replace(/\s+/g, "-").toLowerCase() + "." + ext };
  });
  route("GET", /^\/pages\/([^/]+)$/, (m) => { const pg = db.pages[m[1]]; return pg ? ok({ ...pg, author: user(pg.author_id) }) : err(404, "Not found"); });
  route("PUT", /^\/pages\/([^/]+)$/, (m, b) => {
    const pg = db.pages[m[1]]; if (!pg) return err(404, "Not found");
    if (b.title != null) pg.title = b.title;
    if (b.body != null && b.body !== pg.body) {
      pg.body = b.body; pg.updated_at = new Date().toISOString();
      const nextV = (pg.versions[pg.versions.length - 1]?.version || 0) + 1;
      pg.versions.push({ version: nextV, author_id: ME.id, at: pg.updated_at, size: b.body.length, body: b.body });
    }
    return ok({ ...pg, author: user(pg.author_id) });
  });
  route("DELETE", /^\/pages\/([^/]+)$/, (m) => {
    const pg = db.pages[m[1]]; if (pg) { const sp = db.spaces.find((s) => s.id === pg.space_id); if (sp) sp.page_count--; }
    delete db.pages[m[1]]; return ok({ ok: true });
  });

  // inline comment mutations
  route("PUT", /^\/inline-comments\/([^/]+)$/, (m, b) => {
    for (const k of Object.keys(db.inlineComments)) { const c = db.inlineComments[k].find((x) => x.id === m[1]); if (c) { c.text = b.text ?? c.text; return ok(c); } }
    return err(404, "Not found");
  });
  route("DELETE", /^\/inline-comments\/([^/]+)$/, (m) => {
    for (const k of Object.keys(db.inlineComments)) db.inlineComments[k] = db.inlineComments[k].filter((x) => x.id !== m[1]);
    return ok({ ok: true });
  });
  route("POST", /^\/inline-comments\/([^/]+)\/resolve$/, (m) => {
    for (const k of Object.keys(db.inlineComments)) { const c = db.inlineComments[k].find((x) => x.id === m[1]); if (c) { c.resolved = true; return ok(c); } }
    return err(404, "Not found");
  });
  route("POST", /^\/inline-comments\/([^/]+)\/unresolve$/, (m) => {
    for (const k of Object.keys(db.inlineComments)) { const c = db.inlineComments[k].find((x) => x.id === m[1]); if (c) { c.resolved = false; return ok(c); } }
    return err(404, "Not found");
  });

  // Blog posts
  route("GET", /^\/blog-posts\/([^/]+)$/, (m) => {
    for (const k of Object.keys(db.blogPosts)) { const p = db.blogPosts[k].find((x) => x.id === m[1]); if (p) return ok({ ...p, author: user(p.author_id) }); }
    return err(404, "Not found");
  });
  route("PUT", /^\/blog-posts\/([^/]+)$/, (m, b) => {
    for (const k of Object.keys(db.blogPosts)) { const p = db.blogPosts[k].find((x) => x.id === m[1]); if (p) { Object.assign(p, { title: b.title ?? p.title, body: b.body ?? p.body }); return ok({ ...p, author: user(p.author_id) }); } }
    return err(404, "Not found");
  });
  route("DELETE", /^\/blog-posts\/([^/]+)$/, (m) => {
    for (const k of Object.keys(db.blogPosts)) db.blogPosts[k] = db.blogPosts[k].filter((x) => x.id !== m[1]);
    return ok({ ok: true });
  });
  route("POST", /^\/blog-posts\/([^/]+)\/publish$/, (m) => {
    for (const k of Object.keys(db.blogPosts)) { const p = db.blogPosts[k].find((x) => x.id === m[1]); if (p) { p.published = true; p.published_at = new Date().toISOString(); return ok(p); } }
    return err(404, "Not found");
  });
  route("POST", /^\/blog-posts\/([^/]+)\/unpublish$/, (m) => {
    for (const k of Object.keys(db.blogPosts)) { const p = db.blogPosts[k].find((x) => x.id === m[1]); if (p) { p.published = false; p.published_at = null; return ok(p); } }
    return err(404, "Not found");
  });

  // Notification preferences
  route("GET", /^\/notifications\/preferences$/, () => ok(clone(db.notifPrefs)));
  route("PUT", /^\/notifications\/preferences$/, (m, b) => { Object.assign(db.notifPrefs, b); return ok(clone(db.notifPrefs)); });

  // Activity feed
  route("GET", /^\/activity$/, (m, b, q) => {
    const limit = parseInt(q.limit || "20");
    return ok(db.activity.slice(0, limit).map((a) => ({ ...a, actor: user(a.actor_id) })));
  });

  // Sprint capacity
  route("GET", /^\/sprints\/([^/]+)\/capacity$/, (m) => {
    const members = (db.sprintCapacity[m[1]] || db.sprintCapacity.s24).map((c) => ({ ...c, user: user(c.user_id) }));
    return ok({ members });
  });
  route("PUT", /^\/sprints\/([^/]+)\/capacity$/, (m, b) => {
    db.sprintCapacity[m[1]] = (db.sprintCapacity[m[1]] || db.sprintCapacity.s24).map((c) => {
      const upd = (b.members || []).find((x) => x.user_id === c.user_id);
      return upd ? { ...c, available_hours: upd.available_hours } : c;
    });
    return ok({ members: db.sprintCapacity[m[1]].map((c) => ({ ...c, user: user(c.user_id) })) });
  });

  // Sprint planning (bulk)
  route("GET", /^\/projects\/([^/]+)\/sprint-planning$/, () => ok({ sprints: ["Sprint 24", "Sprint 25"], backlog: P.ISSUES.filter((i) => !i.sprint).map((i) => i.id) }));
  route("POST", /^\/projects\/([^/]+)\/sprint-planning$/, () => ok({ ok: true }));

  // Import
  route("POST", /^\/import\/(jira|trello|linear)$/, (m) => {
    const id = nid("imp");
    db.importJobs[id] = { id, source: m[1], status: "pending", progress: 0, started: Date.now(), summary: null, error: null };
    return ok({ id, status: "pending" }, 202);
  });
  route("GET", /^\/import\/([^/]+)$/, (m) => {
    const j = db.importJobs[m[1]]; if (!j) return err(404, "Not found");
    const elapsed = (Date.now() - j.started) / 1000;
    if (elapsed < 2) { j.status = "pending"; j.progress = 5; }
    else if (elapsed < 9) { j.status = "processing"; j.progress = Math.min(95, Math.round((elapsed / 9) * 100)); }
    else {
      // deterministic: linear "fails" to exercise the error path, others succeed
      if (j.source === "linear") { j.status = "failed"; j.progress = 64; j.error = "Malformed CSV on row 412: missing required column 'state'."; }
      else { j.status = "done"; j.progress = 100; j.summary = { issues_created: j.source === "jira" ? 248 : 96, projects_created: j.source === "jira" ? 4 : 2 }; }
    }
    return ok(j);
  });

  // ─── Dispatcher ───────────────────────────────────────────────────────
  function dispatch(method, path, query, body) {
    for (const [mth, re, fn] of R) {
      if (mth !== method) continue;
      const mt = re.exec(path);
      if (mt) return fn(mt, body, query);
    }
    return err(404, "No mock route for " + method + " " + path);
  }

  const realFetch = window.fetch.bind(window);
  window.fetch = function (input, init) {
    const url = typeof input === "string" ? input : (input && input.url) || "";
    const idx = url.indexOf("/api/v1");
    if (idx === -1) return realFetch(input, init);

    const method = ((init && init.method) || "GET").toUpperCase();
    const rest = url.slice(idx + 7);
    const [rawPath, qs] = rest.split("?");
    const path = rawPath.replace(/\/$/, "") || "/";
    const query = Object.fromEntries(new URLSearchParams(qs || ""));
    const auth = init && init.headers && (init.headers.Authorization || init.headers.authorization);

    let body = null;
    if (init && init.body && typeof init.body === "string") { try { body = JSON.parse(init.body); } catch (e) { body = null; } }
    else if (init && init.body instanceof FormData) { body = { _file: init.body.get("file") }; }

    const latency = 120 + Math.random() * 260;
    return new Promise((resolve) => {
      setTimeout(() => {
        let result;
        if (!auth) result = err(401, "Missing Authorization header");
        else { try { result = dispatch(method, path, query, body); } catch (e) { result = err(500, e.message); } }

        const status = result.status || 200;
        const respBlob = result.blob || new Blob([JSON.stringify(result.data)], { type: "application/json" });
        const resp = new Response(respBlob, { status, statusText: String(status), headers: { "Content-Type": result.blob ? respBlob.type : "application/json" } });
        if (result.filename) Object.defineProperty(resp, "_filename", { value: result.filename });
        resolve(resp);
      }, latency);
    });
  };

  // ─── api() helper + hooks (used by all feature code) ──────────────────
  function api(path, opts) {
    opts = opts || {};
    const headers = { Authorization: "Bearer " + localStorage.getItem("token") };
    if (opts.body !== undefined) headers["Content-Type"] = "application/json";
    return fetch("/api/v1" + path, {
      method: opts.method || (opts.body !== undefined ? "POST" : "GET"),
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    }).then((res) => {
      if (res.status === 401) { window.location.hash = "#/login"; throw new Error("Session expired — please sign in."); }
      if (opts.raw) return res;
      return res.json().catch(() => ({})).then((data) => {
        if (!res.ok) throw new Error((data && data.message) || "Request failed (" + res.status + ")");
        return data;
      });
    });
  }

  // Download a blob endpoint as a file.
  async function apiDownload(path) {
    const res = await api(path, { raw: true });
    if (!res.ok) throw new Error("Export failed (" + res.status + ")");
    const blob = await res.blob();
    const a = document.createElement("a");
    a.href = URL.createObjectURL(blob);
    a.download = res._filename || path.split("/").pop();
    document.body.appendChild(a); a.click(); a.remove();
    setTimeout(() => URL.revokeObjectURL(a.href), 4000);
  }

  // Upload FormData (imports).
  function apiUpload(path, file) {
    const fd = new FormData(); fd.append("file", file);
    return fetch("/api/v1" + path, { method: "POST", headers: { Authorization: "Bearer " + localStorage.getItem("token") }, body: fd })
      .then((r) => r.json());
  }

  // Hook: load data on mount. Returns {data, loading, error, reload, setData}.
  function useApi(path, deps) {
    const [data, setData] = React.useState(null);
    const [loading, setLoading] = React.useState(true);
    const [error, setError] = React.useState(null);
    const [tick, setTick] = React.useState(0);
    React.useEffect(() => {
      let live = true;
      setLoading(true); setError(null);
      if (!path) { setLoading(false); return; }
      api(path).then((d) => { if (live) { setData(d); setLoading(false); } })
        .catch((e) => { if (live) { setError(e.message); setLoading(false); } });
      return () => { live = false; };
    }, [path, tick].concat(deps || []));
    return { data, loading, error, reload: () => setTick((t) => t + 1), setData };
  }

  Object.assign(window, { api, apiDownload, apiUpload, useApi, FORGE_DB: db });
})();
