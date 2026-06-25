import { useState, useEffect, useRef } from 'react';
import DOMPurify from 'dompurify';
import { Icon } from '../components/icons';
import { Avatar, Badge, Button, Modal, Switch, Empty, useToast } from '../components/components';
import { api, apiDownload, useApi } from '../api/api';
import { fmtDate } from '../api/adapters';
import { MiniSpinner } from '../panels/issue';
import { ConfirmDelete } from '../panels/settings';
import { RichEditor } from '../components/RichEditor';
import { ws } from '../lib/ws';
import { useApp } from '../store/AppContext';

const SPACE_TYPE_META = { team: { icon: "users", label: "Team" }, personal: { icon: "user", label: "Personal" }, project: { icon: "briefcase", label: "Project" } };
const REACTION_EMOJI = ["👍", "❤️", "🎉", "🚀", "👀", "😄"];

export function WikiView({ nav, target, onTargetConsumed }) {
  const { me } = useApp();
  const isAdmin = me?.role === "admin";
  const { data: spaces, reload: reloadSpaces } = useApi("/spaces");
  const [spaceId, setSpaceId] = useState(null);
  const [mode, setMode] = useState("page");
  const [pageId, setPageId] = useState(null);
  const [blogId, setBlogId] = useState(null);
  const [newSpaceOpen, setNewSpaceOpen] = useState(false);
  const [spaceSettings, setSpaceSettings] = useState(false);
  const [spaceMenu, setSpaceMenu] = useState(null); // { id, x, y }

  useEffect(() => {
    if (!spaceMenu) return;
    const close = () => setSpaceMenu(null);
    document.addEventListener("click", close);
    return () => document.removeEventListener("click", close);
  }, [spaceMenu]);

  useEffect(() => { if (spaces && spaces.length && !spaceId) setSpaceId(spaces[0].id); }, [spaces]);

  // Navigate to target page from global search
  useEffect(() => {
    if (!target || !spaces) return;
    if (target.spaceId) setSpaceId(target.spaceId);
    if (target.pageId) { setPageId(target.pageId); setMode("page"); }
    onTargetConsumed?.();
  }, [target, spaces]);

  const space = (spaces || []).find((s) => s.id === spaceId);

  return (
    <div className="wiki" style={{ gridTemplateColumns: "264px 1fr" }}>
      <aside className="wiki-tree" style={{ display: "flex", flexDirection: "column", padding: 0 }}>
        <div style={{ padding: "12px 10px 8px", borderBottom: "1px solid var(--border)", display: "flex", flexDirection: "column", minHeight: 0 }}>
          <div className="row" style={{ justifyContent: "space-between", marginBottom: 8 }}>
            <h4 style={{ margin: 0 }}>Spaces</h4>
            {isAdmin && <button className="icon-btn" style={{ width: 22, height: 22 }} title="New space" onClick={() => setNewSpaceOpen(true)}><Icon name="plus" size={13}/></button>}
          </div>
          <div className="stack gap-1" style={{ overflowY: "auto", maxHeight: 260 }}>
            {(spaces || []).map((s) => (
              <div key={s.id} className="row" style={{ position: "relative" }}>
                <button className="tree-item" aria-current={s.id === spaceId ? "page" : undefined} style={{ flex: 1, minWidth: 0 }} onClick={() => { setSpaceId(s.id); setPageId(null); setBlogId(null); setMode("page"); }}>
                  <span style={{ width: 18, height: 18, borderRadius: 5, background: "var(--indigo-50)", color: "var(--indigo-600)", display: "grid", placeItems: "center", flexShrink: 0 }}><Icon name={SPACE_TYPE_META[s.type]?.icon || "notes"} size={11}/></span>
                  <span style={{ flex: 1, textAlign: "left", overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{s.name}</span>
                  <span className="text-xs muted">{s.page_count}</span>
                </button>
                {isAdmin && (
                  <button className="icon-btn" title="Space options" style={{ width: 20, height: 20, flexShrink: 0 }}
                    onClick={(e) => { e.stopPropagation(); setSpaceId(s.id); setSpaceMenu({ id: s.id, x: e.clientX, y: e.clientY }); }}>
                    <Icon name="moreH" size={12}/>
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        {spaceMenu && (
          <div onClick={(e) => e.stopPropagation()} style={{ position: "fixed", top: spaceMenu.y, left: spaceMenu.x, zIndex: 300, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, minWidth: 160 }}>
            <button className="nav-item" style={{ color: "var(--text)", fontSize: 13 }} onClick={() => { setSpaceSettings(true); setSpaceMenu(null); }}>
              <Icon name="settings" size={13}/> Space settings
            </button>
            <button className="nav-item" style={{ color: "var(--danger)", fontSize: 13 }} onClick={async () => {
              setSpaceMenu(null);
              if (!window.confirm("Delete this space and all its pages?")) return;
              try {
                await api("/spaces/" + spaceMenu.id, { method: "DELETE" });
                reloadSpaces(); setSpaceId(null);
              } catch (e) { alert(e.message); }
            }}>
              <Icon name="trash" size={13}/> Delete space
            </button>
          </div>
        )}

        {space && <SpaceNav key={space.id} space={space} mode={mode} setMode={setMode} pageId={pageId} setPageId={setPageId} blogId={blogId} setBlogId={setBlogId} onSettings={() => setSpaceSettings(true)} isAdmin={isAdmin}/>}
      </aside>

      {space ? (
        mode === "blog"
          ? <BlogDoc key={blogId || "blog-empty"} space={space} blogId={blogId} setBlogId={setBlogId}/>
          : <PageDoc key={pageId || "page-empty"} space={space} pageId={pageId} setPageId={setPageId} nav={nav}/>
      ) : (
        <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No space selected" hint={isAdmin ? "Create a space to get started." : "Ask an admin to create a space."}/></div>
      )}

      {isAdmin && <NewSpaceModal open={newSpaceOpen} onClose={() => setNewSpaceOpen(false)} onCreated={(s) => { reloadSpaces(); setSpaceId(s.id); setPageId(null); }}/>}
      {space && isAdmin && <SpaceSettingsModal open={spaceSettings} onClose={() => setSpaceSettings(false)} space={space} onChanged={reloadSpaces} onDeleted={() => { reloadSpaces(); setSpaceId(null); }}/>}
    </div>
  );
}

// ─── Page template picker ──────────────────────────────
const PAGE_TEMPLATES = [
  { id: "blank", title: "Blank page", icon: "📄", content: null },
  {
    id: "meeting",
    title: "Meeting notes",
    icon: "📝",
    content: {
      type: "doc",
      content: [
        { type: "heading", attrs: { level: 1 }, content: [{ type: "text", text: "Meeting Notes" }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Attendees" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Agenda" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Action items" }] },
        { type: "taskList", content: [{ type: "taskItem", attrs: { checked: false }, content: [{ type: "paragraph", content: [] }] }] },
      ],
    },
  },
  {
    id: "spec",
    title: "Technical spec",
    icon: "⚙️",
    content: {
      type: "doc",
      content: [
        { type: "heading", attrs: { level: 1 }, content: [{ type: "text", text: "Technical Specification" }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Overview" }] },
        { type: "paragraph", content: [{ type: "text", text: "Describe the problem and solution." }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Goals" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Design" }] },
        { type: "paragraph", content: [{ type: "text", text: "Architecture details here." }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Open questions" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
      ],
    },
  },
  {
    id: "retrospective",
    title: "Retrospective",
    icon: "🔄",
    content: {
      type: "doc",
      content: [
        { type: "heading", attrs: { level: 1 }, content: [{ type: "text", text: "Sprint Retrospective" }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "✅ What went well" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "❌ What didn't go well" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "💡 Action items" }] },
        { type: "taskList", content: [{ type: "taskItem", attrs: { checked: false }, content: [{ type: "paragraph", content: [] }] }] },
      ],
    },
  },
  {
    id: "how-to",
    title: "How-to guide",
    icon: "📖",
    content: {
      type: "doc",
      content: [
        { type: "heading", attrs: { level: 1 }, content: [{ type: "text", text: "How to …" }] },
        { type: "paragraph", content: [{ type: "text", text: "Brief description of what this covers." }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Prerequisites" }] },
        { type: "bulletList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Steps" }] },
        { type: "orderedList", content: [{ type: "listItem", content: [{ type: "paragraph", content: [] }] }] },
        { type: "heading", attrs: { level: 2 }, content: [{ type: "text", text: "Troubleshooting" }] },
        { type: "paragraph", content: [{ type: "text", text: "Common issues and solutions." }] },
      ],
    },
  },
];

function NewPageModal({ open, onClose, onCreate }) {
  const [title, setTitle] = useState("Untitled");
  const [templateId, setTemplateId] = useState("blank");

  useEffect(() => { if (open) { setTitle("Untitled"); setTemplateId("blank"); } }, [open]);

  function create() {
    const tpl = PAGE_TEMPLATES.find((t) => t.id === templateId);
    onCreate(title || "Untitled", tpl?.content || null);
  }

  return (
    <Modal open={open} onClose={onClose} title="New page"
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="primary" onClick={create}>Create page</Button></>}>
      <div className="stack gap-3">
        <div>
          <label className="label">Page title</label>
          <input className="input" autoFocus value={title} onChange={(e) => setTitle(e.target.value)}
            onKeyDown={(e) => { if (e.key === "Enter") create(); }}/>
        </div>
        <div>
          <label className="label">Start from template</label>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 8 }}>
            {PAGE_TEMPLATES.map((t) => (
              <button key={t.id} onClick={() => setTemplateId(t.id)}
                style={{ padding: "10px 12px", borderRadius: 8, border: "1.5px solid " + (templateId === t.id ? "var(--indigo-500)" : "var(--border)"), background: templateId === t.id ? "var(--indigo-50)" : "var(--bg)", textAlign: "left", cursor: "default", display: "flex", alignItems: "center", gap: 8 }}>
                <span style={{ fontSize: 18 }}>{t.icon}</span>
                <span className="text-sm" style={{ fontWeight: templateId === t.id ? 600 : 400, color: templateId === t.id ? "var(--indigo-700)" : "var(--text)" }}>{t.title}</span>
              </button>
            ))}
          </div>
        </div>
      </div>
    </Modal>
  );
}

// ─── Space nav: page tree (F17) + blog list (F22) ─────────────────────
function SpaceNav({ space, mode, setMode, pageId, setPageId, blogId, setBlogId, onSettings, isAdmin }) {
  const { data: tree, loading, reload } = useApi("/spaces/" + space.id + "/pages/tree", [space.id]);
  const { data: blog, reload: reloadBlog } = useApi("/spaces/" + space.id + "/blog-posts", [space.id]);
  const [expanded, setExpanded] = useState({});
  const [menu, setMenu] = useState(null);
  const [moveNode, setMoveNode] = useState(null);
  const [renaming, setRenaming] = useState(null);
  const [newPageOpen, setNewPageOpen] = useState(false);
  const [newPageParent, setNewPageParent] = useState(null);
  const toast = useToast();

  useEffect(() => { if (tree && tree.length && !pageId && mode === "page") { setPageId(tree[0].id); setExpanded((e) => ({ ...e, [tree[0].id]: true })); } }, [tree]);
  useEffect(() => {
    const close = () => setMenu(null);
    if (menu) { document.addEventListener("click", close); return () => document.removeEventListener("click", close); }
  }, [menu]);

  async function newPage(parentId, title, content) {
    try {
      const body = { title: title || "Untitled", parent_id: parentId || null };
      if (content) body.content = content;
      const p = await api("/spaces/" + space.id + "/pages", { method: "POST", body });
      toast("Page created"); reload(); setMode("page"); setPageId(p.id);
      if (parentId) setExpanded((e) => ({ ...e, [parentId]: true }));
    }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  function openNewPage(parentId) { setNewPageParent(parentId || null); setNewPageOpen(true); }
  async function rename(node, title) {
    setRenaming(null);
    try { await api("/pages/" + node.id, { method: "PUT", body: { title } }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function del(node) {
    try { await api("/pages/" + node.id, { method: "DELETE" }); toast("Page deleted"); if (pageId === node.id) setPageId(null); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function move(node, parentId) {
    setMoveNode(null);
    try { await api("/pages/" + node.id + "/move", { method: "PUT", body: { parent_id: parentId || null } }); toast("Page moved"); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  function flatten(nodes, acc) { (nodes || []).forEach((n) => { acc.push(n); flatten(n.children, acc); }); return acc; }
  const allNodes = flatten(tree, []);

  function TreeNode({ node, depth }) {
    const hasKids = node.children && node.children.length > 0;
    const open = expanded[node.id];
    return (
      <div>
        <div className="tree-item" aria-current={mode === "page" && pageId === node.id ? "page" : undefined} style={{ paddingLeft: 8 + depth * 14 }}
          onClick={() => { setMode("page"); setPageId(node.id); }}
          onContextMenu={(e) => { e.preventDefault(); setMenu({ x: e.clientX, y: e.clientY, node }); }}>
          <span onClick={(e) => { if (hasKids) { e.stopPropagation(); setExpanded((x) => ({ ...x, [node.id]: !x[node.id] })); } }} style={{ display: "grid", placeItems: "center", width: 14 }}>
            {hasKids ? <Icon name={open ? "chevronDown" : "chevronRight"} size={11} color="var(--text-muted)"/> : <Icon name="notes" size={12} color="var(--text-muted)"/>}
          </span>
          {renaming === node.id ? (
            <input className="input" autoFocus defaultValue={node.title} style={{ padding: "1px 5px", height: 22 }} onClick={(e) => e.stopPropagation()} onBlur={(e) => rename(node, e.target.value)} onKeyDown={(e) => { if (e.key === "Enter") rename(node, e.target.value); }}/>
          ) : <span style={{ flex: 1 }}>{node.icon && <span style={{ marginRight: 3 }}>{node.icon}</span>}{node.title}</span>}
          <button className="icon-btn" style={{ width: 18, height: 18 }} title="New child page" onClick={(e) => { e.stopPropagation(); openNewPage(node.id); }}><Icon name="plus" size={11}/></button>
        </div>
        {hasKids && open && node.children.map((c) => <TreeNode key={c.id} node={c} depth={depth + 1}/>)}
      </div>
    );
  }

  return (
    <div style={{ flex: 1, overflowY: "auto", padding: "10px 6px" }}>
      <div className="row" style={{ justifyContent: "space-between", padding: "0 6px 6px" }}>
        <div className="row gap-2" style={{ minWidth: 0 }}>
          <h4 style={{ margin: 0 }}>{space.name}</h4>
        </div>
        <div className="row gap-1">
          {isAdmin && <button className="icon-btn" style={{ width: 22, height: 22 }} title="Space settings" onClick={onSettings}><Icon name="settings" size={13}/></button>}
          <button className="icon-btn" style={{ width: 22, height: 22 }} title="New page" onClick={() => openNewPage(null)}><Icon name="plus" size={13}/></button>
        </div>
      </div>

      {loading ? <div className="row gap-2 text-xs muted" style={{ padding: 8 }}><MiniSpinner/> Loading…</div> : (tree || []).map((n) => <TreeNode key={n.id} node={n} depth={0}/>)}
      {!loading && (tree || []).length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>No pages yet.</div>}

      <div style={{ marginTop: 16, borderTop: "1px solid var(--border)", paddingTop: 10 }}>
        <div className="row" style={{ justifyContent: "space-between", padding: "0 6px 6px" }}>
          <h4 style={{ margin: 0 }}>Blog</h4>
          <button className="icon-btn" style={{ width: 22, height: 22 }} title="New post" onClick={async () => { try { const p = await api("/spaces/" + space.id + "/blog-posts", { method: "POST", body: { title: "Untitled post" } }); reloadBlog(); setMode("blog"); setBlogId(p.id); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }}><Icon name="plus" size={13}/></button>
        </div>
        <div className="stack gap-1">
          {(blog?.items || []).map((p) => (
            <button key={p.id} className="tree-item" aria-current={mode === "blog" && blogId === p.id ? "page" : undefined} style={{ width: "100%", alignItems: "flex-start" }} onClick={() => { setMode("blog"); setBlogId(p.id); }}>
              <Icon name="notes" size={12} color="var(--text-muted)" style={{ marginTop: 3 }}/>
              <span className="stack" style={{ flex: 1, textAlign: "left", lineHeight: 1.3 }}>
                <span className="row gap-2">{p.title}{!p.published && <Badge tone="warning" style={{ fontSize: 9 }}>Draft</Badge>}</span>
                <span className="text-xs muted">{(p.author || {}).name} · {p.published_at ? fmtDate(p.published_at) : "unpublished"}</span>
              </span>
            </button>
          ))}
          {(blog?.items || []).length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>No posts yet.</div>}
        </div>
      </div>

      {menu && (
        <div style={{ position: "fixed", top: menu.y, left: menu.x, zIndex: 200, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, minWidth: 170 }} onClick={(e) => e.stopPropagation()}>
          {[["plus","New child page",() => { openNewPage(menu.node.id); setMenu(null); }],["pencil","Rename",() => { setRenaming(menu.node.id); setMenu(null); }],["arrowRight","Move",() => { setMoveNode(menu.node); setMenu(null); }],["trash","Delete",() => { del(menu.node); setMenu(null); },true]].map(([ic, label, fn, danger], i) => (
            <button key={i} className="nav-item" style={{ color: danger ? "var(--danger)" : "var(--text)", fontSize: 13 }} onClick={fn}><Icon name={ic} size={13}/> {label}</button>
          ))}
        </div>
      )}

      <Modal open={!!moveNode} onClose={() => setMoveNode(null)} title="Move page"
        footer={<Button onClick={() => setMoveNode(null)}>Done</Button>}>
        {moveNode && (
          <div className="stack gap-2">
            <button className="nav-item" style={{ color: "var(--text)" }} onClick={() => move(moveNode, null)}><Icon name="notes" size={13}/> Top level</button>
            {allNodes.filter((n) => n.id !== moveNode.id).map((n) => (
              <button key={n.id} className="nav-item" style={{ color: "var(--text)" }} onClick={() => move(moveNode, n.id)}><Icon name="folder" size={13}/> {n.title}</button>
            ))}
          </div>
        )}
      </Modal>

      <NewPageModal
        open={newPageOpen}
        onClose={() => setNewPageOpen(false)}
        onCreate={(title, content) => { setNewPageOpen(false); newPage(newPageParent, title, content); }}
      />
    </div>
  );
}

// ─── Page document (F18 versions, F19 comments, F20 export, F21 reactions)
function PageDoc({ space, pageId, setPageId, nav }) {
  const { me } = useApp();
  const { data: page, loading, reload, setData } = useApi(pageId ? "/pages/" + pageId : null, [pageId]);
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(null);
  const [right, setRight] = useState(null);
  const [showToc, setShowToc] = useState(false);
  const [exportOpen, setExportOpen] = useState(false);
  const [emojiOpen, setEmojiOpen] = useState(false);
  const toast = useToast();
  const bodyRef = useRef(null);
  const exportRef = useRef(null);

  const { data: comments, reload: reloadComments } = useApi(pageId ? "/pages/" + pageId + "/inline-comments" : null, [pageId]);
  const [sel, setSel] = useState(null);
  const [composing, setComposing] = useState(false);
  const [commentText, setCommentText] = useState("");
  const [activeAnchor, setActiveAnchor] = useState(null);
  const [pageLock, setPageLock] = useState(null); // null | { user_id, user_name }
  const lockIntervalRef = useRef(null);

  useEffect(() => { if (page) setDraft(page.content || null); setEditing(false); setRight(null); setPageLock(null); }, [pageId]);

  // Subscribe to page lock WS events
  useEffect(() => {
    if (!pageId) return;
    ws.subscribe(`page:${pageId}`);
    const offLocked   = ws.on('page.locked',   (msg) => setPageLock(msg.payload));
    const offUnlocked = ws.on('page.unlocked',  ()    => setPageLock(null));
    return () => { offLocked(); offUnlocked(); ws.unsubscribe(`page:${pageId}`); };
  }, [pageId]);

  async function startEditing() {
    try {
      await api(`/pages/${pageId}/lock`, { method: "POST", body: { ttl_seconds: 300 } });
      setPageLock(null);
      setEditing(true);
      // Extend lock every 3 min while editing; store ref so cleanup always clears it.
      lockIntervalRef.current = setInterval(
        () => api(`/pages/${pageId}/lock`, { method: "PUT", body: { ttl_seconds: 300 } }).catch(() => {}),
        180000
      );
    } catch (e) {
      try { const l = await api(`/pages/${pageId}/lock`); setPageLock(l); } catch (_) {}
      toast(e.message || "Page is locked by another user", { icon: "lock", color: "#F87171" });
    }
  }

  async function releaseAndClose() {
    clearInterval(lockIntervalRef.current);
    lockIntervalRef.current = null;
    setEditing(false);
    setDraft(page.content);
    try { await api(`/pages/${pageId}/lock`, { method: "DELETE" }); } catch (_) {}
  }
  useEffect(() => {
    if (!exportOpen) return;
    const h = (e) => { if (exportRef.current && !exportRef.current.contains(e.target)) setExportOpen(false); };
    document.addEventListener("mousedown", h); return () => document.removeEventListener("mousedown", h);
  }, [exportOpen]);

  if (!pageId) return <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No page selected" hint="Pick a page from the tree."/></div>;
  if (loading || !page) return <div className="wiki-doc"><div className="skel" style={{ height: 36, width: "60%", marginBottom: 16 }}/><div className="skel" style={{ height: 14, marginBottom: 8 }}/><div className="skel" style={{ height: 14, width: "80%" }}/></div>;

  async function save() {
    try {
      await api("/pages/" + pageId, { method: "PUT", body: { content: draft } });
      clearInterval(lockIntervalRef.current);
      lockIntervalRef.current = null;
      try { await api(`/pages/${pageId}/lock`, { method: "DELETE" }); } catch (_) {}
      toast("Page saved"); setEditing(false); reload();
    }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function changeIcon(icon) {
    setEmojiOpen(false);
    try { await api("/pages/" + pageId, { method: "PUT", body: { icon } }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function exportAs(fmt) {
    setExportOpen(false);
    try { await apiDownload("/pages/" + pageId + "/export/" + fmt); toast("Exporting as " + fmt.toUpperCase() + "…"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  function onMouseUp() {
    if (editing) return;
    const s = window.getSelection();
    const text = s && s.toString().trim();
    if (text && text.length > 1 && bodyRef.current && bodyRef.current.contains(s.anchorNode)) {
      const rect = s.getRangeAt(0).getBoundingClientRect();
      setSel({ text, rect }); setComposing(false);
    } else if (!composing) { setSel(null); }
  }
  async function postComment() {
    if (!commentText.trim() || !sel) return;
    try { await api("/pages/" + pageId + "/inline-comments", { method: "POST", body: { text: commentText, anchor_text: sel.text, anchor_offset: 0 } }); toast("Comment added"); setCommentText(""); setComposing(false); setSel(null); reloadComments(); setRight("comments"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  const openComments = (comments || []).filter((c) => !c.resolved);
  const bodyHtml = highlightBody(prosemirrorToHtml(page.content), openComments, activeAnchor);

  function onBodyClick(e) {
    const mark = e.target.closest && e.target.closest("[data-anchor]");
    if (mark) { setActiveAnchor(mark.getAttribute("data-anchor")); setRight("comments"); }
  }

  const headings = extractHeadings(page.content);
  const hasToc = headings.length >= 2;

  return (
    <div style={{ display: "grid", gridTemplateColumns: right ? "1fr 340px" : showToc && hasToc ? "1fr 220px" : "1fr", overflow: "hidden" }}>
      <div className="wiki-doc" onMouseUp={onMouseUp} style={{ position: "relative" }}>
        <div className="row gap-2 text-sm muted" style={{ marginBottom: 14 }}>
          <span>{space.name}</span> <Icon name="chevronRight" size={11}/>
          <span style={{ color: "var(--text)" }}>{page.icon && <span style={{ marginRight: 4 }}>{page.icon}</span>}{page.title}</span>
          <div style={{ flex: 1 }}/>
          {editing ? (
            <>
              <button className="btn btn-primary" data-size="sm" onClick={save}><Icon name="check" size={13}/> Save</button>
              <button className="btn btn-ghost" data-size="sm" onClick={releaseAndClose}>Cancel</button>
            </>
          ) : (
            <>
              {pageLock && (
                <span style={{ fontSize: 12, color: "#B45309", background: "#FEF9C3", border: "1px solid #FDE68A", borderRadius: 6, padding: "2px 8px", display: "flex", alignItems: "center", gap: 4 }}>
                  <Icon name="lock" size={12}/> Editing: {pageLock.user_name || pageLock.user_id}
                </span>
              )}
              <button className="btn btn-ghost" data-size="sm" onClick={startEditing} disabled={!!pageLock}><Icon name="pencil" size={13}/> Edit</button>
              {hasToc && <button className="btn btn-ghost" data-size="sm" aria-pressed={showToc} style={showToc ? { background: "var(--bg-subtle)", color: "var(--text)" } : null} onClick={() => { setShowToc(!showToc); setRight(null); }} title="Table of contents">≡ TOC</button>}
              <button className="btn btn-ghost" data-size="sm" aria-pressed={right === "history"} style={right === "history" ? { background: "var(--bg-subtle)", color: "var(--text)" } : null} onClick={() => { setRight(right === "history" ? null : "history"); setShowToc(false); }}><Icon name="history" size={13}/> History</button>
              <button className="btn btn-ghost" data-size="sm" aria-pressed={right === "comments"} style={right === "comments" ? { background: "var(--bg-subtle)", color: "var(--text)" } : null} onClick={() => { setActiveAnchor(null); setRight(right === "comments" ? null : "comments"); setShowToc(false); }}>
                <Icon name="comment" size={13}/> Comments {openComments.length > 0 && <span className="badge" data-tone="info" style={{ padding: "0 5px" }}>{openComments.length}</span>}
              </button>
              <span ref={exportRef} style={{ position: "relative" }}>
                <button className="btn btn-ghost" data-size="sm" onClick={() => setExportOpen((o) => !o)}><Icon name="download" size={13}/> Export</button>
                {exportOpen && (
                  <div style={{ position: "absolute", top: "100%", right: 0, marginTop: 4, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, minWidth: 160 }}>
                    {[["pdf","Export as PDF"],["html","Export as HTML"],["md","Export as Markdown"],["docx","Export as DOCX"]].map(([f, l]) => (
                      <button key={f} className="nav-item" style={{ color: "var(--text)", fontSize: 13 }} onClick={() => exportAs(f)}><Icon name="download" size={13}/> {l}</button>
                    ))}
                  </div>
                )}
              </span>
            </>
          )}
        </div>

        {/* Page title with emoji picker */}
        <div className="row gap-3" style={{ alignItems: "flex-start", marginBottom: 8 }}>
          <div style={{ position: "relative" }}>
            <button className="icon-btn" style={{ width: 40, height: 40, fontSize: 24, borderRadius: 8, border: "1px solid var(--border)", background: page.icon ? "transparent" : "var(--bg-subtle)" }}
              title="Set emoji icon" onClick={() => setEmojiOpen((o) => !o)}>
              {page.icon || <Icon name="notes" size={18} color="var(--text-muted)"/>}
            </button>
            {emojiOpen && <EmojiPicker current={page.icon} onSelect={changeIcon} onClose={() => setEmojiOpen(false)}/>}
          </div>
          <h1 style={{ margin: 0, flex: 1 }}>{page.title}</h1>
        </div>
        <div className="row gap-3 muted text-sm" style={{ marginBottom: 24, paddingBottom: 14, borderBottom: "1px solid var(--border)" }}>
          <Avatar user={page.author} size="sm"/>
          <span>{(page.author || {}).name}</span><span>·</span>
          <span>Edited {fmtDate(page.updated_at)}</span><span>·</span>
          <span>v{page.versions ? page.versions.length : 1}</span>
        </div>

        {editing ? (
          <RichEditor
            content={draft}
            onChange={setDraft}
            editable={true}
            placeholder="Start writing… (paste from anywhere — tables, lists, headings auto-format)"
            minHeight={400}
            pageId={pageId}
            me={me}
          />
        ) : (
          <div ref={bodyRef} onClick={onBodyClick}>
            <RichEditor content={page.content} editable={false} />
          </div>
        )}

        {!editing && <ReactionBar pageId={pageId}/>}

        {sel && !composing && (
          <div style={{ position: "fixed", top: sel.rect.top - 38, left: sel.rect.left + sel.rect.width / 2 - 50, zIndex: 120 }}>
            <button className="btn btn-primary" data-size="sm" onMouseDown={(e) => { e.preventDefault(); setComposing(true); }}><Icon name="comment" size={13}/> Comment</button>
          </div>
        )}
        {sel && composing && (
          <div style={{ position: "fixed", top: sel.rect.bottom + 6, left: Math.min(sel.rect.left, window.innerWidth - 320), zIndex: 120, width: 300, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 10 }}>
            <div className="text-xs muted" style={{ marginBottom: 6 }}>On "<span className="medium" style={{ color: "var(--text)" }}>{sel.text.slice(0, 40)}</span>"</div>
            <textarea className="textarea" autoFocus value={commentText} onChange={(e) => setCommentText(e.target.value)} placeholder="Add a comment…" style={{ minHeight: 60, marginBottom: 8 }}/>
            <div className="row gap-2" style={{ justifyContent: "flex-end" }}>
              <Button data-size="sm" onClick={() => { setComposing(false); setSel(null); setCommentText(""); }}>Cancel</Button>
              <Button data-size="sm" variant="primary" disabled={!commentText.trim()} onClick={postComment}>Comment</Button>
            </div>
          </div>
        )}
      </div>

      {right === "history" && <VersionPanel pageId={pageId} page={page} onClose={() => setRight(null)} onRestored={reload}/>}
      {right === "comments" && <CommentsPanel pageId={pageId} comments={comments} activeAnchor={activeAnchor} onClose={() => { setRight(null); setActiveAnchor(null); }} reload={reloadComments}/>}
      {!right && showToc && hasToc && <TocPanel content={page.content} onClose={() => setShowToc(false)}/>}
    </div>
  );
}

// ─── HTML sanitizer (DOMPurify — battle-tested, replaces custom impl) ───────
function sanitizeHtml(dirty) {
  if (!dirty) return "";
  return DOMPurify.sanitize(dirty, {
    USE_PROFILES: { html: true },
    FORBID_TAGS: ["style", "script", "iframe", "object", "embed", "form", "base"],
    FORBID_ATTR: ["onerror", "onload", "onclick", "onmouseover"],
  });
}

// ─── ProseMirror JSON → HTML ──────────────────────────────────────────
function prosemirrorToHtml(node) {
  if (!node) return "";
  if (typeof node === "string") return node;
  const kids = () => (node.content || []).map(prosemirrorToHtml).join("");
  switch (node.type) {
    case "doc":       return kids();
    case "paragraph": return `<p>${kids() || "<br>"}</p>`;
    case "heading": {
      const l = (node.attrs && node.attrs.level) || 1;
      const text = (node.content || []).map((n) => n.text || "").join("");
      const id = "h-" + text.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "").slice(0, 50);
      return `<h${l} id="${id}">${kids()}</h${l}>`;
    }
    case "bulletList":  return `<ul>${kids()}</ul>`;
    case "orderedList": return `<ol>${kids()}</ol>`;
    case "listItem":    return `<li>${kids()}</li>`;
    case "blockquote":  return `<blockquote>${kids()}</blockquote>`;
    case "codeBlock": {
      const lang = (node.attrs && node.attrs.language) || "";
      const txt  = (node.content || []).map((n) => n.text || "").join("");
      return `<pre><code class="language-${lang}">${txt}</code></pre>`;
    }
    case "hardBreak":       return "<br>";
    case "horizontalRule":  return "<hr>";
    case "table":      return `<table style="border-collapse:collapse;width:100%;margin:16px 0">${kids()}</table>`;
    case "tableRow":   return `<tr>${kids()}</tr>`;
    case "tableHeader": return `<th style="border:1px solid var(--border);padding:8px 12px;background:var(--bg-subtle);text-align:left;font-weight:600">${kids()}</th>`;
    case "tableCell":  return `<td style="border:1px solid var(--border);padding:8px 12px">${kids()}</td>`;
    case "taskList":   return `<ul style="list-style:none;padding-left:4px">${kids()}</ul>`;
    case "taskItem": {
      const checked = node.attrs && node.attrs.checked;
      return `<li style="display:flex;gap:8px;margin:4px 0"><input type="checkbox" ${checked ? "checked" : ""} disabled style="margin-top:3px"> <span>${kids()}</span></li>`;
    }
    case "image": {
      const rawSrc = (node.attrs && node.attrs.src) || "";
      const safeSrc = /^(https?:|\/|data:image\/)/i.test(rawSrc.trim()) ? rawSrc : "";
      const alt = (node.attrs && node.attrs.alt || "").replace(/"/g, "&quot;");
      return safeSrc ? `<img src="${safeSrc.replace(/"/g, "&quot;")}" alt="${alt}" style="max-width:100%">` : "";
    }
    case "text": {
      let t = (node.text || "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
      (node.marks || []).forEach((m) => {
        if (m.type === "bold")      t = `<strong>${t}</strong>`;
        if (m.type === "italic")    t = `<em>${t}</em>`;
        if (m.type === "code")      t = `<code>${t}</code>`;
        if (m.type === "strike")    t = `<s>${t}</s>`;
        if (m.type === "underline") t = `<u>${t}</u>`;
        if (m.type === "link") {
          const rawHref = (m.attrs || {}).href || "";
          const safeHref = /^(https?:|mailto:|\/|#)/i.test(rawHref.trim()) ? rawHref : "#";
          const esc = safeHref.replace(/"/g, "&quot;").replace(/'/g, "&#39;");
          t = `<a href="${esc}" target="_blank" rel="noopener noreferrer">${t}</a>`;
        }
      });
      return t;
    }
    default: return kids();
  }
}

function textToProsemirror(text) {
  const lines = (text || "").split("\n");
  return {
    type: "doc",
    content: lines.map((line) => ({
      type: "paragraph",
      content: line ? [{ type: "text", text: line }] : [],
    })),
  };
}

function highlightBody(html, comments, activeAnchor) {
  let out = html;
  (comments || []).forEach((c) => {
    if (!c.anchor_text) return;
    const safe = c.anchor_text.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const re = new RegExp("(" + safe + ")");
    if (re.test(out) && out.indexOf("data-anchor=\"" + c.id + "\"") === -1) {
      const active = activeAnchor === c.id;
      out = out.replace(re, '<mark data-anchor="' + c.id + '" style="background:' + (active ? "#FDE68A" : "#FEF3C7") + ';border-bottom:2px solid #F59E0B;border-radius:2px;cursor:pointer;padding:0 1px">$1</mark>');
    }
  });
  return out;
}

// ─── Table of Contents ────────────────────────────────────────────────
function extractHeadings(content) {
  const headings = [];
  function walk(node) {
    if (!node) return;
    if (node.type === "heading") {
      const text = (node.content || []).map((n) => n.text || "").join("");
      if (!text.trim()) return;
      const id = "h-" + text.toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-|-$/g, "").slice(0, 50);
      headings.push({ level: (node.attrs && node.attrs.level) || 1, text, id });
    }
    (node.content || []).forEach(walk);
  }
  walk(content);
  return headings;
}

function TocPanel({ content, onClose }) {
  const headings = extractHeadings(content);
  const [active, setActive] = useState(null);

  useEffect(() => {
    const handler = () => {
      const scrollY = window.scrollY || document.documentElement.scrollTop;
      let current = null;
      for (const h of headings) {
        const el = document.getElementById(h.id);
        if (el && el.getBoundingClientRect().top < 120) current = h.id;
      }
      setActive(current);
    };
    const container = document.querySelector(".wiki-doc");
    if (container) container.addEventListener("scroll", handler);
    return () => { if (container) container.removeEventListener("scroll", handler); };
  }, [headings]);

  if (headings.length < 2) return null;

  return (
    <div style={{ width: 220, flexShrink: 0, padding: "16px 12px", borderLeft: "1px solid var(--border)", background: "var(--bg-subtle)", overflowY: "auto", position: "sticky", top: 0, maxHeight: "100vh" }}>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 10 }}>
        <span className="text-xs bold" style={{ color: "var(--text-muted)", textTransform: "uppercase", letterSpacing: "0.06em" }}>On this page</span>
        <button className="icon-btn" style={{ width: 18, height: 18 }} onClick={onClose}><Icon name="x" size={12}/></button>
      </div>
      <div className="stack" style={{ gap: 1 }}>
        {headings.map((h) => (
          <button key={h.id} className="nav-item" style={{
            paddingLeft: 4 + (h.level - 1) * 10,
            fontSize: h.level === 1 ? 13 : 12,
            fontWeight: h.level === 1 ? 500 : 400,
            color: active === h.id ? "var(--indigo-600)" : "var(--text-secondary)",
            background: active === h.id ? "var(--indigo-50)" : "transparent",
            textAlign: "left",
          }} onClick={() => {
            const el = document.getElementById(h.id);
            if (el) el.scrollIntoView({ behavior: "smooth", block: "start" });
          }}>
            {h.text}
          </button>
        ))}
      </div>
    </div>
  );
}

// ─── Emoji picker ──────────────────────────────────────────────────────
const COMMON_EMOJI = ["📝","📄","📋","📌","📎","🗂️","📁","🔧","⚙️","🔍","💡","🚀","✅","🎯","📊","📈","🗺️","📣","🔗","💬","🛡️","🌟","🔥","⚠️","ℹ️","🏷️","🧩","📦","🏗️","🔑","💎","🎉","📢","📖","🧪","🔬","🛠️","🌐","👤","👥","🤖","⏰","📅","🔄","✏️","🗒️","💻","📱","🖥️","🎨"];

function EmojiPicker({ current, onSelect, onClose }) {
  return (
    <div style={{ position: "absolute", zIndex: 200, top: "100%", left: 0, marginTop: 4, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 10, boxShadow: "var(--shadow-lg)", padding: 10, width: 260 }} onClick={(e) => e.stopPropagation()}>
      <div className="row gap-1" style={{ flexWrap: "wrap" }}>
        {current && (
          <button className="icon-btn" title="Remove icon" onClick={() => onSelect("")}
            style={{ width: 28, height: 28, fontSize: 13, border: "1px dashed var(--border)" }}>✕</button>
        )}
        {COMMON_EMOJI.map((e) => (
          <button key={e} className="icon-btn" onClick={() => onSelect(e)}
            style={{ width: 28, height: 28, fontSize: 16, border: current === e ? "2px solid var(--indigo-500)" : "1px solid transparent", borderRadius: 6 }}>
            {e}
          </button>
        ))}
      </div>
    </div>
  );
}

// ─── Reactions bar (F21) ──────────────────────────────────────────────
function ReactionBar({ pageId }) {
  const { data, setData } = useApi("/pages/" + pageId + "/reactions", [pageId]);
  const toast = useToast();
  async function toggle(emoji) {
    try { const r = await api("/pages/" + pageId + "/reactions", { method: "POST", body: { emoji } }); setData(r); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  const counts = (data && data.counts) || {};
  const mine = (data && data.my_reactions) || [];
  return (
    <div className="row gap-2" style={{ marginTop: 40, paddingTop: 20, borderTop: "1px solid var(--border)", flexWrap: "wrap" }}>
      {REACTION_EMOJI.map((emo) => {
        const on = mine.includes(emo);
        const ct = counts[emo] || 0;
        return (
          <button key={emo} onClick={() => toggle(emo)} className="row gap-2" style={{ padding: "4px 10px", borderRadius: 99, border: "1px solid " + (on ? "var(--indigo-500)" : "var(--border)"), background: on ? "var(--indigo-50)" : "var(--bg)" }}>
            <span style={{ fontSize: 15 }}>{emo}</span>
            {ct > 0 && <span className="text-xs bold" style={{ color: on ? "var(--indigo-700)" : "var(--text-secondary)" }}>{ct}</span>}
          </button>
        );
      })}
    </div>
  );
}

// ─── Version history panel (F18) ──────────────────────────────────────
function VersionPanel({ pageId, page, onClose, onRestored }) {
  const { data: versions, loading } = useApi("/pages/" + pageId + "/versions", [pageId]);
  const [preview, setPreview] = useState(null);
  const [compare, setCompare] = useState([]);
  const [diff, setDiff] = useState(null);
  const toast = useToast();

  async function viewVersion(v) {
    try { const d = await api("/pages/" + pageId + "/versions/" + v); setPreview(d); setDiff(null); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function restore(v) {
    try { const d = await api("/pages/" + pageId + "/versions/" + v); await api("/pages/" + pageId, { method: "PUT", body: { content: d.content } }); toast("Restored v" + v); setPreview(null); onRestored(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  function toggleCompare(v) {
    setCompare((c) => c.includes(v) ? c.filter((x) => x !== v) : c.length < 2 ? [...c, v] : [c[1], v]);
  }
  async function runCompare() {
    const [a, b] = compare.slice().sort((x, y) => x - y);
    try { const d = await api("/pages/" + pageId + "/versions/" + a + "/diff/" + b); setDiff(d); setPreview(null); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <aside style={{ borderLeft: "1px solid var(--border)", background: "var(--bg-subtle)", overflowY: "auto", padding: 16 }}>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <h3 style={{ margin: 0, fontSize: 14 }}>Version history</h3>
        <button className="icon-btn" onClick={onClose}><Icon name="x" size={15}/></button>
      </div>

      {compare.length === 2 && <Button data-size="sm" variant="primary" icon="history" style={{ width: "100%", justifyContent: "center", marginBottom: 12 }} onClick={runCompare}>Compare v{compare.slice().sort((a, b) => a - b).join(" ↔ v")}</Button>}

      {loading ? (
        <div className="row gap-2 text-sm muted" style={{ padding: 16 }}><MiniSpinner/> Loading…</div>
      ) : (versions || []).map((v) => (
        <div key={v.version} className="card" style={{ padding: 10, marginBottom: 8, borderColor: compare.includes(v.version) ? "var(--indigo-500)" : "var(--border)" }}>
          <div className="row gap-2" style={{ marginBottom: 6 }}>
            <span className="badge" data-tone="muted">v{v.version}</span>
            <Avatar user={v.author} size="sm"/>
            <span className="text-xs" style={{ flex: 1 }}>{(v.author || {}).name}</span>
            <span className="text-xs muted">{fmtDate(v.at)}</span>
          </div>
          <div className="text-xs muted" style={{ marginBottom: 8 }}>{v.size} bytes</div>
          <div className="row gap-1">
            <Button data-size="sm" onClick={() => viewVersion(v.version)}>Preview</Button>
            <Button data-size="sm" onClick={() => restore(v.version)}>Restore</Button>
            <button className="btn btn-ghost" data-size="sm" aria-pressed={compare.includes(v.version)} style={compare.includes(v.version) ? { background: "var(--indigo-50)", color: "var(--indigo-700)" } : null} onClick={() => toggleCompare(v.version)}>Compare</button>
          </div>
        </div>
      ))}

      <Modal open={!!preview} onClose={() => setPreview(null)} title={preview ? "Version v" + preview.version : ""} size="lg"
        footer={preview && <><Button onClick={() => setPreview(null)}>Close</Button><Button variant="primary" onClick={() => restore(preview.version)}>Restore this version</Button></>}>
        {preview && <RichEditor content={preview.content} editable={false} />}
      </Modal>

      <Modal open={!!diff} onClose={() => setDiff(null)} title={diff ? "Diff v" + diff.from + " → v" + diff.to : ""} size="lg"
        footer={<Button onClick={() => setDiff(null)}>Close</Button>}>
        {diff && (
          <div style={{ fontFamily: "ui-monospace, Menlo, monospace", fontSize: 12.5, lineHeight: 1.7 }}>
            {diff.lines.map((l, i) => (
              <div key={i} style={{ padding: "1px 8px", background: l.op === "insert" ? "rgba(16,185,129,.14)" : l.op === "delete" ? "rgba(239,68,68,.12)" : "transparent", color: l.op === "delete" ? "#B91C1C" : l.op === "insert" ? "#047857" : "var(--text)", borderRadius: 3 }}>
                <span style={{ color: "var(--text-muted)", marginRight: 8 }}>{l.op === "insert" ? "+" : l.op === "delete" ? "−" : " "}</span>{l.text}
              </div>
            ))}
          </div>
        )}
      </Modal>
    </aside>
  );
}

// ─── Inline comments panel (F19) ──────────────────────────────────────
function CommentsPanel({ pageId, comments, activeAnchor, onClose, reload }) {
  const [showResolved, setShowResolved] = useState(false);
  const toast = useToast();

  async function resolve(c) {
    try { await api("/inline-comments/" + c.id + "/" + (c.resolved ? "unresolve" : "resolve"), { method: "POST" }); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function del(c) { try { await api("/inline-comments/" + c.id, { method: "DELETE" }); toast("Comment deleted"); reload(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  let list = (comments || []).filter((c) => showResolved || !c.resolved);
  if (activeAnchor) list = list.filter((c) => c.id === activeAnchor || c.anchor_text);
  const resolvedCount = (comments || []).filter((c) => c.resolved).length;

  return (
    <aside style={{ borderLeft: "1px solid var(--border)", background: "var(--bg-subtle)", overflowY: "auto", padding: 16 }}>
      <div className="row" style={{ justifyContent: "space-between", marginBottom: 12 }}>
        <h3 style={{ margin: 0, fontSize: 14 }}>Comments</h3>
        <button className="icon-btn" onClick={onClose}><Icon name="x" size={15}/></button>
      </div>
      {resolvedCount > 0 && (
        <label className="row gap-2" style={{ marginBottom: 12 }}>
          <Switch on={showResolved} onChange={setShowResolved}/><span className="text-xs muted">Show {resolvedCount} resolved</span>
        </label>
      )}
      {list.length === 0 && <div className="text-sm muted" style={{ padding: 8 }}>No comments. Select text in the page to add one.</div>}
      <div className="stack gap-2">
        {list.map((c) => (
          <div key={c.id} className="card" style={{ padding: 10, opacity: c.resolved ? 0.6 : 1, borderColor: activeAnchor === c.id ? "var(--indigo-500)" : "var(--border)" }}>
            {c.anchor_text && <div style={{ fontSize: 11, color: "#B45309", background: "var(--warning-bg)", borderRadius: 4, padding: "2px 6px", marginBottom: 6, display: "inline-block" }}>"{c.anchor_text.slice(0, 36)}"</div>}
            <div className="row gap-2" style={{ marginBottom: 4 }}>
              <Avatar user={c.user} size="sm"/>
              <span className="text-xs bold">{(c.user || {}).name}</span>
              <span className="text-xs muted">{fmtDate(c.created_at)}</span>
            </div>
            <div className="text-sm" style={{ marginBottom: 8 }}>{c.text}</div>
            <div className="row gap-1">
              <Button data-size="sm" variant="ghost" icon={c.resolved ? "refresh" : "check"} onClick={() => resolve(c)}>{c.resolved ? "Reopen" : "Resolve"}</Button>
              <button className="icon-btn" style={{ width: 24, height: 24 }} onClick={() => del(c)} title="Delete"><Icon name="trash" size={13}/></button>
            </div>
          </div>
        ))}
      </div>
    </aside>
  );
}

// ─── Blog post document (F22) ─────────────────────────────────────────
function BlogDoc({ space, blogId, setBlogId }) {
  const { data: post, loading, reload } = useApi(blogId ? "/blog-posts/" + blogId : null, [blogId]);
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState({ title: "", body: "" });
  const toast = useToast();

  useEffect(() => { if (post) setDraft({ title: post.title, body: post.body }); setEditing(false); }, [blogId, post && post.id]);

  if (!blogId) return <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No post selected" hint="Pick a post or create one."/></div>;
  if (loading || !post) return <div className="wiki-doc"><div className="skel" style={{ height: 36, width: "50%", marginBottom: 16 }}/><div className="skel" style={{ height: 14 }}/></div>;

  async function save() {
    try { await api("/blog-posts/" + blogId, { method: "PUT", body: draft }); toast("Post saved"); setEditing(false); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function togglePublish() {
    try { await api("/blog-posts/" + blogId + "/" + (post.is_published ? "unpublish" : "publish"), { method: "POST" }); toast(post.is_published ? "Unpublished" : "Published"); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="wiki-doc">
      <div className="row gap-2 text-sm muted" style={{ marginBottom: 14 }}>
        <span>{space.name}</span> <Icon name="chevronRight" size={11}/> <span>Blog</span>
        <div style={{ flex: 1 }}/>
        {!post.is_published && <Badge tone="warning">Draft</Badge>}
        {editing ? (
          <><button className="btn btn-primary" data-size="sm" onClick={save}><Icon name="check" size={13}/> Save</button><button className="btn btn-ghost" data-size="sm" onClick={() => setEditing(false)}>Cancel</button></>
        ) : (
          <>
            <button className="btn btn-ghost" data-size="sm" onClick={() => { setDraft({ title: post.title, body: post.body }); setEditing(true); }}><Icon name="pencil" size={13}/> Edit</button>
            <button className={"btn " + (post.is_published ? "btn-ghost" : "btn-primary")} data-size="sm" onClick={togglePublish}>{post.is_published ? "Unpublish" : "Publish"}</button>
          </>
        )}
      </div>

      {editing ? (
        <>
          <input className="input" value={draft.title} onChange={(e) => setDraft({ ...draft, title: e.target.value })} style={{ fontSize: 28, fontWeight: 600, padding: "6px 10px", marginBottom: 14, border: 0 }}/>
          <textarea className="textarea" value={draft.body} onChange={(e) => setDraft({ ...draft, body: e.target.value })} style={{ minHeight: 320, fontFamily: "ui-monospace, Menlo, monospace", fontSize: 13 }}/>
        </>
      ) : (
        <>
          <h1>{post.title}</h1>
          <div className="row gap-3 muted text-sm" style={{ marginBottom: 24, paddingBottom: 14, borderBottom: "1px solid var(--border)" }}>
            <Avatar user={post.author} size="sm"/><span>{(post.author || {}).name}</span><span>·</span>
            <span>{post.published_at ? "Published " + fmtDate(post.published_at) : "Unpublished draft"}</span>
          </div>
          <div dangerouslySetInnerHTML={{ __html: sanitizeHtml(post.body) }}/>
        </>
      )}
    </div>
  );
}

// ─── Space modals (F16) ───────────────────────────────────────────────
function NewSpaceModal({ open, onClose, onCreated }) {
  const [form, setForm]   = useState({ name: "", type: "team", description: "" });
  const [error, setError] = useState(null);
  const toast = useToast();

  useEffect(() => {
    if (open) { setForm({ name: "", type: "team", description: "" }); setError(null); }
  }, [open]);

  const canCreate = form.name.trim().length >= 2;

  async function create() {
    if (!canCreate) return;
    setError(null);
    try {
      const s = await api("/spaces", { method: "POST", body: { name: form.name, type: form.type, description: form.description || undefined } });
      toast("Space created");
      onClose();
      onCreated(s);
    } catch (e) {
      setError(e.message);
    }
  }

  return (
    <Modal open={open} onClose={onClose} title="Create space"
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="primary" disabled={!canCreate} onClick={create}>Create space</Button></>}>
      <div className="stack gap-3">
        {error && <div style={{ padding: "8px 12px", borderRadius: 7, background: "#FEF2F2", color: "#991B1B", fontSize: 13, border: "1px solid #FECACA" }}>{error}</div>}
        <div>
          <label className="label">Name</label>
          <input className="input" autoFocus value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="e.g. Platform Team"
            onKeyDown={(e) => e.key === "Enter" && canCreate && create()}/>
        </div>
        <div>
          <label className="label">Type</label>
          <div className="row gap-2">
            {Object.entries(SPACE_TYPE_META).map(([k, m]) => (
              <button key={k} className="btn" onClick={() => setForm({ ...form, type: k })}
                style={{ flex: 1, justifyContent: "center", border: "1px solid " + (form.type === k ? "var(--indigo-600)" : "var(--border)"), background: form.type === k ? "var(--indigo-50)" : "var(--bg)", color: form.type === k ? "var(--indigo-700)" : "var(--text)" }}>
                <Icon name={m.icon} size={14}/> {m.label}
              </button>
            ))}
          </div>
        </div>
        <div>
          <label className="label">Description <span style={{ fontWeight: 400, color: "var(--text-muted)" }}>(optional)</span></label>
          <textarea className="textarea" value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="What is this space for?" style={{ minHeight: 60 }}/>
        </div>
      </div>
    </Modal>
  );
}

function SpaceSettingsModal({ open, onClose, space, onChanged, onDeleted }) {
  const [tab, setTab] = useState("general");
  const [form, setForm] = useState({ name: space.name, description: space.description || "" });
  const [confirm, setConfirm] = useState(false);
  const toast = useToast();
  useEffect(() => { if (open) { setForm({ name: space.name, description: space.description || "" }); setTab("general"); } }, [open, space]);

  async function save() { try { await api("/spaces/" + space.id, { method: "PUT", body: form }); toast("Space updated"); onChanged(); onClose(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  async function archive() { try { await api("/spaces/" + space.id + "/archive", { method: "POST" }); toast("Space archived"); onChanged(); onClose(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  async function del() { try { await api("/spaces/" + space.id, { method: "DELETE" }); toast("Space deleted"); onClose(); onDeleted(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  const TABS = [["general", "General"], ["members", "Members"]];

  return (
    <Modal open={open} onClose={onClose} title="Space settings"
      footer={tab === "general" ? <><Button onClick={onClose}>Cancel</Button><Button variant="primary" onClick={save}>Save changes</Button></> : <Button onClick={onClose}>Close</Button>}>
      <div className="row gap-1" style={{ marginBottom: 16, borderBottom: "1px solid var(--border)", paddingBottom: 0 }}>
        {TABS.map(([id, label]) => (
          <button key={id} className="tab" aria-selected={tab === id} onClick={() => setTab(id)} style={{ paddingBottom: 10 }}>{label}</button>
        ))}
      </div>

      {tab === "general" && (
        <div className="stack gap-3">
          <div><label className="label">Name</label><input className="input" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}/></div>
          <div><label className="label">Description</label><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} style={{ minHeight: 60 }}/></div>
          <div className="card card-pad" style={{ borderColor: "var(--danger)", background: "var(--danger-bg)" }}>
            <div className="bold text-sm" style={{ color: "#B91C1C", marginBottom: 8 }}>Danger zone</div>
            <div className="row gap-2">
              <Button onClick={archive} icon="archive">Archive space</Button>
              <Button variant="danger" icon="trash" onClick={() => setConfirm(true)}>Delete space</Button>
            </div>
          </div>
        </div>
      )}

      {tab === "members" && <SpaceMembersTab spaceId={space.id} open={open}/>}

      <ConfirmDelete open={confirm} onClose={() => setConfirm(false)} onConfirm={del} title="Delete space" body={`Delete "${space.name}" and all its pages? This cannot be undone.`}/>
    </Modal>
  );
}

const SPACE_ROLES = ["admin", "member", "viewer"];

function SpaceMembersTab({ spaceId, open }) {
  const { data: membersData, reload: reloadMembers } = useApi("/spaces/" + spaceId + "/members?limit=100", [spaceId, open]);
  const { data: allUsers } = useApi("/users?limit=200");
  const [addUserID, setAddUserID] = useState("");
  const [addRole, setAddRole] = useState("member");
  const [search, setSearch] = useState("");
  const [removing, setRemoving] = useState(null);
  const toast = useToast();

  const members = membersData?.items || membersData || [];
  const memberUserIDs = new Set((members).map((m) => m.user_id || m.user?.id));
  const userList = (allUsers?.items || allUsers || []).filter((u) => !memberUserIDs.has(u.id));
  const filtered = userList.filter((u) => {
    const q = search.toLowerCase();
    return !q || (u.full_name || u.name || "").toLowerCase().includes(q) || (u.email || "").toLowerCase().includes(q);
  });

  async function addMember() {
    if (!addUserID) return;
    try {
      await api("/spaces/" + spaceId + "/members", { method: "POST", body: { user_id: addUserID, role: addRole } });
      toast("Member added"); setAddUserID(""); setSearch(""); reloadMembers();
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function changeRole(userID, role) {
    try { await api("/spaces/" + spaceId + "/members/" + userID, { method: "PUT", body: { role } }); reloadMembers(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  async function removeMember(userID) {
    try { await api("/spaces/" + spaceId + "/members/" + userID, { method: "DELETE" }); toast("Member removed"); reloadMembers(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setRemoving(null); }
  }

  const selectedUser = userList.find((u) => u.id === addUserID);

  return (
    <div className="stack gap-4">
      {/* Add member */}
      <div className="card card-pad stack gap-2">
        <div className="bold text-sm">Add member</div>
        <input className="input" placeholder="Search by name or email…" value={search}
          onChange={(e) => { setSearch(e.target.value); setAddUserID(""); }}/>
        {search && !selectedUser && (
          <div style={{ maxHeight: 180, overflowY: "auto", border: "1px solid var(--border)", borderRadius: 7 }}>
            {filtered.length === 0 && <div className="text-xs muted" style={{ padding: 10 }}>No users found</div>}
            {filtered.map((u) => (
              <button key={u.id} className="nav-item" style={{ color: "var(--text)" }}
                onClick={() => { setAddUserID(u.id); setSearch(u.full_name || u.name || u.email); }}>
                <Avatar user={u} size="sm"/>
                <span style={{ flex: 1 }}>{u.full_name || u.name}</span>
                <span className="text-xs muted">{u.email}</span>
              </button>
            ))}
          </div>
        )}
        {selectedUser && (
          <div className="row gap-2" style={{ padding: "6px 10px", borderRadius: 7, background: "var(--bg-subtle)", border: "1px solid var(--border)" }}>
            <Avatar user={selectedUser} size="sm"/>
            <span style={{ flex: 1 }} className="text-sm">{selectedUser.full_name || selectedUser.name}</span>
            <select className="input" style={{ width: 100, padding: "2px 6px" }} value={addRole} onChange={(e) => setAddRole(e.target.value)}>
              {SPACE_ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
            </select>
            <Button variant="primary" data-size="sm" onClick={addMember}>Add</Button>
          </div>
        )}
      </div>

      {/* Members list */}
      <div>
        <div className="bold text-sm" style={{ marginBottom: 8 }}>Current members ({members.length})</div>
        {members.length === 0 && <div className="text-xs muted">No members yet.</div>}
        <div className="stack gap-2">
          {members.map((m) => {
            const user = m.user || {};
            const uid = m.user_id || user.id;
            return (
              <div key={uid} className="row gap-2" style={{ padding: "8px 10px", borderRadius: 7, border: "1px solid var(--border)", background: "var(--bg)" }}>
                <Avatar user={user} size="sm"/>
                <div className="stack" style={{ flex: 1, gap: 1 }}>
                  <span className="text-sm bold">{user.full_name || user.name || "—"}</span>
                  <span className="text-xs muted">{user.email}</span>
                </div>
                <select className="input" style={{ width: 90, padding: "2px 6px", fontSize: 12 }}
                  value={m.role} onChange={(e) => changeRole(uid, e.target.value)}>
                  {SPACE_ROLES.map((r) => <option key={r} value={r}>{r}</option>)}
                </select>
                {removing === uid ? (
                  <div className="row gap-1">
                    <Button data-size="sm" variant="danger" onClick={() => removeMember(uid)}>Remove</Button>
                    <Button data-size="sm" onClick={() => setRemoving(null)}>Cancel</Button>
                  </div>
                ) : (
                  <button className="icon-btn" title="Remove member" onClick={() => setRemoving(uid)}><Icon name="trash" size={13}/></button>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
