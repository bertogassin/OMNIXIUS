/**
 * API base URL. Set via env VITE_API_URL, window.__OMNIXIUS_API_URL__, or localStorage omnixius_api_url.
 * Default: localhost:3000 or same host:3000 when opened from LAN (e.g. phone at 192.168.1.10:5173 → API 192.168.1.10:3000).
 */
function defaultApiUrl(): string {
  if (typeof window === 'undefined') return '';
  const h = window.location.hostname;
  if (h === 'localhost' || h === '127.0.0.1') return 'http://localhost:3000';
  // LAN access (phone, tablet): same host, port 3000
  if (/^192\.168\.|^10\.|^172\.(1[6-9]|2[0-9]|3[01])\./.test(h)) return `http://${h}:3000`;
  return '';
}

export const API_URL =
  (typeof window !== 'undefined' && (window as unknown as { __OMNIXIUS_API_URL__?: string }).__OMNIXIUS_API_URL__) ||
  import.meta.env.VITE_API_URL ||
  (typeof localStorage !== 'undefined' && localStorage.getItem('omnixius_api_url')) ||
  defaultApiUrl();
