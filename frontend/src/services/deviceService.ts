import api from './api';
import type { CreateDeviceData, Device, HeartbeatData, UpdateDeviceData } from '../types';

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
  },

  createDevice: async (deviceData: CreateDeviceData): Promise<Device> => {
    const response = await api.post('/v1/devices', deviceData);
    return response.data;
  },

  updateDevice: async (deviceUuid: string, deviceData: UpdateDeviceData): Promise<Device> => {
    const response = await api.put(`/v1/devices/${deviceUuid}`, deviceData);
    return response.data;
  },

  deleteDevice: async (deviceUuid: string): Promise<void> => {
    await api.delete(`/v1/devices/${deviceUuid}`);
  },

  getDevice: async (deviceUuid: string): Promise<Device> => {
    const response = await api.get(`/v1/devices/${deviceUuid}`);
    return response.data;
  }
};