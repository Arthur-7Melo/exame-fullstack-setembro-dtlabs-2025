import React, { useState, useEffect, type JSX } from 'react';
import { Link } from 'react-router-dom';
import { deviceService } from '../services/deviceService';
import type { Device } from '../types';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { AlertCircle, Cpu, HardDrive, Thermometer, Wifi, WifiOff } from 'lucide-react';

const Home: React.FC = () => {
  const [devices, setDevices] = useState<Device[]>([]);
  const [deviceMetrics, setDeviceMetrics] = useState<{[key: string]: any}>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDevices();
  }, []);

  const loadDevices = async () => {
    try {
      setLoading(true);
      const devicesData = await deviceService.listDevices();
      setDevices(devicesData);
      
      const metrics: {[key: string]: any} = {};
      for (const device of devicesData) {
        try {
          const heartbeat = await deviceService.getLatestDeviceHeartbeat(device.uuid);
          metrics[device.uuid] = {
            cpu: heartbeat.cpu,
            ram: heartbeat.ram,
            temperature: heartbeat.temperature,
            lastUpdate: heartbeat.created_at,
            status: heartbeat.connectivity === 1 ? 'online' : 'offline'
          };
        } catch (err) {
          console.warn(`No heartbeat data for device ${device.uuid}:`, err);
          metrics[device.uuid] = {
            cpu: 0,
            ram: 0,
            temperature: 0,
            lastUpdate: null,
            status: 'offline'
          };
        }
      }
      setDeviceMetrics(metrics);
    } catch (err) {
      setError('Error loading devices');
      console.error('Error loading devices:', err);
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return <Wifi className="h-4 w-4 text-green-500" />;
      case 'offline': return <WifiOff className="h-4 w-4 text-red-500" />;
      case 'warning': return <AlertCircle className="h-4 w-4 text-yellow-500" />;
      default: return <WifiOff className="h-4 w-4 text-gray-500" />;
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'online': return <Badge variant="default" className="bg-green-100 text-green-800">Online</Badge>;
      case 'offline': return <Badge variant="default" className="bg-red-100 text-red-800">Offline</Badge>;
      case 'warning': return <Badge variant="default" className="bg-yellow-100 text-yellow-800">Warning</Badge>;
      default: return <Badge variant="secondary">Unknown</Badge>;
    }
  };

  const getDeviceMetrics = (deviceUuid: string) => {
    return deviceMetrics[deviceUuid] || { 
      cpu: 0, 
      ram: 0, 
      temperature: 0, 
      lastUpdate: null, 
      status: 'offline' 
    };
  };

  if (loading) {
    return <LoadingSkeleton />;
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      <Header devicesCount={devices.length} />
      
      {error && <ErrorAlert error={error} />}

      {devices.length === 0 ? (
        <NoDevicesCard />
      ) : (
        <DevicesGrid 
          devices={devices}
          getDeviceMetrics={getDeviceMetrics}
          getStatusIcon={getStatusIcon}
          getStatusBadge={getStatusBadge}
        />
      )}
    </div>
  );
};

const LoadingSkeleton: React.FC = () => (
  <div className="container mx-auto p-6 space-y-6">
    <div className="flex justify-between items-center">
      <Skeleton className="h-8 w-48" />
      <Skeleton className="h-10 w-32" />
    </div>
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
      {[...Array(6)].map((_, i) => (
        <Card key={i}>
          <CardHeader>
            <Skeleton className="h-6 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
          </CardHeader>
          <CardContent className="space-y-4">
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
            <Skeleton className="h-4 w-full" />
          </CardContent>
        </Card>
      ))}
    </div>
  </div>
);

const Header: React.FC<{ devicesCount: number }> = ({ devicesCount }) => (
  <div className="flex justify-between items-center">
    <div>
      <h2 className="text-3xl font-bold tracking-tight transition-transform transform hover:scale-105 active:scale-95 hover:cursor-pointer">Meus Dispositivos</h2>
      <p className="text-muted-foreground">{devicesCount} dispositivo(s) encontrado(s)</p>
    </div>
    <Button asChild className='transition-transform transform hover:scale-105 hover:cursor-pointer'>
      <Link to="/devices">Ver Detalhes</Link>
    </Button>
  </div>
);

const ErrorAlert: React.FC<{ error: string }> = ({ error }) => (
  <Alert variant="destructive">
    <AlertCircle className="h-4 w-4" />
    <AlertDescription>{error}</AlertDescription>
  </Alert>
);

const NoDevicesCard: React.FC = () => (
  <Card>
    <CardContent className="flex flex-col items-center justify-center p-12 text-center">
      <WifiOff className="h-12 w-12 text-muted-foreground mb-4" />
      <h3 className="text-lg font-semibold">Nenhum dispositivo encontrado</h3>
      <p className="text-muted-foreground mt-2">
        Você ainda não possui dispositivos cadastrados em sua conta.
      </p>
    </CardContent>
  </Card>
);

const DevicesGrid: React.FC<{
  devices: Device[];
  getDeviceMetrics: (uuid: string) => any;
  getStatusIcon: (status: string) => JSX.Element;
  getStatusBadge: (status: string) => JSX.Element;
}> = ({ devices, getDeviceMetrics, getStatusIcon, getStatusBadge }) => (
  <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
    {devices.map((device) => {
      const metrics = getDeviceMetrics(device.uuid);
      
      return (
        <DeviceCard
          key={device.uuid}
          device={device}
          metrics={metrics}
          getStatusIcon={getStatusIcon}
          getStatusBadge={getStatusBadge}
        />
      );
    })}
  </div>
);

const DeviceCard: React.FC<{
  device: Device;
  metrics: any;
  getStatusIcon: (status: string) => JSX.Element;
  getStatusBadge: (status: string) => JSX.Element;
}> = ({ device, metrics, getStatusIcon, getStatusBadge }) => (
  <Card className="hover:shadow-2xl transition-shadow">
    <CardHeader className="pb-3">
      <div className="flex justify-between items-start">
        <div>
          <CardTitle className="text-lg flex items-center gap-2">
            {getStatusIcon(metrics.status)}
            {device.name}
          </CardTitle>
          <CardDescription>{device.location}</CardDescription>
        </div>
        {getStatusBadge(metrics.status)}
      </div>
    </CardHeader>
    <CardContent className="space-y-4">
      <MetricRow 
        icon={<Cpu className="h-4 w-4" />}
        label="CPU"
        value={metrics.cpu}
        unit="%"
      />
      <MetricRow 
        icon={<HardDrive className="h-4 w-4" />}
        label="RAM"
        value={metrics.ram}
        unit="%"
      />
      <MetricRow 
        icon={<Thermometer className="h-4 w-4" />}
        label="Temperatura"
        value={metrics.temperature}
        unit="°C"
      />
      {metrics.lastUpdate && (
        <p className="text-xs text-muted-foreground">
          Última atualização: {new Date(metrics.lastUpdate).toLocaleString()}
        </p>
      )}
    </CardContent>
  </Card>
);

const MetricRow: React.FC<{
  icon: JSX.Element;
  label: string;
  value: number;
  unit: string;
}> = ({ icon, label, value, unit }) => (
  <div className="space-y-2">
    <div className="flex items-center justify-between text-sm">
      <span className="flex items-center gap-2 text-muted-foreground">
        {icon}
        {label}
      </span>
      <span>{value.toFixed(1)}{unit}</span>
    </div>
    <Progress value={value} className="h-2" />
  </div>
);

export default Home;