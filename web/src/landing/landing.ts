/**
 * Landing entry: runs on root HTML pages (index, contact, ecosystem, architecture).
 * Replaces js/professions.js + js/i18n.js + js/main.js + inline profession script.
 */
import { initNav } from './nav';
import { initI18n } from './i18n';
import { professions } from './professions';

function initProfessionGrid(): void {
  const grid = document.getElementById('professionGrid');
  if (!grid) return;

  const currentEl = document.getElementById('professionCurrent');
  const currentNameEl = document.getElementById('professionCurrentName');
  const changeLink = document.getElementById('professionChangeLink');
  const headerLink = document.getElementById('headerProfessionLink');
  const workspaceEl = document.getElementById('professionWorkspaceLink');

  function updateCurrent(): void {
    const id = professions.getCurrentId();
    const name = professions.getCurrentName();
    if (currentNameEl) currentNameEl.textContent = name;
    if (currentEl) currentEl.style.display = id && name ? 'block' : 'none';
    if (workspaceEl) workspaceEl.style.display = id === 'design' ? 'block' : 'none';
    if (headerLink) {
      if (name) {
        headerLink.textContent = name;
        headerLink.style.display = '';
      } else {
        headerLink.style.display = 'none';
      }
    }
  }

  grid.innerHTML = professions.getList()
    .map(
      (p) =>
        `<button type="button" class="profession-card" data-id="${p.id}" data-name="${(professions.getNameFor(p) || '').replace(/"/g, '&quot;')}"><span class="profession-card-icon" aria-hidden="true">${(p.icon || '').replace(/"/g, '&quot;')}</span><span class="profession-card-name">${(professions.getNameFor(p) || '').replace(/</g, '&lt;')}</span></button>`
    )
    .join('');

  grid.querySelectorAll('.profession-card').forEach((btn) => {
    btn.addEventListener('click', () => {
      const id = btn.getAttribute('data-id') || '';
      const name = btn.getAttribute('data-name') || '';
      professions.setCurrent(id, name);
      updateCurrent();
    });
  });

  if (changeLink) {
    changeLink.addEventListener('click', (e) => {
      e.preventDefault();
      document.getElementById('profession')?.scrollIntoView({ behavior: 'smooth' });
    });
  }

  updateCurrent();
}

function run(): void {
  initI18n();
  initNav();
  initProfessionGrid();
}

if (typeof document !== 'undefined' && document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', run);
} else {
  run();
}
