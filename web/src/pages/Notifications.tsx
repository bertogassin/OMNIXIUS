import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Notifications() {
  const [list, setList] = useState<unknown[] | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.notifications.history()
      .then((r) => setList(Array.isArray(r) ? r : []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, []);

  if (err) return <p className="error">{err}</p>;
  if (!list) return <p>Loading…</p>;
  return (
    <div>
      <h1>Notifications</h1>
      {list.length === 0 ? <p>No notifications.</p> : (
        <ul style={{ listStyle: 'none', padding: 0 }}>{list.map((n: unknown, i: number) => (
          <li key={i} style={{ padding: 12, borderBottom: '1px solid var(--border)' }}>{String(JSON.stringify(n))}</li>
        ))}</ul>
      )}
    </div>
  );
}
