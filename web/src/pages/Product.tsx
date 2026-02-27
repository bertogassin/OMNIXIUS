import { useParams, Link } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { api, ApiError } from '../api';

export default function Product() {
  const { id } = useParams<{ id: string }>();
  const [product, setProduct] = useState<unknown>(null);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    api.products.get(id)
      .then(setProduct)
      .catch((e: ApiError) => setErr(e.data?.error || 'Failed'));
  }, [id]);

  if (!id) return <div className="page"><p className="page-error">No product. <Link to="/marketplace">Marketplace</Link></p></div>;
  if (err) return <div className="page"><p className="page-error">{err}</p></div>;
  if (!product) return <div className="page"><h1>Product</h1><p className="page-loading">Loading…</p></div>;
  return (
    <div className="page">
      <p className="page-back"><Link to="/marketplace">← Marketplace</Link></p>
      <header className="page-header">
        <h1>Product {id}</h1>
      </header>
      <div className="page-content">
        <pre style={{ overflow: 'auto', fontSize: 0.85 }}>{JSON.stringify(product, null, 2)}</pre>
      </div>
    </div>
  );
}
