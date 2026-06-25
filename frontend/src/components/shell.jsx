// shell.jsx — Sidebar, Topbar, CommandPalette, NotificationsPopover
import { useState, useEffect, useRef, useMemo, useCallback } from 'react';
import { Icon } from './icons';
import { Avatar, Badge, Button, Modal, useToast } from './components';
import { useApp } from '../store/AppContext';
import { api } from '../api/api';
import { adaptUser } from '../api/adapters';

// ─── Keyboard shortcuts registry ────────────────────────
const SHORTCUTS = [
  { key: "C",          desc: "Create new issue",     group: "Global" },
  { key: "⌘K",        desc: "Open search / command palette", group: "Global" },
  { key: "?",          desc: "Show keyboard shortcuts",       group: "Global" },
  { key: "G D",        desc: "Go to Dashboard",      group: "Navigation" },
  { key: "G B",        desc: "Go to Board",          group: "Navigation" },
  { key: "G L",        desc: "Go to Backlog",        group: "Navigation" },
  { key: "G S",        desc: "Go to Sprints",        group: "Navigation" },
  { key: "G W",        desc: "Go to Wiki",           group: "Navigation" },
  { key: "G M",        desc: "Go to Members",        group: "Navigation" },
  { key: "G R",        desc: "Go to Reports",        group: "Navigation" },
  { key: "Escape",     desc: "Close modal / cancel", group: "Global" },
  { key: "J",          desc: "Next issue (in list)", group: "Issues" },
  { key: "K",          desc: "Previous issue",       group: "Issues" },
  { key: "E",          desc: "Edit issue title",     group: "Issues" },
];

export function KeyboardShortcutsModal({ open, onClose }) {
  const groups = {};
  SHORTCUTS.forEach((s) => { (groups[s.group] ||= []).push(s); });
  return (
    <Modal open={open} onClose={onClose} title="Keyboard shortcuts" footer={<Button onClick={onClose}>Close</Button>}>
      <div className="stack gap-4">
        {Object.entries(groups).map(([g, list]) => (
          <div key={g}>
            <div className="text-xs bold" style={{ color: "var(--text-muted)", textTransform: "uppercase", letterSpacing: ".06em", marginBottom: 8 }}>{g}</div>
            <div className="stack gap-1">
              {list.map((s) => (
                <div key={s.key} className="row gap-3" style={{ padding: "5px 0" }}>
                  <span style={{ flex: 1, fontSize: 13 }}>{s.desc}</span>
                  <div className="row gap-1">
                    {s.key.split(" ").map((k, i) => (
                      <span key={i} className="kbd" style={{ fontSize: 11, padding: "1px 5px", borderRadius: 4, border: "1px solid var(--border)", background: "var(--bg-subtle)", fontFamily: "ui-monospace, monospace" }}>{k}</span>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </Modal>
  );
}

// ─── Global keyboard shortcut handler ───────────────────
export function useKeyboardShortcuts(nav, onOpenCommand) {
  const [shortcutsOpen, setShortcutsOpen] = useState(false);
  const seq = useRef("");
  const seqTimer = useRef(null);

  useEffect(() => {
    function handler(e) {
      const tag = (e.target.tagName || "").toLowerCase();
      const inInput = tag === "input" || tag === "textarea" || tag === "select" || e.target.isContentEditable;

      // ⌘K always works
      if ((e.metaKey || e.ctrlKey) && e.key === "k") { e.preventDefault(); onOpenCommand(); return; }

      if (inInput) return;

      // single keys
      if (!e.metaKey && !e.ctrlKey && !e.altKey) {
        const k = e.key;
        if (k === "?") { setShortcutsOpen(true); return; }
        if (k === "c" || k === "C") { window.dispatchEvent(new CustomEvent("forge:create")); return; }
        if (k === "Escape") return;

        // Two-key sequences: G + X
        clearTimeout(seqTimer.current);
        seq.current += k.toUpperCase();
        if (seq.current.length === 1 && k.toUpperCase() === "G") {
          seqTimer.current = setTimeout(() => { seq.current = ""; }, 1000);
          return;
        }
        if (seq.current === "GD") { nav("dashboard"); seq.current = ""; return; }
        if (seq.current === "GB") { nav("board"); seq.current = ""; return; }
        if (seq.current === "GL") { nav("backlog"); seq.current = ""; return; }
        if (seq.current === "GS") { nav("sprint"); seq.current = ""; return; }
        if (seq.current === "GW") { nav("wiki"); seq.current = ""; return; }
        if (seq.current === "GM") { nav("members"); seq.current = ""; return; }
        if (seq.current === "GR") { nav("reports"); seq.current = ""; return; }
        seq.current = "";
      }
    }
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [nav, onOpenCommand]);

  return { shortcutsOpen, setShortcutsOpen };
}

// ─── Sidebar ────────────────────────────────────────────
export function Sidebar({ route, nav, mobileOpen, onClose }) {
  const { me, projects, activeProjectId, issues } = useApp();
  const [expanded, setExpanded] = useState(true);
  const project = (projects || []).find((p) => p.id === activeProjectId) || (projects && projects[0]) || { name: "…", key: "…", color: "#6366F1" };
  const myIssueCount = (issues || []).filter((i) => me && i.assignee === me.id).length;

  const isActive = (p) => route.path === p;

  function navAndClose(path) { nav(path); if (onClose) onClose(); }

  return (
    <aside className={`sidebar${mobileOpen ? " mobile-open" : ""}`}>
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
          <span style={{ color: "var(--sidebar-text-muted)", fontSize: 11, fontWeight: 400 }}>{project.key}</span>
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
          {myIssueCount > 0 && <span className="nav-count">{myIssueCount}</span>}
        </button>
        <button className="nav-item" aria-current={isActive("projects") ? "page" : undefined} onClick={() => nav("projects")}>
          <Icon name="briefcase" size={16}/> All projects
        </button>
        <button className="nav-item" aria-current={isActive("notifications-page") ? "page" : undefined} onClick={() => nav("notifications-page")}>
          <Icon name="bell" size={16}/> Notifications
        </button>

        <div className="sidebar-section row" style={{ justifyContent: "space-between", paddingRight: 8 }}>
          <span>Project · {project.key}</span>
          <button onClick={() => setExpanded((e) => !e)} style={{ border: 0, background: "transparent", color: "var(--sidebar-text-muted)" }} aria-label="Toggle project nav">
            <Icon name={expanded ? "chevronDown" : "chevronRight"} size={12}/>
          </button>
        </div>
        {expanded && (
          <div>
            {[
              ["board",    "kanban",  "Board"],
              ["backlog",  "list",    "Backlog"],
              ["sprint",   "rocket",  "Sprints"],
              ["roadmap",  "calendar","Roadmap"],
              ["wiki",     "notes",   "Wiki"],
              ["members",  "users",   "Members"],
              ["reports",  "chart",   "Reports"],
              ["settings", "settings","Settings"],
            ].map(([path, icon, label]) => (
              <button key={path} className="nav-item" aria-current={isActive(path) ? "page" : undefined} onClick={() => nav(path)}>
                <Icon name={icon} size={16}/> {label}
              </button>
            ))}
          </div>
        )}

        <div className="sidebar-section">Integrations</div>
        <button className="nav-item" aria-current={isActive("integrations") ? "page" : undefined} onClick={() => nav("integrations")}>
          <Icon name="telegram" size={16} color="#2AABEE"/> Telegram bot
        </button>
        <button className="nav-item" aria-current={isActive("admin") ? "page" : undefined} onClick={() => nav("admin")}>
          <Icon name="shield" size={16}/> Admin
        </button>
      </nav>

      <SidebarFoot nav={nav}/>
    </aside>
  );
}

function SidebarFoot({ nav }) {
  const { me, clearSession } = useApp();
  return (
    <div className="sidebar-foot">
      <Avatar user={me} size="md"/>
      <div className="stack" style={{ flex: 1, minWidth: 0 }}>
        <span className="name" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{me && me.name}</span>
        <span className="role">{me && me.role}</span>
      </div>
      <button className="icon-btn" style={{ color: "var(--sidebar-text-muted)" }} aria-label="Sign out" onClick={clearSession}>
        <Icon name="exit" size={15}/>
      </button>
    </div>
  );
}

// ─── Notifications dropdown ─────────────────────────────
function NotificationsPopover({ onClose }) {
  const ref = useRef(null);
  const [notifs, setNotifs] = useState([]);
  const { me } = useApp();

  useEffect(() => {
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) onClose(); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [onClose]);

  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api("/notifications?limit=20")
      .then((d) => setNotifs(Array.isArray(d) ? d : []))
      .catch(() => setNotifs([]))
      .finally(() => setLoading(false));
  }, []);

  function markAllRead() {
    api("/notifications/mark-all-read", { method: "POST" })
      .then(() => setNotifs((prev) => prev.map((n) => ({ ...n, is_read: true }))))
      .catch(() => {});
  }

  function fmtTime(ts) {
    if (!ts) return "";
    const d = new Date(ts);
    const diff = (Date.now() - d) / 1000;
    if (diff < 60)  return "just now";
    if (diff < 3600) return Math.floor(diff / 60) + "m ago";
    if (diff < 86400) return Math.floor(diff / 3600) + "h ago";
    return Math.floor(diff / 86400) + "d ago";
  }

  const unreadCount = notifs.filter((n) => !n.is_read).length;

  return (
    <div ref={ref} style={{
      position: "absolute", top: 44, right: 0,
      width: 380, background: "var(--bg)", border: "1px solid var(--border)",
      borderRadius: 12, boxShadow: "var(--shadow-lg)", zIndex: 60, overflow: "hidden",
    }}>
      <div style={{ padding: "12px 16px", borderBottom: "1px solid var(--border)", display: "flex", alignItems: "center", justifyContent: "space-between" }}>
        <div style={{ fontWeight: 600, display: "flex", alignItems: "center", gap: 8 }}>
          Notifications
          {unreadCount > 0 && (
            <span style={{ background: "var(--indigo-600)", color: "#fff", borderRadius: 99, fontSize: 11, fontWeight: 700, padding: "1px 7px" }}>{unreadCount}</span>
          )}
        </div>
        {unreadCount > 0 && (
          <button className="btn btn-ghost" data-size="sm" onClick={markAllRead}>Mark all read</button>
        )}
      </div>
      <div style={{ maxHeight: 380, overflowY: "auto" }}>
        {loading && (
          <div style={{ padding: 24, textAlign: "center", color: "var(--text-muted)", fontSize: 13 }}>Loading…</div>
        )}
        {!loading && notifs.length === 0 && (
          <div style={{ padding: 24, textAlign: "center", color: "var(--text-muted)", fontSize: 13 }}>No notifications yet.</div>
        )}
        {notifs.map((n) => {
          const u = n.actor ? adaptUser(n.actor) : { name: "System", initials: "SY", color: "#6366F1" };
          return (
            <div key={n.id} className={"notif-row" + (!n.is_read ? " unread" : "")}
              onClick={() => { if (!n.is_read) api("/notifications/" + n.id + "/read", { method: "POST" }).catch(() => {}); setNotifs((prev) => prev.map((x) => x.id === n.id ? { ...x, is_read: true } : x)); }}>
              <Avatar user={u} size="md"/>
              <div style={{ flex: 1, minWidth: 0 }}>
                <div className="text-sm" style={{ lineHeight: 1.4 }}>
                  <span className="bold">{u.name}</span>{" "}
                  <span className="secondary">{n.message}</span>
                </div>
                <div className="text-xs muted" style={{ marginTop: 2 }}>{fmtTime(n.created_at)}</div>
              </div>
              {!n.is_read && <span style={{ width: 7, height: 7, borderRadius: "50%", background: "var(--indigo-600)", flexShrink: 0, marginTop: 6 }}/>}
            </div>
          );
        })}
      </div>
    </div>
  );
}

// ─── Command palette ────────────────────────────────────
export function CommandPalette({ open, onClose, nav, onWikiPage }) {
  const [q, setQ] = useState("");
  const [active, setActive] = useState(0);
  const [searchResults, setSearchResults] = useState([]);
  const { projects, activeProjectId } = useApp();

  useEffect(() => { if (open) { setQ(""); setActive(0); setSearchResults([]); } }, [open]);

  // Debounced real-time search
  useEffect(() => {
    if (!q.trim() || q.length < 2) { setSearchResults([]); return; }
    const t = setTimeout(async () => {
      try {
        const res = await api("/search?q=" + encodeURIComponent(q) + "&limit=8");
        setSearchResults(res?.items || res?.data || []);
      } catch { setSearchResults([]); }
    }, 200);
    return () => clearTimeout(t);
  }, [q]);

  const navItems = useMemo(() => {
    const all = [
      { g: "Go to", icon: "home",     label: "Dashboard",              shortcut: "G D", to: "dashboard" },
      { g: "Go to", icon: "kanban",   label: "Board",                  shortcut: "G B", to: "board" },
      { g: "Go to", icon: "list",     label: "Backlog",                shortcut: "G L", to: "backlog" },
      { g: "Go to", icon: "rocket",   label: "Sprints",                shortcut: "G S", to: "sprint" },
      { g: "Go to", icon: "notes",    label: "Wiki",                                    to: "wiki" },
      { g: "Go to", icon: "users",    label: "Members",                shortcut: "G M", to: "members" },
      { g: "Go to", icon: "chart",    label: "Reports",                                 to: "reports" },
      { g: "Go to", icon: "telegram", label: "Telegram integration",                    to: "integrations" },
      { g: "Actions", icon: "plus",   label: "Create new issue",       shortcut: "C",   action: "createIssue" },
      { g: "Actions", icon: "moon",   label: "Toggle dark mode",                        action: "toggleDark" },
    ];
    if (!q) return all;
    return all.filter((i) => i.label.toLowerCase().includes(q.toLowerCase()));
  }, [q]);

  const TYPE_ICON = { page: "notes", issue: "checkbox", project: "briefcase", space: "users" };

  function executeItem(it) {
    if (it.to) { nav(it.to); onClose(); }
    else if (it.action === "createIssue") { window.dispatchEvent(new CustomEvent("forge:create")); onClose(); }
    else if (it.action === "toggleDark") { document.documentElement.dataset.theme = document.documentElement.dataset.theme === "dark" ? "light" : "dark"; onClose(); }
    else if (it.type === "page") {
      const spaceId = it.meta?.space_id;
      if (spaceId && onWikiPage) { onWikiPage(spaceId, it.id); onClose(); }
      else { nav("wiki"); onClose(); }
    }
    else if (it.type === "issue") { nav("issue", it.id); onClose(); }
    else if (it.type === "project") { nav("projects"); onClose(); }
    else if (it.type === "space") {
      if (onWikiPage) { onWikiPage(it.id, null); onClose(); }
      else { nav("wiki"); onClose(); }
    }
    else { onClose(); }
  }

  const allItems = [
    ...navItems,
    ...searchResults.map((r) => ({ ...r, g: r.type === "page" ? "Pages" : r.type === "issue" ? "Issues" : r.type === "space" ? "Spaces" : "Results" })),
  ];

  useEffect(() => {
    if (!open) return;
    const h = (e) => {
      if (e.key === "Escape") { onClose(); }
      else if (e.key === "ArrowDown") { e.preventDefault(); setActive((a) => Math.min(allItems.length - 1, a + 1)); }
      else if (e.key === "ArrowUp")   { e.preventDefault(); setActive((a) => Math.max(0, a - 1)); }
      else if (e.key === "Enter") { e.preventDefault(); const it = allItems[active]; if (it) executeItem(it); }
    };
    window.addEventListener("keydown", h);
    return () => window.removeEventListener("keydown", h);
  }, [open, allItems, active, onClose, nav, onWikiPage]);

  if (!open) return null;

  const groups = {};
  allItems.forEach((it) => { (groups[it.g] ||= []).push(it); });

  let cursor = 0;
  return (
    <div className="cmdk" onClick={onClose}>
      <div className="cmdk-box" onClick={(e) => e.stopPropagation()}>
        <input className="cmdk-input" autoFocus placeholder="Search pages, issues, spaces… or type a command"
          value={q} onChange={(e) => { setQ(e.target.value); setActive(0); }}/>
        <div className="cmdk-list">
          {Object.entries(groups).map(([g, list]) => (
            <div key={g}>
              <div className="cmdk-group-label">{g}</div>
              {list.map((it) => {
                const idx = cursor++;
                const iconName = it.type ? (TYPE_ICON[it.type] || "notes") : it.icon;
                return (
                  <div key={it.id || idx} className={"cmdk-item" + (idx === active ? " active" : "")}
                    onMouseEnter={() => setActive(idx)}
                    onClick={() => executeItem(it)}>
                    <Icon name={iconName} size={15}/>
                    <span style={{ flex: 1 }}>{it.label || it.title}</span>
                    {it.shortcut && <span className="meta">{it.shortcut}</span>}
                    {it.excerpt && <span className="meta" style={{ maxWidth: 200, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{it.excerpt}</span>}
                  </div>
                );
              })}
            </div>
          ))}
          {allItems.length === 0 && q.length >= 2 && (
            <div style={{ padding: 24, textAlign: "center", color: "var(--text-muted)" }}>No results for "{q}"</div>
          )}
          {allItems.length === 0 && q.length < 2 && (
            <div style={{ padding: 24, textAlign: "center", color: "var(--text-muted)" }}>Type to search…</div>
          )}
        </div>
      </div>
    </div>
  );
}

// ─── Topbar ─────────────────────────────────────────────
export function Topbar({ nav, route, onOpenCommand }) {
  const [notifOpen, setNotifOpen] = useState(false);
  const [unreadCount, setUnreadCount] = useState(0);
  const { me, projects, activeProjectId } = useApp();

  useEffect(() => {
    api("/notifications?limit=20&is_read=false")
      .then((d) => setUnreadCount(Array.isArray(d) ? d.filter((n) => !n.is_read).length : 0))
      .catch(() => {});
  }, [notifOpen]);
  const proj = (projects || []).find((p) => p.id === activeProjectId);
  const projName = proj ? proj.name : "…";

  const crumbMap = {
    dashboard:           ["Forge", "Dashboard"],
    board:               [projName, "Board"],
    backlog:             [projName, "Backlog"],
    sprint:              [projName, "Sprints"],
    wiki:                [projName, "Wiki"],
    members:             [projName, "Members"],
    reports:             [projName, "Reports"],
    settings:            [projName, "Settings"],
    integrations:        [projName, "Settings", "Integrations"],
    "notifications-page":[projName, "Notification preferences"],
    "my-issues":         ["You", "My issues"],
    projects:            ["Forge", "All projects"],
    issue:               [projName, "Issues", route.rest[0] || ""],
    admin:               ["Forge", "Admin"],
    profile:             ["You", "My profile"],
  };
  const crumbs = crumbMap[route.path] || ["Forge"];

  return (
    <header className="topbar">
      <button className="icon-btn topbar-menu-btn" aria-label="Menu" onClick={() => window.dispatchEvent(new CustomEvent("forge:mobile-menu"))}>
        <Icon name="menu" size={18}/>
      </button>
      <div className="row gap-2 text-sm" style={{ color: "var(--text-muted)" }}>
        {crumbs.map((c, i) => (
          <span key={i} style={{ display: "inline-flex", alignItems: "center", gap: 6 }}>
            {i > 0 && <Icon name="chevronRight" size={12}/>}
            <span style={{ color: i === crumbs.length - 1 ? "var(--text)" : undefined, fontWeight: i === crumbs.length - 1 ? 500 : 400 }}>{c}</span>
          </span>
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
        <button className="icon-btn" aria-label="Keyboard shortcuts (press ?)" title="Keyboard shortcuts" onClick={() => window.dispatchEvent(new CustomEvent("forge:shortcuts"))}><Icon name="bolt" size={17}/></button>
        <button className="icon-btn" aria-label="Notifications" onClick={() => setNotifOpen((o) => !o)} style={{ position: "relative" }}>
          <Icon name="bell" size={17}/>
          {unreadCount > 0 && (
            <span style={{ position: "absolute", top: 2, right: 2, minWidth: 14, height: 14, borderRadius: 99, background: "#EF4444", color: "#fff", fontSize: 9, fontWeight: 700, display: "grid", placeItems: "center", border: "1.5px solid var(--bg)", padding: "0 2px" }}>
              {unreadCount > 9 ? "9+" : unreadCount}
            </span>
          )}
        </button>
        {notifOpen && <NotificationsPopover onClose={() => setNotifOpen(false)}/>}
        <button className="icon-btn" style={{ marginLeft: 4 }} aria-label="Profile" onClick={() => nav("profile")}>
          <Avatar user={me} size="sm"/>
        </button>
      </div>
    </header>
  );
}
