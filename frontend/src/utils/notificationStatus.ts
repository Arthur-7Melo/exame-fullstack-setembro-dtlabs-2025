import type { NotificationResponse, Device } from '@/types/index';

export const hasRemovedDevices = (notification: NotificationResponse, allDevices: Device[]): boolean => {
  if (notification.device_ids.length === 0) return false;

  return notification.device_ids.some(deviceId =>
    !allDevices.some(device => device.uuid === deviceId)
  );
};

export const isNotificationInactiveDueToRemoval = (notification: NotificationResponse, allDevices: Device[]): boolean => {
  if (notification.device_ids.length === 0) return false;

  return notification.device_ids.every(deviceId =>
    !allDevices.some(device => device.uuid === deviceId)
  );
};

export const getDeviceNamesWithStatus = (notification: NotificationResponse, allDevices: Device[]): string[] => {
  if (notification.device_ids.length === 0) return ['Todos os dispositivos'];

  return notification.device_ids.map(deviceId => {
    const device = allDevices.find(d => d.uuid === deviceId);
    if (device) {
      return `${device.name} (${device.sn})`;
    } else {
      return `❌ Dispositivo removido`;
    }
  });
};