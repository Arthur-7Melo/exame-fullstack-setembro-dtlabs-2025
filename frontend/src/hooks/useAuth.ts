import { useState, useCallback } from 'react';
import { jwtDecode } from 'jwt-decode';
import { authService } from '../services/authService';
import type { User, JwtPayload } from '../types';

const getUserFromToken = (token: string, email?: string): User | null => {
  try {
    const payload = jwtDecode<JwtPayload>(token);

    if (!payload.user_id) {
      return null;
    }

    return {
      id: payload.user_id,
      email: email || '',
    };
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
      return JSON.parse(saved);
    }

    if (token) {
      return getUserFromToken(token);
    }

    return null;
  });

  const login = useCallback(async (email: string, password: string) => {
    const response = await authService.login({ email, password });
    localStorage.setItem('authToken', response.token);

    const userData = getUserFromToken(response.token, email);
    if (userData) {
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } else {
      throw new Error('Failed to decode user information from token');
    }
  }, []);

  const signup = useCallback(async (email: string, password: string) => {
    const response = await authService.signup({ email, password });
    localStorage.setItem('authToken', response.token);

    const userData = getUserFromToken(response.token, email);
    if (userData) {
      setUser(userData);
      localStorage.setItem('user', JSON.stringify(userData));
    } else {
      throw new Error('Failed to decode user information from token');
    }
  }, []);

  const logout = useCallback(() => {
    authService.logout();
    setUser(null);
  }, []);

  const isTokenValid = useCallback((): boolean => {
    const token = localStorage.getItem('authToken');
    if (!token) return false;

    try {
      const payload = jwtDecode<JwtPayload>(token);
      const currentTime = Date.now() / 1000;
      return payload.exp > currentTime;
    } catch {
      return false;
    }
  }, []);

  return {
    user,
    login,
    signup,
    logout,
    isTokenValid
  };
};