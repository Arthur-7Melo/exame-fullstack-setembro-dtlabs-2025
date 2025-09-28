import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { RefreshCw, Bell, BellOff } from 'lucide-react';
import type { NotificationResponse, Device } from '@/types/index';
import { formatNumber } from '@/utils/format';

interface NotificationListProps {
  notifications: NotificationResponse[];
  devices: Device[];
  onRefresh: () => void;
  getParameterIcon: (parameter: string) => React.ReactNode;
}

export const NotificationList: React.FC<NotificationListProps> = ({
  notifications,
  devices,
  onRefresh,
  getParameterIcon
}) => {
  if (notifications.length === 0) {
    return (
      <Card>
        <CardContent className="p-6 text-center">
          <Bell className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
          <h3 className="text-lg font-semibold">Nenhuma notificação configurada</h3>
          <p className="text-muted-foreground">
            Crie sua primeira notificação para receber alertas sobre seus dispositivos
          </p>
        </CardContent>
      </Card>
    );
  }

  const getDeviceNames = (deviceIds: string[]) => {
    if (deviceIds.length === 0) return ['Todos os dispositivos'];
    return deviceIds.map(uuid => {
      const device = devices.find(d => d.uuid === uuid);
      return device ? `${device.name} (${device.sn})` : 'Dispositivo não encontrado';
    });
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <div>
          <h3 className="text-lg font-semibold">
            {notifications.length} notificaç{notifications.length === 1 ? 'ão' : 'ões'} configurada{notifications.length === 1 ? '' : 's'}
          </h3>
        </div>
        <Button variant="outline" size="sm" onClick={onRefresh} className='hover:cursor-pointer'>
          <RefreshCw className="h-4 w-4 mr-1" />
          Atualizar
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {notifications.map(notification => (
          <Card key={notification.id} className='hover:shadow-2xl transition-shadow h-full flex flex-col min-h-0'>
            <CardHeader className="pb-3 flex-shrink-0">
              <div className="flex justify-between items-start gap-2">
                <div className="min-w-0 flex-1">
                  <CardTitle className="flex items-center gap-2 text-sm truncate">
                    {notification.enabled ? (
                      <Bell className="h-4 w-4 text-green-500 flex-shrink-0" />
                    ) : (
                      <BellOff className="h-4 w-4 text-gray-400 flex-shrink-0" />
                    )}
                    <span className="truncate">{notification.name}</span>
                  </CardTitle>
                  {notification.description && (
                    <CardDescription className="text-xs line-clamp-2 mt-1">
                      {notification.description}
                    </CardDescription>
                  )}
                </div>
                <Badge 
                  variant={notification.enabled ? "default" : "secondary"}
                  className="flex-shrink-0 text-xs"
                >
                  {notification.enabled ? 'Ativa' : 'Inativa'}
                </Badge>
              </div>
            </CardHeader>
            
            <CardContent className="flex-1 min-h-0 space-y-2">
              <div className="text-xs">
                <span className="font-medium">Dispositivos: </span>
                <span className="text-muted-foreground">
                  {getDeviceNames(notification.device_ids).join(', ')}
                </span>
              </div>
              
              <div className="text-xs">
                <span className="font-medium">Condições: </span>
                <div className="flex flex-wrap gap-1 mt-1">
                  {notification.conditions.map((condition: any, index: number) => (
                    <Badge 
                      key={index} 
                      variant="outline" 
                      className="text-xs flex items-center gap-1 py-0 px-1.5"
                    >
                      {getParameterIcon(condition.parameter)}
                      <span className="truncate text-[10px]">
                        {condition.parameter} {condition.operator} {formatNumber(condition.value)}
                      </span>
                    </Badge>
                  ))}
                </div>
              </div>
              
              <div className="text-[10px] text-muted-foreground pt-1 border-t">
                Criada em: {new Date(notification.created_at).toLocaleDateString('pt-BR')}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
};