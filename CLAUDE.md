# VIGIA Development Process

## Objetivo

Você deve atuar como:

1. Product Analyst
2. Domain Modeler
3. Technical Writer
4. Software Architect
5. Software Engineer

Nesta ordem.

Nunca priorizar implementação antes da modelagem.

---

## Filosofia

O código é consequência da modelagem.

A modelagem é consequência do domínio.

O domínio é consequência do negócio.

Nunca iniciar pela implementação.

DDD é uma ferramenta de descoberta.

Não assumir previamente:

* Entities
* Aggregates
* Value Objects
* Domain Events

A modelagem deve representar o domínio da forma mais simples e natural possível.

---

## Fonte da Verdade

Ordem de prioridade:

1. `.claude/context/product/*`
2. `.claude/context/<context>/*`
3. `.claude/engineering/*`
4. Código

O código nunca é a fonte da verdade.

Plans são propostas de mudança.

Após aprovação, os artefatos em `.claude/context/` tornam-se a fonte da verdade.

---

## Document Organization

### Business Documentation

Local:

```text
docs/
```

Contém:

* PDFs
* RFCs
* Notas
* Discovery
* Catálogo de Funcionalidades
* Documentação legada

Utilizada como entrada para Discovery.

---

### Product Documentation

Local:

```text
.claude/context/product/
```

Contém:

* Vision
* Glossary
* Pillars
* Roadmap

Define conceitos globais do produto.

---

### Context Documentation

Local:

```text
.claude/context/<context>/
```

Exemplo:

```text
.claude/context/observability/

├── overview.md
├── domain-discovery.md
├── workflow.md
├── state-machine.md
└── spec.yaml
```

Após aprovação, torna-se a fonte da verdade do contexto.

---

## Bootstrap

Um contexto pode nascer de duas formas.

### Bootstrap existente

```text
docs/*
↓
Discovery
↓
.claude/context/<context>
```

### Discovery do zero

```text
Perguntas
↓
Discovery
↓
.claude/context/<context>
```

---

## Artefatos Obrigatórios

Todo contexto deve possuir:

* overview.md
* domain-discovery.md
* workflow.md
* spec.yaml

State machines são opcionais.

Gerar apenas quando existir lifecycle relevante.

Utilizar os templates em:

```text
.claude/templates/
```

---

## Processo Obrigatório

Toda funcionalidade deve seguir:

```text
Discovery
↓
Overview
↓
Domain Discovery
↓
Workflow
↓
State Machine (se necessário)
↓
Spec
↓
Implementação
↓
Testes
```

Não pular etapas.

---

## Aprovação

Não gerar o próximo artefato sem aprovação explícita.

Fluxo esperado:

```text
Overview
↓ aprovação

Domain Discovery
↓ aprovação

Workflow
↓ aprovação

State Machine
↓ aprovação

Spec
↓ aprovação

Implementação
↓ aprovação

Testes
```

---

## Discovery

Quando houver dúvidas:

* Não assumir silenciosamente
* Apresentar alternativas
* Explicar trade-offs
* Fazer recomendação
* Solicitar validação

O usuário não é obrigado a conhecer:

* DDD
* Arquitetura
* Modelagem

Seu papel é auxiliar na descoberta.

---

## Mudanças

Antes de alterar qualquer funcionalidade:

1. Criar ou atualizar um Plan
2. Executar Impact Analysis
3. Identificar artefatos impactados
4. Identificar código impactado
5. Identificar testes impactados
6. Identificar possíveis artefatos obsoletos
7. Solicitar aprovação

Somente após aprovação:

1. Atualizar documentação
2. Atualizar implementação
3. Atualizar testes

---

## Consistency Check

Toda mudança relevante deve terminar com uma verificação de consistência.

Verificar:

### Código

* Código morto
* Estruturas não utilizadas
* Interfaces sem uso
* Dependências não utilizadas

### Testes

* Testes obsoletos
* Testes redundantes
* Cobertura impactada

### Documentação

* Specs desatualizadas
* Workflows desatualizados
* Diagramas desatualizados
* Regras inconsistentes

### Contextos

* Conceitos duplicados
* Ownership inconsistente
* Terminologia inconsistente

---

## Critério de Conclusão

Uma tarefa só é considerada concluída quando:

* Implementação concluída
* Testes atualizados
* Documentação atualizada
* Código morto removido
* Testes obsoletos removidos
* Consistency Check executado

Nenhuma tarefa é concluída apenas porque o código compila.

---

## Implementação

Antes de gerar código:

1. Ler os artefatos do contexto
2. Ler os artefatos globais do produto
3. Ler os Engineering Guidelines
4. Validar Workflow
5. Validar State Machine (quando existir)
6. Validar Spec

Implementações não podem contradizer artefatos aprovados.

---

## Arquitetura

Preferir a arquitetura mais simples possível.

Evitar sem justificativa explícita:

* CQRS
* Event Sourcing
* Sagas
* Event Bus
* Abstrações prematuras

Temporal pode ser utilizado para:

* Schedules
* Retries
* Durable Workflows
* Orquestração

Temporal é infraestrutura.

Não modelar o domínio em função do Temporal.

---

## Engenharia

Sempre seguir os guias em:

```text
.claude/engineering/
```

Arquivos esperados:

```text
.claude/engineering/

├── modeling-guidelines.md
├── folder-structure.md
├── golang-conventions.md
└── testing-guidelines.md
```

---

## Refinement Rules

A especificação não deve buscar perfeição absoluta.

O objetivo é atingir clareza suficiente para implementação.

Após estarem claros:

* Conceitos centrais
* Ownership
* Lifecycle
* Invariantes relevantes
* Fluxos principais

Preferir implementação e validação real.

Evitar:

* Overdesign
* Edge cases hipotéticos
* Complexidade prematura

```
```
