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
  const url = (base || '') + path;
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
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
  auth: {
    login: (email: string, password: string) =>
      request<{ token: string; user: { id: number; email?: string; name?: string; role?: string } }>(
        '/api/auth/login',
        { method: 'POST', body: { email, password } }
      ),
    logout: () => {
      setToken(null);
      setUser(null);
    },
  },
  users: {
    me: () =>
      request<{ id: number; email?: string; name?: string; role?: string }>('/api/users/me'),
  },
  products: {
    list: (params: Record<string, string | number | undefined> = {}) => {
      const q = new URLSearchParams();
      Object.entries(params).forEach(([k, v]) => {
        if (v !== undefined && v !== '') q.set(k, String(v));
      });
      const query = q.toString();
      return request<{ products?: Array<{ id: number; name?: string; price?: number; category?: string }> }>(
        '/api/products' + (query ? '?' + query : '')
      );
    },
  },
};
