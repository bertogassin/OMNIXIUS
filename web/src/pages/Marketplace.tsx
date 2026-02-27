import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { api } from '../api';
import '../pages/Marketplace.css';

interface Product {
  id: number;
  name?: string;
  price?: number;
  category?: string;
  is_service?: boolean;
  is_subscription?: boolean;
}

export default function Marketplace() {
  const [searchParams, setSearchParams] = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const q = searchParams.get('q') ?? '';
  const location = searchParams.get('location') ?? '';
  const category = searchParams.get('category') ?? '';
  const service = searchParams.get('service') === '1';
  const subscription = searchParams.get('subscription') === '1';
  const minPrice = searchParams.get('minPrice') ?? '';
  const maxPrice = searchParams.get('maxPrice') ?? '';

  useEffect(() => {
    const params: Record<string, string | number> = {};
    if (q) params.q = q;
    if (location) params.location = location;
    if (category) params.category = category;
    if (service) params.service = '1';
    if (subscription) params.subscription = '1';
    if (minPrice) params.min_price = parseFloat(minPrice) || 0;
    if (maxPrice) params.max_price = parseFloat(maxPrice) || 0;
    setLoading(true);
    api.products.list(params)
      .then((res) => setProducts((res.products ?? []) as Product[]))
      .catch((err: { data?: { error?: string } }) => setError(err?.data?.error ?? 'Failed to load'))
      .finally(() => setLoading(false));
  }, [q, location, category, service, subscription, minPrice, maxPrice]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const form = e.target as HTMLFormElement;
    const next = new URLSearchParams(searchParams);
    next.set('q', (form.q as HTMLInputElement)?.value?.trim() ?? '');
    next.set('location', (form.location as HTMLInputElement)?.value?.trim() ?? '');
    next.set('category', (form.category as HTMLSelectElement)?.value ?? '');
    next.set('service', (form.service as HTMLInputElement)?.checked ? '1' : '');
    next.set('subscription', (form.subscription as HTMLInputElement)?.checked ? '1' : '');
    next.set('minPrice', (form.minPrice as HTMLInputElement)?.value ?? '');
    next.set('maxPrice', (form.maxPrice as HTMLInputElement)?.value ?? '');
    setSearchParams(next);
  };

  return (
    <div className="marketplace-page">
      <p className="marketplace-status">Live. Filters sync to URL.</p>
      <h1>Trove · marketplace</h1>
      <form className="marketplace-filters" onSubmit={handleSubmit}>
        <input type="text" name="q" placeholder="Search" defaultValue={q} />
        <input type="text" name="location" placeholder="City" defaultValue={location} />
        <select name="category"><option value="">All categories</option></select>
        <label className="marketplace-filter-cb"><input type="checkbox" name="service" defaultChecked={service} /> Services only</label>
        <label className="marketplace-filter-cb"><input type="checkbox" name="subscription" defaultChecked={subscription} /> Subscriptions only</label>
        <input type="number" name="minPrice" placeholder="From" min={0} step={0.01} defaultValue={minPrice || undefined} />
        <input type="number" name="maxPrice" placeholder="To" min={0} step={0.01} defaultValue={maxPrice || undefined} />
        <button type="submit" className="marketplace-search-btn">Search</button>
      </form>
      {error && <p className="marketplace-error">{error}</p>}
      {loading ? (
        <p className="marketplace-loading">Loading…</p>
      ) : products.length === 0 ? (
        <p className="marketplace-empty">No products yet. <Link to="/product-create">Add listing</Link></p>
      ) : (
        <ul className="marketplace-list">
          {products.map((p) => (
            <li key={p.id} className="marketplace-item">
              <span className="marketplace-item-name">{p.name ?? 'Product #' + p.id}</span>
              {p.category && <span className="marketplace-item-cat">{p.category}</span>}
              {p.price != null && <span className="marketplace-item-price">{p.price} USD</span>}
              <Link to={`/product/${p.id}`} className="marketplace-item-link">View →</Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
