/**
 * OMNIXIUS — professions (service categories). User can switch at any time on the main page.
 * Stored in localStorage as omnixius_profession_id; used to filter marketplace and show context.
 */
(function () {
  const STORAGE_KEY = 'omnixius_profession_id';
  const STORAGE_NAME_KEY = 'omnixius_profession_name';

  var list = [
    { id: 'cleaning', en: 'Cleaning', ru: 'Уборка' },
    { id: 'programming', en: 'Programming', ru: 'Программирование' },
    { id: 'hairdressing', en: 'Hairdressing', ru: 'Парикмахерские услуги' },
    { id: 'engineering', en: 'Engineering', ru: 'Инженерия' },
    { id: 'tutoring', en: 'Tutoring', ru: 'Репетиторство' },
    { id: 'delivery', en: 'Delivery', ru: 'Доставка' },
    { id: 'plumbing', en: 'Plumbing', ru: 'Сантехника' },
    { id: 'electrician', en: 'Electrician', ru: 'Электрика' },
    { id: 'repair', en: 'Repair & tech', ru: 'Ремонт техники' },
    { id: 'beauty', en: 'Beauty & cosmetics', ru: 'Косметология' },
    { id: 'fitness', en: 'Fitness & training', ru: 'Фитнес, тренировки' },
    { id: 'legal', en: 'Legal', ru: 'Юридические услуги' },
    { id: 'accounting', en: 'Accounting', ru: 'Бухгалтерия' },
    { id: 'design', en: 'Design', ru: 'Дизайн' },
    { id: 'photography', en: 'Photo & video', ru: 'Фото, видеосъёмка' },
    { id: 'moving', en: 'Moving & transport', ru: 'Переезды, грузоперевозки' },
    { id: 'childcare', en: 'Childcare', ru: 'Уход за детьми' },
    { id: 'eldercare', en: 'Elder care', ru: 'Уход за пожилыми' },
    { id: 'cooking', en: 'Cooking & catering', ru: 'Готовка, кейтеринг' },
    { id: 'security', en: 'Security', ru: 'Охрана, безопасность' },
    { id: 'translation', en: 'Translation', ru: 'Переводы' },
    { id: 'marketing', en: 'Marketing & ads', ru: 'Маркетинг, реклама' },
    { id: 'other', en: 'Other', ru: 'Другое' }
  ];

  function getLang() {
    try {
      return (window.omnixius_lang || localStorage.getItem('omnixius_lang') || 'en').slice(0, 2);
    } catch (_) { return 'en'; }
  }

  function getName(p) {
    var lang = getLang();
    return (lang === 'ru' ? p.ru : p.en) || p.en;
  }

  window.OMNIXIUS_PROFESSIONS = {
    list: list,
    getList: function () { return list; },
    getName: function (id) {
      var p = list.find(function (x) { return x.id === id; });
      return p ? getName(p) : id;
    },
    getCurrentId: function () {
      try { return localStorage.getItem(STORAGE_KEY) || ''; } catch (_) { return ''; }
    },
    getCurrentName: function () {
      var id = this.getCurrentId();
      if (!id) return '';
      try { return localStorage.getItem(STORAGE_NAME_KEY) || this.getName(id); } catch (_) { return this.getName(id); }
    },
    setCurrent: function (id, name) {
      try {
        if (id) {
          localStorage.setItem(STORAGE_KEY, id);
          localStorage.setItem(STORAGE_NAME_KEY, name || this.getName(id));
        } else {
          localStorage.removeItem(STORAGE_KEY);
          localStorage.removeItem(STORAGE_NAME_KEY);
        }
      } catch (_) {}
    },
    getNameFor: getName
  };
})();
