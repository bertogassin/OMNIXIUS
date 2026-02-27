import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api';
import '../pages/Login.css';

export default function ProfileEdit() {
  const [name, setName] = useState('');
  const [success, setSuccess] = useState('');
  const [error, setError] = useState('');
  const [pwdError, setPwdError] = useState('');
  const [currentPwd, setCurrentPwd] = useState('');
  const [newPwd, setNewPwd] = useState('');
  const [newPwd2, setNewPwd2] = useState('');
  const [pwdSuccess, setPwdSuccess] = useState('');

  useEffect(() => {
    api.users.me().then((u) => setName(u.name || '')).catch(() => {});
  }, []);

  const handleProfileSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    try {
      await api.users.updateMe({ name: name.trim() });
      setSuccess('Saved.');
    } catch (err: unknown) {
      setError((err as { data?: { error?: string } })?.data?.error ?? 'Failed');
    }
  };

  const handlePwdSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setPwdError('');
    setPwdSuccess('');
    if (newPwd !== newPwd2) {
      setPwdError('Passwords do not match');
      return;
    }
    try {
      await api.auth.changePassword(currentPwd, newPwd);
      setPwdSuccess('Password changed.');
      setCurrentPwd('');
      setNewPwd('');
      setNewPwd2('');
    } catch (err: unknown) {
      setPwdError((err as { data?: { error?: string } })?.data?.error ?? 'Failed');
    }
  };

  return (
    <div className="login-page" style={{ padding: '2rem' }}>
      <div className="login-card" style={{ maxWidth: 500 }}>
        <p><Link to="/settings">Security, sessions → Settings</Link></p>
        <h1>Profile</h1>
        <form onSubmit={handleProfileSubmit} className="login-form">
          {error && <p className="login-error">{error}</p>}
          {success && <p style={{ color: '#00d4aa' }}>{success}</p>}
          <label>Name <input type="text" value={name} onChange={(e) => setName(e.target.value)} className="login-input" /></label>
          <button type="submit" className="login-btn">Save</button>
        </form>
        <h2 style={{ marginTop: '2rem' }}>Change password</h2>
        <form onSubmit={handlePwdSubmit} className="login-form">
          {pwdError && <p className="login-error">{pwdError}</p>}
          {pwdSuccess && <p style={{ color: '#00d4aa' }}>{pwdSuccess}</p>}
          <label>Current password <input type="password" value={currentPwd} onChange={(e) => setCurrentPwd(e.target.value)} className="login-input" /></label>
          <label>New password <input type="password" value={newPwd} onChange={(e) => setNewPwd(e.target.value)} minLength={8} className="login-input" /></label>
          <label>Confirm <input type="password" value={newPwd2} onChange={(e) => setNewPwd2(e.target.value)} minLength={8} className="login-input" /></label>
          <button type="submit" className="login-btn">Save</button>
        </form>
      </div>
    </div>
  );
}
