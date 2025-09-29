import { useState, useCallback } from 'react';
import { jwtDecode } from 'jwt-decode';
import { authService } from '../services/authService';
import type { User } from '../types';

interface JwtPayloadLocal {
  user_id?: string;
  exp?: number;
  [k: string]: any;
}

const getUserFromToken = (token: string): User | null => {
  try {
    const payload = jwtDecode<JwtPayloadLocal>(token);
    if (!payload || !payload.user_id) return null;

    const now = Date.now() / 1000;
    if (payload.exp && payload.exp <= now) return null; // token expirado

    return {
      id: payload.user_id,
    } as User;
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
};

export const useAuth = () => {
  const [user, setUser] = useState<User | null>(() => {
    const saved = localStorage.getItem('user');
    const token = localStorage.getItem('authToken');

    if (saved) {
      try {
        return JSON.parse(saved) as User;
      } catch {
        localStorage.removeItem('user');
      }
    }

    if (token) {
      const userData = getUserFromToken(token);
      if (userData) return userData;

      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
    }

    return null;
  });

  const login = useCallback(async (email: string, password: string): Promise<void> => {
    const response = await authService.login({ email, password });
    localStorage.setItem('authToken', response.token);

    const userData = getUserFromToken(response.token);
    if (userData) {
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } else {
      localStorage.removeItem('authToken');
      throw new Error('Failed to decode user information from token');
    }
  }, []);

  const signup = useCallback(async (email: string, password: string): Promise<void> => {
    const response = await authService.signup({ email, password });
    localStorage.setItem('authToken', response.token);

    const userData = getUserFromToken(response.token);
    if (userData) {
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } else {
      localStorage.removeItem('authToken');
      throw new Error('Failed to decode user information from token');
    }
  }, []);

  const logout = useCallback((): void => {
    authService.logout();
    setUser(null);
  }, []);

  const isTokenValid = useCallback((): boolean => {
    const token = localStorage.getItem('authToken');
    if (!token) return false;

    try {
      const payload = jwtDecode<JwtPayloadLocal>(token);
      const now = Date.now() / 1000;
      return !!(payload && payload.exp && payload.exp > now);
    } catch {
      return false;
    }
  }, []);

  return {
    user,
    login,
    signup,
    logout,
    isTokenValid,
  };
};

export default useAuth;
