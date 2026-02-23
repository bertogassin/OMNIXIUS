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
      const url = API_URL + path;
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
    },
    users: {
      me: () => api.request('/api/users/me'),
      updateMe: (data) => api.request('/api/users/me', { method: 'PATCH', body: data }),
      myOrders: () => api.request('/api/users/me/orders'),
    },
    products: {
      list: (params) => api.request('/api/products?' + new URLSearchParams(params || {}).toString()),
      get: (id) => api.request('/api/products/' + id),
      categories: () => api.request('/api/products/categories'),
      create: (formData) => api.request('/api/products', { method: 'POST', body: formData, headers: {} }),
      update: (id, formData) => api.request('/api/products/' + id, { method: 'PATCH', body: formData, headers: {} }),
      delete: (id) => api.request('/api/products/' + id, { method: 'DELETE' }),
    },
    orders: {
      my: () => api.request('/api/orders/my'),
      create: (product_id) => api.request('/api/orders', { method: 'POST', body: { product_id } }),
      update: (id, status) => api.request('/api/orders/' + id, { method: 'PATCH', body: { status } }),
    },
    conversations: {
      list: () => api.request('/api/conversations'),
      create: (user_id, product_id) => api.request('/api/conversations', { method: 'POST', body: { user_id, product_id } }),
    },
    messages: {
      list: (conversationId) => api.request('/api/messages/conversation/' + conversationId),
      send: (conversationId, body) => api.request('/api/messages/conversation/' + conversationId, { method: 'POST', body: { body } }),
      read: (id) => api.request('/api/messages/' + id + '/read', { method: 'POST' }),
    },
  };
})();
