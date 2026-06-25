import { useEditor, EditorContent } from '@tiptap/react';
import { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { markdownToTiptap } from '../lib/markdownToTiptap';
import StarterKit from '@tiptap/starter-kit';
import { Table } from '@tiptap/extension-table';
import { TableRow } from '@tiptap/extension-table-row';
import { TableHeader } from '@tiptap/extension-table-header';
import { TableCell } from '@tiptap/extension-table-cell';
import { TaskList } from '@tiptap/extension-task-list';
import { TaskItem } from '@tiptap/extension-task-item';
import { Placeholder } from '@tiptap/extension-placeholder';
import { Image } from '@tiptap/extension-image';
import { Link } from '@tiptap/extension-link';
import { Underline } from '@tiptap/extension-underline';
import { Highlight } from '@tiptap/extension-highlight';
import { TextAlign } from '@tiptap/extension-text-align';
import * as Y from 'yjs';
import { HocuspocusProvider } from '@hocuspocus/provider';
import Collaboration from '@tiptap/extension-collaboration';
import CollaborationCursor from '@tiptap/extension-collaboration-cursor';

const COLLAB_URL = import.meta.env.VITE_COLLAB_URL || 'ws://localhost:1234';

// Deterministic color from a string (user id / name)
function stringToColor(str) {
  let h = 0;
  for (let i = 0; i < str.length; i++) h = (h << 5) - h + str.charCodeAt(i);
  const hue = Math.abs(h) % 360;
  return `hsl(${hue},65%,50%)`;
}

// ─── Extensions list ────────────────────────────────────────────────────────
function buildExtensions({ placeholder, ydoc }) {
  const base = [
    ydoc
      ? StarterKit.configure({ history: false, heading: { levels: [1, 2, 3] } })
      : StarterKit.configure({ heading: { levels: [1, 2, 3] } }),
    Underline,
    Highlight.configure({ multicolor: true }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Link.configure({ openOnClick: false, HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' } }),
    Image.configure({ inline: false, allowBase64: true }),
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    TaskList,
    TaskItem.configure({ nested: true }),
    Placeholder.configure({ placeholder: placeholder || 'Write something…' }),
  ];

  if (ydoc) {
    base.push(Collaboration.configure({ document: ydoc }));
  }

  return base;
}

// ─── Toolbar button ──────────────────────────────────────────────────────────
function Btn({ onClick, active, title, children, disabled }) {
  return (
    <button
      type="button"
      title={title}
      disabled={disabled}
      onMouseDown={(e) => { e.preventDefault(); onClick(); }}
      style={{
        padding: '3px 7px',
        borderRadius: 5,
        border: 'none',
        background: active ? 'var(--indigo-100)' : 'transparent',
        color: active ? 'var(--indigo-700)' : 'var(--text)',
        cursor: 'pointer',
        fontSize: 13,
        fontWeight: 500,
        lineHeight: 1.4,
        minWidth: 26,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 3,
      }}
    >
      {children}
    </button>
  );
}

function Sep() {
  return <div style={{ width: 1, background: 'var(--border)', margin: '2px 4px', alignSelf: 'stretch' }} />;
}

// ─── Markdown import modal ───────────────────────────────────────────────────
function MarkdownImportModal({ onImport, onClose }) {
  const [text, setText] = useState('');
  return (
    <div
      style={{
        position: 'fixed', inset: 0, zIndex: 1000,
        background: 'rgba(0,0,0,0.45)',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
      }}
      onClick={(e) => { if (e.target === e.currentTarget) onClose(); }}
    >
      <div style={{
        background: 'var(--bg)', borderRadius: 10,
        padding: 24, width: 640, maxWidth: '95vw',
        boxShadow: '0 8px 40px rgba(0,0,0,0.22)',
        display: 'flex', flexDirection: 'column', gap: 12,
      }}>
        <div style={{ fontWeight: 600, fontSize: 15 }}>Markdown import</div>
        <textarea
          autoFocus
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder="Paste your Markdown here…"
          style={{
            width: '100%', height: 320, resize: 'vertical',
            fontFamily: 'var(--font-mono, monospace)', fontSize: 13,
            padding: 10, borderRadius: 6, border: '1px solid var(--border)',
            background: 'var(--bg-subtle, var(--bg))', color: 'var(--text)',
            boxSizing: 'border-box',
          }}
        />
        <div style={{ display: 'flex', gap: 8, justifyContent: 'flex-end' }}>
          <button
            type="button"
            onClick={onClose}
            style={{
              padding: '6px 16px', borderRadius: 6, border: '1px solid var(--border)',
              background: 'transparent', cursor: 'pointer', fontSize: 13,
            }}
          >Cancel</button>
          <button
            type="button"
            onClick={() => { if (text.trim()) { onImport(text); onClose(); } }}
            style={{
              padding: '6px 16px', borderRadius: 6, border: 'none',
              background: 'var(--indigo-600, #4F46E5)', color: '#fff',
              cursor: 'pointer', fontSize: 13, fontWeight: 500,
            }}
          >Import</button>
        </div>
      </div>
    </div>
  );
}

// ─── Fixed toolbar (editing mode) ───────────────────────────────────────────
function Toolbar({ editor }) {
  const [showMdImport, setShowMdImport] = useState(false);
  if (!editor) return null;

  function setLink() {
    const prev = editor.getAttributes('link').href || '';
    const url = window.prompt('URL:', prev);
    if (url === null) return;
    if (url === '') { editor.chain().focus().unsetLink().run(); return; }
    editor.chain().focus().setLink({ href: url }).run();
  }

  function insertTable() {
    editor.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run();
  }

  function insertImage() {
    const url = window.prompt('Image URL:');
    if (url) editor.chain().focus().setImage({ src: url }).run();
  }

  return (
    <div style={{
      display: 'flex',
      flexWrap: 'wrap',
      gap: 2,
      padding: '6px 8px',
      borderBottom: '1px solid var(--border)',
      background: 'var(--bg)',
      position: 'sticky',
      top: 0,
      zIndex: 10,
      alignItems: 'center',
    }}>
      {/* Heading */}
      <select
        value={
          editor.isActive('heading', { level: 1 }) ? '1' :
          editor.isActive('heading', { level: 2 }) ? '2' :
          editor.isActive('heading', { level: 3 }) ? '3' : '0'
        }
        onChange={(e) => {
          const v = Number(e.target.value);
          if (v === 0) editor.chain().focus().setParagraph().run();
          else editor.chain().focus().toggleHeading({ level: v }).run();
        }}
        style={{
          fontSize: 12, padding: '2px 6px', borderRadius: 5,
          border: '1px solid var(--border)', background: 'var(--bg)',
          color: 'var(--text)', cursor: 'pointer', height: 26,
        }}
      >
        <option value="0">Paragraph</option>
        <option value="1">Heading 1</option>
        <option value="2">Heading 2</option>
        <option value="3">Heading 3</option>
      </select>

      <Sep/>

      {/* Inline formatting */}
      <Btn onClick={() => editor.chain().focus().toggleBold().run()} active={editor.isActive('bold')} title="Bold (Ctrl+B)"><strong>B</strong></Btn>
      <Btn onClick={() => editor.chain().focus().toggleItalic().run()} active={editor.isActive('italic')} title="Italic (Ctrl+I)"><em>I</em></Btn>
      <Btn onClick={() => editor.chain().focus().toggleUnderline().run()} active={editor.isActive('underline')} title="Underline (Ctrl+U)"><u>U</u></Btn>
      <Btn onClick={() => editor.chain().focus().toggleStrike().run()} active={editor.isActive('strike')} title="Strikethrough"><s>S</s></Btn>
      <Btn onClick={() => editor.chain().focus().toggleCode().run()} active={editor.isActive('code')} title="Inline code">{"`c`"}</Btn>
      <Btn onClick={() => editor.chain().focus().toggleHighlight().run()} active={editor.isActive('highlight')} title="Highlight">▌</Btn>
      <Btn onClick={setLink} active={editor.isActive('link')} title="Link">🔗</Btn>

      <Sep/>

      {/* Alignment */}
      <Btn onClick={() => editor.chain().focus().setTextAlign('left').run()} active={editor.isActive({ textAlign: 'left' })} title="Align left">⫷</Btn>
      <Btn onClick={() => editor.chain().focus().setTextAlign('center').run()} active={editor.isActive({ textAlign: 'center' })} title="Align center">☰</Btn>
      <Btn onClick={() => editor.chain().focus().setTextAlign('right').run()} active={editor.isActive({ textAlign: 'right' })} title="Align right">⫸</Btn>

      <Sep/>

      {/* Lists */}
      <Btn onClick={() => editor.chain().focus().toggleBulletList().run()} active={editor.isActive('bulletList')} title="Bullet list">• List</Btn>
      <Btn onClick={() => editor.chain().focus().toggleOrderedList().run()} active={editor.isActive('orderedList')} title="Ordered list">1. List</Btn>
      <Btn onClick={() => editor.chain().focus().toggleTaskList().run()} active={editor.isActive('taskList')} title="Task list">☑ Tasks</Btn>

      <Sep/>

      {/* Blocks */}
      <Btn onClick={() => editor.chain().focus().toggleBlockquote().run()} active={editor.isActive('blockquote')} title="Blockquote">" "</Btn>
      <Btn onClick={() => editor.chain().focus().toggleCodeBlock().run()} active={editor.isActive('codeBlock')} title="Code block">{'</>'}</Btn>
      <Btn onClick={() => editor.chain().focus().setHorizontalRule().run()} title="Divider">—</Btn>

      <Sep/>

      {/* Table */}
      <Btn onClick={insertTable} title="Insert table">⊞ Table</Btn>
      {editor.isActive('table') && (
        <>
          <Btn onClick={() => editor.chain().focus().addColumnAfter().run()} title="Add column">+Col</Btn>
          <Btn onClick={() => editor.chain().focus().addRowAfter().run()} title="Add row">+Row</Btn>
          <Btn onClick={() => editor.chain().focus().deleteColumn().run()} title="Delete column">-Col</Btn>
          <Btn onClick={() => editor.chain().focus().deleteRow().run()} title="Delete row">-Row</Btn>
          <Btn onClick={() => editor.chain().focus().deleteTable().run()} title="Delete table">✕ Table</Btn>
        </>
      )}

      <Sep/>

      {/* Image */}
      <Btn onClick={insertImage} title="Insert image">🖼</Btn>

      <Sep/>

      {/* Undo/Redo */}
      <Btn onClick={() => editor.chain().focus().undo().run()} disabled={!editor.can().undo()} title="Undo">↩</Btn>
      <Btn onClick={() => editor.chain().focus().redo().run()} disabled={!editor.can().redo()} title="Redo">↪</Btn>

      <Sep/>

      {/* Markdown import */}
      <Btn onClick={() => setShowMdImport(true)} title="Import Markdown">MD ↓</Btn>
      {showMdImport && (
        <MarkdownImportModal
          onClose={() => setShowMdImport(false)}
          onImport={(md) => {
            const doc = markdownToTiptap(md);
            editor.commands.setContent(doc, true);
          }}
        />
      )}
    </div>
  );
}


// ─── Slash command definitions ───────────────────────────────────────────────
const SLASH_COMMANDS = [
  { icon: 'H1', label: 'Heading 1',    keywords: ['h1','heading','title'],      action: (e) => e.chain().focus().toggleHeading({ level: 1 }).run() },
  { icon: 'H2', label: 'Heading 2',    keywords: ['h2','heading','section'],    action: (e) => e.chain().focus().toggleHeading({ level: 2 }).run() },
  { icon: 'H3', label: 'Heading 3',    keywords: ['h3','heading','subsection'], action: (e) => e.chain().focus().toggleHeading({ level: 3 }).run() },
  { icon: '•',  label: 'Bullet list',  keywords: ['bullet','list','ul'],        action: (e) => e.chain().focus().toggleBulletList().run() },
  { icon: '1.', label: 'Ordered list', keywords: ['ordered','numbered','ol'],   action: (e) => e.chain().focus().toggleOrderedList().run() },
  { icon: '☑',  label: 'Task list',    keywords: ['task','todo','check'],       action: (e) => e.chain().focus().toggleTaskList().run() },
  { icon: '"',  label: 'Blockquote',   keywords: ['quote','blockquote'],        action: (e) => e.chain().focus().toggleBlockquote().run() },
  { icon: '</>', label: 'Code block',  keywords: ['code','codeblock','pre'],    action: (e) => e.chain().focus().toggleCodeBlock().run() },
  { icon: '—',  label: 'Divider',      keywords: ['divider','hr','rule','line'],action: (e) => e.chain().focus().setHorizontalRule().run() },
  { icon: '⊞',  label: 'Table',        keywords: ['table','grid'],              action: (e) => e.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run() },
  { icon: '🖼',  label: 'Image',        keywords: ['image','photo','img'],       action: (e) => { const url = window.prompt('Image URL:'); if (url) e.chain().focus().setImage({ src: url }).run(); } },
];

function SlashMenu({ pos, filter, onSelect, onClose }) {
  const [activeIdx, setActiveIdx] = useState(0);
  const filtered = SLASH_COMMANDS.filter((c) => {
    if (!filter) return true;
    const f = filter.toLowerCase();
    return c.label.toLowerCase().includes(f) || c.keywords.some((k) => k.includes(f));
  });

  useEffect(() => { setActiveIdx(0); }, [filter]);

  useEffect(() => {
    const h = (e) => {
      if (e.key === 'ArrowDown') { e.preventDefault(); e.stopPropagation(); setActiveIdx((i) => Math.min(filtered.length - 1, i + 1)); }
      else if (e.key === 'ArrowUp') { e.preventDefault(); e.stopPropagation(); setActiveIdx((i) => Math.max(0, i - 1)); }
      else if (e.key === 'Enter') { e.preventDefault(); e.stopPropagation(); if (filtered[activeIdx]) onSelect(filtered[activeIdx]); }
      else if (e.key === 'Escape') { e.preventDefault(); e.stopPropagation(); onClose(); }
    };
    window.addEventListener('keydown', h, true);
    return () => window.removeEventListener('keydown', h, true);
  }, [filtered, activeIdx, onSelect, onClose]);

  if (filtered.length === 0) return null;

  return (
    <div style={{
      position: 'fixed', top: pos.top, left: pos.left, zIndex: 9999,
      background: 'var(--bg)', border: '1px solid var(--border)', borderRadius: 10,
      boxShadow: '0 8px 30px rgba(0,0,0,0.18)', padding: 4, minWidth: 220, maxHeight: 300, overflowY: 'auto',
    }}>
      <div style={{ padding: '4px 10px 6px', fontSize: 11, color: 'var(--text-muted)', fontWeight: 600, textTransform: 'uppercase', letterSpacing: '0.06em' }}>Blocks</div>
      {filtered.map((cmd, i) => (
        <div key={cmd.label}
          style={{
            display: 'flex', alignItems: 'center', gap: 10, padding: '7px 10px',
            borderRadius: 7, cursor: 'pointer', fontSize: 13,
            background: i === activeIdx ? 'var(--indigo-50)' : 'transparent',
            color: i === activeIdx ? 'var(--indigo-700)' : 'var(--text)',
          }}
          onMouseEnter={() => setActiveIdx(i)}
          onMouseDown={(e) => { e.preventDefault(); onSelect(cmd); }}>
          <span style={{ width: 26, textAlign: 'center', fontWeight: 700, fontSize: 12, color: 'var(--text-muted)' }}>{cmd.icon}</span>
          <span>{cmd.label}</span>
        </div>
      ))}
    </div>
  );
}

// ─── Presence bar — shows other editors ─────────────────────────────────────
function PresenceBar({ provider }) {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    if (!provider) return;
    const update = () => {
      const states = [];
      provider.awareness.getStates().forEach((state, clientId) => {
        if (clientId !== provider.awareness.clientID && state.user) {
          states.push(state.user);
        }
      });
      setUsers(states);
    };
    provider.awareness.on('change', update);
    update();
    return () => provider.awareness.off('change', update);
  }, [provider]);

  if (users.length === 0) return null;

  return (
    <div style={{
      display: 'flex', alignItems: 'center', gap: 6,
      padding: '4px 12px', fontSize: 12,
      borderBottom: '1px solid var(--border)',
      background: 'var(--bg-subtle, var(--bg))',
      color: 'var(--text-muted)',
    }}>
      <span>Also editing:</span>
      {users.slice(0, 8).map((u, i) => (
        <span key={i} title={u.name} style={{
          display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
          width: 22, height: 22, borderRadius: '50%',
          background: u.color, color: '#fff',
          fontSize: 10, fontWeight: 700,
        }}>
          {u.name ? u.name[0].toUpperCase() : '?'}
        </span>
      ))}
      {users.length > 8 && <span>+{users.length - 8}</span>}
    </div>
  );
}

// ─── Collab hook — creates Y.Doc + provider, tears them down on unmount ──────
function useCollab(pageId, me) {
  const [collab, setCollab] = useState(null);

  useEffect(() => {
    if (!pageId) return;

    const token = localStorage.getItem('access_token') || '';
    const ydoc = new Y.Doc();
    const provider = new HocuspocusProvider({
      url: COLLAB_URL,
      name: `page:${pageId}`,
      document: ydoc,
      token,
      onConnect() {
        if (me) {
          provider.awareness.setLocalStateField('user', {
            name: me.full_name || me.username || 'Unknown',
            color: stringToColor(me.id || me.username || Math.random().toString()),
          });
        }
      },
    });

    setCollab({ ydoc, provider });

    return () => {
      provider.destroy();
      ydoc.destroy();
      setCollab(null);
    };
  }, [pageId, me?.id]);

  return collab;
}

// ─── Main RichEditor component ───────────────────────────────────────────────
export function RichEditor({ content, onChange, editable = true, placeholder, minHeight = 400, pageId, me }) {
  const collab = useCollab(pageId, me);
  const [slashMenu, setSlashMenu] = useState(null);
  const slashMenuRef = useRef(null);
  slashMenuRef.current = slashMenu;

  const extensions = useMemo(
    () => buildExtensions({ placeholder, ydoc: collab?.ydoc }),
    // Rebuild extensions only when collab ydoc changes (pageId switch)
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [placeholder, collab?.ydoc]
  );

  // Add CollaborationCursor lazily once provider is ready
  const allExtensions = useMemo(() => {
    if (!collab) return extensions;
    return [
      ...extensions,
      CollaborationCursor.configure({
        provider: collab.provider,
        user: me ? {
          name: me.full_name || me.username || 'Unknown',
          color: stringToColor(me.id || me.username || ''),
        } : { name: 'Unknown', color: '#888' },
      }),
    ];
  }, [extensions, collab, me]);

  const checkSlash = useCallback((editor) => {
    if (!editor.isEditable) return;
    const { state } = editor;
    const { selection } = state;
    const { $from } = selection;
    const blockText = $from.parent.textContent.slice(0, $from.parentOffset);

    if (blockText === '/' || (blockText.startsWith('/') && !blockText.includes(' '))) {
      const filter = blockText.slice(1);
      try {
        const domInfo = editor.view.domAtPos($from.pos);
        const el = domInfo.node.nodeType === 3 ? domInfo.node.parentElement : domInfo.node;
        const rect = el.getBoundingClientRect();
        const menuTop = Math.min(rect.bottom + 4, window.innerHeight - 320);
        setSlashMenu({ top: menuTop, left: Math.min(rect.left, window.innerWidth - 240), filter });
      } catch { setSlashMenu(null); }
    } else {
      if (slashMenuRef.current) setSlashMenu(null);
    }
  }, []);

  const editor = useEditor({
    extensions: allExtensions,
    // In collab mode the Y.Doc owns the content — don't seed with prop
    content: collab ? undefined : (content || { type: 'doc', content: [{ type: 'paragraph' }] }),
    editable,
    onUpdate: ({ editor }) => {
      if (!collab) onChange?.(editor.getJSON());
      if (editable) checkSlash(editor);
    },
  }, [collab, editable]);

  function executeSlash(cmd) {
    if (!editor) return;
    const { state } = editor;
    const { $from } = state.selection;
    const blockStart = $from.start();
    const curPos = $from.pos;
    editor.chain().focus().deleteRange({ from: blockStart, to: curPos }).run();
    cmd.action(editor);
    setSlashMenu(null);
  }

  // Close slash menu on click outside
  useEffect(() => {
    if (!slashMenu) return;
    const h = (e) => { if (!e.target.closest('[data-slash-menu]')) setSlashMenu(null); };
    document.addEventListener('mousedown', h);
    return () => document.removeEventListener('mousedown', h);
  }, [slashMenu]);

  // Sync content prop when NOT in collab mode and page changes externally
  if (!collab && editor && !editor.isDestroyed) {
    const current = JSON.stringify(editor.getJSON());
    const incoming = JSON.stringify(content);
    if (incoming && incoming !== current && !editor.isFocused) {
      editor.commands.setContent(content, false);
    }
  }

  return (
    <div className="rich-editor-wrap" data-editable={editable} style={{ position: 'relative' }}>
      {editable && <Toolbar editor={editor} />}
      {collab && <PresenceBar provider={collab.provider} />}
      <EditorContent
        editor={editor}
        style={{ minHeight: editable ? minHeight : undefined }}
      />
      {slashMenu && editable && (
        <div data-slash-menu>
          <SlashMenu
            pos={slashMenu}
            filter={slashMenu.filter}
            onSelect={executeSlash}
            onClose={() => setSlashMenu(null)}
          />
        </div>
      )}
    </div>
  );
}

export default RichEditor;
