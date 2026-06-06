// views_wiki.jsx — Confluence-style wiki. Overrides WikiView from views_misc.
// Features 16 (spaces) · 17 (page tree) · 18 (versions) · 19 (inline comments)
// · 20 (export) · 21 (reactions) · 22 (blog).

const SPACE_TYPE_META = { team: { icon: "users", label: "Team" }, personal: { icon: "user", label: "Personal" }, project: { icon: "briefcase", label: "Project" } };
const REACTION_EMOJI = ["👍", "❤️", "🎉", "🚀", "👀", "😄"];

function WikiView({ nav }) {
  const { data: spaces, reload: reloadSpaces } = useApi("/spaces");
  const [spaceId, setSpaceId] = React.useState(null);
  const [mode, setMode] = React.useState("page"); // page | blog
  const [pageId, setPageId] = React.useState(null);
  const [blogId, setBlogId] = React.useState(null);
  const [newSpaceOpen, setNewSpaceOpen] = React.useState(false);
  const [spaceSettings, setSpaceSettings] = React.useState(false);

  // pick first space once loaded
  React.useEffect(() => { if (spaces && spaces.length && !spaceId) setSpaceId(spaces[0].id); }, [spaces]);

  const space = (spaces || []).find((s) => s.id === spaceId);

  return (
    <div className="wiki" style={{ gridTemplateColumns: "264px 1fr" }}>
      <aside className="wiki-tree" style={{ display: "flex", flexDirection: "column", padding: 0 }}>
        {/* Spaces */}
        <div style={{ padding: "12px 10px 8px", borderBottom: "1px solid var(--border)" }}>
          <div className="row" style={{ justifyContent: "space-between", marginBottom: 8 }}>
            <h4 style={{ margin: 0 }}>Spaces</h4>
            <button className="icon-btn" style={{ width: 22, height: 22 }} title="New space" onClick={() => setNewSpaceOpen(true)}><Icon name="plus" size={13}/></button>
          </div>
          <div className="stack gap-1">
            {(spaces || []).map((s) => (
              <button key={s.id} className="tree-item" aria-current={s.id === spaceId ? "page" : undefined} style={{ width: "100%" }} onClick={() => { setSpaceId(s.id); setPageId(null); setBlogId(null); setMode("page"); }}>
                <span style={{ width: 18, height: 18, borderRadius: 5, background: "var(--indigo-50)", color: "var(--indigo-600)", display: "grid", placeItems: "center" }}><Icon name={SPACE_TYPE_META[s.type].icon} size={11}/></span>
                <span style={{ flex: 1, textAlign: "left" }}>{s.name}</span>
                <span className="text-xs muted">{s.page_count}</span>
              </button>
            ))}
          </div>
        </div>

        {/* Page tree + blog for current space */}
        {space && <SpaceNav key={space.id} space={space} mode={mode} setMode={setMode} pageId={pageId} setPageId={setPageId} blogId={blogId} setBlogId={setBlogId} onSettings={() => setSpaceSettings(true)}/>}
      </aside>

      {space ? (
        mode === "blog"
          ? <BlogDoc key={blogId || "blog-empty"} space={space} blogId={blogId} setBlogId={setBlogId}/>
          : <PageDoc key={pageId || "page-empty"} space={space} pageId={pageId} setPageId={setPageId} nav={nav}/>
      ) : (
        <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No space selected" hint="Create a space to get started."/></div>
      )}

      <NewSpaceModal open={newSpaceOpen} onClose={() => setNewSpaceOpen(false)} onCreated={(s) => { reloadSpaces(); setSpaceId(s.id); setPageId(null); }}/>
      {space && <SpaceSettingsModal open={spaceSettings} onClose={() => setSpaceSettings(false)} space={space} onChanged={reloadSpaces} onDeleted={() => { reloadSpaces(); setSpaceId(null); }}/>}
    </div>
  );
}

// ─── Space nav: page tree (F17) + blog list (F22) ─────────────────────
function SpaceNav({ space, mode, setMode, pageId, setPageId, blogId, setBlogId, onSettings }) {
  const { data: tree, loading, reload } = useApi("/spaces/" + space.id + "/pages/tree", [space.id]);
  const { data: blog, reload: reloadBlog } = useApi("/spaces/" + space.id + "/blog-posts", [space.id]);
  const [expanded, setExpanded] = React.useState({});
  const [menu, setMenu] = React.useState(null); // {x,y,node}
  const [moveNode, setMoveNode] = React.useState(null);
  const [renaming, setRenaming] = React.useState(null);
  const toast = useToast();

  // auto-select first page
  React.useEffect(() => { if (tree && tree.length && !pageId && mode === "page") { setPageId(tree[0].id); setExpanded((e) => ({ ...e, [tree[0].id]: true })); } }, [tree]);
  React.useEffect(() => {
    const close = () => setMenu(null);
    if (menu) { document.addEventListener("click", close); return () => document.removeEventListener("click", close); }
  }, [menu]);

  async function newPage(parentId) {
    try { const p = await api("/spaces/" + space.id + "/pages", { method: "POST", body: { title: "Untitled", parent_id: parentId || null } }); toast("Page created"); reload(); setMode("page"); setPageId(p.id); if (parentId) setExpanded((e) => ({ ...e, [parentId]: true })); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
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
          ) : <span style={{ flex: 1 }}>{node.title}</span>}
          <button className="icon-btn" style={{ width: 18, height: 18 }} title="New child page" onClick={(e) => { e.stopPropagation(); newPage(node.id); }}><Icon name="plus" size={11}/></button>
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
          <button className="icon-btn" style={{ width: 22, height: 22 }} title="Space settings" onClick={onSettings}><Icon name="settings" size={13}/></button>
          <button className="icon-btn" style={{ width: 22, height: 22 }} title="New page" onClick={() => newPage(null)}><Icon name="plus" size={13}/></button>
        </div>
      </div>

      {loading ? <div className="row gap-2 text-xs muted" style={{ padding: 8 }}><MiniSpinner/> Loading…</div> : (tree || []).map((n) => <TreeNode key={n.id} node={n} depth={0}/>)}
      {!loading && (tree || []).length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>No pages yet.</div>}

      {/* Blog */}
      <div style={{ marginTop: 16, borderTop: "1px solid var(--border)", paddingTop: 10 }}>
        <div className="row" style={{ justifyContent: "space-between", padding: "0 6px 6px" }}>
          <h4 style={{ margin: 0 }}>Blog</h4>
          <button className="icon-btn" style={{ width: 22, height: 22 }} title="New post" onClick={async () => { try { const p = await api("/spaces/" + space.id + "/blog-posts", { method: "POST", body: { title: "Untitled post" } }); reloadBlog(); setMode("blog"); setBlogId(p.id); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }}><Icon name="plus" size={13}/></button>
        </div>
        <div className="stack gap-1">
          {(blog || []).map((p) => (
            <button key={p.id} className="tree-item" aria-current={mode === "blog" && blogId === p.id ? "page" : undefined} style={{ width: "100%", alignItems: "flex-start" }} onClick={() => { setMode("blog"); setBlogId(p.id); }}>
              <Icon name="notes" size={12} color="var(--text-muted)" style={{ marginTop: 3 }}/>
              <span className="stack" style={{ flex: 1, textAlign: "left", lineHeight: 1.3 }}>
                <span className="row gap-2">{p.title}{!p.published && <Badge tone="warning" style={{ fontSize: 9 }}>Draft</Badge>}</span>
                <span className="text-xs muted">{(p.author || {}).name} · {p.published_at ? fmtDate(p.published_at) : "unpublished"}</span>
              </span>
            </button>
          ))}
          {(blog || []).length === 0 && <div className="text-xs muted" style={{ padding: 8 }}>No posts yet.</div>}
        </div>
      </div>

      {/* Context menu (F17) */}
      {menu && (
        <div style={{ position: "fixed", top: menu.y, left: menu.x, zIndex: 200, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, minWidth: 170 }} onClick={(e) => e.stopPropagation()}>
          {[["plus", "New child page", () => { newPage(menu.node.id); setMenu(null); }], ["pencil", "Rename", () => { setRenaming(menu.node.id); setMenu(null); }], ["arrowRight", "Move", () => { setMoveNode(menu.node); setMenu(null); }], ["trash", "Delete", () => { del(menu.node); setMenu(null); }, true]].map(([ic, label, fn, danger], i) => (
            <button key={i} className="nav-item" style={{ color: danger ? "var(--danger)" : "var(--text)", fontSize: 13 }} onClick={fn}><Icon name={ic} size={13}/> {label}</button>
          ))}
        </div>
      )}

      {/* Move modal */}
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
    </div>
  );
}

// ─── Page document (F18 versions, F19 comments, F20 export, F21 reactions)
function PageDoc({ space, pageId, setPageId, nav }) {
  const { data: page, loading, reload, setData } = useApi(pageId ? "/pages/" + pageId : null, [pageId]);
  const [editing, setEditing] = React.useState(false);
  const [draft, setDraft] = React.useState("");
  const [right, setRight] = React.useState(null); // null | history | comments
  const [exportOpen, setExportOpen] = React.useState(false);
  const bodyRef = React.useRef(null);
  const exportRef = React.useRef(null);
  const toast = useToast();

  const { data: comments, reload: reloadComments } = useApi(pageId ? "/pages/" + pageId + "/inline-comments" : null, [pageId]);
  const [sel, setSel] = React.useState(null); // {text, rect}
  const [composing, setComposing] = React.useState(false);
  const [commentText, setCommentText] = React.useState("");
  const [activeAnchor, setActiveAnchor] = React.useState(null);

  React.useEffect(() => { if (page) setDraft(page.body); setEditing(false); setRight(null); }, [pageId]);
  React.useEffect(() => {
    if (!exportOpen) return;
    const h = (e) => { if (exportRef.current && !exportRef.current.contains(e.target)) setExportOpen(false); };
    document.addEventListener("mousedown", h); return () => document.removeEventListener("mousedown", h);
  }, [exportOpen]);

  if (!pageId) return <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No page selected" hint="Pick a page from the tree."/></div>;
  if (loading || !page) return <div className="wiki-doc"><div className="skel" style={{ height: 36, width: "60%", marginBottom: 16 }}/><div className="skel" style={{ height: 14, marginBottom: 8 }}/><div className="skel" style={{ height: 14, width: "80%" }}/></div>;

  async function save() {
    try { await api("/pages/" + pageId, { method: "PUT", body: { body: draft } }); toast("Page saved"); setEditing(false); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function exportAs(fmt) {
    setExportOpen(false);
    try { await apiDownload("/pages/" + pageId + "/export/" + fmt); toast("Exporting as " + fmt.toUpperCase() + "…"); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  // inline comment selection (F19)
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
  const bodyHtml = highlightBody(page.body, openComments, activeAnchor);

  function onBodyClick(e) {
    const mark = e.target.closest && e.target.closest("[data-anchor]");
    if (mark) { setActiveAnchor(mark.getAttribute("data-anchor")); setRight("comments"); }
  }

  return (
    <div style={{ display: "grid", gridTemplateColumns: right ? "1fr 340px" : "1fr", overflow: "hidden" }}>
      <div className="wiki-doc" onMouseUp={onMouseUp} style={{ position: "relative" }}>
        {/* toolbar */}
        <div className="row gap-2 text-sm muted" style={{ marginBottom: 14 }}>
          <span>{space.name}</span> <Icon name="chevronRight" size={11}/>
          <span style={{ color: "var(--text)" }}>{page.title}</span>
          <div style={{ flex: 1 }}/>
          {editing ? (
            <>
              <button className="btn btn-primary" data-size="sm" onClick={save}><Icon name="check" size={13}/> Save</button>
              <button className="btn btn-ghost" data-size="sm" onClick={() => { setEditing(false); setDraft(page.body); }}>Cancel</button>
            </>
          ) : (
            <>
              <button className="btn btn-ghost" data-size="sm" onClick={() => { setDraft(page.body); setEditing(true); }}><Icon name="pencil" size={13}/> Edit</button>
              <button className="btn btn-ghost" data-size="sm" aria-pressed={right === "history"} style={right === "history" ? { background: "var(--bg-subtle)", color: "var(--text)" } : null} onClick={() => setRight(right === "history" ? null : "history")}><Icon name="history" size={13}/> History</button>
              <button className="btn btn-ghost" data-size="sm" aria-pressed={right === "comments"} style={right === "comments" ? { background: "var(--bg-subtle)", color: "var(--text)" } : null} onClick={() => { setActiveAnchor(null); setRight(right === "comments" ? null : "comments"); }}>
                <Icon name="comment" size={13}/> Comments {openComments.length > 0 && <span className="badge" data-tone="info" style={{ padding: "0 5px" }}>{openComments.length}</span>}
              </button>
              <span ref={exportRef} style={{ position: "relative" }}>
                <button className="btn btn-ghost" data-size="sm" onClick={() => setExportOpen((o) => !o)}><Icon name="download" size={13}/> Export</button>
                {exportOpen && (
                  <div style={{ position: "absolute", top: "100%", right: 0, marginTop: 4, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50, minWidth: 160 }}>
                    {[["pdf", "Export as PDF"], ["html", "Export as HTML"], ["md", "Export as Markdown"], ["docx", "Export as DOCX"]].map(([f, l]) => (
                      <button key={f} className="nav-item" style={{ color: "var(--text)", fontSize: 13 }} onClick={() => exportAs(f)}><Icon name="download" size={13}/> {l}</button>
                    ))}
                  </div>
                )}
              </span>
            </>
          )}
        </div>

        <h1>{page.title}</h1>
        <div className="row gap-3 muted text-sm" style={{ marginBottom: 24, paddingBottom: 14, borderBottom: "1px solid var(--border)" }}>
          <Avatar user={page.author} size="sm"/>
          <span>{(page.author || {}).name}</span><span>·</span>
          <span>Edited {fmtDate(page.updated_at)}</span><span>·</span>
          <span>v{page.versions ? page.versions.length : 1}</span>
        </div>

        {editing ? (
          <textarea className="textarea" value={draft} onChange={(e) => setDraft(e.target.value)} style={{ minHeight: 360, fontFamily: "ui-monospace, Menlo, monospace", fontSize: 13 }}/>
        ) : (
          <div ref={bodyRef} onClick={onBodyClick} dangerouslySetInnerHTML={{ __html: bodyHtml }}/>
        )}

        {/* Reactions (F21) */}
        {!editing && <ReactionBar pageId={pageId}/>}

        {/* selection tooltip (F19) */}
        {sel && !composing && (
          <div style={{ position: "fixed", top: sel.rect.top - 38, left: sel.rect.left + sel.rect.width / 2 - 50, zIndex: 120 }}>
            <button className="btn btn-primary" data-size="sm" onMouseDown={(e) => { e.preventDefault(); setComposing(true); }}><Icon name="comment" size={13}/> Comment</button>
          </div>
        )}
        {sel && composing && (
          <div style={{ position: "fixed", top: sel.rect.bottom + 6, left: Math.min(sel.rect.left, window.innerWidth - 320), zIndex: 120, width: 300, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8, boxShadow: "var(--shadow-lg)", padding: 10 }}>
            <div className="text-xs muted" style={{ marginBottom: 6 }}>On “<span className="medium" style={{ color: "var(--text)" }}>{sel.text.slice(0, 40)}</span>”</div>
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
    </div>
  );
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
  const [preview, setPreview] = React.useState(null); // version obj content
  const [compare, setCompare] = React.useState([]); // up to 2 version numbers
  const [diff, setDiff] = React.useState(null);
  const toast = useToast();

  async function viewVersion(v) {
    try { const d = await api("/pages/" + pageId + "/versions/" + v); setPreview(d); setDiff(null); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function restore(v) {
    try { const d = await api("/pages/" + pageId + "/versions/" + v); await api("/pages/" + pageId, { method: "PUT", body: { body: d.body } }); toast("Restored v" + v); setPreview(null); onRestored(); }
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

      {loading ? <Loading/> : (versions || []).map((v) => (
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

      {/* Preview modal */}
      <Modal open={!!preview} onClose={() => setPreview(null)} title={preview ? "Version v" + preview.version : ""} size="lg"
        footer={preview && <><Button onClick={() => setPreview(null)}>Close</Button><Button variant="primary" onClick={() => restore(preview.version)}>Restore this version</Button></>}>
        {preview && <div className="wiki-doc" style={{ padding: 0, maxWidth: "none" }} dangerouslySetInnerHTML={{ __html: preview.body }}/>}
      </Modal>

      {/* Diff modal */}
      <Modal open={!!diff} onClose={() => setDiff(null)} title={diff ? "Diff v" + diff.from + " → v" + diff.to : ""} size="lg"
        footer={<Button onClick={() => setDiff(null)}>Close</Button>}>
        {diff && (
          <div style={{ fontFamily: "ui-monospace, Menlo, monospace", fontSize: 12.5, lineHeight: 1.7 }}>
            {diff.lines.map((l, i) => (
              <div key={i} style={{ padding: "1px 8px", background: l.type === "added" ? "rgba(16,185,129,.14)" : l.type === "removed" ? "rgba(239,68,68,.12)" : "transparent", color: l.type === "removed" ? "#B91C1C" : l.type === "added" ? "#047857" : "var(--text)", borderRadius: 3 }}>
                <span style={{ color: "var(--text-muted)", marginRight: 8 }}>{l.type === "added" ? "+" : l.type === "removed" ? "−" : " "}</span>{l.text}
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
  const [showResolved, setShowResolved] = React.useState(false);
  const [reply, setReply] = React.useState({});
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
            {c.anchor_text && <div style={{ fontSize: 11, color: "#B45309", background: "var(--warning-bg)", borderRadius: 4, padding: "2px 6px", marginBottom: 6, display: "inline-block" }}>“{c.anchor_text.slice(0, 36)}”</div>}
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
  const [editing, setEditing] = React.useState(false);
  const [draft, setDraft] = React.useState({ title: "", body: "" });
  const toast = useToast();

  React.useEffect(() => { if (post) setDraft({ title: post.title, body: post.body }); setEditing(false); }, [blogId, post && post.id]);

  if (!blogId) return <div style={{ display: "grid", placeItems: "center" }}><Empty icon="notes" title="No post selected" hint="Pick a post or create one."/></div>;
  if (loading || !post) return <div className="wiki-doc"><div className="skel" style={{ height: 36, width: "50%", marginBottom: 16 }}/><div className="skel" style={{ height: 14 }}/></div>;

  async function save() {
    try { await api("/blog-posts/" + blogId, { method: "PUT", body: draft }); toast("Post saved"); setEditing(false); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  async function togglePublish() {
    try { await api("/blog-posts/" + blogId + "/" + (post.published ? "unpublish" : "publish"), { method: "POST" }); toast(post.published ? "Unpublished" : "Published"); reload(); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  return (
    <div className="wiki-doc">
      <div className="row gap-2 text-sm muted" style={{ marginBottom: 14 }}>
        <span>{space.name}</span> <Icon name="chevronRight" size={11}/> <span>Blog</span>
        <div style={{ flex: 1 }}/>
        {!post.published && <Badge tone="warning">Draft</Badge>}
        {editing ? (
          <><button className="btn btn-primary" data-size="sm" onClick={save}><Icon name="check" size={13}/> Save</button><button className="btn btn-ghost" data-size="sm" onClick={() => setEditing(false)}>Cancel</button></>
        ) : (
          <>
            <button className="btn btn-ghost" data-size="sm" onClick={() => { setDraft({ title: post.title, body: post.body }); setEditing(true); }}><Icon name="pencil" size={13}/> Edit</button>
            <button className={"btn " + (post.published ? "btn-ghost" : "btn-primary")} data-size="sm" onClick={togglePublish}>{post.published ? "Unpublish" : "Publish"}</button>
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
          <div dangerouslySetInnerHTML={{ __html: post.body }}/>
        </>
      )}
    </div>
  );
}

// ─── Space modals (F16) ───────────────────────────────────────────────
function NewSpaceModal({ open, onClose, onCreated }) {
  const [form, setForm] = React.useState({ name: "", key: "", type: "team", description: "" });
  const toast = useToast();
  React.useEffect(() => { if (open) setForm({ name: "", key: "", type: "team", description: "" }); }, [open]);
  async function create() {
    if (!form.name.trim()) return;
    try { const s = await api("/spaces", { method: "POST", body: form }); toast("Space created"); onClose(); onCreated(s); }
    catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }
  return (
    <Modal open={open} onClose={onClose} title="Create space"
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="primary" disabled={!form.name.trim()} onClick={create}>Create space</Button></>}>
      <div className="stack gap-3">
        <div className="row gap-3">
          <div style={{ flex: 2 }}><label className="label">Name</label><input className="input" autoFocus value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="e.g. Platform"/></div>
          <div style={{ flex: 1 }}><label className="label">Key</label><input className="input mono" value={form.key} onChange={(e) => setForm({ ...form, key: e.target.value.toUpperCase() })} maxLength="6" placeholder="PLAT"/></div>
        </div>
        <div><label className="label">Type</label>
          <div className="row gap-2">
            {Object.entries(SPACE_TYPE_META).map(([k, m]) => (
              <button key={k} className="btn" onClick={() => setForm({ ...form, type: k })} style={{ flex: 1, justifyContent: "center", border: "1px solid " + (form.type === k ? "var(--indigo-600)" : "var(--border)"), background: form.type === k ? "var(--indigo-50)" : "var(--bg)", color: form.type === k ? "var(--indigo-700)" : "var(--text)" }}><Icon name={m.icon} size={14}/> {m.label}</button>
            ))}
          </div>
        </div>
        <div><label className="label">Description</label><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} style={{ minHeight: 60 }}/></div>
      </div>
    </Modal>
  );
}

function SpaceSettingsModal({ open, onClose, space, onChanged, onDeleted }) {
  const [form, setForm] = React.useState({ name: space.name, description: space.description });
  const [confirm, setConfirm] = React.useState(false);
  const toast = useToast();
  React.useEffect(() => { if (open) setForm({ name: space.name, description: space.description }); }, [open, space]);

  async function save() { try { await api("/spaces/" + space.id, { method: "PUT", body: form }); toast("Space updated"); onChanged(); onClose(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  async function archive() { try { await api("/spaces/" + space.id + "/archive", { method: "POST" }); toast("Space archived"); onChanged(); onClose(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }
  async function del() { try { await api("/spaces/" + space.id, { method: "DELETE" }); toast("Space deleted"); onClose(); onDeleted(); } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); } }

  return (
    <Modal open={open} onClose={onClose} title="Space settings"
      footer={<><Button onClick={onClose}>Cancel</Button><Button variant="primary" onClick={save}>Save changes</Button></>}>
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
      <ConfirmDelete open={confirm} onClose={() => setConfirm(false)} onConfirm={del} title="Delete space" body={`Delete "${space.name}" and all its pages? This cannot be undone.`}/>
    </Modal>
  );
}

Object.assign(window, { WikiView });
