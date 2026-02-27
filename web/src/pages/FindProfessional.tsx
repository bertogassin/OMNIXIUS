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
    <div>
      <h1>Find professional</h1>
      <section style={{ marginBottom: 24 }}>
        <form onSubmit={(e) => { e.preventDefault(); search(); }} style={{ display: 'flex', flexWrap: 'wrap', gap: 16, alignItems: 'flex-end' }}>
          <div>
            <label>Profession</label>
            <input value={profession} onChange={(e) => setProfession(e.target.value)} placeholder="All" style={{ marginLeft: 8 }} />
          </div>
          <label><input type="checkbox" checked={onlineOnly} onChange={(e) => setOnlineOnly(e.target.checked)} /> Online only</label>
          <div>
            <label>Radius km</label>
            <input type="number" min={0} value={radiusKm} onChange={(e) => setRadiusKm(e.target.value)} style={{ width: 80, marginLeft: 8 }} />
          </div>
          <div>
            <label>Sort</label>
            <select value={sortBy} onChange={(e) => setSortBy(e.target.value)} style={{ marginLeft: 8 }}>
              <option value="">Default</option>
              <option value="rating">By rating</option>
              <option value="distance">By distance</option>
            </select>
          </div>
          <button type="button" className="btn btn-outline" onClick={useMyLocation}>Use my location</button>
          <button type="submit" className="btn btn-primary">Search</button>
        </form>
      </section>
      {loading && <p>Searching…</p>}
      {err && <p className="error">{err}</p>}
      {list && !loading && (
        <ul style={{ listStyle: 'none', padding: 0 }}>
          {list.length === 0 ? <p>No professionals found.</p> : list.map((u) => (
            <li key={u.id} style={{ padding: '12px 0', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', flexWrap: 'wrap' }}>
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
      )}
    </div>
  );
}
