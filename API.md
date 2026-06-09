# Vigia API

Base URL: `http://localhost:8080`

Todas as respostas são `application/json`.  
Datas no formato RFC3339 (`2026-06-08T14:30:00Z`).

---

## Monitors

### Criar monitor

```
POST /monitors
```

**Body:**

```json
{
  "name": "Checkout",
  "description": "Onde o comprador finaliza o pedido.",
  "target": "https://loja.com/checkout",
  "type": "checkout",
  "threshold": 3,
  "interval_seconds": 60,
  "acceptable_response_time_seconds": 3.0
}
```

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `name` | string | ✅ | Nome legível do serviço |
| `description` | string | — | Descrição do que esse serviço faz |
| `target` | string | ✅ | URL a ser monitorada |
| `type` | string | ✅ | `uptime`, `checkout` ou `dependency` |
| `threshold` | int | ✅ | Verificações consecutivas ruins para abrir incident |
| `interval_seconds` | int | ✅ | Intervalo entre verificações (segundos) |
| `acceptable_response_time_seconds` | float | ✅ para `checkout` | Tempo máximo de resposta aceitável |

**Resposta `201`:**

```json
{
  "id": "mon-abc123",
  "name": "Checkout",
  "description": "Onde o comprador finaliza o pedido.",
  "target": "https://loja.com/checkout",
  "type": "checkout",
  "status": "active",
  "threshold": 3,
  "interval_seconds": 60,
  "acceptable_response_time_seconds": 3.0
}
```

---

### Listar todos os monitors

```
GET /monitors
```

**Resposta `200`:**

```json
{
  "monitors": [
    {
      "id": "mon-abc123",
      "name": "Checkout",
      "description": "Onde o comprador finaliza o pedido.",
      "target": "https://loja.com/checkout",
      "type": "checkout",
      "status": "active",
      "threshold": 3,
      "interval_seconds": 60,
      "acceptable_response_time_seconds": 3.0,
      "current_state": "slow",
      "last_checked_at": "2026-06-08T14:30:00Z"
    }
  ]
}
```

| Campo | Valores possíveis | Descrição |
|-------|-------------------|-----------|
| `status` | `active`, `paused`, `disabled` | Status operacional do monitor |
| `current_state` | `functioning`, `slow`, `down` | Estado atual derivado de incidents abertos |
| `last_checked_at` | RFC3339 ou `null` | Timestamp da última verificação realizada |

`current_state`:
- `functioning` — sem incident aberto
- `slow` — incident aberto em monitor do tipo `checkout` (lentidão)
- `down` — incident aberto em monitor do tipo `uptime` ou `dependency`

---

### Pausar monitor

```
POST /monitors/{id}/pause
```

Resposta `204` (sem body).

---

### Retomar monitor

```
POST /monitors/{id}/resume
```

Resposta `204` (sem body).

---

### Desabilitar monitor

```
POST /monitors/{id}/disable
```

Resposta `204` (sem body).

---

### Histórico agregado (todos os monitors)

```
GET /monitors/history?days=30
```

**Parâmetros:**

| Param | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `days` | int | — | Janela em dias (padrão: `30`, máx: `365`) |

**Resposta `200`:**

```json
{
  "availability_percentage": 99.26,
  "total_incidents": 4,
  "total_downtime_seconds": 3420.0,
  "daily": [
    { "date": "2026-05-10", "incidents_count": 0, "downtime_seconds": 0 },
    { "date": "2026-05-11", "incidents_count": 1, "downtime_seconds": 2280.0 }
  ]
}
```

| Campo | Descrição |
|-------|-----------|
| `availability_percentage` | Pior disponibilidade entre todos os monitors no período |
| `total_incidents` | Incidents abertos dentro do período |
| `total_downtime_seconds` | Soma de downtime de todos os monitors no período |
| `daily` | Um item por dia, do mais antigo ao mais recente |

---

### Histórico de um monitor

```
GET /monitors/{id}/history?from=2026-06-01T00:00:00Z&to=2026-06-08T23:59:59Z
```

**Parâmetros:**

| Param | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `from` | RFC3339 | ✅ | Início do período |
| `to` | RFC3339 | ✅ | Fim do período |

**Resposta `200`:**

```json
{
  "incidents": [
    {
      "id": "inc-xyz",
      "monitor_id": "mon-abc123",
      "status": "resolved",
      "opened_at": "2026-06-07T10:00:00Z",
      "resolved_at": "2026-06-07T10:38:00Z",
      "duration_seconds": 2280
    }
  ],
  "availability_percentage": 99.26
}
```

---

## Incidents

### Listar incidents

```
GET /incidents
GET /incidents?status=open
GET /incidents?status=resolved
```

**Parâmetros:**

| Param | Valores | Descrição |
|-------|---------|-----------|
| `status` | `open`, `resolved` | Filtra por status. Omitir retorna todos. |

**Resposta `200`:**

```json
{
  "incidents": [
    {
      "id": "inc-xyz",
      "display_id": "INC-42",
      "monitor_id": "mon-abc123",
      "monitor_name": "Recebimento PIX",
      "status": "open",
      "opened_at": "2026-06-08T12:00:00Z",
      "resolved_at": null,
      "duration_seconds": 0
    },
    {
      "id": "inc-abc",
      "display_id": "INC-41",
      "monitor_id": "mon-def456",
      "monitor_name": "Checkout",
      "status": "resolved",
      "opened_at": "2026-06-07T10:00:00Z",
      "resolved_at": "2026-06-07T10:38:00Z",
      "duration_seconds": 2280
    }
  ]
}
```

| Campo | Descrição |
|-------|-----------|
| `display_id` | Identificador legível sequencial (ex.: `INC-42`) |
| `monitor_name` | Nome do monitor associado |
| `duration_seconds` | `0` enquanto `status == "open"` |

---

## Tipos de monitor

| Tipo | Descrição |
|------|-----------|
| `uptime` | Verifica se o serviço está no ar (alcançabilidade) |
| `checkout` | Verifica performance — incident abre se resposta demorar mais que `acceptable_response_time_seconds` |
| `dependency` | Verifica serviços externos (APIs de terceiros, provedores) |

## Regras de negócio relevantes

- Verificações HTTP: `2xx` e `3xx` = sucesso; `4xx` e `5xx` = falha
- `threshold`: número de verificações consecutivas ruins necessárias para abrir um incident
- Um monitor só tem no máximo 1 incident aberto simultaneamente
- Monitor `checkout` com resposta lenta (acima do ART, abaixo do timeout interno `ART × 3`) → incident abre como lentidão
- Incident fecha automaticamente quando `threshold` verificações consecutivas bem-sucedidas são detectadas
