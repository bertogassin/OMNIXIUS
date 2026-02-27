/**
 * API base URL. Set via env VITE_API_URL (e.g. in .env: VITE_API_URL=http://localhost:3000)
 * or at runtime via window.__OMNIXIUS_API_URL__ for same-origin static app compatibility.
 */
export const API_URL =
  (typeof window !== 'undefined' && (window as unknown as { __OMNIXIUS_API_URL__?: string }).__OMNIXIUS_API_URL__) ||
  import.meta.env.VITE_API_URL ||
  (typeof localStorage !== 'undefined' && localStorage.getItem('omnixius_api_url')) ||
  (typeof window !== 'undefined' &&
    (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
    ? 'http://localhost:3000'
    : '');
