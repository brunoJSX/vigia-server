# Pilares do Produto

## Objetivo

Explicar como o produto está dividido em grandes áreas de capacidade.

## Pilares MVP

### Observabilidade

**Objetivo**: dar visibilidade em tempo real sobre a disponibilidade e o funcionamento técnico
da operação de vendas (site, checkout, infraestrutura, integrações), alertando o influenciador
assim que algo sai do ar ou volta a funcionar.

**Problemas que resolve**: detecção tardia de quedas e lentidão, descoberta do problema só pela
queda nas vendas, falta de histórico operacional para entender o que aconteceu.

**Funcionalidades relacionadas**:
- F001 Monitoramento de Uptime (MVP)
- F002 Alerta de Retorno ao Ar (MVP)
- F003 Monitoramento de Lentidão do Checkout (MVP)
- F006 Histórico de Quedas (MVP)
- F007 Resumo Diário (MVP)
- F004 Monitoramento de Infraestrutura AWS (V1)
- F005 Monitoramento de Integrações (V1)

**Capacidades relacionadas** (da Visão): monitoramento técnico e operacional; alertas em tempo
real (WhatsApp, Telegram, e-mail, ligação) — F002 é a primeira expressão concreta dessa
capacidade de alertas, embutida neste pilar.

## Pilares Futuros

### Revenue Protection

**Objetivo**: monitorar os elos financeiros da cadeia de vendas — geração de PIX, aprovação
bancária, instabilidades de pagamento e cotas presas — para que o influenciador detecte e reaja
antes de perder receita.

**Problemas que resolve**: perda de receita por falhas invisíveis na etapa de pagamento (PIX
não gerado, banco recusando aprovações, cotas travadas), e falta de dimensão do prejuízo
causado por esses problemas.

**Funcionalidades relacionadas**:
- F008 Falhas na Geração de PIX (V1)
- F009 Aprovação por Banco (V1)
- F010 Instabilidade Bancária (V1)
- F011 Cotas Presas (V1)
- F012 Estimativa de Prejuízo (V1)

**Capacidades relacionadas** (da Visão): inteligência de receita — em particular a estimativa
de prejuízo, que conecta falhas operacionais a impacto financeiro concreto.

### Revenue Intelligence

**Objetivo**: ir além da detecção de falhas e ajudar o influenciador a entender e otimizar a
performance dos sorteios através de analytics de conversão e comparações históricas.

**Problemas que resolve**: falta de visibilidade sobre onde e por que os compradores abandonam
o funil, dificuldade em comparar sorteios e aprender com o histórico, ausência de reconhecimento
de marcos de performance.

**Funcionalidades relacionadas**:
- F013 Conversão Clique → PIX (V2)
- F014 Abandono por Etapa (V2)
- F015 Comparativo entre Sorteios (V2)
- F016 Retorno de Compradores (V2)
- F017 Novo Recorde (V2)

**Capacidades relacionadas** (da Visão): analytics de conversão (do clique ao PIX, abandono por
etapa, compradores recorrentes, comparação entre sorteios); inteligência de receita (recordes de
desempenho, histórico operacional e tendências de conversão).

## Roadmap

### MVP
- **Observabilidade**: F001, F002, F003, F006, F007

### V1
- **Observabilidade** (extensão): F004, F005
- **Revenue Protection**: F008, F009, F010, F011, F012

### V2
- **Revenue Intelligence**: F013, F014, F015, F016, F017

## Dependências Entre Pilares

- **Revenue Protection depende de Observabilidade**: a detecção de falhas financeiras (PIX,
  aprovação bancária, cotas presas) pressupõe que a infraestrutura de monitoramento e alertas já
  exista — Revenue Protection estende a vigilância do pilar de Observabilidade para a camada de
  pagamento.
- **Revenue Intelligence depende de Observabilidade e Revenue Protection**: analytics de
  conversão e estimativas de tendência precisam de dados históricos de disponibilidade, falhas e
  perdas já coletados pelos dois pilares anteriores — sem esse histórico operacional, não há base
  para comparações e inteligência de receita.
