import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api, ApiError } from '../api';

export default function Mail() {
  const [list, setList] = useState<unknown[] | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.conversations.list()
      .then((r) => setList(Array.isArray(r) ? r : []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, []);

  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Relay · Mail</h1>
        <p className="page-intro">Conversations and messages.</p>
      </header>
      <div className="page-content">
        {list == null ? <p className="page-loading">Loading…</p> : list.length === 0 ? (
          <p>No conversations yet.</p>
        ) : (
          <ul className="page-list">
            {list.map((c: unknown, i: number) => (
              <li key={i} className="page-list-item">
                <Link to={`/conversation/${(c as { id?: number }).id}`}>Conversation {(c as { id?: number }).id}</Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
