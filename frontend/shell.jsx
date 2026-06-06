// shell.jsx — sidebar, topbar, command palette, notifications dropdown, app frame

const { useState: useStateS, useEffect: useEffectS, useRef: useRefS, useMemo: useMemoS, useCallback: useCallbackS } = React;

// ─── Sidebar ────────────────────────────────────────────
function Sidebar({ route, nav, tg }) {
  const [expanded, setExpanded] = useStateS(true);
  const project = FORGE_DATA.PROJECTS.find((p) => p.id === FORGE_DATA.ACTIVE_PROJECT_ID);

  const isActive = (p) => route.path === p;

  return (
    <aside className="sidebar">
      <div className="sidebar-head">
        <div className="sidebar-logo" aria-hidden="true">F</div>
        <div className="stack" style={{ lineHeight: 1.15 }}>
          <span className="sidebar-brand">Forge</span>
          <span className="sidebar-brand-meta">forge.dev</span>
        </div>
      </div>

      <button className="proj-switcher" onClick={() => nav("projects")}>
        <span className="proj-icon" style={{ background: project.color }}>{project.key.slice(0,2)}</span>
        <span className="stack" style={{ lineHeight: 1.15, textAlign: "left" }}>
          <span>{project.name}</span>
          <span style={{ color: "var(--sidebar-text-muted)", fontSize: 11, fontWeight: 400 }}>{project.key} · 8 members</span>
        </span>
        <span style={{ marginLeft: "auto" }}><Icon name="chevronDown" size={14} color="#94A3B8"/></span>
      </button>

      <nav className="sidebar-nav">
        <div className="sidebar-section">Workspace</div>
        <button className="nav-item" aria-current={isActive("dashboard") ? "page" : undefined} onClick={() => nav("dashboard")}>
          <Icon name="home" size={16}/> Dashboard
        </button>
        <button className="nav-item" aria-current={isActive("my-issues") ? "page" : undefined} onClick={() => nav("my-issues")}>
          <Icon name="checkbox" size={16}/> My issues
          <span className="nav-count">7</span>
        </button>
        <button className="nav-item" aria-current={isActive("projects") ? "page" : undefined} onClick={() => nav("projects")}>
          <Icon name="briefcase" size={16}/> All projects
        </button>
        <button className="nav-item" aria-current={isActive("notifications-page") ? "page" : undefined} onClick={() => nav("notifications-page")}>
          <Icon name="bell" size={16}/> Notifications
          <span className="nav-badge">12</span>
        </button>

        <div className="sidebar-section row" style={{ justifyContent: "space-between", paddingRight: 8 }}>
          <span>Project · {project.key}</span>
          <button onClick={() => setExpanded((e) => !e)} style={{ border: 0, background: "transparent", color: "var(--sidebar-text-muted)" }} aria-label="Toggle project nav">
            <Icon name={expanded ? "chevronDown" : "chevronRight"} size={12}/>
          </button>
        </div>
        {expanded && (
          <div>
            <button className="nav-item" aria-current={isActive("board") ? "page" : undefined} onClick={() => nav("board")}>
              <Icon name="kanban" size={16}/> Board
            </button>
            <button className="nav-item" aria-current={isActive("backlog") ? "page" : undefined} onClick={() => nav("backlog")}>
              <Icon name="list" size={16}/> Backlog
            </button>
            <button className="nav-item" aria-current={isActive("sprint") ? "page" : undefined} onClick={() => nav("sprint")}>
              <Icon name="rocket" size={16}/> Sprints
            </button>
            <button className="nav-item" aria-current={isActive("wiki") ? "page" : undefined} onClick={() => nav("wiki")}>
              <Icon name="notes" size={16}/> Wiki
            </button>
            <button className="nav-item" aria-current={isActive("members") ? "page" : undefined} onClick={() => nav("members")}>
              <Icon name="users" size={16}/> Members
            </button>
            <button className="nav-item" aria-current={isActive("reports") ? "page" : undefined} onClick={() => nav("reports")}>
              <Icon name="chart" size={16}/> Reports
            </button>
            <button className="nav-item" aria-current={isActive("settings") ? "page" : undefined} onClick={() => nav("settings")}>
              <Icon name="settings" size={16}/> Settings
            </button>
          </div>
        )}

        <div className="sidebar-section">Integrations</div>
        <button className="nav-item" aria-current={isActive("integrations") ? "page" : undefined} onClick={() => nav("integrations")}>
          <Icon name="telegram" size={16} color="#2AABEE"/> Telegram bot
          <Badge tone={tg ? "tg" : "muted"} style={{ marginLeft: "auto", fontSize: 10 }}>{tg ? "On" : "Off"}</Badge>
        </button>
        <button className="nav-item" aria-current={isActive("admin") ? "page" : undefined} onClick={() => nav("admin")}>
          <Icon name="shield" size={16}/> Admin
        </button>
      </nav>

      <div className="sidebar-foot">
        <Avatar user={FORGE_DATA.ME} size="md"/>
        <div className="stack" style={{ flex: 1, minWidth: 0 }}>
          <span className="name" style={{ overflow:"hidden", textOverflow:"ellipsis", whiteSpace:"nowrap" }}>{FORGE_DATA.ME.name}</span>
          <span className="role">{FORGE_DATA.ME.role} · {FORGE_DATA.ME.tg || "no telegram"}</span>
        </div>
        <button className="icon-btn" style={{ color: "var(--sidebar-text-muted)" }} aria-label="Sign out" onClick={() => nav("login")}>
          <Icon name="exit" size={15}/>
        </button>
      </div>
    </aside>
  );
}

// ─── Notifications dropdown ─────────────────────────────
const NOTIFS = [
  { id: 1, who: "u3", text: "moved INFRA-220 to In Progress", when: "12m", unread: true, via: "tg" },
  { id: 2, who: "u5", text: "commented on INFRA-232", when: "28m", unread: true, via: "tg" },
  { id: 3, who: "u2", text: "assigned INFRA-242 to you", when: "1h",  unread: true, via: "tg" },
  { id: 4, who: "u1", text: "merged PR for INFRA-201", when: "3h",  unread: false },
  { id: 5, who: "u4", text: "mentioned you in INFRA-222", when: "yesterday", unread: false, via: "tg" },
];

function NotificationsPopover({ onClose, tg }) {
  const ref = useRefS(null);
  useEffectS(() => {
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) onClose(); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [onClose]);
  return (
    <div ref={ref} style={{
      position: "absolute", top: 44, right: 0,
      width: 380, background: "var(--bg)", border: "1px solid var(--border)",
      borderRadius: 12, boxShadow: "var(--shadow-lg)", zIndex: 60, overflow: "hidden"
    }}>
      <div style={{ padding: "12px 16px", borderBottom: "1px solid var(--border)", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <div style={{ fontWeight: 600 }}>Notifications</div>
        <button className="btn btn-ghost" data-size="sm">Mark all read</button>
      </div>
      <div style={{ maxHeight: 380, overflowY: "auto" }}>
        {NOTIFS.map((n) => {
          const u = FORGE_DATA.PEOPLE.find((p) => p.id === n.who);
          return (
            <div key={n.id} className={"notif-row" + (n.unread ? " unread" : "")}>
              <Avatar user={u} size="md"/>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div className="text-sm">
                  <span className="bold">{u.name}</span> <span className="secondary">{n.text}</span>
                </div>
                <div className="text-xs muted row gap-2" style={{ marginTop: 2 }}>
                  {n.via === "tg" && (
                    <span className="row gap-1" style={{ color: "var(--tg)" }}>
                      <Icon name="telegram" size={11}/> via Telegram
                    </span>
                  )}
                  <span>·</span>
                  <span>{n.when}</span>
                </div>
              </div>
              {n.unread && <span style={{ width: 6, height: 6, borderRadius: "50%", background: "var(--indigo-600)", marginTop: 6 }}/>}
            </div>
          );
        })}
      </div>
      <div style={{ padding: 8, borderTop: "1px solid var(--border)", display: "flex", gap: 8, alignItems: "center", justifyContent: "space-between" }}>
        <span className="text-xs muted row gap-1">
          <Icon name="telegram" size={12} color="#2AABEE"/>
          Bot {tg ? <span style={{ color: "var(--tg)" }}>connected</span> : <span>not connected</span>}
        </span>
        <button className="btn btn-ghost" data-size="sm">View all</button>
      </div>
    </div>
  );
}

// ─── Command palette ────────────────────────────────────
function CommandPalette({ open, onClose, nav }) {
  const [q, setQ] = useStateS("");
  const [active, setActive] = useStateS(0);
  useEffectS(() => { if (open) { setQ(""); setActive(0); } }, [open]);

  const items = useMemoS(() => {
    const all = [
      { g: "Go to", icon: "home",      label: "Dashboard",                shortcut: "G D", to: "dashboard" },
      { g: "Go to", icon: "kanban",    label: "Board",                    shortcut: "G B", to: "board" },
      { g: "Go to", icon: "list",      label: "Backlog",                  shortcut: "G L", to: "backlog" },
      { g: "Go to", icon: "rocket",    label: "Sprints",                  shortcut: "G S", to: "sprint" },
      { g: "Go to", icon: "notes",     label: "Wiki",                              to: "wiki" },
      { g: "Go to", icon: "users",     label: "Members",                  shortcut: "G M", to: "members" },
      { g: "Go to", icon: "chart",     label: "Reports",                            to: "reports" },
      { g: "Go to", icon: "telegram",  label: "Telegram integration",               to: "integrations" },
      { g: "Actions", icon: "plus",    label: "Create new issue",          shortcut: "C",   action: "createIssue" },
      { g: "Actions", icon: "user",    label: "Invite member",                       action: "inviteMember" },
      { g: "Actions", icon: "rocket",  label: "Start new sprint",                    action: "newSprint" },
      { g: "Actions", icon: "moon",    label: "Toggle dark mode",                    action: "toggleDark" },
      { g: "Recent", icon: "bug",     label: "INFRA-232 — Investigate Loki ingester OOM", to: "issue/INFRA-232" },
      { g: "Recent", icon: "story",   label: "INFRA-220 — Provision shared etcd cluster", to: "issue/INFRA-220" },
      { g: "Recent", icon: "task",    label: "INFRA-211 — Helm chart: replace deprecated PSP", to: "issue/INFRA-211" },
    ];
    if (!q) return all;
    return all.filter((i) => i.label.toLowerCase().includes(q.toLowerCase()));
  }, [q]);

  const groups = useMemoS(() => {
    const m = {};
    items.forEach((it) => { (m[it.g] ||= []).push(it); });
    return m;
  }, [items]);

  useEffectS(() => {
    if (!open) return;
    const h = (e) => {
      if (e.key === "Escape") { onClose(); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setActive((a) => Math.min(items.length - 1, a + 1)); }
      else if (e.key === "ArrowUp")   { e.preventDefault(); setActive((a) => Math.max(0, a - 1)); }
      else if (e.key === "Enter")     { e.preventDefault(); const it = items[active]; if (it?.to) { nav(it.to); onClose(); } else if (it?.action === "toggleDark") { document.documentElement.dataset.theme = document.documentElement.dataset.theme === "dark" ? "light" : "dark"; onClose(); } else if (it) { onClose(); } }
    };
    window.addEventListener("keydown", h);
    return () => window.removeEventListener("keydown", h);
  }, [open, items, active, onClose, nav]);

  if (!open) return null;
  let cursor = 0;
  return (
    <div className="cmdk" onClick={onClose}>
      <div className="cmdk-box" onClick={(e) => e.stopPropagation()}>
        <input className="cmdk-input" autoFocus placeholder="Type a command or search…" value={q} onChange={(e) => { setQ(e.target.value); setActive(0); }}/>
        <div className="cmdk-list">
          {Object.entries(groups).map(([g, list]) => (
            <div key={g}>
              <div className="cmdk-group-label">{g}</div>
              {list.map((it) => {
                const idx = cursor++;
                return (
                  <div key={idx} className={"cmdk-item" + (idx === active ? " active" : "")}
                    onMouseEnter={() => setActive(idx)}
                    onClick={() => { if (it.to) { nav(it.to); onClose(); } else onClose(); }}>
                    <Icon name={it.icon} size={15}/>
                    <span>{it.label}</span>
                    {it.shortcut && <span className="meta">{it.shortcut}</span>}
                  </div>
                );
              })}
            </div>
          ))}
          {items.length === 0 && <div style={{ padding: 24, textAlign: "center", color: "var(--text-muted)" }}>No results</div>}
        </div>
      </div>
    </div>
  );
}

// ─── Topbar ─────────────────────────────────────────────
function Topbar({ nav, route, onOpenCommand, tg }) {
  const [notifOpen, setNotifOpen] = useStateS(false);

  const crumbMap = {
    dashboard: ["Forge", "Dashboard"],
    board:     ["Core Infrastructure", "Board"],
    backlog:   ["Core Infrastructure", "Backlog"],
    sprint:    ["Core Infrastructure", "Sprints"],
    wiki:      ["Core Infrastructure", "Wiki"],
    members:   ["Core Infrastructure", "Members"],
    reports:   ["Core Infrastructure", "Reports"],
    settings:  ["Core Infrastructure", "Settings"],
    integrations: ["Core Infrastructure", "Settings", "Integrations"],
    "notifications-page": ["You", "Notification preferences"],
    "my-issues": ["You", "My issues"],
    projects:  ["Forge", "All projects"],
    issue:     ["Core Infrastructure", "Issues", route.rest[0] || ""],
  };
  const crumbs = crumbMap[route.path] || ["Forge"];

  return (
    <header className="topbar">
      <div className="row gap-2 text-sm" style={{ color: "var(--text-muted)" }}>
        {crumbs.map((c, i) => (
          <React.Fragment key={i}>
            {i > 0 && <Icon name="chevronRight" size={12}/>}
            <span style={{ color: i === crumbs.length - 1 ? "var(--text)" : undefined, fontWeight: i === crumbs.length - 1 ? 500 : 400 }}>{c}</span>
          </React.Fragment>
        ))}
      </div>
      <div className="spacer"/>
      <button className="search" onClick={onOpenCommand}>
        <Icon name="search" size={14}/>
        <span style={{ flex: 1, textAlign: "left" }}>Search or jump to…</span>
        <span className="kbd">⌘K</span>
      </button>
      <div className="topbar-actions" style={{ position: "relative" }}>
        <Button variant="primary" size="sm" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>
          Create
        </Button>
        <button className="icon-btn" aria-label="Activity"><Icon name="bolt" size={17}/></button>
        <button className="icon-btn" aria-label="Notifications" onClick={() => setNotifOpen((o) => !o)}>
          <Icon name="bell" size={17}/>
          <span className="dot"/>
        </button>
        {notifOpen && <NotificationsPopover onClose={() => setNotifOpen(false)} tg={tg}/>}
        <button className="icon-btn" style={{ marginLeft: 4 }} aria-label="Account">
          <Avatar user={FORGE_DATA.ME} size="sm"/>
        </button>
      </div>
    </header>
  );
}

Object.assign(window, { Sidebar, Topbar, CommandPalette });
