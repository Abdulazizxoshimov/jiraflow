// app.jsx — Forge root: router, state, tweaks

const TWEAK_DEFAULTS = /*EDITMODE-BEGIN*/{
  "dark": false,
  "telegramConnected": true
}/*EDITMODE-END*/;

function App() {
  const [t, setTweak] = useTweaks(TWEAK_DEFAULTS);
  const [route, nav] = useHashRoute();
  const [issues, setIssues] = React.useState(FORGE_DATA.ISSUES);
  const [people, setPeople] = React.useState(FORGE_DATA.PEOPLE);
  const [cmdOpen, setCmdOpen] = React.useState(false);
  const [createOpen, setCreateOpen] = React.useState(false);

  // apply theme
  React.useEffect(() => {
    document.documentElement.dataset.theme = t.dark ? "dark" : "light";
  }, [t.dark]);

  // global hotkeys
  React.useEffect(() => {
    const h = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setCmdOpen(true);
      } else if (e.key === "c" && !e.metaKey && !e.ctrlKey && !e.altKey) {
        const tag = (e.target.tagName || "").toLowerCase();
        if (tag !== "input" && tag !== "textarea" && !e.target.isContentEditable) {
          e.preventDefault();
          setCreateOpen(true);
        }
      }
    };
    window.addEventListener("keydown", h);
    return () => window.removeEventListener("keydown", h);
  }, []);

  // listen for create button anywhere
  React.useEffect(() => {
    const h = () => setCreateOpen(true);
    window.addEventListener("forge:create", h);
    return () => window.removeEventListener("forge:create", h);
  }, []);

  function handleCreate(form) {
    const next = "INFRA-" + (245 + issues.length - FORGE_DATA.ISSUES.length);
    setIssues((p) => [...p, {
      id: next, title: form.title, type: form.type, pri: form.pri,
      status: form.status, assignee: form.assignee, reporter: "u1",
      points: form.points, due: form.due, labels: form.labels ? form.labels.split(",").map((s) => s.trim()).filter(Boolean) : [],
      sprint: "Sprint 24", comments: 0, sub: 0
    }]);
  }

  // ─── unauth view ──────────────────────────────────────
  if (route.path === "login" || route.path === "register" || route.path === "forgot") {
    return (
      <ToastProvider>
        <LoginView nav={nav} mode={route.path}/>
      </ToastProvider>
    );
  }

  return (
    <ToastProvider>
      <div className="app">
        <Sidebar route={route} nav={nav} tg={t.telegramConnected}/>
        <Topbar route={route} nav={nav} onOpenCommand={() => setCmdOpen(true)} tg={t.telegramConnected}/>
        <main className="main">
          {route.path === "dashboard"          && <DashboardView nav={nav} tg={t.telegramConnected}/>}
          {route.path === "my-issues"          && <MyIssuesView nav={nav} issues={issues}/>}
          {route.path === "projects"           && <ProjectsView nav={nav}/>}
          {route.path === "board"              && <BoardView nav={nav} issues={issues} setIssues={setIssues}/>}
          {route.path === "backlog"            && <BacklogView nav={nav} issues={issues}/>}
          {route.path === "sprint"             && <SprintView nav={nav} issues={issues}/>}
          {route.path === "wiki"               && <WikiView nav={nav}/>}
          {route.path === "members"            && <MembersView nav={nav} tg={t.telegramConnected} people={people} setPeople={setPeople}/>}
          {route.path === "reports"            && <ReportsView nav={nav} issues={issues}/>}
          {route.path === "settings"           && <SettingsView nav={nav} tg={t.telegramConnected}/>}
          {route.path === "integrations"       && <IntegrationsView nav={nav} tg={t.telegramConnected} setTg={(v) => setTweak("telegramConnected", v)}/>}
          {route.path === "notifications-page" && <NotificationSettingsView nav={nav} tg={t.telegramConnected}/>}
          {route.path === "admin"              && <AdminView nav={nav}/>}
          {route.path === "issue"              && <IssueView nav={nav} issueId={route.rest[0]} issues={issues} setIssues={setIssues}/>}
        </main>

        <CommandPalette open={cmdOpen} onClose={() => setCmdOpen(false)} nav={nav}/>
        <CreateIssueModal open={createOpen} onClose={() => setCreateOpen(false)} onCreate={handleCreate}/>
      </div>

      <TweaksPanel title="Tweaks">
        <TweakSection label="Appearance"/>
        <TweakToggle label="Dark mode" value={t.dark} onChange={(v) => setTweak("dark", v)}/>
        <TweakSection label="Telegram bot"/>
        <TweakToggle label="Bot connected" value={t.telegramConnected} onChange={(v) => setTweak("telegramConnected", v)}/>
        <div style={{ fontSize: 11, color: "rgba(41,38,27,.55)", lineHeight: 1.5, padding: "2px 2px 6px" }}>
          Toggle to preview how Forge looks with the Telegram bot in either state — members table, notification settings preview, and integrations card all react.
        </div>
        <TweakSection label="Quick jump"/>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 6 }}>
          {[["dashboard","Dashboard"],["board","Board"],["issue/INFRA-232","Issue"],["members","Members"],["notifications-page","Notif. settings"],["integrations","Integrations"],["wiki","Wiki"],["sprint","Sprints"],["reports","Reports"],["login","Login"]].map(([to, label]) => (
            <button key={to} onClick={() => nav(to)} style={{ padding: "5px 8px", border: "0.5px solid rgba(0,0,0,.08)", background: "rgba(255,255,255,.5)", borderRadius: 6, fontSize: 11.5, color: "#29261b", cursor: "default" }}>
              {label}
            </button>
          ))}
        </div>
      </TweaksPanel>
    </ToastProvider>
  );
}

// ─── Admin (Features 12, 13, 26) ────────────────────────
function AdminView({ nav }) {
  const [tab, setTab] = React.useState("apikeys");
  return (
    <div>
      <div className="page-head" style={{ paddingBottom: 0 }}>
        <div>
          <div className="crumbs">Forge <Icon name="chevronRight" size={11}/> <span>Admin</span></div>
          <h1>Admin</h1>
          <p>Workspace-wide settings, access tokens, and data migration.</p>
        </div>
      </div>
      <div className="tabs">
        {[["apikeys", "API keys"], ["audit", "Audit log"], ["import", "Import"]].map(([id, label]) => (
          <button key={id} className="tab" aria-selected={tab === id} onClick={() => setTab(id)}>{label}</button>
        ))}
      </div>
      <div style={{ padding: "24px 32px 40px", maxWidth: 1100 }}>
        {tab === "apikeys" && <ApiKeysPanel/>}
        {tab === "audit" && <AuditLogPanel/>}
        {tab === "import" && <ImportPanel/>}
      </div>
    </div>
  );
}

window.App = App;

// Mount
const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(<App/>);
