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

  if (!id) return <p>No product. <Link to="/marketplace">Marketplace</Link></p>;
  if (err) return <p className="error">{err}</p>;
  if (!product) return <p>Loading…</p>;
  return (
    <div>
      <h1>Product {id}</h1>
      <pre>{JSON.stringify(product, null, 2)}</pre>
      <p><Link to="/marketplace">← Marketplace</Link></p>
    </div>
  );
}
