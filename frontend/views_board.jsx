// views_board.jsx — kanban board with drag-and-drop + swimlanes (Feature 15)

const SWIMLANE_OPTIONS = [{ id: "none", label: "No swimlanes" }, { id: "assignee", label: "By Assignee" }, { id: "epic", label: "By Epic" }];
const EPIC_PALETTE = ["#6366F1", "#10B981", "#F59E0B", "#EF4444", "#8B5CF6", "#06B6D4", "#EC4899", "#F97316"];
const BOARD_ID = FORGE_DATA.ACTIVE_PROJECT_ID;

function BoardView({ nav, issues, setIssues }) {
  const [filter, setFilter] = React.useState({ assignee: "all", type: "all", search: "" });
  const [dragging, setDragging] = React.useState(null); // id
  const [overCol, setOverCol] = React.useState(null); // "laneKey::colId"
  const [swimlane, setSwimlane] = React.useState("none");
  const toast = useToast();

  React.useEffect(() => {
    api("/boards/" + BOARD_ID + "/swimlanes").then((d) => setSwimlane(d.swimlane_type || "none")).catch(() => {});
  }, []);

  async function changeSwimlane(v) {
    setSwimlane(v);
    try { await api("/boards/" + BOARD_ID + "/swimlane-type", { method: "PUT", body: { swimlane_type: v } }); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  const filtered = React.useMemo(() => issues.filter((i) => {
    if (filter.assignee !== "all" && i.assignee !== filter.assignee) return false;
    if (filter.type !== "all" && i.type !== filter.type) return false;
    if (filter.search && !(i.title.toLowerCase().includes(filter.search.toLowerCase()) || i.id.toLowerCase().includes(filter.search.toLowerCase()))) return false;
    return true;
  }), [issues, filter]);

  function onDrop(colId) {
    if (!dragging) return;
    setIssues((prev) => prev.map((i) => i.id === dragging ? { ...i, status: colId } : i));
    setDragging(null); setOverCol(null);
  }

  // Build lanes
  const lanes = React.useMemo(() => {
    if (swimlane === "assignee") {
      const out = [];
      FORGE_DATA.PEOPLE.forEach((u) => {
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
  }, [filtered, swimlane]);

  return (
    <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <div className="page-head" style={{ paddingBottom: 12 }}>
        <div>
          <div className="crumbs"><a href="#/" onClick={(e) => { e.preventDefault(); nav("projects"); }}>Projects</a> <Icon name="chevronRight" size={11}/> <a href="#/">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Board</span></div>
          <h1>Sprint 24 — Edge & Reliability</h1>
          <div className="row gap-3" style={{ marginTop: 6 }}>
            <span className="text-xs muted row gap-1"><Icon name="calendar" size={12}/> Dec 2 → Dec 15</span>
            <span className="text-xs muted">·</span>
            <span className="text-xs muted">5 days left</span>
            <AvatarStack users={FORGE_DATA.PEOPLE.slice(0, 6)} max={5} size="sm"/>
          </div>
        </div>
        <div className="row gap-2">
          <Button icon="filter">Filter</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>Create issue <span className="kbd" style={{ marginLeft: 4 }}>C</span></Button>
        </div>
      </div>

      {/* Filter bar */}
      <div className="row gap-2" style={{ padding: "0 32px 12px", borderBottom: "1px solid var(--border)" }}>
        <div className="search" style={{ width: 240, padding: "4px 10px" }}>
          <Icon name="search" size={13}/>
          <input placeholder="Search board…" value={filter.search} onChange={(e) => setFilter({ ...filter, search: e.target.value })}/>
        </div>
        <Pill label="Assignee" icon="user" value={filter.assignee === "all" ? "All" : FORGE_DATA.PEOPLE.find((p) => p.id === filter.assignee)?.name} options={[{ id: "all", label: "All" }, ...FORGE_DATA.PEOPLE.map((p) => ({ id: p.id, label: p.name }))]} onChange={(v) => setFilter({ ...filter, assignee: v })}/>
        <Pill label="Type" icon="tag" value={filter.type === "all" ? "All" : filter.type} options={[{ id: "all", label: "All" }, ...Object.keys(FORGE_DATA.TYPE_META).map((t) => ({ id: t, label: t }))]} onChange={(v) => setFilter({ ...filter, type: v })}/>
        <Pill label="Swimlanes" icon="layout" value={SWIMLANE_OPTIONS.find((o) => o.id === swimlane).label} options={SWIMLANE_OPTIONS} onChange={changeSwimlane}/>
        <div style={{ flex: 1 }}/>
        <Button variant="ghost" data-size="sm" icon="moreH" style={{ padding: "0 8px" }}/>
      </div>

      {swimlane === "none" ? (
        <KanbanGrid lane={lanes[0]} dragging={dragging} overCol={overCol} setOverCol={setOverCol} onDrop={onDrop} setDragging={setDragging} nav={nav}/>
      ) : (
        <div style={{ overflowY: "auto", flex: 1 }}>
          {lanes.map((lane) => (
            <div key={lane.key}>
              <div className="row gap-3" style={{ padding: "12px 24px 4px", position: "sticky", left: 0 }}>
                {lane.header.type === "user" ? (
                  lane.header.user ? <><Avatar user={lane.header.user} size="sm"/><span className="bold text-sm">{lane.header.user.name}</span></>
                    : <><span className="avatar" data-size="sm" style={{ background: "var(--bg-muted)", color: "var(--text-secondary)" }}>?</span><span className="bold text-sm">Unassigned</span></>
                ) : (
                  <><span style={{ width: 4, height: 16, borderRadius: 2, background: lane.header.color }}/><span className="bold text-sm" style={{ textTransform: "capitalize" }}>{lane.header.title}</span></>
                )}
                <span className="ct" style={{ background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 99, padding: "1px 8px", fontSize: 11, color: "var(--text-secondary)" }}>{lane.header.count}</span>
              </div>
              <KanbanGrid lane={lane} dragging={dragging} overCol={overCol} setOverCol={setOverCol} onDrop={onDrop} setDragging={setDragging} nav={nav} compact/>
            </div>
          ))}
          {lanes.length === 0 && <div style={{ padding: 40 }}><Empty icon="kanban" title="No issues match" hint="Adjust the filters above."/></div>}
        </div>
      )}
    </div>
  );
}

function KanbanGrid({ lane, dragging, overCol, setOverCol, onDrop, setDragging, nav, compact }) {
  const byCol = {};
  FORGE_DATA.COLUMNS.forEach((c) => byCol[c.id] = []);
  lane.issues.forEach((i) => byCol[i.status] && byCol[i.status].push(i));
  return (
    <div className="kanban" style={compact ? { paddingTop: 6, paddingBottom: 14, height: "auto" } : undefined}>
      {FORGE_DATA.COLUMNS.map((c) => {
        const okey = lane.key + "::" + c.id;
        return (
          <div key={c.id} className="kanban-col" style={compact ? { maxHeight: "none" } : undefined}>
            <div className="kanban-col-head">
              <span className="row gap-2">
                <span style={{ width: 8, height: 8, borderRadius: 2, background: { Backlog: "#94A3B8", Todo: "#64748B", "In Progress": "#3B82F6", "In Review": "#8B5CF6", Done: "#10B981" }[c.id] }}/>
                {c.label}
                <span className="ct">{byCol[c.id].length}</span>
              </span>
              <button className="icon-btn" style={{ width: 22, height: 22 }} aria-label="Add to column"><Icon name="plus" size={13}/></button>
            </div>
            <div
              className={"kanban-cards" + (overCol === okey && dragging ? " over" : "")}
              onDragOver={(e) => { e.preventDefault(); setOverCol(okey); }}
              onDragLeave={(e) => { if (e.currentTarget === e.target) setOverCol(null); }}
              onDrop={(e) => { e.preventDefault(); onDrop(c.id); }}
              style={compact ? { minHeight: 40 } : undefined}
            >
              {byCol[c.id].map((i) => (
                <KanbanCard key={i.id} issue={i} isDragging={dragging === i.id} onClick={() => nav("issue/" + i.id)} onDragStart={() => setDragging(i.id)} onDragEnd={() => { setDragging(null); setOverCol(null); }}/>
              ))}
              {byCol[c.id].length === 0 && !compact && (
                <div style={{ padding: "20px 8px", textAlign: "center", color: "var(--text-muted)", fontSize: 12, border: "1px dashed var(--border)", borderRadius: 6 }}>Drop issues here</div>
              )}
            </div>
          </div>
        );
      })}
    </div>
  );
}

function KanbanCard({ issue, onClick, onDragStart, onDragEnd, isDragging }) {
  const assignee = FORGE_DATA.PEOPLE.find((p) => p.id === issue.assignee);
  return (
    <div
      className="kanban-card"
      data-dragging={isDragging}
      draggable
      onDragStart={onDragStart}
      onDragEnd={onDragEnd}
      onClick={onClick}
    >
      <div className="row gap-2" style={{ marginBottom: 6 }}>
        <TypeIcon value={issue.type}/>
        <span className="mono text-xs muted">{issue.id}</span>
        {issue.labels.slice(0, 1).map((l) => (
          <span key={l} className="tag" style={{ marginLeft: "auto", textTransform: "lowercase" }}>{l}</span>
        ))}
      </div>
      <div className="title">{issue.title}</div>
      <div className="foot">
        <div className="row gap-2">
          <PriorityBadge value={issue.pri}/>
          {issue.points != null && (
            <span style={{
              width: 18, height: 18, borderRadius: "50%",
              background: "var(--bg-muted)",
              display: "inline-grid", placeItems: "center",
              fontSize: 10.5, fontWeight: 600, color: "var(--text-secondary)"
            }}>{issue.points}</span>
          )}
        </div>
        <div className="row gap-2">
          {issue.comments > 0 && (
            <span className="text-xs muted row gap-1"><Icon name="comment" size={12}/>{issue.comments}</span>
          )}
          {issue.sub > 0 && (
            <span className="text-xs muted row gap-1"><Icon name="checkbox" size={12}/>{issue.sub}</span>
          )}
          {assignee && <Avatar user={assignee} size="sm"/>}
        </div>
      </div>
    </div>
  );
}

function Pill({ label, icon, value, options, onChange }) {
  const [open, setOpen] = React.useState(false);
  const ref = React.useRef(null);
  React.useEffect(() => {
    if (!open) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [open]);
  return (
    <span ref={ref} style={{ position: "relative" }}>
      <button
        className="btn btn-secondary"
        data-size="sm"
        onClick={() => options && setOpen((o) => !o)}
        style={{ borderStyle: "dashed", color: "var(--text-secondary)" }}
      >
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

// ─── Create issue modal ─────────────────────────────────
function CreateIssueModal({ open, onClose, onCreate }) {
  const [form, setForm] = React.useState({
    title: "", desc: "", type: "Task", pri: "Medium", status: "Todo",
    assignee: "u1", points: 3, due: "", labels: ""
  });
  React.useEffect(() => { if (open) setForm({ title: "", desc: "", type: "Task", pri: "Medium", status: "Todo", assignee: "u1", points: 3, due: "", labels: "" }); }, [open]);

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
                {[["bold", "B"], ["italic", "I"], ["heading", "H"], ["list", "L"], ["code", "Code"], ["link", "Link"], ["picture", "Img"]].map(([ic, name]) => (
                  <button key={name} className="icon-btn" style={{ width: 26, height: 26 }} title={name}><Icon name={ic} size={13}/></button>
                ))}
              </div>
              <textarea className="textarea" style={{ border: 0, borderRadius: 0, minHeight: 140 }} placeholder="Steps to reproduce, acceptance criteria, links…" value={form.desc} onChange={(e) => setForm({ ...form, desc: e.target.value })}/>
            </div>
          </div>
          <div className="row gap-2">
            <Button icon="paperclip">Attach</Button>
            <Button icon="link">Link issue</Button>
            <Button icon="at">Mention</Button>
          </div>
        </div>
        <div className="stack gap-3">
          <Field label="Type">
            <PillSelect value={form.type} onChange={(v) => setForm({ ...form, type: v })} options={Object.keys(FORGE_DATA.TYPE_META).map((t) => ({ id: t, label: t, color: FORGE_DATA.TYPE_META[t].color, icon: FORGE_DATA.TYPE_META[t].icon }))}/>
          </Field>
          <Field label="Status">
            <PillSelect value={form.status} onChange={(v) => setForm({ ...form, status: v })} options={FORGE_DATA.COLUMNS.map((c) => ({ id: c.id, label: c.id }))}/>
          </Field>
          <Field label="Priority">
            <PillSelect value={form.pri} onChange={(v) => setForm({ ...form, pri: v })} options={Object.keys(FORGE_DATA.PRIORITY_META).map((p) => ({ id: p, label: p, color: FORGE_DATA.PRIORITY_META[p].color, icon: FORGE_DATA.PRIORITY_META[p].icon }))}/>
          </Field>
          <Field label="Assignee">
            <PillSelect value={form.assignee} onChange={(v) => setForm({ ...form, assignee: v })} options={FORGE_DATA.PEOPLE.map((u) => ({ id: u.id, label: u.name }))}/>
          </Field>
          <Field label="Sprint">
            <PillSelect value="s24" onChange={() => {}} options={[{ id: "s24", label: "Sprint 24" }, { id: "bl", label: "Backlog" }]}/>
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

function Field({ label, children }) {
  return (
    <div>
      <label className="label" style={{ fontSize: 11.5, color: "var(--text-muted)", textTransform: "uppercase", letterSpacing: ".04em", fontWeight: 600 }}>{label}</label>
      {children}
    </div>
  );
}

function PillSelect({ value, onChange, options }) {
  const [open, setOpen] = React.useState(false);
  const ref = React.useRef(null);
  React.useEffect(() => {
    if (!open) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [open]);
  const opt = options.find((o) => o.id === value);
  return (
    <span ref={ref} style={{ position: "relative", display: "block" }}>
      <button onClick={() => setOpen((o) => !o)} className="btn btn-secondary" style={{ width: "100%", justifyContent: "space-between" }}>
        <span className="row gap-2">{opt?.color && <span style={{ width: 10, height: 10, borderRadius: 3, background: opt.color }}/>}{opt?.label}</span>
        <Icon name="chevronDown" size={12} color="var(--text-muted)"/>
      </button>
      {open && (
        <div style={{ position: "absolute", top: "calc(100% + 4px)", left: 0, right: 0, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, maxHeight: 240, overflowY: "auto" }}>
          {options.map((o) => (
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

Object.assign(window, { BoardView, CreateIssueModal });
