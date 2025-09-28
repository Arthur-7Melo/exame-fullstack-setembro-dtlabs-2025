import React, { useState, useEffect, useMemo } from 'react';
import { deviceService } from '../services/deviceService';
import type { Device, HeartbeatData } from '../types';
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { CalendarIcon, Filter, AlertCircle, Wifi, WifiOff, AlertTriangle } from 'lucide-react';
import { format, subDays, startOfDay, endOfDay } from 'date-fns';
import { ptBR } from 'date-fns/locale';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface DeviceFilters {
  status?: string;
  dateRange?: {
    start: Date;
    end: Date;
  };
}

const COLOR_PALETTE = [
  '#3B82F6', // Azul
  '#EF4444', // Vermelho
  '#10B981', // Verde
  '#F59E0B', // Amarelo
  '#8B5CF6', // Roxo
  '#EC4899', // Rosa
  '#06B6D4', // Ciano
  '#84CC16', // Verde Lima
  '#F97316', // Laranja
  '#6366F1', // Índigo
  '#14B8A6', // Turquesa
  '#F43F5E', // Rosa Escuro
];

const getDeviceColor = (deviceId: string): string => {
  let hash = 0;
  for (let i = 0; i < deviceId.length; i++) {
    hash = deviceId.charCodeAt(i) + ((hash << 5) - hash);
  }
  const colorIndex = Math.abs(hash) % COLOR_PALETTE.length;
  return COLOR_PALETTE[colorIndex];
};

const darkenColor = (color: string, percent: number): string => {
  const num = parseInt(color.replace("#", ""), 16);
  const amt = Math.round(2.55 * percent);
  const R = (num >> 16) - amt;
  const G = (num >> 8 & 0x00FF) - amt;
  const B = (num & 0x0000FF) - amt;
  return "#" + (
    0x1000000 +
    (R < 255 ? R < 1 ? 0 : R : 255) * 0x10000 +
    (G < 255 ? G < 1 ? 0 : G : 255) * 0x100 +
    (B < 255 ? B < 1 ? 0 : B : 255)
  ).toString(16).slice(1);
};

const Devices: React.FC = () => {
  const [devices, setDevices] = useState<Device[]>([]);
  const [selectedDevices, setSelectedDevices] = useState<string[]>([]);
  const [heartbeatData, setHeartbeatData] = useState<HeartbeatData[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingChart, setLoadingChart] = useState(false);
  const [error, setError] = useState('');
  const [filters, setFilters] = useState<DeviceFilters>({
    status: 'all',
    dateRange: {
      start: subDays(new Date(), 7),
      end: new Date()
    }
  });

  const deviceColors = useMemo(() => {
    const colors: { [key: string]: string } = {};
    devices.forEach((device) => {
      colors[device.uuid] = getDeviceColor(device.uuid);
    });
    return colors;
  }, [devices]);

  useEffect(() => {
    loadDevices();
  }, []);

  const loadDevices = async () => {
    try {
      setLoading(true);
      const devicesData = await deviceService.listDevices();
      setDevices(devicesData);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Erro ao carregar dispositivos');
      console.error('Error loading devices:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDeviceToggle = (deviceId: string) => {
    setSelectedDevices(prev => {
      if (prev.includes(deviceId)) {
        return prev.filter(id => id !== deviceId);
      } else {
        return [...prev, deviceId];
      }
    });
  };

  const handleSelectAll = () => {
    if (selectedDevices.length === devices.length) {
      setSelectedDevices([]);
    } else {
      setSelectedDevices(devices.map(d => d.uuid));
    }
  };

  const handleLoadData = async () => {
    if (selectedDevices.length === 0) {
      setError('Selecione pelo menos um dispositivo');
      return;
    }

    if (!filters.dateRange) {
      setError('Selecione um período válido');
      return;
    }

    try {
      setLoadingChart(true);
      setError('');

      const allHeartbeats: HeartbeatData[] = [];
      
      for (const deviceUuid of selectedDevices) {
        try {
          const heartbeats = await deviceService.getDeviceHeartbeats(
            deviceUuid,
            startOfDay(filters.dateRange.start),
            endOfDay(filters.dateRange.end)
          );
          allHeartbeats.push(...heartbeats);
        } catch (err: any) {
          console.error(`Error loading heartbeats for device ${deviceUuid}:`, err);
        }
      }

      setHeartbeatData(allHeartbeats);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Erro ao carregar dados dos dispositivos');
      console.error('Error loading heartbeat data:', err);
    } finally {
      setLoadingChart(false);
    }
  };

  const filteredDevices = devices.filter(device => {
    if (filters.status && filters.status !== 'all' && device.status !== filters.status) {
      return false;
    }
    return true;
  });

  const chartData = React.useMemo(() => {
    const groupedData: { [key: string]: any } = {};
    
    heartbeatData.forEach(item => {
      const dateKey = format(new Date(item.created_at), 'dd/MM HH:mm');
      
      if (!groupedData[dateKey]) {
        groupedData[dateKey] = { timestamp: dateKey };
      }
      
      groupedData[dateKey][`cpu_${item.device_id}`] = item.cpu;
      groupedData[dateKey][`ram_${item.device_id}`] = item.ram;
      groupedData[dateKey][`temp_${item.device_id}`] = item.temperature;
    });
    
    return Object.values(groupedData);
  }, [heartbeatData]);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'online': return <Wifi className="h-4 w-4 text-green-500" />;
      case 'offline': return <WifiOff className="h-4 w-4 text-red-500" />;
      case 'warning': return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      default: return <WifiOff className="h-4 w-4 text-gray-500" />;
    }
  };

  const getDeviceName = (deviceId: string) => {
    const device = devices.find(d => d.uuid === deviceId);
    return device ? device.name : `Device ${deviceId.substring(0, 8)}...`;
  };

  const CustomLegend: React.FC = () => {
    const visibleDevices = devices.filter(device => selectedDevices.includes(device.uuid));
    
    if (visibleDevices.length === 0) return null;

    return (
      <div className="flex flex-wrap gap-3 mb-6 p-4 bg-muted/30 rounded-lg border">
        <span className="text-sm font-medium text-muted-foreground w-full mb-2">
          Dispositivos no gráfico:
        </span>
        {visibleDevices.map(device => (
          <div 
            key={device.uuid}
            className="flex items-center gap-2 px-3 py-2 rounded-full border cursor-pointer hover:shadow-sm transition-all"
            style={{ 
              borderColor: deviceColors[device.uuid],
              backgroundColor: `${deviceColors[device.uuid]}15`
            }}
            onClick={() => handleDeviceToggle(device.uuid)}
            title="Clique para remover"
          >
            <div 
              className="w-3 h-3 rounded-full"
              style={{ backgroundColor: deviceColors[device.uuid] }}
            />
            <span className="text-sm font-medium">{device.name}</span>
            <span className="text-xs text-muted-foreground hover:text-foreground ml-1">
              ×
            </span>
          </div>
        ))}
      </div>
    );
  };

  if (loading) {
    return (
      <div className="container mx-auto p-6 space-y-6">
        <Skeleton className="h-8 w-48" />
        <div className="grid gap-6 md:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
        <Skeleton className="h-96" />
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4">
        <div>
         <h2 className="text-2xl sm:text-3xl font-bold tracking-tight transition-transform transform hover:scale-105 active:scale-95 hover:cursor-pointer">Análise de Dispositivos</h2>
        </div>
        <Button
         onClick={handleLoadData}
         disabled={loadingChart || selectedDevices.length === 0}
         className="w-full sm:w-auto transition-transform transform hover:scale-105 hover:cursor-pointer"
         >
          {loadingChart ? 'Carregando...' : 'Carregar Dados'}
        </Button>
      </div>

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      <Card className='hover:shadow-2xl transition-shadow'>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filtros
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
        <div className="flex flex-col gap-4">
          <div className="space-y-2">
            <Label>Status</Label>
            <Select
              value={filters.status}
              onValueChange={(value: string) => 
                setFilters(prev => ({ ...prev, status: value }))
              }
            >
              <SelectTrigger className="w-full sm:w-[180px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos</SelectItem>
                <SelectItem value="online">Online</SelectItem>
                <SelectItem value="offline">Offline</SelectItem>
                <SelectItem value="warning">Atenção</SelectItem>
              </SelectContent>
            </Select>
          </div>

        <div className="space-y-2">
          <Label>Período</Label>
          <div className="flex flex-col sm:flex-row gap-2">
            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" className="w-full sm:w-[140px] justify-start text-left font-normal">
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {filters.dateRange ? format(filters.dateRange.start, 'dd/MM/yyyy') : 'Início'}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="single"
                  selected={filters.dateRange?.start}
                  onSelect={(date) => date && setFilters(prev => ({
                    ...prev,
                    dateRange: { ...prev.dateRange!, start: date }
                  }))}
                  locale={ptBR}
                />
              </PopoverContent>
            </Popover>

            <Popover>
              <PopoverTrigger asChild>
                <Button variant="outline" className="w-full sm:w-[140px] justify-start text-left font-normal">
                  <CalendarIcon className="mr-2 h-4 w-4" />
                  {filters.dateRange ? format(filters.dateRange.end, 'dd/MM/yyyy') : 'Fim'}
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-auto p-0" align="start">
                <Calendar
                  mode="single"
                  selected={filters.dateRange?.end}
                  onSelect={(date) => date && setFilters(prev => ({
                    ...prev,
                    dateRange: { ...prev.dateRange!, end: date }
                  }))}
                  locale={ptBR}
                />
              </PopoverContent>
            </Popover>
            </div>
          </div>
        </div>
        </CardContent>
        </Card>

      <Card className='hover:shadow-2xl transition-shadow'>
        <CardHeader>
          <CardTitle>Dispositivos ({filteredDevices.length})</CardTitle>
         <CardDescription className="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4">
          <span>Selecione os dispositivos para análise</span>
          <div className="flex items-center gap-2 sm:gap-4">
            <Button variant="outline" size="sm" className='hover:cursor-pointer' onClick={handleSelectAll}>
              {selectedDevices.length === devices.length ? 'Desmarcar Todos' : 'Selecionar Todos'}
            </Button>
            <span className="text-sm">
              {selectedDevices.length} dispositivo(s) selecionado(s)
            </span>
          </div>
        </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {filteredDevices.map((device) => (
              <div 
                key={device.uuid} 
                className="flex items-center space-x-3 p-3 border rounded-lg hover:shadow-md transition-shadow"
                style={{
                  borderLeftColor: deviceColors[device.uuid],
                  borderLeftWidth: '4px'
                }}
              >
                <Checkbox
                  checked={selectedDevices.includes(device.uuid)}
                  onCheckedChange={() => handleDeviceToggle(device.uuid)}
                  className='cursor-pointer'
                />
                <div 
                  className="w-3 h-3 rounded-full flex-shrink-0"
                  style={{ backgroundColor: deviceColors[device.uuid] }}
                />
                {getStatusIcon(device.status)}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium truncate">{device.name}</p>
                  <p className="text-xs text-muted-foreground truncate">{device.location}</p>
                  <p className="text-xs text-muted-foreground">ID: {device.uuid}</p>
                  <p className="text-xs text-muted-foreground">SN: {device.sn}</p>
                </div>
                <Badge variant={device.status === 'online' ? 'default' : 'secondary'}>
                  {device.status === 'online' ? 'Online' : 
                   device.status === 'offline' ? 'Offline' : 'Atenção'}
                </Badge>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className='hover:shadow-2xl transition-shadow'>
        <CardHeader>
          <CardTitle>Dados de Monitoramento</CardTitle>
          <CardDescription>
            Dados históricos de CPU, RAM e Temperatura para os dispositivos selecionados
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loadingChart ? (
            <Skeleton className="h-96" />
          ) : selectedDevices.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              Selecione pelo menos um dispositivo e clique em "Carregar Dados"
            </div>
          ) : chartData.length === 0 ? (
            <div className="text-center py-12 text-muted-foreground">
              Nenhum dado encontrado para o período selecionado
            </div>
          ) : (
            <div className="space-y-8">
              <CustomLegend />

              <div>
                <h3 className="text-lg font-semibold mb-4">Uso de CPU (%)</h3>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                    <XAxis 
                      dataKey="timestamp" 
                      stroke="#666"
                      fontSize={12}
                    />
                    <YAxis 
                      domain={[0, 100]} 
                      stroke="#666"
                      fontSize={12}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'white', 
                        border: '1px solid #e5e5e5',
                        borderRadius: '6px',
                        boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
                      }}
                    />
                    <Legend />
                    {selectedDevices.map(deviceId => {
                      const color = deviceColors[deviceId];
                      
                      return (
                        <Line
                          key={`cpu-${deviceId}`}
                          type="monotone"
                          dataKey={`cpu_${deviceId}`}
                          stroke={color}
                          strokeWidth={3}
                          dot={{ fill: color, strokeWidth: 2, r: 4 }}
                          activeDot={{ 
                            r: 6, 
                            fill: darkenColor(color, 20),
                            stroke: darkenColor(color, 30),
                            strokeWidth: 2
                          }}
                          name={`CPU - ${getDeviceName(deviceId)}`}
                        />
                      );
                    })}
                  </LineChart>
                </ResponsiveContainer>
              </div>

              <div>
                <h3 className="text-lg font-semibold mb-4">Uso de RAM (%)</h3>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                    <XAxis 
                      dataKey="timestamp" 
                      stroke="#666"
                      fontSize={12}
                    />
                    <YAxis 
                      domain={[0, 100]} 
                      stroke="#666"
                      fontSize={12}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'white', 
                        border: '1px solid #e5e5e5',
                        borderRadius: '6px',
                        boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
                      }}
                    />
                    <Legend />
                    {selectedDevices.map(deviceId => {
                      const color = deviceColors[deviceId];
                      
                      return (
                        <Line
                          key={`ram-${deviceId}`}
                          type="monotone"
                          dataKey={`ram_${deviceId}`}
                          stroke={color}
                          strokeWidth={3}
                          dot={{ fill: color, strokeWidth: 2, r: 4 }}
                          activeDot={{ 
                            r: 6, 
                            fill: darkenColor(color, 20),
                            stroke: darkenColor(color, 30),
                            strokeWidth: 2
                          }}
                          name={`RAM - ${getDeviceName(deviceId)}`}
                        />
                      );
                    })}
                  </LineChart>
                </ResponsiveContainer>
              </div>

              <div>
                <h3 className="text-lg font-semibold mb-4">Temperatura (°C)</h3>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                    <XAxis 
                      dataKey="timestamp" 
                      stroke="#666"
                      fontSize={12}
                    />
                    <YAxis 
                      stroke="#666"
                      fontSize={12}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'white', 
                        border: '1px solid #e5e5e5',
                        borderRadius: '6px',
                        boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)'
                      }}
                    />
                    <Legend />
                    {selectedDevices.map(deviceId => {
                      const color = deviceColors[deviceId];
                      
                      return (
                        <Line
                          key={`temp-${deviceId}`}
                          type="monotone"
                          dataKey={`temp_${deviceId}`}
                          stroke={color}
                          strokeWidth={3}
                          dot={{ fill: color, strokeWidth: 2, r: 4 }}
                          activeDot={{ 
                            r: 6, 
                            fill: darkenColor(color, 20),
                            stroke: darkenColor(color, 30),
                            strokeWidth: 2
                          }}
                          name={`Temp - ${getDeviceName(deviceId)}`}
                        />
                      );
                    })}
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default Devices;