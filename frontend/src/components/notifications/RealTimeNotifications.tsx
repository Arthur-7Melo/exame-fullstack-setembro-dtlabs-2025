import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { useWebSocket } from '@/hooks/useWebSocket';
import { Wifi, WifiOff } from 'lucide-react';
import { formatDiskSpace, formatLatency, formatNumber, formatPercentage, formatTemperature } from '@/utils/format';

interface RealTimeNotificationsProps {
  getParameterIcon: (parameter: string) => React.ReactNode;
}

export const RealTimeNotifications: React.FC<RealTimeNotificationsProps> = ({ getParameterIcon }) => {
  const { messages, isConnected } = useWebSocket();

  if (!isConnected) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <WifiOff className="h-5 w-5 text-red-500" />
            Conexão em Tempo Real
          </CardTitle>
          <CardDescription>
            Conectando ao serviço de notificações em tempo real...
          </CardDescription>
        </CardHeader>
      </Card>
    );
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Wifi className="h-5 w-5 text-green-500" />
            Notificações em Tempo Real
          </CardTitle>
          <CardDescription>
            Alertas disparados pelos seus dispositivos
          </CardDescription>
        </CardHeader>
      </Card>

      {messages.length === 0 ? (
        <Card>
          <CardContent className="p-6 text-center">
            <p className="text-muted-foreground">
              Nenhuma notificação em tempo real no momento
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          {messages.map((message: any, index: number) => (
            <Card key={index} className="border-l-4 border-l-red-500 h-full flex flex-col min-h-0">
              <CardHeader className="pb-3 flex-shrink-0">
                <div className="flex justify-between items-start gap-2">
                  <div className="min-w-0 flex-1">
                    <CardTitle className="text-base truncate">{message.name}</CardTitle>
                    {message.description && (
                      <CardDescription className="text-xs line-clamp-2 mt-1">
                        {message.description}
                      </CardDescription>
                    )}
                  </div>
                  <Badge variant="destructive" className="flex-shrink-0">
                    Alerta
                  </Badge>
                </div>
              </CardHeader>
              
              <CardContent className="flex-1 min-h-0 space-y-3">
               <div className="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <span className="font-medium text-xs">Dispositivo:</span>
                    <div className="text-xs text-muted-foreground truncate">{message.device_sn}</div>
                  </div>
                  <div>
                    <span className="font-medium text-xs">Valor:</span>
                    <div className="text-xs text-muted-foreground">{formatNumber(message.triggered_value)}</div>
                  </div>
                  <div className="col-span-2">
                    <span className="font-medium text-xs">Horário:</span>
                    <div className="text-xs text-muted-foreground">
                      {new Date(message.timestamp).toLocaleString('pt-BR')}
                    </div>
                  </div>
                </div>
                
                <div>
                  <span className="font-medium text-xs">Dados do Heartbeat:</span>
                    <div className="grid grid-cols-2 gap-1 mt-2">
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('cpu')}
                        CPU: {formatPercentage(message.heartbeat_data.cpu)}
                      </Badge>
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('ram')}
                        RAM: {formatPercentage(message.heartbeat_data.ram)}
                      </Badge>
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('disk_free')}
                        Disco: {formatDiskSpace(message.heartbeat_data.disk_free)}
                      </Badge>
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('temperature')}
                        Temp: {formatTemperature(message.heartbeat_data.temperature)}
                      </Badge>
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('latency')}
                        Lat: {formatLatency(message.heartbeat_data.latency)}
                      </Badge>
                      <Badge variant="secondary" className="text-xs flex items-center gap-1 justify-center py-1">
                        {getParameterIcon('connectivity')}
                        Conex: {formatPercentage(message.heartbeat_data.connectivity)}
                      </Badge>
                    </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
};