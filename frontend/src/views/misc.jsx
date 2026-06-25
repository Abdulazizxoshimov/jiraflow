import { useState, useEffect } from 'react';
import { Icon } from '../components/icons';
import { Avatar, AvatarStack, Badge, Button, Modal, StatusBadge, PriorityBadge, TypeIcon, Empty, Skeleton, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { PillSelect } from './board';
import { WorkflowTab, ComponentsTab, VersionsTab, CustomFieldsTab, LabelsTab, WebhooksTab, SprintCapacity, AutomationTab } from '../panels/settings';
import { api, apiUpload, useApi } from '../api/api';

// ─── Projects list ──────────────────────────────────────
export function ProjectsView({ nav }) {
  const { projects, people, me, switchProject, projectsLoading } = useApp();
  const toast = useToast();
  const [view, setView] = useState("grid");
  const [openCreate, setOpenCreate] = useState(false);
  const [newProj, setNewProj] = useState({ name: "", key: "", desc: "", lead: "", template: "Scrum" });
  const [creating, setCreating] = useState(false);

  function autoKey(name) {
    return name.toUpperCase().replace(/[^A-Z]/g, "").slice(0, 5) || "PROJ";
  }
  async function createProject() {
    if (!newProj.name.trim()) return;
    setCreating(true);
    try {
      const created = await api("/projects", {
        method: "POST",
        body: {
          name: newProj.name.trim(),
          key: (newProj.key.trim() || autoKey(newProj.name)).toUpperCase(),
          description: newProj.desc,
          lead_id: newProj.lead || (me ? me._raw?.id || me.id : ""),
        },
      });
      setOpenCreate(false);
      setNewProj({ name: "", key: "", desc: "", lead: "", template: "Scrum" });
      if (created.id) { await switchProject(created.id); nav("board"); }
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setCreating(false); }
  }

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">Forge <Icon name="chevronRight" size={11}/> <span>All projects</span></div>
          <h1>Projects</h1>
          <p>{projectsLoading ? "Loading…" : `${projects.length} projects · across the workspace.`}</p>
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
        {projectsLoading ? (
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 14 }}>
            {[1, 2, 3, 4, 5, 6].map((i) => (
              <div key={i} className="card" style={{ padding: 0, overflow: "hidden" }}>
                <Skeleton w="100%" h={64} radius="0"/>
                <div style={{ padding: 16, display: "flex", flexDirection: "column", gap: 10 }}>
                  <div className="row gap-2"><Skeleton w={120} h={14}/><Skeleton w={40} h={12}/></div>
                  <Skeleton w="85%" h={12}/>
                  <Skeleton w="65%" h={12}/>
                  <div className="row" style={{ justifyContent: "space-between", marginTop: 4 }}>
                    <Skeleton w={60} h={12}/>
                    <Skeleton w={24} h={24} radius="50%"/>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : !projects.length ? (
          <div style={{ padding: "60px 0" }}>
            <Empty icon="grid" title="No projects yet" hint="Create your first project to get started."
              action={
                <button className="btn btn-primary" onClick={() => setOpenCreate(true)}>
                  <Icon name="plus" size={14}/> Create project
                </button>
              }
            />
          </div>
        ) : view === "grid" ? (
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(280px, 1fr))", gap: 14 }}>
            {projects.map((p) => {
              const lead = p.leadUser || people.find((u) => u.id === p.lead);
              return (
                <div key={p.id} className="card" style={{ padding: 0, overflow: "hidden", cursor: "default" }}
                  onMouseEnter={(e) => e.currentTarget.style.borderColor = "var(--border-strong)"}
                  onMouseLeave={(e) => e.currentTarget.style.borderColor = ""}
                  onClick={() => { switchProject(p.id); nav("board"); }}>
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
                {projects.map((p) => {
                  const lead = p.leadUser || people.find((u) => u.id === p.lead);
                  return (
                    <tr key={p.id} onClick={() => { switchProject(p.id); nav("board"); }}>
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
        footer={
          <>
            <Button onClick={() => setOpenCreate(false)}>Cancel</Button>
            <Button variant="primary" disabled={!newProj.name.trim() || creating} onClick={createProject}>
              {creating ? "Creating…" : "Create project"}
            </Button>
          </>
        }>
        <div className="stack gap-3">
          <div>
            <label className="label">Project name</label>
            <input className="input" autoFocus placeholder="e.g. Platform Migrations" value={newProj.name}
              onChange={(e) => setNewProj({ ...newProj, name: e.target.value, key: autoKey(e.target.value) })}
              onKeyDown={(e) => { if (e.key === "Enter") createProject(); }}/>
          </div>
          <div className="row gap-3">
            <div style={{ flex: 1 }}>
              <label className="label">Key</label>
              <input className="input mono" placeholder="PLAT" maxLength="5" value={newProj.key}
                onChange={(e) => setNewProj({ ...newProj, key: e.target.value.toUpperCase().replace(/[^A-Z]/g, "").slice(0, 5) })}/>
            </div>
            <div style={{ flex: 2 }}>
              <label className="label">Lead</label>
              <PillSelect value={newProj.lead} onChange={(v) => setNewProj({ ...newProj, lead: v })}
                options={[{ id: "", label: "— Assign to me" }, ...people.map((p) => ({ id: p.id, label: p.name }))]}/>
            </div>
          </div>
          <div>
            <label className="label">Description</label>
            <textarea className="textarea" placeholder="What does this team own?" value={newProj.desc}
              onChange={(e) => setNewProj({ ...newProj, desc: e.target.value })}/>
          </div>
          <div>
            <label className="label">Template</label>
            <div className="row gap-2" style={{ flexWrap: "wrap" }}>
              {["Kanban","Scrum","Bug tracking","Empty"].map((t) => (
                <button key={t} className="btn" onClick={() => setNewProj({ ...newProj, template: t })}
                  style={{ border: "1px solid " + (newProj.template === t ? "var(--indigo-600)" : "var(--border)"), background: newProj.template === t ? "var(--indigo-50)" : "var(--bg)", color: newProj.template === t ? "var(--indigo-700)" : "var(--text)" }}>
                  {t}
                </button>
              ))}
            </div>
          </div>
        </div>
      </Modal>
    </div>
  );
}

// ─── Sprint planning ────────────────────────────────────
export function SprintView({ nav }) {
  const { issues, people, activeProjectId, projects } = useApp();
  const proj = projects.find((p) => p.id === activeProjectId);
  const projName = proj ? proj.name : "Project";
  const { data: sprintsData } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/sprints" : null, [activeProjectId]);
  const activeSprint = (sprintsData?.items || sprintsData || []).find((s) => s.status === "active");
  const activeSprintId = activeSprint?.id || null;

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">{projName}</a> <Icon name="chevronRight" size={11}/> <span>Sprints</span></div>
          <h1>Sprints</h1>
          <p>Plan, run, and review iterations.</p>
        </div>
        <div className="row gap-2">
          <Button icon="chart">Velocity report</Button>
          <Button variant="primary" icon="plus">Start sprint</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px" }}>
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

        <SprintCapacity sprintId={activeSprintId}/>

        <div className="card" style={{ marginBottom: 16 }}>
          <div className="card-head">
            <div className="row gap-3">
              <Badge tone="info" dot>Active</Badge>
              <h3>Sprint 24 — Edge &amp; Reliability</h3>
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
                const u = i.assigneeUser || people.find((p) => p.id === i.assignee);
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

        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
          <div className="card">
            <div className="card-head">
              <div className="row gap-3"><Badge tone="muted" dot>Planned</Badge><h3>Sprint 25 — Observability deep clean</h3></div>
              <Button data-size="sm">Start</Button>
            </div>
            <div style={{ padding: 16 }} className="text-sm secondary">23 issues · 52 points committed · Dec 16 → Dec 29</div>
          </div>
          <div className="card">
            <div className="card-head">
              <div className="row gap-3"><Badge tone="success" dot>Completed</Badge><h3>Sprint 23 — IAM cleanup</h3></div>
              <Button data-size="sm" icon="chart">Report</Button>
            </div>
            <div style={{ padding: 16 }} className="text-sm secondary">19 / 21 issues delivered · 47 / 51 points · Nov 18 → Dec 1</div>
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
export function ReportsView({ nav }) {
  const { issues, people, activeProjectId, projects } = useApp();
  const proj = projects.find((p) => p.id === activeProjectId);
  const projName = proj ? proj.name : "Project";
  const [reportTab, setReportTab] = useState("overview");

  // Fetch active sprint for burndown
  const { data: sprintsData } = useApi(activeProjectId ? `/projects/${activeProjectId}/sprints?status=active&limit=1` : null, [activeProjectId]);
  const activeSprint = (sprintsData?.items || sprintsData || [])[0] || null;

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">{projName}</a> <Icon name="chevronRight" size={11}/> <span>Reports</span></div>
          <h1>Reports &amp; analytics</h1>
          <p>Insights across velocity, workload and throughput.</p>
        </div>
        <div className="row gap-2">
          <Button icon="calendar">Last 30 days</Button>
          <Button icon="download">Export PDF</Button>
        </div>
      </div>

      <div className="tabs">
        {[["overview","Overview"],["gantt","Gantt chart"],["cfd","Cumulative flow"]].map(([id, label]) => (
          <button key={id} className="tab" aria-selected={reportTab === id} onClick={() => setReportTab(id)}>{label}</button>
        ))}
      </div>

      {reportTab === "gantt" && <GanttView projectId={activeProjectId} issues={issues}/>}
      {reportTab === "cfd"   && <CfdView projectId={activeProjectId}/>}

      {reportTab === "overview" && <div style={{ padding: "0 32px 32px", display: "grid", gridTemplateColumns: "repeat(12, 1fr)", gap: 16 }}>
        <div className="card" style={{ gridColumn: "span 8" }}>
          <div className="card-head">
            <h3>Burndown{activeSprint ? ` — ${activeSprint.name}` : ''}</h3>
          </div>
          <div style={{ padding: 16, height: 320 }}>
            <Burndown sprintId={activeSprint?.id} />
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
              {[["Story",38,"#8B5CF6"],["Task",32,"#3B82F6"],["Bug",22,"#EF4444"],["Epic",8,"#F97316"]].map(([l,v,c]) => (
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
            {people.slice(0, 7).map((u) => {
              const pts = issues.filter((i) => i.assignee === u.id && i.sprint).reduce((s, i) => s + (i.points || 0), 0);
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
              ["Critical", issues.filter((i) => i.pri === "Critical").length, "#DC2626"],
              ["High",     issues.filter((i) => i.pri === "High").length,     "#EA580C"],
              ["Medium",   issues.filter((i) => i.pri === "Medium").length,   "#3B82F6"],
              ["Low",      issues.filter((i) => i.pri === "Low").length,      "#64748B"],
            ].map(([l, v, c]) => {
              const max = Math.max(...[issues.filter((i) => i.pri === "Critical").length, issues.filter((i) => i.pri === "High").length, issues.filter((i) => i.pri === "Medium").length, 1]);
              return (
                <div key={l} className="row gap-3" style={{ marginBottom: 10 }}>
                  <span className="text-sm" style={{ width: 80, color: c, fontWeight: 500 }}>{l}</span>
                  <div style={{ flex: 1, height: 22, background: "var(--bg-subtle)", borderRadius: 4, overflow: "hidden", position: "relative" }}>
                    <div style={{ width: (v / max) * 100 + "%", height: "100%", background: c, opacity: .8 }}/>
                  </div>
                  <span className="bold text-sm" style={{ width: 24, textAlign: "right" }}>{v}</span>
                </div>
              );
            })}
          </div>
        </div>

        <div className="card" style={{ gridColumn: "span 12" }}>
          <div className="card-head">
            <h3>Velocity (story points per sprint)</h3>
          </div>
          <div style={{ padding: 16 }}>
            <Throughput projectId={activeProjectId}/>
          </div>
        </div>
      </div>}
    </div>
  );
}

// ─── Gantt chart view ─────────────────────────────────
function GanttView({ projectId, issues }) {
  const { data: gantt, loading, error } = useApi(projectId ? "/projects/" + projectId + "/gantt" : null, [projectId]);

  const rows = gantt?.items || gantt || [];
  const fallback = issues.filter((i) => i.due).slice(0, 20).map((i) => ({
    id: i.id, title: i.title, type: i.type,
    start_date: i.due, end_date: i.due, status: i.status,
  }));
  const data = rows.length ? rows : fallback;

  const today = new Date();
  const startMs = data.length
    ? Math.min(...data.map((r) => new Date(r.start_date || r.end_date).getTime()))
    : today.getTime() - 7 * 86400000;
  const endMs = data.length
    ? Math.max(...data.map((r) => new Date(r.end_date || r.start_date).getTime()))
    : today.getTime() + 14 * 86400000;
  const totalDays = Math.max(1, Math.round((endMs - startMs) / 86400000)) + 2;

  function dayOffset(date) {
    return Math.round((new Date(date).getTime() - startMs) / 86400000);
  }

  const STATUS_COLORS = { done: "#10B981", "in-progress": "#6366F1", todo: "#94A3B8" };

  if (loading) return <div className="text-sm muted" style={{ padding: 40, textAlign: "center" }}>Loading Gantt…</div>;
  if (error) return <div className="text-sm" style={{ padding: 40, textAlign: "center", color: "var(--danger)" }}>{error}</div>;
  if (!data.length) return (
    <div style={{ padding: 40 }}>
      <Empty icon="calendar" title="No Gantt data" hint="Issues with due dates will appear here."/>
    </div>
  );

  const todayOffset = dayOffset(today);

  return (
    <div style={{ padding: "0 32px 32px", overflowX: "auto" }}>
      <div style={{ minWidth: Math.max(800, totalDays * 28 + 220) }}>
        {/* Header — day labels */}
        <div className="row" style={{ marginBottom: 2 }}>
          <div style={{ width: 220, flexShrink: 0 }}/>
          <div style={{ flex: 1, position: "relative", height: 24 }}>
            {Array.from({ length: Math.ceil(totalDays / 7) }).map((_, wi) => {
              const d = new Date(startMs + wi * 7 * 86400000);
              return (
                <div key={wi} style={{ position: "absolute", left: wi * 7 * 28, fontSize: 10, color: "var(--text-muted)", fontWeight: 500 }}>
                  {d.toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                </div>
              );
            })}
          </div>
        </div>
        {/* Rows */}
        {data.map((row) => {
          const left = Math.max(0, dayOffset(row.start_date || row.end_date)) * 28;
          const width = Math.max(28, (dayOffset(row.end_date || row.start_date) - dayOffset(row.start_date || row.end_date) + 1) * 28);
          const color = STATUS_COLORS[row.status] || "#6366F1";
          return (
            <div key={row.id} className="row" style={{ marginBottom: 4, alignItems: "center" }}>
              <div style={{ width: 220, flexShrink: 0, paddingRight: 12, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                <span className="mono text-xs muted" style={{ marginRight: 6 }}>{row.id}</span>
                <span className="text-sm">{row.title}</span>
              </div>
              <div style={{ flex: 1, position: "relative", height: 26, background: "var(--bg-subtle)", borderRadius: 4 }}>
                {todayOffset >= 0 && (
                  <div style={{ position: "absolute", left: todayOffset * 28, top: 0, bottom: 0, width: 2, background: "#EF4444", zIndex: 2 }}/>
                )}
                <div style={{ position: "absolute", left, top: 4, height: 18, width, background: color, borderRadius: 4, opacity: 0.85, display: "flex", alignItems: "center", paddingLeft: 6 }}>
                  <span style={{ fontSize: 10, color: "#fff", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{row.title}</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>
      <div className="row gap-4 text-xs muted" style={{ marginTop: 12 }}>
        <span className="row gap-2"><span style={{ width: 12, height: 12, borderRadius: 3, background: "#10B981" }}/> Done</span>
        <span className="row gap-2"><span style={{ width: 12, height: 12, borderRadius: 3, background: "#6366F1" }}/> In progress</span>
        <span className="row gap-2"><span style={{ width: 12, height: 12, borderRadius: 3, background: "#94A3B8" }}/> To do</span>
        <span className="row gap-2"><span style={{ width: 2, height: 12, background: "#EF4444" }}/> Today</span>
      </div>
    </div>
  );
}

// ─── CFD — Cumulative Flow Diagram ─────────────────────
function CfdView({ projectId }) {
  const { data: cfd, loading, error } = useApi(projectId ? "/projects/" + projectId + "/cfd" : null, [projectId]);

  const COLORS = ["#10B981", "#6366F1", "#F59E0B", "#94A3B8"];
  const fallbackData = {
    columns: ["Done", "In Review", "In Progress", "To Do"],
    rows: [
      { date: "W1",  values: [2, 1, 3, 8] },
      { date: "W2",  values: [4, 2, 3, 7] },
      { date: "W3",  values: [7, 2, 4, 6] },
      { date: "W4",  values: [10, 3, 4, 5] },
      { date: "W5",  values: [13, 3, 5, 4] },
      { date: "W6",  values: [16, 2, 5, 3] },
      { date: "W7",  values: [19, 3, 4, 2] },
      { date: "W8",  values: [22, 2, 4, 2] },
    ],
  };

  const chartData = (cfd && cfd.rows && cfd.rows.length) ? cfd : fallbackData;
  const { columns: cols, rows } = chartData;

  const w = 720, h = 240, padL = 40, padR = 20, padT = 20, padB = 30;
  const maxTotal = Math.max(...rows.map((r) => r.values.reduce((s, v) => s + v, 0)));
  const xStep = (w - padL - padR) / Math.max(1, rows.length - 1);
  const yScale = (v) => h - padB - (v / maxTotal) * (h - padT - padB);

  function getStack(rowIdx) {
    const stack = [];
    let acc = 0;
    for (let ci = cols.length - 1; ci >= 0; ci--) {
      acc += rows[rowIdx].values[ci];
      stack.push(acc);
    }
    return stack.reverse();
  }

  if (loading) return <div className="text-sm muted" style={{ padding: 40, textAlign: "center" }}>Loading CFD…</div>;

  return (
    <div style={{ padding: "0 32px 32px" }}>
      <div className="card">
        <div className="card-head">
          <h3>Cumulative Flow Diagram</h3>
          <div className="row gap-4 text-xs muted">
            {cols.map((c, i) => (
              <span key={c} className="row gap-2"><span style={{ width: 12, height: 12, borderRadius: 3, background: COLORS[i % COLORS.length] }}/>{c}</span>
            ))}
          </div>
        </div>
        <div style={{ padding: 16 }}>
          <svg viewBox={`0 0 ${w} ${h}`} style={{ width: "100%", height: 240 }}>
            {/* Y axis */}
            {[0, 25, 50, 75, 100].map((pct) => {
              const v = maxTotal * pct / 100;
              return (
                <g key={pct}>
                  <line x1={padL} y1={yScale(v)} x2={w - padR} y2={yScale(v)} stroke="var(--border)" strokeDasharray="2 4"/>
                  <text x={padL - 6} y={yScale(v) + 4} fontSize="9" fill="var(--text-muted)" textAnchor="end">{Math.round(v)}</text>
                </g>
              );
            })}
            {/* X labels */}
            {rows.map((r, i) => (
              <text key={i} x={padL + i * xStep} y={h - padB + 14} fontSize="9" fill="var(--text-muted)" textAnchor="middle">{r.date}</text>
            ))}
            {/* Stacked area */}
            {cols.map((col, ci) => {
              const upperVals = rows.map((_, ri) => getStack(ri)[ci]);
              const lowerVals = ci === 0 ? rows.map(() => 0) : rows.map((_, ri) => getStack(ri)[ci - 1]);
              const upper = upperVals.map((v, i) => `${i === 0 ? "M" : "L"} ${padL + i * xStep} ${yScale(v)}`).join(" ");
              const lower = lowerVals.map((v, i) => `L ${padL + (rows.length - 1 - i) * xStep} ${yScale(v)}`).join(" ");
              return (
                <path key={col} d={`${upper} ${lower} Z`} fill={COLORS[ci % COLORS.length]} fillOpacity="0.7"/>
              );
            })}
          </svg>
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

function Burndown({ sprintId }) {
  const { data: bd, loading } = useApi(sprintId ? `/sprints/${sprintId}/burndown` : null, [sprintId]);

  const days   = bd?.days   || [];
  const ideal  = days.map((d) => d.ideal  ?? d.remaining_ideal ?? 0);
  const actual = days.map((d) => d.actual ?? d.remaining        ?? 0);

  if (loading) return <div className="text-sm muted" style={{ paddingTop: 80, textAlign: "center" }}>Loading burndown…</div>;
  if (!days.length) return <div className="text-sm muted" style={{ paddingTop: 80, textAlign: "center" }}>No active sprint data</div>;

  const w = 800, h = 260, padL = 44, padR = 16, padT = 16, padB = 30;
  const max = Math.max(...ideal, ...actual, 1);
  const stepX = (w - padL - padR) / Math.max(1, ideal.length - 1);
  const y = (v) => h - padB - (v / max) * (h - padT - padB);
  const yTicks = [0, 0.25, 0.5, 0.75, 1].map((f) => Math.round(f * max));

  return (
    <svg viewBox={`0 0 ${w} ${h}`} style={{ width: "100%", height: "100%" }}>
      {yTicks.map((v) => (
        <g key={v}>
          <line x1={padL} y1={y(v)} x2={w - padR} y2={y(v)} stroke="var(--border)" strokeDasharray="2 4"/>
          <text x={padL - 8} y={y(v) + 4} fontSize="10" fill="var(--text-muted)" textAnchor="end">{v}</text>
        </g>
      ))}
      {ideal.map((_, i) => i % 2 === 0 && (
        <text key={i} x={padL + i * stepX} y={h - padB + 16} fontSize="10" fill="var(--text-muted)" textAnchor="middle">D{i + 1}</text>
      ))}
      <path d={ideal.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ")} fill="none" stroke="var(--text-muted)" strokeDasharray="4 4"/>
      {actual.length > 1 && <>
        <path d={actual.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ") + ` L ${padL + (actual.length - 1) * stepX} ${h - padB} L ${padL} ${h - padB} Z`} fill="#6366F1" fillOpacity=".12"/>
        <path d={actual.map((v, i) => `${i ? "L" : "M"} ${padL + i * stepX} ${y(v)}`).join(" ")} fill="none" stroke="#6366F1" strokeWidth="2.5"/>
        {actual.map((v, i) => <circle key={i} cx={padL + i * stepX} cy={y(v)} r="3" fill="#6366F1"/>)}
      </>}
      <text x={w - padR} y={padT + 4} fontSize="11" fill="var(--text-muted)" textAnchor="end">— — Ideal</text>
      <text x={w - padR} y={padT + 18} fontSize="11" fill="#6366F1" textAnchor="end">Actual</text>
    </svg>
  );
}

function Throughput({ projectId }) {
  const { data: vel, loading } = useApi(projectId ? `/projects/${projectId}/velocity` : null, [projectId]);
  const sprints = vel?.sprints || vel?.items || [];

  if (loading) return <div className="text-sm muted" style={{ height: 140, display: "flex", alignItems: "center", justifyContent: "center" }}>Loading velocity…</div>;
  if (!sprints.length) return <div className="text-sm muted" style={{ height: 140, display: "flex", alignItems: "center", justifyContent: "center" }}>No sprint data yet</div>;

  const maxPts = Math.max(...sprints.map((s) => s.completed_points ?? s.completed ?? 0), 1);
  return (
    <div className="row" style={{ alignItems: "flex-end", gap: 6, height: 140, padding: "0 8px" }}>
      {sprints.slice(-12).map((s, i) => {
        const v = s.completed_points ?? s.completed ?? 0;
        return (
          <div key={i} style={{ flex: 1, display: "flex", flexDirection: "column", alignItems: "center" }}>
            <div style={{ width: "70%", height: (v / maxPts) * 110, background: "linear-gradient(180deg, #818CF8, #6366F1)", borderRadius: "3px 3px 0 0" }}/>
            <div className="text-xs muted" style={{ marginTop: 4, fontSize: 10 }} title={s.name}>S{i + 1}</div>
          </div>
        );
      })}
    </div>
  );
}

// ─── Settings ───────────────────────────────────────────
export function SettingsView({ nav }) {
  const { me, people, activeProjectId, projects, columns } = useApp();
  const toast = useToast();
  const [tab, setTab] = useState("general");
  const proj = projects.find((p) => p.id === activeProjectId);
  const projName = proj ? proj.name : "Project";
  const [sName, setSName] = useState(projName);
  const [sDesc, setSDesc] = useState(proj ? proj.desc : "");
  const [sLead, setSLead] = useState(proj ? proj.lead : me ? me.id : "");
  const [savingSettings, setSavingSettings] = useState(false);

  useEffect(() => {
    setSName(proj ? proj.name : "");
    setSDesc(proj ? proj.desc : "");
    setSLead(proj ? proj.lead : me ? me.id : "");
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeProjectId]);

  async function saveSettings() {
    if (!proj) return;
    setSavingSettings(true);
    try {
      await api("/projects/" + (proj._raw?.id || proj.id), { method: "PUT", body: { name: sName, description: sDesc, lead_id: sLead } });
      toast("Settings saved");
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSavingSettings(false); }
  }

  return (
    <div>
      <div className="page-head" style={{ paddingBottom: 0 }}>
        <div>
          <div className="crumbs"><a href="#/board">{projName}</a> <Icon name="chevronRight" size={11}/> <span>Settings</span></div>
          <h1>Project settings</h1>
        </div>
      </div>
      <div className="tabs">
        {[["general","General"],["workflow","Workflow"],["components","Components"],["versions","Versions"],["custom-fields","Custom fields"],["labels","Labels"],["webhooks","Webhooks"],["automation","Automation"],["members","Members"],["notifications","Notifications"],["integrations","Integrations"],["permissions","Roles & permissions"],["permission-schemes","Permission Schemes"],["issue-type-schemes","Issue Type Schemes"],["field-configurations","Field Config"],["notification-schemes","Notification Schemes"],["audit","Audit log"]].map(([id, label]) => (
          <button key={id} className="tab" aria-selected={tab === id} onClick={() => {
            if (id === "integrations") nav("integrations");
            else if (id === "notifications") nav("notifications-page");
            else if (id === "members") nav("members");
            else setTab(id);
          }}>{label}</button>
        ))}
      </div>

      <div style={{ padding: "24px 32px", maxWidth: 1040 }}>
        {tab === "workflow"      && <WorkflowTab/>}
        {tab === "components"   && <ComponentsTab/>}
        {tab === "versions"     && <VersionsTab/>}
        {tab === "custom-fields"&& <CustomFieldsTab/>}
        {tab === "labels"       && <LabelsTab/>}
        {tab === "webhooks"     && <WebhooksTab/>}
        {tab === "automation"            && <AutomationTab/>}
        {tab === "permission-schemes"    && <PermissionSchemesTab/>}
        {tab === "issue-type-schemes"    && <IssueTypeSchemesTab/>}
        {tab === "field-configurations"  && <FieldConfigurationsTab/>}
        {tab === "notification-schemes"  && <NotificationSchemesTab/>}

        {tab === "general" && (
          <div className="stack gap-4">
            <div className="card card-pad">
              <h3 style={{ margin: "0 0 12px" }}>About this project</h3>
              <div className="stack gap-3">
                <div className="row gap-3" style={{ alignItems: "flex-start" }}>
                  <div style={{ width: 64, height: 64, borderRadius: 12, background: proj ? proj.color : "var(--indigo-600)", color: "#fff", display: "grid", placeItems: "center" }}>
                    <Icon name={proj ? proj.icon : "briefcase"} size={28}/>
                  </div>
                  <div className="stack gap-2" style={{ flex: 1 }}>
                    <div><label className="label">Name</label><input className="input" value={sName} onChange={(e) => setSName(e.target.value)}/></div>
                    <div className="row gap-3">
                      <div style={{ flex: 1 }}><label className="label">Key</label><input className="input mono" defaultValue={proj ? proj.key : ""} disabled/></div>
                      <div style={{ flex: 2 }}><label className="label">Lead</label>
                        <PillSelect value={sLead} onChange={setSLead} options={people.map((u) => ({ id: u.id, label: u.name }))}/>
                      </div>
                    </div>
                    <div><label className="label">Description</label><textarea className="textarea" value={sDesc} onChange={(e) => setSDesc(e.target.value)}/></div>
                    <div className="row" style={{ justifyContent: "flex-end", paddingTop: 4 }}>
                      <Button variant="primary" disabled={savingSettings} onClick={saveSettings}>{savingSettings ? "Saving…" : "Save changes"}</Button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="card card-pad">
              <h3 style={{ margin: "0 0 12px" }}>Workflow</h3>
              <div className="text-sm secondary" style={{ marginBottom: 12 }}>The columns that appear on the Kanban board.</div>
              <div className="row gap-2" style={{ flexWrap: "wrap" }}>
                {columns.map((c) => {
                  const TONE_COLOR = { muted: "#94A3B8", info: "#3B82F6", success: "#10B981", warning: "#F59E0B", danger: "#EF4444", purple: "#8B5CF6" };
                  return (
                    <div key={c.id} className="row gap-2" style={{ padding: "6px 12px", border: "1px solid var(--border)", borderRadius: 6, background: "var(--bg-subtle)" }}>
                      <span style={{ width: 8, height: 8, borderRadius: 2, background: TONE_COLOR[c.tone] || "#94A3B8" }}/>
                      {c.label}
                    </div>
                  );
                })}
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

        {tab === "permissions" && <PermissionsMatrix/>}

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
                  ["yesterday","u1","project.settings.update",projName,"10.0.4.21"],
                ].map(([when, w, ev, tgt, ip], i) => {
                  const u = people.find((p) => p.id === w) || { name: w, initials: "?", color: "#94A3B8" };
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
  const perms = ["View board & issues","Create & edit issues","Delete issues","Manage sprints","Edit wiki pages","Invite members","Manage roles","Manage integrations","Delete project"];
  const roles = ["Admin","Manager","Developer","Viewer"];
  const grid = [[1,1,1,1],[1,1,1,0],[1,1,0,0],[1,1,0,0],[1,1,1,0],[1,1,0,0],[1,0,0,0],[1,1,0,0],[1,0,0,0]];
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

// ─── Profile / Account settings ───────────────────────────────────────────────
export function ProfileView({ nav }) {
  const { me, clearSession } = useApp();
  const toast = useToast();

  const [name,        setName]        = useState(me?.name     || "");
  const [email,       setEmail]       = useState(me?.email    || "");
  const [avatarUrl,   setAvatarUrl]   = useState(me?.avatar   || null);
  const [pwdOpen,     setPwdOpen]     = useState(false);
  const [oldPwd,      setOldPwd]      = useState("");
  const [newPwd,      setNewPwd]      = useState("");
  const [saving,      setSaving]      = useState(false);
  const [avatarBusy,  setAvatarBusy]  = useState(false);
  const avatarInputRef = { current: null };

  const { data: meData } = useApi("/users/me");
  useEffect(() => {
    if (meData) {
      setName(meData.full_name || meData.name || "");
      setEmail(meData.email || "");
      setAvatarUrl(meData.avatar_url || null);
    }
  }, [meData]);

  async function saveProfile() {
    setSaving(true);
    try {
      await api("/users/me", { method: "PUT", body: { full_name: name } });
      toast("Profile updated");
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSaving(false); }
  }

  async function changePassword() {
    if (!oldPwd || !newPwd) return;
    try {
      await api("/users/me/password", { method: "PUT", body: { old_password: oldPwd, new_password: newPwd } });
      toast("Password changed");
      setPwdOpen(false); setOldPwd(""); setNewPwd("");
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function handleAvatarFile(file) {
    if (!file) return;
    setAvatarBusy(true);
    try {
      const updated = await apiUpload("/users/me/avatar", file);
      setAvatarUrl(updated.avatar_url || null);
      toast("Avatar updated");
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setAvatarBusy(false); }
  }

  const avatar = { ...(me || {}), name: name || "?", initials: (name || "?")[0], color: me?.color || "#6366F1", avatar: avatarUrl };

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><Icon name="user" size={11}/> <span>Profile</span></div>
          <h1>My account</h1>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", maxWidth: 680 }}>

        {/* Avatar + name card */}
        <div className="card card-pad" style={{ marginBottom: 16 }}>
          <div className="row gap-4" style={{ alignItems: "center", marginBottom: 20 }}>
            <div style={{ position: "relative", cursor: "pointer" }}
              onClick={() => !avatarBusy && avatarInputRef.current?.click()}>
              <Avatar user={avatar} size="xl" style={{ width: 72, height: 72, fontSize: 28 }}/>
              <div style={{
                position: "absolute", inset: 0, borderRadius: "50%",
                background: "rgba(0,0,0,.45)", display: "flex", alignItems: "center",
                justifyContent: "center", opacity: avatarBusy ? 1 : 0, transition: "opacity .15s",
              }} className="avatar-overlay">
                {avatarBusy ? <span style={{ width: 18, height: 18, borderRadius: "50%", border: "2px solid #fff", borderTopColor: "transparent", animation: "forge-spin .7s linear infinite", display: "block" }}/> : <Icon name="camera" size={18} color="#fff"/>}
              </div>
              <input ref={(el) => { avatarInputRef.current = el; }} type="file" accept="image/*"
                style={{ display: "none" }}
                onChange={(e) => { handleAvatarFile(e.target.files[0]); e.target.value = ""; }}/>
            </div>
            <div className="stack" style={{ flex: 1 }}>
              <span style={{ fontWeight: 700, fontSize: 20 }}>{name || "—"}</span>
              <span className="text-sm secondary">{email}</span>
              <Badge tone="info" style={{ alignSelf: "flex-start", marginTop: 4 }}>{me?.role || "member"}</Badge>
            </div>
          </div>

          <div className="stack gap-3">
            <div>
              <label className="label">Full name</label>
              <input className="input" value={name} onChange={(e) => setName(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && saveProfile()}/>
            </div>
            <div>
              <label className="label">Email</label>
              <input className="input" value={email} disabled style={{ opacity: 0.6 }}/>
              <div className="text-xs muted" style={{ marginTop: 4 }}>Email cannot be changed here</div>
            </div>
            <div className="row gap-2" style={{ justifyContent: "flex-end" }}>
              <Button onClick={saveProfile} variant="primary" disabled={saving}>{saving ? "Saving…" : "Save changes"}</Button>
            </div>
          </div>
        </div>

        {/* Password */}
        <div className="card card-pad" style={{ marginBottom: 16 }}>
          <div className="row" style={{ justifyContent: "space-between", alignItems: "center" }}>
            <div>
              <div className="bold" style={{ marginBottom: 2 }}>Password</div>
              <div className="text-sm secondary">Change your login password</div>
            </div>
            <Button icon="lock" onClick={() => setPwdOpen(true)}>Change password</Button>
          </div>
        </div>

        {/* Telegram status */}
        <div className="card card-pad" style={{ marginBottom: 16 }}>
          <div className="row" style={{ justifyContent: "space-between", alignItems: "center" }}>
            <div className="row gap-3">
              <div style={{ width: 40, height: 40, borderRadius: 10, background: "linear-gradient(135deg,#2AABEE,#229ED9)", color: "#fff", display: "grid", placeItems: "center" }}>
                <Icon name="telegram" size={20}/>
              </div>
              <div>
                <div className="bold" style={{ marginBottom: 2 }}>Telegram</div>
                <div className="text-sm secondary">Notifications via @jira_flowbot</div>
              </div>
            </div>
            <Button icon="telegram" onClick={() => nav("integrations")}>Manage</Button>
          </div>
        </div>

        {/* Danger zone */}
        <div className="card card-pad" style={{ borderColor: "var(--danger)", background: "var(--danger-bg)" }}>
          <div className="bold" style={{ color: "#B91C1C", marginBottom: 4 }}>Sign out</div>
          <div className="text-sm secondary" style={{ marginBottom: 12 }}>You will be logged out from this device.</div>
          <Button variant="danger" icon="exit" onClick={clearSession}>Sign out</Button>
        </div>
      </div>

      {/* Password change modal */}
      <Modal open={pwdOpen} onClose={() => { setPwdOpen(false); setOldPwd(""); setNewPwd(""); }} title="Change password"
        footer={<><Button onClick={() => setPwdOpen(false)}>Cancel</Button><Button variant="primary" onClick={changePassword} disabled={!oldPwd || newPwd.length < 6}>Change</Button></>}>
        <div className="stack gap-3">
          <div>
            <label className="label">Current password</label>
            <input className="input" type="password" value={oldPwd} onChange={(e) => setOldPwd(e.target.value)}/>
          </div>
          <div>
            <label className="label">New password</label>
            <input className="input" type="password" value={newPwd} onChange={(e) => setNewPwd(e.target.value)}
              placeholder="Min. 6 characters"/>
          </div>
        </div>
      </Modal>
    </div>
  );
}

// ─── Enterprise Settings Tabs ──────────────────────────────────────────────────

function EnterpriseSchemeTab({ title, description, apiPath, renderRow, renderForm }) {
  const { data, loading, reload } = useApi(apiPath);
  const [form, setForm] = useState(null);
  const [saving, setSaving] = useState(false);
  const toast = useToast();
  const items = data?.items || data || [];

  async function save() {
    if (!form?.name?.trim()) return;
    setSaving(true);
    try {
      if (form._id) {
        await api(`${apiPath}/${form._id}`, { method: "PUT", body: form });
        toast("Updated");
      } else {
        await api(apiPath, { body: form });
        toast("Created");
      }
      reload();
      setForm(null);
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSaving(false); }
  }

  async function remove(id) {
    if (!window.confirm("Delete this scheme?")) return;
    try { await api(`${apiPath}/${id}`, { method: "DELETE" }); reload(); toast("Deleted"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="stack gap-4">
      <div className="row" style={{ justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <h3 style={{ margin: 0 }}>{title}</h3>
          <p className="text-sm secondary" style={{ margin: "4px 0 0" }}>{description}</p>
        </div>
        <Button variant="primary" icon="plus" onClick={() => setForm({ name: "", description: "" })}>New scheme</Button>
      </div>

      {loading ? <div className="text-sm muted">Loading…</div> : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table className="table">
            <thead><tr><th>Name</th><th>Description</th><th style={{ width: 80 }}></th></tr></thead>
            <tbody>
              {items.length === 0 && (
                <tr><td colSpan={3} style={{ textAlign: "center", color: "var(--text-muted)", padding: 24 }}>No schemes yet</td></tr>
              )}
              {items.map((item) => (
                <tr key={item.id}>
                  <td style={{ fontWeight: 600 }}>{item.name}</td>
                  <td className="text-sm secondary">{item.description || "—"}</td>
                  <td>
                    <div className="row gap-1">
                      <button className="icon-btn" title="Edit" onClick={() => setForm({ ...item, _id: item.id })}><Icon name="pencil" size={13}/></button>
                      <button className="icon-btn" title="Delete" style={{ color: "var(--danger)" }} onClick={() => remove(item.id)}><Icon name="trash" size={13}/></button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal open={!!form} onClose={() => setForm(null)} title={form?._id ? "Edit scheme" : "New scheme"}
        footer={<><Button onClick={() => setForm(null)}>Cancel</Button><Button variant="primary" disabled={saving} onClick={save}>{saving ? "Saving…" : "Save"}</Button></>}>
        {form && (
          <div className="stack gap-3">
            <div><label className="label">Name *</label>
              <input className="input" value={form.name} onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))} autoFocus/>
            </div>
            <div><label className="label">Description</label>
              <textarea className="textarea" value={form.description || ""} onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))} rows={3}/>
            </div>
            {renderForm && renderForm(form, setForm)}
          </div>
        )}
      </Modal>
    </div>
  );
}

function PermissionSchemesTab() {
  return (
    <EnterpriseSchemeTab
      title="Permission Schemes"
      description="Define who can view, create, edit, and delete issues within projects."
      apiPath="/permission-schemes"
    />
  );
}

function IssueTypeSchemesTab() {
  const { data: typesData } = useApi("/issue-types");
  const issueTypes = typesData?.items || typesData || [];

  return (
    <EnterpriseSchemeTab
      title="Issue Type Schemes"
      description="Control which issue types are available in each project."
      apiPath="/issue-type-schemes"
      renderForm={(form, setForm) => (
        <div>
          <label className="label">Issue types in this scheme</label>
          <div className="stack gap-1" style={{ marginTop: 6 }}>
            {issueTypes.map((t) => {
              const selected = (form.issue_type_ids || []).includes(t.id);
              return (
                <label key={t.id} className="row gap-2" style={{ cursor: "pointer", fontSize: 13 }}>
                  <input type="checkbox" checked={selected} onChange={(e) => {
                    const ids = form.issue_type_ids || [];
                    setForm((f) => ({
                      ...f,
                      issue_type_ids: e.target.checked ? [...ids, t.id] : ids.filter((id) => id !== t.id)
                    }));
                  }}/>
                  {t.name}
                </label>
              );
            })}
          </div>
        </div>
      )}
    />
  );
}

function FieldConfigurationsTab() {
  const { data, loading, reload } = useApi("/field-configurations");
  const [selected, setSelected] = useState(null);
  const { data: fieldsData, reload: reloadFields } = useApi(
    selected ? `/field-configurations/${selected}/fields` : null, [selected]
  );
  const toast = useToast();
  const configs = data?.items || data || [];
  const fields = fieldsData?.items || fieldsData || [];

  async function toggleField(fieldId, hidden) {
    try {
      await api(`/field-configurations/${selected}/fields/${fieldId}`, { method: "PUT", body: { hidden: !hidden } });
      reloadFields();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function createConfig() {
    const name = window.prompt("Configuration name:");
    if (!name?.trim()) return;
    try { await api("/field-configurations", { body: { name } }); reload(); toast("Created"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="stack gap-4">
      <div className="row" style={{ justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <h3 style={{ margin: 0 }}>Field Configurations</h3>
          <p className="text-sm secondary" style={{ margin: "4px 0 0" }}>Control which fields are visible and required in each project.</p>
        </div>
        <Button variant="primary" icon="plus" onClick={createConfig}>New configuration</Button>
      </div>
      <div className="row gap-4" style={{ alignItems: "flex-start" }}>
        <div className="card" style={{ width: 240, overflow: "hidden", flexShrink: 0 }}>
          {loading ? <div className="text-sm muted" style={{ padding: 16 }}>Loading…</div> : (
            configs.map((c) => (
              <button key={c.id} className="nav-item" aria-current={selected === c.id ? "page" : undefined}
                style={{ width: "100%", textAlign: "left", padding: "10px 14px", fontSize: 13 }}
                onClick={() => setSelected(c.id)}>
                {c.name}
              </button>
            ))
          )}
        </div>
        {selected ? (
          <div className="card card-pad" style={{ flex: 1 }}>
            <h4 style={{ margin: "0 0 12px" }}>Fields</h4>
            <table className="table">
              <thead><tr><th>Field</th><th>Type</th><th>Visible</th></tr></thead>
              <tbody>
                {fields.map((f) => (
                  <tr key={f.id}>
                    <td style={{ fontWeight: 500 }}>{f.name}</td>
                    <td className="text-sm secondary">{f.field_type}</td>
                    <td>
                      <label className="row gap-2" style={{ cursor: "pointer" }}>
                        <input type="checkbox" checked={!f.hidden}
                          onChange={() => toggleField(f.id, f.hidden)}/>
                        <span className="text-sm">{f.hidden ? "Hidden" : "Visible"}</span>
                      </label>
                    </td>
                  </tr>
                ))}
                {fields.length === 0 && <tr><td colSpan={3} style={{ textAlign: "center", color: "var(--text-muted)", padding: 16 }}>No fields</td></tr>}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="card card-pad" style={{ flex: 1, textAlign: "center", color: "var(--text-muted)", padding: 40 }}>
            Select a configuration to manage fields
          </div>
        )}
      </div>
    </div>
  );
}

function NotificationSchemesTab() {
  const { data, loading, reload } = useApi("/notification-schemes");
  const [selected, setSelected] = useState(null);
  const { data: rulesData, reload: reloadRules } = useApi(
    selected ? `/notification-schemes/${selected}/rules` : null, [selected]
  );
  const toast = useToast();
  const schemes = data?.items || data || [];
  const rules = rulesData?.items || rulesData || [];

  const EVENT_TYPES = [
    "issue_created","issue_updated","issue_assigned","issue_commented",
    "issue_status_changed","sprint_started","sprint_completed","page_commented",
  ];

  async function createScheme() {
    const name = window.prompt("Scheme name:");
    if (!name?.trim()) return;
    try { await api("/notification-schemes", { body: { name } }); reload(); toast("Created"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function addRule(eventType) {
    try {
      await api(`/notification-schemes/${selected}/rules`, { body: { event_type: eventType, notify_type: "all_members" } });
      reloadRules();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function removeRule(ruleId) {
    try { await api(`/notification-schemes/${selected}/rules/${ruleId}`, { method: "DELETE" }); reloadRules(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="stack gap-4">
      <div className="row" style={{ justifyContent: "space-between", alignItems: "center" }}>
        <div>
          <h3 style={{ margin: 0 }}>Notification Schemes</h3>
          <p className="text-sm secondary" style={{ margin: "4px 0 0" }}>Configure which events trigger notifications for project members.</p>
        </div>
        <Button variant="primary" icon="plus" onClick={createScheme}>New scheme</Button>
      </div>
      <div className="row gap-4" style={{ alignItems: "flex-start" }}>
        <div className="card" style={{ width: 240, overflow: "hidden", flexShrink: 0 }}>
          {loading ? <div className="text-sm muted" style={{ padding: 16 }}>Loading…</div> : (
            schemes.map((s) => (
              <button key={s.id} className="nav-item" aria-current={selected === s.id ? "page" : undefined}
                style={{ width: "100%", textAlign: "left", padding: "10px 14px", fontSize: 13 }}
                onClick={() => setSelected(s.id)}>
                {s.name}
              </button>
            ))
          )}
        </div>
        {selected ? (
          <div className="card card-pad" style={{ flex: 1 }}>
            <h4 style={{ margin: "0 0 12px" }}>Notification rules</h4>
            <div className="stack gap-2">
              {EVENT_TYPES.map((evt) => {
                const rule = rules.find((r) => r.event_type === evt);
                return (
                  <div key={evt} className="row gap-3" style={{ alignItems: "center", padding: "8px 0", borderBottom: "1px solid var(--border)" }}>
                    <label className="row gap-2" style={{ cursor: "pointer", flex: 1 }}>
                      <input type="checkbox" checked={!!rule}
                        onChange={() => rule ? removeRule(rule.id) : addRule(evt)}/>
                      <span style={{ fontSize: 13 }}>{evt.replace(/_/g, " ")}</span>
                    </label>
                    {rule && <span className="badge" style={{ fontSize: 11 }}>{rule.notify_type || "all_members"}</span>}
                  </div>
                );
              })}
            </div>
          </div>
        ) : (
          <div className="card card-pad" style={{ flex: 1, textAlign: "center", color: "var(--text-muted)", padding: 40 }}>
            Select a scheme to configure rules
          </div>
        )}
      </div>
    </div>
  );
}
