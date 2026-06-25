import { useState, useEffect, Component } from 'react';
import { AppProvider, useApp } from './store/AppContext';
import { Icon } from './components/icons';
import { ToastProvider, useHashRoute } from './components/components';
import { useTweaks, TweaksPanel, TweakSection, TweakToggle } from './components/tweaks-panel';
import { api, clearTokens } from './api/api';
import { adaptIssue, toBackendType, toBackendPriority } from './api/adapters';
import { Sidebar, Topbar, CommandPalette, useKeyboardShortcuts, KeyboardShortcutsModal } from './components/shell';
import { LoginView, AcceptInviteView } from './views/auth';
import { DashboardView } from './views/dashboard';
import { BoardView, CreateIssueModal } from './views/board';
import { IssueView, BacklogView, MyIssuesView } from './views/issue';
import { SprintView, ProjectsView, ReportsView, SettingsView, ProfileView } from './views/misc';
import { RoadmapView } from './views/roadmap';
import { WikiView } from './views/wiki';
import { MembersView, NotificationSettingsView, IntegrationsView } from './views/members';
import { ApiKeysPanel, AuditLogPanel, ImportPanel, UsersPanel } from './panels/admin';

class ErrorBoundary extends Component {
  state = { error: null };
  static getDerivedStateFromError(error) { return { error }; }
  render() {
    if (this.state.error) {
      return (
        <div style={{ padding: 48, textAlign: "center", color: "var(--text-secondary)" }}>
          <div style={{ fontSize: 32, marginBottom: 12 }}>⚠</div>
          <div style={{ fontWeight: 600, marginBottom: 8, color: "var(--text)" }}>Something went wrong</div>
          <div style={{ fontSize: 13, marginBottom: 20 }}>{this.state.error.message}</div>
          <button
            onClick={() => this.setState({ error: null })}
            style={{ padding: "8px 20px", background: "var(--indigo-600)", color: "#fff", border: "none", borderRadius: "var(--radius)", cursor: "pointer", fontSize: 13 }}
          >
            Try again
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}

const TWEAK_DEFAULTS = { dark: false };

export function App() {
  return (
    <ToastProvider>
      <AppProvider>
        <AppInner/>
      </AppProvider>
    </ToastProvider>
  );
}

function AppInner() {
  const { authReady, me, projects, issues, setIssues, columns, activeProjectId, people, switchProject, clearSession } = useApp();
  const [t, setTweak] = useTweaks(TWEAK_DEFAULTS);
  const [route, nav] = useHashRoute();
  const [cmdOpen, setCmdOpen] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);
  const [wikiTarget, setWikiTarget] = useState(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);

  const openCommand = () => setCmdOpen(true);
  const { shortcutsOpen, setShortcutsOpen } = useKeyboardShortcuts(nav, openCommand);

  useEffect(() => {
    document.documentElement.dataset.theme = t.dark ? "dark" : "light";
  }, [t.dark]);

  useEffect(() => {
    const h = () => setCreateOpen(true);
    window.addEventListener("forge:create", h);
    return () => window.removeEventListener("forge:create", h);
  }, []);

  useEffect(() => {
    const h = () => setShortcutsOpen(true);
    window.addEventListener("forge:shortcuts", h);
    return () => window.removeEventListener("forge:shortcuts", h);
  }, [setShortcutsOpen]);

  useEffect(() => {
    const h = () => setSidebarOpen((o) => !o);
    window.addEventListener("forge:mobile-menu", h);
    return () => window.removeEventListener("forge:mobile-menu", h);
  }, []);

  async function handleCreate(form) {
    if (!activeProjectId) return;
    try {
      const proj = projects.find((p) => p.id === activeProjectId);
      const body = {
        project_id:   activeProjectId,
        title:        form.title,
        type:         toBackendType(form.type),
        priority:     toBackendPriority(form.pri),
        assignee_id:  form.assignee || undefined,
        story_points: form.points ? parseInt(form.points) : undefined,
      };
      if (form.status) {
        const col = columns.find((c) => c.id === form.status);
        if (col) body.status_id = col._id;
      }
      const created = await api("/issues", { body });
      const adapted = adaptIssue(created, proj ? proj.key : "");
      const map = {};
      (columns || []).forEach((c) => { map[c._id] = c.id; });
      setIssues((prev) => [{ ...adapted, status: map[adapted.status_id] || adapted.status }, ...prev]);
    } catch (e) {
      console.error("Create issue failed:", e.message);
    }
  }

  if (!authReady) {
    return (
      <div style={{ position: "fixed", inset: 0, display: "grid", placeItems: "center", background: "var(--bg)" }}>
        <div style={{ color: "var(--text-muted)", fontSize: 14 }}>Loading…</div>
      </div>
    );
  }

  if (route.path === "accept-invite") {
    return <AcceptInviteView nav={nav}/>;
  }

  if (!me || route.path === "login" || route.path === "register" || route.path === "forgot") {
    return (
      <LoginView
        nav={nav}
        mode={route.path === "login" || route.path === "register" || route.path === "forgot" ? route.path : "login"}
      />
    );
  }

  return (
    <>
      <div className="app">
        <div className={`mobile-overlay${sidebarOpen ? " visible" : ""}`} onClick={() => setSidebarOpen(false)}/>
        <Sidebar route={route} nav={nav} mobileOpen={sidebarOpen} onClose={() => setSidebarOpen(false)}/>
        <Topbar  route={route} nav={nav} onOpenCommand={() => setCmdOpen(true)}/>
        <main className="main">
          <ErrorBoundary key={route.path}>
            {route.path === "dashboard"          && <DashboardView nav={nav}/>}
            {route.path === "my-issues"          && <MyIssuesView nav={nav}/>}
            {route.path === "projects"           && <ProjectsView nav={nav}/>}
            {route.path === "board"              && <BoardView nav={nav}/>}
            {route.path === "backlog"            && <BacklogView nav={nav}/>}
            {route.path === "sprint"             && <SprintView nav={nav}/>}
            {route.path === "wiki"               && <WikiView nav={nav} target={wikiTarget} onTargetConsumed={() => setWikiTarget(null)}/>}
            {route.path === "members"            && <MembersView nav={nav}/>}
            {route.path === "roadmap"            && <RoadmapView nav={nav}/>}
            {route.path === "reports"            && <ReportsView nav={nav}/>}
            {route.path === "settings"           && <SettingsView nav={nav}/>}
            {route.path === "integrations"       && <IntegrationsView nav={nav}/>}
            {route.path === "notifications-page" && <NotificationSettingsView nav={nav}/>}
            {route.path === "admin"              && <AdminView nav={nav}/>}
            {route.path === "issue"              && <IssueView nav={nav} issueId={route.rest[0]}/>}
            {route.path === "profile"            && <ProfileView nav={nav}/>}
          </ErrorBoundary>
        </main>

        <CommandPalette open={cmdOpen} onClose={() => setCmdOpen(false)} nav={nav}
          onWikiPage={(spaceId, pageId) => { setWikiTarget({ spaceId, pageId }); nav("wiki"); }}/>
        <CreateIssueModal open={createOpen} onClose={() => setCreateOpen(false)} onCreate={handleCreate}/>
        <KeyboardShortcutsModal open={shortcutsOpen} onClose={() => setShortcutsOpen(false)}/>
      </div>

      <TweaksPanel title="Tweaks">
        <TweakSection label="Appearance"/>
        <TweakToggle label="Dark mode" value={t.dark} onChange={(v) => setTweak("dark", v)}/>
        <TweakSection label="Quick jump"/>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 6 }}>
          {[["dashboard","Dashboard"],["board","Board"],["members","Members"],["wiki","Wiki"],["sprint","Sprints"],["reports","Reports"],["settings","Settings"],["admin","Admin"],["login","Log out"]].map(([to, label]) => (
            <button key={to} onClick={() => { if (to === "login") clearSession(); else nav(to); }}
              style={{ padding: "5px 8px", border: "0.5px solid rgba(0,0,0,.08)", background: "rgba(255,255,255,.5)", borderRadius: 6, fontSize: 11.5, color: "#29261b", cursor: "default" }}>
              {label}
            </button>
          ))}
        </div>
      </TweaksPanel>
    </>
  );
}

function AdminView({ nav }) {
  const { me } = useApp();
  const [tab, setTab] = useState("users");

  useEffect(() => {
    if (me && me.role !== "admin") nav("dashboard");
  }, [me]);

  if (!me || me.role !== "admin") return null;

  return (
    <div>
      <div className="page-head" style={{ paddingBottom: 0 }}>
        <div>
          <div className="crumbs">Forge <Icon name="chevronRight" size={11}/> <span>Admin</span></div>
          <h1>Admin</h1>
          <p>Workspace-wide settings, user management, access tokens, and data migration.</p>
        </div>
      </div>
      <div className="tabs">
        {[["users","Users"],["apikeys","API keys"],["audit","Audit log"],["import","Import"]].map(([id, label]) => (
          <button key={id} className="tab" aria-selected={tab === id} onClick={() => setTab(id)}>{label}</button>
        ))}
      </div>
      <div style={{ padding: "24px 32px 40px", maxWidth: 1100 }}>
        {tab === "users"   && <UsersPanel/>}
        {tab === "apikeys" && <ApiKeysPanel/>}
        {tab === "audit"   && <AuditLogPanel/>}
        {tab === "import"  && <ImportPanel/>}
      </div>
    </div>
  );
}
