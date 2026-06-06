// views_issue.jsx — Issue detail page + backlog/list views

function IssueView({ nav, issueId, issues, setIssues }) {
  const issue = issues.find((i) => i.id === issueId) || issues[8]; // fallback
  const [tab, setTab] = React.useState("comments");
  const [editingDesc, setEditingDesc] = React.useState(false);
  const [editingTitle, setEditingTitle] = React.useState(false);
  const [title, setTitle] = React.useState(issue.title);
  const [desc, setDesc] = React.useState("");
  const [comment, setComment] = React.useState("");
  const [comments, setComments] = React.useState([
    { id: 1, who: "u5", time: "2h ago", text: "Reproduced on staging. Looks like the chunk-encoding path introduced in v3.2 is leaking buffers past the 256MB GC root. Rolling staging back to v3.1.2 confirmed.", reactions: { "👀": 3, "🔥": 1 } },
    { id: 2, who: "u1", time: "1h ago", text: "Nice find @priya. Can you pin which commit between v3.1.2 and v3.2 introduced it? I want to be sure before we file upstream.", reactions: {} },
    { id: 3, who: "u3", time: "32m ago", text: "Bisected to `loki/storage/chunks: buffer pool reuse` (sha c0b9d44). Will open upstream issue today.", reactions: { "🎯": 2 } },
  ]);
  const assignee = FORGE_DATA.PEOPLE.find((p) => p.id === issue.assignee);
  const reporter = FORGE_DATA.PEOPLE.find((p) => p.id === issue.reporter);

  function postComment() {
    if (!comment.trim()) return;
    setComments([...comments, { id: Date.now(), who: "u1", time: "just now", text: comment, reactions: {} }]);
    setComment("");
  }

  const subtasks = [
    { id: "sub-1", title: "Bisect commit range", done: true },
    { id: "sub-2", title: "Pin upstream issue", done: true },
    { id: "sub-3", title: "Patch chunk-encoder buffer pool", done: false },
    { id: "sub-4", title: "Add ingester memory regression test", done: false },
  ];

  return (
    <div className="detail">
      <div className="detail-main">
        <div className="crumbs row gap-2" style={{ marginBottom: 12 }}>
          <a href="#/projects" onClick={(e) => { e.preventDefault(); nav("projects"); }}>Projects</a>
          <Icon name="chevronRight" size={11}/>
          <a href="#/board" onClick={(e) => { e.preventDefault(); nav("board"); }}>Core Infrastructure</a>
          <Icon name="chevronRight" size={11}/>
          <a href="#/board" onClick={(e) => { e.preventDefault(); nav("board"); }}>Board</a>
          <Icon name="chevronRight" size={11}/>
          <span className="mono" style={{ color: "var(--text-secondary)" }}>{issue.id}</span>
          <div className="row gap-1" style={{ marginLeft: "auto" }}>
            <Button variant="ghost" data-size="sm" icon="copy" title="Copy link"/>
            <Button variant="ghost" data-size="sm" icon="star" title="Watch"/>
            <Button variant="ghost" data-size="sm" icon="moreH"/>
          </div>
        </div>

        <div className="row gap-2" style={{ marginBottom: 12 }}>
          <TypeIcon value={issue.type}/>
          <span className="mono text-sm muted">{issue.id}</span>
          <Badge tone="danger">CVE risk</Badge>
          <Badge tone="muted">Sprint 24</Badge>
        </div>

        {editingTitle ? (
          <input className="input" autoFocus value={title} onChange={(e) => setTitle(e.target.value)}
            onBlur={() => setEditingTitle(false)}
            onKeyDown={(e) => { if (e.key === "Enter") setEditingTitle(false); }}
            style={{ fontSize: 24, fontWeight: 600, padding: "6px 10px", marginBottom: 16 }}
          />
        ) : (
          <h1 onClick={() => setEditingTitle(true)} style={{ fontSize: 24, fontWeight: 600, letterSpacing: "-.015em", margin: "0 0 16px", cursor: "text", padding: "2px 0" }}>
            {title}
          </h1>
        )}

        {/* Action buttons */}
        <div className="row gap-2" style={{ marginBottom: 24 }}>
          <Button variant="secondary" icon="plus">Add subtask</Button>
          <Button variant="secondary" icon="link">Link issue</Button>
          <Button variant="secondary" icon="paperclip">Attach</Button>
        </div>

        <h3 style={{ fontSize: 13, fontWeight: 600, letterSpacing: ".02em", color: "var(--text-secondary)", textTransform: "uppercase", margin: "0 0 8px" }}>Description</h3>
        {editingDesc ? (
          <div style={{ marginBottom: 16 }}>
            <textarea className="textarea" value={desc} onChange={(e) => setDesc(e.target.value)} autoFocus style={{ minHeight: 120 }}/>
            <div className="row gap-2" style={{ marginTop: 8 }}>
              <Button variant="primary" onClick={() => setEditingDesc(false)}>Save</Button>
              <Button onClick={() => setEditingDesc(false)}>Cancel</Button>
            </div>
          </div>
        ) : (
          <div onClick={() => setEditingDesc(true)} style={{ background: "var(--bg-subtle)", border: "1px solid var(--border)", borderRadius: 8, padding: "14px 16px", marginBottom: 24, cursor: "text" }}>
            <p style={{ margin: "0 0 10px" }}>The Loki ingester pods in <code style={{ background: "var(--bg)", padding: "1px 5px", borderRadius: 4, fontSize: 12.5 }}>observability/loki-ingester</code> are getting OOMKilled every 30–45 minutes under steady-state log volume.</p>
            <p style={{ margin: "0 0 10px" }}><strong>Impact:</strong> log ingestion lag spikes to 8+ minutes during restarts. Affects oncall visibility during peak.</p>
            <p style={{ margin: "0 0 4px", fontWeight: 600 }}>Acceptance criteria:</p>
            <ul style={{ margin: "0 0 10px 18px", paddingLeft: 0, lineHeight: 1.7 }}>
              <li>Ingester memory stays under 4GB p99 over 6h soak</li>
              <li>No ingester restarts attributable to memory pressure for 24h</li>
              <li>Regression test added to release pipeline</li>
            </ul>
            <p style={{ margin: 0, color: "var(--text-muted)", fontSize: 12 }}>Click to edit description</p>
          </div>
        )}

        {/* Subtasks */}
        <h3 style={{ fontSize: 13, fontWeight: 600, letterSpacing: ".02em", color: "var(--text-secondary)", textTransform: "uppercase", margin: "0 0 8px" }}>
          Subtasks <span style={{ color: "var(--text-muted)", fontWeight: 400, textTransform: "none", letterSpacing: 0, marginLeft: 4 }}>{subtasks.filter((s) => s.done).length} / {subtasks.length}</span>
        </h3>
        <div style={{ border: "1px solid var(--border)", borderRadius: 8, marginBottom: 24 }}>
          {subtasks.map((s, i) => (
            <div key={s.id} className="row gap-3" style={{ padding: "10px 14px", borderBottom: i < subtasks.length - 1 ? "1px solid var(--border)" : 0 }}>
              <span style={{
                width: 16, height: 16, borderRadius: 4,
                border: "1.5px solid " + (s.done ? "var(--indigo-600)" : "var(--border-strong)"),
                background: s.done ? "var(--indigo-600)" : "transparent",
                display: "grid", placeItems: "center", color: "#fff"
              }}>{s.done && <Icon name="check" size={11} strokeWidth={3}/>}</span>
              <span className="mono text-xs muted">{s.id.toUpperCase()}</span>
              <span style={{ flex: 1, textDecoration: s.done ? "line-through" : "none", color: s.done ? "var(--text-muted)" : "var(--text)" }}>{s.title}</span>
              <Avatar user={FORGE_DATA.PEOPLE[i % 6]} size="sm"/>
            </div>
          ))}
        </div>

        {/* Tabs */}
        <div style={{ borderBottom: "1px solid var(--border)", marginBottom: 16, marginLeft: -32, marginRight: -32, paddingLeft: 32, paddingRight: 32 }}>
          <div className="row gap-1">
            {[
              ["comments", "Comments", comments.length],
              ["activity", "Activity", 14],
              ["history", "History", 6],
            ].map(([id, label, ct]) => (
              <button key={id} className="tab" aria-selected={tab === id} onClick={() => setTab(id)}>
                {label} <span className="muted" style={{ marginLeft: 4, fontSize: 11.5 }}>{ct}</span>
              </button>
            ))}
          </div>
        </div>

        {tab === "comments" && (
          <div>
            {comments.map((c) => {
              const u = FORGE_DATA.PEOPLE.find((p) => p.id === c.who);
              return (
                <div key={c.id} className="row gap-3" style={{ marginBottom: 16, alignItems: "flex-start" }}>
                  <Avatar user={u} size="md"/>
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div className="row gap-2" style={{ marginBottom: 4 }}>
                      <span className="bold text-sm">{u.name}</span>
                      <span className="text-xs muted">{c.time}</span>
                    </div>
                    <div className="text-sm" style={{ lineHeight: 1.6 }}>{c.text.split(/(@\w+)/).map((part, i) => part.startsWith("@") ? <span key={i} style={{ background: "var(--indigo-50)", color: "var(--indigo-700)", padding: "0 4px", borderRadius: 3, fontWeight: 500 }}>{part}</span> : part)}</div>
                    {Object.keys(c.reactions).length > 0 && (
                      <div className="row gap-1" style={{ marginTop: 6 }}>
                        {Object.entries(c.reactions).map(([emo, ct]) => (
                          <span key={emo} style={{ fontSize: 12, padding: "1px 7px", border: "1px solid var(--border)", borderRadius: 12, background: "var(--bg-subtle)" }}>
                            {emo} <span className="muted">{ct}</span>
                          </span>
                        ))}
                        <button className="icon-btn" style={{ width: 22, height: 22 }} title="Add reaction"><span style={{ fontSize: 12 }}>＋</span></button>
                      </div>
                    )}
                  </div>
                </div>
              );
            })}

            {/* Comment box with @mention */}
            <CommentBox value={comment} onChange={setComment} onSubmit={postComment}/>
          </div>
        )}

        {tab === "activity" && (
          <div style={{ marginBottom: 32 }}>
            {[
              { who: "u1", text: "created this issue", time: "5d ago" },
              { who: "u1", text: "added label `loki`", time: "5d ago" },
              { who: "u2", text: "moved to Sprint 24", time: "4d ago" },
              { who: "u5", text: "changed status from Backlog → Todo", time: "2d ago" },
              { who: "u3", text: "changed priority Medium → Critical", time: "yesterday" },
              { who: "u5", text: "linked PR #4471 from forge/loki", time: "1h ago" },
            ].map((a, i) => {
              const u = FORGE_DATA.PEOPLE.find((p) => p.id === a.who);
              return (
                <div key={i} className="row gap-3" style={{ padding: "8px 0", borderBottom: i < 5 ? "1px solid var(--border)" : 0 }}>
                  <Avatar user={u} size="sm"/>
                  <span className="text-sm secondary"><span className="bold" style={{ color: "var(--text)" }}>{u.name}</span> {a.text}</span>
                  <span className="text-xs muted" style={{ marginLeft: "auto" }}>{a.time}</span>
                </div>
              );
            })}
          </div>
        )}

        {tab === "history" && <HistoryTab issueId={issue.id}/>}
      </div>

      {/* RIGHT SIDEBAR */}
      <aside className="detail-aside">
        <div className="row gap-2" style={{ marginBottom: 16 }}>
          <PillSelect value={issue.status} onChange={(v) => setIssues((p) => p.map((i) => i.id === issue.id ? { ...i, status: v } : i))} options={FORGE_DATA.COLUMNS.map((c) => ({ id: c.id, label: c.id }))}/>
          <Button variant="secondary" icon="moreH" data-size="sm" style={{ padding: "0 8px" }}/>
        </div>

        <AssigneesPanel issueId={issue.id}/>
        <WatchVotePanel issueId={issue.id}/>
        <TimeTrackingPanel issueId={issue.id}/>
        <IssueLinksPanel issueId={issue.id}/>

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
            <dd>{issue.due ? <span className="row gap-1"><Icon name="calendar" size={12} color="var(--text-muted)"/>{issue.due}, 2024</span> : "—"}</dd>
            <dt>Labels</dt>
            <dd className="row gap-1" style={{ flexWrap: "wrap" }}>{issue.labels.map((l) => <span key={l} className="tag">{l}</span>)}</dd>
            <dt>Created</dt>
            <dd className="text-xs muted">Nov 27, 09:12</dd>
            <dt>Updated</dt>
            <dd className="text-xs muted">12m ago</dd>
          </dl>
        </div>

        {/* Notifications via Telegram */}
        <div className="tg-card" style={{ padding: 12, marginBottom: 18 }}>
          <div className="row gap-2" style={{ marginBottom: 6 }}>
            <Icon name="telegram" size={14} color="#2AABEE"/>
            <span className="bold text-sm">Telegram notifications</span>
          </div>
          <div className="text-xs secondary" style={{ marginBottom: 8 }}>
            {assignee?.tg ? <>Assignee <span className="bold">{assignee.tg}</span> will receive updates.</> : "Assignee hasn't connected Telegram."}
          </div>
          <Button variant="ghost" data-size="sm" icon="bell" style={{ padding: 0, color: "var(--tg)" }}>Manage</Button>
        </div>

        <div>
          <h4 style={{ fontSize: 11, fontWeight: 600, letterSpacing: ".06em", color: "var(--text-muted)", textTransform: "uppercase", margin: "0 0 8px" }}>Linked work</h4>
          <div style={{ display: "grid", gap: 6 }}>
            <LinkedItem icon="branch" label="forge/loki#4471" sub="PR · Open · +84 −12" tone="info"/>
            <LinkedItem icon="branch" label="forge/loki#4452" sub="PR · Merged" tone="purple"/>
            <LinkedItem icon="bug" label="INFRA-198" sub="Blocks · Done" tone="muted"/>
          </div>
        </div>
      </aside>
    </div>
  );
}

function LinkedItem({ icon, label, sub, tone }) {
  return (
    <div className="row gap-3" style={{ padding: "8px 10px", background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
      <Icon name={icon} size={14}/>
      <div className="stack" style={{ lineHeight: 1.2 }}>
        <span className="text-sm medium">{label}</span>
        <span className="text-xs muted">{sub}</span>
      </div>
      <Icon name="externalLink" size={12} color="var(--text-muted)" style={{ marginLeft: "auto" }}/>
    </div>
  );
}

function CommentBox({ value, onChange, onSubmit }) {
  const [mentioning, setMentioning] = React.useState(false);
  const [mentionQ, setMentionQ] = React.useState("");

  React.useEffect(() => {
    const at = value.lastIndexOf("@");
    if (at >= 0 && (at === 0 || value[at - 1] === " ")) {
      const q = value.slice(at + 1);
      if (!q.includes(" ") && q.length < 20) {
        setMentioning(true);
        setMentionQ(q);
        return;
      }
    }
    setMentioning(false);
  }, [value]);

  const matches = mentioning
    ? FORGE_DATA.PEOPLE.filter((p) => p.name.toLowerCase().includes(mentionQ.toLowerCase())).slice(0, 5)
    : [];

  function pick(u) {
    const at = value.lastIndexOf("@");
    const newVal = value.slice(0, at) + "@" + u.name.split(" ")[0].toLowerCase() + " ";
    onChange(newVal);
    setMentioning(false);
  }

  return (
    <div className="row gap-3" style={{ alignItems: "flex-start", marginBottom: 32, position: "relative" }}>
      <Avatar user={FORGE_DATA.ME} size="md"/>
      <div style={{ flex: 1, border: "1px solid var(--border)", borderRadius: 8, overflow: "hidden", background: "var(--bg)" }}>
        <textarea
          className="textarea"
          style={{ border: 0, borderRadius: 0, minHeight: 80 }}
          placeholder="Add a comment… Use @ to mention"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={(e) => { if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) onSubmit(); }}
        />
        <div className="row" style={{ padding: 6, borderTop: "1px solid var(--border)", background: "var(--bg-subtle)" }}>
          <div className="row gap-1">
            {[["bold", "Bold"], ["italic", "Italic"], ["code", "Code"], ["link", "Link"], ["picture", "Image"], ["at", "Mention"]].map(([ic, name]) => (
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

// ─── Backlog (table view) — Feature 14: saved filters ───
const EMPTY_FILTER = { status: "all", assignee: "all", priority: "all", type: "all", label: "all", sprint: "all" };
// saved-filter filter_json uses the same keys as our internal filter state.
function mapFilterJson(json) { return json || {}; }
function unmapKey(k) { return k; }

function BacklogView({ nav, issues }) {
  const [editing, setEditing] = React.useState(null);
  const [filter, setFilter] = React.useState(EMPTY_FILTER);
  const [activeSaved, setActiveSaved] = React.useState(null); // saved filter id
  const { data: saved, reload: reloadSaved, setData: setSaved } = useApi("/saved-filters");
  const [saveOpen, setSaveOpen] = React.useState(false);
  const [savedOpen, setSavedOpen] = React.useState(false);
  const [name, setName] = React.useState("");
  const savedRef = React.useRef(null);
  const toast = useToast();

  React.useEffect(() => {
    if (!savedOpen) return;
    const h = (e) => { if (savedRef.current && !savedRef.current.contains(e.target)) setSavedOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [savedOpen]);

  const labels = React.useMemo(() => Array.from(new Set(issues.flatMap((i) => i.labels))), [issues]);
  const sprints = React.useMemo(() => Array.from(new Set(issues.map((i) => i.sprint).filter(Boolean))), [issues]);

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
    setFilter({ ...EMPTY_FILTER, ...mapFilterJson(sf.filter_json) });
    setActiveSaved(sf.id);
    setSavedOpen(false);
  }
  async function saveFilter() {
    if (!name.trim()) return;
    const filter_json = {};
    Object.entries(filter).forEach(([k, v]) => { if (v !== "all") filter_json[unmapKey(k)] = v; });
    try {
      const sf = await api("/saved-filters", { method: "POST", body: { name, filter_json } });
      toast("Filter saved");
      setSaveOpen(false); setName(""); setActiveSaved(sf.id); reloadSaved();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function delSaved(id, e) {
    e.stopPropagation();
    setSaved((saved || []).filter((s) => s.id !== id));
    if (activeSaved === id) setActiveSaved(null);
    try { await api("/saved-filters/" + id, { method: "DELETE" }); toast("Filter deleted"); }
    catch (err) { toast(err.message, { icon: "x", color: "#F87171" }); reloadSaved(); }
  }

  const activeName = activeSaved && (saved || []).find((s) => s.id === activeSaved);

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Backlog</span></div>
          <h1>Backlog</h1>
          <p>{filtered.length} of {issues.length} issues · drag to plan into upcoming sprints.</p>
        </div>
        <div className="row gap-2">
          <Button icon="download">Export CSV</Button>
          <Button variant="primary" icon="plus" onClick={() => window.dispatchEvent(new CustomEvent("forge:create"))}>Create</Button>
        </div>
      </div>

      {/* Filter bar */}
      <div className="row gap-2" style={{ padding: "0 32px 14px", flexWrap: "wrap" }}>
        <Pill label="Status" icon="kanban" value={filter.status === "all" ? "All" : filter.status} options={[{ id: "all", label: "All" }, ...FORGE_DATA.COLUMNS.map((c) => ({ id: c.id, label: c.id }))]} onChange={(v) => set("status", v)}/>
        <Pill label="Assignee" icon="user" value={filter.assignee === "all" ? "All" : filter.assignee === "none" ? "Unassigned" : FORGE_DATA.PEOPLE.find((p) => p.id === filter.assignee)?.name} options={[{ id: "all", label: "All" }, { id: "none", label: "Unassigned" }, ...FORGE_DATA.PEOPLE.map((p) => ({ id: p.id, label: p.name }))]} onChange={(v) => set("assignee", v)}/>
        <Pill label="Priority" icon="flag" value={filter.priority === "all" ? "All" : filter.priority} options={[{ id: "all", label: "All" }, ...Object.keys(FORGE_DATA.PRIORITY_META).map((p) => ({ id: p, label: p }))]} onChange={(v) => set("priority", v)}/>
        <Pill label="Type" icon="tag" value={filter.type === "all" ? "All" : filter.type} options={[{ id: "all", label: "All" }, ...Object.keys(FORGE_DATA.TYPE_META).map((t) => ({ id: t, label: t }))]} onChange={(v) => set("type", v)}/>
        <Pill label="Label" icon="tag" value={filter.label === "all" ? "All" : filter.label} options={[{ id: "all", label: "All" }, ...labels.map((l) => ({ id: l, label: l }))]} onChange={(v) => set("label", v)}/>
        <Pill label="Sprint" icon="rocket" value={filter.sprint === "all" ? "All" : filter.sprint} options={[{ id: "all", label: "All" }, ...sprints.map((s) => ({ id: s, label: s }))]} onChange={(v) => set("sprint", v)}/>

        <div style={{ flex: 1 }}/>

        {/* Active saved filter chip */}
        {activeName && (
          <span className="badge" data-tone="info" style={{ height: 28, padding: "0 8px", gap: 6 }}>
            <Icon name="star" size={12}/> {activeName.name}
            <button className="icon-btn" style={{ width: 18, height: 18 }} onClick={() => { setFilter(EMPTY_FILTER); setActiveSaved(null); }}><Icon name="x" size={11}/></button>
          </span>
        )}

        {/* Saved filters dropdown */}
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
        <Button data-size="sm" variant="secondary" icon="star" disabled={!dirty} onClick={() => setSaveOpen(true)} title="Save current filter">Save filter</Button>
      </div>

      <Modal open={saveOpen} onClose={() => setSaveOpen(false)} title="Save filter"
        footer={<><Button onClick={() => setSaveOpen(false)}>Cancel</Button><Button variant="primary" disabled={!name.trim()} onClick={saveFilter}>Save</Button></>}>
        <label className="label">Filter name</label>
        <input className="input" autoFocus placeholder="e.g. My critical bugs" value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => { if (e.key === "Enter") saveFilter(); }}/>
        <div className="help" style={{ marginTop: 8 }}>Saves the {Object.values(filter).filter((v) => v !== "all").length} active filter condition(s) currently applied.</div>
      </Modal>

      <div style={{ padding: "0 32px 32px" }}>
        <div style={{ border: "1px solid var(--border)", borderRadius: 10, overflow: "hidden", background: "var(--bg)" }}>
          <div style={{ padding: "10px 16px", borderBottom: "1px solid var(--border)", background: "var(--bg-subtle)", display: "flex", alignItems: "center", gap: 12 }}>
            <Icon name="rocket" size={14} color="var(--indigo-600)"/>
            <span className="bold">Sprint 25 — Observability deep clean</span>
            <Badge tone="muted">Planned · Dec 16 → Dec 29</Badge>
            <span className="text-xs muted">23 issues · 52 points</span>
            <div style={{ flex: 1 }}/>
            <Button data-size="sm">Start sprint</Button>
          </div>
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
              {filtered.map((i) => {
                const u = FORGE_DATA.PEOPLE.find((p) => p.id === i.assignee);
                return (
                  <tr key={i.id}>
                    <td><TypeIcon value={i.type}/></td>
                    <td><a href={"#/issue/" + i.id} className="mono text-xs" onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.id}</a></td>
                    <td onClick={() => setEditing(i.id)}>
                      {editing === i.id ? (
                        <input className="input" autoFocus defaultValue={i.title} onBlur={() => setEditing(null)} onKeyDown={(e) => { if (e.key === "Enter") setEditing(null); }} style={{ padding: "2px 6px" }}/>
                      ) : (
                        <a href={"#/issue/" + i.id} onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.title}</a>
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
      </div>
    </div>
  );
}

// ─── My issues ──────────────────────────────────────────
function MyIssuesView({ nav, issues }) {
  const mine = issues.filter((i) => i.assignee === FORGE_DATA.ME.id);
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
                  <td><a href={"#/issue/" + i.id} className="mono text-xs" onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.id}</a></td>
                  <td><a href={"#/issue/" + i.id} onClick={(e) => { e.preventDefault(); nav("issue/" + i.id); }}>{i.title}</a></td>
                  <td><StatusBadge value={i.status}/></td>
                  <td><PriorityBadge value={i.pri}/></td>
                  <td className="text-xs muted">{i.due || "—"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

Object.assign(window, { IssueView, BacklogView, MyIssuesView });
