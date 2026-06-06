// views_members.jsx — Members table, add modal, member profile, integrations, notification settings

function MembersView({ nav, tg, people, setPeople }) {
  const [openAdd, setOpenAdd] = React.useState(false);
  const [openProfile, setOpenProfile] = React.useState(null); // user id
  const [search, setSearch] = React.useState("");
  const [filter, setFilter] = React.useState("all"); // role
  const toast = useToast();

  const list = people.filter((p) => {
    if (search && !(p.name.toLowerCase().includes(search.toLowerCase()) || p.email.toLowerCase().includes(search.toLowerCase()))) return false;
    if (filter !== "all" && p.role !== filter) return false;
    return true;
  });

  const connected = people.filter((p) => p.tg).length;

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/board">Core Infrastructure</a> <Icon name="chevronRight" size={11}/> <span>Members</span></div>
          <h1>Members</h1>
          <p>Manage who can access this project and how they receive notifications.</p>
        </div>
        <div className="row gap-2">
          <Button icon="download">Export</Button>
          <Button variant="primary" icon="plus" onClick={() => setOpenAdd(true)}>Invite member</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px" }}>
        {/* Telegram callout strip */}
        <div className="tg-card" style={{ padding: 14, marginBottom: 16, display: "flex", alignItems: "center", gap: 14 }}>
          <div style={{ width: 40, height: 40, borderRadius: 10, background: "var(--tg-bg)", display: "grid", placeItems: "center", color: "var(--tg)" }}>
            <Icon name="telegram" size={20}/>
          </div>
          <div style={{ flex: 1 }}>
            <div className="bold text-sm">Telegram notifications are {tg ? "live" : "not connected"}</div>
            <div className="text-xs secondary">
              {connected} of {people.length} members are connected to <span style={{ color: "var(--tg)", fontWeight: 500 }}>@forge_team_bot</span>. {!tg && "Turn on the bot in Tweaks panel to see the live state."}
            </div>
          </div>
          <Button variant="ghost" data-size="sm" icon="telegram" onClick={() => nav("integrations")}>Configure</Button>
        </div>

        <div className="card" style={{ overflow: "hidden" }}>
          <div className="row" style={{ padding: 12, gap: 8, borderBottom: "1px solid var(--border)" }}>
            <div className="search" style={{ width: 280, padding: "4px 10px" }}>
              <Icon name="search" size={13}/>
              <input placeholder="Search by name or email…" value={search} onChange={(e) => setSearch(e.target.value)}/>
            </div>
            <Pill label="Role" icon="shield" value={filter === "all" ? "All" : filter} options={[{ id: "all", label: "All roles" }, ...Object.keys(FORGE_DATA.ROLE_META).map((r) => ({ id: r, label: r }))]} onChange={setFilter}/>
            <Pill label="Status" icon="check" value="All" options={[{ id: "all", label: "All" }, { id: "Active", label: "Active" }, { id: "Pending", label: "Pending" }, { id: "Inactive", label: "Inactive" }]} onChange={() => {}}/>
            <Pill label="Telegram" icon="telegram" value="All" options={[{ id: "all", label: "All" }, { id: "y", label: "Connected" }, { id: "n", label: "Not connected" }]} onChange={() => {}}/>
            <div style={{ flex: 1 }}/>
            <span className="text-xs muted">{list.length} members</span>
          </div>
          <table className="table">
            <thead>
              <tr>
                <th>Member</th>
                <th style={{ width: 130 }}>Role</th>
                <th style={{ width: 220 }}>Telegram</th>
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
                  <td><Badge tone={FORGE_DATA.ROLE_META[u.role].tone}>{u.role}</Badge></td>
                  <td>
                    {u.tg && tg ? (
                      <div className="row gap-2">
                        <Icon name="telegram" size={14} color="#2AABEE"/>
                        <span className="text-sm">{u.tg}</span>
                        <span className="text-xs muted mono">·{u.tgId.slice(-6)}</span>
                      </div>
                    ) : u.tg && !tg ? (
                      <div className="row gap-2">
                        <Icon name="telegram" size={14} color="var(--text-muted)"/>
                        <span className="text-sm secondary">{u.tg}</span>
                        <Badge tone="warning" style={{ fontSize: 10 }}>bot off</Badge>
                      </div>
                    ) : (
                      <span className="row gap-2 text-sm muted">
                        <Icon name="telegram" size={14}/> Not connected
                        <Button variant="ghost" data-size="sm" style={{ padding: "0 6px", color: "var(--tg)" }} onClick={(e) => { e.stopPropagation(); toast("Invite sent to " + u.name); }}>Invite</Button>
                      </span>
                    )}
                  </td>
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
                        { icon: "user", label: "View profile", onClick: () => setOpenProfile(u.id) },
                        { icon: "telegram", label: "Send test message" },
                        { icon: "shield", label: "Change role" },
                        { divider: true },
                        { icon: "trash", label: "Remove from project", danger: true, onClick: () => setPeople((p) => p.filter((x) => x.id !== u.id)) },
                      ]}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <AddMemberModal open={openAdd} onClose={() => setOpenAdd(false)} onAdd={(m) => {
        setPeople((p) => [...p, { ...m, id: "u" + (p.length + 1), initials: m.name.split(" ").map((s) => s[0]).slice(0, 2).join("").toUpperCase(), color: FORGE_DATA.COLORS[p.length % FORGE_DATA.COLORS.length], status: "Pending", joined: "Just now" }]);
        toast("Invitation sent to " + m.email);
      }}/>
      <MemberProfileDrawer open={!!openProfile} userId={openProfile} onClose={() => setOpenProfile(null)} tg={tg}/>
    </div>
  );
}

// ─── Add member modal ───────────────────────────────────
function AddMemberModal({ open, onClose, onAdd }) {
  const [form, setForm] = React.useState({ name: "", email: "", role: "Developer", tg: "", tgId: "", sendInvite: true });
  React.useEffect(() => { if (open) setForm({ name: "", email: "", role: "Developer", tg: "", tgId: "", sendInvite: true }); }, [open]);

  return (
    <Modal open={open} onClose={onClose} title="Invite a new member"
      footer={
        <>
          <Button onClick={onClose}>Cancel</Button>
          <Button variant="primary" disabled={!form.email} onClick={() => { onAdd(form); onClose(); }}>Send invite</Button>
        </>
      }
    >
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12, marginBottom: 12 }}>
        <div>
          <label className="label">Full name</label>
          <input className="input" placeholder="Mei Yamazaki" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })}/>
        </div>
        <div>
          <label className="label">Email</label>
          <input className="input" type="email" placeholder="mei@forge.dev" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })}/>
        </div>
      </div>
      <div style={{ marginBottom: 12 }}>
        <label className="label">Role</label>
        <div className="row gap-2">
          {Object.keys(FORGE_DATA.ROLE_META).map((r) => (
            <button key={r} onClick={() => setForm({ ...form, role: r })}
              className="btn"
              style={{
                border: "1px solid " + (form.role === r ? "var(--indigo-600)" : "var(--border)"),
                background: form.role === r ? "var(--indigo-50)" : "var(--bg)",
                color: form.role === r ? "var(--indigo-700)" : "var(--text)",
                flex: 1, justifyContent: "center"
              }}>
              {r}
            </button>
          ))}
        </div>
        <div className="help">
          {form.role === "Admin" && "Full access to all settings, members, and projects."}
          {form.role === "Manager" && "Can manage sprints, members, and issues."}
          {form.role === "Developer" && "Can create and edit issues, push commits, edit wiki."}
          {form.role === "Viewer" && "Read-only access to boards and docs."}
        </div>
      </div>

      {/* Telegram section */}
      <div style={{ borderTop: "1px solid var(--border)", paddingTop: 14, marginTop: 14 }}>
        <div className="row gap-2" style={{ marginBottom: 4 }}>
          <Icon name="telegram" size={16} color="#2AABEE"/>
          <span className="bold">Telegram (optional)</span>
        </div>
        <div className="help" style={{ marginBottom: 12, padding: 10, background: "var(--tg-bg)", borderRadius: 6, color: "var(--text)", display: "flex", gap: 8, alignItems: "flex-start" }}>
          <Icon name="bell" size={13} color="#2AABEE" style={{ flexShrink: 0, marginTop: 2 }}/>
          <span>Members receive instant notifications via our Telegram bot. They can complete this step themselves after accepting the invite.</span>
        </div>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
          <div>
            <label className="label">Telegram username</label>
            <input className="input" placeholder="@username" value={form.tg} onChange={(e) => setForm({ ...form, tg: e.target.value })}/>
          </div>
          <div>
            <label className="label">Telegram ID</label>
            <input className="input" placeholder="847291043" value={form.tgId} onChange={(e) => setForm({ ...form, tgId: e.target.value })}/>
            <div className="help">Numeric ID from @userinfobot</div>
          </div>
        </div>
      </div>

      <label className="row gap-2" style={{ marginTop: 16, cursor: "default" }}>
        <Switch on={form.sendInvite} onChange={(v) => setForm({ ...form, sendInvite: v })}/>
        <div>
          <div className="text-sm medium">Send invite email and Telegram link</div>
          <div className="help">They'll get a one-click link to join the workspace and connect their Telegram.</div>
        </div>
      </label>
    </Modal>
  );
}

// ─── Member profile drawer ──────────────────────────────
function MemberProfileDrawer({ open, userId, onClose, tg }) {
  const u = FORGE_DATA.PEOPLE.find((p) => p.id === userId);
  if (!u) return null;
  return (
    <Modal open={open} onClose={onClose} title={u.name} size="lg">
      <div style={{ display: "grid", gridTemplateColumns: "180px 1fr", gap: 24 }}>
        <div style={{ textAlign: "center" }}>
          <Avatar user={u} size="xl" style={{ margin: "0 auto" }}/>
          <div className="bold" style={{ marginTop: 12, fontSize: 16 }}>{u.name}</div>
          <div className="text-xs muted" style={{ marginBottom: 8 }}>{u.email}</div>
          <Badge tone={FORGE_DATA.ROLE_META[u.role].tone}>{u.role}</Badge>
          <div className="row gap-2" style={{ justifyContent: "center", marginTop: 16 }}>
            <Button data-size="sm" icon="mail">Email</Button>
            {u.tg && <Button data-size="sm" icon="telegram" style={{ color: "var(--tg)" }}>Message</Button>}
          </div>
        </div>
        <div>
          <h4 style={{ margin: "0 0 10px", fontSize: 13, fontWeight: 600 }}>Basic info</h4>
          <dl style={{ display: "grid", gridTemplateColumns: "140px 1fr", gap: 8, fontSize: 13, margin: 0 }}>
            <dt className="muted">Role</dt><dd>{u.role}</dd>
            <dt className="muted">Status</dt><dd><Badge tone={u.status === "Active" ? "success" : u.status === "Pending" ? "warning" : "muted"} dot>{u.status}</Badge></dd>
            <dt className="muted">Joined</dt><dd>{u.joined}</dd>
            <dt className="muted">Local time</dt><dd>14:32 GMT+1</dd>
            <dt className="muted">Projects</dt><dd>Core Infrastructure, Observability</dd>
          </dl>

          <h4 style={{ margin: "20px 0 10px", fontSize: 13, fontWeight: 600, display: "flex", alignItems: "center", gap: 8 }}>
            <Icon name="telegram" size={14} color="#2AABEE"/> Telegram integration
          </h4>
          {u.tg ? (
            <div className="tg-card" style={{ padding: 12 }}>
              <div className="row gap-3">
                <div style={{ width: 36, height: 36, borderRadius: "50%", background: "var(--tg)", color: "#fff", display: "grid", placeItems: "center" }}>
                  <Icon name="telegram" size={18}/>
                </div>
                <div className="stack" style={{ lineHeight: 1.3, flex: 1 }}>
                  <span className="bold text-sm">{u.tg}</span>
                  <span className="text-xs muted mono">ID {u.tgId}</span>
                </div>
                <Badge tone={tg ? "success" : "warning"} dot>{tg ? "Connected" : "Bot off"}</Badge>
              </div>
              <div className="text-xs muted" style={{ marginTop: 8 }}>Last notification delivered 12 minutes ago</div>
            </div>
          ) : (
            <div style={{ padding: 12, background: "var(--bg-subtle)", border: "1px dashed var(--border)", borderRadius: 8 }}>
              <div className="text-sm secondary">Not connected to Telegram bot.</div>
              <Button variant="telegram" data-size="sm" icon="telegram" style={{ marginTop: 8 }}>Send connect link</Button>
            </div>
          )}

          <h4 style={{ margin: "20px 0 10px", fontSize: 13, fontWeight: 600 }}>Recent activity</h4>
          <div style={{ fontSize: 13 }}>
            {[
              "resolved INFRA-200 — 2d ago",
              "commented on INFRA-232 — 28m ago",
              "moved INFRA-220 to In Progress — 12m ago",
            ].map((t, i) => (
              <div key={i} style={{ padding: "6px 0", borderBottom: i < 2 ? "1px solid var(--border)" : 0 }} className="secondary">{t}</div>
            ))}
          </div>
        </div>
      </div>
    </Modal>
  );
}

// ─── Notification settings page ─────────────────────────
// Feature 23: email preferences backed by /notifications/preferences
const EMAIL_PREFS = [
  ["email_assigned", "Issue assigned to me", "When someone assigns an issue to you."],
  ["email_mentioned", "Mentioned in comment", "Someone @mentions you in a comment."],
  ["email_commented", "Comment on watched issue", "New comments on issues you watch."],
  ["email_status", "Issue status changed", "Status changes on issues you watch."],
  ["email_watcher", "Updates on watched items", "Any update to items you watch."],
];

function EmailPrefs() {
  const [prefs, setPrefs] = React.useState(null);
  const [saving, setSaving] = React.useState(false);
  const toast = useToast();
  const timer = React.useRef(null);

  React.useEffect(() => {
    let live = true;
    api("/notifications/preferences").then((d) => { if (live) setPrefs(d); }).catch((e) => toast(e.message, { icon: "x", color: "#F87171" }));
    return () => { live = false; clearTimeout(timer.current); };
  }, []);

  function toggle(key, val) {
    setPrefs((p) => ({ ...p, [key]: val }));
    clearTimeout(timer.current);
    setSaving(true);
    timer.current = setTimeout(async () => {
      try { await api("/notifications/preferences", { method: "PUT", body: { [key]: val } }); toast("Email preferences saved"); }
      catch (e) { toast(e.message, { icon: "x", color: "#F87171" }); }
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

function NotificationSettingsView({ nav, tg }) {
  const [prefs, setPrefs] = React.useState(() => Object.fromEntries(FORGE_DATA.NOTIF_EVENTS.map((e) => [e.key, e.default])));
  const [channel, setChannel] = React.useState({ tg: true, email: true, inapp: true });
  const [digestTime, setDigestTime] = React.useState("09:00");
  const toast = useToast();

  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs">You <Icon name="chevronRight" size={11}/> <span>Notification preferences</span></div>
          <h1>Notification preferences</h1>
          <p>Control what reaches you, on which channel, and when.</p>
        </div>
        <div className="row gap-2">
          <Button onClick={() => toast("Settings saved")}>Save changes</Button>
          <Button variant="primary" icon="send" onClick={() => toast("Test notification sent to " + (FORGE_DATA.ME.tg || "email"), { icon: "telegram", color: "#2AABEE" })}>Send test notification</Button>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", display: "grid", gridTemplateColumns: "1fr 360px", gap: 24, maxWidth: 1280 }}>
        <div>
          {/* Channels */}
          <div className="card" style={{ marginBottom: 16 }}>
            <div className="card-head">
              <h3>Channels</h3>
              <span className="text-xs muted">Pick where each event below is delivered</span>
            </div>
            <div style={{ padding: 16, display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 12 }}>
              <ChannelToggle icon="telegram" iconColor="#2AABEE" label="Telegram" sub={tg && FORGE_DATA.ME.tg ? FORGE_DATA.ME.tg : "Not connected"} on={channel.tg && tg} onChange={(v) => setChannel({ ...channel, tg: v })} disabled={!tg}/>
              <ChannelToggle icon="mail" label="Email" sub="maya@forge.dev" on={channel.email} onChange={(v) => setChannel({ ...channel, email: v })}/>
              <ChannelToggle icon="bell" label="In-app" sub="Always on for mentions" on={channel.inapp} onChange={(v) => setChannel({ ...channel, inapp: v })}/>
            </div>
          </div>

          {/* Feature 23: API-backed email preferences */}
          <EmailPrefs/>

          {/* Events */}
          <div className="card" style={{ marginTop: 16 }}>
            <div className="card-head">
              <h3>Notify me about</h3>
              <Menu align="right" trigger={<Button variant="ghost" data-size="sm">Bulk actions <Icon name="chevronDown" size={12}/></Button>} items={[
                { label: "Turn all on", onClick: () => setPrefs(Object.fromEntries(FORGE_DATA.NOTIF_EVENTS.map((e) => [e.key, true]))) },
                { label: "Turn all off", onClick: () => setPrefs(Object.fromEntries(FORGE_DATA.NOTIF_EVENTS.map((e) => [e.key, false]))) },
              ]}/>
            </div>
            <div>
              {FORGE_DATA.NOTIF_EVENTS.map((ev, i) => (
                <div key={ev.key} className="row gap-4" style={{ padding: "14px 20px", borderBottom: i < FORGE_DATA.NOTIF_EVENTS.length - 1 ? "1px solid var(--border)" : 0 }}>
                  <div style={{ flex: 1 }}>
                    <div className="bold text-sm">{ev.label}</div>
                    <div className="text-xs muted" style={{ marginTop: 2 }}>{ev.desc}</div>
                    {ev.key === "digest" && prefs.digest && (
                      <div className="row gap-2 text-xs" style={{ marginTop: 8 }}>
                        <span className="muted">Send at</span>
                        <input className="input" type="time" value={digestTime} onChange={(e) => setDigestTime(e.target.value)} style={{ width: 110, padding: "3px 8px" }}/>
                        <span className="muted">your local time</span>
                      </div>
                    )}
                  </div>
                  <div className="row gap-3">
                    {channel.tg && tg && <Icon name="telegram" size={14} color={prefs[ev.key] ? "#2AABEE" : "var(--text-muted)"} title="Telegram"/>}
                    {channel.email && <Icon name="mail" size={14} color={prefs[ev.key] ? "var(--info)" : "var(--text-muted)"} title="Email"/>}
                    {channel.inapp && <Icon name="bell" size={14} color={prefs[ev.key] ? "var(--indigo-600)" : "var(--text-muted)"} title="In-app"/>}
                  </div>
                  <Switch on={prefs[ev.key]} onChange={(v) => setPrefs({ ...prefs, [ev.key]: v })}/>
                </div>
              ))}
            </div>
          </div>

          <div style={{ marginTop: 24, padding: 16, background: "var(--bg-subtle)", borderRadius: 10, border: "1px solid var(--border)" }}>
            <h3 style={{ margin: "0 0 4px", fontSize: 14 }}>Quiet hours</h3>
            <p className="text-sm secondary" style={{ margin: "0 0 12px" }}>Mute non-critical notifications during these hours. Mentions and incidents still come through.</p>
            <div className="row gap-3">
              <input className="input" type="time" defaultValue="20:00" style={{ width: 130 }}/>
              <span className="muted">→</span>
              <input className="input" type="time" defaultValue="08:00" style={{ width: 130 }}/>
              <Switch on={false} onChange={() => {}}/>
            </div>
          </div>
        </div>

        {/* Preview panel */}
        <div>
          <div style={{ position: "sticky", top: 0 }}>
            <h3 style={{ fontSize: 13, fontWeight: 600, color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: ".04em", margin: "0 0 12px" }}>Live preview</h3>
            <div style={{
              background: "linear-gradient(180deg, #C7D2FE 0%, #A5B4FC 100%)",
              padding: 24,
              borderRadius: 14,
              minHeight: 480,
              position: "relative",
              boxShadow: "var(--shadow-sm)"
            }}>
              {/* Telegram chat header */}
              <div style={{ background: "rgba(255,255,255,.92)", borderRadius: 10, padding: "10px 12px", marginBottom: 12, display: "flex", alignItems: "center", gap: 10, boxShadow: "var(--shadow-xs)" }}>
                <div style={{ width: 34, height: 34, borderRadius: "50%", background: "#2AABEE", color: "#fff", display: "grid", placeItems: "center", boxShadow: "0 1px 3px rgba(0,0,0,.2)" }}>
                  <Icon name="telegram" size={16}/>
                </div>
                <div className="stack" style={{ lineHeight: 1.2, color: "#0F172A" }}>
                  <span className="bold text-sm">@forge_team_bot</span>
                  <span className="text-xs" style={{ color: "#10B981" }}>● online</span>
                </div>
              </div>

              {/* Message bubbles */}
              <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
                <div className="tg-msg">
                  <div className="bold" style={{ marginBottom: 6, color: "#0F172A" }}>📋 New task assigned to you</div>
                  <div style={{ color: "#334155" }}>
                    <div><span style={{ color: "#64748B" }}>Project:</span> <span className="bold">Core Infrastructure</span></div>
                    <div><span style={{ color: "#64748B" }}>Task:</span> Fix login bug <span className="mono">#INFRA-234</span></div>
                    <div><span style={{ color: "#64748B" }}>Priority:</span> 🔴 Critical</div>
                    <div><span style={{ color: "#64748B" }}>Due:</span> Dec 15, 2024</div>
                  </div>
                  <div style={{ marginTop: 8, paddingTop: 8, borderTop: "1px solid rgba(15,23,42,.08)" }}>
                    <a href="#" style={{ color: "#2AABEE", textDecoration: "none", fontWeight: 500 }}>→ Open in Forge</a>
                  </div>
                  <div style={{ marginTop: 4, color: "#94A3B8", fontSize: 11, textAlign: "right" }}>14:32 ✓✓</div>
                </div>

                <div className="tg-msg">
                  <div className="bold" style={{ marginBottom: 6, color: "#0F172A" }}>💬 Diego mentioned you</div>
                  <div style={{ color: "#334155" }}>
                    <span className="mono text-xs">INFRA-222</span> · Roll out cgroup v2 to nodepools
                  </div>
                  <div style={{ marginTop: 6, padding: "8px 10px", background: "rgba(255,255,255,.5)", borderRadius: 6, fontSize: 12.5, color: "#475569" }}>
                    "@maya can you take the eu-west-1 node group? I'll handle us-east-1."
                  </div>
                  <div style={{ marginTop: 4, color: "#94A3B8", fontSize: 11, textAlign: "right" }}>14:34 ✓✓</div>
                </div>

                <div className="tg-msg">
                  <div className="bold" style={{ marginBottom: 4, color: "#0F172A" }}>🌅 Your morning digest</div>
                  <div style={{ color: "#334155", fontSize: 12.5 }}>
                    7 issues need attention today · 3 PRs awaiting your review · Sprint 24 is 68% complete with 5 days left.
                  </div>
                </div>
              </div>
            </div>
            <p className="text-xs muted" style={{ marginTop: 10, textAlign: "center" }}>
              This is exactly how messages will appear in your Telegram.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

function ChannelToggle({ icon, iconColor, label, sub, on, onChange, disabled }) {
  return (
    <div style={{ padding: 12, border: "1px solid " + (on ? "var(--indigo-500)" : "var(--border)"), borderRadius: 8, background: on ? "var(--indigo-50)" : "var(--bg)", opacity: disabled ? .6 : 1 }}>
      <div className="row gap-3" style={{ marginBottom: 8 }}>
        <div style={{ width: 28, height: 28, borderRadius: 7, background: iconColor || "var(--indigo-600)", color: "#fff", display: "grid", placeItems: "center" }}>
          <Icon name={icon} size={15}/>
        </div>
        <div className="stack" style={{ lineHeight: 1.3, flex: 1, minWidth: 0 }}>
          <span className="bold text-sm">{label}</span>
          <span className="text-xs muted" style={{ overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>{sub}</span>
        </div>
        <Switch on={on} onChange={onChange}/>
      </div>
    </div>
  );
}

// ─── Integrations page ──────────────────────────────────
function IntegrationsView({ nav, tg, setTg }) {
  const [code, setCode] = React.useState("");
  const toast = useToast();
  return (
    <div>
      <div className="page-head">
        <div>
          <div className="crumbs"><a href="#/settings">Settings</a> <Icon name="chevronRight" size={11}/> <span>Integrations</span></div>
          <h1>Integrations</h1>
          <p>Connect Forge to the tools your team already uses.</p>
        </div>
      </div>

      <div style={{ padding: "0 32px 32px", maxWidth: 1100 }}>
        {/* Telegram hero card */}
        <div className="card" style={{ padding: 0, marginBottom: 24, borderLeft: "4px solid var(--tg)", overflow: "hidden" }}>
          <div style={{ padding: 24, display: "grid", gridTemplateColumns: "auto 1fr auto", gap: 20, alignItems: "center" }}>
            <div style={{ width: 64, height: 64, borderRadius: 14, background: "linear-gradient(135deg, #2AABEE, #229ED9)", color: "#fff", display: "grid", placeItems: "center", boxShadow: "0 4px 12px rgba(42,171,238,.3)" }}>
              <Icon name="telegram" size={32}/>
            </div>
            <div>
              <div className="row gap-2" style={{ marginBottom: 4 }}>
                <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Telegram bot</h2>
                <Badge tone={tg ? "success" : "muted"} dot>{tg ? "Active" : "Not connected"}</Badge>
                <Badge tone="tg">Featured</Badge>
              </div>
              <p className="secondary" style={{ margin: 0, fontSize: 13.5 }}>
                Instant push notifications via Telegram. Send issue assignments, mentions, sprint events, and a daily digest — directly to <span className="bold mono" style={{ color: "var(--tg)" }}>@forge_team_bot</span>.
              </p>
              <div className="row gap-4" style={{ marginTop: 12 }}>
                <Stat label="Members connected" value={tg ? "7 / 10" : "0 / 10"}/>
                <Stat label="Sent today" value={tg ? "143" : "0"}/>
                <Stat label="Avg delivery" value={tg ? "0.8s" : "—"}/>
                <Stat label="Bot uptime (30d)" value={tg ? "99.98%" : "—"}/>
              </div>
            </div>
            <div className="stack gap-2">
              <Button variant={tg ? "secondary" : "telegram"} icon="telegram" onClick={() => setTg(!tg)}>
                {tg ? "Disconnect" : "Connect bot"}
              </Button>
              <Button data-size="sm" icon="send" onClick={() => toast("Test message sent to your Telegram", { icon: "telegram", color: "#2AABEE" })}>Send test</Button>
            </div>
          </div>

          {/* Connection guide */}
          <div style={{ borderTop: "1px solid var(--border)", padding: "20px 24px", background: "var(--bg-subtle)" }}>
            <h4 style={{ margin: "0 0 12px", fontSize: 13, fontWeight: 600, color: "var(--text-secondary)", textTransform: "uppercase", letterSpacing: ".04em" }}>How to connect</h4>
            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr 1fr", gap: 14 }}>
              <Step n={1} title="Open the bot" desc={<>Tap <span className="mono" style={{ color: "var(--tg)", fontWeight: 500 }}>@forge_team_bot</span> on Telegram and press Start.</>} cta={<Button data-size="sm" icon="externalLink" onClick={() => toast("Copied bot link to clipboard", { icon: "copy" })}>Open @forge_team_bot</Button>}/>
              <Step n={2} title="Get your Telegram ID" desc="The bot will reply with a 6-digit verification code and your Telegram ID."/>
              <Step n={3} title="Enter the code" desc={
                <div className="row gap-2" style={{ marginTop: 4 }}>
                  <input className="input mono" maxLength="6" placeholder="XXXXXX" value={code} onChange={(e) => setCode(e.target.value.toUpperCase())} style={{ width: 110, letterSpacing: ".2em", textAlign: "center" }}/>
                  <Button variant="primary" data-size="sm" onClick={() => { if (code.length === 6) { setTg(true); toast("Telegram connected"); } }} disabled={code.length !== 6}>Verify</Button>
                </div>
              }/>
            </div>
          </div>
        </div>

        {/* Other integrations grid */}
        <h3 style={{ fontSize: 14, fontWeight: 600, margin: "0 0 12px" }}>More integrations</h3>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 12 }}>
          {[
            { name: "GitHub", desc: "Link PRs and commits to issues automatically.", icon: "branch", status: "Active", tone: "success" },
            { name: "Slack",   desc: "Slash commands and channel digests.",            icon: "comment", status: "Active", tone: "success" },
            { name: "Datadog", desc: "Create issues from alerts.",                    icon: "chart",  status: "Setup",  tone: "warning" },
            { name: "PagerDuty", desc: "Sync on-call rotations & incidents.",         icon: "bell",   status: "Setup",  tone: "warning" },
            { name: "Linear",  desc: "Migrate from your previous tracker.",            icon: "list",   status: "—",      tone: "muted" },
            { name: "Webhook", desc: "Stream all events to your URL.",                 icon: "code",   status: "—",      tone: "muted" },
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

function Stat({ label, value }) {
  return (
    <div>
      <div className="text-xs muted">{label}</div>
      <div className="bold" style={{ fontSize: 16 }}>{value}</div>
    </div>
  );
}

Object.assign(window, { MembersView, NotificationSettingsView, IntegrationsView });
