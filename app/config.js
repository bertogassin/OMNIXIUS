// Backend API URL. Local: http://localhost:3000
// On GitHub Pages / production: set via "Set API URL" on login/register, or set window.API_URL before config.js loads.
(function () {
  if (window.API_URL !== undefined && window.API_URL !== '') return;
  if (location.hostname === 'localhost' || location.hostname === '127.0.0.1') {
    window.API_URL = 'http://localhost:3000';
    return;
  }
  window.API_URL = localStorage.getItem('omnixius_api_url') || '';
})();
