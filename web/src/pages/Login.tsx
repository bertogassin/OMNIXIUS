import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../auth';
import { API_URL } from '../config';
import '../pages/Login.css';

const API_STORAGE_KEY = 'omnixius_api_url';
const TEST_USER_EMAIL = 'test@test.com';
const TEST_USER_PASSWORD = 'Test123!';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [remember, setRemember] = useState(true);
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [apiUrl, setApiUrl] = useState('');
  const [showApiBlock, setShowApiBlock] = useState(false);
  const [apiCheckResult, setApiCheckResult] = useState('');
  const [seedResult, setSeedResult] = useState('');
  const [seeding, setSeeding] = useState(false);
  const { login, token } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const from = (location.state as { from?: { pathname: string } })?.from?.pathname ?? '/';
  const isLocal = typeof window !== 'undefined' && /^(localhost|127\.0\.0\.1)$/.test(window.location.hostname);

  useEffect(() => {
    if (typeof window === 'undefined') return;
    const stored = localStorage.getItem(API_STORAGE_KEY);
    if (!isLocal && !stored) setShowApiBlock(true);
    setApiUrl(stored || '');
  }, [isLocal]);

  // Auto-enter on localhost: seed test user (if empty DB) then log in so browser opens "already inside"
  const didAutoEnter = React.useRef(false);
  useEffect(() => {
    if (token || !isLocal || didAutoEnter.current) return;
    const apiBase = 'http://localhost:3000';
    didAutoEnter.current = true;
    fetch(apiBase + '/api/seed-test-user', { method: 'POST' })
      .then((r) => r.json().catch(() => ({})))
      .then(() => login(TEST_USER_EMAIL, TEST_USER_PASSWORD, true))
      .then(() => navigate(from || '/', { replace: true }))
      .catch(() => { didAutoEnter.current = false; });
  }, [token, isLocal, login, navigate, from]);

  if (token) {
    navigate(from, { replace: true });
    return null;
  }

  const handleSaveApiUrl = async () => {
    const url = apiUrl.trim().replace(/\/+$/, '') || apiUrl.trim();
    if (!url) {
      setApiCheckResult('Enter URL');
      return;
    }
    setApiCheckResult('…');
    localStorage.setItem(API_STORAGE_KEY, url);
    try {
      const r = await fetch(url + '/health', { method: 'GET' });
      if (r.ok) {
        window.location.reload();
      } else {
        setApiCheckResult('Backend not running. Run: go run . in backend-go');
      }
    } catch {
      setApiCheckResult('Backend not reachable. Run: go run . in backend-go');
    }
  };

  const handleSeedTestUser = async () => {
    const base =
      (typeof window !== 'undefined' && (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') ? 'http://localhost:3000' : '') ||
      API_URL ||
      apiUrl.trim().replace(/\/+$/, '') ||
      (typeof window !== 'undefined' && localStorage.getItem(API_STORAGE_KEY)) ||
      '';
    const apiBase = String(base || 'http://localhost:3000').replace(/\/+$/, '');
    setSeedResult('');
    setSeeding(true);
    setError('');
    try {
      const r = await fetch(apiBase + '/api/seed-test-user', { method: 'POST' });
      const data = await r.json().catch(() => ({}));
      if (r.status === 201) {
        const d = data as { token?: string; user?: unknown };
        if (d.token && d.user) {
          try {
            await login(TEST_USER_EMAIL, TEST_USER_PASSWORD, true);
            navigate(from, { replace: true });
            return;
          } catch {
            setEmail(TEST_USER_EMAIL);
            setPassword(TEST_USER_PASSWORD);
            setSeedResult('Test user created. Click Sign in.');
          }
        } else {
          setEmail(TEST_USER_EMAIL);
          setPassword(TEST_USER_PASSWORD);
          setSeedResult('Test user created. Click Sign in.');
        }
      } else {
        setSeedResult((data as { message?: string }).message || (data as { error?: string }).error || 'Done');
      }
    } catch {
      setSeedResult('Backend not reachable. Is it running? go run . in backend-go');
    } finally {
      setSeeding(false);
    }
  };

  const [copied, setCopied] = useState(false);
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text).then(() => { setCopied(true); setTimeout(() => setCopied(false), 2000); }).catch(() => {});
  };
  const appUrl = typeof window !== 'undefined' ? window.location.origin + '/app/' : 'http://localhost:5173/app/';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await login(email.trim(), password, remember);
      navigate(from, { replace: true });
    } catch (err: unknown) {
      const apiErr = err as { status?: number; data?: { error?: string }; message?: string };
      if (apiErr?.status === 0 || (apiErr?.message && /fetch|Failed/i.test(apiErr.message)))
        setError(apiErr?.data?.error ?? 'Cannot reach server. Set API URL?');
      else
        setError(apiErr?.data?.error ?? 'Sign in failed');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page login-page">
      <div className="login-card">
        <h1>Sign in</h1>
        {showApiBlock && (
          <div className="login-api-block">
            <p className="login-api-title">Backend URL</p>
            <p className="login-api-hint">Run <code>go run .</code> in backend-go, then paste URL and click Save.</p>
            <div className="login-api-row">
              <input
                type="url"
                value={apiUrl}
                onChange={(e) => setApiUrl(e.target.value)}
                placeholder="http://localhost:3000"
                className="login-input"
              />
              <button type="button" className="login-btn" onClick={handleSaveApiUrl}>Save</button>
            </div>
            {apiCheckResult && <p className="login-api-result">{apiCheckResult}</p>}
          </div>
        )}
        <div className="login-api-block" style={{ marginTop: '0.5rem' }}>
          <p className="login-api-title">Test user (empty DB only)</p>
          <button type="button" className="login-btn" onClick={handleSeedTestUser} disabled={seeding}>
            {seeding ? 'Creating…' : 'Create test user (test@test.com / Test123!)'}
          </button>
          {seedResult && <p className="login-api-result" style={{ color: seedResult.includes('reachable') ? '#e57373' : undefined }}>{seedResult}</p>}
        </div>
        {error && <p className="login-error">{error}</p>}
        <form onSubmit={handleSubmit} className="login-form">
          <label>
            Email
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
              className="login-input"
            />
          </label>
          <label>
            Password
            <div className="login-pwd-wrap">
              <input
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="current-password"
                className="login-input"
              />
              <button
                type="button"
                className="login-pwd-toggle"
                aria-label={showPassword ? 'Hide password' : 'Show password'}
                onClick={() => setShowPassword((v) => !v)}
              >
                {showPassword ? 'Hide' : 'Show'}
              </button>
            </div>
          </label>
          <label className="login-remember">
            <input
              type="checkbox"
              checked={remember}
              onChange={(e) => setRemember(e.target.checked)}
            />
            Keep me signed in
          </label>
          <button type="submit" className="login-btn" disabled={submitting}>
            {submitting ? 'Signing in…' : 'Sign in'}
          </button>
        </form>
        <p className="login-links">
          <Link to="/forgot-password">Forgot password?</Link>
          {' · '}
          <Link to="/register">Register</Link>
        </p>
        <div className="login-api-block" style={{ marginTop: '1rem' }}>
          <p className="login-api-title">App link (copy and send)</p>
          <div className="login-api-row">
            <input type="text" readOnly value={appUrl} className="login-input" style={{ flex: 1, cursor: 'text' }} onFocus={(e) => e.target.select()} />
            <button type="button" className="login-btn" onClick={() => copyToClipboard(appUrl)}>{copied ? 'Copied' : 'Copy'}</button>
          </div>
          {isLocal && (
            <>
              <p className="login-api-title" style={{ marginTop: '0.75rem' }}>Backend URL</p>
              <div className="login-api-row">
                <input type="text" readOnly value="http://localhost:3000" className="login-input" style={{ flex: 1, cursor: 'text' }} onFocus={(e) => e.target.select()} />
                <button type="button" className="login-btn" onClick={() => copyToClipboard('http://localhost:3000')}>Copy</button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
