import api from './api';
import type { Device, HeartbeatData } from '../types';

export const deviceService = {
  listDevices: async (): Promise<Device[]> => {
    const response = await api.get('/v1/devices');
    return response.data;
  },

  getDeviceHeartbeats: async (
    deviceUuid: string,
    startDate: Date,
    endDate: Date
  ): Promise<HeartbeatData[]> => {
    const response = await api.get(`/v1/devices/${deviceUuid}/heartbeats`, {
      params: {
        start: startDate.toISOString(),
        end: endDate.toISOString()
      }
    });
    return response.data;
  },

  getLatestDeviceHeartbeat: async (deviceUuid: string): Promise<HeartbeatData> => {
    const response = await api.get(`/v1/devices/${deviceUuid}/heartbeats/latest`);
    return response.data;
  }
};