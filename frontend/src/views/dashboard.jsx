// dashboard.jsx — Dashboard view
import { useState, useEffect } from 'react';
import { Icon } from '../components/icons';
import { Avatar, AvatarStack, Badge, Button, TypeIcon, PriorityBadge, StatusBadge, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api } from '../api/api';
import { adaptUser, fmtDate } from '../api/adapters';

function ActivityFeed() {
  const [items, setItems] = useState(null);
  const toast = useToast();

  useEffect(() => {
    let live = true;
    const load = () => api("/activity?limit=20")
      .then((d) => { if (live) setItems(d); })
      .catch((e) => { if (live && items === null) toast && toast(e.message, { icon: "x", color: "#F87171" }); });
    load();
    const t = setInterval(load, 60000);
    return () => { live = false; clearInterval(t); };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const isIssue = (id) => /^[A-Z]+-\d+$/.test(id || "");
  const phrase = (a) => {
    function Ent() {
      if (isIssue(a.entity_id)) return <a href={"#/issue/" + a.entity_id} style={{ color: "var(--indigo-600)", textDecoration: "none", fontWeight: 500 }}>{a.entity_id}</a>;
      return <span className="medium" style={{ color: "var(--text)" }}>{a.entity_id}</span>;
    }
    switch (a.action) {
      case "moved":        return <><span className="secondary">moved </span><Ent/><span className="secondary"> to <span className="bold">{(a.meta && a.meta.to) || "a column"}</span></span></>;
      case "assigned":     return <><span className="secondary">assigned </span><Ent/><span className="secondary"> to <span className="bold">{(a.meta && a.meta.to) || "someone"}</span></span></>;
      case "commented":    return <><span className="secondary">commented on </span><Ent/></>;
      case "created":      return <><span className="secondary">created {a.entity_type} </span><Ent/></>;
      case "closed":       return <><span className="secondary">closed </span><Ent/></>;
      case "resolved":     return <><span className="secondary">resolved </span><Ent/></>;
      case "started":      return <><span className="secondary">started </span><Ent/></>;
      case "updated":      return <><span className="secondary">updated </span><Ent/></>;
      default:             return <><span className="secondary">{a.action} </span><Ent/></>;
    }
  };

  if (items === null) {
    return (
      <div style={{ padding: "8px 16px 16px" }}>
        {[0, 1, 2, 3, 4].map((i) => (
          <div key={i} className="row gap-3" style={{ padding: "8px 0" }}>
            <div className="skel" style={{ width: 22, height: 22, borderRadius: "50%" }}/>
            <div style={{ flex: 1 }}>
              <div className="skel" style={{ height: 11, width: "80%", marginBottom: 6 }}/>
              <div className="skel" style={{ height: 9, width: 60 }}/>
            </div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div style={{ padding: "8px 16px 16px" }}>
      {items.map((a, idx) => {
        const u = a.actor ? adaptUser(a.actor) : { name: "Someone", initials: "?", color: "#94A3B8" };
        return (
          <div key={a.id} className="row gap-3" style={{ padding: "8px 0", position: "relative" }}>
            <div style={{ position: "relative" }}>
              <Avatar user={u} size="sm"/>
              {idx < items.length - 1 && <div style={{ position: "absolute", left: 10, top: 26, bottom: -8, width: 1, background: "var(--border)" }}/>}
            </div>
            <div style={{ flex: 1, paddingTop: 1 }}>
              <div className="text-sm">
                <span className="bold">{u.name}</span> {phrase(a)}
              </div>
              <div className="text-xs muted" style={{ marginTop: 2 }}>{fmtDate(a.created_at)}</div>
            </div>
          </div>
        );
      })}
      {items.length === 0 && <div className="text-sm muted" style={{ padding: 12, textAlign: "center" }}>No recent activity.</div>}
    </div>
  );
}

function BurndownMini({ sprintId, totalPoints }) {
  const [data, setData] = useState(null);

  useEffect(() => {
    if (!sprintId) return;
    api("/sprints/" + sprintId + "/burndown")
      .then((d) => setData(d))
      .catch(() => {});
  }, [sprintId]);

  const ideal  = data?.ideal  || [];
  const actual = data?.actual || [];
  const max    = data?.total  || totalPoints || 1;

  if (!sprintId || (!ideal.length && !actual.length)) {
    return (
      <div style={{ height: 110, display: "grid", placeItems: "center", color: "var(--text-muted)", fontSize: 13 }}>
        No active sprint
      </div>
    );
  }

  const pts = Math.max(ideal.length, actual.length, 2);
  const w = 600, h = 100, pad = 20;
  const stepX = (w - pad * 2) / (pts - 1);
  const y = (v) => h - pad - (Math.max(0, v) / max) * (h - pad * 2);
  const idealPath = ideal.map((v, i) => `${i ? "L" : "M"} ${pad + i * stepX} ${y(v)}`).join(" ");
  const actualPts = actual.filter((v) => v !== null && v !== undefined);
  const actualPath = actualPts.map((v, i) => `${i ? "L" : "M"} ${pad + i * stepX} ${y(v)}`).join(" ");

  return (
    <svg viewBox={`0 0 ${w} ${h}`} style={{ width: "100%", height: 110 }}>
      <defs>
        <linearGradient id="bd" x1="0" x2="0" y1="0" y2="1">
          <stop offset="0%" stopColor="#6366F1" stopOpacity=".3"/>
          <stop offset="100%" stopColor="#6366F1" stopOpacity="0"/>
        </linearGradient>
      </defs>
      <line x1={pad} y1={h - pad} x2={w - pad} y2={h - pad} stroke="var(--border)"/>
      <line x1={pad} y1={pad} x2={pad} y2={h - pad} stroke="var(--border)"/>
      {idealPath && <path d={idealPath} fill="none" stroke="var(--text-muted)" strokeDasharray="4 4"/>}
      {actualPath && actualPts.length > 1 && (
        <>
          <path d={actualPath + ` L ${pad + (actualPts.length - 1) * stepX} ${h - pad} L ${pad} ${h - pad} Z`} fill="url(#bd)"/>
          <path d={actualPath} fill="none" stroke="#6366F1" strokeWidth="2"/>
          {actualPts.map((v, i) => <circle key={i} cx={pad + i * stepX} cy={y(v)} r="2.5" fill="#6366F1"/>)}
        </>
      )}
      <text x={w - pad} y={pad} fill="var(--text-muted)" fontSize="10" textAnchor="end">ideal</text>
      <text x={w - pad} y={pad + 12} fill="#6366F1" fontSize="10" textAnchor="end">actual</text>
    </svg>
  );
}

export function DashboardView({ nav }) {
  const { me, projects, people, issues, activeProjectId } = useApp();
  const myIssues  = issues.filter((i) => me && i.assignee === me.id);
  const openCount = issues.filter((i) => i.status !== "Done").length;
  const [activeSprint, setActiveSprint] = useState(null);

  useEffect(() => {
    if (!activeProjectId) return;
    api("/projects/" + activeProjectId + "/sprints")
      .then((d) => {
        const list = d?.items || d || [];
        const running = list.find((s) => s.status === "active");
        setActiveSprint(running || null);
      })
      .catch(() => {});
  }, [activeProjectId]);

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/">Forge</a> <Icon name="chevronRight" size={11}/> <span>Dashboard</span></div>
          <h1>Hello, {me ? me.name.split(" ")[0] : "there"}</h1>
          <p>Here's what's happening across your projects today.</p>
        </div>
        <div className="row gap-2">
          <Button icon="calendar">This week</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>Create</Button>
        </div>
      </div>

      <div className="dash">
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic"><Icon name="briefcase" size={15}/></div>
          <div className="label">Active projects</div>
          <div className="num">{projects.length}</div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--warning-bg)", color: "var(--warning)" }}><Icon name="checkbox" size={15}/></div>
          <div className="label">Open issues</div>
          <div className="num">{openCount}</div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--info-bg)", color: "var(--info)" }}><Icon name="user" size={15}/></div>
          <div className="label">Assigned to me</div>
          <div className="num">{myIssues.length}</div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--success-bg)", color: "var(--success)" }}><Icon name="rocket" size={15}/></div>
          <div className="label">People</div>
          <div className="num">{people.length}</div>
        </div>

        {/* My issues */}
        <div className="card" style={{ gridColumn: "span 8" }}>
          <div className="card-head">
            <h3>My open issues</h3>
            <button className="btn btn-ghost" data-size="sm" onClick={() => nav("my-issues")}>View all <Icon name="arrowRight" size={13}/></button>
          </div>
          <div>
            {myIssues.slice(0, 5).map((i) => (
              <div key={i.id} onClick={() => nav("issue/" + i.id)} style={{
                display: "grid", gridTemplateColumns: "auto auto 1fr auto auto auto", gap: 12,
                alignItems: "center", padding: "12px 20px",
                borderBottom: "1px solid var(--border)", cursor: "default",
              }} onMouseEnter={(e) => e.currentTarget.style.background = "var(--bg-subtle)"} onMouseLeave={(e) => e.currentTarget.style.background = ""}>
                <TypeIcon value={i.type}/>
                <span className="mono text-xs muted" style={{ width: 76 }}>{i.id}</span>
                <span className="text-sm" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{i.title}</span>
                <PriorityBadge value={i.pri}/>
                <StatusBadge value={i.status}/>
                <span className="text-xs muted" style={{ width: 60, textAlign: "right" }}>{i.due || "—"}</span>
              </div>
            ))}
            {myIssues.length === 0 && <div className="text-sm muted" style={{ padding: 16, textAlign: "center" }}>No open issues assigned to you.</div>}
          </div>
        </div>

        {/* Activity */}
        <div className="card" style={{ gridColumn: "span 4" }}>
          <div className="card-head">
            <h3>Recent activity</h3>
            <span className="text-xs muted row gap-1"><Icon name="refresh" size={12}/> live</span>
          </div>
          <ActivityFeed/>
        </div>

        {/* Projects */}
        <div className="card" style={{ gridColumn: "span 8" }}>
          <div className="card-head">
            <h3>Your projects</h3>
            <button className="btn btn-ghost" data-size="sm" onClick={() => nav("projects")}>All projects</button>
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 0 }}>
            {projects.slice(0, 6).map((p, idx) => (
              <div key={p.id} onClick={() => nav("board")} style={{
                padding: 16,
                borderRight: idx % 3 !== 2 ? "1px solid var(--border)" : 0,
                borderTop: idx >= 3 ? "1px solid var(--border)" : 0,
                cursor: "default",
              }} onMouseEnter={(e) => e.currentTarget.style.background = "var(--bg-subtle)"} onMouseLeave={(e) => e.currentTarget.style.background = ""}>
                <div className="row gap-3" style={{ marginBottom: 8 }}>
                  <div style={{ width: 32, height: 32, borderRadius: 8, background: p.color, color: "#fff", display: "grid", placeItems: "center" }}>
                    <Icon name={p.icon} size={16}/>
                  </div>
                  <div className="stack" style={{ lineHeight: 1.2 }}>
                    <span className="bold text-sm">{p.name}</span>
                    <span className="text-xs muted">{p.key} · {p.openIssues} open</span>
                  </div>
                </div>
                <div className="text-xs secondary" style={{ minHeight: 32, marginBottom: 8 }}>{p.desc}</div>
                <div className="row" style={{ justifyContent: "space-between" }}>
                  <AvatarStack users={people.slice(0, Math.min(p.members || 4, 5))} max={4}/>
                  <span className="text-xs muted">Updated {p.updated}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="card" style={{ gridColumn: "span 7" }}>
          <div className="card-head">
            <h3>{activeSprint ? activeSprint.name || "Current sprint" : "Sprint burndown"}</h3>
            {activeSprint && <span className="badge" data-tone="info">Active</span>}
          </div>
          <div style={{ padding: 16 }}>
            <BurndownMini sprintId={activeSprint?.id} totalPoints={activeSprint?.total_points}/>
          </div>
        </div>

        {/* Team members */}
        <div className="card" style={{ gridColumn: "span 5" }}>
          <div className="card-head">
            <h3>Team members</h3>
            <button className="btn btn-ghost" data-size="sm" onClick={() => nav("members")}>View all</button>
          </div>
          <div>
            {people.slice(0, 6).map((u, idx) => (
              <div key={u.id} className="row gap-3" style={{ padding: "10px 16px", borderBottom: idx < Math.min(people.length, 6) - 1 ? "1px solid var(--border)" : 0 }}>
                <Avatar user={u} size="sm"/>
                <div style={{ flex: 1, minWidth: 0 }}>
                  <div className="text-sm bold">{u.name}</div>
                  <div className="text-xs muted">{u.role ? u.role.charAt(0).toUpperCase() + u.role.slice(1) : "Member"} · {u.email}</div>
                </div>
                <Badge tone={u.status === "Active" ? "success" : "muted"} dot>{u.status}</Badge>
              </div>
            ))}
            {people.length === 0 && <div className="text-sm muted" style={{ padding: 16, textAlign: "center" }}>No team members loaded yet.</div>}
          </div>
        </div>
      </div>
    </div>
  );
}
