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

  if (err) return <p className="error">{err}</p>;
  return (
    <div>
      <h1>Remittances</h1>
      <p>My remittances and cross-border transfers.</p>
      {list && <ul>{list.map((r: unknown, i: number) => <li key={i}>{JSON.stringify(r)}</li>)}</ul>}
    </div>
  );
}
