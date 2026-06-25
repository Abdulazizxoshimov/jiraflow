// board.jsx — Kanban board with drag-and-drop + swimlanes
import { useState, useEffect, useRef, useMemo } from 'react';
import { Icon } from '../components/icons';
import { Avatar, AvatarStack, Button, Modal, TypeIcon, PriorityBadge, StatusBadge, Empty, Skeleton, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api } from '../api/api';
import { TYPE_META, PRIORITY_META } from '../store/data';
import { ws } from '../lib/ws';

const SWIMLANE_OPTIONS = [
  { id: "none",     label: "No swimlanes" },
  { id: "assignee", label: "By Assignee" },
  { id: "epic",     label: "By Epic" },
];
const EPIC_PALETTE = ["#6366F1","#10B981","#F59E0B","#EF4444","#8B5CF6","#06B6D4","#EC4899","#F97316"];

// Derive a color from column category/tone
const TONE_COLOR = {
  muted:   "#94A3B8",
  info:    "#3B82F6",
  success: "#10B981",
  warning: "#F59E0B",
  danger:  "#EF4444",
  purple:  "#8B5CF6",
};
function colColor(col) {
  if (col.tone && TONE_COLOR[col.tone]) return TONE_COLOR[col.tone];
  return "#94A3B8";
}

export function BoardView({ nav }) {
  const { issues, setIssues, columns, people, activeProjectId: projectId, projects, issuesLoading } = useApp();
  const [filter, setFilter]   = useState({ assignee: "all", type: "all", search: "" });
  const [dragging, setDragging] = useState(null);
  const [overCol, setOverCol] = useState(null);
  const [swimlane, setSwimlane] = useState("none");
  const [boardId, setBoardId] = useState(null);
  const toast = useToast();

  // Sync columns → issue statuses whenever columns update
  useEffect(() => {
    if (!columns || columns.length === 0) return;
    const map = {};
    columns.forEach((c) => { map[c._id] = c.id; });
    setIssues((prev) => prev.map((i) => ({ ...i, status: map[i.status_id] || i.status })));
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [columns]);

  // Keep a ref to latest columns so WS handlers always see current columns
  // without re-registering handlers on every column change.
  const columnsRef = useRef(columns);
  useEffect(() => { columnsRef.current = columns; }, [columns]);

  // Real-time: listen for issue.moved / issue.updated from other users
  useEffect(() => {
    if (!projectId) return;
    const offMoved = ws.on('issue.moved', (msg) => {
      const p = msg.payload;
      if (!p || p.project_id !== projectId) return;
      const col = (columnsRef.current || []).find((c) => c._id === p.status_id);
      if (col) {
        setIssues((prev) => prev.map((i) =>
          i._id === p.id ? { ...i, status: col.id, status_id: p.status_id } : i
        ));
      }
    });
    const offUpdated = ws.on('issue.updated', (msg) => {
      const p = msg.payload;
      if (!p || p.project_id !== projectId) return;
      const col = (columnsRef.current || []).find((c) => c._id === p.status_id);
      if (col) {
        setIssues((prev) => prev.map((i) =>
          i._id === p.id ? { ...i, status: col.id, status_id: p.status_id } : i
        ));
      }
    });
    return () => { offMoved(); offUpdated(); };
  }, [projectId, setIssues]);

  // Load swimlane preference for this board
  useEffect(() => {
    if (!projectId) return;
    api("/projects/" + projectId + "/boards")
      .then((res) => {
        const boards = res.items || res || [];
        if (!boards.length) return;
        const id = boards[0].id;
        setBoardId(id);
        return api("/boards/" + id + "/swimlanes");
      })
      .then((d) => { if (d) setSwimlane(d.swimlane_type || "none"); })
      .catch(() => {});
  }, [projectId]);

  async function changeSwimlane(v) {
    setSwimlane(v);
    if (!boardId) return;
    try { await api("/boards/" + boardId + "/swimlane-type", { method: "PUT", body: { swimlane_type: v } }); }
    catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function onDrop(colId) {
    if (!dragging) return;
    const issue = issues.find((i) => i.id === dragging);
    setDragging(null); setOverCol(null);
    if (!issue) return;
    const col = (columns || []).find((c) => c.id === colId);
    setIssues((prev) => prev.map((i) => i.id === dragging ? { ...i, status: colId } : i));
    if (col && issue._id) {
      try { await api("/issues/" + issue._id, { method: "PUT", body: { status_id: col._id } }); }
      catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
    }
  }

  const filtered = useMemo(() => (issues || []).filter((i) => {
    if (filter.assignee !== "all" && i.assignee !== filter.assignee) return false;
    if (filter.type !== "all" && i.type !== filter.type) return false;
    if (filter.search && !(i.title.toLowerCase().includes(filter.search.toLowerCase()) || i.id.toLowerCase().includes(filter.search.toLowerCase()))) return false;
    return true;
  }), [issues, filter]);

  const lanes = useMemo(() => {
    if (swimlane === "assignee") {
      const out = [];
      (people || []).forEach((u) => {
        const its = filtered.filter((i) => i.assignee === u.id);
        if (its.length) out.push({ key: u.id, header: { type: "user", user: u, count: its.length }, issues: its });
      });
      const un = filtered.filter((i) => !i.assignee);
      if (un.length) out.push({ key: "none", header: { type: "user", user: null, count: un.length }, issues: un });
      return out;
    }
    if (swimlane === "epic") {
      const epics = Array.from(new Set(filtered.map((i) => i.labels[0]).filter(Boolean)));
      const out = epics.map((ep, idx) => ({ key: ep, header: { type: "epic", title: ep, color: EPIC_PALETTE[idx % EPIC_PALETTE.length], count: filtered.filter((i) => i.labels[0] === ep).length }, issues: filtered.filter((i) => i.labels[0] === ep) }));
      const noEpic = filtered.filter((i) => !i.labels[0]);
      if (noEpic.length) out.push({ key: "_none", header: { type: "epic", title: "No epic", color: "var(--text-muted)", count: noEpic.length }, issues: noEpic });
      return out;
    }
    return [{ key: "all", header: null, issues: filtered }];
  }, [filtered, swimlane, people]);

  const proj = (projects || []).find((p) => p.id === projectId);

  return (
    <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <div className="page-head" style={{ paddingBottom: 12 }}>
        <div>
          <div className="crumbs">
            <a href="#/" onClick={(e) => { e.preventDefault(); nav("projects"); }}>Projects</a>
            <Icon name="chevronRight" size={11}/>
            <span>{proj ? proj.name : "…"}</span>
            <Icon name="chevronRight" size={11}/>
            <span>Board</span>
          </div>
          <h1>Board</h1>
          <div className="row gap-3" style={{ marginTop: 6 }}>
            <AvatarStack users={(people || []).slice(0, 6)} max={5} size="sm"/>
          </div>
        </div>
        <div className="row gap-2">
          <Button icon="filter">Filter</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>
            Create issue <span className="kbd" style={{ marginLeft: 4 }}>C</span>
          </Button>
        </div>
      </div>

      <div className="row gap-2" style={{ padding: "0 32px 12px", borderBottom: "1px solid var(--border)" }}>
        <div className="search" style={{ width: 240, padding: "4px 10px" }}>
          <Icon name="search" size={13}/>
          <input placeholder="Search board…" value={filter.search} onChange={(e) => setFilter({ ...filter, search: e.target.value })}/>
        </div>
        <Pill label="Assignee" icon="user" value={filter.assignee === "all" ? "All" : (people || []).find((p) => p.id === filter.assignee)?.name} options={[{ id: "all", label: "All" }, ...(people || []).map((p) => ({ id: p.id, label: p.name }))]} onChange={(v) => setFilter({ ...filter, assignee: v })}/>
        <Pill label="Type" icon="tag" value={filter.type === "all" ? "All" : filter.type} options={[{ id: "all", label: "All" }, ...Object.keys(TYPE_META).map((t) => ({ id: t, label: t }))]} onChange={(v) => setFilter({ ...filter, type: v })}/>
        <Pill label="Swimlanes" icon="layout" value={(SWIMLANE_OPTIONS.find((o) => o.id === swimlane) || SWIMLANE_OPTIONS[0]).label} options={SWIMLANE_OPTIONS} onChange={changeSwimlane}/>
        <div style={{ flex: 1 }}/>
        <Button variant="ghost" data-size="sm" icon="moreH" style={{ padding: "0 8px" }}/>
      </div>

      {issuesLoading ? (
        <div className="kanban">
          {[1, 2, 3, 4].map((col) => (
            <div key={col} className="kanban-col">
              <div className="kanban-col-head">
                <Skeleton w={100} h={14}/>
              </div>
              <div className="kanban-cards" style={{ gap: 8 }}>
                {Array.from({ length: 3 + (col % 3) }).map((_, i) => (
                  <div key={i} className="kanban-card" style={{ gap: 10 }}>
                    <div className="row gap-2"><Skeleton w={16} h={16} radius="4px"/><Skeleton w={60} h={12}/></div>
                    <Skeleton w="90%" h={14}/>
                    <Skeleton w="70%" h={12}/>
                    <div className="row gap-2" style={{ marginTop: 4, justifyContent: "space-between" }}>
                      <Skeleton w={48} h={18} radius="99px"/>
                      <Skeleton w={24} h={24} radius="50%"/>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      ) : swimlane === "none" ? (
        <KanbanGrid lane={lanes[0]} columns={columns} people={people} dragging={dragging} overCol={overCol} setOverCol={setOverCol} onDrop={onDrop} setDragging={setDragging} nav={nav}/>
      ) : (
        <div style={{ overflowY: "auto", flex: 1 }}>
          {lanes.map((lane) => (
            <div key={lane.key}>
              <div className="row gap-3" style={{ padding: "12px 24px 4px", position: "sticky", left: 0 }}>
                {lane.header.type === "user" ? (
                  lane.header.user
                    ? <><Avatar user={lane.header.user} size="sm"/><span className="bold text-sm">{lane.header.user.name}</span></>
                    : <><span className="avatar" data-size="sm" style={{ background: "var(--bg-muted)", color: "var(--text-secondary)" }}>?</span><span className="bold text-sm">Unassigned</span></>
                ) : (
                  <><span style={{ width: 4, height: 16, borderRadius: 2, background: lane.header.color }}/><span className="bold text-sm" style={{ textTransform: "capitalize" }}>{lane.header.title}</span></>
                )}
                <span className="ct" style={{ background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 99, padding: "1px 8px", fontSize: 11, color: "var(--text-secondary)" }}>{lane.header.count}</span>
              </div>
              <KanbanGrid lane={lane} columns={columns} people={people} dragging={dragging} overCol={overCol} setOverCol={setOverCol} onDrop={onDrop} setDragging={setDragging} nav={nav} compact/>
            </div>
          ))}
          {lanes.length === 0 && <div style={{ padding: 40 }}><Empty icon="kanban" title="No issues match" hint="Adjust the filters above."/></div>}
        </div>
      )}
    </div>
  );
}

function KanbanGrid({ lane, columns, people, dragging, overCol, setOverCol, onDrop, setDragging, nav, compact }) {
  const cols = columns || [];
  const byCol = {};
  cols.forEach((c) => { byCol[c.id] = []; });
  (lane ? lane.issues : []).forEach((i) => { if (byCol[i.status] !== undefined) byCol[i.status].push(i); });
  return (
    <div className="kanban" style={compact ? { paddingTop: 6, paddingBottom: 14, height: "auto" } : undefined}>
      {cols.map((c) => {
        const okey = (lane ? lane.key : "all") + "::" + c.id;
        return (
          <div key={c.id} className="kanban-col" style={compact ? { maxHeight: "none" } : undefined}>
            <div className="kanban-col-head">
              <span className="row gap-2">
                <span style={{ width: 8, height: 8, borderRadius: 2, background: colColor(c) }}/>
                {c.label}
                <span className="ct">{byCol[c.id] ? byCol[c.id].length : 0}</span>
              </span>
              <button className="icon-btn" style={{ width: 22, height: 22 }} aria-label="Add issue to column"><Icon name="plus" size={13}/></button>
            </div>
            <div
              className={"kanban-cards" + (overCol === okey && dragging ? " over" : "")}
              onDragOver={(e) => { e.preventDefault(); setOverCol(okey); }}
              onDragLeave={(e) => { if (e.currentTarget === e.target) setOverCol(null); }}
              onDrop={(e) => { e.preventDefault(); onDrop(c.id); }}
              style={compact ? { minHeight: 40 } : undefined}
            >
              {(byCol[c.id] || []).map((i) => (
                <KanbanCard key={i.id} issue={i} people={people} isDragging={dragging === i.id} onClick={() => nav("issue/" + i.id)} onDragStart={() => setDragging(i.id)} onDragEnd={() => { setDragging(null); setOverCol(null); }}/>
              ))}
              {(!byCol[c.id] || byCol[c.id].length === 0) && !compact && (
                <div style={{ padding: "20px 8px", textAlign: "center", color: "var(--text-muted)", fontSize: 12, border: "1px dashed var(--border)", borderRadius: 6 }}>Drop issues here</div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}

function KanbanCard({ issue, people, onClick, onDragStart, onDragEnd, isDragging }) {
  const assignee = issue.assigneeUser || (people || []).find((p) => p.id === issue.assignee);
  return (
    <div className="kanban-card" data-dragging={isDragging} draggable onDragStart={onDragStart} onDragEnd={onDragEnd} onClick={onClick}>
      <div className="row gap-2" style={{ marginBottom: 6 }}>
        <TypeIcon value={issue.type}/>
        <span className="mono text-xs muted">{issue.id}</span>
        {issue.labels && issue.labels.slice(0, 1).map((l) => (
          <span key={l} className="tag" style={{ marginLeft: "auto", textTransform: "lowercase" }}>{l}</span>
        ))}
      </div>
      <div className="title">{issue.title}</div>
      <div className="foot">
        <div className="row gap-2">
          <PriorityBadge value={issue.pri}/>
          {issue.points != null && (
            <span style={{ width: 18, height: 18, borderRadius: "50%", background: "var(--bg-muted)", display: "inline-grid", placeItems: "center", fontSize: 10.5, fontWeight: 600, color: "var(--text-secondary)" }}>{issue.points}</span>
          )}
        </div>
        <div className="row gap-2">
          {issue.comments > 0 && <span className="text-xs muted row gap-1"><Icon name="comment" size={12}/>{issue.comments}</span>}
          {issue.sub > 0 && <span className="text-xs muted row gap-1"><Icon name="checkbox" size={12}/>{issue.sub}</span>}
          {assignee && <Avatar user={assignee} size="sm"/>}
        </div>
      </div>
    </div>
  );
}

export function Pill({ label, icon, value, options, onChange }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);
  useEffect(() => {
    if (!open) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [open]);
  return (
    <span ref={ref} style={{ position: "relative" }}>
      <button className="btn btn-secondary" data-size="sm" onClick={() => options && setOpen((o) => !o)} style={{ borderStyle: "dashed", color: "var(--text-secondary)" }}>
        {icon && <Icon name={icon} size={12}/>}
        <span className="text-xs">{label}:</span>
        <span className="text-xs bold" style={{ color: "var(--text)" }}>{value}</span>
      </button>
      {open && options && (
        <div style={{ position: "absolute", top: "calc(100% + 4px)", left: 0, minWidth: 180, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, maxHeight: 280, overflowY: "auto" }}>
          {options.map((o) => (
            <button key={o.id} onClick={() => { onChange(o.id); setOpen(false); }} className="nav-item" style={{ color: "var(--text)", fontSize: 13 }}>
              {o.label}
            </button>
          ))}
        </div>
      )}
    </span>
  );
}

export function PillSelect({ value, onChange, options }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);
  useEffect(() => {
    if (!open) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [open]);
  const opt = (options || []).find((o) => o.id === value);
  return (
    <span ref={ref} style={{ position: "relative", display: "block" }}>
      <button onClick={() => setOpen((o) => !o)} className="btn btn-secondary" style={{ width: "100%", justifyContent: "space-between" }}>
        <span className="row gap-2">{opt && opt.color && <span style={{ width: 10, height: 10, borderRadius: 3, background: opt.color }}/>}{opt ? opt.label : "—"}</span>
        <Icon name="chevronDown" size={12} color="var(--text-muted)"/>
      </button>
      {open && (
        <div style={{ position: "absolute", top: "calc(100% + 4px)", left: 0, right: 0, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, maxHeight: 240, overflowY: "auto" }}>
          {(options || []).map((o) => (
            <button key={o.id} onClick={() => { onChange(o.id); setOpen(false); }} className="nav-item" style={{ color: "var(--text)", fontSize: 13 }}>
              {o.color && <span style={{ width: 10, height: 10, borderRadius: 3, background: o.color }}/>}
              {o.label}
              {o.id === value && <Icon name="check" size={13} color="var(--indigo-600)" style={{ marginLeft: "auto" }}/>}
            </button>
          ))}
        </div>
      )}
    </span>
  );
}

function Field({ label, children }) {
  return (
    <div>
      <label className="label" style={{ fontSize: 11.5, color: "var(--text-muted)", textTransform: "uppercase", letterSpacing: ".04em", fontWeight: 600 }}>{label}</label>
      {children}
    </div>
  );
}

export function CreateIssueModal({ open, onClose, onCreate }) {
  const { people, columns } = useApp();
  const [form, setForm] = useState({
    title: "", desc: "", type: "Task", pri: "Medium", status: "",
    assignee: "", points: 3, due: "", labels: ""
  });
  useEffect(() => {
    if (open) {
      const defStatus = columns && columns[0] ? columns[0].id : "";
      setForm({ title: "", desc: "", type: "Task", pri: "Medium", status: defStatus, assignee: "", points: 3, due: "", labels: "" });
    }
  }, [open, columns]);

  return (
    <Modal open={open} onClose={onClose} title="Create issue" size="lg"
      footer={
        <>
          <span className="text-xs muted" style={{ marginRight: "auto" }} aria-hidden="true">
            <span className="kbd">⌘</span> + <span className="kbd">↵</span> to submit
          </span>
          <Button onClick={onClose}>Cancel</Button>
          <Button variant="primary" disabled={!form.title} onClick={() => { onCreate(form); onClose(); }}>Create issue</Button>
        </>
      }
    >
      <div style={{ display: "grid", gridTemplateColumns: "1fr 200px", gap: 24 }}>
        <div>
          <div style={{ marginBottom: 14 }}>
            <label className="label">Title</label>
            <input className="input" autoFocus placeholder="What needs to happen?" value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })}/>
          </div>
          <div style={{ marginBottom: 14 }}>
            <label className="label">Description</label>
            <div style={{ border: "1px solid var(--border)", borderRadius: 6, overflow: "hidden" }}>
              <div className="row gap-1" style={{ padding: 6, borderBottom: "1px solid var(--border)", background: "var(--bg-subtle)" }}>
                {[["bold","B"],["italic","I"],["heading","H"],["list","L"],["code","Code"],["link","Link"],["picture","Img"]].map(([ic, name]) => (
                  <button key={name} className="icon-btn" style={{ width: 26, height: 26 }} title={name}><Icon name={ic} size={13}/></button>
                ))}
              </div>
              <textarea className="textarea" style={{ border: 0, borderRadius: 0, minHeight: 140 }} placeholder="Steps to reproduce, acceptance criteria, links…" value={form.desc} onChange={(e) => setForm({ ...form, desc: e.target.value })}/>
            </div>
          </div>
        </div>
        <div className="stack gap-3">
          <Field label="Type">
            <PillSelect value={form.type} onChange={(v) => setForm({ ...form, type: v })} options={Object.keys(TYPE_META).map((t) => ({ id: t, label: t, color: TYPE_META[t].color, icon: TYPE_META[t].icon }))}/>
          </Field>
          <Field label="Status">
            <PillSelect value={form.status} onChange={(v) => setForm({ ...form, status: v })} options={(columns || []).map((c) => ({ id: c.id, label: c.label || c.id }))}/>
          </Field>
          <Field label="Priority">
            <PillSelect value={form.pri} onChange={(v) => setForm({ ...form, pri: v })} options={Object.keys(PRIORITY_META).map((p) => ({ id: p, label: p, color: PRIORITY_META[p].color, icon: PRIORITY_META[p].icon }))}/>
          </Field>
          <Field label="Assignee">
            <PillSelect value={form.assignee} onChange={(v) => setForm({ ...form, assignee: v })} options={[{ id: "", label: "Unassigned" }, ...(people || []).map((u) => ({ id: u.id, label: u.name }))]}/>
          </Field>
          <Field label="Story points">
            <input className="input" type="number" min="0" max="13" value={form.points} onChange={(e) => setForm({ ...form, points: Number(e.target.value) })}/>
          </Field>
          <Field label="Due date">
            <input className="input" type="text" placeholder="MMM DD" value={form.due} onChange={(e) => setForm({ ...form, due: e.target.value })}/>
          </Field>
          <Field label="Labels">
            <input className="input" placeholder="terraform, oncall…" value={form.labels} onChange={(e) => setForm({ ...form, labels: e.target.value })}/>
          </Field>
        </div>
      </div>
    </Modal>
  );
}
