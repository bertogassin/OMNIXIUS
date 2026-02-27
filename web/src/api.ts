import { API_URL } from './config';

const TOKEN_KEY = 'omnixius_token';
const USER_KEY = 'omnixius_user';

export function getToken(): string | null {
  return sessionStorage.getItem(TOKEN_KEY) || localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string | null, persistent = true): void {
  sessionStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(TOKEN_KEY);
  if (token) {
    const storage = persistent ? localStorage : sessionStorage;
    storage.setItem(TOKEN_KEY, token);
  }
  sessionStorage.removeItem(USER_KEY);
  localStorage.removeItem(USER_KEY);
}

export function getUser(): { id: number; email?: string; name?: string; role?: string } | null {
  try {
    const raw = sessionStorage.getItem(USER_KEY) || localStorage.getItem(USER_KEY);
    return raw ? JSON.parse(raw) : null;
  } catch {
    return null;
  }
}

export function setUser(u: { id: number; email?: string; name?: string; role?: string } | null): void {
  sessionStorage.removeItem(USER_KEY);
  localStorage.removeItem(USER_KEY);
  const token = getToken();
  const storage = token && sessionStorage.getItem(TOKEN_KEY) ? sessionStorage : localStorage;
  if (u) storage.setItem(USER_KEY, JSON.stringify(u));
}

export interface ApiError {
  status: number;
  data: { error?: string };
}

type RequestOptions = Omit<RequestInit, 'body'> & { body?: Record<string, unknown> | FormData };

async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const base = API_URL;
  if (!base && typeof window !== 'undefined') {
    const err = new Error('API URL not set') as Error & ApiError;
    err.status = 0;
    err.data = { error: 'Set VITE_API_URL in .env or window.__OMNIXIUS_API_URL__' };
    throw err;
  }
  const url = (base || 'http://localhost:3000') + path;
  const headers: HeadersInit = {
    ...(options.body instanceof FormData ? {} : { 'Content-Type': 'application/json' }),
    ...(getToken() ? { Authorization: `Bearer ${getToken()}` } : {}),
    ...(options.headers as Record<string, string>),
  };
  let body: BodyInit | undefined;
  if (options.body && typeof options.body === 'object' && !(options.body instanceof FormData)) {
    body = JSON.stringify(options.body);
  } else {
    body = options.body as BodyInit | undefined;
  }
  const { body: _b, ...rest } = options;
  const res = await fetch(url, { ...rest, headers, body });
  const text = await res.text();
  let data: unknown;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    data = null;
  }
  if (!res.ok) {
    const e: ApiError = { status: res.status, data: (data as { error?: string }) || { error: text } };
    throw e;
  }
  return data as T;
}

export const api = {
  getToken,
  setToken,
  get user() { return getUser(); },
  set user(u) { setUser(u); },
  request: request as <T>(path: string, options?: RequestOptions) => Promise<T>,
  auth: {
    login: (email: string, password: string) =>
      request<{ token: string; user: { id: number; email?: string; name?: string; role?: string } }>(
        '/api/auth/login',
        { method: 'POST', body: { email, password } }
      ),
    register: (email: string, password: string, name: string) =>
      request<{ token: string; user: unknown }>('/api/auth/register', { method: 'POST', body: { email, password, name } }),
    logout: () => { setToken(null); setUser(null); },
    forgotPassword: (email: string) =>
      request<unknown>('/api/auth/forgot-password', { method: 'POST', body: { email } }),
    resetPassword: (token: string, password: string) =>
      request<unknown>('/api/auth/reset-password', { method: 'POST', body: { token, password } }),
    changePassword: (currentPassword: string, newPassword: string) =>
      request<unknown>('/api/auth/change-password', { method: 'POST', body: { current_password: currentPassword, new_password: newPassword } }),
    sessions: () => request<unknown>('/api/auth/sessions'),
    sessionDelete: (id: string) => request<unknown>('/api/auth/sessions/' + id, { method: 'DELETE' }),
    devices: () => request<unknown>('/api/auth/devices'),
    deviceDelete: (id: string) => request<unknown>('/api/auth/devices/' + id, { method: 'DELETE' }),
    recoveryGenerate: (recoveryHash: string) =>
      request<unknown>('/api/auth/recovery/generate', { method: 'POST', body: { recoveryHash } }),
  },
  users: {
    me: () => request<{ id: number; email?: string; name?: string; role?: string }>('/api/users/me'),
    get: (id: string) => request<unknown>('/api/users/' + id),
    updateMe: (data: Record<string, unknown>) => request<unknown>('/api/users/me', { method: 'PATCH', body: data }),
    heartbeat: () => request<unknown>('/api/users/me/heartbeat', { method: 'POST' }),
    myOrders: () => request<{ asBuyer?: unknown[]; asSeller?: unknown[] }>('/api/users/me/orders'),
    balance: () => request<unknown>('/api/users/me/balance'),
    balanceCredit: (amount: number) => request<unknown>('/api/users/me/balance/credit', { method: 'POST', body: { amount } }),
  },
  products: {
    list: (params: Record<string, string | number | undefined> = {}) => {
      const q = new URLSearchParams();
      Object.entries(params).forEach(([k, v]) => { if (v !== undefined && v !== '') q.set(k, String(v)); });
      return request<{ products?: unknown[] }>('/api/products' + (q.toString() ? '?' + q.toString() : ''));
    },
    get: (id: string) => request<unknown>('/api/products/' + id),
    categories: () => request<unknown>('/api/products/categories'),
    create: (formData: Record<string, unknown> | FormData) =>
      request<unknown>('/api/products', { method: 'POST', body: formData as Record<string, unknown>, headers: {} }),
    update: (id: string, formData: Record<string, unknown>) =>
      request<unknown>('/api/products/' + id, { method: 'PATCH', body: formData, headers: {} }),
    delete: (id: string) => request<unknown>('/api/products/' + id, { method: 'DELETE' }),
    slots: (productId: string) => request<unknown>('/api/products/' + productId + '/slots'),
    addSlot: (productId: string, slot_at: string) =>
      request<unknown>('/api/products/' + productId + '/slots', { method: 'POST', body: { slot_at } }),
    bookSlot: (productId: string, slotId: string) =>
      request<unknown>('/api/products/' + productId + '/slots/' + slotId + '/book', { method: 'POST' }),
    closedContent: (productId: string) => request<unknown>('/api/products/' + productId + '/closed-content'),
  },
  orders: {
    my: () => request<{ asBuyer?: unknown[]; asSeller?: unknown[] }>('/api/orders/my'),
    get: (id: string) => request<unknown>('/api/orders/' + id),
    create: (product_id: number, data?: { urgent?: boolean }) =>
      request<unknown>('/api/orders', { method: 'POST', body: { product_id, ...data } }),
    update: (id: string, data: Record<string, unknown> | string) =>
      request<unknown>('/api/orders/' + id, { method: 'PATCH', body: typeof data === 'string' ? { status: data } : data }),
  },
  professionals: {
    search: (params: Record<string, string>) =>
      request<{ professionals?: unknown[] }>('/api/professionals/search?' + new URLSearchParams(params).toString()),
    get: (id: string) => request<unknown>('/api/professionals/' + id),
  },
  remittances: {
    my: () => request<unknown>('/api/remittances/my'),
    create: (to_identifier: string, amount: number, currency?: string) =>
      request<unknown>('/api/remittances', { method: 'POST', body: { to_identifier, amount, currency: currency || 'USD' } }),
  },
  wallet: {
    balances: () => request<unknown>('/api/wallet/balances'),
    transactions: (params?: Record<string, string>) =>
      request<unknown>('/api/wallet/transactions' + (params ? '?' + new URLSearchParams(params).toString() : '')),
    transfer: (to_user_id: number, currency: string, amount: number) =>
      request<unknown>('/api/wallet/transfer', { method: 'POST', body: { to_user_id, currency, amount } }),
    transferVerify: (to_user_id: number) =>
      request<unknown>('/api/wallet/transfer/verify', { method: 'POST', body: { to_user_id } }),
  },
  conversations: {
    list: () => request<unknown[]>('/api/conversations'),
    get: (id: string) => request<unknown>('/api/conversations/' + id),
    unreadCount: () => request<{ unread?: number }>('/api/conversations/unread-count'),
    create: (user_id: number, product_id?: number) =>
      request<unknown>('/api/conversations', { method: 'POST', body: { user_id, product_id } }),
  },
  messages: {
    list: (conversationId: string) => request<unknown>('/api/messages/conversation/' + conversationId),
    send: (conversationId: string, body: string) =>
      request<unknown>('/api/messages/conversation/' + conversationId, { method: 'POST', body: { body } }),
    read: (id: string) => request<unknown>('/api/messages/' + id + '/read', { method: 'POST' }),
  },
  notifications: {
    settings: () => request<unknown>('/api/notifications/settings'),
    updateSettings: (data: Record<string, unknown>) =>
      request<unknown>('/api/notifications/settings', { method: 'PATCH', body: data }),
    history: () => request<unknown>('/api/notifications/history'),
    get: (id: string) => request<unknown>('/api/notifications/history/' + id),
    markRead: (id: string) => request<unknown>('/api/notifications/history/' + id + '/read', { method: 'POST' }),
  },
  admin: {
    stats: () => request<unknown>('/api/admin/stats'),
    reportsList: (params?: Record<string, string>) =>
      request<unknown>('/api/admin/reports' + (params ? '?' + new URLSearchParams(params).toString() : '')),
    reportGet: (id: string) => request<unknown>('/api/admin/reports/' + id),
    reportAssign: (id: string, assigned_to: number) =>
      request<unknown>('/api/admin/reports/' + id + '/assign', { method: 'POST', body: { assigned_to } }),
    reportResolve: (id: string, resolution: string, status?: string) =>
      request<unknown>('/api/admin/reports/' + id + '/resolve', { method: 'POST', body: { resolution, status: status || 'resolved' } }),
    userGet: (id: string) => request<unknown>('/api/admin/users/' + id),
    userBan: (id: string, reason: string, expires_at?: number) =>
      request<unknown>('/api/admin/users/' + id + '/ban', { method: 'POST', body: expires_at != null ? { reason, expires_at } : { reason } }),
    userUnban: (id: string) => request<unknown>('/api/admin/users/' + id + '/unban', { method: 'POST' }),
  },
};
