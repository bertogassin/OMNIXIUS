import { useState } from 'react';
import { Link } from 'react-router-dom';
import { api, ApiError } from '../api';

interface Pro {
  id: number;
  name?: string;
  profession_id?: string;
  online?: boolean;
  rating_avg?: number;
  rating_count?: number;
  distance_km?: number;
}

export default function FindProfessional() {
  const [profession, setProfession] = useState('');
  const [onlineOnly, setOnlineOnly] = useState(false);
  const [radiusKm, setRadiusKm] = useState('');
  const [sortBy, setSortBy] = useState('');
  const [list, setList] = useState<Pro[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string | null>(null);
  const [lat, setLat] = useState('');
  const [lng, setLng] = useState('');

  const search = () => {
    setLoading(true);
    setErr(null);
    const params: Record<string, string> = {};
    if (profession) params.profession = profession;
    if (onlineOnly) params.online = '1';
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
      (pos) => { setLat(String(pos.coords.latitude)); setLng(String(pos.coords.longitude)); },
      () => {}
    );
  };

  return (
    <div className="page">
      <header className="page-header">
        <h1>Find professional</h1>
        <p className="page-intro">Search by profession, online, distance.</p>
        <p style={{ marginTop: '0.75rem' }}>
          <Link to="/profession/cleaning" className="btn btn-outline" style={{ fontSize: '0.9rem' }}>🧹 Cleaning — dedicated screen</Link>
        </p>
      </header>
      <div className="page-content">
        <section className="page-section">
          <form onSubmit={(e) => { e.preventDefault(); search(); }} className="page-form-row">
            <label>Profession <input value={profession} onChange={(e) => setProfession(e.target.value)} placeholder="All" /></label>
            <label><input type="checkbox" checked={onlineOnly} onChange={(e) => setOnlineOnly(e.target.checked)} /> Online only</label>
            <label>Radius km <input type="number" min={0} value={radiusKm} onChange={(e) => setRadiusKm(e.target.value)} style={{ width: 80 }} /></label>
            <label>Sort <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}><option value="">Default</option><option value="rating">By rating</option><option value="distance">By distance</option></select></label>
            <button type="button" className="btn btn-outline" onClick={useMyLocation}>Use my location</button>
            <button type="submit" className="btn btn-primary">Search</button>
          </form>
        </section>
        {loading && <p className="page-loading">Searching…</p>}
        {err && <p className="page-error">{err}</p>}
        {list != null && !loading && (
          list.length === 0 ? <p>No professionals found.</p> : (
            <ul className="page-list">
              {list.map((u) => (
                <li key={u.id} className="page-list-item">
                  <span>
                    {u.name || 'User #' + u.id}
                    {u.online ? <span style={{ background: 'var(--accent)', fontSize: 12, padding: '2px 6px', borderRadius: 4, marginLeft: 8 }}>Online</span> : null}
                    {u.rating_avg != null && u.rating_count ? ' ★ ' + Number(u.rating_avg).toFixed(1) + ' (' + u.rating_count + ')' : ''}
                    {u.distance_km != null ? ' · ' + u.distance_km + ' km' : ''}
                  </span>
                  <Link to={'/marketplace?category=' + (u.profession_id || '')}>View</Link>
                </li>
              ))}
            </ul>
          )
        )}
      </div>
    </div>
  );
}
