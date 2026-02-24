(function () {
  if (!window.api || !api.getToken()) return;
  api.conversations.unreadCount().then(function (r) {
    var n = r && r.unread ? parseInt(r.unread, 10) : 0;
    if (n <= 0) return;
    var el = document.querySelector('nav a[href="mail.html"]');
    if (!el) return;
    var badge = el.querySelector('.nav-unread');
    if (!badge) {
      badge = document.createElement('span');
      badge.className = 'nav-unread';
      el.appendChild(badge);
    }
    badge.textContent = n > 99 ? '99+' : String(n);
  }).catch(function () {});
})();
