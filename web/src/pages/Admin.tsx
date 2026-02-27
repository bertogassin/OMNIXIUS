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

  if (err) return <p className="error">{err}</p>;
  return (
    <div>
      <h1>Admin</h1>
      <p>Stats, reports, user lookup, ban/unban. Screen preserved.</p>
      {stats != null ? <pre>{JSON.stringify(stats, null, 2)}</pre> : null}
    </div>
  );
}
