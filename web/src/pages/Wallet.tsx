import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Wallet() {
  const [balances, setBalances] = useState<unknown[] | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.wallet.balances()
      .then((r: unknown) => setBalances((r as { balances?: unknown[] })?.balances ?? []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, []);

  if (err) return <p className="error">{err}</p>;
  return (
    <div>
      <h1>Wallet</h1>
      <p>Balances and transfers (Forge · trade & finance).</p>
      {balances && (
        <ul>{Array.isArray(balances) && balances.map((b: unknown, i: number) => (
          <li key={i}>{JSON.stringify(b)}</li>
        ))}</ul>
      )}
    </div>
  );
}
