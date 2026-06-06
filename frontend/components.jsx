// components.jsx — Forge shared UI primitives

const { useState, useEffect, useRef, useMemo, useCallback, useLayoutEffect } = React;

// ─── Avatar ─────────────────────────────────────────────
function Avatar({ user, size = "" }) {
  if (!user) return null;
  return (
    <span className="avatar" data-size={size} style={{ background: user.color }}>
      {user.initials}
    </span>
  );
}

function AvatarStack({ users, max = 4, size = "sm" }) {
  const shown = users.slice(0, max);
  const extra = users.length - shown.length;
  return (
    <span className="avatar-stack">
      {shown.map((u) => <Avatar key={u.id} user={u} size={size} />)}
      {extra > 0 && (
        <span className="avatar" data-size={size} style={{ background: "var(--bg-muted)", color: "var(--text-secondary)", fontWeight: 600 }}>
          +{extra}
        </span>
      )}
    </span>
  );
}

// ─── Badge ──────────────────────────────────────────────
function Badge({ children, tone = "muted", dot = false, style }) {
  return (
    <span className="badge" data-tone={tone} style={style}>
      {dot && <span className="dot"/>}
      {children}
    </span>
  );
}

function PriorityBadge({ value }) {
  const m = FORGE_DATA.PRIORITY_META[value];
  if (!m) return null;
  return (
    <span title={value + " priority"} style={{ display: "inline-flex", alignItems: "center", gap: 4, color: m.color, fontSize: 12, fontWeight: 500 }}>
      <Icon name={m.icon} size={14} strokeWidth={2.4}/>
      {value}
    </span>
  );
}

function TypeIcon({ value, size = 14 }) {
  const m = FORGE_DATA.TYPE_META[value];
  if (!m) return null;
  return (
    <span title={value} style={{ width: 18, height: 18, borderRadius: 4, background: m.color, color: "#fff", display: "inline-grid", placeItems: "center", flexShrink: 0 }}>
      <Icon name={m.icon} size={size - 2} strokeWidth={2.2}/>
    </span>
  );
}

function StatusBadge({ value }) {
  const m = FORGE_DATA.STATUS_META[value] || { tone: "muted" };
  return <Badge tone={m.tone} dot>{value}</Badge>;
}

// ─── Button ─────────────────────────────────────────────
function Button({ children, variant = "secondary", size, icon, iconRight, onClick, type = "button", disabled, style, ...rest }) {
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      data-size={size}
      className={"btn btn-" + variant}
      style={style}
      {...rest}
    >
      {icon && <Icon name={icon} size={15}/>}
      {children}
      {iconRight && <Icon name={iconRight} size={15}/>}
    </button>
  );
}

// ─── Switch ─────────────────────────────────────────────
function Switch({ on, onChange, label }) {
  return (
    <button
      role="switch"
      aria-checked={on}
      className="switch"
      data-on={on}
      onClick={() => onChange && onChange(!on)}
      style={{ border: 0 }}
      aria-label={label}
    />
  );
}

// ─── Modal ──────────────────────────────────────────────
function Modal({ open, onClose, title, children, footer, size }) {
  useEffect(() => {
    if (!open) return;
    const h = (e) => e.key === "Escape" && onClose && onClose();
    window.addEventListener("keydown", h);
    return () => window.removeEventListener("keydown", h);
  }, [open, onClose]);
  if (!open) return null;
  return (
    <div className="modal-backdrop" onClick={onClose}>
      <div className="modal" data-size={size} onClick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
        <div className="modal-head">
          <h2>{title}</h2>
          <button className="icon-btn" onClick={onClose} aria-label="Close"><Icon name="x" size={16}/></button>
        </div>
        <div className="modal-body">{children}</div>
        {footer && <div className="modal-foot">{footer}</div>}
      </div>
    </div>
  );
}

// ─── Toast ──────────────────────────────────────────────
const ToastCtx = React.createContext(null);
function ToastProvider({ children }) {
  const [toasts, set] = useState([]);
  const push = useCallback((msg, opt = {}) => {
    const id = Math.random().toString(36).slice(2);
    set((t) => [...t, { id, msg, ...opt }]);
    setTimeout(() => set((t) => t.filter((x) => x.id !== id)), opt.duration || 2600);
  }, []);
  return (
    <ToastCtx.Provider value={push}>
      {children}
      <div className="toast-region">
        {toasts.map((t) => (
          <div key={t.id} className="toast">
            {t.icon !== false && <Icon name={t.icon || "check"} size={16} color={t.color || "#34D399"}/>}
            <span>{t.msg}</span>
          </div>
        ))}
      </div>
    </ToastCtx.Provider>
  );
}
const useToast = () => React.useContext(ToastCtx);

// ─── Dropdown menu (very simple) ────────────────────────
function Menu({ trigger, items, align = "left" }) {
  const [open, setOpen] = useState(false);
  const ref = useRef(null);
  useEffect(() => {
    if (!open) return;
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false); };
    document.addEventListener("mousedown", h);
    return () => document.removeEventListener("mousedown", h);
  }, [open]);
  return (
    <span ref={ref} style={{ position: "relative", display: "inline-block" }}>
      <span onClick={() => setOpen((o) => !o)}>{trigger}</span>
      {open && (
        <div style={{
          position: "absolute", top: "100%", marginTop: 6,
          [align]: 0,
          minWidth: 180, background: "var(--bg)",
          border: "1px solid var(--border)", borderRadius: 8,
          boxShadow: "var(--shadow-lg)", padding: 4, zIndex: 50
        }}>
          {items.map((it, i) => it.divider ? (
            <div key={i} style={{ height: 1, background: "var(--border)", margin: 4 }}/>
          ) : (
            <button key={i} className="nav-item" onClick={() => { setOpen(false); it.onClick && it.onClick(); }}
              style={{ color: it.danger ? "var(--danger)" : "var(--text)", fontSize: 13 }}>
              {it.icon && <Icon name={it.icon} size={14}/>} {it.label}
            </button>
          ))}
        </div>
      )}
    </span>
  );
}

// ─── Hook: hash router ──────────────────────────────────
function useHashRoute() {
  const get = () => {
    const h = window.location.hash.replace(/^#\/?/, "");
    const [path, qs] = h.split("?");
    const params = Object.fromEntries(new URLSearchParams(qs || ""));
    const parts = (path || "dashboard").split("/").filter(Boolean);
    return { path: parts[0] || "dashboard", rest: parts.slice(1), params, raw: h };
  };
  const [route, setRoute] = useState(get);
  useEffect(() => {
    const h = () => setRoute(get());
    window.addEventListener("hashchange", h);
    return () => window.removeEventListener("hashchange", h);
  }, []);
  const nav = useCallback((path) => {
    window.location.hash = "#/" + path.replace(/^#?\/?/, "");
  }, []);
  return [route, nav];
}

// ─── Empty placeholder ──────────────────────────────────
function Empty({ icon = "inbox", title, hint, action }) {
  return (
    <div style={{ display: "grid", placeItems: "center", padding: "48px 24px", textAlign: "center" }}>
      <div style={{ width: 56, height: 56, borderRadius: 14, background: "var(--bg-subtle)", border: "1px solid var(--border)", display: "grid", placeItems: "center", color: "var(--text-muted)", marginBottom: 12 }}>
        <Icon name={icon} size={24}/>
      </div>
      <div style={{ fontWeight: 600, marginBottom: 4 }}>{title}</div>
      {hint && <div className="muted text-sm" style={{ marginBottom: 12 }}>{hint}</div>}
      {action}
    </div>
  );
}

// Expose
Object.assign(window, {
  Avatar, AvatarStack, Badge, PriorityBadge, TypeIcon, StatusBadge,
  Button, Switch, Modal, ToastProvider, useToast, Menu, useHashRoute, Empty,
});
