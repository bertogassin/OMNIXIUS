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

  if (sent) return <p>Check your email for reset link. <Link to="/login">Sign in</Link></p>;
  return (
    <div>
      <h1>Forgot password</h1>
      <form onSubmit={submit}>
        <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        <button type="submit" className="btn btn-primary">Send reset link</button>
      </form>
      {err && <p className="error">{err}</p>}
      <p><Link to="/login">Sign in</Link></p>
    </div>
  );
}
