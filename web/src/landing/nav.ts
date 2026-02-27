/**
 * Landing nav: menu toggle and dropdowns. Replaces js/main.js behaviour.
 */
function closeAllDrops(): void {
  document.querySelectorAll('.nav-drop.open').forEach((d) => {
    d.classList.remove('open');
    const t = d.querySelector('.nav-drop-trigger');
    if (t) t.setAttribute('aria-expanded', 'false');
  });
}

export function initNav(): void {
  const menuToggle = document.querySelector('.menu-toggle');
  const nav = document.querySelector('.nav');

  if (menuToggle && nav) {
    (menuToggle as HTMLButtonElement).type = 'button';
    menuToggle.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      nav.classList.toggle('open');
      closeAllDrops();
    });
    nav.querySelectorAll('a').forEach((a) => {
      a.addEventListener('click', () => {
        nav.classList.remove('open');
        closeAllDrops();
      });
    });
  }

  document.querySelectorAll('.nav-drop').forEach((drop) => {
    const trigger = drop.querySelector('.nav-drop-trigger');
    if (!trigger) return;
    trigger.addEventListener('click', (e) => {
      e.preventDefault();
      e.stopPropagation();
      const isOpen = drop.classList.toggle('open');
      trigger.setAttribute('aria-expanded', isOpen ? 'true' : 'false');
    });
  });

  document.addEventListener('click', (e) => {
    if (!(e.target as Element).closest?.('.nav-drop')) closeAllDrops();
  });

  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      closeAllDrops();
      if (nav?.classList.contains('open')) nav.classList.remove('open');
    }
  });
}
