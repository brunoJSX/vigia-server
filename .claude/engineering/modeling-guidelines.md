# Modeling Guidelines

## Objetivo

O objetivo da modelagem é compreender o domínio e representar o negócio da forma mais clara possível.

A modelagem deve ser consequência do domínio.

Não da arquitetura.

Não dos frameworks.

Não dos padrões.

---

## Princípios

### O domínio vem primeiro

Não iniciar pela implementação.

Não iniciar pela estrutura de pastas.

Não iniciar pelo banco de dados.

Primeiro compreender:

* Problema
* Usuários
* Fluxos
* Regras de negócio
* Conceitos centrais

---

### Simplicidade

Escolher sempre o modelo mais simples capaz de representar corretamente o domínio.

Evitar:

* Overdesign
* Abstrações prematuras
* Complexidade desnecessária

---

### DDD é uma ferramenta

DDD é uma ferramenta de descoberta.

Não é obrigatório utilizar:

* Entities
* Aggregates
* Value Objects
* Domain Events

Utilizar apenas quando agregarem clareza.

---

## Possíveis Abordagens

### Entity-Centric

Utilizar quando o domínio for centrado em:

* entidades
* ownership
* lifecycle
* invariantes

Perguntas típicas:

* Quem é o dono?
* Quem pode alterar?
* Qual é o estado atual?

Exemplos:

* Client
* Subscription
* Billing
* User

Estrutura comum:

```text
Client
Subscription
Invoice
```

---

### Workflow-Centric

Utilizar quando o domínio for centrado em:

* processos
* aprovações
* etapas
* transições

Perguntas típicas:

* Qual o próximo passo?
* Quem executa cada etapa?
* Quais são as decisões?

Exemplos:

* Onboarding
* Approval Flow
* Provisionamento

Estrutura comum:

```text
Request
↓
Validation
↓
Approval
↓
Execution
```

---

### Data-Oriented

Utilizar quando o domínio for centrado em:

* coleta
* transformação
* análise
* detecção

Perguntas típicas:

* Qual dado entra?
* Como ele é transformado?
* Qual resultado é produzido?

Exemplos:

* Monitoring
* Observability
* Analytics
* Fraud Detection

Estrutura comum:

```text
Monitor
↓
Collector
↓
Metric
↓
Analyzer
↓
Incident
```

---

## Como Escolher

Durante o Domain Discovery identificar:

1. Conceitos centrais
2. Fluxos principais
3. Regras relevantes
4. Ownership
5. Lifecycle

Somente após isso escolher a abordagem.

Não assumir previamente.

---

## Ownership

Ownership deve existir apenas quando for relevante para o domínio.

Exemplos:

```text
Client
 └── Subscription
```

```text
Monitor
 └── Incident
```

Não criar ownership artificial.

---

## Lifecycle

Nem todo conceito possui lifecycle.

Criar state machines apenas quando existir estado explícito.

Exemplos válidos:

```text
Incident

Open
↓
Resolved
```

```text
Subscription

Trial
↓
Active
↓
Cancelled
```

Evitar state machines para conceitos sem estado próprio.

---

## Eventos

Eventos devem representar fatos do domínio.

Não criar eventos apenas porque uma arquitetura orientada a eventos está sendo utilizada.

Exemplos:

* Incident Opened
* Incident Resolved
* Subscription Activated

---

## Infraestrutura

Infraestrutura não define o domínio.

Exemplos:

* PostgreSQL
* Redis
* Temporal
* Kafka
* RabbitMQ
* Playwright

São detalhes de implementação.

A modelagem deve permanecer válida mesmo que essas tecnologias sejam substituídas.

---

## Temporal

Temporal é uma ferramenta de orquestração.

Não modelar o domínio em função do Temporal.

Exemplo correto:

```text
Monitor
↓
Collector
↓
Analyzer
↓
Incident
```

Temporal executa o fluxo.

Temporal não define o fluxo.

---

## Regra Principal

Modelar o negócio.

Não modelar padrões arquiteturais.

Se uma abstração não melhora o entendimento do domínio, ela provavelmente não deveria existir.
