# GoPayFlow

**GoPayFlow** é uma arquitetura backend de processamento de pagamentos construída em **Golang** com foco em **Clean Architecture**, **processamento assíncrono**, **escalabilidade horizontal** e **infraestrutura containerizada**.

O projeto simula o fluxo real de um **sistema de pagamentos**, separando claramente:

* API de entrada
* fila de processamento
* workers assíncronos
* banco de dados
* balanceamento de carga

O objetivo não é apenas criar uma API, mas **demonstrar uma arquitetura de backend completa**, preparada para crescimento e alto volume de requisições.

---

# Arquitetura do Sistema

O sistema segue o seguinte fluxo:

```
Client (Frontend)
        ↓
   Nginx (Load Balancer)
        ↓
   API (replicas)
        ↓
   PostgreSQL (persistência)
        ↓
   Redis / RabbitMQ (fila)
        ↓
---------------------------
        ↓
        Worker
        ↓
 Simulação de pagamento
        ↓
 Atualização do pedido
```

### Explicação

1. O cliente envia uma requisição para criar um pagamento.
2. O **Nginx distribui a carga** entre múltiplas instâncias da API.
3. A **API registra o pedido no banco** com status `PENDING`.
4. A API publica um evento em uma **fila de processamento**.
5. Um **Worker assíncrono** consome a fila.
6. O Worker simula o processamento do pagamento.
7. O Worker atualiza o status do pedido no banco.

Essa arquitetura é amplamente utilizada em sistemas reais de pagamento.

---

# Componentes da Arquitetura

## Client

Interface que consome a API.

Pode ser:

* Web App
* Mobile App
* Sistemas terceiros
* Integrações via REST

---

## Nginx (Load Balancer)

Responsável por:

* balancear requisições
* distribuir tráfego entre múltiplas APIs
* aumentar disponibilidade

Exemplo:

```
Client
   ↓
 Nginx
  ↓ ↓ ↓ ↓
API API API API
```

---

## API

A API é responsável por:

* receber requisições HTTP
* validar dados
* registrar pedidos
* publicar eventos na fila
* retornar resposta imediata ao cliente

Tecnologias utilizadas:

* **Golang**
* **Gin**
* **Clean Architecture**

A API é **stateless**, permitindo escalabilidade horizontal.

---

## PostgreSQL

Banco de dados principal do sistema.

Responsável por armazenar:

* pedidos
* status do pagamento
* histórico

Tabela principal:

```
orders
```

Campos principais:

| Campo      | Tipo      |
| ---------- | --------- |
| id         | UUID      |
| amount     | NUMERIC   |
| status     | VARCHAR   |
| created_at | TIMESTAMP |

Status possíveis:

```
PENDING
APPROVED
FAILED
```

---

## Sistema de Fila

A fila desacopla o **tempo de resposta da API** do **processamento do pagamento**.

Isso permite:

* maior throughput
* processamento assíncrono
* resiliência
* escalabilidade

No projeto são utilizados:

* RabbitMQ ou Redis

A API envia mensagens para a fila e os Workers consomem.

---

## Worker

O Worker é responsável pelo processamento assíncrono.

Funções:

* consumir mensagens da fila
* simular processamento do pagamento
* atualizar status no banco

Fluxo:

```
Fila → Worker → Banco
```

Exemplo de simulação:

```
sleep(2s)
status = APPROVED
```

Workers podem ser escalados horizontalmente.

---

# Estrutura do Projeto

```
gopayflow
│
├── backend
│
│   ├── cmd
│   │   └── api
│   │
│   ├── internal
│   │
│   │   ├── domain
│   │   │
│   │   ├── application
│   │   │
│   │   ├── infrastructure
│   │   │   ├── database
│   │   │   └── messaging
│   │   │
│   │   └── interfaces
│   │       └── http
│   │
│   ├── migrations
│   │
│   ├── Dockerfile
│   ├── go.mod
│   │
│
├── docker-compose.yml
│
└── nginx
```

---

# Clean Architecture

O projeto segue **Clean Architecture**, separando responsabilidades.

Camadas principais:

```
Interfaces
    ↓
Application
    ↓
Domain
    ↓
Infrastructure
```

---

## Domain

Contém as **entidades e contratos do sistema**.

Não depende de nenhuma outra camada.

Exemplo:

```
Order
OrderRepository
```

---

## Application

Implementa os **casos de uso do sistema**.

Exemplos:

* CreateOrder
* GetOrderByID
* ListOrders

Essa camada orquestra a lógica de negócio.

---

## Infrastructure

Responsável por integrações externas.

Exemplos:

* PostgreSQL
* RabbitMQ
* Redis
* serviços externos

---

## Interfaces

Responsável pela comunicação com o mundo externo.

Exemplo:

```
HTTP Handlers
```

Recebem requests e chamam os UseCases.

---

# API Endpoints

### Criar pedido

```
POST /orders
```

Body:

```
{
  "amount": 100
}
```

Resposta:

```
{
  "id": "uuid"
}
```

---

### Listar pedidos

```
GET /orders
```

---

### Buscar pedido

```
GET /orders/:id
```

---

# Infraestrutura

Toda a stack roda em containers.

Containers utilizados:

* API
* PostgreSQL
* RabbitMQ
* Migration Runner

---

# Docker Compose

O ambiente pode ser iniciado com:

```
docker compose up --build
```

Serviços disponíveis:

| Serviço     | Porta |
| ----------- | ----- |
| API         | 8080  |
| PostgreSQL  | 5432  |
| RabbitMQ    | 5672  |
| RabbitMQ UI | 15672 |

---

# Escalabilidade

O sistema foi projetado para escalar horizontalmente.

Exemplo:

```
        Nginx
         ↓
  API  API  API  API
         ↓
      PostgreSQL
         ↓
        Queue
         ↓
  Worker Worker Worker
```

---

# Benefícios da Arquitetura

* desacoplamento entre API e processamento
* melhor throughput
* resiliência
* facilidade de escalar workers
* APIs stateless
* separação clara de responsabilidades

---

# Próximos Passos

Possíveis evoluções do sistema:

* rate limiting
* autenticação
* observabilidade
* métricas
* tracing distribuído
* retry automático de pagamentos
* idempotência de requisições
* dead letter queues
* monitoramento de filas
* testes automatizados
* CI/CD

---

# Objetivo do Projeto

Este projeto demonstra como construir um backend moderno em Go com:

* arquitetura limpa
* mensageria
* processamento assíncrono
* infraestrutura containerizada
* escalabilidade horizontal

Ele serve como **base para sistemas reais de pagamentos ou processamento de eventos**.

---
