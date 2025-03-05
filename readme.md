# Como criar um bate-papo em tempo real utilizando microsserviços

Para criar um bate-papo em tempo real usando microsserviços, é necessário dividir as funcionalidades em serviços independentes, escaláveis e acoplados de forma flexível. Aqui está uma abordagem detalhada:

---

## 1. Arquitetura Geral

Divida o sistema em microsserviços especializados, utilizando tecnologias adequadas para cada função:

| **Microsserviço**        | **Responsabilidade**                          | **Tecnologias Sugeridas**                |
| ------------------------ | --------------------------------------------- | ---------------------------------------- |
| **WebSocket Service**    | Gerar conexões em tempo real (WebSocket)      | Node.js + Socket.io, WebSocket API (AWS) |
| **Message Service**      | Armazenar/recuperar mensagens                 | MongoDB, PostgreSQL, Redis               |
| **Presence Service**     | Gerenciar status de usuários (online/offline) | Redis, Cassandra                         |
| **Auth Service**         | Autenticar usuários e validar tokens          | JWT, OAuth2, Keycloak                    |
| **Notification Service** | Notificar eventos (ex: novas mensagens)       | RabbitMQ, Kafka, AWS SNS                 |
| **API Gateway**          | Roteamento e gerenciamento de requisições     | NGINX, Kong, Spring Cloud Gateway        |

---

## 2. Fluxo de Funcionamento

### **a. Conexão Inicial do Usuário**

1. **Autenticação**:
   - O cliente envia credenciais para o **Auth Service** via HTTP (ex: login).
   - O Auth Service retorna um **JWT** para acesso aos demais serviços.
2. **Conexão WebSocket**:
   - O cliente estabelece uma conexão WebSocket com o **WebSocket Service**, incluindo o JWT no handshake.
   - O WebSocket Service valida o token com o **Auth Service**.

### **b. Envio de Mensagem**

1. O cliente envia uma mensagem via WebSocket.
2. O **WebSocket Service** publica a mensagem em um **message broker** (ex: Kafka/RabbitMQ).
3. O **Message Service** consome a mensagem do broker e a persiste no banco de dados.
4. O **WebSocket Service** distribui a mensagem em tempo real para os destinatários via WebSocket.

### **c. Status de Presença**

1. Quando um usuário se conecta, o **Presence Service** atualiza seu status para "online" (ex: usando Redis para armazenar estado).
2. Ao desconectar, o WebSocket Service notifica o **Presence Service** para atualizar o status para "offline".

### **d. Notificações**

- O **Notification Service** escuta eventos do broker (ex: nova mensagem) e envia notificações push (ex: via Firebase Cloud Messaging).

---

## 3. Comunicação Entre Serviços

- **Síncrona (HTTP/REST/gRPC)**:
  - Validação de tokens, consulta de histórico de mensagens.
- **Assíncrona (Message Broker)**:
  - Eventos de mensagens, atualizações de status e notificações.
- **WebSocket**:
  - Comunicação em tempo real entre cliente e servidor.

---

## 4. Escalabilidade e Tolerância a Falhas

- **WebSocket Service**:
  - Use **Redis Pub/Sub** ou **ElastiCache** para sincronizar conexões WebSocket em múltiplas instâncias.
- **Message Broker**:
  - Utilize clusters do Kafka ou RabbitMQ para garantir entrega de mensagens.
- **Banco de Dados**:
  - Escalone horizontalmente (ex: MongoDB Sharding) ou use bancos otimizados para leitura/escrita (ex: Cassandra para presença).
- **API Gateway**:
  - Balanceamento de carga e rate limiting para evitar sobrecarga.

---

## 5. Segurança

- **WebSocket Secure (WSS)**:
  - Criptografe a comunicação WebSocket com TLS.
- **Validação de Tokens**:
  - O **Auth Service** deve validar JWT em todas as requisições.
- **Proteção de Dados**:
  - Criptografe mensagens sensíveis (ex: end-to-end encryption).

---

## 6. Monitoramento e Logs

- **Ferramentas**:
  - Prometheus + Grafana (métricas), ELK Stack (logs), Jaeger (tracing distribuído).
- **Métricas Chave**:
  - Latência de mensagens, conexões ativas, taxa de erros.

---

## 7. Exemplo de Implementação

```text
Cliente (Web/Mobile)
  │
  ├── HTTP → API Gateway → Auth Service (valida JWT)
  │
  └── WebSocket → WebSocket Service (Node.js + Socket.io)
        │
        ├── Publica mensagem → Kafka → Message Service (persiste no MongoDB)
        │
        ├── Atualiza status → Redis (Presence Service)
        │
        └── Notificação → Firebase (Notification Service)
```

![image](./architecture.png)
