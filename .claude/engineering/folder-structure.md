# Folder Structure

## Objetivo

A estrutura de diretórios deve refletir o domínio e as responsabilidades do sistema.

Não organizar o código em função de padrões arquiteturais.

Organizar o código em função dos conceitos identificados durante a modelagem.

---

## Estrutura Raiz

```text
cmd/
internal/
```

---

## cmd/

Contém os entrypoints da aplicação.

Responsabilidades:

* bootstrap
* dependency injection
* carregamento de configuração
* inicialização de infraestrutura
* inicialização de workers
* inicialização do HTTP server

Não deve conter:

* regras de negócio
* casos de uso
* lógica de domínio

Exemplo:

```text
cmd/

└── vigia/
    └── main.go
```

---

## internal/

Contém toda a implementação da aplicação.

Organizado por contexto.

Exemplo:

```text
internal/

├── observability/
├── notifications/
├── billing/
└── shared/
```

---

## Organização por Contexto

Cada contexto deve ser isolado.

Estrutura mínima:

```text
internal/<context>/

├── application/
├── infrastructure/
└── interfaces/
```

Os demais diretórios devem surgir da modelagem.

---

## application/

Responsável por:

* casos de uso
* coordenação
* orquestração

Pode utilizar:

* repositories
* services
* gateways
* collectors
* analyzers

Não deve conter:

* acesso direto ao banco
* acesso direto a APIs externas

Exemplos:

```text
application/

├── create_monitor.go
├── pause_monitor.go
├── disable_monitor.go
└── query_incidents.go
```

---

## infrastructure/

Responsável por detalhes técnicos.

Exemplos:

* PostgreSQL
* Redis
* Temporal
* Playwright
* HTTP Clients
* SMTP
* WhatsApp API

Exemplo:

```text
infrastructure/

├── persistence/
├── temporal/
├── probes/
├── gateways/
└── clients/
```

---

## interfaces/

Responsável por entradas e saídas do sistema.

Exemplos:

* HTTP
* CLI
* gRPC
* Webhooks

Não deve conter regras de negócio.

Exemplo:

```text
interfaces/

└── http/
```

---

## Contextos Data-Oriented

Quando o domínio for centrado em coleta, transformação e análise de dados.

Exemplo:

```text
internal/observability/

├── monitor/
├── collector/
├── analyzer/
├── incident/
├── application/
├── infrastructure/
└── interfaces/
```

Fluxo típico:

```text
Monitor
↓
Collector
↓
Analyzer
↓
Incident
```

---

## Contextos Entity-Centric

Quando o domínio for centrado em entidades, ownership e lifecycle.

Exemplo:

```text
internal/billing/

├── subscription/
├── invoice/
├── payment/
├── application/
├── infrastructure/
└── interfaces/
```

Não é obrigatório utilizar:

```text
domain/
entities/
value_objects/
aggregates/
```

Utilizar apenas quando agregarem clareza.

---

## Shared

Criar apenas para componentes realmente compartilhados.

Exemplos:

```text
internal/shared/

├── clock/
├── logger/
├── uuid/
└── errors/
```

Evitar mover código para shared prematuramente.

---

## Temporal

Temporal pertence à infraestrutura.

Exemplo:

```text
infrastructure/

└── temporal/
    ├── workflows/
    ├── activities/
    └── schedules/
```

O domínio não deve conhecer Temporal.

---

## Regra Principal

A estrutura deve refletir o domínio.

Não criar diretórios porque um padrão recomenda.

Criar diretórios apenas quando existir uma responsabilidade clara para eles.
