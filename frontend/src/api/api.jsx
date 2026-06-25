// api.jsx — Real HTTP client for /api/v1/*
import { useState, useEffect } from 'react';

const BASE = "/api/v1";

// ─── Token storage ──────────────────────────────────────────────────────
export function getToken()       { return localStorage.getItem("access_token"); }
export function getRefresh()     { return localStorage.getItem("refresh_token"); }
export function saveTokens(pair) {
  localStorage.setItem("access_token",  pair.access_token);
  localStorage.setItem("refresh_token", pair.refresh_token);
}
export function clearTokens() {
  localStorage.removeItem("access_token");
  localStorage.removeItem("refresh_token");
  localStorage.removeItem("current_user");
  localStorage.removeItem("active_project_id");
}

// ─── Token refresh ──────────────────────────────────────────────────────
let _refreshPromise = null;
async function refreshTokens() {
  if (_refreshPromise) return _refreshPromise;
  _refreshPromise = fetch(BASE + "/auth/refresh", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: getRefresh() }),
  })
    .then((r) => r.json())
    .then((body) => {
      if (body.data) {
        saveTokens(body.data);
        return body.data.access_token;
      }
      throw new Error("refresh failed");
    })
    .finally(() => { _refreshPromise = null; });
  return _refreshPromise;
}

// ─── Core api() helper ──────────────────────────────────────────────────
// opts: { method, body, raw, formData }
export async function api(path, opts, _retry) {
  opts = opts || {};
  const token = getToken();
  const headers = {};
  if (token) headers["Authorization"] = "Bearer " + token;
  if (opts.body !== undefined) headers["Content-Type"] = "application/json";

  const res = await fetch(BASE + path, {
    method: opts.method || (opts.body !== undefined ? "POST" : "GET"),
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
  });

  // 401 → try refresh once
  if (res.status === 401 && !_retry) {
    try {
      await refreshTokens();
      return api(path, opts, true);
    } catch {
      clearTokens();
      window.location.hash = "#/login";
      throw new Error("Session expired — please sign in.");
    }
  }

  if (opts.raw) return res;

  const data = await res.json().catch(() => ({}));

  if (!res.ok) {
    throw new Error((data && data.message) || "Request failed (" + res.status + ")");
  }
  return data.data !== undefined ? data.data : data;
}

// ─── Download blob endpoint as a file ───────────────────────────────────
export async function apiDownload(path) {
  const res = await api(path, { raw: true });
  if (!res.ok) throw new Error("Export failed (" + res.status + ")");
  const blob = await res.blob();
  const cd = res.headers.get("Content-Disposition") || "";
  const match = cd.match(/filename="?([^"]+)"?/);
  const filename = match ? match[1] : path.split("/").pop();
  const a = document.createElement("a");
  a.href = URL.createObjectURL(blob);
  a.download = filename;
  document.body.appendChild(a); a.click(); a.remove();
  setTimeout(() => URL.revokeObjectURL(a.href), 4000);
}

// ─── Upload FormData (attachments) ───────────────────────────────────────
export async function apiUpload(path, file) {
  const fd = new FormData();
  fd.append("file", file);
  const token = getToken();
  const res = await fetch(BASE + path, {
    method: "POST",
    headers: token ? { Authorization: "Bearer " + token } : {},
    body: fd,
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error((body && body.message) || "Upload failed");
  return body.data !== undefined ? body.data : body;
}

// ─── Import raw file body (import endpoints expect raw bytes) ─────────────
export async function apiImport(path, file) {
  const token = getToken();
  const res = await fetch(BASE + path, {
    method: "POST",
    headers: token ? { Authorization: "Bearer " + token } : {},
    body: file,
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error((body && body.message) || "Import failed");
  return body.data !== undefined ? body.data : body;
}

// ─── useApi hook — load data on mount ────────────────────────────────────
// Returns { data, loading, error, reload, setData }
export function useApi(path, deps) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [tick, setTick] = useState(0);
  useEffect(() => {
    let live = true;
    setLoading(true);
    setError(null);
    if (!path) { setLoading(false); return; }
    api(path)
      .then((d) => { if (live) { setData(d); setLoading(false); } })
      .catch((e) => { if (live) { setError(e.message); setLoading(false); } });
    return () => { live = false; };
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [path, tick, ...(deps || [])]);
  return { data, loading, error, reload: () => setTick((t) => t + 1), setData };
}
