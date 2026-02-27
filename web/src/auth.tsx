import React, { createContext, useCallback, useContext, useEffect, useState } from 'react';
import { getToken, setToken, getUser, setUser, api } from './api';

export interface User {
  id: number;
  email?: string;
  name?: string;
  role?: string;
}

interface AuthState {
  token: string | null;
  user: User | null;
  loading: boolean;
  login: (email: string, password: string, remember?: boolean) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = useState<string | null>(() => getToken());
  const [user, setUserState] = useState<User | null>(() => getUser());
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    if (!getToken()) {
      setUserState(null);
      return;
    }
    try {
      const u = await api.users.me();
      setUserState(u);
      setUser(u);
    } catch {
      setToken(null);
      setUser(null);
      setUserState(null);
      setTokenState(null);
    }
  }, []);

  useEffect(() => {
    if (token) {
      refreshUser().finally(() => setLoading(false));
    } else {
      setUserState(null);
      setLoading(false);
    }
  }, [token, refreshUser]);

  const login = useCallback(
    async (email: string, password: string, remember = true) => {
      const res = await api.auth.login(email, password);
      setToken(res.token, remember);
      setUser(res.user);
      setTokenState(res.token);
      setUserState(res.user);
    },
    []
  );

  const logout = useCallback(() => {
    api.auth.logout();
    setTokenState(null);
    setUserState(null);
  }, []);

  const value: AuthState = {
    token,
    user,
    loading,
    login,
    logout,
    refreshUser,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}
