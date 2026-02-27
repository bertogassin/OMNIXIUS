import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api, ApiError } from '../api';

export default function Order() {
  const { id } = useParams<{ id: string }>();
  const [order, setOrder] = useState<Record<string, unknown> | null>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    Promise.all([api.users.me(), api.orders.get(id)])
      .then(([, o]) => {
        setOrder(o as Record<string, unknown>);
      })
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, [id]);

  if (!id) return <div className="page"><p className="page-error">No order ID. <Link to="/orders">Orders</Link></p></div>;
  if (err) return <div className="page"><p className="page-error">{err}. <Link to="/orders">Orders</Link></p></div>;
  if (!order) return <div className="page"><h1>Order</h1><p className="page-loading">Loading…</p></div>;

  const status = (order.status as string) || '';
  const isSeller = order.seller_id === (api.user?.id ?? 0);
  const price = typeof order.price === 'number' ? order.price.toFixed(2) : String(order.price ?? '—');
  const date = order.created_at ? new Date((order.created_at as number) * 1000).toLocaleString() : '—';
  const urgent = order.urgent === true || order.urgent === 1;

  const accept = () => {
    api.orders.update(id, { status: 'confirmed' }).then(() => setOrder((o) => o ? { ...o, status: 'confirmed' } : o)).catch(() => {});
  };
  const decline = () => {
    if (!window.confirm('Decline this order?')) return;
    api.orders.update(id, { status: 'cancelled' }).then(() => setOrder((o) => o ? { ...o, status: 'cancelled' } : o)).catch(() => {});
  };

  return (
    <div className="page">
      <p className="page-back"><Link to="/orders">← Back to orders</Link></p>
      <header className="page-header">
        <h1>{String(order.title || 'Order')} {urgent ? <span style={{ background: '#c00', color: '#fff', fontSize: 12, padding: '2px 8px', borderRadius: 4 }}>Urgent</span> : null}</h1>
      </header>
      <div className="page-content">
        <article className="page-card">
          <p><strong>Price:</strong> {price}</p>
          <p><strong>Status:</strong> <span className={'order-status order-status-' + status.toLowerCase()}>{status}</span></p>
          <p className="page-intro">{date}</p>
          {isSeller && status === 'pending' && (
            <div style={{ marginTop: 16 }}>
              <p>Accept or decline:</p>
              <button type="button" className="btn btn-primary" onClick={accept}>Accept</button>
              <button type="button" className="btn btn-outline" onClick={decline}>Decline</button>
            </div>
          )}
        </article>
      </div>
    </div>
  );
}
