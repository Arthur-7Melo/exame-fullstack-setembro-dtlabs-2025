import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Bell, Plus, Cpu, HardDrive, Thermometer, Clock, Network } from 'lucide-react';
import { NotificationForm } from '@/components/notifications/NotificationForm';
import { NotificationList } from '@/components/notifications/NotificationList';
import { RealTimeNotifications } from '@/components/notifications/RealTimeNotifications';
import { notificationService } from '@/services/notificationService';
import type { NotificationResponse, Device } from '@/types/index';
import { toast } from 'sonner';

const Memory = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
    <rect x="4" y="4" width="16" height="16" rx="2" />
    <path d="M10 10h4" />
    <path d="M10 14h4" />
    <path d="M10 18h4" />
  </svg>
);

const Notifications = () => {
  const [notifications, setNotifications] = useState<NotificationResponse[]>([]);
  const [devices, setDevices] = useState<Device[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(true);

  const loadNotifications = async () => {
    try {
      const [notificationsData, devicesData] = await Promise.all([
        notificationService.getNotifications(),
        notificationService.getDevices()
      ]);
      setNotifications(notificationsData);
      setDevices(devicesData);
    } catch (error) {
      toast.error('Falha ao carregar notificações');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadNotifications();
  }, []);

  const handleCreateNotification = async (data: any) => {
    try {
      await notificationService.createNotification(data);
      toast.success('Notificação criada com sucesso');
      setShowForm(false);
      loadNotifications();
    } catch (error) {
      toast.error('Falha ao criar notificação');
    }
  };

  const getParameterIcon = (parameter: string) => {
    switch (parameter) {
      case 'cpu': return <Cpu className="h-4 w-4" />;
      case 'ram': return <Memory />;
      case 'disk_free': return <HardDrive className="h-4 w-4" />;
      case 'temperature': return <Thermometer className="h-4 w-4" />;
      case 'latency': return <Clock className="h-4 w-4" />;
      case 'connectivity': return <Network className="h-4 w-4" />;
      default: return <Bell className="h-4 w-4" />;
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto py-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">Carregando...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto py-6 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight transition-transform transform hover:scale-105 active:scale-95 hover:cursor-pointer">Notificações</h2>
          <p className="text-muted-foreground">
            Configure alertas baseados no status dos seus dispositivos
          </p>
        </div>
        <Button className='hover:cursor-pointer' onClick={() => setShowForm(!showForm)}>
          <Plus className="h-4 w-4 mr-2" />
          Nova Notificação
        </Button>
      </div>

      {showForm && (
        <Card>
          <CardHeader>
            <CardTitle>Criar Nova Notificação</CardTitle>
            <CardDescription>
              Configure alertas que serão disparados quando os parâmetros dos dispositivos atingirem certos valores
            </CardDescription>
          </CardHeader>
          <CardContent>
            <NotificationForm 
              devices={devices} 
              onSubmit={handleCreateNotification}
              onCancel={() => setShowForm(false)}
            />
          </CardContent>
        </Card>
      )}

      <Tabs defaultValue="configurations" className="space-y-6">
        <TabsList>
          <TabsTrigger value="configurations" className='cursor-pointer'>Minhas Configurações</TabsTrigger>
          <TabsTrigger value="realtime" className='cursor-pointer'>Notificações em Tempo Real</TabsTrigger>
        </TabsList>

        <TabsContent value="configurations" className="space-y-4">
          <NotificationList 
            notifications={notifications} 
            devices={devices}
            onRefresh={loadNotifications}
            getParameterIcon={getParameterIcon}
          />
        </TabsContent>

        <TabsContent value="realtime">
          <RealTimeNotifications getParameterIcon={getParameterIcon} />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Notifications;