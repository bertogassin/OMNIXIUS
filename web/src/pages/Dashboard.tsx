import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../auth';
import { api } from '../api';
import '../pages/Dashboard.css';

interface OrderRow {
  id?: number;
  title?: string;
  status?: string;
  created_at?: number;
}
interface RemittanceRow {
  amount?: number;
  currency?: string;
  to_identifier?: string;
  status?: string;
  created_at?: number;
}

export default function Dashboard() {
  const { user } = useAuth();
  const [profile, setProfile] = useState<{ name?: string; email?: string; role?: string; avatar_path?: string } | null>(null);
  const [balance, setBalance] = useState<number | null>(null);
  const [balanceLoading, setBalanceLoading] = useState(true);
  const [orders, setOrders] = useState<OrderRow[]>([]);
  const [remittances, setRemittances] = useState<RemittanceRow[]>([]);
  const [creditAmount, setCreditAmount] = useState('100');
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState('');

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const [u, ords, rem] = await Promise.all([
          api.users.me(),
          api.users.myOrders().then((r) => {
            const a = (r.asBuyer || []) as OrderRow[];
            const b = (r.asSeller || []) as OrderRow[];
            return [...a, ...b].slice(0, 20);
          }),
          api.remittances.my().then((r) => (Array.isArray(r) ? (r as RemittanceRow[]) : [])).catch(() => []),
        ]);
        if (!cancelled) {
          setProfile(u);
          setOrders(ords);
          setRemittances(rem.slice(0, 10));
        }
        try {
          const bal = await api.users.balance();
          const b = (bal as { balance?: number })?.balance;
          if (!cancelled) setBalance(typeof b === 'number' ? b : 0);
        } catch {
          if (!cancelled) setBalance(null);
        } finally {
          if (!cancelled) setBalanceLoading(false);
        }
      } catch (e) {
        if (!cancelled) setErr((e as { data?: { error?: string } })?.data?.error ?? 'Failed to load');
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, []);

  const handleBalanceCredit = async () => {
    const amt = parseFloat(creditAmount);
    if (!(amt > 0)) return;
    try {
      const res = await api.users.balanceCredit(amt) as { balance?: number };
      setBalance(res.balance != null ? res.balance : balance ?? 0);
    } catch {
      /* ignore */
    }
  };

  const displayName = profile?.name || profile?.email || user?.name || user?.email || 'User';
  const avatarUrl = profile?.avatar_path
    ? (typeof window !== 'undefined' && (window as unknown as { __OMNIXIUS_API_URL__?: string }).__OMNIXIUS_API_URL__ || 'http://localhost:3000') + '/uploads/' + profile.avatar_path
    : '';

  if (loading && !profile) {
    return (
      <div className="dashboard-page">
        <h1>Dashboard</h1>
        <p className="dashboard-loading">Loading…</p>
      </div>
    );
  }

  if (err) {
    return (
      <div className="dashboard-page">
        <h1>Dashboard</h1>
        <p className="dashboard-error">{err}</p>
      </div>
    );
  }

  return (
    <div className="dashboard-page">
      <p className="dashboard-status">Live. Settings, password recovery available.</p>
      <div className="dashboard-welcome">
        <h1 className="dashboard-welcome-title">{displayName.length > 30 ? displayName.slice(0, 27) + '…' : displayName}</h1>
      </div>
      <div className="dashboard-quick-links">
        <Link to="/marketplace" className="dashboard-quick-link">Trove · marketplace</Link>
        <Link to="/mail" className="dashboard-quick-link">Relay · mail</Link>
        <Link to="/orders" className="dashboard-quick-link">Orders</Link>
        <Link to="/vault" className="dashboard-quick-link">Crate · files</Link>
        <Link to="/ai" className="dashboard-quick-link">Oracle · AI</Link>
        <a href="/#map" className="dashboard-quick-link">Map</a>
      </div>
      <h2 className="dashboard-section">Profile</h2>
      <p>
        <Link to="/profile-edit" className="dashboard-btn">Edit profile</Link>
        {' '}
        <Link to="/settings" className="dashboard-btn">Settings</Link>
      </p>
      <div className="dashboard-profile-block">
        {avatarUrl ? <img src={avatarUrl} alt="" className="dashboard-avatar" /> : <div className="dashboard-avatar-placeholder" />}
        <p><strong>{profile?.name || '—'}</strong> · {profile?.email}</p>
        <p className="dashboard-muted">Role: {profile?.role || 'user'}</p>
        <p><Link to="/profile-edit" className="dashboard-btn dashboard-btn-primary">Edit</Link></p>
      </div>
      <h2 className="dashboard-section">Balance</h2>
      <div className="dashboard-balance-block">
        {balanceLoading ? (
          <p className="dashboard-muted">Loading…</p>
        ) : (
          <>
            <p className="dashboard-balance-amount">{balance != null ? balance.toFixed(2) : '0.00'} <span className="dashboard-balance-unit">units</span></p>
            <p className="dashboard-balance-stub">
              <button type="button" className="dashboard-btn-sm" onClick={handleBalanceCredit}>Add test credit</button>
              {' '}
              <input type="number" min="0.01" step="0.01" value={creditAmount} onChange={(e) => setCreditAmount(e.target.value)} style={{ width: 80 }} />
              <span className="dashboard-muted"> (stub for testing)</span>
            </p>
          </>
        )}
      </div>
      <h2 className="dashboard-section">My orders</h2>
      <div className="dashboard-orders-list">
        {orders.length === 0 ? (
          <p className="dashboard-muted">No orders yet. <Link to="/marketplace">Trove · marketplace</Link></p>
        ) : (
          <ul className="dashboard-list">
            {orders.map((o, i) => (
              <li key={o.id ?? i}>
                <span>{o.title ?? 'Order'} — {o.status}</span>
                <span>{o.created_at ? new Date(o.created_at * 1000).toLocaleDateString() : ''}</span>
              </li>
            ))}
          </ul>
        )}
      </div>
      <h2 className="dashboard-section">My remittances</h2>
      <div className="dashboard-remittances-list">
        {remittances.length === 0 ? (
          <p className="dashboard-muted">No remittances yet. <Link to="/remittances">Create one</Link>.</p>
        ) : (
          <>
            <ul className="dashboard-list">
              {remittances.map((r, i) => (
                <li key={i}>
                  <span>{r.amount != null ? Number(r.amount).toFixed(2) : ''} {r.currency || 'USD'} → {r.to_identifier || ''} · {r.status || 'pending'}</span>
                  <span>{r.created_at ? new Date(r.created_at * 1000).toLocaleDateString() : '—'}</span>
                </li>
              ))}
            </ul>
            <p><Link to="/remittances">View all → Remittances</Link></p>
          </>
        )}
      </div>
    </div>
  );
}
