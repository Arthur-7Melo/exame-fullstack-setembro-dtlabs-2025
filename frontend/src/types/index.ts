export interface User {
  id: string;
  email: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
}

export interface JwtPayload {
  user_id: string;
  exp: number;
  iat?: number;
  iss?: string;
}

export interface Device {
  uuid: string;
  name: string;
  location: string;
  sn: string;
  description: string;
  user_id: string;
  created_at: string;
  updated_at: string;
  status: 'online' | 'offline' | 'warning';
  lastHeartbeat?: string;
  cpuUsage: number;
  ramUsage: number;
  temperature: number;
}

export interface HeartbeatData {
  id: string;
  device_id: string;
  cpu: number;
  ram: number;
  disk_free: number;
  temperature: number;
  latency: number;
  connectivity: number;
  boot_time: string;
  created_at: string;
}

export interface DeviceFilters {
  status?: string;
  dateRange?: {
    start: Date;
    end: Date;
  };
}

export interface NotificationCondition {
  parameter: 'cpu' | 'ram' | 'disk_free' | 'temperature' | 'latency' | 'connectivity';
  operator: '>' | '<' | '>=' | '<=' | '==' | '!=';
  value: number;
}

export interface CreateNotificationRequest {
  name: string;
  description?: string;
  enabled: boolean;
  conditions: NotificationCondition[];
  device_ids: string[];
}

export interface NotificationResponse {
  id: string;
  user_id: string;
  name: string;
  description?: string;
  enabled: boolean;
  conditions: NotificationCondition[];
  device_ids: string[];
  created_at: string;
  updated_at: string;
}
