// panels/issue.jsx — Issue detail aside panels: links, assignees, watchers/votes, time tracking, history
import { useState, useEffect, useRef } from 'react';
import { Icon } from '../components/icons';
import { Avatar, AvatarStack, Badge, Button, Modal, TypeIcon, StatusBadge, Empty, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api, apiUpload, useApi } from '../api/api';
import { adaptUser } from '../api/adapters';

const LINK_TYPES = ["blocks", "is blocked by", "relates to", "duplicates", "clones"];
const LINK_GROUP_ORDER = ["is blocked by", "blocks", "relates to", "duplicates", "is duplicated by", "clones", "is cloned by"];

// Spinner keyframes injected once
const _style = document.createElement("style");
_style.textContent = "@keyframes forge-spin{to{transform:rotate(360deg)}}";
document.head.appendChild(_style);

export function MiniSpinner({ size = 14 }) {
  return (
    <span style={{
      width: size, height: size, borderRadius: "50%",
      border: "2px solid var(--border)", borderTopColor: "var(--indigo-600)",
      display: "inline-block", animation: "forge-spin .7s linear infinite",
    }}/>
  );
}

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

// ─── Issue Links ──────────────────────────────────────────
export function IssueLinksPanel({ issueId, issues, triggerAdd }) {
  const { data: links, loading, reload, setData } = useApi(issueId ? "/issues/" + issueId + "/links" : null);
  const [adding, setAdding] = useState(false);

  useEffect(() => { if (triggerAdd > 0) setAdding(true); }, [triggerAdd]);
  const [linkType, setLinkType] = useState("blocks");
  const [targetId, setTargetId] = useState("");
  const [busy, setBusy] = useState(false);
  const toast = useToast();

  const others = (issues || []).filter((i) => i._id !== issueId);

  async function add() {
    if (!targetId) return;
    setBusy(true);
    try {
      await api("/issues/" + issueId + "/links", { method: "POST", body: { linked_issue_id: targetId, link_type: linkType } });
      toast && toast("Issue linked");
      setAdding(false); setTargetId("");
      reload();
    } catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
    setBusy(false);
  }
  async function remove(id) {
    setData((links || []).filter((l) => l.id !== id));
    try { await api("/issues/links/" + id, { method: "DELETE" }); toast && toast("Link removed"); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }

  const grouped = {};
  (links || []).forEach((l) => { (grouped[l.link_type] ||= []).push(l); });
  const groupKeys = LINK_GROUP_ORDER.filter((k) => grouped[k]);

  return (
    <AsideSection title="Issue links" count={links ? links.length : null}
      action={<button className="icon-btn" style={{ width: 22, height: 22 }} title="Link issue" onClick={() => setAdding((a) => !a)}><Icon name={adding ? "x" : "plus"} size={13}/></button>}>
      {loading && <div className="row gap-2 text-xs muted"><MiniSpinner/> Loading links…</div>}
      {!loading && groupKeys.length === 0 && !adding && <div className="text-xs muted" style={{ padding: "4px 0" }}>No linked issues yet.</div>}
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
                  <button className="icon-btn" style={{ width: 20, height: 20 }} title="Remove" onClick={() => remove(l.id)}><Icon name="x" size={12}/></button>
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
            {others.map((i) => <option key={i._id || i.id} value={i._id || i.id}>{i.id} — {i.title.slice(0, 40)}</option>)}
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

// ─── Assignees Panel ──────────────────────────────────────
export function AssigneesPanel({ issueId, people, me }) {
  const { data: assignees, loading, reload, setData } = useApi(issueId ? "/issues/" + issueId + "/assignees" : null);
  const [picking, setPicking] = useState(false);
  const ref = useRef(null);
  const toast = useToast();

  useEffect(() => {
    if (!picking) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setPicking(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [picking]);

  const ids = (assignees || []).map((u) => u.id);
  const candidates = (people || []).filter((p) => !ids.includes(p.id));

  async function addUser(uid) {
    const next = [...ids, uid];
    setData((people || []).filter((p) => next.includes(p.id)));
    setPicking(false);
    try { await api("/issues/" + issueId + "/assignees", { method: "PUT", body: { user_ids: next } }); toast && toast("Assignee added"); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }
  async function removeUser(uid) {
    setData((assignees || []).filter((u) => u.id !== uid));
    try { await api("/issues/" + issueId + "/assignees/" + uid, { method: "DELETE" }); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
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

// ─── Watch / Vote ─────────────────────────────────────────
export function WatchVotePanel({ issueId }) {
  const { me } = useApp();
  const { data: watchers, reload: reloadW, setData: setW } = useApi(issueId ? "/issues/" + issueId + "/watchers" : null);
  const { data: votes, setData: setV } = useApi(issueId ? "/issues/" + issueId + "/votes" : null);
  const toast = useToast();
  const watching = me && me.id && (watchers || []).some((u) => u.id === me.id);

  async function toggleWatch() {
    try {
      if (watching) { setW((watchers || []).filter((u) => u.id !== me.id)); await api("/issues/" + issueId + "/watchers", { method: "DELETE" }); }
      else { setW([...(watchers || []), me]); await api("/issues/" + issueId + "/watchers", { method: "POST", body: { user_id: me.id } }); toast && toast("Watching this issue"); }
    } catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); reloadW(); }
  }
  async function toggleVote() {
    try { const r = await api("/issues/" + issueId + "/votes", { method: "POST" }); setV(r); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
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
        {watchers ? (watchers.length ? <AvatarStack users={watchers} max={5} size="sm"/> : <span className="text-xs muted">No watchers</span>) : <MiniSpinner/>}
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

// ─── Time Tracking ────────────────────────────────────────
function fmtHM(min) {
  if (min == null) return "0h";
  const h = Math.floor(min / 60), m = min % 60;
  return (h ? h + "h" : "") + (m ? " " + m + "m" : "") || "0h";
}
function relTime(iso) {
  if (!iso) return "";
  const d = new Date(iso), now = new Date();
  const diff = (now - d) / 1000;
  if (diff < 3600) return Math.max(1, Math.floor(diff / 60)) + "m ago";
  if (diff < 86400) return Math.floor(diff / 3600) + "h ago";
  if (diff < 86400 * 7) return Math.floor(diff / 86400) + "d ago";
  return d.toLocaleDateString(undefined, { month: "short", day: "numeric" });
}

export function TimeTrackingPanel({ issueId, me }) {
  const { people } = useApp();
  const { data: summary, reload: reloadSum } = useApi(issueId ? "/issues/" + issueId + "/time-summary" : null);
  const { data: logs,    reload: reloadLogs } = useApi(issueId ? "/issues/" + issueId + "/worklogs" : null);
  const [modal, setModal]     = useState(false);
  const [editing, setEditing] = useState(null);
  const toast = useToast();

  function refresh() { reloadSum(); reloadLogs(); }

  async function del(id) {
    try { await api("/issues/" + issueId + "/worklogs/" + id, { method: "DELETE" }); toast && toast("Worklog deleted"); refresh(); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  const logged = summary ? summary.time_spent : 0;
  const orig   = summary ? summary.original_estimate : 0;
  const pct    = orig ? Math.min(100, (logged / orig) * 100) : 0;
  const over   = orig && logged > orig;

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
              const u = w.user ? adaptUser(w.user) : (people || []).find((p) => p.id === w.user_id) || { name: "Unknown", initials: "?", color: "#94A3B8" };
              return (
                <div key={w.id} className="row gap-2" style={{ alignItems: "flex-start", padding: 8, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
                  <Avatar user={u} size="sm"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="row gap-2">
                      <span className="text-xs bold">{w.time_spent}</span>
                      <span className="text-xs muted">{relTime(w.started_at)}</span>
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
  const [timeSpent, setTimeSpent] = useState("");
  const [desc, setDesc]           = useState("");
  const [date, setDate]           = useState("");
  const [busy, setBusy]           = useState(false);
  const toast = useToast();

  useEffect(() => {
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
      toast && toast(editing ? "Worklog updated" : "Time logged");
      onClose(); onSaved();
    } catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
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

// ─── History tab ─────────────────────────────────────────
export function HistoryTab({ issueId }) {
  const { data: history, loading, error } = useApi(issueId ? "/issues/" + issueId + "/history" : null);
  if (loading) return <div className="row gap-2 text-sm muted" style={{ padding: 16 }}><MiniSpinner/> Loading history…</div>;
  if (error) return <div className="text-sm" style={{ color: "var(--danger)", padding: 16 }}>{error}</div>;
  if (!history || history.length === 0) return <Empty icon="history" title="No history yet" hint="Field changes will appear here."/>;
  return (
    <div style={{ marginBottom: 32 }}>
      {history.map((h, i) => {
        const u = h.user ? adaptUser(h.user) : { name: "Unknown", initials: "?", color: "#94A3B8" };
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
              <div className="text-xs muted" style={{ marginTop: 3 }}>{relTime(h.at)}</div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

// ─── Attachments Panel ────────────────────────────────────
export function AttachmentsPanel({ issueId, triggerUpload }) {
  const { data: attachments, loading, reload } = useApi(issueId ? `/issues/${issueId}/attachments` : null);
  const [uploading, setUploading] = useState(false);
  const [lightbox, setLightbox] = useState(null);
  const inputRef = useRef(null);
  const toast = useToast();

  useEffect(() => { if (triggerUpload > 0) inputRef.current?.click(); }, [triggerUpload]);

  async function handleFiles(files) {
    if (!files || !files.length) return;
    setUploading(true);
    for (const file of Array.from(files)) {
      try {
        await apiUpload(`/issues/${issueId}/attachments`, file);
      } catch (e) {
        toast(e.message, { icon: "x", color: "#F87171" });
      }
    }
    setUploading(false);
    reload();
  }

  async function remove(id) {
    try {
      await api(`/attachments/${id}`, { method: "DELETE" });
      reload();
    } catch (e) {
      toast(e.message, { icon: "x", color: "#F87171" });
    }
  }

  function isImage(m) { return m?.startsWith("image/"); }
  function isVideo(m) { return m?.startsWith("video/"); }
  function fmtSize(b) {
    if (b < 1024) return b + " B";
    if (b < 1048576) return (b / 1024).toFixed(1) + " KB";
    return (b / 1048576).toFixed(1) + " MB";
  }

  const list = Array.isArray(attachments) ? attachments : [];

  return (
    <>
      <AsideSection title="Attachments" count={list.length || undefined}
        action={
          <button className="icon-btn" title="Attach file" disabled={uploading}
            onClick={() => inputRef.current?.click()}>
            {uploading ? <MiniSpinner size={13}/> : <Icon name="plus" size={13}/>}
          </button>
        }>

        <input ref={inputRef} type="file" multiple style={{ display: "none" }}
          onChange={(e) => { handleFiles(e.target.files); e.target.value = ""; }}/>

        {loading && <div className="row gap-2 text-sm muted"><MiniSpinner size={12}/> Loading…</div>}

        {!loading && list.length === 0 && (
          <div className="text-sm muted" style={{ padding: "8px 0" }}>No attachments yet.</div>
        )}

        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 6, marginTop: list.length ? 4 : 0 }}>
          {list.map((a) => (
            <div key={a.id} className="attachment-card" style={{
              position: "relative", borderRadius: 6, overflow: "hidden",
              border: "1px solid var(--border)", background: "var(--bg-raised)",
            }}>
              {isImage(a.mime_type) ? (
                <img src={a.download_url} alt={a.file_name}
                  onClick={() => setLightbox({ url: a.download_url, type: "image" })}
                  style={{ width: "100%", height: 72, objectFit: "cover", cursor: "zoom-in", display: "block" }}/>
              ) : isVideo(a.mime_type) ? (
                <div onClick={() => setLightbox({ url: a.download_url, type: "video" })}
                  style={{ width: "100%", height: 72, background: "#111", cursor: "pointer",
                    display: "grid", placeItems: "center" }}>
                  <Icon name="play" size={24} color="#fff"/>
                </div>
              ) : (
                <div style={{ padding: "8px 6px", display: "flex", alignItems: "center", gap: 5, minHeight: 56 }}>
                  <Icon name="attachment" size={16} color="var(--text-muted)"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="text-xs bold" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{a.file_name}</div>
                    <div className="text-xs muted">{fmtSize(a.file_size)}</div>
                  </div>
                </div>
              )}

              <div style={{ position: "absolute", top: 3, right: 3, display: "flex", gap: 2 }}>
                <a href={a.download_url} download={a.file_name} target="_blank" rel="noreferrer"
                  style={{ background: "rgba(0,0,0,.6)", borderRadius: 4, padding: "2px 5px", display: "flex", alignItems: "center" }}>
                  <Icon name="download" size={10} color="#fff"/>
                </a>
                <button onClick={() => remove(a.id)}
                  style={{ background: "rgba(0,0,0,.6)", borderRadius: 4, padding: "2px 5px", border: "none", cursor: "pointer", display: "flex", alignItems: "center" }}>
                  <Icon name="x" size={10} color="#fff"/>
                </button>
              </div>

              {(isImage(a.mime_type) || isVideo(a.mime_type)) && (
                <div className="text-xs muted" style={{ padding: "2px 5px", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap", fontSize: 10 }}>
                  {a.file_name}
                </div>
              )}
            </div>
          ))}
        </div>
      </AsideSection>

      {lightbox && (
        <div onClick={() => setLightbox(null)} style={{
          position: "fixed", inset: 0, background: "rgba(0,0,0,.88)", zIndex: 9999,
          display: "grid", placeItems: "center", cursor: "zoom-out",
        }}>
          {lightbox.type === "image"
            ? <img src={lightbox.url} onClick={(e) => e.stopPropagation()}
                style={{ maxWidth: "90vw", maxHeight: "90vh", borderRadius: 8, cursor: "default" }}/>
            : <video src={lightbox.url} controls autoPlay onClick={(e) => e.stopPropagation()}
                style={{ maxWidth: "90vw", maxHeight: "90vh", borderRadius: 8 }}/>
          }
          <button onClick={() => setLightbox(null)}
            style={{ position: "fixed", top: 18, right: 18, background: "rgba(255,255,255,.12)",
              border: "none", borderRadius: "50%", width: 36, height: 36, cursor: "pointer",
              color: "#fff", fontSize: 20, display: "grid", placeItems: "center" }}>×</button>
        </div>
      )}
    </>
  );
}
