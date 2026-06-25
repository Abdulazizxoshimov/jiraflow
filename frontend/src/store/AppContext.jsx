// AppContext.jsx — Centralized app state: auth, projects, issues, people.
import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { api, getToken, clearTokens } from '../api/api';
import { adaptUser, adaptProject, adaptIssue, adaptStatus, adaptList } from '../api/adapters';
import { ws } from '../lib/ws';

const AppCtx = createContext(null);
export const useApp = () => useContext(AppCtx);

export function AppProvider({ children }) {
  const [authReady, setAuthReady]             = useState(false);
  const [me, setMe]                           = useState(null);
  const [projects, setProjects]               = useState([]);
  const [projectsLoading, setProjectsLoading] = useState(false);
  const [activeProjectId, setActiveProjectId] = useState(null);
  const [people, setPeople]                   = useState([]);
  const [issues, setIssues]                   = useState([]);
  const [issuesTotal, setIssuesTotal]         = useState(0);
  const [issuesLoading, setIssuesLoading]     = useState(false);
  const [columns, setColumns]                 = useState([]);

  // ── helpers ──────────────────────────────────────────────────────────────
  const loadColumns = useCallback(async (projId, projectList) => {
    try {
      const list = projectList || projects;
      const proj = list.find((p) => p.id === projId);
      if (!proj || !proj._raw || !proj._raw.workflow_id) return [];
      const wf = await api("/workflows/" + proj._raw.workflow_id);
      const cols = (wf.statuses || []).map(adaptStatus);
      setColumns(cols);
      return cols;
    } catch { return []; }
  }, [projects]);

  const loadIssues = useCallback(async (projId, projectList, cols) => {
    try {
      setIssuesLoading(true);
      const list = projectList || projects;
      const proj = list.find((p) => p.id === projId);
      const projKey = proj ? proj.key : "";
      const res = await api("/issues?project_id=" + projId + "&limit=100");
      const items = res.items || res || [];
      const total = res.total ?? items.length;
      const adapted = adaptList(items, adaptIssue, projKey);
      const colMap = {};
      (cols || []).forEach((c) => { colMap[c._id] = c.id; });
      const mapped = adapted.map((i) => ({ ...i, status: colMap[i.status_id] || i.status }));
      setIssues(mapped);
      setIssuesTotal(total);
    } catch {
      setIssues([]);
    } finally {
      setIssuesLoading(false);
    }
  }, [projects]);

  // Load more issues (cursor pagination via offset)
  const loadMoreIssues = useCallback(async () => {
    if (issuesLoading || issues.length >= issuesTotal || !activeProjectId) return;
    try {
      setIssuesLoading(true);
      const proj = projects.find((p) => p.id === activeProjectId);
      const projKey = proj ? proj.key : "";
      const res = await api("/issues?project_id=" + activeProjectId + "&limit=100&offset=" + issues.length);
      const items = res.items || res || [];
      const adapted = adaptList(items, adaptIssue, projKey);
      const colMap = {};
      columns.forEach((c) => { colMap[c._id] = c.id; });
      const mapped = adapted.map((i) => ({ ...i, status: colMap[i.status_id] || i.status }));
      setIssues((prev) => [...prev, ...mapped]);
      setIssuesTotal(res.total ?? issuesTotal);
    } catch {} finally {
      setIssuesLoading(false);
    }
  }, [issuesLoading, issues.length, issuesTotal, activeProjectId, projects, columns]);

  // ── shared session loader (eliminates duplication between boot & onLoggedIn) ─
  const loadSession = useCallback(async () => {
    const rawMe = await api("/auth/me");
    const adapted = adaptUser(rawMe);
    setMe(adapted);

    setProjectsLoading(true);
    const projRes = await api("/projects");
    const projectList = adaptList(projRes.items || projRes, adaptProject);
    setProjects(projectList);
    setProjectsLoading(false);

    const savedId = localStorage.getItem("active_project_id");
    const activeProjId = (savedId && projectList.find((p) => p.id === savedId))
      ? savedId
      : (projectList[0] ? projectList[0].id : null);
    setActiveProjectId(activeProjId);

    const usersRes = await api("/users?limit=100");
    setPeople(adaptList(usersRes.items || usersRes, adaptUser));

    if (activeProjId) {
      const cols = await loadColumns(activeProjId, projectList);
      await loadIssues(activeProjId, projectList, cols || []);
    }
    return adapted;
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // ── boot ─────────────────────────────────────────────────────────────────
  useEffect(() => {
    async function boot() {
      const token = getToken();
      if (!token) { setAuthReady(true); return; }
      try {
        await loadSession();
        ws.connect(token);
      } catch {
        setProjectsLoading(false);
        setIssuesLoading(false);
        clearTokens();
      }
      setAuthReady(true);
    }
    boot();
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // ── actions ───────────────────────────────────────────────────────────────
  const switchProject = useCallback(async (projId) => {
    // Unsubscribe from old project room, subscribe to new one
    if (activeProjectId) ws.unsubscribe(`project:${activeProjectId}`);
    ws.subscribe(`project:${projId}`);
    setActiveProjectId(projId);
    localStorage.setItem("active_project_id", projId);
    setIssues([]);
    setIssuesTotal(0);
    const cols = await loadColumns(projId, projects);
    await loadIssues(projId, projects, cols || []);
  }, [activeProjectId, projects, loadColumns, loadIssues]);

  const clearSession = useCallback(() => {
    ws.disconnect();
    clearTokens();
    setMe(null);
    setProjects([]);
    setIssues([]);
    setIssuesTotal(0);
    setColumns([]);
    setPeople([]);
    window.location.hash = "#/login";
  }, []);

  const onLoggedIn = useCallback(async () => {
    try {
      await loadSession();
      ws.connect(getToken());
    } catch {}
  }, [loadSession]);

  const value = {
    authReady,
    me,
    projects,
    projectsLoading,
    activeProjectId,
    people,
    issues,   setIssues,
    issuesTotal,
    issuesLoading,
    columns,
    setPeople,
    switchProject,
    clearSession,
    onLoggedIn,
    loadMoreIssues,
    reloadIssues: () => loadIssues(activeProjectId, projects, columns),
  };

  return <AppCtx.Provider value={value}>{children}</AppCtx.Provider>;
}
