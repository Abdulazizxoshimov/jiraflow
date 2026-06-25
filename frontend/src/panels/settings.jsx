import { useState } from 'react';
import { Icon } from '../components/icons';
import { Avatar, Badge, Button, Switch, Modal, Empty, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api, useApi } from '../api/api';
import { adaptUser, fmtDate } from '../api/adapters';
import { PillSelect } from '../views/board';
import { MiniSpinner } from './issue';

const PRESET_COLORS = ["#6366F1","#06B6D4","#10B981","#F59E0B","#EF4444","#8B5CF6","#EC4899","#14B8A6","#F97316","#3B82F6","#64748B","#A855F7"];
const CAT_META = { todo: { label: "To do", tone: "muted" }, in_progress: { label: "In progress", tone: "info" }, done: { label: "Done", tone: "success" } };

function Loading({ label }) {
  return <div className="row gap-2 text-sm muted" style={{ padding: 24 }}><MiniSpinner/> {label || "Loading…"}</div>;
}
function ErrorNote({ msg }) {
  return <div className="card card-pad text-sm" style={{ color: "var(--danger)", borderColor: "var(--danger)" }}><Icon name="x" size={14}/> {msg}</div>;
}
export function ColorPicker({ value, onChange }) {
  return (
    <div className="row gap-2" style={{ flexWrap: "wrap" }}>
      {PRESET_COLORS.map((c) => (
        <button key={c} onClick={() => onChange(c)} title={c}
          style={{ width: 26, height: 26, borderRadius: 7, background: c, border: value === c ? "2px solid var(--text)" : "2px solid transparent", boxShadow: value === c ? "0 0 0 2px var(--bg) inset" : "none" }}/>
      ))}
    </div>
  );
}
export function ConfirmDelete({ open, onClose, onConfirm, title, body }) {
  return (
    <Modal open={open} onClose={onClose} title={title || "Confirm delete"}
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="danger" icon="trash" onClick={() => { onConfirm(); onClose(); }}>Delete</Button></>}>
      <p className="text-sm secondary" style={{ margin: 0 }}>{body || "This action cannot be undone."}</p>
    </Modal>
  );
}

// ─── Feature 7: Workflow management ───────────────────────────────────
export function WorkflowTab() {
  const { data: workflows, loading, error, reload } = useApi("/workflows");
  const [expanded, setExpanded] = useState(null);
  const [newOpen, setNewOpen] = useState(false);
  const [name, setName] = useState("");
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  async function create() {
    if (!name.trim()) return;
    try { await api("/workflows", { method: "POST", body: { name } }); toast("Workflow created"); setNewOpen(false); setName(""); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function del(id) {
    try { await api("/workflows/" + id, { method: "DELETE" }); toast("Workflow deleted"); if (expanded === id) setExpanded(null); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  if (loading) return <Loading label="Loading workflows…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div className="stack gap-3">
      <div className="row" style={{ justifyContent: "space-between" }}>
        <div className="text-sm secondary">Define the statuses and transitions issues move through.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setNewOpen(true)}>New workflow</Button>
      </div>
      {(workflows || []).map((w) => (
        <div key={w.id} className="card" style={{ overflow: "hidden" }}>
          <div className="row gap-3" style={{ padding: "12px 16px", cursor: "default" }} onClick={() => setExpanded((e) => e === w.id ? null : w.id)}>
            <Icon name={expanded === w.id ? "chevronDown" : "chevronRight"} size={14} color="var(--text-muted)"/>
            <span className="bold">{w.name}</span>
            <Badge tone="muted">{w.status_count} statuses</Badge>
            <div style={{ flex: 1 }}/>
            <button className="icon-btn" title="Delete workflow" onClick={(e) => { e.stopPropagation(); setConfirm(w); }}><Icon name="trash" size={15}/></button>
          </div>
          {expanded === w.id && <WorkflowDetail id={w.id}/>}
        </div>
      ))}
      {(workflows || []).length === 0 && <Empty icon="layout" title="No workflows" hint="Create your first workflow to get started."/>}

      <Modal open={newOpen} onClose={() => setNewOpen(false)} title="New workflow"
        footer={<><Button onClick={() => setNewOpen(false)}>Cancel</Button><Button variant="primary" disabled={!name.trim()} onClick={create}>Create</Button></>}>
        <label className="label">Workflow name</label>
        <input className="input" autoFocus placeholder="e.g. Incident response" value={name} onChange={(e) => setName(e.target.value)} onKeyDown={(e) => { if (e.key === "Enter") create(); }}/>
      </Modal>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} title="Delete workflow" body={confirm ? `Delete "${confirm.name}" and all its statuses and transitions?` : ""}/>
    </div>
  );
}

function WorkflowDetail({ id }) {
  const { data: wf, loading, reload } = useApi("/workflows/" + id);
  const [addOpen, setAddOpen] = useState(false);
  const [sName, setSName] = useState(""); const [sCat, setSCat] = useState("todo"); const [sColor, setSColor] = useState(PRESET_COLORS[0]);
  const [trFrom, setTrFrom] = useState(""); const [trTo, setTrTo] = useState("");
  const toast = useToast();

  if (loading || !wf) return <div style={{ borderTop: "1px solid var(--border)" }}><Loading/></div>;
  const statuses = wf.statuses || [];
  const nameOf = (sid) => (statuses.find((s) => s.id === sid) || {}).name || "?";

  async function addStatus() {
    if (!sName.trim()) return;
    try { await api("/workflows/" + id + "/statuses", { method: "POST", body: { name: sName, category: sCat, color: sColor } }); toast("Status added"); setAddOpen(false); setSName(""); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function delStatus(sid) {
    try { await api("/workflows/statuses/" + sid, { method: "DELETE" }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function move(idx, dir) {
    const order = statuses.map((s) => s.id);
    const j = idx + dir; if (j < 0 || j >= order.length) return;
    [order[idx], order[j]] = [order[j], order[idx]];
    try { await api("/workflows/" + id, { method: "PUT", body: { status_order: order } }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function addTransition() {
    if (!trFrom || !trTo || trFrom === trTo) return;
    try { await api("/workflows/" + id + "/transitions", { method: "POST", body: { from_id: trFrom, to_id: trTo } }); toast("Transition added"); setTrFrom(""); setTrTo(""); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function delTransition(tid) {
    try { await api("/workflows/transitions/" + tid, { method: "DELETE" }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div style={{ borderTop: "1px solid var(--border)", padding: 16, background: "var(--bg-subtle)" }}>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 8 }}>
        <h4 style={{ margin: 0, fontSize: 12, fontWeight: 600, textTransform: "uppercase", letterSpacing: ".04em", color: "var(--text-secondary)" }}>Statuses</h4>
        <Button data-size="sm" icon="plus" onClick={() => setAddOpen(true)}>Add status</Button>
      </div>
      <div className="stack gap-1" style={{ marginBottom: 16 }}>
        {statuses.map((s, i) => (
          <div key={s.id} className="row gap-3" style={{ padding: "8px 10px", background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
            <span style={{ width: 10, height: 10, borderRadius: 3, background: s.color }}/>
            <span className="medium text-sm" style={{ flex: 1 }}>{s.name}</span>
            <Badge tone={CAT_META[s.category].tone}>{CAT_META[s.category].label}</Badge>
            <div className="row gap-1">
              <button className="icon-btn" style={{ width: 22, height: 22 }} disabled={i === 0} onClick={() => move(i, -1)} title="Move up"><Icon name="arrowUp" size={12}/></button>
              <button className="icon-btn" style={{ width: 22, height: 22 }} disabled={i === statuses.length - 1} onClick={() => move(i, 1)} title="Move down"><Icon name="arrowDown" size={12}/></button>
              <button className="icon-btn" style={{ width: 22, height: 22 }} onClick={() => delStatus(s.id)} title="Delete"><Icon name="trash" size={12}/></button>
            </div>
          </div>
        ))}
      </div>

      <h4 style={{ margin: "0 0 8px", fontSize: 12, fontWeight: 600, textTransform: "uppercase", letterSpacing: ".04em", color: "var(--text-secondary)" }}>Transitions</h4>
      <div className="stack gap-1" style={{ marginBottom: 12 }}>
        {(wf.transitions || []).map((t) => (
          <div key={t.id} className="row gap-2" style={{ padding: "6px 10px", background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 6 }}>
            <span className="text-sm">{nameOf(t.from_id)}</span>
            <Icon name="arrowRight" size={13} color="var(--text-muted)"/>
            <span className="text-sm">{nameOf(t.to_id)}</span>
            <div style={{ flex: 1 }}/>
            <button className="icon-btn" style={{ width: 22, height: 22 }} onClick={() => delTransition(t.id)} title="Delete"><Icon name="x" size={12}/></button>
          </div>
        ))}
        {(wf.transitions || []).length === 0 && <div className="text-xs muted">No transitions defined.</div>}
      </div>
      <div className="row gap-2">
        <select className="select" value={trFrom} onChange={(e) => setTrFrom(e.target.value)} style={{ flex: 1 }}>
          <option value="">From…</option>{statuses.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
        </select>
        <Icon name="arrowRight" size={14} color="var(--text-muted)"/>
        <select className="select" value={trTo} onChange={(e) => setTrTo(e.target.value)} style={{ flex: 1 }}>
          <option value="">To…</option>{statuses.map((s) => <option key={s.id} value={s.id}>{s.name}</option>)}
        </select>
        <Button data-size="sm" variant="primary" disabled={!trFrom || !trTo || trFrom === trTo} onClick={addTransition}>Add</Button>
      </div>

      <Modal open={addOpen} onClose={() => setAddOpen(false)} title="Add status"
        footer={<><Button onClick={() => setAddOpen(false)}>Cancel</Button><Button variant="primary" disabled={!sName.trim()} onClick={addStatus}>Add status</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Name</label><input className="input" autoFocus value={sName} onChange={(e) => setSName(e.target.value)} placeholder="e.g. Blocked"/></div>
          <div><label className="label">Category</label>
            <div className="row gap-2">
              {Object.entries(CAT_META).map(([k, m]) => (
                <button key={k} className="btn" onClick={() => setSCat(k)} style={{ flex: 1, justifyContent: "center", border: "1px solid " + (sCat === k ? "var(--indigo-600)" : "var(--border)"), background: sCat === k ? "var(--indigo-50)" : "var(--bg)", color: sCat === k ? "var(--indigo-700)" : "var(--text)" }}>{m.label}</button>
              ))}
            </div>
          </div>
          <div><label className="label">Color</label><ColorPicker value={sColor} onChange={setSColor}/></div>
        </div>
      </Modal>
    </div>
  );
}

// ─── Feature 8a: Components ───────────────────────────────────────────
export function ComponentsTab() {
  const { activeProjectId, people } = useApp();
  const { data: comps, loading, error, reload } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/components" : null);
  const [modal, setModal] = useState(false);
  const [editing, setEditing] = useState(null);
  const [form, setForm] = useState({ name: "", description: "", lead_id: "" });
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  function open(c) { setEditing(c); setForm(c ? { name: c.name, description: c.description, lead_id: c.lead_id || "" } : { name: "", description: "", lead_id: "" }); setModal(true); }
  async function save() {
    if (!form.name.trim()) return;
    try {
      if (editing) await api("/components/" + editing.id, { method: "PUT", body: form });
      else await api("/projects/" + activeProjectId + "/components", { method: "POST", body: form });
      toast(editing ? "Component updated" : "Component added"); setModal(false); reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function del(id) { try { await api("/components/" + id, { method: "DELETE" }); toast("Component deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  if (loading) return <Loading label="Loading components…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Group issues by area of the codebase or system.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => open(null)}>Add component</Button>
      </div>
      <div className="card" style={{ overflow: "hidden" }}>
        <table className="table">
          <thead><tr><th>Name</th><th>Description</th><th style={{ width: 200 }}>Lead</th><th style={{ width: 90 }}>Issues</th><th style={{ width: 80 }}/></tr></thead>
          <tbody>
            {(comps || []).map((c) => (
              <tr key={c.id}>
                <td className="bold">{c.name}</td>
                <td className="text-sm secondary">{c.description || "—"}</td>
                <td>{c.lead ? <span className="row gap-2"><Avatar user={c.lead} size="sm"/>{c.lead.name}</span> : <span className="muted">Unassigned</span>}</td>
                <td>{c.issue_count}</td>
                <td><div className="row gap-1">
                  <button className="icon-btn" onClick={() => open(c)} title="Edit"><Icon name="pencil" size={14}/></button>
                  <button className="icon-btn" onClick={() => setConfirm(c)} title="Delete"><Icon name="trash" size={14}/></button>
                </div></td>
              </tr>
            ))}
          </tbody>
        </table>
        {(comps || []).length === 0 && <Empty icon="folder" title="No components" hint="Add a component to organize issues."/>}
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title={editing ? "Edit component" : "Add component"}
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" disabled={!form.name.trim()} onClick={save}>Save</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Name</label><input className="input" autoFocus value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. API gateway"/></div>
          <div><label className="label">Description</label><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} style={{ minHeight: 60 }}/></div>
          <div><label className="label">Lead</label><PillSelect value={form.lead_id} onChange={(v) => setForm({ ...form, lead_id: v })} options={people.map((u) => ({ id: u.id, label: u.name }))}/></div>
        </div>
      </Modal>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} body={confirm ? `Delete component "${confirm.name}"?` : ""}/>
    </div>
  );
}

// ─── Feature 8b: Versions ─────────────────────────────────────────────
export function VersionsTab() {
  const { activeProjectId } = useApp();
  const { data: versions, loading, error, reload } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/versions" : null);
  const [modal, setModal] = useState(false);
  const [form, setForm] = useState({ name: "", start_date: "", release_date: "", description: "" });
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();
  const tone = { unreleased: "info", released: "success", archived: "muted" };

  async function save() {
    if (!form.name.trim()) return;
    try { await api("/projects/" + activeProjectId + "/versions", { method: "POST", body: form }); toast("Version created"); setModal(false); setForm({ name: "", start_date: "", release_date: "", description: "" }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function release(id) { try { await api("/versions/" + id + "/release", { method: "POST" }); toast("Version released"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  async function del(id) { try { await api("/versions/" + id, { method: "DELETE" }); toast("Version deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  if (loading) return <Loading label="Loading versions…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Track releases and the issues that ship in them.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setModal(true)}>Create version</Button>
      </div>
      <div className="card" style={{ overflow: "hidden" }}>
        <table className="table">
          <thead><tr><th>Name</th><th style={{ width: 120 }}>Start</th><th style={{ width: 120 }}>Release</th><th style={{ width: 120 }}>Status</th><th style={{ width: 160 }}/></tr></thead>
          <tbody>
            {(versions || []).map((v) => (
              <tr key={v.id}>
                <td><div className="bold">{v.name}</div>{v.description && <div className="text-xs muted">{v.description}</div>}</td>
                <td className="text-sm secondary">{v.start_date || "—"}</td>
                <td className="text-sm secondary">{v.release_date || "—"}</td>
                <td><Badge tone={tone[v.status]} dot>{v.status}</Badge></td>
                <td><div className="row gap-1" style={{ justifyContent: "flex-end" }}>
                  {v.status === "unreleased" && <Button data-size="sm" icon="rocket" onClick={() => release(v.id)}>Release</Button>}
                  <button className="icon-btn" onClick={() => setConfirm(v)} title="Delete"><Icon name="trash" size={14}/></button>
                </div></td>
              </tr>
            ))}
          </tbody>
        </table>
        {(versions || []).length === 0 && <Empty icon="rocket" title="No versions" hint="Create a version to plan a release."/>}
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title="Create version"
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" disabled={!form.name.trim()} onClick={save}>Create</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Name</label><input className="input" autoFocus value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. 2025.02 — Q1 hardening"/></div>
          <div className="row gap-3">
            <div style={{ flex: 1 }}><label className="label">Start date</label><input className="input" type="date" value={form.start_date} onChange={(e) => setForm({ ...form, start_date: e.target.value })}/></div>
            <div style={{ flex: 1 }}><label className="label">Release date</label><input className="input" type="date" value={form.release_date} onChange={(e) => setForm({ ...form, release_date: e.target.value })}/></div>
          </div>
          <div><label className="label">Description</label><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} style={{ minHeight: 60 }}/></div>
        </div>
      </Modal>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} body={confirm ? `Delete version "${confirm.name}"?` : ""}/>
    </div>
  );
}

// ─── Feature 9: Custom fields ─────────────────────────────────────────
const FIELD_TYPES = ["text", "number", "date", "select", "user", "url"];
export function CustomFieldsTab() {
  const { activeProjectId } = useApp();
  const { data: fields, loading, error, reload, setData } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/custom-fields" : null);
  const [modal, setModal] = useState(false);
  const [form, setForm] = useState({ name: "", type: "text", description: "", required: false, options: [] });
  const [optInput, setOptInput] = useState("");
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  async function save() {
    if (!form.name.trim()) return;
    try { await api("/projects/" + activeProjectId + "/custom-fields", { method: "POST", body: form }); toast("Field added"); setModal(false); setForm({ name: "", type: "text", description: "", required: false, options: [] }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function toggleReq(f) {
    setData((fields || []).map((x) => x.id === f.id ? { ...x, required: !x.required } : x));
    try { await api("/custom-fields/" + f.id, { method: "PUT", body: { required: !f.required } }); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }
  async function move(idx, dir) {
    const arr = (fields || []).slice(); const j = idx + dir; if (j < 0 || j >= arr.length) return;
    [arr[idx], arr[j]] = [arr[j], arr[idx]]; setData(arr);
    try { await api("/custom-fields/reorder", { method: "PUT", body: { reorder: arr.map((f) => f.id) } }); } catch (e) { reload(); }
  }
  async function del(id) { try { await api("/custom-fields/" + id, { method: "DELETE" }); toast("Field deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  if (loading) return <Loading label="Loading custom fields…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Capture extra structured data on every issue.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setModal(true)}>Add field</Button>
      </div>
      <div className="card" style={{ overflow: "hidden" }}>
        <table className="table">
          <thead><tr><th style={{ width: 70 }}>Order</th><th>Name</th><th style={{ width: 110 }}>Type</th><th>Description</th><th style={{ width: 90 }}>Required</th><th style={{ width: 50 }}/></tr></thead>
          <tbody>
            {(fields || []).map((f, i) => (
              <tr key={f.id}>
                <td><div className="row gap-1">
                  <button className="icon-btn" style={{ width: 22, height: 22 }} disabled={i === 0} onClick={() => move(i, -1)}><Icon name="arrowUp" size={12}/></button>
                  <button className="icon-btn" style={{ width: 22, height: 22 }} disabled={i === fields.length - 1} onClick={() => move(i, 1)}><Icon name="arrowDown" size={12}/></button>
                </div></td>
                <td className="bold">{f.name}{f.type === "select" && f.options.length > 0 && <span className="text-xs muted" style={{ fontWeight: 400, marginLeft: 6 }}>{f.options.join(", ")}</span>}</td>
                <td><Badge tone="muted">{f.type}</Badge></td>
                <td className="text-sm secondary">{f.description || "—"}</td>
                <td><Switch on={f.required} onChange={() => toggleReq(f)}/></td>
                <td><button className="icon-btn" onClick={() => setConfirm(f)} title="Delete"><Icon name="trash" size={14}/></button></td>
              </tr>
            ))}
          </tbody>
        </table>
        {(fields || []).length === 0 && <Empty icon="list" title="No custom fields" hint="Add a field to capture more data."/>}
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title="Add custom field"
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" disabled={!form.name.trim()} onClick={save}>Add field</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Field name</label><input className="input" autoFocus value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. Story owner"/></div>
          <div><label className="label">Type</label>
            <PillSelect value={form.type} onChange={(v) => setForm({ ...form, type: v })} options={FIELD_TYPES.map((t) => ({ id: t, label: t }))}/>
          </div>
          {form.type === "select" && (
            <div>
              <label className="label">Options</label>
              <div className="row gap-2" style={{ marginBottom: 8 }}>
                <input className="input" value={optInput} onChange={(e) => setOptInput(e.target.value)} placeholder="Add an option" onKeyDown={(e) => { if (e.key === "Enter" && optInput.trim()) { setForm({ ...form, options: [...form.options, optInput.trim()] }); setOptInput(""); } }}/>
                <Button data-size="sm" onClick={() => { if (optInput.trim()) { setForm({ ...form, options: [...form.options, optInput.trim()] }); setOptInput(""); } }}>Add</Button>
              </div>
              <div className="row gap-2" style={{ flexWrap: "wrap" }}>
                {form.options.map((o, i) => <span key={i} className="tag" style={{ gap: 4 }}>{o}<button className="icon-btn" style={{ width: 16, height: 16 }} onClick={() => setForm({ ...form, options: form.options.filter((_, j) => j !== i) })}><Icon name="x" size={10}/></button></span>)}
              </div>
            </div>
          )}
          <div><label className="label">Description</label><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} style={{ minHeight: 50 }}/></div>
          <label className="row gap-2"><Switch on={form.required} onChange={(v) => setForm({ ...form, required: v })}/><span className="text-sm">Required field</span></label>
        </div>
      </Modal>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} body={confirm ? `Delete field "${confirm.name}"?` : ""}/>
    </div>
  );
}

// ─── Feature 10: Labels ───────────────────────────────────────────────
export function LabelsTab() {
  const { activeProjectId } = useApp();
  const { data: labels, loading, error, reload, setData } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/labels" : null);
  const [adding, setAdding] = useState(false);
  const [name, setName] = useState(""); const [color, setColor] = useState(PRESET_COLORS[0]);
  const [editing, setEditing] = useState(null);
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  async function create() {
    if (!name.trim()) return;
    try { await api("/projects/" + activeProjectId + "/labels", { method: "POST", body: { name, color } }); toast("Label created"); setAdding(false); setName(""); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function saveEdit(l) {
    setEditing(null);
    try { await api("/labels/" + l.id, { method: "PUT", body: { name: l.name, color: l.color } }); toast("Label updated"); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }
  async function del(id) { try { await api("/labels/" + id, { method: "DELETE" }); toast("Label deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  if (loading) return <Loading label="Loading labels…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Tag issues with reusable colored labels.</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setAdding((a) => !a)}>Add label</Button>
      </div>

      {adding && (
        <div className="card card-pad" style={{ marginBottom: 12 }}>
          <div className="row gap-3" style={{ marginBottom: 10 }}>
            <input className="input" autoFocus value={name} onChange={(e) => setName(e.target.value)} placeholder="Label name" style={{ maxWidth: 240 }} onKeyDown={(e) => { if (e.key === "Enter") create(); }}/>
            <span className="badge" style={{ background: color + "22", color: color }}><span className="dot" style={{ background: color }}/>{name || "preview"}</span>
            <div style={{ flex: 1 }}/>
            <Button data-size="sm" onClick={() => setAdding(false)}>Cancel</Button>
            <Button data-size="sm" variant="primary" disabled={!name.trim()} onClick={create}>Create</Button>
          </div>
          <ColorPicker value={color} onChange={setColor}/>
        </div>
      )}

      <div className="card card-pad">
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(220px, 1fr))", gap: 10 }}>
          {(labels || []).map((l) => (
            <div key={l.id} className="row gap-2" style={{ padding: "8px 10px", border: "1px solid var(--border)", borderRadius: 8 }}>
              {editing === l.id ? (
                <>
                  <input className="input" value={l.name} onChange={(e) => setData(labels.map((x) => x.id === l.id ? { ...x, name: e.target.value } : x))} style={{ padding: "2px 6px" }} autoFocus onKeyDown={(e) => { if (e.key === "Enter") saveEdit(l); }}/>
                  <button className="icon-btn" style={{ width: 22, height: 22 }} onClick={() => saveEdit(l)} title="Save"><Icon name="check" size={13} color="var(--success)"/></button>
                </>
              ) : (
                <>
                  <span className="badge" style={{ background: l.color + "22", color: l.color, flex: 1 }}><span className="dot" style={{ background: l.color }}/>{l.name}</span>
                  <span className="text-xs muted">{l.issue_count}</span>
                  <button className="icon-btn" style={{ width: 22, height: 22 }} onClick={() => setEditing(l.id)} title="Edit"><Icon name="pencil" size={12}/></button>
                  <button className="icon-btn" style={{ width: 22, height: 22 }} onClick={() => setConfirm(l)} title="Delete"><Icon name="trash" size={12}/></button>
                </>
              )}
            </div>
          ))}
        </div>
        {(labels || []).length === 0 && <Empty icon="tag" title="No labels" hint="Create your first label."/>}
      </div>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} body={confirm ? `Delete label "${confirm.name}"? It will be removed from all issues.` : ""}/>
    </div>
  );
}

// ─── Feature 11: Webhooks ─────────────────────────────────────────────
const WEBHOOK_EVENTS = ["issue.created","issue.updated","issue.transition","issue.assigned","sprint.started","sprint.completed","page.created","page.updated"];
export function WebhooksTab() {
  const { activeProjectId } = useApp();
  const { data: hooks, loading, error, reload, setData } = useApi(activeProjectId ? "/projects/" + activeProjectId + "/webhooks" : null);
  const [modal, setModal] = useState(false);
  const [form, setForm] = useState({ url: "", secret: "", events: [], is_active: true });
  const [expanded, setExpanded] = useState(null);
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  async function create() {
    if (!form.url.trim()) return;
    try { await api("/webhooks", { method: "POST", body: form }); toast("Webhook created"); setModal(false); setForm({ url: "", secret: "", events: [], is_active: true }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function toggleActive(w) {
    setData((hooks || []).map((x) => x.id === w.id ? { ...x, is_active: !x.is_active } : x));
    try { await api("/webhooks/" + w.id, { method: "PUT", body: { is_active: !w.is_active } }); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); reload(); }
  }
  async function del(id) { try { await api("/webhooks/" + id, { method: "DELETE" }); toast("Webhook deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  function toggleEvent(ev) { setForm((f) => ({ ...f, events: f.events.includes(ev) ? f.events.filter((x) => x !== ev) : [...f.events, ev] })); }

  if (loading) return <Loading label="Loading webhooks…"/>;
  if (error) return <ErrorNote msg={error}/>;

  return (
    <div>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <div className="text-sm secondary">Send event payloads to external URLs (HMAC-signed if a secret is set).</div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setModal(true)}>Add webhook</Button>
      </div>
      <div className="stack gap-2">
        {(hooks || []).map((w) => (
          <div key={w.id} className="card" style={{ overflow: "hidden" }}>
            <div className="row gap-3" style={{ padding: "12px 16px" }}>
              <span style={{ width: 8, height: 8, borderRadius: "50%", background: w.last_status === "failed" ? "var(--danger)" : w.last_status === "success" ? "var(--success)" : "var(--text-muted)" }} title={w.last_status || "no deliveries"}/>
              <span className="mono text-sm" style={{ flex: 1, minWidth: 0, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }} title={w.url}>{w.url}</span>
              <div className="row gap-1" style={{ flexWrap: "wrap", maxWidth: 280, justifyContent: "flex-end" }}>
                {w.events.slice(0, 3).map((e) => <Badge key={e} tone="muted">{e}</Badge>)}
                {w.events.length > 3 && <Badge tone="muted">+{w.events.length - 3}</Badge>}
              </div>
              <Switch on={w.is_active} onChange={() => toggleActive(w)}/>
              <button className="btn btn-ghost" data-size="sm" onClick={() => setExpanded((x) => x === w.id ? null : w.id)}><Icon name="history" size={14}/> Deliveries</button>
              <button className="icon-btn" onClick={() => setConfirm(w)} title="Delete"><Icon name="trash" size={15}/></button>
            </div>
            {expanded === w.id && <WebhookDeliveries id={w.id}/>}
          </div>
        ))}
        {(hooks || []).length === 0 && <Empty icon="code" title="No webhooks" hint="Add a webhook to stream events."/>}
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title="Add webhook"
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" disabled={!form.url.trim()} onClick={create}>Create webhook</Button></>}>
        <div className="stack gap-3">
          <div><label className="label">Payload URL</label><input className="input mono" autoFocus value={form.url} onChange={(e) => setForm({ ...form, url: e.target.value })} placeholder="https://example.com/webhook"/></div>
          <div><label className="label">Secret <span className="muted" style={{ fontWeight: 400 }}>(optional)</span></label><input className="input mono" value={form.secret} onChange={(e) => setForm({ ...form, secret: e.target.value })} placeholder="whsec_…"/><div className="help">Used to sign payloads with an HMAC SHA-256 signature.</div></div>
          <div>
            <label className="label">Events</label>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 6 }}>
              {WEBHOOK_EVENTS.map((ev) => (
                <label key={ev} className="row gap-2" style={{ padding: "5px 8px", border: "1px solid " + (form.events.includes(ev) ? "var(--indigo-500)" : "var(--border)"), borderRadius: 6, background: form.events.includes(ev) ? "var(--indigo-50)" : "var(--bg)", cursor: "default" }}>
                  <span style={{ width: 16, height: 16, borderRadius: 4, border: "1.5px solid " + (form.events.includes(ev) ? "var(--indigo-600)" : "var(--border-strong)"), background: form.events.includes(ev) ? "var(--indigo-600)" : "transparent", display: "grid", placeItems: "center", color: "#fff" }} onClick={() => toggleEvent(ev)}>{form.events.includes(ev) && <Icon name="check" size={11} strokeWidth={3}/>}</span>
                  <span className="mono text-xs" onClick={() => toggleEvent(ev)}>{ev}</span>
                </label>
              ))}
            </div>
          </div>
          <label className="row gap-2"><Switch on={form.is_active} onChange={(v) => setForm({ ...form, is_active: v })}/><span className="text-sm">Active</span></label>
        </div>
      </Modal>
      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)} body={confirm ? "Delete this webhook? Event delivery will stop immediately." : ""}/>
    </div>
  );
}

function WebhookDeliveries({ id }) {
  const { data: deliveries, loading } = useApi("/webhooks/" + id + "/deliveries");
  if (loading) return <div style={{ borderTop: "1px solid var(--border)" }}><Loading/></div>;
  return (
    <div style={{ borderTop: "1px solid var(--border)", background: "var(--bg-subtle)" }}>
      <table className="table">
        <thead><tr><th>When</th><th>Event</th><th style={{ width: 110 }}>Status</th><th style={{ width: 90 }}>Code</th><th style={{ width: 90 }}>Duration</th></tr></thead>
        <tbody>
          {(deliveries || []).map((d) => (
            <tr key={d.id}>
              <td className="text-xs muted">{fmtDate(d.at)}</td>
              <td className="mono text-xs">{d.event}</td>
              <td><Badge tone={d.status === "success" ? "success" : "danger"} dot>{d.status}</Badge></td>
              <td className="mono text-xs">{d.code}</td>
              <td className="text-xs muted">{d.duration_ms}ms</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// ─── Feature 25: Sprint capacity ──────────────────────────────────────
export function SprintCapacity({ sprintId }) {
  const { people } = useApp();
  const { data, loading, error, reload, setData } = useApi(sprintId ? "/sprints/" + sprintId + "/capacity" : null);
  const [modal, setModal] = useState(false);
  const [draft, setDraft] = useState([]);
  const toast = useToast();

  function openModal() { setDraft((data.members || []).map((m) => ({ user_id: m.user_id, available_hours: m.available_hours }))); setModal(true); }
  async function save() {
    try { const r = await api("/sprints/" + sprintId + "/capacity", { method: "PUT", body: { members: draft } }); setData(r); toast("Capacity updated"); setModal(false); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  if (loading) return <div className="card" style={{ marginBottom: 16 }}><div className="card-head"><h3>Capacity</h3></div><Loading/></div>;
  if (error) return <ErrorNote msg={error}/>;
  if (!data) return null;

  const members = data.members || [];
  const totals = members.reduce((a, m) => ({ avail: a.avail + m.available_hours, logged: a.logged + m.logged_hours }), { avail: 0, logged: 0 });

  return (
    <div className="card" style={{ marginBottom: 16 }}>
      <div className="card-head">
        <h3>Capacity</h3>
        <Button data-size="sm" icon="clock" onClick={openModal}>Set capacity</Button>
      </div>
      <div style={{ padding: 16 }}>
        <table className="table" style={{ border: "1px solid var(--border)", borderRadius: 8, overflow: "hidden" }}>
          <thead><tr><th>Member</th><th style={{ width: 120 }}>Available</th><th style={{ width: 120 }}>Logged</th><th style={{ width: 120 }}>Remaining</th><th style={{ width: 220 }}>Utilization</th></tr></thead>
          <tbody>
            {members.map((m) => {
              const u = m.user ? adaptUser(m.user) : people.find((p) => p.id === m.user_id) || { name: m.user_id, initials: "?", color: "#94A3B8" };
              const util = m.available_hours ? Math.round((m.logged_hours / m.available_hours) * 100) : 0;
              const col = util > 100 ? "var(--danger)" : util >= 80 ? "var(--warning)" : "var(--success)";
              return (
                <tr key={m.user_id}>
                  <td><span className="row gap-2"><Avatar user={u} size="sm"/>{u.name}</span></td>
                  <td className="mono text-sm">{m.available_hours}h</td>
                  <td className="mono text-sm">{m.logged_hours}h</td>
                  <td className="mono text-sm" style={{ color: m.available_hours - m.logged_hours < 0 ? "var(--danger)" : undefined }}>{m.available_hours - m.logged_hours}h</td>
                  <td>
                    <div className="row gap-2">
                      <div style={{ flex: 1, height: 8, background: "var(--bg-muted)", borderRadius: 99, overflow: "hidden" }}>
                        <div style={{ width: Math.min(100, util) + "%", height: "100%", background: col }}/>
                      </div>
                      <span className="text-xs bold" style={{ width: 38, textAlign: "right", color: col }}>{util}%</span>
                    </div>
                  </td>
                </tr>
              );
            })}
            <tr style={{ background: "var(--bg-subtle)" }}>
              <td className="bold">Total</td>
              <td className="mono bold">{totals.avail}h</td>
              <td className="mono bold">{totals.logged}h</td>
              <td className="mono bold">{totals.avail - totals.logged}h</td>
              <td><span className="text-xs bold">{totals.avail ? Math.round((totals.logged / totals.avail) * 100) : 0}% team utilization</span></td>
            </tr>
          </tbody>
        </table>
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title="Set sprint capacity"
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" onClick={save}>Save capacity</Button></>}>
        <div className="stack gap-2">
          {draft.map((d, i) => {
            const u = people.find((p) => p.id === d.user_id) || { name: d.user_id, initials: "?", color: "#94A3B8" };
            return (
              <div key={d.user_id} className="row gap-3" style={{ padding: "6px 0" }}>
                <Avatar user={u} size="sm"/>
                <span className="text-sm" style={{ flex: 1 }}>{u.name}</span>
                <input className="input" type="number" min="0" value={d.available_hours} onChange={(e) => setDraft(draft.map((x, j) => j === i ? { ...x, available_hours: Number(e.target.value) } : x))} style={{ width: 110 }}/>
                <span className="text-xs muted">hours</span>
              </div>
            );
          })}
        </div>
      </Modal>
    </div>
  );
}

// ─── Automation rules tab ─────────────────────────────
const TRIGGER_LABELS = {
  "issue.created":          "Issue created",
  "issue.status_changed":   "Status changed",
  "issue.assigned":         "Issue assigned",
  "sprint.started":         "Sprint started",
  "sprint.completed":       "Sprint completed",
  "comment.added":          "Comment added",
};
const ACTION_LABELS = {
  "assign":                 "Assign to user",
  "notify":                 "Send notification",
  "set_label":              "Add label",
  "set_priority":           "Set priority",
  "move_to_status":         "Move to status",
  "create_issue":           "Create linked issue",
};

export function AutomationTab() {
  const { activeProjectId } = useApp();
  const { data: rules, loading, error, reload } = useApi(
    activeProjectId ? "/projects/" + activeProjectId + "/automation-rules" : "/automation-rules",
    [activeProjectId]
  );
  const [modal, setModal] = useState(false);
  const [form, setForm] = useState({ name: "", trigger: "issue.created", action: "notify", condition: "", enabled: true });
  const [saving, setSaving] = useState(false);
  const [confirm, setConfirm] = useState(null);
  const toast = useToast();

  async function create() {
    if (!form.name.trim()) return;
    setSaving(true);
    try {
      await api(activeProjectId ? "/projects/" + activeProjectId + "/automation-rules" : "/automation-rules", {
        method: "POST", body: form,
      });
      toast("Automation rule created");
      setModal(false);
      setForm({ name: "", trigger: "issue.created", action: "notify", condition: "", enabled: true });
      reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSaving(false); }
  }

  async function toggle(rule) {
    try {
      await api("/automation-rules/" + rule.id, { method: "PUT", body: { enabled: !rule.enabled } });
      reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function del(id) {
    try {
      await api("/automation-rules/" + id, { method: "DELETE" });
      toast("Rule deleted"); reload();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setConfirm(null); }
  }

  if (loading) return <Loading label="Loading automation rules…"/>;

  const list = rules?.items || rules || [];

  return (
    <div className="stack gap-3">
      <div className="row" style={{ justifyContent: "space-between" }}>
        <div>
          <div className="bold">Automation rules</div>
          <div className="text-sm secondary">Automate repetitive tasks with triggers and actions.</div>
        </div>
        <Button variant="primary" icon="plus" data-size="sm" onClick={() => setModal(true)}>New rule</Button>
      </div>

      {list.length === 0 && !error && (
        <Empty icon="bolt" title="No automation rules" hint="Create a rule to automate workflows — e.g. assign new issues, send notifications, or move status."/>
      )}
      {error && <ErrorNote msg={error}/>}

      <div className="stack gap-2">
        {list.map((rule) => (
          <div key={rule.id} className="card" style={{ padding: "12px 16px" }}>
            <div className="row gap-3">
              <div style={{ flex: 1, minWidth: 0 }}>
                <div className="row gap-2" style={{ marginBottom: 4 }}>
                  <span className="bold text-sm">{rule.name}</span>
                  {!rule.enabled && <Badge tone="muted">Disabled</Badge>}
                </div>
                <div className="text-xs secondary row gap-2">
                  <span className="badge" data-tone="info" style={{ padding: "1px 7px" }}>{TRIGGER_LABELS[rule.trigger] || rule.trigger}</span>
                  <Icon name="arrowRight" size={11} color="var(--text-muted)"/>
                  <span className="badge" data-tone="success" style={{ padding: "1px 7px" }}>{ACTION_LABELS[rule.action] || rule.action}</span>
                  {rule.last_run_at && <span className="muted">· last ran {fmtDate(rule.last_run_at)}</span>}
                </div>
              </div>
              <div className="row gap-2">
                <Switch on={rule.enabled} onChange={() => toggle(rule)}/>
                <button className="icon-btn" title="Delete rule" onClick={() => setConfirm(rule)}>
                  <Icon name="trash" size={14}/>
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      <Modal open={modal} onClose={() => setModal(false)} title="Create automation rule"
        footer={<><Button onClick={() => setModal(false)}>Cancel</Button><Button variant="primary" disabled={!form.name.trim() || saving} onClick={create}>{saving ? "Creating…" : "Create rule"}</Button></>}>
        <div className="stack gap-3">
          <div>
            <label className="label">Rule name</label>
            <input className="input" autoFocus value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. Auto-assign bugs to QA"/>
          </div>
          <div>
            <label className="label">When (trigger)</label>
            <select className="select" value={form.trigger} onChange={(e) => setForm({ ...form, trigger: e.target.value })}>
              {Object.entries(TRIGGER_LABELS).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
            </select>
          </div>
          <div>
            <label className="label">Then (action)</label>
            <select className="select" value={form.action} onChange={(e) => setForm({ ...form, action: e.target.value })}>
              {Object.entries(ACTION_LABELS).map(([k, v]) => <option key={k} value={k}>{v}</option>)}
            </select>
          </div>
          <div>
            <label className="label">Condition <span style={{ fontWeight: 400, color: "var(--text-muted)" }}>(optional)</span></label>
            <input className="input" value={form.condition} onChange={(e) => setForm({ ...form, condition: e.target.value })} placeholder="e.g. priority = Critical"/>
          </div>
        </div>
      </Modal>

      <ConfirmDelete open={!!confirm} onClose={() => setConfirm(null)} onConfirm={() => del(confirm.id)}
        title="Delete rule" body={confirm ? `Delete "${confirm.name}"? This cannot be undone.` : ""}/>
    </div>
  );
}
