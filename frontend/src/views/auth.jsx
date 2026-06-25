// auth.jsx — Login, Register, Forgot password, Accept invite
import { useState, useEffect } from 'react';
import { Icon } from '../components/icons';
import { Button, useToast } from '../components/components';
import { api, saveTokens } from '../api/api';
import { useApp } from '../store/AppContext';

export function LoginView({ nav, mode = "login" }) {
  const [fullName, setFullName] = useState("");
  const [email, setEmail]       = useState("");
  const [pwd, setPwd]           = useState("");
  const [showPwd, setShowPwd]   = useState(false);
  const [loading, setLoading]   = useState(false);
  const [error, setError]       = useState(null);
  const [sent, setSent]         = useState(false);
  const toast = useToast();
  const { onLoggedIn } = useApp();

  async function handleSubmit() {
    setError(null);
    setLoading(true);
    try {
      if (mode === "login") {
        const data = await api("/auth/login", { body: { email, password: pwd } });
        saveTokens(data);
        await onLoggedIn();
        nav("dashboard");
      } else if (mode === "register") {
        const data = await api("/auth/register", { body: { full_name: fullName, email, password: pwd } });
        saveTokens(data);
        await onLoggedIn();
        nav("dashboard");
      } else if (mode === "forgot") {
        await api("/auth/forgot-password", { body: { email } });
        setSent(true);
      }
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth" style={{ position: "absolute", inset: 0 }}>
      <div className="auth-side">
        <div className="row gap-3" style={{ color: "#fff" }}>
          <div style={{
            width: 36, height: 36, borderRadius: 9,
            background: "rgba(255,255,255,.12)",
            display: "grid", placeItems: "center",
            fontWeight: 700, letterSpacing: "-.02em",
            boxShadow: "inset 0 0 0 1px rgba(255,255,255,.18)",
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
              ["kanban",   "Boards built for runbooks, deploys and on-call rotations"],
              ["telegram", "Native Telegram bot — no Zapier, no glue code"],
              ["notes",    "Confluence-style wiki that lives next to your tickets"],
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

          {sent && mode === "forgot" && (
            <div style={{ padding: "12px 16px", borderRadius: 8, background: "#DCFCE7", color: "#166534", fontSize: 14, marginBottom: 20 }}>
              Reset link sent — check your inbox.
            </div>
          )}

          {error && (
            <div style={{ padding: "10px 14px", borderRadius: 8, background: "#FEF2F2", color: "#991B1B", fontSize: 13, marginBottom: 16, border: "1px solid #FECACA" }}>
              {error}
            </div>
          )}

          {mode === "register" && (
            <div style={{ padding: "16px", borderRadius: 10, background: "#EEF2FF", border: "1px solid #C7D2FE", marginBottom: 20 }}>
              <div style={{ fontWeight: 600, fontSize: 14, color: "#3730A3", marginBottom: 6 }}>Registration is invite-only</div>
              <div style={{ fontSize: 13, color: "#4338CA", lineHeight: 1.5 }}>
                This workspace only allows invited members. Ask an admin to send you an invite link, or{" "}
                <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }} style={{ color: "#4F46E5", fontWeight: 600 }}>sign in</a>{" "}
                if you already have an account.
              </div>
            </div>
          )}
          {mode !== "register" && (
            <>
              <div style={{ marginBottom: 14 }}>
                <label className="label">Email</label>
                <input className="input" type="email" value={email} onChange={(e) => setEmail(e.target.value)} placeholder="you@company.com"
                  onKeyDown={(e) => e.key === "Enter" && !loading && handleSubmit()}/>
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
                    <input className="input" type={showPwd ? "text" : "password"} value={pwd}
                      onChange={(e) => setPwd(e.target.value)} placeholder="••••••••"
                      onKeyDown={(e) => e.key === "Enter" && !loading && handleSubmit()}/>
                    <button onClick={() => setShowPwd(!showPwd)} type="button" style={{ position: "absolute", right: 8, top: "50%", transform: "translateY(-50%)", background: "transparent", border: 0, color: "var(--text-muted)" }} aria-label="Toggle password visibility">
                      <Icon name={showPwd ? "eyeOff" : "eye"} size={15}/>
                    </button>
                  </div>
                </div>
              )}

              <Button variant="primary" size="lg" onClick={handleSubmit} disabled={loading} style={{ width: "100%", justifyContent: "center", opacity: loading ? .7 : 1 }}>
                {loading ? "Please wait…" : mode === "forgot" ? "Send reset link" : "Sign in"}
              </Button>
            </>
          )}

          {mode !== "forgot" && (
            <>
              <div style={{ display: "flex", alignItems: "center", gap: 12, margin: "20px 0", color: "var(--text-muted)", fontSize: 12 }}>
                <div style={{ flex: 1, height: 1, background: "var(--border)" }}/>
                <span>or</span>
                <div style={{ flex: 1, height: 1, background: "var(--border)" }}/>
              </div>
              <div style={{ display: "grid", gap: 8 }}>
                <Button variant="secondary" size="lg" style={{ width: "100%", justifyContent: "center" }}
                  onClick={() => { window.location.href = "/api/v1/auth/google"; }}>
                  <svg width="16" height="16" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.5 12.3c0-.8-.1-1.6-.2-2.3H12v4.4h5.9c-.3 1.4-1 2.6-2.2 3.4v2.8h3.6c2.1-1.9 3.2-4.8 3.2-8.3z"/><path fill="#34A853" d="M12 23c2.9 0 5.3-1 7.1-2.6l-3.6-2.8c-1 .7-2.3 1.1-3.5 1.1-2.7 0-5-1.8-5.8-4.3H2.6v2.9C4.4 20.6 7.9 23 12 23z"/><path fill="#FBBC04" d="M6.2 14.4c-.2-.6-.3-1.2-.3-1.9s.1-1.3.3-1.9V7.7H2.6C1.9 9 1.5 10.5 1.5 12s.4 3 1.1 4.3l3.6-1.9z"/><path fill="#EA4335" d="M12 5.4c1.5 0 2.9.5 4 1.5l3.1-3C17.3 2.2 14.9 1 12 1 7.9 1 4.4 3.4 2.6 7l3.6 2.8c.8-2.5 3.1-4.4 5.8-4.4z"/></svg>
                  Continue with Google
                </Button>
              </div>
            </>
          )}

          <div className="text-sm secondary" style={{ marginTop: 24, textAlign: "center" }}>
            {mode === "login" && (<>Don't have access? Ask an admin for an invite link.</>)}
            {mode === "register" && (<>Already have one? <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }} style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Sign in</a></>)}
            {mode === "forgot" && (<>Remembered it? <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }} style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Back to sign in</a></>)}
          </div>

          {import.meta.env.DEV && (
            <div style={{ marginTop: 28, padding: "12px 14px", borderRadius: 8, background: "var(--surface-2, #F8FAFC)", fontSize: 12, color: "var(--text-muted)" }}>
              <strong style={{ display: "block", marginBottom: 4 }}>Demo accounts</strong>
              admin@jiraflow.com / Admin123! &nbsp;·&nbsp; member1@jiraflow.com / Member123!
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// ─── Accept Invite ────────────────────────────────────────────────────────────
export function AcceptInviteView({ nav }) {
  const { onLoggedIn } = useApp();
  const [fullName, setFullName] = useState("");
  const [pwd, setPwd]           = useState("");
  const [pwd2, setPwd2]         = useState("");
  const [showPwd, setShowPwd]   = useState(false);
  const [loading, setLoading]   = useState(false);
  const [error, setError]       = useState(null);
  const [token, setToken]       = useState("");

  useEffect(() => {
    const hash = window.location.hash; // e.g. "#/accept-invite?token=abc"
    const m = hash.match(/[?&]token=([^&]+)/);
    if (m) setToken(decodeURIComponent(m[1]));
  }, []);

  async function handleSubmit() {
    if (!token) { setError("Invalid invite link — token is missing."); return; }
    if (pwd !== pwd2) { setError("Passwords don't match."); return; }
    if (pwd.length < 8) { setError("Password must be at least 8 characters."); return; }
    setError(null);
    setLoading(true);
    try {
      const data = await api("/invites/accept", { body: { token, full_name: fullName, password: pwd } });
      saveTokens(data);
      await onLoggedIn();
      nav("dashboard");
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth" style={{ position: "absolute", inset: 0 }}>
      <div className="auth-side">
        <div className="row gap-3" style={{ color: "#fff" }}>
          <div style={{ width: 36, height: 36, borderRadius: 9, background: "rgba(255,255,255,.12)", display: "grid", placeItems: "center", fontWeight: 700, boxShadow: "inset 0 0 0 1px rgba(255,255,255,.18)" }}>F</div>
          <div style={{ fontWeight: 600, fontSize: 16 }}>Forge</div>
        </div>
        <div style={{ maxWidth: 460 }}>
          <h1 style={{ fontSize: 34, lineHeight: 1.1, letterSpacing: "-.02em", fontWeight: 600, margin: "0 0 16px" }}>
            You've been invited to join Forge
          </h1>
          <p style={{ color: "#C7D2FE", fontSize: 15, lineHeight: 1.6 }}>
            Set up your account to get started. Your email and role have already been configured by the admin.
          </p>
        </div>
        <div aria-hidden="true" style={{ position: "absolute", right: -120, top: -120, width: 380, height: 380, borderRadius: "50%", background: "radial-gradient(circle, #818CF8 0%, transparent 60%)", filter: "blur(20px)", opacity: .5 }}/>
      </div>

      <div className="auth-form-wrap">
        <div className="auth-form">
          <h2 style={{ fontSize: 24, fontWeight: 600, letterSpacing: "-.01em", margin: "0 0 6px" }}>Create your account</h2>
          <p className="secondary" style={{ margin: "0 0 28px" }}>Complete your profile to accept the invite.</p>

          {!token && (
            <div style={{ padding: "12px 16px", borderRadius: 8, background: "#FEF2F2", color: "#991B1B", fontSize: 13, marginBottom: 16, border: "1px solid #FECACA" }}>
              Invalid invite link. Please ask your admin to resend it.
            </div>
          )}

          {error && (
            <div style={{ padding: "10px 14px", borderRadius: 8, background: "#FEF2F2", color: "#991B1B", fontSize: 13, marginBottom: 16, border: "1px solid #FECACA" }}>
              {error}
            </div>
          )}

          <div style={{ marginBottom: 14 }}>
            <label className="label">Full name</label>
            <input className="input" placeholder="Aisha Toshmatova" value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && !loading && handleSubmit()}/>
          </div>

          <div style={{ marginBottom: 14 }}>
            <div style={{ position: "relative" }}>
              <label className="label">Password</label>
              <input className="input" type={showPwd ? "text" : "password"} value={pwd}
                placeholder="At least 8 characters"
                onChange={(e) => setPwd(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && !loading && handleSubmit()}/>
              <button onClick={() => setShowPwd(!showPwd)} type="button" style={{ position: "absolute", right: 8, bottom: 9, background: "transparent", border: 0, color: "var(--text-muted)" }}>
                <Icon name={showPwd ? "eyeOff" : "eye"} size={15}/>
              </button>
            </div>
          </div>

          <div style={{ marginBottom: 20 }}>
            <label className="label">Confirm password</label>
            <input className="input" type={showPwd ? "text" : "password"} value={pwd2}
              placeholder="••••••••"
              onChange={(e) => setPwd2(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && !loading && handleSubmit()}/>
          </div>

          <Button variant="primary" size="lg" onClick={handleSubmit}
            disabled={loading || !token || !fullName || !pwd}
            style={{ width: "100%", justifyContent: "center", opacity: (loading || !token) ? .7 : 1 }}>
            {loading ? "Creating account…" : "Accept invite & join"}
          </Button>

          <div className="text-sm secondary" style={{ marginTop: 24, textAlign: "center" }}>
            Already have an account?{" "}
            <a href="#/login" onClick={(e) => { e.preventDefault(); nav("login"); }}
              style={{ color: "var(--indigo-600)", fontWeight: 500, textDecoration: "none" }}>Sign in</a>
          </div>
        </div>
      </div>
    </div>
  );
}
