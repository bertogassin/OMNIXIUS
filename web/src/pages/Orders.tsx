import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api, ApiError } from '../api';

interface OrderRow {
  id: string | number;
  title?: string;
  status?: string;
  price?: number;
  seller_name?: string;
  buyer_name?: string;
  created_at?: number;
  urgent?: boolean;
  installment_plan?: string;
}

export default function Orders() {
  const [data, setData] = useState<{ asBuyer?: OrderRow[]; asSeller?: OrderRow[] } | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    api.users.myOrders()
      .then((r) => setData(r as { asBuyer?: OrderRow[]; asSeller?: OrderRow[] }))
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed to load orders'));
  }, []);

  if (err) return <p className="error">{err}</p>;
  if (!data) return <p>Loading…</p>;

  const asBuyer = data.asBuyer || [];
  const asSeller = data.asSeller || [];

  const renderCard = (o: OrderRow) => (
    <article key={String(o.id)} style={{ border: '1px solid var(--border)', borderRadius: 8, padding: 12, marginBottom: 8 }}>
      <Link to={`/order/${o.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
        <h3>{o.title || 'Order'} {o.urgent ? <span style={{ background: '#c00', color: '#fff', fontSize: 12, padding: '2px 6px', borderRadius: 4 }}>Urgent</span> : null}</h3>
        <p>{o.seller_name || o.buyer_name || ''} · {o.created_at ? new Date(o.created_at * 1000).toLocaleDateString() : ''}</p>
        <p>{typeof o.price === 'number' ? o.price.toFixed(2) : o.price}</p>
        <span className={'order-status order-status-' + (o.status || '').toLowerCase()}>{o.status || ''}</span>
      </Link>
    </article>
  );

  return (
    <div>
      <h1>My orders</h1>
      <section>
        <h2>As buyer</h2>
        {asBuyer.length ? <div className="orders-grid">{asBuyer.map((o) => renderCard(o))}</div> : <p>No orders yet.</p>}
      </section>
      <section>
        <h2>As seller</h2>
        {asSeller.length ? <div className="orders-grid">{asSeller.map((o) => renderCard(o))}</div> : <p>No orders yet.</p>}
      </section>
    </div>
  );
}
