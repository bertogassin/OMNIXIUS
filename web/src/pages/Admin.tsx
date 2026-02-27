import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Admin() {
  const [stats, setStats] = useState<unknown>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.admin.stats()
      .then(setStats)
      .catch((e: ApiError) => setErr(e.data?.error || 'Admin only'));
  }, []);

  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Admin</h1>
        <p className="page-intro">Stats, reports, user lookup, ban/unban.</p>
      </header>
      <div className="page-content">
        {stats == null ? <p className="page-loading">Loading…</p> : (
          <pre style={{ overflow: 'auto', fontSize: 0.85 }}>{JSON.stringify(stats, null, 2)}</pre>
        )}
      </div>
    </div>
  );
}
