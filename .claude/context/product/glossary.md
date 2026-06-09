# Glossário — VIGIA

## Objetivo

Definir a linguagem ubíqua global do produto VIGIA. Estes termos são usados de forma consistente
em toda a documentação de produto e de contextos.

---

## Conceitos Centrais do Sistema

Conceitos globais usados por múltiplos contextos. Definidos sem entrar em modelagem específica
de cada contexto — cada contexto refina seu lifecycle e invariantes em sua própria
domain-discovery, mas deve usar exatamente estes nomes e significados de base.

**Client**
Cliente do VIGIA — entidade contratante da plataforma. Pode possuir `Monitor`s, receber
`Notification`s e visualizar métricas. Não confundir com `Influenciador`: `Influenciador` é o
perfil de negócio do público-alvo; `Client` é o papel desse perfil dentro do sistema VIGIA.

**Monitor**
Configuração que define algo que deve ser observado continuamente — ex.: disponibilidade,
performance, PIX, banco, integração. Um `Monitor` produz observações que podem resultar em
`Incident`s.

**Incident**
Situação operacional identificada pelo sistema que requer atenção do `Client` — ex.:
indisponibilidade, lentidão, falha de PIX. Pode ser aberto e posteriormente encerrado.

**Notification**
Comunicação enviada ao `Client` sobre eventos relevantes — ex.: `Incident` aberto, recuperação,
resumo diário. Pode ser entregue por diferentes canais (WhatsApp, Telegram, e-mail, ligação).
Conceito oficial de domínio; "Alerta" é apenas seu termo de UX/comunicação (ver `Alerta` e seção
"Regras de Linguagem").

---

## Negócio e Domínio

**Influenciador**
Pessoa com audiência própria que realiza rifas online regularmente, depende de campanhas digitais
para vender e opera sobre plataformas de rifas de terceiros. Público-alvo inicial do VIGIA.

**Operação**
Conjunto de componentes interligados que sustentam a venda de uma rifa online — site, checkout,
PIX, gateway de pagamento, bancos, infraestrutura e integrações externas.

**Sorteio**
Evento de rifa realizado pelo influenciador — unidade central sobre a qual a operação acontece e
sobre a qual o VIGIA mede disponibilidade, conversão, receita e desempenho.

**Cota**
Unidade de participação em um sorteio, adquirida pelo comprador através do fluxo de compra.

**Campanha**
Período de divulgação e venda intensiva de um sorteio, frequentemente acompanhado de lives —
momento crítico em que falhas operacionais geram maior impacto.

**Receita**
Faturamento gerado pela venda de cotas em sorteios — eixo central pelo qual o VIGIA avalia o
impacto de problemas operacionais.

**Conversão**
Capacidade da operação de transformar interesse (cliques, acessos) em vendas concluídas (PIX
pago) — medida ao longo da jornada do comprador.

---

## Monitoramento e Observabilidade

**Monitoramento Operacional**
Acompanhamento contínuo da operação de vendas sob a ótica da receita e da conversão — diferencial
do VIGIA frente a ferramentas tradicionais que olham apenas para infraestrutura.

**Disponibilidade**
Condição de um componente da operação (site, checkout, infraestrutura, integrações) estar no ar e
funcional.

**Lentidão**
Degradação de desempenho de um componente da operação — em particular do checkout — que prejudica
a experiência de compra sem necessariamente derrubar o serviço.

**Queda**
Evento em que um componente da operação para de funcionar ou fica indisponível — gera histórico
operacional e dispara alertas.

**Alerta**
Termo de UX/comunicação usado para representar uma `Notification` urgente perante o
influenciador — ex.: "você recebeu um alerta de queda". Não é conceito de domínio próprio; é
sinônimo de superfície para `Notification`. Ver definição de `Notification` abaixo e seção
"Regras de Linguagem".

**Histórico Operacional**
Registro de eventos passados (quedas, falhas, recuperações) que permite ao influenciador entender
o que aconteceu e fundamentar decisões futuras.

**Resumo Diário**
Síntese periódica do estado e desempenho da operação, entregue ao influenciador como apoio à
visibilidade contínua.

---

## Pagamentos e Receita

**PIX**
Meio de pagamento instantâneo sobre o qual a operação de vendas do influenciador é baseada —
ponto crítico de monitoramento dentro da cadeia de pagamento.

**Gateway de Pagamento**
Serviço externo responsável por processar pagamentos da operação — componente da cadeia cuja
falha impacta diretamente a receita.

**Aprovação Bancária**
Processo pelo qual um banco confirma ou recusa uma transação de pagamento — fonte de falhas
invisíveis que podem causar perda de receita.

**Instabilidade Bancária**
Comportamento irregular ou degradado de bancos na cadeia de pagamento, capaz de gerar falhas
intermitentes na aprovação de transações.

**Cota Presa**
Cota cuja compra foi iniciada mas não concluída por travamento no fluxo — representa receita
potencialmente perdida.

**Estimativa de Prejuízo**
Cálculo que conecta falhas operacionais a impacto financeiro concreto, dimensionando quanto o
influenciador deixou de ganhar.

**Inteligência de Receita**
Capacidade do VIGIA de transformar dados operacionais e de conversão em estimativas de prejuízo,
recordes de desempenho, histórico e tendências — apoiando decisões e otimização de sorteios.

---

## Conversão e Analytics

**Funil de Conversão**
Sequência de etapas que um comprador percorre desde o clique inicial até a confirmação do
pagamento via PIX.

**Abandono**
Saída do comprador do funil de conversão antes de concluir a compra — analisado por etapa para
identificar pontos de perda.

**Comprador Recorrente**
Pessoa que realiza mais de uma compra ao longo de diferentes sorteios — métrica usada para
entender fidelização e retorno de público.

**Comparativo entre Sorteios**
Análise que confronta desempenho de diferentes sorteios para apoiar aprendizado e otimização de
campanhas futuras.

**Recorde**
Marco de desempenho atingido pela operação (ex.: maior conversão, maior receita em um sorteio) —
reconhecido pelo VIGIA como referência histórica.

---

## Produto e Estrutura

**VIGIA**
Plataforma de monitoramento operacional especializada no mercado de rifas online — atua como
central que dá visibilidade em tempo real sobre disponibilidade, pagamentos, conversão e receita.

**Pilar**
Grande área de capacidade do produto, que agrupa funcionalidades relacionadas por objetivo comum.
Pilares do VIGIA: Observabilidade, Revenue Protection, Revenue Intelligence.

**Observabilidade**
Pilar responsável por dar visibilidade em tempo real sobre disponibilidade e funcionamento técnico
da operação de vendas, alertando o influenciador sobre quedas e retornos.

**Revenue Protection**
Pilar responsável por monitorar os elos financeiros da cadeia de vendas (PIX, aprovação bancária,
instabilidade, cotas presas) para evitar perda de receita.

**Revenue Intelligence**
Pilar responsável por ajudar o influenciador a entender e otimizar a performance dos sorteios
através de analytics de conversão e comparações históricas.

**Funcionalidade**
Capacidade concreta do produto, identificada por um código (ex.: F001) e classificada por fase de
entrega (MVP, V1, V2).

**MVP / V1 / V2**
Fases de entrega do roadmap do produto, em ordem de prioridade — MVP cobre o pilar de
Observabilidade; V1 estende Observabilidade e introduz Revenue Protection; V2 introduz Revenue
Intelligence.

---

## Regras de Linguagem

Estas regras fixam a nomenclatura oficial. Todos os contextos (Observability, Notifications,
Revenue Protection, etc.) devem seguir exatamente estes termos — evitar sinônimos, traduções
livres ou termos técnicos concorrentes.

* Usar `Incident`, não "Problema", "Issue" ou "Falha".
* Usar `Monitor`, não "Checker", "Probe" ou "Sensor".
* Usar `Notification` como conceito oficial de comunicação ao cliente — não criar conceito
  paralelo de domínio para "Alerta", "Aviso" ou "Mensagem".
* Usar `Client` para o cliente do VIGIA — não confundir com `Influenciador` (perfil de negócio
  do público-alvo) nem usar "Usuário", "Conta" ou "Tenant" como sinônimo de domínio.
* Usar `Influenciador` para o usuário final do negócio (dono da operação de rifas), reservando
  `Client` para o papel desse usuário dentro do sistema VIGIA.
* Usar "Alerta" apenas como termo de UX/comunicação para representar uma `Notification` urgente
  — nunca como conceito de domínio com modelagem própria.
* Usar `Sorteio` para o evento de rifa — não usar "Campanha" como sinônimo (`Campanha` é o
  período de divulgação/venda em torno de um `Sorteio`, não o evento em si).
* Usar `Cota` para a unidade de participação adquirida — não usar "Bilhete", "Ticket" ou
  "Número".
* Usar `Queda` para indisponibilidade observada de um componente, e `Incident` para a situação
  operacional aberta a partir dessa observação — os dois não são sinônimos: uma `Queda` pode
  gerar um `Incident`, mas nem toda observação vira `Incident`.
* Manter `MVP`, `V1`, `V2` como identificadores de fase de roadmap — não renomear para "Fase 1",
  "Fase 2" etc. em nenhum artefato.
