import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Remittances() {
  const [list, setList] = useState<unknown[] | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.remittances.my()
      .then((r) => setList(Array.isArray(r) ? r : []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, []);

  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Remittances</h1>
        <p className="page-intro">My remittances and cross-border transfers.</p>
      </header>
      <div className="page-content">
        {list == null ? <p className="page-loading">Loading…</p> : list.length === 0 ? (
          <p>No remittances yet.</p>
        ) : (
          <ul className="page-list">{list.map((r: unknown, i: number) => (
            <li key={i} className="page-list-item">{JSON.stringify(r)}</li>
          ))}</ul>
        )}
      </div>
    </div>
  );
}
