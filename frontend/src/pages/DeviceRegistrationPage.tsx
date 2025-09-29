import React, { useState, useEffect } from 'react';
import { deviceService } from '@/services/deviceService';
import type { Device, CreateDeviceData, UpdateDeviceData } from '@/types/index';
import { useAuthContext } from '@/contexts/authContext';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Loader2, Pencil, Trash2, Search } from 'lucide-react';
import { toast } from 'sonner';
import { getErrorMessage } from '@/utils/errorHandler';

const DeviceRegistrationPage: React.FC = () => {
  const { user } = useAuthContext();
  
  const [devices, setDevices] = useState<Device[]>([]);
  const [filteredDevices, setFilteredDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [currentDevice, setCurrentDevice] = useState<Device | null>(null);
  const [errors, setErrors] = useState<{sn?: string}>({});
  const [formError, setFormError] = useState('');

  const [formData, setFormData] = useState<CreateDeviceData>({
    name: '',
    location: '',
    sn: '',
    description: ''
  });

  const validateSerialNumber = (sn: string): boolean => {
    return /^\d{12}$/.test(sn);
  };

  const loadDevices = async () => {
    try {
      setLoading(true);
      const deviceList = await deviceService.listDevices();
      setDevices(deviceList);
      setFilteredDevices(deviceList);
    } catch (error) {
      const errorMessage = getErrorMessage(error);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (user) {
      loadDevices();
    }
  }, [user]);

  useEffect(() => {
    if (!searchTerm) {
      setFilteredDevices(devices);
    } else {
      const filtered = devices.filter(device =>
        device.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        device.location.toLowerCase().includes(searchTerm.toLowerCase()) ||
        device.sn.toLowerCase().includes(searchTerm.toLowerCase())
      );
      setFilteredDevices(filtered);
    }
  }, [searchTerm, devices]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));

    if (name === 'sn' && errors.sn) {
      setErrors(prev => ({ ...prev, sn: undefined }));
    }
    if (formError) {
      setFormError('');
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      location: '',
      sn: '',
      description: ''
    });
    setIsEditing(false);
    setCurrentDevice(null);
    setErrors({});
    setFormError(''); 
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    setErrors({});
    setFormError('');

    if (!formData.name || !formData.location || !formData.sn) {
      setFormError('Nome, localização e número de série são obrigatórios');
      return;
    }

    if (!validateSerialNumber(formData.sn)) {
      setErrors({ sn: 'Número de série deve conter exatamente 12 números' });
      setFormError('Número de série deve conter exatamente 12 números');
      return;
    }

    try {
      setSubmitting(true);
      
      if (isEditing && currentDevice) {
        const updateData: UpdateDeviceData = {
          name: formData.name,
          location: formData.location,
          description: formData.description
        };
        const updatedDevice = await deviceService.updateDevice(currentDevice.uuid, updateData);
        setDevices(devices.map(device => device.uuid === updatedDevice.uuid ? updatedDevice : device));
        toast.success('Dispositivo atualizado com sucesso');
      } else {
        const newDevice = await deviceService.createDevice(formData);
        setDevices([...devices, newDevice]);
        toast.success('Dispositivo criado com sucesso');
      }
      
      resetForm();
    } catch (error) {
      const errorMessage = getErrorMessage(error);
      setFormError(errorMessage);
    } finally {
      setSubmitting(false);
    }
  };

  const handleEdit = (device: Device) => {
    setFormData({
      name: device.name,
      location: device.location,
      sn: device.sn,
      description: device.description
    });
    setIsEditing(true);
    setCurrentDevice(device);
    setErrors({});
    setFormError('');
    };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Tem certeza que deseja excluir este dispositivo?')) {
      return;
    }

    try {
      await deviceService.deleteDevice(id);
      setDevices(devices.filter(device => device.uuid !== id));
      toast.success('Dispositivo excluído com sucesso');
    } catch (error) {
      const errorMessage = getErrorMessage(error);
      toast.error(errorMessage);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  if (!user) {
    return (
      <div className="container mx-auto p-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">Verificando autenticação...</div>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="container mx-auto p-6">
        <div className="flex items-center justify-center h-64">
          <div className="text-center">Carregando dispositivos...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl sm:text-3xl font-bold tracking-tight transition-transform transform hover:scale-105 active:scale-95 hover:cursor-pointer">Dispositivos</h2>
          <p className="text-muted-foreground">Gerencie seus dispositivos registrados</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>
              {isEditing ? 'Editar Dispositivo' : 'Novo Dispositivo'}
            </CardTitle>
            <CardDescription>
              {isEditing 
                ? 'Atualize as informações do dispositivo' 
                : 'Adicione um novo dispositivo ao sistema'
              }
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              {formError && (
                <div className="text-red-700 text-md text-center py-2">
                  {formError}
                </div>
              )}

              <div className="space-y-2">
                <label htmlFor="name" className="text-sm font-medium">
                  Nome *
                </label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Nome do dispositivo"
                  value={formData.name}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="location" className="text-sm font-medium">
                  Localização *
                </label>
                <Input
                  id="location"
                  name="location"
                  placeholder="Localização do dispositivo"
                  value={formData.location}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="space-y-2">
                <label htmlFor="sn" className="text-sm font-medium">
                  Número de Série *
                </label>
                <Input
                  id="sn"
                  name="sn"
                  placeholder="Número de série (12 dígitos)"
                  value={formData.sn}
                  onChange={handleInputChange}
                  required
                  disabled={isEditing}
                  className={errors.sn ? 'border-red-500' : ''}
                />
                {errors.sn && (
                  <p className="text-xs text-red-500">{errors.sn}</p>
                )}
                {isEditing && (
                  <p className="text-xs text-gray-500">
                    Número de série não pode ser alterado
                  </p>
                )}
              </div>

              <div className="space-y-2">
                <label htmlFor="description" className="text-sm font-medium">
                  Descrição
                </label>
                <Input
                  id="description"
                  name="description"
                  placeholder="Descrição opcional"
                  value={formData.description}
                  onChange={handleInputChange}
                />
              </div>

              <div className="flex gap-2 pt-4">
                <Button 
                  type="submit" 
                  className="flex-1 hover:cursor-pointer"
                  disabled={submitting}
                >
                  {submitting && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
                  {isEditing ? 'Atualizar' : 'Adicionar'}
                </Button>
                
                {isEditing && (
                  <Button 
                    type="button" 
                    variant="outline" 
                    onClick={resetForm}
                    disabled={submitting}
                    className='hover:cursor-pointer'
                  >
                    Cancelar
                  </Button>
                )}
              </div>
            </form>
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader>
            <div className="flex justify-between items-center">
              <div>
                <CardTitle>Dispositivos Registrados</CardTitle>
                <CardDescription>
                  {filteredDevices.length} dispositivo(s) encontrado(s)
                </CardDescription>
              </div>
              <div className="relative w-64">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  placeholder="Buscar dispositivos..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
          </CardHeader>
          <CardContent>
            {filteredDevices.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                {devices.length === 0 
                  ? 'Nenhum dispositivo registrado. Adicione o primeiro dispositivo usando o formulário ao lado.'
                  : 'Nenhum dispositivo encontrado com os critérios de busca.'
                }
              </div>
            ) : (
              <div className="rounded-md border">
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Nome</TableHead>
                      <TableHead>Localização</TableHead>
                      <TableHead>Nº Série</TableHead>
                      <TableHead>Descrição</TableHead>
                      <TableHead>Cadastrado em</TableHead>
                      <TableHead className="w-24">Ações</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {filteredDevices.map((device) => (
                      <TableRow key={device.uuid}>
                        <TableCell className="font-medium">{device.name}</TableCell>
                        <TableCell>
                          <Badge variant="secondary">{device.location}</Badge>
                        </TableCell>
                        <TableCell className="font-mono text-sm">
                          {device.sn}
                        </TableCell>
                        <TableCell className="max-w-xs truncate">
                          {device.description || '-'}
                        </TableCell>
                        <TableCell>
                          {formatDate(device.created_at)}
                        </TableCell>
                        <TableCell>
                          <div className="flex gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleEdit(device)}
                              className='hover:cursor-pointer'
                            >
                              <Pencil className="w-4 h-4" />
                            </Button>
                            <Button
                              variant="destructive"
                              size="sm"
                              onClick={() => handleDelete(device.uuid)}
                              className='hover:cursor-pointer'
                            >
                              <Trash2 className="w-4 h-4" />
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default DeviceRegistrationPage;