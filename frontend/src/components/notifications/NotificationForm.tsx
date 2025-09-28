import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Checkbox } from '@/components/ui/checkbox';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Card, CardContent } from '@/components/ui/card';
import { Plus, Trash2, X } from 'lucide-react';
import type { Device, NotificationCondition } from '@/types/index';

type ParameterType = 'cpu' | 'ram' | 'disk_free' | 'temperature' | 'latency' | 'connectivity';
type OperatorType = '>' | '<' | '>=' | '<=' | '==' | '!=';

const notificationSchema = z.object({
  name: z.string().min(1, 'Nome é obrigatório'),
  description: z.string().optional(),
  enabled: z.boolean(),
  device_ids: z.array(z.string()),
  conditions: z.array(
    z.object({
      parameter: z.enum(['cpu', 'ram', 'disk_free', 'temperature', 'latency', 'connectivity']),
      operator: z.enum(['>', '<', '>=', '<=', '==', '!=']),
      value: z.number().min(0)
    })
  ).min(1, 'Pelo menos uma condição é necessária')
});

type NotificationFormData = z.infer<typeof notificationSchema>;

interface NotificationFormProps {
  devices: Device[];
  onSubmit: (data: NotificationFormData) => void;
  onCancel: () => void;
}

export const NotificationForm: React.FC<NotificationFormProps> = ({ 
  devices, 
  onSubmit, 
  onCancel 
}) => {
  const [selectedDevices, setSelectedDevices] = useState<string[]>([]);
  
  const { register, handleSubmit, watch, setValue, formState: { errors } } = useForm<NotificationFormData>({
    resolver: zodResolver(notificationSchema),
    defaultValues: {
      name: '',
      description: '',
      enabled: true,
      device_ids: [],
      conditions: [{ parameter: 'cpu', operator: '>', value: 70 }]
    }
  });

  const conditions = watch('conditions') || [];
  const enabledValue = watch('enabled');

  const addCondition = () => {
    const newCondition: NotificationCondition = { 
      parameter: 'cpu', 
      operator: '>', 
      value: 70 
    };
    const newConditions = [...conditions, newCondition];
    setValue('conditions', newConditions);
  };

  const removeCondition = (index: number) => {
    const newConditions = conditions.filter((_, i) => i !== index);
    setValue('conditions', newConditions);
  };

  const updateCondition = (index: number, field: keyof NotificationCondition, value: any) => {
    const updatedConditions = [...conditions];
    updatedConditions[index] = { 
      ...updatedConditions[index], 
      [field]: field === 'value' ? Number(value) : value 
    };
    setValue('conditions', updatedConditions);
  };

  const toggleDevice = (deviceUuid: string) => {
    const updated = selectedDevices.includes(deviceUuid)
      ? selectedDevices.filter(uuid => uuid !== deviceUuid)
      : [...selectedDevices, deviceUuid];
    
    setSelectedDevices(updated);
    setValue('device_ids', updated);
  };

  const onSubmitForm = handleSubmit((data) => {
    onSubmit(data);
  });

  return (
    <form onSubmit={onSubmitForm} className="space-y-6">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="space-y-4">
          <div>
            <Label htmlFor="name">Nome da Notificação *</Label>
            <Input
              id="name"
              {...register('name')}
              placeholder="Ex: Alerta de CPU Alta"
            />
            {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
          </div>

          <div>
            <Label htmlFor="description">Descrição</Label>
            <Textarea
              id="description"
              {...register('description')}
              placeholder="Descreva o propósito desta notificação..."
            />
          </div>

          <div className="flex items-center space-x-2">
            <Checkbox
              id="enabled"
              checked={enabledValue}
              onCheckedChange={(checked) => setValue('enabled', checked === true)}
              className='hover:cursor-pointer'
            />
            <Label htmlFor="enabled">Notificação ativa</Label>
          </div>
        </div>

        <div className="space-y-4">
          <div>
            <Label>Dispositivos</Label>
            <p className="text-sm text-muted-foreground mb-2">
              Selecione dispositivos específicos ou deixe vazio para todos
            </p>
            <Card>
              <CardContent className="p-4 max-h-40 overflow-y-auto">
                {devices.map((device: Device) => (
                  <div key={device.uuid} className="flex items-center space-x-2 py-1">
                    <Checkbox
                      id={`device-${device.uuid}`}
                      checked={selectedDevices.includes(device.uuid)}
                      onCheckedChange={() => toggleDevice(device.uuid)}
                      className='hover:cursor-pointer'
                    />
                    <Label htmlFor={`device-${device.uuid}`} className="flex-1 text-sm">
                      {device.name} ({device.sn})
                    </Label>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Label>Condições de Alerta *</Label>
          <Button type="button" variant="outline" size="sm" className='hover:cursor-pointer' onClick={addCondition}>
            <Plus className="h-4 w-4 mr-1" />
            Adicionar Condição
          </Button>
        </div>

        {conditions.map((condition, index) => (
          <Card key={index}>
            <CardContent className="p-4">
              <div className="flex flex-col sm:flex-row gap-3 sm:items-end">
                <div className="flex-1 min-w-0">
                  <Label className="text-sm">Parâmetro</Label>
                  <Select
                    value={condition.parameter}
                    onValueChange={(value: ParameterType) => 
                      updateCondition(index, 'parameter', value)
                    }
                  >
                    <SelectTrigger className='cursor-pointer text-sm h-9'>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="cpu" className="text-sm">Uso de CPU (%)</SelectItem>
                      <SelectItem value="ram" className="text-sm">Uso de RAM (%)</SelectItem>
                      <SelectItem value="disk_free" className="text-sm">Disco Livre (GB)</SelectItem>
                      <SelectItem value="temperature" className="text-sm">Temperatura (°C)</SelectItem>
                      <SelectItem value="latency" className="text-sm">Latência (ms)</SelectItem>
                      <SelectItem value="connectivity" className="text-sm">Conectividade</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex-1 min-w-0">
                  <Label className="text-sm">Operador</Label>
                  <Select
                    value={condition.operator}
                    onValueChange={(value: OperatorType) => 
                      updateCondition(index, 'operator', value)
                    }
                  >
                    <SelectTrigger className='cursor-pointer text-sm h-9'>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value=">" className="text-sm">Maior que</SelectItem>
                      <SelectItem value="<" className="text-sm">Menor que</SelectItem>
                      <SelectItem value=">=" className="text-sm">Maior ou igual</SelectItem>
                      <SelectItem value="<=" className="text-sm">Menor ou igual</SelectItem>
                      <SelectItem value="==" className="text-sm">Igual a</SelectItem>
                      <SelectItem value="!=" className="text-sm">Diferente de</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div className="flex-1 min-w-0">
                  <Label className="text-sm">Valor</Label>
                  <Input
                    type="number"
                    value={condition.value}
                    onChange={(e) => updateCondition(index, 'value', parseFloat(e.target.value) || 0)}
                    step="0.1"
                    className="h-9 text-sm"
                    min="0"
                  />
                </div>

                <div className="flex sm:flex-none justify-end sm:justify-start pt-1 sm:pt-0">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeCondition(index)}
                    className='hover:cursor-pointer h-9 w-9 p-0'
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        ))}

        {errors.conditions && (
          <p className="text-sm text-red-500">{errors.conditions.message}</p>
        )}
      </div>

      <div className="flex flex-col sm:flex-row justify-end gap-2">
        <Button type="button" variant="outline" onClick={onCancel} className='hover:cursor-pointer order-2 sm:order-1'>
          <X className="h-4 w-4 mr-1" />
          Cancelar
        </Button>
        <Button type="submit" className='hover:cursor-pointer order-1 sm:order-2'>
          Criar Notificação
        </Button>
      </div>
    </form>
  );
};