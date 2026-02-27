import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Profile() {
  const [me, setMe] = useState<{ id?: number; email?: string; name?: string; role?: string } | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.users.me()
      .then((r) => {
        const u = r as { id?: number; email?: string; name?: string; role?: string };
        setMe(u);
        if (typeof u?.id === 'number') api.user = { id: u.id, email: u.email, name: u.name, role: u.role };
      })
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, []);

  if (err) return <p className="error">{err}</p>;
  if (!me) return <p>Loading…</p>;
  return (
    <div>
      <h1>Profile</h1>
      <p><strong>Name:</strong> {me.name ?? '—'}</p>
      <p><strong>Email:</strong> {me.email ?? '—'}</p>
      <p><strong>Role:</strong> {me.role ?? '—'}</p>
    </div>
  );
}
