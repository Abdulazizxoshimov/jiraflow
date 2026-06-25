// members.jsx — Members table, add modal, member profile, integrations, notification settings
import { useState, useEffect, useRef } from 'react';
import { Icon } from '../components/icons';
import { Avatar, Badge, Button, Modal, Switch, Menu, useToast } from '../components/components';
import { useApp } from '../store/AppContext';
import { api } from '../api/api';
import { adaptUser, userInitials } from '../api/adapters';
import { COLORS } from '../store/data';
import { Pill } from './board';
import { MiniSpinner } from '../panels/issue';

const ROLE_META = {
  admin:     { tone: "danger" },
  manager:   { tone: "warning" },
  developer: { tone: "info" },
  viewer:    { tone: "muted" },
  member:    { tone: "info" },
  owner:     { tone: "purple" },
};
const roleLabel = (r) => r ? r.charAt(0).toUpperCase() + r.slice(1) : "Member";

const NOTIF_EVENTS = [
  { key: "assigned",   label: "Issue assigned to me",     desc: "Triggered when someone assigns an issue directly to you.",                    default: true },
  { key: "mentioned",  label: "Mentioned in comment",     desc: "Someone @mentions you in a comment or description.",                          default: true },
  { key: "commented",  label: "Comment on watched issue", desc: "New comments on issues you are watching.",                                    default: false },
  { key: "status",     label: "Issue status changed",     desc: "Status changes on issues you created or are watching.",                       default: false },
  { key: "sprint",     label: "Sprint events",            desc: "Sprint started, completed, or a goal changed.",                               default: true },
  { key: "digest",     label: "Daily digest",             desc: "A morning summary of open issues, PRs awaiting review, and sprint progress.", default: true },
];

const EMAIL_PREFS = [
  ["email_assigned",  "Issue assigned to me",     "When someone assigns an issue to you."],
  ["email_mentioned", "Mentioned in comment",     "Someone @mentions you in a comment."],
  ["email_commented", "Comment on watched issue", "New comments on issues you watch."],
  ["email_status",    "Issue status changed",     "Status changes on issues you watch."],
  ["email_watcher",   "Updates on watched items", "Any update to items you watch."],
];

const CHANGE_ROLE_OPTIONS = [
  { value: "admin",  label: "Admin",  desc: "Full access including user management" },
  { value: "member", label: "Member", desc: "Can create and manage issues and pages" },
  { value: "viewer", label: "Viewer", desc: "Read-only access to all content" },
];

export function MembersView({ nav }) {
  const { people, setPeople, me } = useApp();
  const isAdmin = me?.role === "admin";
  const [openAdd, setOpenAdd] = useState(false);
  const [openProfile, setOpenProfile] = useState(null);
  const [search, setSearch]   = useState("");
  const [filter, setFilter]   = useState("all");
  const [roleTarget, setRoleTarget] = useState(null); // user for role modal
  const [roleValue, setRoleValue]   = useState("member");
  const [savingRole, setSavingRole] = useState(false);
  const toast = useToast();

  const list = people.filter((p) => {
    if (search && !(p.name.toLowerCase().includes(search.toLowerCase()) || p.email.toLowerCase().includes(search.toLowerCase()))) return false;
    if (filter !== "all" && p.role !== filter) return false;
    return true;
  });

  async function handleRemove(u) {
    try {
      await api("/users/" + u.id + "/deactivate", { method: "POST" });
      setPeople((p) => p.filter((x) => x.id !== u.id));
      toast(u.name + " removed from workspace");
    } catch (e) {
      toast(e.message, { icon: "x", color: "#F87171" });
    }
  }

  function openRoleChange(u) { setRoleTarget(u); setRoleValue(u.role?.toLowerCase() || "member"); }

  async function saveRole() {
    if (!roleTarget) return;
    setSavingRole(true);
    try {
      await api("/users/" + roleTarget.id, { method: "PUT", body: { role: roleValue } });
      setPeople((p) => p.map((x) => x.id === roleTarget.id ? { ...x, role: roleValue.toLowerCase() } : x));
      toast("Role updated");
      setRoleTarget(null);
    } catch (e) {
      toast(e.message, { icon: "x", color: "#F87171" });
    } finally {
      setSavingRole(false);
    }
  }

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><Icon name="users" size={11}/> <span>Members</span></div>
          <h1>Members</h1>
          <p>Manage who can access this project and how they receive notifications.</p>
        </div>
        <div className="row gap-2">
          <Button icon="download">Export</Button>
          {isAdmin && <Button variant="primary" icon="plus" onClick={() => setOpenAdd(true)}>Invite member</Button>}
        </div>
      </div>

      <div style={{ padding: "0 32px 32px" }}>
        <div className="card" style={{ overflow: "hidden" }}>
          <div className="row" style={{ padding: 12, gap: 8, borderBottom: "1px solid var(--border)" }}>
            <div className="search" style={{ width: 280, padding: "4px 10px" }}>
              <Icon name="search" size={13}/>
              <input placeholder="Search by name or email…" value={search} onChange={(e) => setSearch(e.target.value)}/>
            </div>
            <Pill label="Role" icon="shield" value={filter === "all" ? "All" : filter} options={[{ id: "all", label: "All roles" }, ...Object.keys(ROLE_META).map((r) => ({ id: r, label: r }))]} onChange={setFilter}/>
            <div style={{ flex: 1 }}/>
            <span className="text-xs muted">{list.length} members</span>
          </div>
          <table className="table">
            <thead>
              <tr>
                <th>Member</th>
                <th style={{ width: 130 }}>Role</th>
                <th style={{ width: 110 }}>Status</th>
                <th style={{ width: 110 }}>Joined</th>
                <th style={{ width: 60 }}/>
              </tr>
            </thead>
            <tbody>
              {list.map((u) => (
                <tr key={u.id} onClick={() => setOpenProfile(u.id)} style={{ cursor: "default" }}>
                  <td>
                    <div className="row gap-3">
                      <Avatar user={u} size="md"/>
                      <div className="stack" style={{ lineHeight: 1.25 }}>
                        <span className="bold">{u.name}</span>
                        <span className="text-xs muted">{u.email}</span>
                      </div>
                    </div>
                  </td>
                  <td><Badge tone={(ROLE_META[u.role] || ROLE_META.member).tone}>{roleLabel(u.role)}</Badge></td>
                  <td>
                    <Badge tone={u.status === "Active" ? "success" : u.status === "Pending" ? "warning" : "muted"} dot>
                      {u.status}
                    </Badge>
                  </td>
                  <td className="text-sm secondary">{u.joined}</td>
                  <td onClick={(e) => e.stopPropagation()}>
                    <Menu
                      align="right"
                      trigger={<button className="icon-btn"><Icon name="moreH" size={15}/></button>}
                      items={[
                        { icon: "user",  label: "View profile", onClick: () => setOpenProfile(u.id) },
                        ...(isAdmin && u.id !== me?.id ? [
                          { icon: "shield", label: "Change role", onClick: () => openRoleChange(u) },
                          { divider: true },
                          { icon: "trash", label: "Remove from workspace", danger: true, onClick: () => handleRemove(u) },
                        ] : []),
                      ]}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <AddMemberModal open={openAdd} onClose={() => setOpenAdd(false)} onSuccess={(invite) => {
        setPeople((p) => [...p, {
          id: "pending-" + invite.id,
          name: invite.email,
          email: invite.email,
          initials: invite.email[0].toUpperCase(),
          color: COLORS[p.length % COLORS.length],
          role: invite.role,
          status: "Pending",
          joined: "Invited",
        }]);
        toast && toast("Invite link created for " + invite.email);
      }}/>
      <MemberProfileDrawer open={!!openProfile} userId={openProfile} onClose={() => setOpenProfile(null)} people={people}/>

      <Modal open={!!roleTarget} onClose={() => setRoleTarget(null)} title="Change role"
        footer={<>
          <Button onClick={() => setRoleTarget(null)}>Cancel</Button>
          <Button variant="primary" disabled={savingRole} onClick={saveRole}>{savingRole ? "Saving…" : "Save"}</Button>
        </>}>
        {roleTarget && (
          <div>
            <div className="text-sm secondary" style={{ marginBottom: 14 }}>
              Change role for <strong>{roleTarget.name}</strong>
            </div>
            <div style={{ display: "grid", gap: 8 }}>
              {CHANGE_ROLE_OPTIONS.map((r) => (
                <label key={r.value} style={{ display: "flex", alignItems: "center", gap: 10, padding: "10px 12px", borderRadius: 8, border: "1px solid " + (roleValue === r.value ? "var(--indigo-400)" : "var(--border)"), background: roleValue === r.value ? "var(--indigo-50)" : "var(--bg)", cursor: "pointer" }}>
                  <input type="radio" name="member-role" value={r.value} checked={roleValue === r.value} onChange={() => setRoleValue(r.value)}/>
                  <div>
                    <div style={{ fontWeight: 500, fontSize: 13 }}>{r.label}</div>
                    <div className="text-xs muted">{r.desc}</div>
                  </div>
                </label>
              ))}
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}

const INVITE_ROLES = [
  { label: "Admin",  value: "admin",  desc: "Full access including user management" },
  { label: "Member", value: "member", desc: "Can create and manage issues and pages" },
  { label: "Viewer", value: "viewer", desc: "Read-only access to all content" },
];

function AddMemberModal({ open, onClose, onSuccess }) {
  const [email, setEmail]       = useState("");
  const [role, setRole]         = useState("member");
  const [loading, setLoading]   = useState(false);
  const [error, setError]       = useState(null);
  const [invite, setInvite]     = useState(null);
  const [copied, setCopied]     = useState(false);

  useEffect(() => {
    if (open) { setEmail(""); setRole("member"); setLoading(false); setError(null); setInvite(null); setCopied(false); }
  }, [open]);

  async function handleSend() {
    if (!email) return;
    setError(null);
    setLoading(true);
    try {
      const result = await api("/invites", { body: { email, role } });
      setInvite(result);
      onSuccess(result);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  function handleCopy() {
    if (invite?.invite_url) {
      navigator.clipboard.writeText(invite.invite_url).then(() => {
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      });
    }
  }

  if (invite) {
    return (
      <Modal open={open} onClose={onClose} title="Invite created"
        footer={<Button variant="primary" onClick={onClose}>Done</Button>}
      >
        <div style={{ padding: "4px 0" }}>
          <div style={{ fontSize: 14, marginBottom: 16, color: "var(--text-muted)" }}>
            Invite link for <strong>{invite.email}</strong> ({invite.role}) — expires in 7 days.
            {" "}Copy and send it manually, or wait for the email if SMTP is configured.
          </div>
          <div style={{ display: "flex", gap: 8, alignItems: "stretch" }}>
            <input
              readOnly
              value={invite.invite_url || ""}
              className="input"
              style={{ flex: 1, fontFamily: "var(--font-mono, monospace)", fontSize: 12 }}
              onClick={(e) => e.target.select()}
            />
            <button
              type="button"
              onClick={handleCopy}
              style={{
                padding: "0 16px", borderRadius: 7, border: "1px solid var(--border)",
                background: copied ? "var(--green-50, #F0FDF4)" : "var(--bg)",
                color: copied ? "var(--green-700, #15803D)" : "var(--text)",
                cursor: "pointer", fontSize: 13, fontWeight: 500, whiteSpace: "nowrap",
              }}
            >
              {copied ? "Copied!" : "Copy link"}
            </button>
          </div>
        </div>
      </Modal>
    );
  }

  return (
    <Modal open={open} onClose={onClose} title="Invite a new member"
      footer={
        <>
          <Button onClick={onClose}>Cancel</Button>
          <Button variant="primary" disabled={!email || loading} onClick={handleSend}>
            {loading ? "Sending…" : "Create invite link"}
          </Button>
        </>
      }
    >
      {error && (
        <div style={{ padding: "10px 14px", borderRadius: 8, background: "#FEF2F2", color: "#991B1B", fontSize: 13, marginBottom: 14, border: "1px solid #FECACA" }}>
          {error}
        </div>
      )}
      <div style={{ marginBottom: 14 }}>
        <label className="label">Email address</label>
        <input className="input" type="email" placeholder="colleague@company.com" value={email}
          onChange={(e) => setEmail(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && !loading && handleSend()}/>
      </div>
      <div>
        <label className="label">Role</label>
        <div style={{ display: "grid", gap: 8, marginTop: 4 }}>
          {INVITE_ROLES.map((r) => (
            <label key={r.value} style={{ display: "flex", alignItems: "flex-start", gap: 10, padding: "10px 12px", borderRadius: 8, border: "1px solid " + (role === r.value ? "var(--indigo-400)" : "var(--border)"), background: role === r.value ? "var(--indigo-50, #EEF2FF)" : "var(--bg)", cursor: "pointer" }}>
              <input type="radio" name="invite-role" value={r.value} checked={role === r.value}
                onChange={() => setRole(r.value)} style={{ marginTop: 2 }}/>
              <div>
                <div style={{ fontWeight: 500, fontSize: 13, color: "var(--text)" }}>{r.label}</div>
                <div style={{ fontSize: 12, color: "var(--text-muted)", marginTop: 1 }}>{r.desc}</div>
              </div>
            </label>
          ))}
        </div>
      </div>
    </Modal>
  );
}

function MemberProfileDrawer({ open, userId, onClose, people }) {
  const u = (people || []).find((p) => p.id === userId);
  if (!u) return null;
  return (
    <Modal open={open} onClose={onClose} title={u.name} size="lg">
      <div style={{ display: "grid", gridTemplateColumns: "180px 1fr", gap: 24 }}>
        <div style={{ textAlign: "center" }}>
          <Avatar user={u} size="xl" style={{ margin: "0 auto" }}/>
          <div className="bold" style={{ marginTop: 12, fontSize: 16 }}>{u.name}</div>
          <div className="text-xs muted" style={{ marginBottom: 8 }}>{u.email}</div>
          <Badge tone={(ROLE_META[u.role] || ROLE_META.member).tone}>{roleLabel(u.role)}</Badge>
          <div className="row gap-2" style={{ justifyContent: "center", marginTop: 16 }}>
            <Button data-size="sm" icon="mail">Email</Button>
          </div>
        </div>
        <div>
          <h4 style={{ margin: "0 0 10px", fontSize: 13, fontWeight: 600 }}>Basic info</h4>
          <dl style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8, fontSize: 13, margin: 0 }}>
            <dt className="muted">Role</dt><dd>{roleLabel(u.role)}</dd>
            <dt className="muted">Status</dt><dd><Badge tone={u.status === "Active" ? "success" : u.status === "Pending" ? "warning" : "muted"} dot>{u.status}</Badge></dd>
            <dt className="muted">Joined</dt><dd>{u.joined}</dd>
          </dl>
        </div>
      </div>
    </Modal>
  );
}

// ─── Notification Settings ────────────────────────────────
function EmailPrefs() {
  const [prefs, setPrefs] = useState(null);
  const [saving, setSaving] = useState(false);
  const toast = useToast();
  const timer = useRef(null);

  useEffect(() => {
    let live = true;
    api("/notifications/preferences").then((d) => { if (live) setPrefs(d); }).catch(() => {});
    return () => { live = false; clearTimeout(timer.current); };
  }, []);

  function toggle(key, val) {
    setPrefs((p) => ({ ...p, [key]: val }));
    clearTimeout(timer.current);
    setSaving(true);
    timer.current = setTimeout(async () => {
      try { await api("/notifications/preferences", { method: "PUT", body: { [key]: val } }); toast && toast("Email preferences saved"); }
      catch (e) { toast && toast(e.message, { icon: "x", color: "#F87171" }); }
      setSaving(false);
    }, 600);
  }

  return (
    <div className="card">
      <div className="card-head">
        <h3 className="row gap-2"><Icon name="mail" size={15} color="var(--info)"/> Email notifications</h3>
        {saving ? <span className="text-xs muted row gap-2"><MiniSpinner size={12}/> Saving…</span> : prefs && <span className="text-xs muted">Saved automatically</span>}
      </div>
      <div>
        {prefs === null ? (
          [0, 1, 2, 3, 4].map((i) => (
            <div key={i} className="row gap-4" style={{ padding: "14px 20px", borderBottom: i < 4 ? "1px solid var(--border)" : 0 }}>
              <div style={{ flex: 1 }}><div className="skel" style={{ height: 12, width: "40%", marginBottom: 6 }}/><div className="skel" style={{ height: 10, width: "70%" }}/></div>
              <div className="skel" style={{ width: 32, height: 18, borderRadius: 99 }}/>
            </div>
          ))
        ) : EMAIL_PREFS.map(([key, label, desc], i) => (
          <div key={key} className="row gap-4" style={{ padding: "14px 20px", borderBottom: i < EMAIL_PREFS.length - 1 ? "1px solid var(--border)" : 0 }}>
            <div style={{ flex: 1 }}>
              <div className="bold text-sm">{label}</div>
              <div className="text-xs muted" style={{ marginTop: 2 }}>{desc}</div>
            </div>
            <Switch on={!!prefs[key]} onChange={(v) => toggle(key, v)}/>
          </div>
        ))}
      </div>
    </div>
  );
}

export function NotificationSettingsView({ nav }) {
  const { me } = useApp();
  const [prefs, setPrefs]     = useState(() => Object.fromEntries(NOTIF_EVENTS.map((e) => [e.key, e.default])));
  const [channel, setChannel] = useState({ email: true, inapp: true, telegram: true });
  const [digestTime, setDigestTime] = useState("09:00");
  const [savingDigest, setSavingDigest] = useState(false);
  const [tgStatus, setTgStatus]     = useState(null); // null=loading, false=not connected, object=connected
  const [savingCh, setSavingCh]     = useState(false);
  const toast = useToast();

  // Load telegram connection status + notification preferences
  useEffect(() => {
    api("/auth/telegram/status").then((d) => setTgStatus(d || false)).catch(() => setTgStatus(false));
    api("/notifications/preferences").then((d) => {
      if (d) {
        setChannel((c) => ({ ...c, telegram: d.telegram_enabled !== false }));
        if (d.digest_time) setDigestTime(d.digest_time);
      }
    }).catch(() => {});
  }, []);

  async function saveDigestTime(time) {
    setSavingDigest(true);
    try {
      await api("/notifications/preferences", { method: "PUT", body: { digest_time: time } });
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSavingDigest(false); }
  }

  async function saveChannel(key, val) {
    const next = { ...channel, [key]: val };
    setChannel(next);
    setSavingCh(true);
    try {
      await api("/notifications/preferences", { method: "PUT", body: { telegram_enabled: next.telegram } });
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
    finally { setSavingCh(false); }
  }

  const tgConnected = tgStatus && tgStatus.connected;
  const tgSub = tgConnected
    ? (tgStatus.username ? "@" + tgStatus.username : "Connected")
    : "Not connected";

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">You <Icon name="chevronRight" size={11}/> <span>Notification preferences</span></div>
          <h1>Notification preferences</h1>
          <p>Control what reaches you, on which channel, and when.</p>
        </div>
        <div className="row gap-2">
          {savingCh && <span className="text-xs muted row gap-2"><MiniSpinner size={12}/> Saving…</span>}
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", display: "grid", gridTemplateColumns: "1fr 360px", gap: 24, maxWidth: 1280 }}>
        <div>
          <div className="card" style={{ marginBottom: 16 }}>
            <div className="card-head">
              <h3>Channels</h3>
            </div>
            <div style={{ padding: 16, display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 12 }}>
              <ChannelToggle icon="mail" label="Email" sub={me ? me.email : "—"} on={channel.email} onChange={(v) => setChannel({ ...channel, email: v })}/>
              <ChannelToggle icon="bell" label="In-app" sub="Always on for mentions" on={channel.inapp} onChange={(v) => setChannel({ ...channel, inapp: v })}/>
              <ChannelToggle
                icon="telegram" iconColor="linear-gradient(135deg,#2AABEE,#229ED9)"
                label="Telegram" sub={tgSub}
                on={channel.telegram && !!tgConnected}
                disabled={!tgConnected}
                onChange={(v) => saveChannel("telegram", v)}
                onClickDisabled={() => nav("integrations")}
              />
            </div>
            {!tgConnected && (
              <div style={{ padding: "0 16px 14px", display: "flex", alignItems: "center", gap: 8 }}>
                <Icon name="telegram" size={13} color="#2AABEE"/>
                <span className="text-xs muted">Telegram not connected —</span>
                <button className="text-xs" style={{ color: "var(--indigo-600)", background: "none", border: "none", cursor: "pointer", padding: 0 }}
                  onClick={() => nav("integrations")}>Connect now</button>
              </div>
            )}
          </div>
          <EmailPrefs/>
          <div className="card" style={{ marginTop: 16 }}>
            <div className="card-head">
              <h3>Notify me about</h3>
              <Menu align="right" trigger={<Button variant="ghost" data-size="sm">Bulk actions <Icon name="chevronDown" size={12}/></Button>} items={[
                { label: "Turn all on",  onClick: () => setPrefs(Object.fromEntries(NOTIF_EVENTS.map((e) => [e.key, true]))) },
                { label: "Turn all off", onClick: () => setPrefs(Object.fromEntries(NOTIF_EVENTS.map((e) => [e.key, false]))) },
              ]}/>
            </div>
            <div>
              {NOTIF_EVENTS.map((ev, i) => (
                <div key={ev.key} className="row gap-4" style={{ padding: "14px 20px", borderBottom: i < NOTIF_EVENTS.length - 1 ? "1px solid var(--border)" : 0 }}>
                  <div style={{ flex: 1 }}>
                    <div className="bold text-sm">{ev.label}</div>
                    <div className="text-xs muted" style={{ marginTop: 2 }}>{ev.desc}</div>
                    {ev.key === "digest" && prefs.digest && (
                      <div className="row gap-2 text-xs" style={{ marginTop: 8 }}>
                        <span className="muted">Send at</span>
                        <input className="input" type="time" value={digestTime}
                          onChange={(e) => setDigestTime(e.target.value)}
                          onBlur={(e) => saveDigestTime(e.target.value)}
                          style={{ width: 110, padding: "3px 8px" }}/>
                        {savingDigest && <span className="text-xs muted">Saving…</span>}
                        <span className="muted">your local time</span>
                      </div>
                    )}
                  </div>
                  <Switch on={prefs[ev.key]} onChange={(v) => setPrefs({ ...prefs, [ev.key]: v })}/>
                </div>
              ))}
            </div>
          </div>
        </div>
        <div/>
      </div>
    </div>
  );
}

function ChannelToggle({ icon, iconColor, label, sub, on, onChange, disabled, onClickDisabled }) {
  return (
    <div onClick={disabled && onClickDisabled ? onClickDisabled : undefined}
      style={{ padding: 12, border: "1px solid " + (on ? "var(--indigo-500)" : "var(--border)"), borderRadius: 8,
        background: on ? "var(--indigo-50)" : "var(--bg)", opacity: disabled ? .55 : 1,
        cursor: disabled && onClickDisabled ? "pointer" : "default" }}>
      <div className="row gap-3" style={{ marginBottom: 8 }}>
        <div style={{ width: 28, height: 28, borderRadius: 7, background: iconColor || "var(--indigo-600)", color: "#fff", display: "grid", placeItems: "center" }}>
          <Icon name={icon} size={15}/>
        </div>
        <div className="stack" style={{ lineHeight: 1.3, flex: 1, minWidth: 0 }}>
          <span className="bold text-sm">{label}</span>
          <span className="text-xs muted" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{sub}</span>
        </div>
        <Switch on={on} onChange={disabled ? undefined : onChange}/>
      </div>
    </div>
  );
}

// ─── Integrations page ────────────────────────────────────
export function IntegrationsView({ nav }) {
  const [code, setCode]           = useState("");
  const [status, setStatus]       = useState(null);   // { connected, username, verified_at }
  const [loading, setLoading]     = useState(true);
  const [verifying, setVerifying] = useState(false);
  const toast = useToast();

  async function loadStatus() {
    try {
      const data = await api("/auth/telegram/status");
      setStatus(data);
    } catch { setStatus({ connected: false }); }
    finally { setLoading(false); }
  }
  useEffect(() => { loadStatus(); }, []);

  async function verify() {
    if (code.length !== 6) return;
    setVerifying(true);
    try {
      await api("/auth/telegram/verify", { method: "POST", body: { code } });
      toast("✅ Telegram connected!", { icon: "telegram", color: "#2AABEE" });
      setCode("");
      loadStatus();
    } catch (e) {
      toast(e.message || "Invalid or expired code", { icon: "x", color: "#F87171" });
    } finally { setVerifying(false); }
  }

  async function disconnect() {
    try {
      await api("/auth/telegram/disconnect", { method: "DELETE" });
      toast("Telegram disconnected");
      setStatus({ connected: false });
    } catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
  }

  const connected = status?.connected;

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><Icon name="settings" size={11}/> <span>Integrations</span></div>
          <h1>Integrations</h1>
          <p>Connect JiraFlow to the tools your team already uses.</p>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", maxWidth: 1100 }}>
        <div className="card" style={{ padding: 0, marginBottom: 24, borderLeft: "4px solid var(--tg)", overflow: "hidden" }}>
          <div style={{ padding: 24, display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 20, alignItems: "center" }}>
            <div style={{ width: 64, height: 64, borderRadius: 14, background: "linear-gradient(135deg, #2AABEE, #229ED9)", color: "#fff", display: "grid", placeItems: "center", boxShadow: "0 4px 12px rgba(42,171,238,.3)" }}>
              <Icon name="telegram" size={32}/>
            </div>
            <div>
              <div className="row gap-2" style={{ marginBottom: 4 }}>
                <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Telegram bot</h2>
                {loading
                  ? <Badge tone="muted" dot>Checking…</Badge>
                  : <Badge tone={connected ? "success" : "muted"} dot>{connected ? "Active" : "Not connected"}</Badge>}
              </div>
              <p className="secondary" style={{ margin: 0, fontSize: 13.5 }}>
                Instant push notifications via Telegram. Get issue assignments, mentions, and sprint events straight to your phone.
              </p>
              {connected && status?.username && (
                <p style={{ margin: "6px 0 0", fontSize: 13, color: "var(--tg)", fontWeight: 500 }}>
                  @{status.username}
                </p>
              )}
            </div>
            {connected && (
              <Button variant="secondary" icon="x" onClick={disconnect}>Disconnect</Button>
            )}
          </div>

          {!connected && (
            <div style={{ borderTop: "1px solid var(--border)", padding: "20px 24px", background: "var(--bg-subtle)" }}>
              <h4 style={{ margin: "0 0 12px", fontSize: 13, fontWeight: 600, color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: ".04em" }}>How to connect</h4>
              <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 14 }}>
                <Step n={1} title="Open the bot" desc={
                  <><a href="https://t.me/jira_flowbot" target="_blank" rel="noreferrer" style={{ color: "var(--tg)", fontWeight: 500 }}>@jira_flowbot</a>{" "}ni Telegram'da oching va <span className="mono bold">/start</span> yuboring.</>
                }/>
                <Step n={2} title="Kodni oling" desc="Bot sizga 6 xonali tasdiqlash kodini yuboradi. Kod 10 daqiqa amal qiladi."/>
                <Step n={3} title="Kodni kiriting" desc="Botdan olgan kodni quyida kiriting:" cta={
                  <div className="row gap-2" style={{ marginTop: 4 }}>
                    <input
                      className="input mono"
                      maxLength="6"
                      placeholder="000000"
                      value={code}
                      onChange={(e) => setCode(e.target.value.replace(/\D/g, ""))}
                      onKeyDown={(e) => e.key === "Enter" && verify()}
                      style={{ width: 110, letterSpacing: ".2em", textAlign: "center", fontSize: 18 }}
                    />
                    <Button variant="primary" data-size="sm" onClick={verify} disabled={code.length !== 6 || verifying}>
                      {verifying ? "…" : "Verify"}
                    </Button>
                  </div>
                }/>
              </div>
            </div>
          )}
        </div>

        <h3 style={{ fontSize: 14, fontWeight: 600, margin: "0 0 12px" }}>More integrations</h3>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 12 }}>
          {[
            { name: "GitHub",    desc: "Link PRs and commits to issues automatically.", icon: "branch",  status: "Active", tone: "success" },
            { name: "Slack",     desc: "Slash commands and channel digests.",            icon: "comment", status: "Active", tone: "success" },
            { name: "Datadog",   desc: "Create issues from alerts.",                    icon: "chart",   status: "Setup",  tone: "warning" },
            { name: "PagerDuty", desc: "Sync on-call rotations & incidents.",           icon: "bell",    status: "Setup",  tone: "warning" },
            { name: "Linear",    desc: "Migrate from your previous tracker.",            icon: "list",    status: "—",      tone: "muted" },
            { name: "Webhook",   desc: "Stream all events to your URL.",                 icon: "code",    status: "—",      tone: "muted" },
          ].map((x) => (
            <div key={x.name} className="card card-pad">
              <div className="row gap-3" style={{ marginBottom: 8 }}>
                <div style={{ width: 36, height: 36, borderRadius: 8, background: "var(--bg-subtle)", border: "1px solid var(--border)", display: "grid", placeItems: "center", color: "var(--text-secondary)" }}>
                  <Icon name={x.icon} size={16}/>
                </div>
                <div className="stack" style={{ lineHeight: 1.25, flex: 1 }}>
                  <span className="bold">{x.name}</span>
                  <Badge tone={x.tone} dot style={{ alignSelf: "flex-start", marginTop: 1 }}>{x.status}</Badge>
                </div>
              </div>
              <div className="text-sm secondary" style={{ minHeight: 32, marginBottom: 10 }}>{x.desc}</div>
              <Button data-size="sm" style={{ width: "100%", justifyContent: "center" }}>Configure</Button>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function Step({ n, title, desc, cta }) {
  return (
    <div style={{ padding: 12, background: "var(--bg)", border: "1px solid var(--border)", borderRadius: 8 }}>
      <div className="row gap-2" style={{ marginBottom: 6 }}>
        <span style={{ width: 20, height: 20, borderRadius: "50%", background: "var(--tg)", color: "#fff", display: "grid", placeItems: "center", fontSize: 11, fontWeight: 700 }}>{n}</span>
        <span className="bold text-sm">{title}</span>
      </div>
      <div className="text-sm secondary" style={{ marginBottom: cta ? 8 : 0 }}>{desc}</div>
      {cta}
    </div>
  );
}
