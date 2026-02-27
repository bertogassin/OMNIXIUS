import { strings } from './strings';

const STORAGE_KEY = 'omnixius_lang';
const defaultLang = 'en';

export function getLang(): string {
  if (typeof localStorage === 'undefined') return defaultLang;
  return localStorage.getItem(STORAGE_KEY) || defaultLang;
}

export function setLang(lang: string): boolean {
  if (!strings[lang]) return false;
  if (typeof localStorage !== 'undefined') localStorage.setItem(STORAGE_KEY, lang);
  if (typeof document !== 'undefined') document.documentElement.lang = lang;
  apply();
  return true;
}

export function t(key: string): string {
  const lang = getLang();
  const L = strings[lang] || strings[defaultLang];
  return (L && L[key]) || (strings[defaultLang] && strings[defaultLang][key]) || key;
}

export function apply(): void {
  if (typeof document === 'undefined') return;
  document.querySelectorAll<HTMLElement>('[data-i18n]').forEach((el) => {
    const key = el.getAttribute('data-i18n');
    if (key && t(key) !== key) el.textContent = t(key);
  });
  document.querySelectorAll<HTMLInputElement | HTMLTextAreaElement>('[data-i18n-placeholder]').forEach((el) => {
    const key = el.getAttribute('data-i18n-placeholder');
    if (key) el.placeholder = t(key);
  });
  const cur = getLang();
  document.querySelectorAll('.lang-switcher a[data-lang]').forEach((a) => {
    a.classList.toggle('active', a.getAttribute('data-lang') === cur);
  });
}

export function bindLangSwitcher(): void {
  if (typeof document === 'undefined') return;
  document.body.addEventListener('click', (e) => {
    const a = (e.target as Element)?.closest?.('.lang-switcher a[data-lang]');
    if (a) {
      e.preventDefault();
      setLang(a.getAttribute('data-lang') || 'en');
    }
  }, true);
}

export function initI18n(): void {
  if (typeof document === 'undefined') return;
  document.documentElement.lang = getLang();
  apply();
  bindLangSwitcher();
}

if (typeof window !== 'undefined') {
  (window as unknown as { i18n?: { getLang: typeof getLang; setLang: typeof setLang; t: typeof t; apply: typeof apply } }).i18n = { getLang, setLang, t, apply };
}
