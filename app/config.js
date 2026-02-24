// Backend API URL. Local / file: http://localhost:3000 (run "go run ." in backend-go). Production: set on login/register or via api_url in URL.
window.normalizeApiUrl = function (s) {
  if (!s || typeof s !== 'string') return '';
  s = s.trim().replace(/\/+$/, '');
  if (!s) return '';
  if (/^https?:\/\//i.test(s)) return s;
  if (/^localhost(:\d+)?$/i.test(s) || /^127\.0\.0\.1(:\d+)?$/.test(s)) return 'http://' + s;
  return 'https://' + s;
};
(function () {
  if (window.API_URL !== undefined && window.API_URL !== '') return;
  var fromUrl = new URLSearchParams(typeof location !== 'undefined' ? location.search : '').get('api_url');
  var stored = localStorage.getItem('omnixius_api_url');
  if (fromUrl) { window.API_URL = window.normalizeApiUrl(fromUrl); return; }
  if (stored) { window.API_URL = window.normalizeApiUrl(stored); return; }
  var isLocal = !location.hostname || location.hostname === 'localhost' || location.hostname === '127.0.0.1' || location.protocol === 'file:';
  if (isLocal) window.API_URL = 'http://localhost:3000';
  else window.API_URL = '';
})();

// AI service URL (OMNIXIUS AI — root of our own AI). Local: http://localhost:8000
(function () {
  if (window.AI_URL !== undefined && window.AI_URL !== '') return;
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    window.AI_URL = 'http://localhost:8000';
    return;
  }
  window.AI_URL = localStorage.getItem('omnixius_ai_url') || '';
})();
