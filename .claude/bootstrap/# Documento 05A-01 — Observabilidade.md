# Documento 05A-01 — Observabilidade

## Objetivo

O contexto de Observabilidade é responsável por detectar problemas operacionais, registrar incidentes, notificar usuários e fornecer visibilidade sobre a saúde da operação.

Seu objetivo é permitir que o influenciador descubra problemas antes que eles impactem significativamente suas vendas e campanhas.

---

# Escopo

Este contexto contempla:

* Monitoramento de Uptime
* Recuperação de Serviço
* Monitoramento de Performance do Checkout
* Monitoramento de Dependências Externas
* Histórico de Incidentes
* Resumo Diário Operacional

---

# F001 — Monitoramento de Uptime

## Objetivo

Detectar indisponibilidade do site monitorado e alertar o cliente rapidamente para minimizar impactos operacionais e financeiros.

## Problemas Resolvidos

* P00 — Descobrir problemas tarde demais
* P01 — Site fora do ar

## Jornadas Relacionadas

* J03 — Monitoramento Diário
* J04 — Resposta a Incidentes
* J06 — Recebimento de Alertas
* J08 — Investigação de Problemas

## Ator Principal

Sistema

## Gatilho

Execução automática periódica do monitor.

## Entradas

### Obrigatórias

* URL monitorada

### Configuráveis

* Intervalo de verificação
* Timeout
* Quantidade de falhas consecutivas
* Canais de notificação

## Saídas

### Operacionais

* Registro de verificação

### Negócio

* Incidente
* Alerta

## Fluxo Principal

1. Sistema executa verificação da URL.
2. Site responde dentro dos critérios esperados.
3. Resultado é registrado.
4. Nenhum incidente é criado.

## Fluxo Alternativo — Falha Detectada

1. Sistema executa verificação.
2. Site não responde corretamente.
3. Falha é registrada.
4. Sistema verifica quantidade de falhas consecutivas.
5. Threshold atingido.
6. Incidente é aberto.
7. Notificação é enviada.

## Fluxo Alternativo — Incidente Já Aberto

1. Sistema executa nova verificação.
2. Site continua indisponível.
3. Resultado é registrado.
4. Incidente existente é atualizado.
5. Nenhum novo incidente é criado.

## Regras de Negócio

### RN-001

Um incidente não deve ser criado na primeira falha.

### RN-002

Só pode existir um incidente ativo por monitor.

### RN-003

Todas as verificações devem ser armazenadas para histórico.

### RN-004

Toda abertura de incidente deve gerar uma notificação.

### RN-005

Toda recuperação deve encerrar o incidente ativo.

## Critérios de Aceitação

### CA-001

Dado um monitor ativo

Quando o site responder corretamente

Então a verificação deve ser registrada como sucesso.

### CA-002

Dado um monitor ativo

Quando ocorrerem falhas consecutivas suficientes para atingir o threshold

Então um incidente deve ser criado.

### CA-003

Dado um incidente ativo

Quando novas falhas ocorrerem

Então nenhum novo incidente deve ser criado.

## Decisões

* Disponibilidade baseada apenas em HTTP.
* 2xx e 3xx são considerados sucesso.
* Threshold padrão: 2 falhas consecutivas.
* Timeout padrão: 10 segundos.
* Intervalo padrão: 1 minuto.

---

# F002 — Recuperação de Serviço (Retorno ao Ar)

## Objetivo

Detectar quando um serviço monitorado volta a operar normalmente após um período de indisponibilidade.

## Problemas Resolvidos

* P00
* P01

## Jornadas Relacionadas

* J04
* J06
* J08

## Ator Principal

Sistema

## Gatilho

Verificação bem-sucedida em monitor com incidente ativo.

## Entradas

* Resultado da verificação
* Incidente ativo

## Saídas

* Encerramento do incidente
* Notificação de recuperação

## Fluxo Principal

1. Existe um incidente ativo.
2. Sistema executa nova verificação.
3. Serviço responde corretamente.
4. Incidente é encerrado.
5. Duração total é calculada.
6. Notificação é enviada.

## Regras de Negócio

### RN-010

Somente incidentes ativos podem ser encerrados.

### RN-011

Uma única verificação bem-sucedida encerra o incidente.

### RN-012

A duração do incidente deve ser calculada automaticamente.

### RN-013

Toda recuperação gera notificação.

## Critérios de Aceitação

### CA-008

Dado um incidente ativo

Quando uma verificação for bem-sucedida

Então o incidente deve ser encerrado.

### CA-009

Dado um incidente encerrado

Então sua duração deve ser calculada.

### CA-010

Dado um incidente encerrado

Então uma notificação deve ser enviada.

# F003 — Monitoramento de Lentidão do Checkout

## Objetivo

Detectar degradação de performance no checkout antes que ela gere abandono de compra e perda de conversão.

## Problemas Resolvidos

* P00 — Descobrir problemas tarde demais
* P02 — Checkout lento

## Jornadas Relacionadas

* J03 — Monitoramento Diário
* J04 — Resposta a Incidentes
* J06 — Recebimento de Alertas
* J08 — Investigação de Problemas

## Ator Principal

Sistema

## Gatilho

Execução automática periódica do monitor de checkout.

## Entradas

### Obrigatórias

* URL do checkout

### Configuráveis

* Threshold de performance
* Timeout
* Intervalo de execução

## Saídas

### Operacionais

* Registro de medição

### Negócio

* Incidente de Performance
* Notificação

## Fluxo Principal

1. Sistema inicia execução do monitor.
2. Navegador automatizado acessa a URL.
3. Tempo de carregamento é medido.
4. Resultado é registrado.
5. Nenhuma ação adicional é executada.

## Fluxo Alternativo — Lentidão Detectada

1. Tempo excede o threshold.
2. Ocorrência é registrada.
3. Sistema verifica ocorrências consecutivas.
4. Threshold atingido.
5. Incidente é criado.
6. Notificação é enviada.

## Fluxo Alternativo — Recuperação

1. Nova medição retorna ao padrão esperado.
2. Incidente é encerrado.
3. Notificação de recuperação é enviada.

## Regras de Negócio

### RN-024

Checkout é um tipo especializado de Monitor.

### RN-025

Medições de performance devem ser armazenadas.

### RN-026

O monitor deve utilizar navegador automatizado.

### RN-027

A primeira ocorrência não abre incidente.

### RN-028

Somente um incidente ativo pode existir por checkout.

## Critérios de Aceitação

### CA-018

Dado um checkout monitorado

Quando o tempo estiver abaixo do threshold

Então a medição deve ser registrada como saudável.

### CA-019

Dado um checkout monitorado

Quando ocorrerem lentidões consecutivas suficientes

Então um incidente deve ser criado.

### CA-020

Dado um incidente ativo

Quando a medição voltar ao padrão

Então o incidente deve ser encerrado.

## Decisões

* Utilizar Playwright.
* Threshold padrão: 5 segundos.
* Threshold para incidente: 2 ocorrências consecutivas.
* Todas as medições devem ser armazenadas.
* Checkout gera Incidentes de Performance.
* Dashboard exibe última medição, média 24h e status.

---

# F005 — Monitoramento de Dependências Externas

## Objetivo

Detectar indisponibilidade em serviços externos críticos para a operação.

## Problemas Resolvidos

* P00 — Descobrir problemas tarde demais
* P03 — Dependências indisponíveis

## Jornadas Relacionadas

* J03
* J04
* J06
* J08

## Ator Principal

Sistema

## Gatilho

Execução automática periódica do monitor.

## Entradas

### Obrigatórias

* Nome da dependência
* Tipo
* URL

### Configuráveis

* Timeout
* Intervalo
* Threshold
* Severidade

## Saídas

### Operacionais

* Registro de verificação

### Negócio

* Incidente
* Notificação

## Fluxo Principal

1. Sistema verifica a dependência.
2. Dependência responde corretamente.
3. Resultado é registrado.

## Fluxo Alternativo — Falha

1. Dependência não responde.
2. Falha é registrada.
3. Threshold atingido.
4. Incidente é aberto.
5. Notificação é enviada.

## Fluxo Alternativo — Recuperação

1. Dependência volta a responder.
2. Incidente é encerrado.
3. Notificação de recuperação é enviada.

## Regras de Negócio

### RN-019

Toda dependência possui tipo.

### RN-020

Toda dependência possui severidade.

### RN-021

Apenas um incidente ativo por dependência.

### RN-022

Dependências informativas não geram alertas.

### RN-023

Toda dependência monitorada deve possuir URL válida.

## Critérios de Aceitação

### CA-014

Dado uma dependência monitorada

Quando a verificação for bem-sucedida

Então o resultado deve ser registrado.

### CA-015

Dado uma dependência monitorada

Quando o threshold for atingido

Então um incidente deve ser criado.

### CA-016

Dado um incidente ativo

Quando a dependência voltar a responder

Então o incidente deve ser encerrado.

## Decisões

* Threshold padrão: 2 falhas consecutivas.
* Severidades: Crítica, Importante e Informativa.
* V1 monitora apenas disponibilidade.
* Dependências informativas não geram alertas.

## Observação de Discovery

Nem todas as dependências possuem endpoints públicos ou mecanismos padronizados de monitoramento.

Algumas poderão exigir monitoramento indireto através de eventos operacionais.

-------

# F006 — Histórico de Quedas

## Objetivo

Permitir consulta histórica dos incidentes registrados.

## Problemas Resolvidos

* P00
* P01

## Jornadas Relacionadas

* J05
* J08
* J09

## Ator Principal

Influenciador

## Gatilho

Consulta do usuário.

## Entradas

* Incidentes registrados
* Período selecionado

## Saídas

* Histórico de incidentes
* Métricas de disponibilidade

## Regras de Negócio

### RN-015

Todo incidente encerrado permanece disponível para consulta.

### RN-016

O histórico é somente leitura.

### RN-017

Disponibilidade é calculada a partir dos incidentes registrados.

### RN-018

Incidentes ativos também aparecem no histórico.

## Critérios de Aceitação

### CA-011

Dado um incidente encerrado

Quando o usuário consultar o histórico

Então o incidente deve ser exibido.

### CA-012

Dado um período selecionado

Quando o histórico for consultado

Então apenas incidentes daquele período devem ser exibidos.

## Decisões

* Histórico somente leitura.
* Incidentes não podem ser removidos manualmente.
* Disponibilidade calculada automaticamente.

---

# F007 — Resumo Diário

## Objetivo

Fornecer uma visão consolidada da operação.

## Problemas Resolvidos

* P00

## Jornadas Relacionadas

* J09

## Ator Principal

Sistema

## Gatilho

Execução diária agendada.

## Entradas

* Incidentes do período
* Métricas de disponibilidade
* Métricas de performance

## Saídas

* Resumo operacional
* Notificação

## Fluxo Principal

1. Sistema consolida os dados do dia.
2. Métricas são calculadas.
3. Resumo é gerado.
4. Resumo é enviado pelos canais configurados.

## Regras de Negócio

### RN-030

Resumo deve ser enviado diariamente.

### RN-031

Resumo deve ser enviado mesmo sem incidentes.

### RN-032

Resumo utiliza os mesmos canais de notificação.

### RN-033

Resumo consolida dados do dia anterior.

## Critérios de Aceitação

### CA-021

Dado um dia encerrado

Quando a rotina diária executar

Então um resumo deve ser gerado.

### CA-022

Dado um resumo gerado

Quando existirem canais configurados

Então o resumo deve ser enviado.

## Decisões

* Envio diário.
* Utiliza os mesmos canais de notificação.
* Enviado mesmo sem incidentes.
* Somente leitura.

---

# Descobertas do Domínio

## Conceitos Emergentes

### Monitor

Representa algo monitorado pelo sistema.

Tipos identificados:

* Uptime Monitor
* Checkout Monitor
* Dependency Monitor

### MonitorExecution

Representa uma execução individual de um monitor.

### Incident

Representa um problema detectado.

Tipos identificados:

* Availability Incident
* Performance Incident

### Notification

Representa uma comunicação enviada ao usuário.

### Dependency

Representa um serviço externo necessário para a operação.

---

# Resumo

O contexto de Observabilidade é responsável por detectar problemas, registrar incidentes, notificar usuários e fornecer visibilidade operacional sobre a saúde da operação monitorada.

# Regras e Decisões Descobertas Durante o Refinamento

## RN-034 — Falha de Notificação Não Impede o Incidente

A falha no envio de uma notificação não pode impedir a criação, atualização ou encerramento de um incidente.

### Motivação

O incidente representa um fato ocorrido na operação.

A notificação é apenas um mecanismo de comunicação desse fato.

---

## RN-035 — Dependências Críticas Geram Alertas Imediatos

Dependências classificadas como Críticas devem gerar alertas assim que um incidente for aberto.

### Exemplos

* Gateway PIX
* Banco
* Gateway de Pagamento

---

## RN-036 — Dependências Informativas Não Geram Alertas

Dependências classificadas como Informativas devem ser monitoradas e registradas, porém não devem gerar notificações para o usuário.

### Exemplos

* Analytics
* Ferramentas auxiliares

---

# Perguntas em Aberto

## PA-023 — Retenção de Execuções

Por quanto tempo as execuções dos monitores devem ser armazenadas?

### Opções

* 30 dias
* 90 dias
* 180 dias
* 1 ano
* Indefinidamente

### Impacto

* Custo de armazenamento
* Capacidade de auditoria
* Relatórios históricos

---

## PA-024 — Retenção de Incidentes

Por quanto tempo os incidentes devem permanecer disponíveis para consulta?

### Opções

* 1 ano
* 2 anos
* Indefinidamente

### Impacto

* Histórico operacional
* Auditoria
* Relatórios de longo prazo

---

## PA-025 — Limite de Monitores por Cliente

Existe limite de monitores por cliente?

### Exemplos

Plano Básico:

* 5 monitores

Plano Profissional:

* 20 monitores

Plano Enterprise:

* Ilimitado

### Impacto

* Modelo de monetização
* Custos operacionais
* Escalabilidade da plataforma


# Regras e Decisões Descobertas Durante o Refinamento

## RN-037 — Todo Monitor Possui Status

Todo monitor deve possuir um status operacional.

### Possíveis Status

* Active
* Paused
* Disabled

### Regras

* Monitores Active participam normalmente das execuções.
* Monitores Paused não executam verificações.
* Monitores Disabled permanecem cadastrados para fins históricos, porém não executam verificações.

---

## RN-038 — Todo Incidente Possui Status

Todo incidente deve possuir um status que represente seu ciclo de vida.

### Possíveis Status

* Open
* Resolved

### Regras

* Incidentes Open representam problemas ainda ativos.
* Incidentes Resolved representam problemas encerrados.
* Somente incidentes Open podem ser encerrados.

---

## RN-039 — Toda Notificação Possui Estado de Entrega

Toda notificação deve registrar seu estado de entrega.

### Informações Mínimas

* Canal
* Conteúdo
* Data de envio
* Status de entrega

### Possíveis Status

* Pending
* Sent
* Failed

---

## PA-026 — Quem é o Proprietário de um Monitor?

Qual entidade é responsável por possuir um Monitor?

### Possibilidades

* Cliente
* Organização
* Workspace

### Observação

Esta decisão impacta diretamente o Domain Discovery e a modelagem de multi-tenancy da plataforma.

---

## PA-027 — Incidentes Possuem Severidade?

Desejamos classificar incidentes por severidade?

### Possibilidades

* Critical
* Warning
* Info

### Exemplos

Site fora do ar:

* Critical

Checkout lento:

* Warning

Dependência informativa indisponível:

* Info

---

## PA-028 — Primeira Execução do Monitor

Quando um monitor for criado, quando deve ocorrer sua primeira execução?

### Possibilidades

* Imediatamente após criação
* No próximo ciclo agendado

### Impacto

* Experiência do usuário
* Tempo até primeira validação
* Consumo de recursos

