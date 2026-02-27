import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { api } from '../api';
import '../pages/Login.css';

export default function ResetPassword() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token') ?? '';
  const [password, setPassword] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const noToken = !token.trim();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSubmitting(true);
    try {
      await api.auth.resetPassword(token, password);
      setSuccess(true);
    } catch (err: unknown) {
      setError((err as { data?: { error?: string } })?.data?.error ?? 'Invalid or expired link. Request a new reset from forgot password.');
    } finally {
      setSubmitting(false);
    }
  };

  if (noToken) {
    return (
      <div className="login-page">
        <div className="login-card">
          <h1>Set new password</h1>
          <p className="login-error">Invalid or expired link. Request a new reset from forgot password.</p>
          <p className="login-links"><Link to="/forgot-password">Forgot password?</Link></p>
        </div>
      </div>
    );
  }

  if (success) {
    return (
      <div className="login-page">
        <div className="login-card">
          <h1>Set new password</h1>
          <p style={{ color: 'var(--accent, #00d4aa)' }}>Password updated. You can sign in.</p>
          <p className="login-links"><Link to="/login">← Sign in</Link></p>
        </div>
      </div>
    );
  }

  return (
    <div className="login-page">
      <div className="login-card">
        <h1>Set new password</h1>
        {error && <p className="login-error">{error}</p>}
        <form onSubmit={handleSubmit} className="login-form">
          <label>
            New password (min 8)
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={8}
              autoComplete="new-password"
              className="login-input"
            />
          </label>
          <button type="submit" className="login-btn" disabled={submitting}>
            {submitting ? 'Saving…' : 'Set password'}
          </button>
        </form>
        <p className="login-links"><Link to="/forgot-password">Forgot password?</Link></p>
      </div>
    </div>
  );
}
