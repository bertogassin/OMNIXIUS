(function () {
  var menuToggle = document.querySelector('.menu-toggle');
  var nav = document.querySelector('.nav');

  if (menuToggle && nav) {
    menuToggle.type = 'button';
    menuToggle.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      nav.classList.toggle('open');
      closeAllDrops();
    });
    nav.querySelectorAll('a').forEach(function (a) {
      a.addEventListener('click', function () {
        nav.classList.remove('open');
        closeAllDrops();
      });
    });
  }

  function closeAllDrops() {
    document.querySelectorAll('.nav-drop.open').forEach(function (d) {
      d.classList.remove('open');
      var t = d.querySelector('.nav-drop-trigger');
      if (t) t.setAttribute('aria-expanded', 'false');
    });
  }

  document.querySelectorAll('.nav-drop').forEach(function (drop) {
    var trigger = drop.querySelector('.nav-drop-trigger');
    var panel = drop.querySelector('.nav-drop-panel');
    if (!trigger || !panel) return;
    trigger.addEventListener('click', function (e) {
      e.preventDefault();
      e.stopPropagation();
      var isOpen = drop.classList.toggle('open');
      trigger.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
    });
  });

  document.addEventListener('click', function (e) {
    if (!e.target.closest('.nav-drop')) closeAllDrops();
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') {
      closeAllDrops();
      if (nav && nav.classList.contains('open')) nav.classList.remove('open');
    }
  });
})();
