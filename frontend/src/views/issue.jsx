// issue.jsx — Issue detail, Backlog table, My Issues
import { useState, useEffect, useRef, useMemo } from 'react';
import { Icon } from '../components/icons';
import { Avatar, Button, Badge, TypeIcon, PriorityBadge, StatusBadge, Modal, Empty, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api, useApi } from '../api/api';
import { adaptIssue, adaptComment, fmtDate } from '../api/adapters';
import { TYPE_META, PRIORITY_META } from '../store/data';
import { Pill, PillSelect } from './board';
import { IssueLinksPanel, AssigneesPanel, WatchVotePanel, TimeTrackingPanel, HistoryTab, AttachmentsPanel } from '../panels/issue';

export function IssueView({ nav, issueId }) {
  const { issues, setIssues, people, me, columns, activeProjectId: projectId, projects } = useApp();
  const issue = issues.find((i) => i.id === issueId) ?? null;
  const [tab, setTab]             = useState("comments");
  const [editingDesc, setEditingDesc] = useState(false);
  const [editingTitle, setEditingTitle] = useState(false);
  const [title, setTitle]         = useState(issue ? issue.title : "");
  const [desc, setDesc]           = useState("");
  const [comment, setComment]     = useState("");
  const [comments, setComments]   = useState([]);
  const [attachTrigger, setAttachTrigger] = useState(0);
  const [linkTrigger, setLinkTrigger]     = useState(0);
  const [subtaskOpen, setSubtaskOpen]     = useState(false);
  const [cloning, setCloning]             = useState(false);
  const toast = useToast();
  const proj = (projects || []).find((p) => p.id === (projectId || (issue && issue.project_id)));

  useEffect(() => {
    if (!issue || !issue._id) return;
    api("/issues/" + issue._id + "/comments").then((res) => {
      const items = res.items || res || [];
      setComments(items.map(adaptComment));
    }).catch(() => {});
    api("/issues/" + issue._id).then((detail) => {
      const adapted = adaptIssue(detail, proj ? proj.key : "");
      setDesc(detail.description || "");
      setTitle(adapted.title);
    }).catch(() => {});
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [issueId]);

  if (!issue) {
    if (!issues.length) return (
      <div style={{ padding: 40, textAlign: "center", color: "var(--text-muted)", fontSize: 14 }}>Loading…</div>
    );
    return (
      <div style={{ padding: 40 }}>
        <Empty icon="checkbox" title="Issue not found" hint="It may have been deleted or moved."/>
      </div>
    );
  }

  const assignee = issue.assigneeUser || (people || []).find((p) => p.id === issue.assignee);
  const reporter = issue.reporterUser || (people || []).find((p) => p.id === issue.reporter);

  async function saveTitle() {
    setEditingTitle(false);
    if (!issue._id || title === issue.title) return;
    try {
      await api("/issues/" + issue._id, { method: "PUT", body: { title } });
      setIssues((prev) => prev.map((i) => i.id === issue.id ? { ...i, title } : i));
    } catch (e) {
      toast && toast(e.message, { icon: "x", color: "#F87171" });
    }
  }

  async function cloneIssue() {
    if (!issue._id || cloning) return;
    setCloning(true);
    try {
      const cloned = await api("/issues/" + issue._id + "/clone", { method: "POST" });
      const adapted = adaptIssue(cloned, proj ? proj.key : "");
      setIssues((prev) => [adapted, ...prev]);
      toast && toast("Issue cloned: " + adapted.id);
      nav("issue/" + adapted.id);
    } catch (e) {
      toast && toast(e.message, { icon: "x", color: "#F87171" });
    } finally {
      setCloning(false);
    }
  }

  async function saveDesc() {
    setEditingDesc(false);
    if (!issue._id) return;
    try {
      await api("/issues/" + issue._id, { method: "PUT", body: { description: desc } });
    } catch (e) {
      toast && toast(e.message, { icon: "x", color: "#F87171" });
    }
  }

  async function postComment() {
    if (!comment.trim() || !issue._id) return;
    try {
      const created = await api("/issues/" + issue._id + "/comments", { body: { body: comment } });
      setComments((prev) => [...prev, adaptComment(created)]);
      setComment("");
    } catch (e) {
      toast && toast(e.message, { icon: "x", color: "#F87171" });
    }
  }

  return (
    <div className="detail">
      <div className="detail-main">
        <div className="crumbs row gap-2" style={{ marginBottom: 12 }}>
          <a href="#/projects" onClick={(e) => { e.preventDefault(); nav("projects"); }}>Projects</a>
          <Icon name="chevronRight" size={11}/>
          <a href="#/board" onClick={(e) => { e.preventDefault(); nav("board"); }}>{proj ? proj.name : "Board"}</a>
          <Icon name="chevronRight" size={11}/>
          <span className="mono" style={{ color: "var(--text-secondary)" }}>{issue.id}</span>
          <div className="row gap-1" style={{ marginLeft: "auto" }}>
            <Button variant="ghost" data-size="sm" icon="link" title="Copy link" onClick={() => { navigator.clipboard?.writeText(window.location.href); toast && toast("Link copied"); }}/>
            <Button variant="ghost" data-size="sm" icon="copy" title="Clone issue" disabled={cloning} onClick={cloneIssue}>{cloning ? "Cloning…" : "Clone"}</Button>
            <Button variant="ghost" data-size="sm" icon="star" title="Watch"/>
            <Button variant="ghost" data-size="sm" icon="moreH"/>
          </div>
        </div>

        <div className="row gap-2" style={{ marginBottom: 12 }}>
          <TypeIcon value={issue.type}/>
          <span className="mono text-sm muted">{issue.id}</span>
        </div>

        {editingTitle ? (
          <input className="input" autoFocus value={title} onChange={(e) => setTitle(e.target.value)}
            onBlur={saveTitle}
            onKeyDown={(e) => { if (e.key === "Enter") saveTitle(); if (e.key === "Escape") { setTitle(issue.title); setEditingTitle(false); } }}
            style={{ fontSize: 24, fontWeight: 600, padding: "6px 10px", marginBottom: 16 }}
          />
        ) : (
          <h1 onClick={() => setEditingTitle(true)} style={{ fontSize: 24, fontWeight: 600, letterSpacing: "-.015em", margin: "0 0 16px", cursor: "text", padding: "2px 0" }}>
            {title}
          </h1>
        )}

        <div className="row gap-2" style={{ marginBottom: 24 }}>
          <Button variant="secondary" icon="plus" onClick={() => setSubtaskOpen(true)}>Add subtask</Button>
          <Button variant="secondary" icon="link" onClick={() => setLinkTrigger((t) => t + 1)}>Link issue</Button>
          <Button variant="secondary" icon="paperclip" onClick={() => setAttachTrigger((t) => t + 1)}>Attach</Button>
        </div>

        <h3 style={{ fontSize: 13, fontWeight: 600, letterSpacing: ".02em", color: "var(--text-secondary)", textTransform: "uppercase", margin: "0 0 8px" }}>Description</h3>
        {editingDesc ? (
          <div style={{ marginBottom: 16 }}>
            <textarea className="textarea" value={desc} onChange={(e) => setDesc(e.target.value)} autoFocus style={{ minHeight: 120 }}/>
            <div className="row gap-2" style={{ marginTop: 8 }}>
              <Button variant="primary" onClick={saveDesc}>Save</Button>
              <Button onClick={() => setEditingDesc(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div onClick={() => setEditingDesc(true)} style={{ background: "var(--bg-subtle)", border: "1px solid var(--border)", borderRadius: 8, padding: "14px 16px", marginBottom: 24, cursor: "text", minHeight: 60 }}>
            {desc
              ? <p style={{ margin: 0, whiteSpace: "pre-wrap", lineHeight: 1.6 }}>{desc}</p>
              : <p style={{ margin: 0, color: "var(--text-muted)", fontSize: 13 }}>Click to add a description…</p>}
          </div>
        )}

        {/* Tabs */}
        <div style={{ borderBottom: "1px solid var(--border)", marginBottom: 16, marginLeft: -32, marginRight: -32, paddingLeft: 32, paddingRight: 32 }}>
          <div className="row gap-1">
            {[
              ["comments", "Comments", comments.length],
              ["activity", "Activity", null],
              ["history",  "History",  null],
            ].map(([id, label, ct]) => (
              <button key={id} className="tab" aria-selected={tab === id} onClick={() => setTab(id)}>
                {label} {ct != null && <span className="muted" style={{ marginLeft: 4, fontSize: 11.5 }}>{ct}</span>}
              </button>
            ))}
          </div>
        </div>

        {tab === "comments" && (
          <div>
            {comments.map((c) => {
              const u = c.author || (people || []).find((p) => p.id === c.who) || { name: "Unknown", initials: "?", color: "#94A3B8" };
              return (
                <div key={c.id} className="row gap-3" style={{ marginBottom: 16, alignItems: "flex-start" }}>
                  <Avatar user={u} size="md"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="row gap-2" style={{ marginBottom: 4 }}>
                      <span className="bold text-sm">{u.name}</span>
                      {/* FIX: use c.at (from adaptComment), not c.time */}
                      <span className="text-xs muted">{c.at}</span>
                    </div>
                    {/* FIX: use c.body (from adaptComment), not c.text */}
                    <div className="text-sm" style={{ lineHeight: 1.6 }}>
                      {(c.body || "").split(/(@\w+)/).map((part, i) =>
                        part.startsWith("@")
                          ? <span key={i} style={{ background: "var(--indigo-50)", color: "var(--indigo-700)", padding: "0 4px", borderRadius: 3, fontWeight: 500 }}>{part}</span>
                          : part
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
            <CommentBox value={comment} onChange={setComment} onSubmit={postComment} people={people} me={me}/>
          </div>
        )}

        {tab === "activity" && (
          <div style={{ marginBottom: 32 }}>
            <div className="text-sm secondary" style={{ padding: "20px 0", textAlign: "center" }}>Activity log coming soon.</div>
          </div>
        )}

        {tab === "history" && <HistoryTab issueId={issue._id}/>}
      </div>

      {/* RIGHT SIDEBAR */}
      <aside className="detail-aside">
        <div className="row gap-2" style={{ marginBottom: 16 }}>
          <PillSelect value={issue.status} onChange={async (v) => {
            setIssues((p) => p.map((i) => i.id === issue.id ? { ...i, status: v } : i));
            const col = (columns || []).find((c) => c.id === v);
            if (col && issue._id) {
              try { await api("/issues/" + issue._id, { method: "PUT", body: { status_id: col._id } }); }
              catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
            }
          }} options={(columns || []).map((c) => ({ id: c.id, label: c.label || c.id }))}/>
          <Button variant="secondary" icon="moreH" data-size="sm" style={{ padding: "0 8px" }}/>
        </div>

        <AssigneesPanel issueId={issue._id} people={people} me={me}/>
        <WatchVotePanel issueId={issue._id}/>
        <TimeTrackingPanel issueId={issue._id} me={me}/>
        <AttachmentsPanel issueId={issue._id} triggerUpload={attachTrigger}/>
        <IssueLinksPanel issueId={issue._id} issues={issues} triggerAdd={linkTrigger}/>

        <div style={{ marginBottom: 18 }}>
          <h4 style={{ fontSize: 11, fontWeight: 600, letterSpacing: ".06em", color: "var(--text-muted)", textTransform: "uppercase", margin: "0 0 10px" }}>Details</h4>
          <dl>
            <dt>Reporter</dt>
            <dd className="row gap-2">{reporter && <Avatar user={reporter} size="sm"/>}<span>{reporter?.name}</span></dd>
            <dt>Priority</dt>
            <dd><PriorityBadge value={issue.pri}/></dd>
            <dt>Type</dt>
            <dd><span className="row gap-2"><TypeIcon value={issue.type}/>{issue.type}</span></dd>
            <dt>Sprint</dt>
            <dd>{issue.sprint || "—"}</dd>
            <dt>Story points</dt>
            <dd>{issue.points}</dd>
            <dt>Due date</dt>
            {/* FIX: removed hardcoded ", 2024" */}
            <dd>{issue.due ? <span className="row gap-1"><Icon name="calendar" size={12} color="var(--text-muted)"/>{issue.due}</span> : "—"}</dd>
            <dt>Labels</dt>
            <dd className="row gap-1" style={{ flexWrap: "wrap" }}>{(issue.labels || []).map((l) => <span key={l} className="tag">{l}</span>)}</dd>
          </dl>
        </div>

        {assignee?.tg && (
          <div className="tg-card" style={{ padding: 12, marginBottom: 18 }}>
            <div className="row gap-2" style={{ marginBottom: 6 }}>
              <Icon name="telegram" size={14} color="#2AABEE"/>
              <span className="bold text-sm">Telegram notifications</span>
            </div>
            <div className="text-xs secondary" style={{ marginBottom: 8 }}>
              Assignee <span className="bold">{assignee.tg}</span> will receive updates.
            </div>
            <Button variant="ghost" data-size="sm" icon="bell" style={{ padding: 0, color: "var(--tg)" }}>Manage</Button>
          </div>
        )}
      </aside>

      <SubtaskModal
        open={subtaskOpen}
        onClose={() => setSubtaskOpen(false)}
        parentIssue={issue}
        projectId={projectId || issue.project_id}
        columns={columns}
        onCreated={(sub) => {
          setIssues((prev) => [sub, ...prev]);
          toast && toast("Subtask created: " + sub.title);
        }}
      />
    </div>
  );
}

function SubtaskModal({ open, onClose, parentIssue, projectId, columns, onCreated }) {
  const [title, setTitle]   = useState("");
  const [type, setType]     = useState("subtask");
  const [saving, setSaving] = useState(false);
  const toast = useToast();

  useEffect(() => { if (open) { setTitle(""); setType("subtask"); } }, [open]);

  async function create() {
    if (!title.trim() || !projectId) return;
    setSaving(true);
    try {
      const col = columns && columns[0];
      const body = {
        project_id: projectId,
        title: title.trim(),
        type,
        parent_id: parentIssue._id,
        ...(col ? { status_id: col._id } : {}),
      };
      const created = await api("/issues", { body });
      onCreated(adaptIssue(created, parentIssue.id.split("-")[0]));
      onClose();
    } catch (e) {
      toast && toast(e.message, { icon: "x", color: "#F87171" });
    } finally {
      setSaving(false);
    }
  }

  return (
    <Modal open={open} onClose={onClose} title={"Add subtask to " + (parentIssue ? parentIssue.id : "")}
      footer={<>
        <Button onClick={onClose}>Cancel</Button>
        <Button variant="primary" disabled={!title.trim() || saving} onClick={create}>
          {saving ? "Creating…" : "Create subtask"}
        </Button>
      </>}>
      <div style={{ marginBottom: 14 }}>
        <label className="label">Title</label>
        <input className="input" autoFocus value={title} onChange={(e) => setTitle(e.target.value)}
          placeholder="Subtask title…"
          onKeyDown={(e) => e.key === "Enter" && create()}/>
      </div>
      <div>
        <label className="label">Type</label>
        <select className="select" value={type} onChange={(e) => setType(e.target.value)}>
          <option value="subtask">Subtask</option>
          <option value="task">Task</option>
          <option value="bug">Bug</option>
        </select>
      </div>
    </Modal>
  );
}

function CommentBox({ value, onChange, onSubmit, people, me }) {
  const [mentioning, setMentioning] = useState(false);
  const [mentionQ, setMentionQ]     = useState("");

  useEffect(() => {
    const at = value.lastIndexOf("@");
    if (at >= 0 && (at === 0 || value[at - 1] === " ")) {
      const q = value.slice(at + 1);
      if (!q.includes(" ") && q.length < 20) {
        setMentioning(true); setMentionQ(q); return;
      }
    }
    setMentioning(false);
  }, [value]);

  const matches = mentioning ? (people || []).filter((p) => p.name.toLowerCase().includes(mentionQ.toLowerCase())).slice(0, 5) : [];

  function pick(u) {
    const at = value.lastIndexOf("@");
    onChange(value.slice(0, at) + "@" + u.name.split(" ")[0].toLowerCase() + " ");
    setMentioning(false);
  }

  return (
    <div className="row gap-3" style={{ alignItems: "flex-start", marginBottom: 32, position: "relative" }}>
      <Avatar user={me} size="md"/>
      <div style={{ flex: 1, border: "1px solid var(--border)", borderRadius: 8, overflow: "hidden", background: "var(--bg)" }}>
        <textarea className="textarea" style={{ border: 0, borderRadius: 0, minHeight: 80 }}
          placeholder="Add a comment… Use @ to mention" value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={(e) => { if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) onSubmit(); }}
        />
        <div className="row" style={{ padding: 6, borderTop: "1px solid var(--border)", background: "var(--bg-subtle)" }}>
          <div className="row gap-1">
            {[["bold","Bold"],["italic","Italic"],["code","Code"],["link","Link"],["picture","Image"],["at","Mention"]].map(([ic, name]) => (
              <button key={name} className="icon-btn" style={{ width: 26, height: 26 }} title={name}><Icon name={ic} size={13}/></button>
            ))}
          </div>
          <div style={{ marginLeft: "auto" }} className="row gap-2">
            <span className="text-xs muted"><span className="kbd">⌘</span><span className="kbd" style={{ marginLeft: 2 }}>↵</span> to submit</span>
            <Button variant="primary" data-size="sm" icon="send" onClick={onSubmit} disabled={!value.trim()}>Comment</Button>
          </div>
        </div>

        {mentioning && matches.length > 0 && (
          <div style={{ position: "absolute", top: 80, left: 56, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 30, width: 260 }}>
            <div className="text-xs muted" style={{ padding: "4px 8px" }}>Mention a member</div>
            {matches.map((u) => (
              <button key={u.id} onClick={() => pick(u)} className="nav-item" style={{ color: "var(--text)" }}>
                <Avatar user={u} size="sm"/>
                <span className="text-sm">{u.name}</span>
                <span className="text-xs muted" style={{ marginLeft: "auto" }}>{u.role}</span>
              </button>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

// ─── Backlog ─────────────────────────────────────────────
const EMPTY_FILTER = { status: "all", assignee: "all", priority: "all", type: "all", label: "all", sprint: "all" };

export function BacklogView({ nav }) {
  const { issues, setIssues, columns, people, activeProjectId: projectId, projects, issuesTotal, issuesLoading, loadMoreIssues } = useApp();
  const [editing, setEditing]   = useState(null);
  const [editVal, setEditVal]   = useState("");
  const [filter, setFilter]     = useState(EMPTY_FILTER);
  const [activeSaved, setActiveSaved] = useState(null);
  const { data: saved, reload: reloadSaved, setData: setSaved } = useApi("/saved-filters");
  const [saveOpen, setSaveOpen] = useState(false);
  const [savedOpen, setSavedOpen] = useState(false);
  const [name, setName]         = useState("");
  const savedRef = useRef(null);
  const toast = useToast();
  const proj = (projects || []).find((p) => p.id === projectId);

  useEffect(() => {
    if (!savedOpen) return;
    const h = (e) => { if (savedRef.current && !savedRef.current.contains(e.target)) setSavedOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [savedOpen]);

  const labels  = useMemo(() => Array.from(new Set(issues.flatMap((i) => i.labels))), [issues]);
  const sprints = useMemo(() => Array.from(new Set(issues.map((i) => i.sprint).filter(Boolean))), [issues]);

  function set(k, v) { setFilter((f) => ({ ...f, [k]: v })); setActiveSaved(null); }

  const filtered = issues.filter((i) => {
    if (filter.status !== "all" && i.status !== filter.status) return false;
    if (filter.assignee !== "all" && (filter.assignee === "none" ? i.assignee : i.assignee !== filter.assignee)) return false;
    if (filter.priority !== "all" && i.pri !== filter.priority) return false;
    if (filter.type !== "all" && i.type !== filter.type) return false;
    if (filter.label !== "all" && !i.labels.includes(filter.label)) return false;
    if (filter.sprint !== "all" && i.sprint !== filter.sprint) return false;
    return true;
  });

  const dirty = JSON.stringify(filter) !== JSON.stringify(EMPTY_FILTER);

  function applySaved(sf) {
    setFilter({ ...EMPTY_FILTER, ...(sf.filter_json || {}) });
    setActiveSaved(sf.id); setSavedOpen(false);
  }
  async function saveFilter() {
    if (!name.trim()) return;
    const filter_json = {};
    Object.entries(filter).forEach(([k, v]) => { if (v !== "all") filter_json[k] = v; });
    try {
      const sf = await api("/saved-filters", { method: "POST", body: { name, filter_json } });
      toast && toast("Filter saved");
      setSaveOpen(false); setName(""); setActiveSaved(sf.id); reloadSaved();
    } catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function delSaved(id, e) {
    e.stopPropagation();
    setSaved((saved || []).filter((s) => s.id !== id));
    if (activeSaved === id) setActiveSaved(null);
    try { await api("/saved-filters/" + id, { method: "DELETE" }); toast && toast("Filter deleted"); }
    catch (err) { toast && toast(err.message, { icon: "x", color: "#F87171" }); reloadSaved(); }
  }

  async function saveInlineEdit(i) {
    if (!editVal.trim() || !i._id) { setEditing(null); return; }
    try {
      await api("/issues/" + i._id, { method: "PUT", body: { title: editVal } });
      setIssues((prev) => prev.map((x) => x.id === i.id ? { ...x, title: editVal } : x));
    } catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
    setEditing(null);
  }

  const activeName = activeSaved && (saved || []).find((s) => s.id === activeSaved);

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">
            <a href="#/board" onClick={(e) => { e.preventDefault(); nav("board"); }}>{proj ? proj.name : "Project"}</a>
            <Icon name="chevronRight" size={11}/> <span>Backlog</span>
          </div>
          <h1>Backlog</h1>
          <p>{filtered.length} of {issues.length} issues · drag to plan into upcoming sprints.</p>
        </div>
        <div className="row gap-2">
          <Button icon="download">Export CSV</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>Create</Button>
        </div>
      </div>

      <div className="row gap-2" style={{ padding: "0 32px 14px", flexWrap: "wrap" }}>
        <Pill label="Status" icon="kanban" value={filter.status === "all" ? "All" : filter.status} options={[{ id: "all", label: "All" }, ...(columns || []).map((c) => ({ id: c.id, label: c.label || c.id }))]} onChange={(v) => set("status", v)}/>
        <Pill label="Assignee" icon="user" value={filter.assignee === "all" ? "All" : filter.assignee === "none" ? "Unassigned" : (people || []).find((p) => p.id === filter.assignee)?.name} options={[{ id: "all", label: "All" }, { id: "none", label: "Unassigned" }, ...(people || []).map((p) => ({ id: p.id, label: p.name }))]} onChange={(v) => set("assignee", v)}/>
        <Pill label="Priority" icon="flag" value={filter.priority === "all" ? "All" : filter.priority} options={[{ id: "all", label: "All" }, ...Object.keys(PRIORITY_META).map((p) => ({ id: p, label: p }))]} onChange={(v) => set("priority", v)}/>
        <Pill label="Type" icon="tag" value={filter.type === "all" ? "All" : filter.type} options={[{ id: "all", label: "All" }, ...Object.keys(TYPE_META).map((t) => ({ id: t, label: t }))]} onChange={(v) => set("type", v)}/>
        <Pill label="Label" icon="tag" value={filter.label === "all" ? "All" : filter.label} options={[{ id: "all", label: "All" }, ...labels.map((l) => ({ id: l, label: l }))]} onChange={(v) => set("label", v)}/>
        <Pill label="Sprint" icon="rocket" value={filter.sprint === "all" ? "All" : filter.sprint} options={[{ id: "all", label: "All" }, ...sprints.map((s) => ({ id: s, label: s }))]} onChange={(v) => set("sprint", v)}/>
        <div style={{ flex: 1 }}/>
        {activeName && (
          <span className="badge" data-tone="info" style={{ height: 28, padding: "0 8px", gap: 6 }}>
            <Icon name="star" size={12}/> {activeName.name}
            <button className="icon-btn" style={{ width: 18, height: 18 }} onClick={() => { setFilter(EMPTY_FILTER); setActiveSaved(null); }}><Icon name="x" size={11}/></button>
          </span>
        )}
        <span ref={savedRef} style={{ position: "relative" }}>
          <Button data-size="sm" icon="star" iconRight="chevronDown" onClick={() => setSavedOpen((o) => !o)}>Saved filters</Button>
          {savedOpen && (
            <div style={{ position: "absolute", top: "calc(100% + 4px)", right: 0, minWidth: 240, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50 }}>
              {(saved || []).length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>No saved filters yet.</div>}
              {(saved || []).map((sf) => (
                <div key={sf.id} className="nav-item" style={{ color: "var(--text)", cursor: "default" }} onClick={() => applySaved(sf)}>
                  <Icon name="star" size={13} color={activeSaved === sf.id ? "var(--indigo-600)" : "var(--text-muted)"}/>
                  <span className="text-sm" style={{ flex: 1 }}>{sf.name}</span>
                  <button className="icon-btn" style={{ width: 20, height: 20 }} onClick={(e) => delSaved(sf.id, e)} title="Delete"><Icon name="trash" size={12}/></button>
                </div>
              ))}
            </div>
          )}
        </span>
        <Button data-size="sm" variant="secondary" icon="star" disabled={!dirty} onClick={() => setSaveOpen(true)}>Save filter</Button>
      </div>

      <Modal open={saveOpen} onClose={() => setSaveOpen(false)} title="Save filter"
        footer={<><Button onClick={() => setSaveOpen(false)}>Cancel</Button><Button variant="primary" disabled={!name.trim()} onClick={saveFilter}>Save</Button></>}>
        <label className="label">Filter name</label>
        <input className="input" autoFocus placeholder="e.g. My critical bugs" value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => { if (e.key === "Enter") saveFilter(); }}/>
      </Modal>

      <div style={{ padding: "0 32px 32px" }}>
        <div style={{ border: "1px solid var(--border)", borderRadius: 10, overflow: "hidden", background: "var(--bg)" }}>
          <table className="table">
            <thead>
              <tr>
                <th style={{ width: 28 }}/>
                <th style={{ width: 100 }}>Key</th>
                <th>Title</th>
                <th style={{ width: 110 }}>Status</th>
                <th style={{ width: 100 }}>Priority</th>
                <th style={{ width: 60 }}>Pts</th>
                <th style={{ width: 100 }}>Due</th>
                <th style={{ width: 50 }}>Assignee</th>
              </tr>
            </thead>
            <tbody>
              {filtered.length === 0 ? (
                <tr>
                  <td colSpan={8} style={{ padding: 0, border: 0 }}>
                    <Empty icon="inbox"
                      title={dirty ? "No issues match the current filter" : "Backlog is empty"}
                      hint={dirty ? "Try changing or clearing the filters." : "Create your first issue to get started."}
                    />
                  </td>
                </tr>
              ) : filtered.map((i) => {
                const u = i.assigneeUser || (people || []).find((p) => p.id === i.assignee);
                return (
                  <tr key={i.id}>
                    <td><TypeIcon value={i.type}/></td>
                    <td><a href={"/issue/" + i.id} className="mono text-xs" onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.id}</a></td>
                    <td onClick={() => { setEditing(i.id); setEditVal(i.title); }}>
                      {editing === i.id ? (
                        <input className="input" autoFocus value={editVal}
                          onChange={(e) => setEditVal(e.target.value)}
                          onBlur={() => saveInlineEdit(i)}
                          onKeyDown={(e) => { if (e.key === "Enter") saveInlineEdit(i); if (e.key === "Escape") setEditing(null); }}
                          style={{ padding: "2px 6px" }}/>
                      ) : (
                        <a href={"/issue/" + i.id} onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.title}</a>
                      )}
                    </td>
                    <td><StatusBadge value={i.status}/></td>
                    <td><PriorityBadge value={i.pri}/></td>
                    <td><span className="mono text-xs">{i.points}</span></td>
                    <td className="text-xs muted">{i.due || "—"}</td>
                    <td>{u && <Avatar user={u} size="sm"/>}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        {issues.length < issuesTotal && (
          <div style={{ textAlign: "center", paddingTop: 16 }}>
            <Button icon={issuesLoading ? undefined : "chevronDown"} disabled={issuesLoading} onClick={loadMoreIssues}>
              {issuesLoading ? "Loading…" : `Load more  ·  ${issues.length} / ${issuesTotal}`}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

// ─── My issues ───────────────────────────────────────────
export function MyIssuesView({ nav }) {
  const { issues, me } = useApp();
  const mine = issues.filter((i) => me && i.assignee === me.id);
  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">You <Icon name="chevronRight" size={11}/> <span>My issues</span></div>
          <h1>My issues</h1>
          <p>{mine.length} issues assigned to you across all projects.</p>
        </div>
      </div>
      <div style={{ padding: "0 32px 32px" }}>
        {mine.length === 0 ? (
          <Empty icon="checkbox" title="No issues assigned to you" hint="Issues assigned to you will appear here."/>
        ) : (
          <div style={{ border: "1px solid var(--border)", borderRadius: 10, overflow: "hidden", background: "var(--bg)" }}>
            <table className="table">
              <thead>
                <tr>
                  <th style={{ width: 28 }}/>
                  <th style={{ width: 100 }}>Key</th>
                  <th>Title</th>
                  <th style={{ width: 130 }}>Status</th>
                  <th style={{ width: 100 }}>Priority</th>
                  <th style={{ width: 100 }}>Due</th>
                </tr>
              </thead>
              <tbody>
                {mine.map((i) => (
                  <tr key={i.id}>
                    <td><TypeIcon value={i.type}/></td>
                    <td><a href={"/issue/" + i.id} className="mono text-xs" onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.id}</a></td>
                    <td><a href={"/issue/" + i.id} onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.title}</a></td>
                    <td><StatusBadge value={i.status}/></td>
                    <td><PriorityBadge value={i.pri}/></td>
                    <td className="text-xs muted">{i.due || "—"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
