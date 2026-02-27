import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { api, ApiError } from '../api';

export default function Register() {
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [err, setErr] = useState<string | null>(null);

  const submit = (e: React.FormEvent) => {
    e.preventDefault();
    setErr(null);
    api.auth.register(email, password, name)
      .then((r: { token?: string; user?: unknown }) => {
        if (r.token) api.setToken(r.token);
        if (r.user) api.user = r.user as { id: number; email?: string; name?: string; role?: string };
        navigate('/');
      })
      .catch((e: ApiError) => setErr(e.data?.error || 'Registration failed'));
  };

  return (
    <div className="page">
      <header className="page-header">
        <h1>Create account</h1>
        <p className="page-intro">Create your OMNIXIUS account.</p>
      </header>
      <div className="page-content">
        <form onSubmit={submit}>
          <div className="page-form-row">
            <label><input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required /></label>
            <label><input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required /></label>
            <label><input type="text" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} /></label>
            <button type="submit" className="btn btn-primary">Register</button>
          </div>
        </form>
        {err && <p className="page-error">{err}</p>}
        <p className="page-back"><Link to="/login">← Sign in</Link></p>
      </div>
    </div>
  );
}
