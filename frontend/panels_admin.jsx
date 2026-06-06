// panels_admin.jsx — AdminView panels: API keys (12), Audit log (13), Import (26)

// ─── Feature 12: API keys ─────────────────────────────────────────────
function ApiKeysPanel() {
  const { data: keys, loading, error, reload } = useApi("/api-keys");
  const [genOpen, setGenOpen] = React.useState(false);
  const [name, setName] = React.useState("");
  const [created, setCreated] = React.useState(null); // {plain_key,...}
  const [confirm, setConfirm] = React.useState(null);
  const [copied, setCopied] = React.useState(false);
  const toast = useToast();

  async function generate() {
    if (!name.trim()) return;
    try {
      const k = await api("/api-keys", { method: "POST", body: { name } });
      setCreated(k); setName(""); setGenOpen(false); reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function revoke(id) { try { await api("/api-keys/" + id, { method: "DELETE" }); toast("API key revoked"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  function copyKey() { navigator.clipboard && navigator.clipboard.writeText(created.plain_key); setCopied(true); setTimeout(() => setCopied(false), 1500); }

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Programmatic access tokens for CI, Terraform, and other services.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setGenOpen(true)}>Generate API key</Button>
      </div>

      {loading ? <div className="row gap-2 text-sm muted" style={{ padding: 24 }}><MiniSpinner/> Loading keys…</div> : error ? (
        <div className="card card-pad text-sm" style={{ color: "var(--danger)" }}>{error}</div>
      ) : (
        <div className="card" style={{ overflow: "hidden" }}>
          <table className="table">
            <thead><tr><th>Name</th><th style={{ width: 160 }}>Prefix</th><th style={{ width: 140 }}>Created</th><th style={{ width: 140 }}>Last used</th><th style={{ width: 100 }}/></tr></thead>
            <tbody>
              {(keys || []).map((k) => (
                <tr key={k.id}>
                  <td className="bold">{k.name}</td>
                  <td><span className="mono text-xs">{k.prefix}…</span></td>
                  <td className="text-xs muted">{fmtDate(k.created_at)}</td>
                  <td className="text-xs muted">{k.last_used_at ? fmtDate(k.last_used_at) : "Never"}</td>
                  <td><Button data-size="sm" variant="ghost" style={{ color: "var(--danger)" }} onClick={() => setConfirm(k)}>Revoke</Button></td>
                </tr>
              ))}
            </tbody>
          </table>
          {(keys || []).length === 0 && <Empty icon="lock" title="No API keys" hint="Generate a key to access the API."/>}
        </div>
      )}

      {/* Generate modal */}
      <Modal open={genOpen} onClose={() => setGenOpen(false)} title="Generate API key"
        footer={<><Button onClick={() => setGenOpen(false)}>Cancel</Button><Button variant="primary" disabled={!name.trim()} onClick={generate}>Generate</Button></>}>
        <label className="label">Key name</label>
        <input className="input" autoFocus value={name} onChange={(e) => setName(e.target.value)} placeholder="e.g. CI pipeline" onKeyDown={(e) => { if (e.key === "Enter") generate(); }}/>
        <div className="help" style={{ marginTop: 8 }}>Give the key a descriptive name so you can identify it later.</div>
      </Modal>

      {/* Show-once modal */}
      <Modal open={!!created} onClose={() => setCreated(null)} title="API key created"
        footer={<Button variant="primary" onClick={() => setCreated(null)}>Done</Button>}>
        <div className="card card-pad" style={{ background: "var(--warning-bg)", borderColor: "var(--warning)", marginBottom: 14 }}>
          <div className="row gap-2 text-sm" style={{ color: "#B45309" }}><Icon name="shield" size={15}/><span className="bold">Save this key — it won't be shown again.</span></div>
        </div>
        {created && (
          <>
            <label className="label">{created.name}</label>
            <div className="row gap-2" style={{ padding: "10px 12px", background: "var(--bg-subtle)", border: "1px solid var(--border)", borderRadius: 8 }}>
              <span className="mono text-xs" style={{ flex: 1, wordBreak: "break-all" }}>{created.plain_key}</span>
              <Button data-size="sm" icon={copied ? "check" : "copy"} onClick={copyKey}>{copied ? "Copied" : "Copy"}</Button>
            </div>
          </>
        )}
      </Modal>

      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => revoke(confirm.id)} title="Revoke API key" body={confirm ? `Revoke "${confirm.name}"? Any service using it will lose access immediately.` : ""}/>
    </div>
  );
}

// ─── Feature 13: Audit log ────────────────────────────────────────────
const ACTION_CAT_TONE = { create: "success", update: "info", delete: "danger", auth: "purple" };
function AuditLogPanel() {
  const [filters, setFilters] = React.useState({ user_id: "", action: "", from: "", to: "" });
  const [page, setPage] = React.useState(1);
  const [rows, setRows] = React.useState([]);
  const [total, setTotal] = React.useState(0);
  const [loading, setLoading] = React.useState(true);
  const toast = useToast();
  const limit = 50;

  const qs = React.useMemo(() => {
    const p = new URLSearchParams({ page: String(page), limit: String(limit) });
    Object.entries(filters).forEach(([k, v]) => { if (v) p.set(k, v); });
    return p.toString();
  }, [filters, page]);

  React.useEffect(() => {
    let live = true; setLoading(true);
    api("/audit-logs?" + qs).then((d) => {
      if (!live) return;
      setRows((prev) => page === 1 ? d.items : [...prev, ...d.items]);
      setTotal(d.total); setLoading(false);
    }).catch((e) => { if (live) { toast(e.message, { icon: "x", color: "#F87171" }); setLoading(false); } });
    return () => { live = false; };
  }, [qs]);

  function setF(k, v) { setFilters((f) => ({ ...f, [k]: v })); setPage(1); }
  async function exportCsv() { try { await apiDownload("/audit-logs/export"); toast("Exporting audit log…"); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  return (
    <div>
      {/* Filter bar */}
      <div className="card card-pad" style={{ marginBottom: 12 }}>
        <div className="row gap-3" style={{ flexWrap: "wrap", alignItems: "flex-end" }}>
          <div><label className="label">From</label><input className="input" type="date" value={filters.from} onChange={(e) => setF("from", e.target.value)} style={{ width: 150 }}/></div>
          <div><label className="label">To</label><input className="input" type="date" value={filters.to} onChange={(e) => setF("to", e.target.value)} style={{ width: 150 }}/></div>
          <div><label className="label">Action</label>
            <select className="select" value={filters.action} onChange={(e) => setF("action", e.target.value)} style={{ width: 150 }}>
              <option value="">All actions</option>
              {Object.keys(ACTION_CAT_TONE).map((c) => <option key={c} value={c}>{c}</option>)}
            </select>
          </div>
          <div><label className="label">User</label>
            <select className="select" value={filters.user_id} onChange={(e) => setF("user_id", e.target.value)} style={{ width: 170 }}>
              <option value="">All users</option>
              {FORGE_DATA.PEOPLE.map((u) => <option key={u.id} value={u.id}>{u.name}</option>)}
            </select>
          </div>
          <div style={{ flex: 1 }}/>
          {(filters.user_id || filters.action || filters.from || filters.to) && <Button data-size="sm" variant="ghost" onClick={() => { setFilters({ user_id: "", action: "", from: "", to: "" }); setPage(1); }}>Clear</Button>}
          <Button data-size="sm" icon="download" onClick={exportCsv}>Export CSV</Button>
        </div>
      </div>

      <div className="card" style={{ overflow: "hidden" }}>
        <table className="table">
          <thead><tr><th style={{ width: 130 }}>When</th><th style={{ width: 200 }}>Actor</th><th style={{ width: 160 }}>Action</th><th>Resource</th><th style={{ width: 120 }}>IP</th></tr></thead>
          <tbody>
            {rows.map((l) => (
              <tr key={l.id}>
                <td className="text-xs muted">{fmtDate(l.created_at)}</td>
                <td>{l.actor ? <span className="row gap-2"><Avatar user={l.actor} size="sm"/>{l.actor.name}</span> : l.actor_id}</td>
                <td><Badge tone={ACTION_CAT_TONE[l.category] || "muted"} dot>{l.action}</Badge></td>
                <td><span className="text-sm"><span className="mono text-xs muted">{l.resource_type}</span> {l.resource_id}</span></td>
                <td className="mono text-xs muted">{(l.meta || {}).ip || "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {loading && <div className="row gap-2 text-sm muted" style={{ padding: 16 }}><MiniSpinner/> Loading…</div>}
        {!loading && rows.length === 0 && <Empty icon="history" title="No audit events" hint="Try adjusting the filters."/>}
        {!loading && rows.length < total && (
          <div style={{ padding: 12, textAlign: "center", borderTop: "1px solid var(--border)" }}>
            <Button data-size="sm" onClick={() => setPage((p) => p + 1)}>Load more <span className="muted">({rows.length} of {total})</span></Button>
          </div>
        )}
      </div>
    </div>
  );
}

// ─── Feature 26: Data import ──────────────────────────────────────────
const IMPORT_SOURCES = [
  { key: "jira", name: "Jira", format: "XML", icon: "briefcase", desc: "Import your Jira projects, issues, and workflows." },
  { key: "trello", name: "Trello", format: "JSON", icon: "layout", desc: "Bring boards and cards over from Trello." },
  { key: "linear", name: "Linear", format: "CSV", icon: "list", desc: "Migrate issues and teams from Linear." },
];
function ImportPanel() {
  return (
    <div>
      <div className="text-sm secondary" style={{ marginBottom: 14 }}>Migrate from another tool. Files are processed in the background.</div>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 14 }}>
        {IMPORT_SOURCES.map((s) => <ImportCard key={s.key} source={s}/>)}
      </div>
    </div>
  );
}

function ImportCard({ source }) {
  const [file, setFile] = React.useState(null);
  const [job, setJob] = React.useState(null);
  const [dragOver, setDragOver] = React.useState(false);
  const inputRef = React.useRef(null);
  const pollRef = React.useRef(null);
  const toast = useToast();

  React.useEffect(() => () => clearInterval(pollRef.current), []);

  async function start() {
    if (!file) return;
    try {
      const res = await apiUpload("/import/" + source.key, file);
      setJob({ id: res.id, status: res.status || "pending", progress: 0 });
      clearInterval(pollRef.current);
      pollRef.current = setInterval(() => poll(res.id), 3000);
      poll(res.id);
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function poll(id) {
    try {
      const j = await api("/import/" + id);
      setJob(j);
      if (j.status === "done" || j.status === "failed") {
        clearInterval(pollRef.current);
        toast(j.status === "done" ? source.name + " import complete" : source.name + " import failed", j.status === "failed" ? { icon: "x", color: "#F87171" } : {});
      }
    } catch (e) { clearInterval(pollRef.current); }
  }
  function reset() { clearInterval(pollRef.current); setJob(null); setFile(null); }

  const running = job && (job.status === "pending" || job.status === "processing");

  return (
    <div className="card card-pad">
      <div className="row gap-3" style={{ marginBottom: 10 }}>
        <div style={{ width: 36, height: 36, borderRadius: 8, background: "var(--bg-subtle)", border: "1px solid var(--border)", display: "grid", placeItems: "center", color: "var(--indigo-600)" }}><Icon name={source.icon} size={18}/></div>
        <div className="stack" style={{ lineHeight: 1.25 }}><span className="bold">{source.name}</span><span className="text-xs muted">{source.format} export</span></div>
      </div>
      <div className="text-sm secondary" style={{ minHeight: 36, marginBottom: 10 }}>{source.desc}</div>

      {!job && (
        <>
          <div onClick={() => inputRef.current && inputRef.current.click()}
            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
            onDragLeave={() => setDragOver(false)}
            onDrop={(e) => { e.preventDefault(); setDragOver(false); if (e.dataTransfer.files[0]) setFile(e.dataTransfer.files[0]); }}
            style={{ border: "1.5px dashed " + (dragOver ? "var(--indigo-500)" : "var(--border-strong)"), borderRadius: 8, padding: "18px 12px", textAlign: "center", background: dragOver ? "var(--indigo-50)" : "var(--bg-subtle)", marginBottom: 10, cursor: "default" }}>
            <Icon name="upload" size={20} color="var(--text-muted)"/>
            <div className="text-xs" style={{ marginTop: 6 }}>{file ? <span className="medium">{file.name}</span> : <>Drop {source.format} file or <span style={{ color: "var(--indigo-600)" }}>browse</span></>}</div>
            <input ref={inputRef} type="file" hidden onChange={(e) => setFile(e.target.files[0])}/>
          </div>
          <Button variant="primary" data-size="sm" icon="upload" disabled={!file} onClick={start} style={{ width: "100%", justifyContent: "center" }}>Start import</Button>
        </>
      )}

      {job && (
        <div>
          <div className="row gap-2" style={{ marginBottom: 8 }}>
            {running && <MiniSpinner/>}
            <Badge tone={job.status === "done" ? "success" : job.status === "failed" ? "danger" : "info"} dot>{job.status}</Badge>
            <span className="text-xs muted" style={{ marginLeft: "auto" }}>{job.progress != null ? job.progress + "%" : ""}</span>
          </div>
          <div style={{ height: 8, background: "var(--bg-muted)", borderRadius: 99, overflow: "hidden", marginBottom: 10 }}>
            <div style={{ width: (job.progress || 0) + "%", height: "100%", background: job.status === "failed" ? "var(--danger)" : job.status === "done" ? "var(--success)" : "var(--indigo-600)", transition: "width .4s" }}/>
          </div>
          {job.status === "done" && job.summary && (
            <div className="text-sm" style={{ marginBottom: 10 }}>{job.summary.issues_created} issues · {job.summary.projects_created} projects created.</div>
          )}
          {job.status === "failed" && (
            <div className="text-xs" style={{ color: "var(--danger)", marginBottom: 10 }}>{job.error}</div>
          )}
          {!running && <Button data-size="sm" onClick={reset} style={{ width: "100%", justifyContent: "center" }}>{job.status === "failed" ? "Try again" : "Import another"}</Button>}
        </div>
      )}
    </div>
  );
}

Object.assign(window, { ApiKeysPanel, AuditLogPanel, ImportPanel });
