(function () {
  const API_URL = window.API_URL || '';
  const TOKEN_KEY = 'omnixius_token';
  const USER_KEY = 'omnixius_user';
  const PERSIST_KEY = 'omnixius_remember';

  function getStorage(persistent) {
    return persistent ? localStorage : sessionStorage;
  }

  function getToken() {
    return sessionStorage.getItem(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY);
  }

  function setToken(token, persistent) {
    sessionStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(TOKEN_KEY);
    if (token) {
      var storage = persistent !== false ? localStorage : sessionStorage;
      storage.setItem(TOKEN_KEY, token);
      try { storage.setItem(PERSIST_KEY, persistent !== false ? '1' : '0'); } catch (_) {}
    } else {
      try { localStorage.removeItem(PERSIST_KEY); sessionStorage.removeItem(PERSIST_KEY); } catch (_) {}
    }
  }

  function getAuthHeaders() {
    const token = getToken();
    return {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    };
  }

  window.api = {
    getToken,
    setToken,
    get user() {
      try {
        var raw = sessionStorage.getItem(USER_KEY) || localStorage.getItem(USER_KEY);
        return raw ? JSON.parse(raw) : null;
      } catch (_) {
        return null;
      }
    },
    set user(u) {
      var storage = sessionStorage.getItem(TOKEN_KEY) ? sessionStorage : localStorage;
      sessionStorage.removeItem(USER_KEY);
      localStorage.removeItem(USER_KEY);
      if (u) storage.setItem(USER_KEY, JSON.stringify(u));
    },
    async request(path, options = {}) {
      const base = window.API_URL || API_URL || '';
      const isLocal = (typeof location !== 'undefined') && (!location.hostname || location.hostname === 'localhost' || location.hostname === '127.0.0.1' || location.protocol === 'file:');
      if (!base && !isLocal) {
        var err = new Error('API URL not set');
        err.status = 0;
        err.data = { error: 'Set API URL on login/register page. The app cannot call the backend until you set it.' };
        throw err;
      }
      const url = (base || 'http://localhost:3000') + path;
      const headers = { ...getAuthHeaders(), ...options.headers };
      if (options.body && typeof options.body === 'object' && !(options.body instanceof FormData)) {
        headers['Content-Type'] = 'application/json';
        options.body = JSON.stringify(options.body);
      }
      if (options.body instanceof FormData) delete headers['Content-Type'];
      const res = await fetch(url, { ...options, headers });
      const text = await res.text();
      let data;
      try {
        data = text ? JSON.parse(text) : null;
      } catch (_) {
        data = null;
      }
      if (!res.ok) throw { status: res.status, data: data || { error: text } };
      return data;
    },
    auth: {
      register: (email, password, name) => api.request('/api/auth/register', { method: 'POST', body: { email, password, name } }),
      login: (email, password) => api.request('/api/auth/login', { method: 'POST', body: { email, password } }),
      logout: () => { api.setToken(null); api.user = null; },
      forgotPassword: (email) => api.request('/api/auth/forgot-password', { method: 'POST', body: { email } }),
      changePassword: (currentPassword, newPassword) => api.request('/api/auth/change-password', { method: 'POST', body: { current_password: currentPassword, new_password: newPassword } }),
    },
    users: {
      me: () => api.request('/api/users/me'),
      get: (id) => api.request('/api/users/' + id),
      updateMe: (data) => api.request('/api/users/me', { method: 'PATCH', body: data }),
      myOrders: () => api.request('/api/users/me/orders'),
      balance: () => api.request('/api/users/me/balance'),
      balanceCredit: (amount) => api.request('/api/users/me/balance/credit', { method: 'POST', body: { amount } }),
    },
    subscriptions: {
      create: (product_id) => api.request('/api/subscriptions', { method: 'POST', body: { product_id } }),
      my: () => api.request('/api/subscriptions/my'),
    },
    products: {
      list: (params) => api.request('/api/products?' + new URLSearchParams(params || {}).toString()),
      get: (id) => api.request('/api/products/' + id),
      categories: () => api.request('/api/products/categories'),
      create: (formData) => api.request('/api/products', { method: 'POST', body: formData, headers: {} }),
      update: (id, formData) => api.request('/api/products/' + id, { method: 'PATCH', body: formData, headers: {} }),
      delete: (id) => api.request('/api/products/' + id, { method: 'DELETE' }),
      slots: (productId) => api.request('/api/products/' + productId + '/slots'),
      addSlot: (productId, slot_at) => api.request('/api/products/' + productId + '/slots', { method: 'POST', body: { slot_at } }),
      bookSlot: (productId, slotId) => api.request('/api/products/' + productId + '/slots/' + slotId + '/book', { method: 'POST' }),
      closedContent: (productId) => api.request('/api/products/' + productId + '/closed-content'),
    },
    orders: {
      my: () => api.request('/api/orders/my'),
      create: (product_id, data) => api.request('/api/orders', { method: 'POST', body: Object.assign({ product_id }, data || {}) }),
      update: (id, data) => api.request('/api/orders/' + id, { method: 'PATCH', body: typeof data === 'string' ? { status: data } : (data || {}) }),
    },
    remittances: {
      my: () => api.request('/api/remittances/my'),
      create: (to_identifier, amount, currency) => api.request('/api/remittances', { method: 'POST', body: { to_identifier, amount, currency: currency || 'USD' } }),
    },
    conversations: {
      list: () => api.request('/api/conversations'),
      unreadCount: () => api.request('/api/conversations/unread-count'),
      create: (user_id, product_id) => api.request('/api/conversations', { method: 'POST', body: { user_id, product_id } }),
    },
    messages: {
      list: (conversationId) => api.request('/api/messages/conversation/' + conversationId),
      send: (conversationId, body) => api.request('/api/messages/conversation/' + conversationId, { method: 'POST', body: { body } }),
      read: (id) => api.request('/api/messages/' + id + '/read', { method: 'POST' }),
    },
  };
})();
