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
    <div>
      <h1>Create account</h1>
      <form onSubmit={submit}>
        <input type="email" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} required />
        <input type="password" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} required />
        <input type="text" placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} />
        <button type="submit" className="btn btn-primary">Register</button>
      </form>
      {err && <p className="error">{err}</p>}
      <p><Link to="/login">Sign in</Link></p>
    </div>
  );
}
