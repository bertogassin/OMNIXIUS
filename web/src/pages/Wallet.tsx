import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Wallet() {
  const [balances, setBalances] = useState<unknown[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    api.wallet.balances()
      .then((r: unknown) => setBalances((r as { balances?: unknown[] })?.balances ?? []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'))
      .finally(() => setLoading(false));
  }, []);

  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  return (
    <div className="page">
      <header className="page-header">
        <h1>Wallet</h1>
        <p className="page-intro">Balances and transfers (Forge · trade & finance).</p>
      </header>
      <div className="page-content">
        {loading ? <p className="page-loading">Loading…</p> : (
          balances != null && balances.length > 0 ? (
            <ul className="page-list">{balances.map((b: unknown, i: number) => (
              <li key={i} className="page-list-item">{JSON.stringify(b)}</li>
            ))}</ul>
          ) : <p>No balances yet.</p>
        )}
      </div>
    </div>
  );
}
