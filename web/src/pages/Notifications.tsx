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

  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  if (!list) return <div className="page"><h1>Notifications</h1><p className="page-loading">Loading…</p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Notifications</h1>
        <p className="page-intro">History and alerts.</p>
      </header>
      <div className="page-content">
        {list.length === 0 ? <p>No notifications.</p> : (
          <ul className="page-list">{list.map((n: unknown, i: number) => (
            <li key={i} className="page-list-item">{String(JSON.stringify(n))}</li>
          ))}</ul>
        )}
      </div>
    </div>
  );
}
