import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api } from '../api';
import '../pages/Marketplace.css';

interface Product {
  id: number;
  name?: string;
  price?: number;
  category?: string;
}

export default function Marketplace() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    api.products
      .list({})
      .then((res) => setProducts((res.products ?? []) as Product[]))
      .catch((err: { status?: number; data?: { error?: string } }) => {
        setError(err?.data?.error ?? 'Failed to load products');
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="marketplace-page">
        <h1>Marketplace</h1>
        <p className="marketplace-loading">Loading…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="marketplace-page">
        <h1>Marketplace</h1>
        <p className="marketplace-error">{error}</p>
      </div>
    );
  }

  return (
    <div className="marketplace-page">
      <h1>Marketplace</h1>
      <p className="marketplace-intro">Products and services (Trove).</p>
      {products.length === 0 ? (
        <p className="marketplace-empty">No products yet.</p>
      ) : (
        <ul className="marketplace-list">
          {products.map((p) => (
            <li key={p.id} className="marketplace-item">
              <span className="marketplace-item-name">{p.name ?? 'Product #' + p.id}</span>
              {p.category && (
                <span className="marketplace-item-cat">{p.category}</span>
              )}
              {p.price != null && (
                <span className="marketplace-item-price">{p.price} USD</span>
              )}
              <Link to={`/product/${p.id}`} className="marketplace-item-link">
                View →
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
