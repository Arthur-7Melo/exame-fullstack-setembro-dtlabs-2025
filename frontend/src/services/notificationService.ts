import api from './api';
import type { CreateNotificationRequest, NotificationResponse } from '@/types/index';

export const notificationService = {
  createNotification: async (data: CreateNotificationRequest): Promise<NotificationResponse> => {
    const response = await api.post<NotificationResponse>('/v1/notifications', data);
    return response.data;
  },

  getNotifications: async (): Promise<NotificationResponse[]> => {
    const response = await api.get<NotificationResponse[]>('/v1/notifications');
    return response.data;
  },

  getDevices: async (): Promise<any[]> => {
    const response = await api.get('/v1/devices');
    return response.data;
  }
};