import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { api, ApiError } from '../api';
import './Profession.css';

interface Pro {
  id: number;
  name?: string;
  profession_id?: string;
  online?: boolean;
  rating_avg?: number;
  rating_count?: number;
  distance_km?: number;
  bio?: string;
}

export default function Professional() {
  const { id } = useParams<{ id: string }>();
  const [pro, setPro] = useState<Pro | null>(null);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;
    setLoading(true);
    api.professionals
      .get(id)
      .then((r) => setPro(r as Pro))
      .catch((e: ApiError) => {
        if (e.status === 404) setPro(null);
        else setErr(e.data?.error || 'Failed to load');
      })
      .finally(() => setLoading(false));
  }, [id]);

  if (!id) return <div className="page"><p className="page-error">No professional ID.</p></div>;
  if (err) return <div className="page"><p className="page-error">{err}</p><p className="page-back"><Link to="/profession/cleaning">← Cleaning</Link></p></div>;
  if (loading) return <div className="page profession-page"><h1>Professional</h1><p className="page-loading">Loading…</p></div>;
  if (!pro) {
    return (
      <div className="page profession-page">
        <p className="page-back"><Link to="/profession/cleaning">← Cleaning</Link></p>
        <p className="page-error">Professional not found.</p>
        <p><Link to="/marketplace">Browse marketplace</Link> or <Link to="/find-professional">find professional</Link>.</p>
      </div>
    );
  }

  const professionLabel = (pro.profession_id === 'cleaning' && 'Cleaning') || pro.profession_id || 'Professional';

  return (
    <div className="page profession-page professional-detail">
      <p className="page-back"><Link to="/profession/cleaning">← Cleaning</Link></p>

      <header className="profession-detail-header">
        <span className="profession-icon profession-icon-lg" aria-hidden="true">🧹</span>
        <div className="profession-detail-heading">
          <span className="profession-card-badge">{professionLabel}</span>
          <h1>{pro.name || 'Professional #' + pro.id}</h1>
          <div className="profession-card-meta">
            {pro.online && <span className="profession-card-online">Online</span>}
            {pro.rating_avg != null && pro.rating_count != null && (
              <span className="profession-card-rating">★ {Number(pro.rating_avg).toFixed(1)} ({pro.rating_count} reviews)</span>
            )}
            {pro.distance_km != null && <span className="profession-card-distance">{pro.distance_km} km away</span>}
          </div>
        </div>
      </header>

      <div className="page-content profession-content">
        {pro.bio && <section className="profession-section"><h2>About</h2><p>{pro.bio}</p></section>}

        <section className="profession-actions">
          <Link to={`/marketplace?category=${pro.profession_id || 'cleaning'}`} className="btn btn-primary">
            View listings
          </Link>
          <Link to="/mail" className="btn btn-outline">Contact</Link>
        </section>
      </div>
    </div>
  );
}
