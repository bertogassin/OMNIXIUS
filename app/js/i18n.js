(function () {
  const STORAGE_KEY = 'omnixius_lang';
  const defaultLang = 'en';
  const strings = {
    en: {
      navMarketplace: 'Marketplace',
      navMail: 'Mail',
      navOrders: 'Orders',
      navProfile: 'Profile',
      navSite: 'Site',
      signIn: 'Sign in',
      signOut: 'Sign out',
      register: 'Register',
      ordersTitle: 'My orders',
      ordersHint: 'Purchases and sales in one place.',
      ordersAsBuyer: 'As buyer',
      ordersAsSeller: 'As seller',
      noOrders: 'No orders yet.',
      setApiUrl: 'Set API URL',
      apiUrlHint: 'Backend URL not set. You cannot sign in or register until you set it.',
      email: 'Email',
      password: 'Password',
      passwordMin: 'Password (min 8 characters)',
      name: 'Name',
      loginTitle: 'Sign in',
      registerTitle: 'Register',
      forgotPassword: 'Forgot password?',
      noAccount: 'No account? Register',
      haveAccount: 'Already have an account? Sign in',
      signInFailed: 'Sign in failed',
      registerFailed: 'Registration failed',
      cannotConnect: 'Cannot reach server. Set API URL?',
    },
    ru: {
      navMarketplace: 'Маркетплейс',
      navMail: 'Почта',
      navOrders: 'Заказы',
      navProfile: 'Профиль',
      navSite: 'Сайт',
      signIn: 'Войти',
      signOut: 'Выйти',
      register: 'Регистрация',
      ordersTitle: 'Мои заказы',
      ordersHint: 'Покупки и продажи в одном месте.',
      ordersAsBuyer: 'Как покупатель',
      ordersAsSeller: 'Как продавец',
      noOrders: 'Пока заказов нет.',
      setApiUrl: 'Указать URL API',
      apiUrlHint: 'URL бэкенда не задан. Вход и регистрация недоступны.',
      email: 'Email',
      password: 'Пароль',
      passwordMin: 'Пароль (не менее 8 символов)',
      name: 'Имя',
      loginTitle: 'Вход',
      registerTitle: 'Регистрация',
      forgotPassword: 'Забыли пароль?',
      noAccount: 'Нет аккаунта? Зарегистрироваться',
      haveAccount: 'Уже есть аккаунт? Войти',
      signInFailed: 'Ошибка входа',
      registerFailed: 'Ошибка регистрации',
      cannotConnect: 'Нет связи с сервером. Указать URL API?',
    },
    fr: {
      navMarketplace: 'Marketplace',
      navMail: 'Mail',
      navOrders: 'Commandes',
      navProfile: 'Profil',
      navSite: 'Site',
      signIn: 'Connexion',
      signOut: 'Déconnexion',
      register: 'Inscription',
      ordersTitle: 'Mes commandes',
      ordersHint: 'Achats et ventes au même endroit.',
      ordersAsBuyer: 'En tant qu\'acheteur',
      ordersAsSeller: 'En tant que vendeur',
      noOrders: 'Aucune commande pour l\'instant.',
      setApiUrl: 'Définir l\'URL API',
      apiUrlHint: 'URL du backend non définie. Connexion et inscription indisponibles.',
      email: 'Email',
      password: 'Mot de passe',
      passwordMin: 'Mot de passe (min. 8 caractères)',
      name: 'Nom',
      loginTitle: 'Connexion',
      registerTitle: 'Inscription',
      forgotPassword: 'Mot de passe oublié ?',
      noAccount: 'Pas de compte ? S\'inscrire',
      haveAccount: 'Déjà un compte ? Connexion',
      signInFailed: 'Échec de la connexion',
      registerFailed: 'Échec de l\'inscription',
      cannotConnect: 'Impossible de joindre le serveur. Définir l\'URL API ?',
    },
  };

  function getLang() {
    return localStorage.getItem(STORAGE_KEY) || defaultLang;
  }
  function setLang(lang) {
    if (strings[lang]) {
      localStorage.setItem(STORAGE_KEY, lang);
      document.documentElement.lang = lang;
      if (typeof apply === 'function') apply();
      return true;
    }
    return false;
  }
  function t(key) {
    const lang = getLang();
    return (strings[lang] && strings[lang][key]) || (strings[defaultLang] && strings[defaultLang][key]) || key;
  }
  function apply() {
    document.querySelectorAll('[data-i18n]').forEach(function (el) {
      const key = el.getAttribute('data-i18n');
      if (key && t(key) !== key) el.textContent = t(key);
    });
    document.querySelectorAll('[data-i18n-placeholder]').forEach(function (el) {
      const key = el.getAttribute('data-i18n-placeholder');
      if (key) el.placeholder = t(key);
    });
    const cur = getLang();
    document.querySelectorAll('.lang-switcher a[data-lang]').forEach(function (a) {
      a.classList.toggle('active', a.getAttribute('data-lang') === cur);
    });
  }

  window.i18n = {
    getLang,
    setLang,
    t,
    apply,
    defaultLang,
    supported: ['en', 'ru', 'fr'],
  };
  document.documentElement.lang = getLang();
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', function () {
      apply();
      document.querySelectorAll('.lang-switcher a[data-lang]').forEach(function (a) {
        a.addEventListener('click', function (e) { e.preventDefault(); setLang(a.getAttribute('data-lang')); apply(); });
      });
    });
  } else { apply(); }
})();
