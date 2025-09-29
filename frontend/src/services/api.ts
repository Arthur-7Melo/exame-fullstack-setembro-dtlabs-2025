import axios from 'axios';
import { toast } from 'sonner';

const API_BASE_URL = import.meta.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const status = error.response.status;
      if (status === 401 || status === 403) {
        toast.error('Sessão inválida ou expirada. Faça login novamente.');
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
        window.location.href = '/login';
      } else if (status >= 400 && status < 500) {
        const msg = error.response.data?.message || 'Erro na requisição';
        toast.error(msg);
      }
    } else {
      toast.error('Erro de rede. Verifique sua conexão.');
    }
    return Promise.reject(error);
  }
);

export default api;
