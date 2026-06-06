// panels_issue.jsx — Issue detail aside panels (links, watchers, votes,
// assignees, time tracking) + History tab. Features 1–6.

const LINK_TYPES = ["blocks", "is blocked by", "relates to", "duplicates", "clones"];
const LINK_GROUP_ORDER = ["is blocked by", "blocks", "relates to", "duplicates", "is duplicated by", "clones", "is cloned by"];

function AsideSection({ title, count, action, children }) {
  return (
    <div style={{ marginBottom: 18 }}>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 10 }}>
        <h4 style={{ fontSize: 11, fontWeight: 600, letterSpacing: ".06em", color: "var(--text-muted)", textTransform: "uppercase", margin: 0 }}>
          {title}{count != null && <span style={{ marginLeft: 6, color: "var(--text-muted)" }}>{count}</span>}
        </h4>
        {action}
      </div>
      {children}
    </div>
  );
}

function MiniSpinner({ size = 14 }) {
  return (
    <span style={{
      width: size, height: size, borderRadius: "50%",
      border: "2px solid var(--border)", borderTopColor: "var(--indigo-600)",
      display: "inline-block", animation: "forge-spin .7s linear infinite",
    }}/>
  );
}
// inject keyframes once
if (!document.getElementById("forge-spin-kf")) {
  const s = document.createElement("style"); s.id = "forge-spin-kf";
  s.textContent = "@keyframes forge-spin{to{transform:rotate(360deg)}}";
  document.head.appendChild(s);
}

// ─── Feature 1: Issue Links ───────────────────────────────────────────
function IssueLinksPanel({ issueId }) {
  const { data: links, loading, reload, setData } = useApi("/issues/" + issueId + "/links");
  const [adding, setAdding] = React.useState(false);
  const [linkType, setLinkType] = React.useState("blocks");
  const [targetId, setTargetId] = React.useState("");
  const [busy, setBusy] = React.useState(false);
  const toast = useToast();

  const others = FORGE_DATA.ISSUES.filter((i) => i.id !== issueId);

  async function add() {
    if (!targetId) return;
    setBusy(true);
    try {
      await api("/issues/" + issueId + "/links", { method: "POST", body: { linked_issue_id: targetId, link_type: linkType } });
      toast("Issue linked");
      setAdding(false); setTargetId("");
      reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    setBusy(false);
  }
  async function remove(id) {
    setData((links || []).filter((l) => l.id !== id));
    try { await api("/issues/links/" + id, { method: "DELETE" }); toast("Link removed"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }

  const grouped = {};
  (links || []).forEach((l) => { (grouped[l.link_type] ||= []).push(l); });
  const groupKeys = LINK_GROUP_ORDER.filter((k) => grouped[k]);

  return (
    <AsideSection title="Issue links" count={links ? links.length : null}
      action={<button className="icon-btn" style={{ width: 22, height: 22 }} title="Link issue" onClick={() => setAdding((a) => !a)}><Icon name={adding ? "x" : "plus"} size={13}/></button>}>
      {loading && <div className="row gap-2 text-xs muted"><MiniSpinner/> Loading links…</div>}
      {!loading && groupKeys.length === 0 && !adding && (
        <div className="text-xs muted" style={{ padding: "4px 0" }}>No linked issues yet.</div>
      )}
      <div className="stack gap-3">
        {groupKeys.map((k) => (
          <div key={k}>
            <div className="text-xs muted" style={{ textTransform: "capitalize", marginBottom: 4 }}>{k}</div>
            <div className="stack gap-1">
              {grouped[k].map((l) => (
                <div key={l.id} className="row gap-2" style={{ padding: "6px 8px", background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
                  <TypeIcon value={l.issue.type}/>
                  <a href={"#/issue/" + l.issue.id} className="mono text-xs" style={{ color: "var(--indigo-600)", textDecoration: "none" }}>{l.issue.id}</a>
                  <span className="text-xs" style={{ flex: 1, minWidth: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }} title={l.issue.title}>{l.issue.title}</span>
                  <StatusBadge value={l.issue.status}/>
                  <button className="icon-btn" style={{ width: 20, height: 20 }} title="Remove link" onClick={() => remove(l.id)}><Icon name="x" size={12}/></button>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {adding && (
        <div style={{ marginTop: 10, padding: 10, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8 }}>
          <select className="select" value={linkType} onChange={(e) => setLinkType(e.target.value)} style={{ marginBottom: 8 }}>
            {LINK_TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
          </select>
          <select className="select" value={targetId} onChange={(e) => setTargetId(e.target.value)} style={{ marginBottom: 8 }}>
            <option value="">Select an issue…</option>
            {others.map((i) => <option key={i.id} value={i.id}>{i.id} — {i.title.slice(0, 40)}</option>)}
          </select>
          <div className="row gap-2">
            <Button variant="primary" data-size="sm" disabled={!targetId || busy} onClick={add}>{busy ? "Linking…" : "Add link"}</Button>
            <Button data-size="sm" onClick={() => setAdding(false)}>Cancel</Button>
          </div>
        </div>
      )}
    </AsideSection>
  );
}

// ─── Feature 4: Multiple Assignees ────────────────────────────────────
function AssigneesPanel({ issueId }) {
  const { data: assignees, loading, reload, setData } = useApi("/issues/" + issueId + "/assignees");
  const [picking, setPicking] = React.useState(false);
  const ref = React.useRef(null);
  const toast = useToast();

  React.useEffect(() => {
    if (!picking) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setPicking(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [picking]);

  const ids = (assignees || []).map((u) => u.id);
  const candidates = FORGE_DATA.PEOPLE.filter((p) => !ids.includes(p.id));

  async function addUser(uid) {
    const next = [...ids, uid];
    setData(FORGE_DATA.PEOPLE.filter((p) => next.includes(p.id)));
    setPicking(false);
    try { await api("/issues/" + issueId + "/assignees", { method: "PUT", body: { user_ids: next } }); toast("Assignee added"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }
  async function removeUser(uid) {
    setData((assignees || []).filter((u) => u.id !== uid));
    try { await api("/issues/" + issueId + "/assignees/" + uid, { method: "DELETE" }); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }

  return (
    <AsideSection title="Assignees" count={assignees ? assignees.length : null}>
      {loading ? <div className="row gap-2 text-xs muted"><MiniSpinner/> Loading…</div> : (
        <div className="row gap-2" style={{ flexWrap: "wrap" }}>
          {(assignees || []).map((u) => (
            <button key={u.id} onClick={() => removeUser(u.id)} title={"Remove " + u.name}
              className="row gap-2" style={{ padding: "3px 8px 3px 3px", border: "1px solid var(--border)", borderRadius: 99, background: "var(--bg)" }}>
              <Avatar user={u} size="sm"/>
              <span className="text-xs">{u.name.split(" ")[0]}</span>
              <Icon name="x" size={11} color="var(--text-muted)"/>
            </button>
          ))}
          <span ref={ref} style={{ position: "relative" }}>
            <button className="icon-btn" style={{ width: 28, height: 28, border: "1px dashed var(--border-strong)", borderRadius: "50%" }} title="Add assignee" onClick={() => setPicking((p) => !p)}>
              <Icon name="plus" size={14}/>
            </button>
            {picking && (
              <div style={{ position: "absolute", top: "calc(100% + 4px)", left: 0, minWidth: 220, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, maxHeight: 260, overflowY: "auto" }}>
                {candidates.length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>Everyone is assigned.</div>}
                {candidates.map((u) => (
                  <button key={u.id} className="nav-item" style={{ color: "var(--text)" }} onClick={() => addUser(u.id)}>
                    <Avatar user={u} size="sm"/><span className="text-sm">{u.name}</span>
                    <span className="text-xs muted" style={{ marginLeft: "auto" }}>{u.role}</span>
                  </button>
                ))}
              </div>
            )}
          </span>
        </div>
      )}
    </AsideSection>
  );
}

// ─── Feature 2 + 3: Watchers and Votes ────────────────────────────────
function WatchVotePanel({ issueId }) {
  const me = FORGE_DATA.ME;
  const { data: watchers, reload: reloadW, setData: setW } = useApi("/issues/" + issueId + "/watchers");
  const { data: votes, setData: setV } = useApi("/issues/" + issueId + "/votes");
  const toast = useToast();
  const watching = (watchers || []).some((u) => u.id === me.id);

  async function toggleWatch() {
    try {
      if (watching) { setW((watchers || []).filter((u) => u.id !== me.id)); await api("/issues/" + issueId + "/watchers", { method: "DELETE" }); }
      else { setW([...(watchers || []), me]); await api("/issues/" + issueId + "/watchers", { method: "POST", body: { user_id: me.id } }); toast("Watching this issue"); }
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reloadW(); }
  }
  async function toggleVote() {
    try { const r = await api("/issues/" + issueId + "/votes", { method: "POST" }); setV(r); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="row gap-2" style={{ marginBottom: 18 }}>
      <div style={{ flex: 1, border: "1px solid var(--border)", borderRadius: 8, padding: 10, background: "var(--bg)" }}>
        <div className="row" style={{ justifyContent: "space-between", marginBottom: 8 }}>
          <span className="text-xs muted">Watchers</span>
          <button className="btn btn-ghost" data-size="sm" style={{ padding: "0 6px", height: 22, color: watching ? "var(--indigo-600)" : "var(--text-secondary)" }} onClick={toggleWatch}>
            <Icon name={watching ? "eye" : "eyeOff"} size={13}/> {watching ? "Watching" : "Watch"}
          </button>
        </div>
        {watchers ? (
          watchers.length ? <AvatarStack users={watchers} max={5} size="sm"/> : <span className="text-xs muted">No watchers</span>
        ) : <MiniSpinner/>}
      </div>
      <button onClick={toggleVote} title="Vote for this issue"
        style={{ border: "1px solid " + (votes && votes.voted_by_me ? "var(--indigo-600)" : "var(--border)"), background: votes && votes.voted_by_me ? "var(--indigo-50)" : "var(--bg)", borderRadius: 8, padding: "8px 12px", display: "grid", placeItems: "center", minWidth: 64 }}>
        <Icon name="arrowUp" size={16} color={votes && votes.voted_by_me ? "var(--indigo-600)" : "var(--text-secondary)"} strokeWidth={2.4}/>
        <span className="bold" style={{ fontSize: 15, color: votes && votes.voted_by_me ? "var(--indigo-700)" : "var(--text)" }}>{votes ? votes.count : "·"}</span>
        <span className="text-xs muted">votes</span>
      </button>
    </div>
  );
}

// ─── Feature 5: Time Tracking ─────────────────────────────────────────
function fmtHM(min) {
  if (min == null) return "0h";
  const h = Math.floor(min / 60), m = min % 60;
  return (h ? h + "h" : "") + (m ? " " + m + "m" : "") || "0h";
}
function fmtDate(iso) {
  if (!iso) return "";
  const d = new Date(iso), now = new Date();
  const diff = (now - d) / 1000;
  if (diff < 3600) return Math.max(1, Math.floor(diff / 60)) + "m ago";
  if (diff < 86400) return Math.floor(diff / 3600) + "h ago";
  if (diff < 86400 * 7) return Math.floor(diff / 86400) + "d ago";
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

function TimeTrackingPanel({ issueId }) {
  const { data: summary, reload: reloadSum } = useApi("/issues/" + issueId + "/time-summary");
  const { data: logs, reload: reloadLogs } = useApi("/issues/" + issueId + "/worklogs");
  const [modal, setModal] = React.useState(false);
  const [editing, setEditing] = React.useState(null);
  const toast = useToast();

  function refresh() { reloadSum(); reloadLogs(); }

  async function del(id) {
    try { await api("/issues/" + issueId + "/worklogs/" + id, { method: "DELETE" }); toast("Worklog deleted"); refresh(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  const logged = summary ? summary.time_spent : 0;
  const orig = summary ? summary.original_estimate : 0;
  const pct = orig ? Math.min(100, (logged / orig) * 100) : 0;
  const over = orig && logged > orig;

  return (
    <AsideSection title="Time tracking"
      action={<button className="btn btn-ghost" data-size="sm" style={{ padding: "0 6px", height: 22, color: "var(--indigo-600)" }} onClick={() => { setEditing(null); setModal(true); }}><Icon name="plus" size={13}/> Log</button>}>
      {!summary ? <div className="row gap-2 text-xs muted"><MiniSpinner/> Loading…</div> : (
        <div>
          <div style={{ height: 8, background: "var(--bg-muted)", borderRadius: 99, overflow: "hidden", marginBottom: 6 }}>
            <div style={{ width: pct + "%", height: "100%", background: over ? "var(--danger)" : "var(--indigo-600)" }}/>
          </div>
          <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
            <span className="text-xs"><span className="bold">{fmtHM(logged)}</span> <span className="muted">logged</span></span>
            <span className="text-xs"><span className="bold">{fmtHM(summary.time_remaining)}</span> <span className="muted">remaining</span></span>
          </div>
          <div className="text-xs muted" style={{ marginBottom: 8 }}>Estimate {fmtHM(orig)}</div>
          <div className="stack gap-2">
            {(logs || []).map((w) => {
              const u = w.user || FORGE_DATA.PEOPLE.find((p) => p.id === w.user_id);
              return (
                <div key={w.id} className="row gap-2" style={{ alignItems: "flex-start", padding: "8px", background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
                  <Avatar user={u} size="sm"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="row gap-2">
                      <span className="text-xs bold">{w.time_spent}</span>
                      <span className="text-xs muted">{fmtDate(w.started_at)}</span>
                    </div>
                    {w.description && <div className="text-xs secondary" style={{ marginTop: 2 }}>{w.description}</div>}
                  </div>
                  <div className="row gap-1">
                    <button className="icon-btn" style={{ width: 20, height: 20 }} title="Edit" onClick={() => { setEditing(w); setModal(true); }}><Icon name="pencil" size={11}/></button>
                    <button className="icon-btn" style={{ width: 20, height: 20 }} title="Delete" onClick={() => del(w.id)}><Icon name="trash" size={11}/></button>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}
      <LogTimeModal open={modal} onClose={() => setModal(false)} issueId={issueId} editing={editing} onSaved={refresh}/>
    </AsideSection>
  );
}

function LogTimeModal({ open, onClose, issueId, editing, onSaved }) {
  const [timeSpent, setTimeSpent] = React.useState("");
  const [desc, setDesc] = React.useState("");
  const [date, setDate] = React.useState("");
  const [busy, setBusy] = React.useState(false);
  const toast = useToast();

  React.useEffect(() => {
    if (!open) return;
    setTimeSpent(editing ? editing.time_spent : "");
    setDesc(editing ? editing.description : "");
    setDate(editing && editing.started_at ? editing.started_at.slice(0, 10) : new Date().toISOString().slice(0, 10));
  }, [open, editing]);

  async function save() {
    if (!timeSpent.trim()) return;
    setBusy(true);
    const body = { time_spent: timeSpent, description: desc, started_at: new Date(date).toISOString() };
    try {
      if (editing) await api("/issues/" + issueId + "/worklogs/" + editing.id, { method: "PUT", body });
      else await api("/issues/" + issueId + "/worklogs", { method: "POST", body });
      toast(editing ? "Worklog updated" : "Time logged");
      onClose(); onSaved();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    setBusy(false);
  }

  return (
    <Modal open={open} onClose={onClose} title={editing ? "Edit worklog" : "Log time"}
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="primary" disabled={!timeSpent.trim() || busy} onClick={save}>{busy ? "Saving…" : "Save"}</Button></>}>
      <div className="stack gap-3">
        <div>
          <label className="label">Time spent</label>
          <input className="input" placeholder="e.g. 2h 30m" value={timeSpent} onChange={(e) => setTimeSpent(e.target.value)} autoFocus/>
          <div className="help">Use the format 2h 30m, 45m, or 1h.</div>
        </div>
        <div>
          <label className="label">Date started</label>
          <input className="input" type="date" value={date} onChange={(e) => setDate(e.target.value)}/>
        </div>
        <div>
          <label className="label">Description</label>
          <textarea className="textarea" placeholder="What did you work on?" value={desc} onChange={(e) => setDesc(e.target.value)}/>
        </div>
      </div>
    </Modal>
  );
}

// ─── Feature 6: History / Changelog tab ───────────────────────────────
function HistoryTab({ issueId }) {
  const { data: history, loading, error } = useApi("/issues/" + issueId + "/history");
  if (loading) return <div className="row gap-2 text-sm muted" style={{ padding: 16 }}><MiniSpinner/> Loading history…</div>;
  if (error) return <div className="text-sm" style={{ color: "var(--danger)", padding: 16 }}>{error}</div>;
  if (!history || history.length === 0) return <Empty icon="history" title="No history yet" hint="Field changes will appear here."/>;
  return (
    <div style={{ marginBottom: 32 }}>
      {history.map((h, i) => {
        const u = h.user || FORGE_DATA.PEOPLE.find((p) => p.id === h.user_id);
        return (
          <div key={h.id} className="row gap-3" style={{ padding: "10px 0", borderBottom: i < history.length - 1 ? "1px solid var(--border)" : 0, alignItems: "flex-start" }}>
            <Avatar user={u} size="sm"/>
            <div style={{ flex: 1 }}>
              <div className="text-sm">
                <span className="bold">{u.name}</span>{" "}
                {h.from == null ? (
                  <span className="secondary">{h.field === "Issue" ? "created this issue" : <>set <span className="medium" style={{ color: "var(--text)" }}>{h.field}</span> to <Badge tone="muted">{h.to}</Badge></>}</span>
                ) : (
                  <span className="secondary">changed <span className="medium" style={{ color: "var(--text)" }}>{h.field}</span> from <Badge tone="muted">{h.from}</Badge> <Icon name="arrowRight" size={11} style={{ verticalAlign: "middle" }}/> <Badge tone="info">{h.to}</Badge></span>
                )}
              </div>
              <div className="text-xs muted" style={{ marginTop: 3 }}>{fmtDate(h.at)}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

Object.assign(window, { AsideSection, MiniSpinner, IssueLinksPanel, AssigneesPanel, WatchVotePanel, TimeTrackingPanel, LogTimeModal, HistoryTab, fmtHM, fmtDate });
