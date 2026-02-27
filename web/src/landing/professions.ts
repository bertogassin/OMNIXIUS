/** OMNIXIUS professions — replaces js/professions.js */
const STORAGE_KEY = 'omnixius_profession_id';
const STORAGE_NAME_KEY = 'omnixius_profession_name';

export interface Profession {
  id: string;
  icon: string;
  en: string;
  ru: string;
}

export const list: Profession[] = [
  { id: 'cleaning', icon: '\uD83D\uDDA7', en: 'Cleaning', ru: '\u0423\u0431\u043E\u0440\u043A\u0430' },
  { id: 'programming', icon: '\uD83D\uDCBB', en: 'Programming', ru: '\u041F\u0440\u043E\u0433\u0440\u0430\u043C\u043C\u0438\u0440\u043E\u0432\u0430\u043D\u0438\u0435' },
  { id: 'hairdressing', icon: '\u2702\uFE0F', en: 'Hairdressing', ru: '\u041F\u0430\u0440\u0438\u043A\u043C\u0430\u0445\u0435\u0440\u0441\u043A\u0438\u0435 \u0443\u0441\u043B\u0443\u0433\u0438' },
  { id: 'engineering', icon: '\u2699\uFE0F', en: 'Engineering', ru: '\u0418\u043D\u0436\u0435\u043D\u0435\u0440\u0438\u044F' },
  { id: 'tutoring', icon: '\uD83D\uDCDA', en: 'Tutoring', ru: '\u0420\u0435\u043F\u0435\u0442\u0438\u0442\u043E\u0440\u0441\u0442\u0432\u043E' },
  { id: 'delivery', icon: '\uD83D\uDCE6', en: 'Delivery', ru: '\u0414\u043E\u0441\u0442\u0430\u0432\u043A\u0430' },
  { id: 'plumbing', icon: '\uD83D\uDEA7', en: 'Plumbing', ru: '\u0421\u0430\u043D\u0442\u0435\u0445\u043D\u0438\u043A\u0430' },
  { id: 'electrician', icon: '\u26A1', en: 'Electrician', ru: '\u042D\u043B\u0435\u043A\u0442\u0440\u0438\u043A\u0430' },
  { id: 'repair', icon: '\uD83D\uDEE0\uFE0F', en: 'Repair & tech', ru: '\u0420\u0435\u043C\u043E\u043D\u0442 \u0442\u0435\u0445\u043D\u0438\u043A\u0438' },
  { id: 'beauty', icon: '\uD83D\uDC84', en: 'Beauty & cosmetics', ru: '\u041A\u043E\u0441\u043C\u0435\u0442\u043E\u043B\u043E\u0433\u0438\u044F' },
  { id: 'fitness', icon: '\uD83D\uDCAA', en: 'Fitness & training', ru: '\u0424\u0438\u0442\u043D\u0435\u0441, \u0442\u0440\u0435\u043D\u0438\u0440\u043E\u0432\u043A\u0438' },
  { id: 'legal', icon: '\u2696\uFE0F', en: 'Legal', ru: '\u042E\u0440\u0438\u0434\u0438\u0447\u0435\u0441\u043A\u0438\u0435 \u0443\u0441\u043B\u0443\u0433\u0438' },
  { id: 'accounting', icon: '\uD83D\uDCCA', en: 'Accounting', ru: '\u0411\u0443\u0445\u0433\u0430\u043B\u0442\u0435\u0440\u0438\u044F' },
  { id: 'design', icon: '\uD83C\uDFA8', en: 'Design', ru: '\u0414\u0438\u0437\u0430\u0439\u043D' },
  { id: 'photography', icon: '\uD83D\uDCF7', en: 'Photo & video', ru: '\u0424\u043E\u0442\u043E, \u0432\u0438\u0434\u0435\u043E\u0441\u044A\u0451\u043C\u043A\u0430' },
  { id: 'moving', icon: '\uD83D\uDE9A', en: 'Moving & transport', ru: '\u041F\u0435\u0440\u0435\u0435\u0437\u0434\u044B, \u0433\u0440\u0443\u0437\u043E\u043F\u0435\u0440\u0435\u0432\u043E\u0437\u043A\u0438' },
  { id: 'childcare', icon: '\uD83D\uDC76', en: 'Childcare', ru: '\u0423\u0445\u043E\u0434 \u0437\u0430 \u0434\u0435\u0442\u044C\u043C\u0438' },
  { id: 'eldercare', icon: '\uD83D\uDC74', en: 'Elder care', ru: '\u0423\u0445\u043E\u0434 \u0437\u0430 \u043F\u043E\u0436\u0438\u043B\u044B\u043C\u0438' },
  { id: 'cooking', icon: '\uD83C\uDF73', en: 'Cooking & catering', ru: '\u0413\u043E\u0442\u043E\u0432\u043A\u0430, \u043A\u0435\u0439\u0442\u0435\u0440\u0438\u043D\u0433' },
  { id: 'security', icon: '\uD83D\uDEE1\uFE0F', en: 'Security', ru: '\u041E\u0445\u0440\u0430\u043D\u0430, \u0431\u0435\u0437\u043E\u043F\u0430\u0441\u043D\u043E\u0441\u0442\u044C' },
  { id: 'translation', icon: '\uD83C\uDF10', en: 'Translation', ru: '\u041F\u0435\u0440\u0435\u0432\u043E\u0434\u044B' },
  { id: 'marketing', icon: '\uD83D\uDCE2', en: 'Marketing & ads', ru: '\u041C\u0430\u0440\u043A\u0435\u0442\u0438\u043D\u0433, \u0440\u0435\u043A\u043B\u0430\u043C\u0430' },
  { id: 'other', icon: '\uD83D\uDCCC', en: 'Other', ru: '\u0414\u0440\u0443\u0433\u043E\u0435' },
];

function getLang(): string {
  try {
    const w = typeof window !== 'undefined' ? (window as unknown as { omnixius_lang?: string }) : null;
    const stored = typeof localStorage !== 'undefined' ? localStorage.getItem('omnixius_lang') : null;
    return (w?.omnixius_lang || stored || 'en').slice(0, 2);
  } catch {
    return 'en';
  }
}

function getName(p: Profession): string {
  return (getLang() === 'ru' ? p.ru : p.en) || p.en;
}

export const professions = {
  list,
  getList: () => list,
  getName: (id: string) => { const p = list.find((x) => x.id === id); return p ? getName(p) : id; },
  getCurrentId: (): string => {
    try { return (typeof localStorage !== 'undefined' && localStorage.getItem(STORAGE_KEY)) || ''; } catch { return ''; }
  },
  getCurrentName: (): string => {
    const id = professions.getCurrentId();
    if (!id) return '';
    try { return (typeof localStorage !== 'undefined' && localStorage.getItem(STORAGE_NAME_KEY)) || professions.getName(id); } catch { return professions.getName(id); }
  },
  setCurrent: (id: string, name?: string): void => {
    try {
      if (typeof localStorage === 'undefined') return;
      if (id) {
        localStorage.setItem(STORAGE_KEY, id);
        localStorage.setItem(STORAGE_NAME_KEY, name || professions.getName(id));
      } else {
        localStorage.removeItem(STORAGE_KEY);
        localStorage.removeItem(STORAGE_NAME_KEY);
      }
    } catch { /* noop */ }
  },
  getNameFor: (p: Profession) => getName(p),
};

if (typeof window !== 'undefined') (window as unknown as { OMNIXIUS_PROFESSIONS?: typeof professions }).OMNIXIUS_PROFESSIONS = professions;
