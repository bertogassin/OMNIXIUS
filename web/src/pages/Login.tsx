import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../auth';
import '../pages/Login.css';

const API_STORAGE_KEY = 'omnixius_api_url';

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
        <p className="login-note">
          API URL: set in .env <code>VITE_API_URL</code> (e.g. http://localhost:3000) or use the block above.
        </p>
      </div>
    </div>
  );
}
