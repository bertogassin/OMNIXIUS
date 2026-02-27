import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { api, ApiError } from '../api';
import './Profession.css';

const PROFESSION_ID = 'cleaning';
const TITLE = 'Cleaning';
const INTRO = 'Find a cleaning professional: home, office, one-time or regular. Filter by distance and rating.';

interface Pro {
  id: number;
  name?: string;
  profession_id?: string;
  online?: boolean;
  rating_avg?: number;
  rating_count?: number;
  distance_km?: number;
}

export default function ProfessionCleaning() {
  const [list, setList] = useState<Pro[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [err, setErr] = useState<string | null>(null);
  const [radiusKm, setRadiusKm] = useState('');
  const [sortBy, setSortBy] = useState('');
  const [lat, setLat] = useState('');
  const [lng, setLng] = useState('');

  const search = () => {
    setLoading(true);
    setErr(null);
    const params: Record<string, string> = { profession: PROFESSION_ID };
    if (radiusKm) params.radius_km = radiusKm;
    if (sortBy) params.sort = sortBy;
    if (lat) params.lat = lat;
    if (lng) params.lng = lng;
    api.professionals.search(params)
      .then((r) => setList((r.professionals as Pro[]) || []))
      .catch((e: ApiError) => setErr(e.data?.error || 'Search failed'))
      .finally(() => setLoading(false));
  };

  const useMyLocation = () => {
    if (!navigator.geolocation) return;
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        setLat(String(pos.coords.latitude));
        setLng(String(pos.coords.longitude));
      },
      () => {}
    );
  };

  useEffect(() => {
    search();
  }, []);

  return (
    <div className="page profession-page">
      <header className="page-header profession-header">
        <p className="page-back"><Link to="/find-professional">← All professions</Link></p>
        <span className="profession-icon" aria-hidden="true">🧹</span>
        <h1>{TITLE}</h1>
        <p className="page-intro">{INTRO}</p>
      </header>

      <div className="page-content profession-content">
        <section className="profession-filters">
          <form onSubmit={(e) => { e.preventDefault(); search(); }} className="profession-form">
            <label className="profession-field">
              <span className="profession-label">Radius (km)</span>
              <input type="number" min={0} value={radiusKm} onChange={(e) => setRadiusKm(e.target.value)} placeholder="Any" className="profession-input" />
            </label>
            <label className="profession-field">
              <span className="profession-label">Sort</span>
              <select value={sortBy} onChange={(e) => setSortBy(e.target.value)} className="profession-select">
                <option value="">Default</option>
                <option value="rating">By rating</option>
                <option value="distance">By distance</option>
              </select>
            </label>
            <button type="button" className="btn btn-outline profession-btn" onClick={useMyLocation}>Use my location</button>
            <button type="submit" className="btn btn-primary profession-btn">Search</button>
          </form>
        </section>

        {err && <p className="page-error">{err}</p>}
        {loading && <p className="page-loading">Loading…</p>}

        {!loading && list !== null && (
          list.length === 0 ? (
            <p className="profession-empty">No cleaning professionals found. Try a larger radius or <Link to="/find-professional">all professions</Link>.</p>
          ) : (
            <ul className="profession-grid">
              {list.map((pro) => (
                <li key={pro.id} className="profession-card">
                  <Link to={`/professional/${pro.id}`} className="profession-card-link">
                    <span className="profession-card-badge">Cleaning</span>
                    <h3 className="profession-card-name">{pro.name || 'Professional #' + pro.id}</h3>
                    <div className="profession-card-meta">
                      {pro.online && <span className="profession-card-online">Online</span>}
                      {pro.rating_avg != null && pro.rating_count != null && (
                        <span className="profession-card-rating">★ {Number(pro.rating_avg).toFixed(1)} ({pro.rating_count})</span>
                      )}
                      {pro.distance_km != null && <span className="profession-card-distance">{pro.distance_km} km</span>}
                    </div>
                    <span className="profession-card-cta">View profile →</span>
                  </Link>
                </li>
              ))}
            </ul>
          )
        )}
      </div>
    </div>
  );
}
