import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api, ApiError } from '../api';

export default function ForgotPassword() {
  const [email, setEmail] = useState('');
  const [sent, setSent] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    setErr(null);
    api.auth.forgotPassword(email)
      .then(() => setSent(true))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  };

  if (sent) return <div className="page"><header className="page-header"><h1>Forgot password</h1></header><p>Check your email for reset link. <Link to="/login">Sign in</Link></p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Forgot password</h1>
        <p className="page-intro">We’ll send a reset link to your email.</p>
      </header>
      <div className="page-content">
        <form onSubmit={submit}>
          <div className="page-form-row">
            <label><input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required /></label>
            <button type="submit" className="btn btn-primary">Send reset link</button>
          </div>
        </form>
        {err && <p className="page-error">{err}</p>}
        <p className="page-back"><Link to="/login">← Sign in</Link></p>
      </div>
    </div>
  );
}
