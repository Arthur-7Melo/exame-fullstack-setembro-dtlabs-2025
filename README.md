# IoT Monitoring — README

> Documentação para rodar o case técnico (frontend React + backend Go + simulator + infra)

---

## Sumário

- Visão geral
- Arquitetura
- Pré-requisitos
- Variáveis de ambiente (exemplo `.env`)
- Como rodar (Docker Compose)
- Rodando localmente (backend, frontend, simulator)
- Registrar devices e atualizar simulator
- Endpoints importantes / Swagger
- Testes
- Como testar notificações em tempo real
- Dicas de troubleshooting

---

## Visão geral

Plataforma para monitoramento de dispositivos IoT. Dispositivos (reais ou simulados) enviam *heartbeats* (telemetria) a cada 1 minuto com os seguintes campos principais:

- `cpu` (%)
- `ram` (%)
- `disk_free` (% disponível)
- `temperature` (°C)
- `latency_ms` (latência para 8.8.8.8 em ms)
- `connectivity` (0 ou 1)
- `boot_time` (timestamp UTC+00)

Fluxo simplificado (como implementado no projeto):

1. Dispositivo/Simulator → publica em RabbitMQ (fila `heartbeats`)
2. HeartbeatConsumer (Go) consome a fila
3. HeartbeatConsumer salva no PostgreSQL
4. HeartbeatConsumer chama NotificationService (avalida regras)
5. Se regra satisfeita → publica em Redis (canal `notifications`)
6. Redis → notifica frontend via WebSocket
7. Frontend exibe notificações em tempo real

---

## Arquitetura (resumida)

```
[Simulator(s)] -> RabbitMQ (heartbeats queue)
                       |
                       v
                HeartbeatConsumer (Go)
                 /           \
                v             v
          PostgreSQL        NotificationService
                                |
                                v
                             Redis PUB/SUB
                                |
                                v
                            Frontend (WebSocket)
```

---

## Pré-requisitos

- Docker & Docker Compose
- Go >= 1.20 (para rodar localmente se necessário)
- Node.js >= 18 (para frontend)

> Observação: o repositório já contém `Dockerfile` e `docker-compose.yml` para rodar tudo em containers.

---

## Exemplo de `.env`

Coloque este arquivo na raiz do monorepo (onde está o `docker-compose.yml`).

```.env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=iotplatform
SSL_MODE=disable

# App
PORT=8080
JWT_SECRET=62774aa06a16f84f7acefe1c0be66aca07b665743eb459f90db56afd4deace4b

#Rabbitmq
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
AMQP_URL=amqp://guest:guest@rabbitmq:5672/

#Redis
REDIS_URL=redis://redis:6379

#DEVIDE_TESTS_FOR_SIMULATOR
DEVICE_IDS=uuid1,uuid2...

```

> **Importante:** nunca versionar senhas reais. Versionado apenas para facilitar testes de case-técnico. Mantenha `.env` no `.gitignore`.

---

## Como rodar (Docker Compose)

1. Clone o repositório com: `git clone https://github.com/Arthur-7Melo/exame-fullstack-setembro-dtlabs-2025.git`

2. Verifique se o `.env` está no diretório raíz do projeto

2. Subir todos os serviços (app, frontend, postgres, rabbitmq, redis, simulator):

```bash
docker-compose up --build
```

3. Verifique logs (se necessário):

```bash
docker-compose logs -f iotplatform_app
docker-compose logs -f iotplatform_simulator
```

4. Para derrubar os containers:

```bash
docker-compose down
```

5. O Frontend estará rodando na porta `:5173`. Acesse: `http://localhost:5173` para ter acesso a UI interativa
---

## Rodando localmente (sem Docker)

### Backend (Go)

```bash
# instale dependências
go mod download

# executar servidor (ex: cmd/server)
cd cmd/server
go run main.go
```

> O server por padrão roda na porta `:8080` (ajuste com `PORT` no .env).

### Frontend (React + Vite)

```bash
cd frontend
npm install
npm run dev
# ou para build: npm run build
```

> O Frontend roda na porta `:5173`. Acesse em: `http://localhost:5173`

### Simulator (Go)

O projeto inclui `cmd/simulator/main.go` e um `Dockerfile.simulator` para rodar via container. Se quiser rodar local:

```bash
cd cmd/simulator
go run main.go
```

O simulator lê `DEVICE_IDS` do `.env` para saber quais devices simular. Ele publica heartbeats com pequena aleatoriedade para gerar eventos de notificação reais.

---

## Registrar devices e atualizar simulator

1. Crie/registre device(s) via frontend ou API (rota privada): `/api/v1/devices` (CRUD).
2. Após cadastrar os devices, pare a stack do docker-compose para atualizar o `SIMULATOR_DEVICE_IDS` no `.env` com os `DEVICE_IDS` gerados.

```bash
# pare a stack
docker-compose down

# editar .env -> Remova o "#" do DEVICE_IDS, acrescente o "=" e cole os IDS separados por ","
# ex: DEVICE_IDS=11111111-1111-1111-1111-111111111111,22222222-2222-2222-2222-222222222222

# subir novamente
docker-compose up --build
```

> O motivo: o simulator carrega os IDs na inicialização (por isso precisa reiniciar para pegar os novos ids). Alternativamente o simulator pode aceitar reload via sinal ou endpoint, mas a instrução acima cumpre o requisito do teste.

---

## Endpoints importantes / Swagger

A documentação swagger roda em:

```
http://localhost:8080/swagger/*any
```

Principais endpoints (exemplos):

- `POST /api/auth/register` — registrar usuário (body: `email`, `password`) (retorna jwt)
- `POST /api/auth/login` — autenticar (retorna JWT)
- `GET /api/v1/devices` — listar devices do usuário
- `POST /api/v1/devices` — criar device
- `GET /api/v1/devices/:id/heartbeats` — listar heartbeats
- `POST /api/v1/notifications` — criar regra de notificação
- WebSocket: `ws://localhost:8080/ws/notifications?user_id=<USER_UUID>` — conexão para receber notificações em tempo real

---

## Payload de exemplo — Heartbeat

```json
{
  "device_id": "11111111-1111-1111-1111-111111111111",
  "sn": "123456789012",        
  "cpu": 78.4,
  "ram": 65.1,
  "disk_free": 24.5,
  "temperature": 73.2,
  "latency_ms": 22,
  "connectivity": 1,
  "boot_time": "2025-09-29T08:00:00Z",
  "timestamp": "2025-09-29T09:01:00Z"
}
```

---

## Testes

- Testes unitários do backend (handlers e services):

```bash
# rodar todos os testes go
go test ./... -run Test

# ou dentro do package específico
go test ./internal/services/ -v
go test ./internal/handlers/ -v
```

---

## Como testar notificações em tempo real

1. Abra o frontend e faça login com um usuário
2. Vá em `Notifications` e crie uma nova regra. Exemplo: CPU &gt; 70% para device X
3. Garanta que o simulator está rodando e que `DEVICE_IDS` contém o device alvo.
4. O simulator envia heartbeats com aleatoriedade: quando um heartbeat exceder a regra, o backend publica a notificação no Redis e o frontend conectado via WebSocket receberá o evento.
5. Verifique logs do backend para confirmar que a regra foi avaliada.

Dicas de debug:

- Verifique se o HeartbeatConsumer está consumindo a fila `heartbeats` no RabbitMQ.
- Verifique se os registros estão sendo salvos no PostgreSQL (tabela `heartbeats`).
- Verifique se a NotificationService publicou em `notifications` (Redis) quando condição satisfeita.
- Verifique se o frontend está conectado ao WebSocket correto e autenticado.

---

## Boas práticas implementadas

- Senhas hashed (bcrypt)
- Uso correto de status codes e tratamento de erro centralizado.
- Endpoints documentados por Swagger.
- Testes unitários nos packages `handlers` e `services`.

---


## Troubleshooting rápido

- **Erro 404 em `/api/...`**: verifique `PORT` no `.env` e logs do backend.
- **WebSocket fecha antes da conexão**: confirmar URL do WS, query param `user_id`, e se o token/JWT está sendo enviado se necessário.
- **Simulator não envia heartbeats**: conferir variáveis `DEVICE_IDS` e reiniciar o container do simulator.

---


