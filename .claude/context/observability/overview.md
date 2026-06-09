# Overview — Observability

## Objetivo

Dar visibilidade em tempo real sobre a disponibilidade e o funcionamento técnico da operação de
vendas (site, checkout, infraestrutura, integrações), detectando problemas operacionais,
registrando incidentes e notificando o influenciador assim que algo sai do ar, fica lento ou
volta a funcionar — antes que o impacto apareça como queda nas vendas.

## Problemas Resolvidos

- Detecção tardia de quedas e lentidão — hoje o influenciador só percebe o problema pela queda
  nas vendas, quando já perdeu oportunidade de venda, engajamento e faturamento.
- Falta de visibilidade sobre a saúde técnica da operação (site, checkout, infraestrutura,
  integrações externas).
- Ausência de histórico operacional que explique o que aconteceu e por quanto tempo durou.
- Falta de uma visão consolidada e recorrente do estado da operação.

## Usuários Impactados

- **Influenciador** — usuário final que opera o sorteio; é informado sobre incidentes e
  recuperações, consulta histórico de quedas e resumo diário.

## Responsabilidades

- Monitorar continuamente a disponibilidade de sites/serviços (Uptime).
- Detectar quando um serviço previamente indisponível volta a operar normalmente (Recuperação).
- Monitorar performance/lentidão do checkout.
- Monitorar dependências externas críticas para a operação.
- Identificar situações operacionais que requerem atenção e abrir/encerrar Incidents conforme
  thresholds e regras de consecutividade definidas.
- Emitir eventos relevantes quando incidentes são detectados, atualizados ou resolvidos.
- Disponibilizar informações necessárias para alertas e comunicações ao cliente (sem definir
  como/por onde essas comunicações são entregues).
- Manter histórico consultável (somente leitura) de incidentes e métricas de disponibilidade.
- Consolidar dados da operação em um resumo diário, mesmo na ausência de incidentes.

## Fora de Escopo

- Monitoramento da camada financeira/pagamentos (PIX, aprovação bancária, instabilidade
  bancária, cotas presas, estimativa de prejuízo) — pertence ao pilar Revenue Protection.
- Analytics de conversão, abandono por etapa, comparação entre sorteios, recordes de
  desempenho — pertence ao pilar Revenue Intelligence.
- Tudo que está fora de escopo do produto VIGIA (plataforma de rifas, gateway de pagamento,
  checkout, CRM, automação de marketing, criação de páginas).
- Definição, envio e entrega de notificações/comunicações ao cliente — ownership de
  `Notification` ainda não está decidido; este contexto apenas emite os eventos e disponibiliza
  as informações necessárias (ver Ambiguidades).
- Modelagem de entidades, banco de dados, APIs ou qualquer detalhe de implementação.

## Capacidades

- Monitoramento de Uptime
- Recuperação / Retorno ao Ar
- Monitoramento de Lentidão de Checkout
- Monitoramento de Dependências Externas
- Histórico de Incidentes (Quedas)
- Resumo Diário Operacional

## Integração com Outros Contextos

- Observability é o pilar fundacional do produto: **Revenue Protection depende dele** — estende
  a mesma vigilância para a camada de pagamento (PIX, bancos, cotas presas).
- **Revenue Intelligence depende de Observability** (e de Revenue Protection) — usa o histórico
  operacional de disponibilidade e falhas coletado aqui como base para analytics, comparações e
  tendências.
- Observability é a primeira expressão concreta da capacidade de "alertas em tempo real"
  descrita na Visão do produto (WhatsApp, Telegram, e-mail, ligação) — fornecendo os eventos e
  informações que viabilizam essas comunicações.

## Funcionalidades Relacionadas

MVP (cobertas por este contexto):
- F001 — Monitoramento de Uptime
- F002 — Alerta de Retorno ao Ar
- F003 — Monitoramento de Lentidão do Checkout
- F006 — Histórico de Quedas
- F007 — Resumo Diário

V1 (extensão futura do pilar Observabilidade):
- F004 — Monitoramento de Infraestrutura AWS
- F005 — Monitoramento de Integrações
