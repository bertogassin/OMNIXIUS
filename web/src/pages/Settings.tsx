import { useState } from 'react';
import { api, ApiError } from '../api';

export default function Settings() {
  const [msg, setMsg] = useState<string | null>(null);

  const changePassword = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const form = e.currentTarget;
    const current = (form.elements.namedItem('current') as HTMLInputElement).value;
    const newP = (form.elements.namedItem('new') as HTMLInputElement).value;
    if (!current || !newP) { setMsg('Fill both fields'); return; }
    api.auth.changePassword(current, newP)
      .then(() => setMsg('Password updated'))
      .catch((e: ApiError) => setMsg(e.data?.error || 'Failed'));
  };

  return (
    <div className="page">
      <header className="page-header">
        <h1>Settings</h1>
        <p className="page-intro">Security and preferences.</p>
      </header>
      <div className="page-content">
        <section className="page-section">
          <h2>Change password</h2>
          <form onSubmit={changePassword}>
            <div className="page-form-row">
              <label><input name="current" type="password" placeholder="Current password" required /></label>
              <label><input name="new" type="password" placeholder="New password" required /></label>
              <button type="submit" className="btn btn-primary">Update</button>
            </div>
          </form>
        </section>
        {msg && <p>{msg}</p>}
      </div>
    </div>
  );
}
