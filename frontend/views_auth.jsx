// views_auth.jsx — login, register, forgot password

function LoginView({ nav, mode = "login", tg }) {
  const [email, setEmail] = React.useState("maya@forge.dev");
  const [pwd, setPwd] = React.useState("••••••••••");
  const [showPwd, setShowPwd] = React.useState(false);

  return (
    <div className="auth" style={{ position: "absolute", inset: 0 }}>
      <div className="auth-side">
        <div className="row gap-3" style={{ color: "#fff" }}>
          <div style={{
            width: 36, height: 36, borderRadius: 9,
            background: "rgba(255,255,255,.12)",
            display: "grid", placeItems: "center",
            fontWeight: 700, letterSpacing: "-.02em",
            boxShadow: "inset 0 0 0 1px rgba(255,255,255,.18)"
          }}>F</div>
          <div style={{ fontWeight: 600, fontSize: 16 }}>Forge</div>
        </div>

        <div style={{ maxWidth: 460 }}>
          <div style={{ display: "inline-flex", alignItems: "center", gap: 6, padding: "4px 10px", borderRadius: 999, background: "rgba(255,255,255,.10)", fontSize: 12, color: "#C7D2FE", marginBottom: 18 }}>
            <Icon name="sparkle" size={12}/> New in 2.4 — Telegram digests, mobile boards
          </div>
          <h1 style={{ fontSize: 38, lineHeight: 1.1, letterSpacing: "-.02em", fontWeight: 600, margin: "0 0 16px" }}>
            Ship infrastructure with the cadence of product.
          </h1>
          <p style={{ color: "#C7D2FE", fontSize: 15, lineHeight: 1.6, margin: 0 }}>
            Forge is the single workspace where your platform team plans sprints, runs incidents, writes runbooks — and pings everyone on Telegram, not just email.
          </p>

          <div style={{ marginTop: 36, display: "grid", gap: 12 }}>
            {[
              ["kanban",  "Boards built for runbooks, deploys and on-call rotations"],
              ["telegram","Native Telegram bot — no Zapier, no glue code"],
              ["notes",   "Confluence-style wiki that lives next to your tickets"],
            ].map(([ic, t]) => (
              <div key={t} className="row gap-3">
                <div style={{ width: 28, height: 28, borderRadius: 7, background: "rgba(255,255,255,.12)", display: "grid", placeItems: "center", color: "#A5B4FC", flexShrink: 0 }}>
                  <Icon name={ic} size={15}/>
                </div>
                <div style={{ color: "#E0E7FF", fontSize: 13.5 }}>{t}</div>
              </div>
            ))}
          </div>
        </div>

        <div style={{ color: "rgba(199,210,254,.7)", fontSize: 12 }}>
          Trusted by infra teams at <strong style={{ color: "#fff" }}>Linear-ish</strong>, <strong style={{ color: "#fff" }}>Vapor</strong>, <strong style={{ color: "#fff" }}>Hexstack</strong> and 380+ others.
        </div>

        {/* decorative gradient blob */}
        <div aria-hidden="true" style={{ position: "absolute", right: -120, top: -120, width: 380, height: 380, borderRadius: "50%", background: "radial-gradient(circle, #818CF8 0%, transparent 60%)", filter: "blur(20px)", opacity: .5 }}/>
        <div aria-hidden="true" style={{ position: "absolute", right: 40, bottom: -160, width: 320, height: 320, borderRadius: "50%", background: "radial-gradient(circle, #4338CA 0%, transparent 60%)", filter: "blur(20px)", opacity: .6 }}/>
      </div>

      <div className="auth-form-wrap">
        <div className="auth-form">
          <h2 style={{ fontSize: 24, fontWeight: 600, letterSpacing: "-.01em", margin: "0 0 6px" }}>
            {mode === "register" ? "Create your account" : mode === "forgot" ? "Reset your password" : "Welcome back"}
          </h2>
          <p className="secondary" style={{ margin: "0 0 28px" }}>
            {mode === "register" ? "Spin up a workspace in 30 seconds." : mode === "forgot" ? "We'll send a reset link to your email." : "Sign in to continue to Forge."}
          </p>

          {mode === "register" && (
            <div style={{ marginBottom: 14 }}>
              <label className="label">Full name</label>
              <input className="input" placeholder="Maya Chen"/>
            </div>
          )}

          <div style={{ marginBottom: 14 }}>
            <label className="label">Email</label>
            <input className="input" type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@company.com"/>
          </div>

          {mode !== "forgot" && (
            <div style={{ marginBottom: 14 }}>
              <div className="row" style={{ justifyContent: "space-between" }}>
                <label className="label">Password</label>
                {mode === "login" && (
                  <a href="#/forgot" onClick={(e) => { e.preventDefault(); nav("forgot"); }} style={{ fontSize: 12, color: "var(--indigo-600)", textDecoration: "none" }}>Forgot?</a>
                )}
              </div>
              <div style={{ position: "relative" }}>
                <input className="input" type={showPwd ? "text" : "password"} value={pwd} onChange={(e) => setPwd(e.target.value)} placeholder="••••••••"/>
                <button onClick={() => setShowPwd(!showPwd)} type="button" style={{ position: "absolute", right: 8, top: "50%", transform: "translateY(-50%)", background: "transparent", border: 0, color: "var(--text-muted)" }} aria-label="Toggle password visibility">
                  <Icon name={showPwd ? "eyeOff" : "eye"} size={15}/>
                </button>
              </div>
            </div>
          )}

          {mode === "login" && (
            <label className="row gap-2 text-sm" style={{ margin: "12px 0 18px", color: "var(--text-secondary)" }}>
              <input type="checkbox" defaultChecked style={{ accentColor: "var(--indigo-600)" }}/> Keep me signed in for 30 days
            </label>
          )}

          <Button variant="primary" size="lg" onClick={() => nav("dashboard")} style={{ width: "100%", justifyContent: "center" }}>
            {mode === "register" ? "Create account" : mode === "forgot" ? "Send reset link" : "Sign in"}
          </Button>

          {mode !== "forgot" && (
            <>
              <div style={{ display: "flex", alignItems: "center", gap: 12, margin: "20px 0", color: "var(--text-muted)", fontSize: 12 }}>
                <div style={{ flex: 1, height: 1, background: "var(--border)" }}/>
                <span>or</span>
                <div style={{ flex: 1, height: 1, background: "var(--border)" }}/>
              </div>
              <div style={{ display: "grid", gap: 8 }}>
                <Button variant="secondary" size="lg" style={{ width: "100%", justifyContent: "center" }} onClick={() => nav("dashboard")}>
                  <svg width="16" height="16" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.4h5.9c-.3 1.4-1 2.6-2.2 3.4v2.8h3.6c2.1-1.9 3.2-4.8 3.2-8.3z"/><path fill="#34A853" d="M12 23c2.9 0 5.3-1 7.1-2.6l-3.6-2.8c-1 .7-2.3 1.1-3.5 1.1-2.7 0-5-1.8-5.8-4.3H2.6v2.9C4.4 20.6 7.9 23 12 23z"/><path fill="#FBBC04" d="M6.2 14.4c-.2-.6-.3-1.2-.3-1.9s.1-1.3.3-1.9V7.7H2.6C1.9 9 1.5 10.5 1.5 12s.4 3 1.1 4.3l3.6-1.9z"/><path fill="#EA4335" d="M12 5.4c1.5 0 2.9.5 4 1.5l3.1-3C17.3 2.2 14.9 1 12 1 7.9 1 4.4 3.4 2.6 7l3.6 2.8c.8-2.5 3.1-4.4 5.8-4.4z"/></svg>
                  Continue with Google
                </Button>
                <Button variant="secondary" size="lg" style={{ width: "100%", justifyContent: "center" }} onClick={() => nav("dashboard")}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 .5C5.7.5.5 5.7.5 12c0 5.1 3.3 9.4 7.8 10.9.6.1.8-.2.8-.6v-2c-3.2.7-3.8-1.4-3.8-1.4-.5-1.3-1.3-1.6-1.3-1.6-1-.7.1-.7.1-.7 1.1.1 1.7 1.1 1.7 1.1 1 1.8 2.7 1.3 3.3 1 .1-.7.4-1.3.7-1.6-2.5-.3-5.2-1.3-5.2-5.7 0-1.3.5-2.3 1.2-3.1-.1-.3-.5-1.5.1-3.1 0 0 1-.3 3.2 1.2.9-.3 1.9-.4 2.9-.4s2 .1 2.9.4c2.2-1.5 3.2-1.2 3.2-1.2.6 1.6.2 2.8.1 3.1.8.8 1.2 1.9 1.2 3.1 0 4.4-2.7 5.4-5.2 5.7.4.4.8 1.1.8 2.2v3.2c0 .3.2.7.8.6 4.5-1.5 7.8-5.8 7.8-10.9C23.5 5.7 18.3.5 12 .5z"/></svg>
                  Continue with GitHub
                </Button>
              </div>
            </>
          )}

          <div className="text-sm secondary" style={{ marginTop: 24, textAlign: "center" }}>
            {mode === "login" && (<>New to Forge? <a href="#/register" onClick={(e) => { e.preventDefault(); nav("register"); }} style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Create an account</a></>)}
            {mode === "register" && (<>Already have one? <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }} style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Sign in</a></>)}
            {mode === "forgot" && (<>Remembered it? <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }} style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Back to sign in</a></>)}
          </div>
        </div>
      </div>
    </div>
  );
}

window.LoginView = LoginView;
