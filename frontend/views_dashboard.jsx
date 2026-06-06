// views_dashboard.jsx

// Feature 24: live activity feed (auto-refresh, skeleton, formatted strings)
function ActivityFeed() {
  const [items, setItems] = React.useState(null);
  const toast = useToast();

  React.useEffect(() => {
    let live = true;
    const load = () => api("/activity?limit=20").then((d) => { if (live) setItems(d); }).catch((e) => { if (live && items === null) toast(e.message, { icon: "x", color: "#F87171" }); });
    load();
    const t = setInterval(load, 60000); // auto-refresh every 60s
    return () => { live = false; clearInterval(t); };
  }, []);

  if (items === null) {
    return (
      <div style={{ padding: "8px 16px 16px" }}>
        {[0, 1, 2, 3, 4].map((i) => (
          <div key={i} className="row gap-3" style={{ padding: "8px 0" }}>
            <div className="skel" style={{ width: 22, height: 22, borderRadius: "50%" }}/>
            <div style={{ flex: 1 }}><div className="skel" style={{ height: 11, width: "80%", marginBottom: 6 }}/><div className="skel" style={{ height: 9, width: 60 }}/></div>
          </div>
        ))}
      </div>
    );
  }

  const isIssue = (id) => /^[A-Z]+-\d+$/.test(id || "");
  const phrase = (a) => {
    switch (a.action) {
      case "moved": return <>moved <Ent a={a}/> to <span className="bold">{(a.meta && a.meta.to) || "a column"}</span></>;
      case "assigned": return <>assigned <Ent a={a}/> to <span className="bold">{(a.meta && a.meta.to) || "someone"}</span></>;
      case "commented": return <>commented on <Ent a={a}/></>;
      case "logged time on": return <>logged time on <Ent a={a}/></>;
      case "created": return <>created {a.entity_type} <Ent a={a}/></>;
      case "closed": return <>closed <Ent a={a}/></>;
      case "resolved": return <>resolved <Ent a={a}/></>;
      case "started": return <>started <Ent a={a}/></>;
      case "released": return <>released version <Ent a={a}/></>;
      case "updated": return <>updated <Ent a={a}/></>;
      case "watched": return <>started watching <Ent a={a}/></>;
      default: return <>{a.action} <Ent a={a}/></>;
    }
  };

  function Ent({ a }) {
    if (isIssue(a.entity_id)) return <a href={"#/issue/" + a.entity_id} style={{ color: "var(--indigo-600)", textDecoration: "none", fontWeight: 500 }}>{a.entity_id}</a>;
    return <span className="medium" style={{ color: "var(--text)" }}>{a.entity_id}</span>;
  }

  return (
    <div style={{ padding: "8px 16px 16px" }}>
      {items.map((a, idx) => {
        const u = a.actor || FORGE_DATA.PEOPLE.find((p) => p.id === a.actor_id);
        return (
          <div key={a.id} className="row gap-3" style={{ padding: "8px 0", position: "relative" }}>
            <div style={{ position: "relative" }}>
              <Avatar user={u} size="sm"/>
              {idx < items.length - 1 && <div style={{ position: "absolute", left: 10, top: 26, bottom: -8, width: 1, background: "var(--border)" }}/>}
            </div>
            <div style={{ flex: 1, paddingTop: 1 }}>
              <div className="text-sm"><span className="bold">{u ? u.name : "Someone"}</span> <span className="secondary">{phrase(a)}</span></div>
              <div className="text-xs muted" style={{ marginTop: 2 }}>{fmtDate(a.created_at)}</div>
            </div>
          </div>
        );
      })}
      {items.length === 0 && <div className="text-sm muted" style={{ padding: 12, textAlign: "center" }}>No recent activity.</div>}
    </div>
  );
}

function DashboardView({ nav, tg }) {
  const issues = FORGE_DATA.ISSUES;
  const myIssues = issues.filter((i) => i.assignee === FORGE_DATA.ME.id);
  const openIssues = issues.filter((i) => i.status !== "Done");

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/">Forge</a> <Icon name="chevronRight" size={11}/> <span>Dashboard</span></div>
          <h1>Good afternoon, Maya</h1>
          <p>Here's what's happening across Core Infrastructure today.</p>
        </div>
        <div className="row gap-2">
          <Button icon="calendar">This week</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>Create</Button>
        </div>
      </div>

      <div className="dash">
        {/* Stats */}
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic"><Icon name="briefcase" size={15}/></div>
          <div className="label">Active projects</div>
          <div className="num">6</div>
          <div className="trend up"><Icon name="arrowUp" size={12}/> 1 this month</div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--warning-bg)", color: "var(--warning)" }}><Icon name="checkbox" size={15}/></div>
          <div className="label">Open issues</div>
          <div className="num">{openIssues.length + 145}</div>
          <div className="trend down"><Icon name="arrowDown" size={12}/> 8 vs last week</div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--info-bg)", color: "var(--info)" }}><Icon name="user" size={15}/></div>
          <div className="label">Assigned to me</div>
          <div className="num">{myIssues.length}</div>
          <div className="trend"><span className="muted">2 due this week</span></div>
        </div>
        <div className="stat" style={{ gridColumn: "span 3" }}>
          <div className="ic" style={{ background: "var(--success-bg)", color: "var(--success)" }}><Icon name="rocket" size={15}/></div>
          <div className="label">Sprint 24 progress</div>
          <div className="num">68<span style={{ fontSize: 16, fontWeight: 500, color: "var(--text-muted)" }}>%</span></div>
          <div style={{ height: 6, background: "var(--bg-muted)", borderRadius: 99, marginTop: 8, overflow: "hidden" }}>
            <div style={{ width: "68%", height: "100%", background: "linear-gradient(90deg, #6366F1, #818CF8)" }}/>
          </div>
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
                borderBottom: "1px solid var(--border)", cursor: "default"
              }} onMouseEnter={(e) => e.currentTarget.style.background = "var(--bg-subtle)"} onMouseLeave={(e) => e.currentTarget.style.background = ""}>
                <TypeIcon value={i.type}/>
                <span className="mono text-xs muted" style={{ width: 76 }}>{i.id}</span>
                <span className="text-sm" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{i.title}</span>
                <PriorityBadge value={i.pri}/>
                <StatusBadge value={i.status}/>
                <span className="text-xs muted" style={{ width: 60, textAlign: "right" }}>{i.due || "—"}</span>
              </div>
            ))}
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

        {/* Recent projects */}
        <div className="card" style={{ gridColumn: "span 8" }}>
          <div className="card-head">
            <h3>Your projects</h3>
            <button className="btn btn-ghost" data-size="sm" onClick={() => nav("projects")}>All projects</button>
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 0 }}>
            {FORGE_DATA.PROJECTS.slice(0, 6).map((p, idx) => (
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
                  <AvatarStack users={FORGE_DATA.PEOPLE.slice(0, Math.min(p.members, 5))} max={4}/>
                  <span className="text-xs muted">Updated {p.updated}</span>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Telegram card */}
        <div className="card tg-card" style={{ gridColumn: "span 4", padding: 0, overflow: "hidden" }}>
          <div style={{ padding: "14px 16px", borderBottom: "1px solid var(--border)", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
            <div className="row gap-2">
              <Icon name="telegram" size={18} color="#2AABEE"/>
              <span className="bold">Telegram bot</span>
            </div>
            <Badge tone={tg ? "tg" : "muted"} dot>{tg ? "Connected" : "Not connected"}</Badge>
          </div>
          <div style={{ padding: 16 }}>
            <div className="text-sm secondary" style={{ marginBottom: 12 }}>
              {tg
                ? "Bot is live as @forge_team_bot. 7 of 8 members are receiving notifications."
                : "Connect Telegram so your team gets notified instantly — assignments, comments, deploys."}
            </div>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8, marginBottom: 12 }}>
              <div style={{ padding: 10, background: "var(--bg-subtle)", borderRadius: 6 }}>
                <div className="text-xs muted">Members connected</div>
                <div style={{ fontSize: 22, fontWeight: 600 }}>{tg ? "7 / 8" : "0 / 8"}</div>
              </div>
              <div style={{ padding: 10, background: "var(--bg-subtle)", borderRadius: 6 }}>
                <div className="text-xs muted">Sent today</div>
                <div style={{ fontSize: 22, fontWeight: 600 }}>{tg ? "143" : "0"}</div>
              </div>
            </div>
            <Button variant={tg ? "secondary" : "telegram"} icon="telegram" style={{ width: "100%", justifyContent: "center" }} onClick={() => nav("integrations")}>
              {tg ? "Configure" : "Connect bot"}
            </Button>
          </div>
        </div>

        {/* Sprint */}
        <div className="card" style={{ gridColumn: "span 7" }}>
          <div className="card-head">
            <h3>Sprint 24 — Edge & Reliability</h3>
            <span className="text-xs muted">Dec 2 → Dec 15 · 5 days left</span>
          </div>
          <div style={{ padding: 16 }}>
            <div className="row gap-3" style={{ marginBottom: 12 }}>
              <div style={{ flex: 1 }}>
                <div className="text-xs muted">Committed</div>
                <div style={{ fontSize: 20, fontWeight: 600 }}>62 pts</div>
              </div>
              <div style={{ flex: 1 }}>
                <div className="text-xs muted">Completed</div>
                <div style={{ fontSize: 20, fontWeight: 600, color: "var(--success)" }}>42 pts</div>
              </div>
              <div style={{ flex: 1 }}>
                <div className="text-xs muted">In progress</div>
                <div style={{ fontSize: 20, fontWeight: 600, color: "var(--info)" }}>13 pts</div>
              </div>
              <div style={{ flex: 1 }}>
                <div className="text-xs muted">Remaining</div>
                <div style={{ fontSize: 20, fontWeight: 600, color: "var(--text-secondary)" }}>7 pts</div>
              </div>
            </div>
            {/* tiny burndown */}
            <BurndownMini/>
          </div>
        </div>

        {/* Commits */}
        <div className="card" style={{ gridColumn: "span 5" }}>
          <div className="card-head">
            <h3>Linked commits</h3>
            <button className="btn btn-ghost" data-size="sm"><Icon name="branch" size={13}/> main</button>
          </div>
          <div>
            {FORGE_DATA.COMMITS.map((c, idx) => {
              const u = FORGE_DATA.PEOPLE.find((p) => p.id === c.who);
              return (
                <div key={c.sha} className="row gap-3" style={{ padding: "10px 16px", borderBottom: idx < FORGE_DATA.COMMITS.length - 1 ? "1px solid var(--border)" : 0 }}>
                  <Avatar user={u} size="sm"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="text-sm" style={{ whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>{c.msg}</div>
                    <div className="text-xs muted row gap-2">
                      <span className="mono">{c.sha}</span> · <span>{u.name}</span> · <span>{c.when}</span>
                    </div>
                  </div>
                  <Icon name="branch" size={14} color="var(--text-muted)"/>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}

function BurndownMini() {
  // synthetic burndown
  const ideal = [62, 58, 53, 48, 44, 40, 35, 31, 27, 22, 18, 13, 9, 4, 0];
  const actual = [62, 60, 58, 51, 48, 44, 41, 38, 34, 28, 21, null, null, null, null];
  const max = 62;
  const w = 600, h = 100, pad = 20;
  const stepX = (w - pad * 2) / (ideal.length - 1);
  const y = (v) => h - pad - (v / max) * (h - pad * 2);
  const idealPath = ideal.map((v, i) => `${i ? "L" : "M"} ${pad + i * stepX} ${y(v)}`).join(" ");
  const actualPts = actual.filter((v) => v !== null);
  const actualPath = actualPts.map((v, i) => `${i ? "L" : "M"} ${pad + i * stepX} ${y(v)}`).join(" ");

  return (
    <svg viewBox={`0 0 ${w} ${h}`} style={{ width: "100%", height: 110 }}>
      <defs>
        <linearGradient id="bd" x1="0" x2="0" y1="0" y2="1">
          <stop offset="0%" stopColor="#6366F1" stopOpacity=".3"/>
          <stop offset="100%" stopColor="#6366F1" stopOpacity="0"/>
        </linearGradient>
      </defs>
      {/* axis */}
      <line x1={pad} y1={h-pad} x2={w-pad} y2={h-pad} stroke="var(--border)"/>
      <line x1={pad} y1={pad} x2={pad} y2={h-pad} stroke="var(--border)"/>
      {/* ideal dashed */}
      <path d={idealPath} fill="none" stroke="var(--text-muted)" strokeDasharray="4 4"/>
      {/* actual area */}
      <path d={actualPath + ` L ${pad + (actualPts.length - 1) * stepX} ${h - pad} L ${pad} ${h - pad} Z`} fill="url(#bd)"/>
      <path d={actualPath} fill="none" stroke="#6366F1" strokeWidth="2"/>
      {actualPts.map((v, i) => (
        <circle key={i} cx={pad + i * stepX} cy={y(v)} r="2.5" fill="#6366F1"/>
      ))}
      <text x={w - pad} y={pad} fill="var(--text-muted)" fontSize="10" textAnchor="end">ideal</text>
      <text x={w - pad} y={pad + 12} fill="#6366F1" fontSize="10" textAnchor="end">actual</text>
    </svg>
  );
}

window.DashboardView = DashboardView;
