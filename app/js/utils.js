/**
 * OMNIXIUS — общие утилиты (один раз, без дубликатов в страницах).
 */
(function () {
  function escapeHtml(s) {
    return (s || '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }
  function formatDate(ts, options) {
    if (!ts) return '';
    var d = new Date(ts * 1000);
    return options && options.dateOnly ? d.toLocaleDateString(undefined, { dateStyle: 'medium' }) : d.toLocaleString();
  }
  window.escapeHtml = escapeHtml;
  window.formatDate = formatDate;
})();
