// views_misc.jsx — Wiki, Projects, Sprint, Reports, Settings, Admin

// ─── Wiki view ──────────────────────────────────────────
function WikiView({ nav }) {
  const [activePage, setActivePage] = React.useState("rb-loki");
  const [editing, setEditing] = React.useState(false);
  const [expanded, setExpanded] = React.useState({ root: true, onboarding: true, runbooks: true });

  function renderNode(node, depth = 0) {
    return (
      <div key={node.id}>
        <div className="tree-item" aria-current={activePage === node.id ? "page" : undefined} style={{ paddingLeft: 8 + depth * 12 }} onClick={() => {
          if (node.children) setExpanded((e) => ({ ...e, [node.id]: !e[node.id] }));
          else setActivePage(node.id);
        }}>
          {node.children ? (
            <Icon name={expanded[node.id] ? "chevronDown" : "chevronRight"} size={11} color="var(--text-muted)"/>
          ) : (
            <Icon name="notes" size={13}/>
          )}
          <span>{node.title}</span>
        </div>
        {node.children && expanded[node.id] && node.children.map((c) => renderNode(c, depth + 1))}
      </div>
    );
  }

  return (
    <div className="wiki">
      <aside className="wiki-tree">
        <div className="row" style={{ justifyContent: "space-between", padding: "0 8px 10px" }}>
          <h4>Pages</h4>
          <button className="icon-btn" style={{ width: 22, height: 22 }} aria-label="New page"><Icon name="plus" size={13}/></button>
        </div>
        <div className="search" style={{ margin: "0 8px 10px", padding: "4px 10px", background: "var(--bg)" }}>
          <Icon name="search" size={13}/>
          <input placeholder="Search pages…"/>
        </div>
        {FORGE_DATA.WIKI_TREE.map((n) => renderNode(n))}
      </aside>

      <div className="wiki-doc">
        <div className="row gap-2 text-sm muted" style={{ marginBottom: 14 }}>
          <span>Engineering</span> <Icon name="chevronRight" size={11}/>
          <span>Runbooks</span> <Icon name="chevronRight" size={11}/>
          <span style={{ color: "var(--text)" }}>Loki ingester OOM</span>
          <div style={{ flex: 1 }}/>
          <button className="btn btn-ghost" data-size="sm" onClick={() => setEditing((e) => !e)}>
            <Icon name={editing ? "eye" : "pencil"} size={13}/> {editing ? "Preview" : "Edit"}
          </button>
          <button className="btn btn-ghost" data-size="sm"><Icon name="history" size={13}/> History</button>
          <button className="btn btn-ghost" data-size="sm"><Icon name="moreH" size={14}/></button>
        </div>

        {editing && (
          <div style={{ border: "1px solid var(--border)", borderRadius: 8, padding: 6, marginBottom: 16, display: "flex", gap: 2, background: "var(--bg-subtle)" }}>
            {[["heading","H"],["bold","B"],["italic","I"],["list","List"],["table","Table"],["code","Code"],["link","Link"],["picture","Img"],["paperclip","File"]].map(([ic,n]) => (
              <button key={n} className="icon-btn" style={{ width: 28, height: 28 }} title={n}><Icon name={ic} size={14}/></button>
            ))}
          </div>
        )}

        <h1 contentEditable={editing} suppressContentEditableWarning>Runbook — Loki ingester OOM</h1>
        <div className="row gap-3 muted text-sm" style={{ marginBottom: 24, paddingBottom: 14, borderBottom: "1px solid var(--border)" }}>
          <Avatar user={FORGE_DATA.PEOPLE[2]} size="sm"/>
          <span>Priya Raman</span>
          <span>·</span>
          <span>Last edited 28m ago</span>
          <span>·</span>
          <span>Linked to <a href="#/issue/INFRA-232" onClick={(e) => { e.preventDefault(); nav("issue/INFRA-232"); }} className="mono" style={{ color: "var(--indigo-600)", textDecoration: "none" }}>INFRA-232</a></span>
        </div>

        <p contentEditable={editing} suppressContentEditableWarning>
          This runbook covers the symptoms, diagnostics, and mitigation steps for Loki ingester pods getting OOMKilled under steady-state log volume in production.
        </p>

        <blockquote contentEditable={editing} suppressContentEditableWarning>
          <strong>TL;DR</strong> — If ingester memory climbs above 4GB and stays there, scale ingesters horizontally first and rotate chunk store afterwards. Rollback to v3.1.2 is the nuclear option.
        </blockquote>

        <h2>Symptoms</h2>
        <ul style={{ lineHeight: 1.8, paddingLeft: 20 }}>
          <li><code>kubectl get pods -n observability</code> shows <code>OOMKilled</code> in restart history</li>
          <li>Log ingestion lag spikes to 8+ minutes (alert: <code>LokiIngestionLagHigh</code>)</li>
          <li>Grafana panel "Loki — ingester memory" shows ramps past 4GB</li>
        </ul>

        <h2>Diagnostics</h2>
        <p>Run these in order. Stop as soon as you find the root cause.</p>
        <pre contentEditable={editing} suppressContentEditableWarning>{`# 1. Confirm OOM, not eviction
kubectl describe pod -n observability loki-ingester-0 | grep -A2 "Last State"

# 2. Heap profile (port-forward first)
go tool pprof -top -nodecount=20 http://localhost:9095/debug/pprof/heap

# 3. Check chunk encoder buffer count
curl -s localhost:3100/metrics | grep loki_chunk_encoder_pool`}</pre>

        <h2>Mitigation</h2>
        <h3>1. Horizontal scale (5 minutes)</h3>
        <p>
          Scale ingesters from 6 to 9 replicas. This is the safest first move and rarely makes things worse.
        </p>
        <pre>{`kubectl scale -n observability statefulset/loki-ingester --replicas=9`}</pre>

        <h3>2. Reduce chunk encoder buffer pool (15 minutes)</h3>
        <p>
          If horizontal scale alone doesn't hold, the chunk encoder buffer pool is leaking past <code>chunk_encoder_pool_max_buffers</code>. Patch the config:
        </p>
        <pre>{`# values.yaml
ingester:
  chunk_encoder_pool_max_buffers: 256  # was 1024`}</pre>

        <h3>3. Rollback to v3.1.2 (last resort)</h3>
        <p>
          See <a href="#" style={{ color: "var(--indigo-600)" }}>Postmortem: Nov 8 etcd leader churn</a> for the rollback playbook.
        </p>

        <h2>Owners</h2>
        <div className="row gap-2" style={{ marginBottom: 32 }}>
          <Avatar user={FORGE_DATA.PEOPLE[2]} size="md"/>
          <Avatar user={FORGE_DATA.PEOPLE[4]} size="md"/>
          <Avatar user={FORGE_DATA.PEOPLE[0]} size="md"/>
        </div>

        {/* Comments on page */}
        <div style={{ borderTop: "1px solid var(--border)", paddingTop: 16, marginTop: 32 }}>
          <h3 style={{ fontSize: 13, fontWeight: 600, color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: ".04em", margin: "0 0 12px" }}>2 page comments</h3>
          <div className="row gap-3" style={{ marginBottom: 14, alignItems: "flex-start" }}>
            <Avatar user={FORGE_DATA.PEOPLE[1]} size="md"/>
            <div>
              <div className="bold text-sm">Diego Alvarez <span className="muted" style={{ fontWeight: 400, fontSize: 11.5 }}>· yesterday</span></div>
              <div className="text-sm" style={{ marginTop: 2 }}>Step 2 should mention to also check the cgroup memory.limit, not just the container limit. Different numbers depending on whether you're on cgroup v1 or v2.</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// ─── Projects list ──────────────────────────────────────
function ProjectsView({ nav }) {
  const [view, setView] = React.useState("grid");
  const [openCreate, setOpenCreate] = React.useState(false);

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">Forge <Icon name="chevronRight" size={11}/> <span>All projects</span></div>
          <h1>Projects</h1>
          <p>{FORGE_DATA.PROJECTS.length} projects · 162 open issues across the workspace.</p>
        </div>
        <div className="row gap-2">
          <div className="row" style={{ border: "1px solid var(--border)", borderRadius: 6, overflow: "hidden" }}>
            <button onClick={() => setView("grid")} className="btn btn-ghost" data-size="sm" style={{ borderRadius: 0, background: view === "grid" ? "var(--bg-subtle)" : undefined, color: view === "grid" ? "var(--text)" : undefined }}>
              <Icon name="grid" size={14}/>
            </button>
            <button onClick={() => setView("list")} className="btn btn-ghost" data-size="sm" style={{ borderRadius: 0, background: view === "list" ? "var(--bg-subtle)" : undefined, color: view === "list" ? "var(--text)" : undefined }}>
              <Icon name="list" size={14}/>
            </button>
          </div>
          <Button variant="primary" icon="plus" onClick={() => setOpenCreate(true)}>New project</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px" }}>
        {view === "grid" ? (
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 14 }}>
            {FORGE_DATA.PROJECTS.map((p) => {
              const lead = FORGE_DATA.PEOPLE.find((u) => u.id === p.lead);
              return (
                <div key={p.id} className="card" style={{ padding: 0, overflow: "hidden", cursor: "default" }}
                  onMouseEnter={(e) => e.currentTarget.style.borderColor = "var(--border-strong)"}
                  onMouseLeave={(e) => e.currentTarget.style.borderColor = ""}
                  onClick={() => nav("board")}>
                  <div style={{ height: 64, background: "linear-gradient(135deg, " + p.color + " 0%, " + p.color + "AA 100%)", position: "relative", display: "flex", alignItems: "center", padding: "0 16px" }}>
                    <div style={{ width: 40, height: 40, borderRadius: 9, background: "rgba(255,255,255,.20)", color: "#fff", display: "grid", placeItems: "center", boxShadow: "inset 0 0 0 1px rgba(255,255,255,.25)" }}>
                      <Icon name={p.icon} size={20}/>
                    </div>
                  </div>
                  <div style={{ padding: 16 }}>
                    <div className="row gap-2" style={{ marginBottom: 4 }}>
                      <span className="bold">{p.name}</span>
                      <span className="mono text-xs muted">{p.key}</span>
                    </div>
                    <div className="text-sm secondary" style={{ minHeight: 40, marginBottom: 12 }}>{p.desc}</div>
                    <div className="row" style={{ justifyContent: "space-between" }}>
                      <div className="row gap-2 text-xs muted">
                        <Icon name="checkbox" size={12}/> {p.openIssues}
                        <span>·</span>
                        <Icon name="users" size={12}/> {p.members}
                      </div>
                      {lead && <Avatar user={lead} size="sm" title={"Lead: " + lead.name}/>}
                    </div>
                  </div>
                </div>
              );
            })}
            {/* Create card */}
            <div className="card" style={{ display: "grid", placeItems: "center", borderStyle: "dashed", padding: 24, cursor: "default", minHeight: 196 }} onClick={() => setOpenCreate(true)}>
              <div style={{ textAlign: "center", color: "var(--text-muted)" }}>
                <div style={{ width: 40, height: 40, borderRadius: 10, background: "var(--bg-subtle)", display: "grid", placeItems: "center", margin: "0 auto 8px" }}>
                  <Icon name="plus" size={20}/>
                </div>
                <div className="bold text-sm" style={{ color: "var(--text)" }}>New project</div>
                <div className="text-xs muted">Start from scratch or a template</div>
              </div>
            </div>
          </div>
        ) : (
          <div className="card" style={{ overflow: "hidden" }}>
            <table className="table">
              <thead><tr><th>Name</th><th>Key</th><th>Lead</th><th>Members</th><th>Open</th><th>Updated</th></tr></thead>
              <tbody>
                {FORGE_DATA.PROJECTS.map((p) => {
                  const lead = FORGE_DATA.PEOPLE.find((u) => u.id === p.lead);
                  return (
                    <tr key={p.id} onClick={() => nav("board")}>
                      <td><div className="row gap-3"><div style={{ width: 28, height: 28, borderRadius: 7, background: p.color, color: "#fff", display: "grid", placeItems: "center" }}><Icon name={p.icon} size={14}/></div><span>{p.name}</span></div></td>
                      <td className="mono text-xs">{p.key}</td>
                      <td>{lead && <span className="row gap-2"><Avatar user={lead} size="sm"/>{lead.name}</span>}</td>
                      <td>{p.members}</td>
                      <td>{p.openIssues}</td>
                      <td className="muted text-xs">{p.updated}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <Modal open={openCreate} onClose={() => setOpenCreate(false)} title="Create new project"
        footer={<><Button onClick={() => setOpenCreate(false)}>Cancel</Button><Button variant="primary" onClick={() => setOpenCreate(false)}>Create project</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Project name</label><input className="input" placeholder="e.g. Platform Migrations"/></div>
          <div className="row gap-3">
            <div style={{ flex: 1 }}><label className="label">Key</label><input className="input mono" placeholder="PLAT" maxLength="5"/></div>
            <div style={{ flex: 2 }}><label className="label">Lead</label><PillSelect value="u1" onChange={() => {}} options={FORGE_DATA.PEOPLE.map((p) => ({ id: p.id, label: p.name }))}/></div>
          </div>
          <div><label className="label">Description</label><textarea className="textarea" placeholder="What does this team own?"/></div>
          <div><label className="label">Template</label>
            <div className="row gap-2" style={{ flexWrap: "wrap" }}>
              {["Kanban","Scrum","Bug tracking","Empty"].map((t, i) => (
                <button key={t} className="btn" style={{ border: "1px solid " + (i === 1 ? "var(--indigo-600)" : "var(--border)"), background: i === 1 ? "var(--indigo-50)" : "var(--bg)", color: i === 1 ? "var(--indigo-700)" : "var(--text)" }}>{t}</button>
              ))}
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
}

// ─── Sprint planning ────────────────────────────────────
function SprintView({ nav, issues }) {
  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Sprints</span></div>
          <h1>Sprints</h1>
          <p>Plan, run, and review iterations.</p>
        </div>
        <div className="row gap-2">
          <Button icon="chart">Velocity report</Button>
          <Button variant="primary" icon="plus">Start sprint</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px" }}>
        {/* Velocity */}
        <div className="card" style={{ marginBottom: 16 }}>
          <div className="card-head">
            <h3>Velocity — last 6 sprints</h3>
            <div className="row gap-3 text-xs muted">
              <span className="row gap-2"><span style={{ width: 10, height: 10, borderRadius: 2, background: "#6366F1" }}/>Committed</span>
              <span className="row gap-2"><span style={{ width: 10, height: 10, borderRadius: 2, background: "#10B981" }}/>Completed</span>
            </div>
          </div>
          <div style={{ padding: 24 }}>
            <Velocity/>
          </div>
        </div>

        {/* Feature 25: Capacity */}
        <SprintCapacity sprintId="s24"/>

        {/* Active sprint */}
        <div className="card" style={{ marginBottom: 16 }}>
          <div className="card-head">
            <div className="row gap-3">
              <Badge tone="info" dot>Active</Badge>
              <h3>Sprint 24 — Edge & Reliability</h3>
              <span className="text-xs muted">Dec 2 → Dec 15</span>
            </div>
            <Button data-size="sm" variant="primary">Complete sprint</Button>
          </div>
          <div style={{ padding: 16 }}>
            <div style={{ display: "grid", gridTemplateColumns: "repeat(4, 1fr)", gap: 12, marginBottom: 16 }}>
              <SprintStat label="Committed" value="62 pts" trend=""/>
              <SprintStat label="Completed" value="42 pts" trend="68%" color="var(--success)"/>
              <SprintStat label="In progress" value="13 pts" color="var(--info)"/>
              <SprintStat label="Days remaining" value="5"/>
            </div>
            <h4 style={{ margin: "0 0 8px", fontSize: 12, fontWeight: 600, color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: ".04em" }}>Sprint issues</h4>
            <div style={{ border: "1px solid var(--border)", borderRadius: 8, overflow: "hidden", maxHeight: 320, overflowY: "auto" }}>
              {issues.filter((i) => i.sprint).map((i, idx) => {
                const u = FORGE_DATA.PEOPLE.find((p) => p.id === i.assignee);
                return (
                  <div key={i.id} className="row gap-3" style={{ padding: "8px 12px", borderBottom: idx < issues.filter((i) => i.sprint).length - 1 ? "1px solid var(--border)" : 0 }} onClick={() => nav("issue/" + i.id)}>
                    <TypeIcon value={i.type}/>
                    <span className="mono text-xs muted" style={{ width: 80 }}>{i.id}</span>
                    <span className="text-sm" style={{ flex: 1, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>{i.title}</span>
                    <PriorityBadge value={i.pri}/>
                    <StatusBadge value={i.status}/>
                    <span className="mono text-xs muted" style={{ width: 20, textAlign: "right" }}>{i.points}</span>
                    {u && <Avatar user={u} size="sm"/>}
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Planned + done */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
          <div className="card">
            <div className="card-head">
              <div className="row gap-3"><Badge tone="muted" dot>Planned</Badge><h3>Sprint 25 — Observability deep clean</h3></div>
              <Button data-size="sm">Start</Button>
            </div>
            <div style={{ padding: 16 }} className="text-sm secondary">
              23 issues · 52 points committed · Dec 16 → Dec 29
            </div>
          </div>
          <div className="card">
            <div className="card-head">
              <div className="row gap-3"><Badge tone="success" dot>Completed</Badge><h3>Sprint 23 — IAM cleanup</h3></div>
              <Button data-size="sm" icon="chart">Report</Button>
            </div>
            <div style={{ padding: 16 }} className="text-sm secondary">
              19 / 21 issues delivered · 47 / 51 points · Nov 18 → Dec 1
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function SprintStat({ label, value, trend, color }) {
  return (
    <div style={{ padding: 12, background: "var(--bg-subtle)", borderRadius: 8 }}>
      <div className="text-xs muted">{label}</div>
      <div style={{ fontSize: 22, fontWeight: 600, color: color || "var(--text)" }}>{value}</div>
      {trend && <div className="text-xs" style={{ color: color || "var(--text-muted)" }}>{trend}</div>}
    </div>
  );
}

function Velocity() {
  const data = [
    { label: "S19", committed: 48, completed: 46 },
    { label: "S20", committed: 52, completed: 49 },
    { label: "S21", committed: 55, completed: 52 },
    { label: "S22", committed: 60, completed: 54 },
    { label: "S23", committed: 51, completed: 47 },
    { label: "S24", committed: 62, completed: 42 },
  ];
  const max = 70;
  const h = 200;
  return (
    <div style={{ display: "flex", gap: 32, alignItems: "flex-end", height: h + 30, justifyContent: "space-around" }}>
      {data.map((d, i) => (
        <div key={i} style={{ display: "flex", flexDirection: "column", alignItems: "center", flex: 1 }}>
          <div className="row gap-1" style={{ alignItems: "flex-end", height: h, width: "100%", justifyContent: "center" }}>
            <div style={{ width: "28%", maxWidth: 30, height: (d.committed / max) * h, background: "#6366F1", borderRadius: "4px 4px 0 0" }}/>
            <div style={{ width: "28%", maxWidth: 30, height: (d.completed / max) * h, background: "#10B981", borderRadius: "4px 4px 0 0" }}/>
          </div>
          <div className="text-xs muted" style={{ marginTop: 8 }}>{d.label}</div>
          <div className="text-xs bold" style={{ color: i === 5 ? "var(--warning)" : "var(--text)" }}>{d.completed}/{d.committed}</div>
        </div>
      ))}
    </div>
  );
}

// ─── Reports ────────────────────────────────────────────
function ReportsView({ nav, issues }) {
  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Reports</span></div>
          <h1>Reports & analytics</h1>
          <p>Insights across velocity, workload and throughput.</p>
        </div>
        <div className="row gap-2">
          <Button icon="calendar">Last 30 days</Button>
          <Button icon="download">Export PDF</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", display: "grid", gridTemplateColumns: "repeat(12, 1fr)", gap: 16 }}>
        <div className="card" style={{ gridColumn: "span 8" }}>
          <div className="card-head">
            <h3>Burndown — Sprint 24</h3>
            <span className="text-xs muted">Updated 12m ago</span>
          </div>
          <div style={{ padding: 16, height: 320 }}>
            <Burndown/>
          </div>
        </div>
        <div className="card" style={{ gridColumn: "span 4" }}>
          <div className="card-head"><h3>Issue distribution by type</h3></div>
          <div style={{ padding: 16, display: "grid", gridTemplateColumns: "180px 1fr", gap: 12, alignItems: "center" }}>
            <Donut data={[
              { label: "Story", value: 38, color: "#8B5CF6" },
              { label: "Task",  value: 32, color: "#3B82F6" },
              { label: "Bug",   value: 22, color: "#EF4444" },
              { label: "Epic",  value: 8,  color: "#F97316" },
            ]}/>
            <div className="stack gap-2">
              {[
                ["Story",38,"#8B5CF6"],["Task",32,"#3B82F6"],["Bug",22,"#EF4444"],["Epic",8,"#F97316"]
              ].map(([l,v,c]) => (
                <div key={l} className="row gap-2 text-sm">
                  <span style={{ width: 10, height: 10, borderRadius: 2, background: c }}/>
                  <span style={{ flex: 1 }}>{l}</span>
                  <span className="bold">{v}%</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className="card" style={{ gridColumn: "span 7" }}>
          <div className="card-head"><h3>Team workload (this sprint)</h3></div>
          <div style={{ padding: 16 }}>
            {FORGE_DATA.PEOPLE.slice(0, 7).map((u, i) => {
              const pts = [13, 11, 14, 9, 8, 5, 2][i];
              const cap = 13;
              const pct = Math.min(100, (pts / cap) * 100);
              return (
                <div key={u.id} className="row gap-3" style={{ marginBottom: 10 }}>
                  <Avatar user={u} size="sm"/>
                  <span className="text-sm" style={{ width: 140 }}>{u.name}</span>
                  <div style={{ flex: 1, height: 8, background: "var(--bg-muted)", borderRadius: 99, overflow: "hidden" }}>
                    <div style={{ width: pct + "%", height: "100%", background: pct > 100 ? "var(--danger)" : pct > 85 ? "var(--warning)" : "var(--indigo-600)" }}/>
                  </div>
                  <span className="text-xs muted mono" style={{ width: 60, textAlign: "right" }}>{pts}/{cap}pt</span>
                </div>
              );
            })}
          </div>
        </div>

        <div className="card" style={{ gridColumn: "span 5" }}>
          <div className="card-head"><h3>By priority</h3></div>
          <div style={{ padding: 16 }}>
            {[
              ["Critical", 3,  "#DC2626"],
              ["High",     12, "#EA580C"],
              ["Medium",   24, "#3B82F6"],
              ["Low",      8,  "#64748B"],
            ].map(([l, v, c]) => (
              <div key={l} className="row gap-3" style={{ marginBottom: 10 }}>
                <span className="text-sm" style={{ width: 80, color: c, fontWeight: 500 }}>{l}</span>
                <div style={{ flex: 1, height: 22, background: "var(--bg-subtle)", borderRadius: 4, overflow: "hidden", position: "relative" }}>
                  <div style={{ width: (v / 24) * 100 + "%", height: "100%", background: c, opacity: .8 }}/>
                </div>
                <span className="bold text-sm" style={{ width: 24, textAlign: "right" }}>{v}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="card" style={{ gridColumn: "span 12" }}>
          <div className="card-head">
            <h3>Throughput (issues closed per week)</h3>
            <span className="text-xs muted">12 weeks</span>
          </div>
          <div style={{ padding: 16 }}>
            <Throughput/>
          </div>
        </div>
      </div>
    </div>
  );
}

function Donut({ data }) {
  const total = data.reduce((s, d) => s + d.value, 0);
  let acc = 0;
  const r = 60, c = 70;
  return (
    <svg viewBox="0 0 140 140" width="160" height="160" style={{ margin: "auto", display: "block" }}>
      {data.map((d, i) => {
        const start = (acc / total) * Math.PI * 2 - Math.PI / 2;
        acc += d.value;
        const end = (acc / total) * Math.PI * 2 - Math.PI / 2;
        const x1 = c + Math.cos(start) * r, y1 = c + Math.sin(start) * r;
        const x2 = c + Math.cos(end)   * r, y2 = c + Math.sin(end)   * r;
        const large = d.value / total > 0.5 ? 1 : 0;
        return <path key={i} d={`M ${c} ${c} L ${x1} ${y1} A ${r} ${r} 0 ${large} 1 ${x2} ${y2} Z`} fill={d.color}/>;
      })}
      <circle cx={c} cy={c} r="36" fill="var(--bg)"/>
      <text x={c} y={c - 2} textAnchor="middle" fontSize="20" fontWeight="600" fill="var(--text)">{total}</text>
      <text x={c} y={c + 14} textAnchor="middle" fontSize="10" fill="var(--text-muted)">issues</text>
    </svg>
  );
}

function Burndown() {
  const ideal   = [62, 58, 53, 48, 44, 40, 35, 31, 27, 22, 18, 13, 9, 4, 0];
  const actual  = [62, 60, 58, 51, 48, 44, 41, 38, 34, 28, 21];
  const w = 800, h = 260, padL = 40, padR = 16, padT = 16, padB = 30;
  const max = 64;
  const stepX = (w - padL - padR) / (ideal.length - 1);
  const y = (v) => h - padB - (v / max) * (h - padT - padB);
  return (
    <svg viewBox={`0 0 ${w} ${h}`} style={{ width: "100%", height: "100%" }}>
      {/* grid */}
      {[0, 16, 32, 48, 64].map((v) => (
        <g key={v}>
          <line x1={padL} y1={y(v)} x2={w - padR} y2={y(v)} stroke="var(--border)" strokeDasharray="2 4"/>
          <text x={padL - 8} y={y(v) + 4} fontSize="10" fill="var(--text-muted)" textAnchor="end">{v}</text>
        </g>
      ))}
      {/* x labels */}
      {ideal.map((_, i) => i % 2 === 0 && (
        <text key={i} x={padL + i * stepX} y={h - padB + 16} fontSize="10" fill="var(--text-muted)" textAnchor="middle">Day {i + 1}</text>
      ))}
      <path d={ideal.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ")} fill="none" stroke="var(--text-muted)" strokeDasharray="4 4"/>
      <path d={actual.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ") + ` L ${padL + (actual.length - 1) * stepX} ${h - padB} L ${padL} ${h - padB} Z`} fill="#6366F1" fillOpacity=".12"/>
      <path d={actual.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ")} fill="none" stroke="#6366F1" strokeWidth="2.5"/>
      {actual.map((v, i) => (
        <circle key={i} cx={padL + i * stepX} cy={y(v)} r="3" fill="#6366F1"/>
      ))}
      <text x={w - padR} y={padT + 4} fontSize="11" fill="var(--text-muted)" textAnchor="end">— — Ideal</text>
      <text x={w - padR} y={padT + 18} fontSize="11" fill="#6366F1" textAnchor="end">Actual</text>
    </svg>
  );
}

function Throughput() {
  const data = [4, 6, 5, 8, 7, 9, 11, 8, 10, 12, 14, 13];
  const max = Math.max(...data);
  return (
    <div className="row" style={{ alignItems: "flex-end", gap: 6, height: 140, padding: "0 8px" }}>
      {data.map((v, i) => (
        <div key={i} style={{ flex: 1, display: "flex", flexDirection: "column", alignItems: "center" }}>
          <div style={{ width: "70%", height: (v / max) * 110, background: `linear-gradient(180deg, #818CF8, #6366F1)`, borderRadius: "3px 3px 0 0" }}/>
          <div className="text-xs muted" style={{ marginTop: 4, fontSize: 10 }}>W{i + 1}</div>
        </div>
      ))}
    </div>
  );
}

// ─── Settings ───────────────────────────────────────────
function SettingsView({ nav, tg }) {
  const [tab, setTab] = React.useState("general");
  return (
    <div>
      <div className="page-head" style={{ paddingBottom: 0 }}>
        <div>
          <div className="crumbs"><a href="#/board">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Settings</span></div>
          <h1>Project settings</h1>
        </div>
      </div>
      <div className="tabs">
        {[["general","General"],["workflow","Workflow"],["components","Components"],["versions","Versions"],["custom-fields","Custom fields"],["labels","Labels"],["webhooks","Webhooks"],["members","Members"],["notifications","Notifications"],["integrations","Integrations"],["permissions","Roles & permissions"],["audit","Audit log"]].map(([id, label]) => (
          <button key={id} className="tab" aria-selected={tab === id} onClick={() => { if (id === "integrations") nav("integrations"); else if (id === "notifications") nav("notifications-page"); else if (id === "members") nav("members"); else setTab(id); }}>{label}</button>
        ))}
      </div>

      <div style={{ padding: "24px 32px", maxWidth: 1040 }}>
        {tab === "workflow" && <WorkflowTab/>}
        {tab === "components" && <ComponentsTab/>}
        {tab === "versions" && <VersionsTab/>}
        {tab === "custom-fields" && <CustomFieldsTab/>}
        {tab === "labels" && <LabelsTab/>}
        {tab === "webhooks" && <WebhooksTab/>}
        {tab === "general" && (
          <div className="stack gap-4">
            <div className="card card-pad">
              <h3 style={{ margin: "0 0 12px" }}>About this project</h3>
              <div className="stack gap-3">
                <div className="row gap-3" style={{ alignItems: "flex-start" }}>
                  <div style={{ width: 64, height: 64, borderRadius: 12, background: "var(--indigo-600)", color: "#fff", display: "grid", placeItems: "center" }}>
                    <Icon name="server" size={28}/>
                  </div>
                  <div className="stack gap-2" style={{ flex: 1 }}>
                    <div><label className="label">Name</label><input className="input" defaultValue="Core Infrastructure"/></div>
                    <div className="row gap-3">
                      <div style={{ flex: 1 }}><label className="label">Key</label><input className="input mono" defaultValue="INFRA"/></div>
                      <div style={{ flex: 2 }}><label className="label">Lead</label>
                        <PillSelect value="u1" onChange={() => {}} options={FORGE_DATA.PEOPLE.map((u) => ({ id: u.id, label: u.name }))}/>
                      </div>
                    </div>
                    <div><label className="label">Description</label><textarea className="textarea" defaultValue="Kubernetes clusters, networking, baseline infra."/></div>
                  </div>
                </div>
              </div>
            </div>

            <div className="card card-pad">
              <h3 style={{ margin: "0 0 12px" }}>Workflow</h3>
              <div className="text-sm secondary" style={{ marginBottom: 12 }}>The columns that appear on the Kanban board.</div>
              <div className="row gap-2" style={{ flexWrap: "wrap" }}>
                {FORGE_DATA.COLUMNS.map((c) => (
                  <div key={c.id} className="row gap-2" style={{ padding: "6px 12px", border: "1px solid var(--border)", borderRadius: 6, background: "var(--bg-subtle)" }}>
                    <span style={{ width: 8, height: 8, borderRadius: 2, background: { Backlog: "#94A3B8", Todo: "#64748B", "In Progress": "#3B82F6", "In Review": "#8B5CF6", Done: "#10B981" }[c.id] }}/>
                    {c.label}
                  </div>
                ))}
                <Button data-size="sm" icon="plus">Add column</Button>
              </div>
            </div>

            <div className="card card-pad" style={{ borderColor: "var(--danger)", background: "var(--danger-bg)" }}>
              <h3 style={{ margin: "0 0 4px", color: "#B91C1C" }}>Danger zone</h3>
              <p className="text-sm secondary" style={{ margin: "0 0 12px" }}>Archive or delete this project. This cannot be undone.</p>
              <div className="row gap-2">
                <Button>Archive project</Button>
                <Button variant="danger" icon="trash">Delete project</Button>
              </div>
            </div>
          </div>
        )}

        {tab === "permissions" && (
          <PermissionsMatrix/>
        )}

        {tab === "audit" && (
          <div className="card" style={{ overflow: "hidden" }}>
            <table className="table">
              <thead><tr><th>When</th><th>Who</th><th>Event</th><th>Target</th><th>IP</th></tr></thead>
              <tbody>
                {[
                  ["12 min ago","u1","members.invite","aisha@forge.dev","10.0.4.21"],
                  ["1h ago","u2","sprint.complete","Sprint 23","10.0.4.34"],
                  ["2h ago","u1","integrations.telegram.connect","@forge_team_bot","10.0.4.21"],
                  ["yesterday","u3","members.role.change","Hana Suzuki → Developer","10.0.4.18"],
                  ["yesterday","u1","project.settings.update","Core Infrastructure","10.0.4.21"],
                ].map(([when, w, ev, tgt, ip], i) => {
                  const u = FORGE_DATA.PEOPLE.find((p) => p.id === w);
                  return (
                    <tr key={i}>
                      <td className="text-xs muted">{when}</td>
                      <td><span className="row gap-2"><Avatar user={u} size="sm"/>{u.name}</span></td>
                      <td><span className="mono text-xs">{ev}</span></td>
                      <td>{tgt}</td>
                      <td className="mono text-xs muted">{ip}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

function PermissionsMatrix() {
  const perms = [
    "View board & issues",
    "Create & edit issues",
    "Delete issues",
    "Manage sprints",
    "Edit wiki pages",
    "Invite members",
    "Manage roles",
    "Manage integrations",
    "Delete project",
  ];
  const roles = ["Admin","Manager","Developer","Viewer"];
  const grid = [
    [1,1,1,1],
    [1,1,1,0],
    [1,1,0,0],
    [1,1,0,0],
    [1,1,1,0],
    [1,1,0,0],
    [1,0,0,0],
    [1,1,0,0],
    [1,0,0,0],
  ];
  return (
    <div className="card" style={{ overflow: "hidden" }}>
      <table className="table">
        <thead>
          <tr>
            <th>Permission</th>
            {roles.map((r) => <th key={r} style={{ textAlign: "center" }}>{r}</th>)}
          </tr>
        </thead>
        <tbody>
          {perms.map((p, i) => (
            <tr key={p}>
              <td>{p}</td>
              {grid[i].map((v, j) => (
                <td key={j} style={{ textAlign: "center" }}>
                  {v ? <Icon name="check" size={16} color="var(--success)" strokeWidth="2.5"/> : <Icon name="x" size={14} color="var(--text-muted)"/>}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

Object.assign(window, { WikiView, ProjectsView, SprintView, ReportsView, SettingsView });
