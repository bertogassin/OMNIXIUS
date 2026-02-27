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

  if (err) return <p className="error">{err}</p>;
  return (
    <div>
      <h1>Relay · Mail</h1>
      <p>Conversations and messages.</p>
      {list && (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {list.map((c: unknown, i: number) => (
            <li key={i} style={{ padding: 8, borderBottom: '1px solid var(--border)' }}>
              <Link to={`/conversation/${(c as { id?: number }).id}`}>Conversation {(c as { id?: number }).id}</Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
