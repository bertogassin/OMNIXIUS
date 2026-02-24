// Backend API URL. Local: http://localhost:3000
// On GitHub Pages / production: set via "Set API URL" on login/register, or from Go redirect (api_url in URL).
(function () {
  if (window.API_URL !== undefined && window.API_URL !== '') return;
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    window.API_URL = 'http://localhost:3000';
    return;
  }
  var fromUrl = new URLSearchParams(typeof location !== 'undefined' ? location.search : '').get('api_url');
  window.API_URL = fromUrl || localStorage.getItem('omnixius_api_url') || '';
})();
