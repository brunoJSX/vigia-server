# Domain Discovery — Observability

## Objetivo

Descrever a compreensão atual dos conceitos centrais do contexto Observability — sem assumir
Entity, Aggregate, Value Object ou Domain Event, e sem descrever comportamento, ordem de
execução ou fluxo temporal (isso pertence ao `workflow.md`). Este documento responde apenas: o
que existe no domínio, quais conceitos o negócio reconhece, quais são centrais, ambíguos,
mecanismos internos ou dados operacionais — e como esses conceitos se relacionam entre si de
forma estática.

## Conceitos Centrais do Domínio

### Monitor

Representa a configuração de uma verificação contínua — define o que deve ser observado
(disponibilidade, performance de checkout, dependências externas), com que frequência e threshold,
e mantém status operacional (Active, Paused, Disabled — RN-037). Monitor não executa a
verificação; é a configuração que o pipeline usa para executá-la. Reconhecido pelo glossário e
pelo overview como conceito central; é o termo que o influenciador usa para falar do que está
sendo vigiado.

Monitor do tipo Checkout carrega adicionalmente `AcceptableResponseTime` — limiar de tempo com
significado de domínio: define o ponto a partir do qual o fluxo de checkout é considerado lento
(Lentidão, glossário). O tempo de espera do Collector para Uptime e Dependency é detalhe técnico
interno, sem atributo correspondente no domínio.

Invariantes relevantes: todo Monitor tem um status (RN-037); Checkout é um tipo especializado de
Monitor (RN-024 — ver Conceitos Ambíguos para os demais tipos); para Checkout,
`AcceptableResponseTime` é obrigatório (RN-025).

### Incident

Representa uma situação operacional relevante para o cliente — indisponibilidade, lentidão,
falha de dependência. Possui lifecycle conceitual com dois estados: Open e Resolved (RN-038); um
Incident em estado Resolved não retorna a Open. Reconhecido pelo glossário e pelo overview como
conceito central; é o termo que o influenciador usa para falar de quedas, problemas e
recuperações.

Invariantes relevantes: existe no máximo um Incident em estado Open por Monitor (RN-002, RN-028);
um Incident carrega consigo a duração do problema que representa (RN-012).

## Conceitos Ambíguos (Ainda em Descoberta)

### Dependency

Aparece nas fontes com atributos e regras próprias — tipo, severidade, URL, e a distinção entre
dependências críticas e informativas (RN-019 a RN-023, RN-035, RN-036). Tem peso de domínio
suficiente para não ser ignorado, mas não está claro se é um conceito de primeira classe
(registrado e gerenciado com identidade própria) ou apenas a configuração de um "Dependency
Monitor" — uma variação de `Monitor`, não um conceito paralelo.

### DailySummary (Resumo Diário)

F007 consolida incidentes e métricas de um período em uma visão recorrente enviada ao cliente.
Nenhuma fonte o nomeia como conceito — pode ser um conceito de suporte próprio ou simplesmente
uma agregação/consulta sobre `Incident` e dados operacionais já existentes, sem identidade
própria.

## Mecanismos Internos

Papéis técnicos que participam de como um `Monitor` é verificado e de como se chega a uma
decisão sobre `Incident` — não são conceitos de domínio: nenhuma fonte os nomeia, nenhum
influenciador os reconheceria, e não possuem identidade, lifecycle ou regra de negócio próprios.
As fontes descrevem esses passos de forma genérica ("Sistema executa verificação", "navegador
automatizado acessa a URL"). Os nomes abaixo são rótulos provisórios, úteis para conversas
técnicas, e podem ser renomeados ou reorganizados livremente em Workflow ou documentação técnica
sem impacto no domínio.

- **Collector** — papel de obter o dado bruto de uma verificação (ex.: resposta HTTP, tempo de
  carregamento, status de uma dependência). Considera sucesso apenas respostas com status HTTP
  2xx ou 3xx; respostas 4xx e 5xx são registradas como falha (`Sample.Success = false`). O
  timeout interno é derivado automaticamente a partir da configuração do Monitor: para Checkout,
  a partir de `AcceptableResponseTime` (com fator multiplicador técnico interno); para Uptime e
  Dependency, a partir de um valor padrão interno. Esse timeout não é conceito de domínio e não
  é persistido como atributo de `Monitor`.
- **Analyzer** — papel de interpretar dados coletados contra threshold e regras de
  consecutividade, produzindo uma decisão sobre o estado de um `Incident`. O Analyzer é
  type-aware: para Checkout, avalia latência do Sample contra `AcceptableResponseTime` do
  Monitor (detecta Lentidão — RN-025); para Uptime e Dependency, avalia alcançabilidade
  (`Sample.Success`).

## Dados Operacionais

Resultado bruto de cada verificação realizada por um `Monitor` (ex.: dados coletados, samples,
resultados de verificação, métricas). Têm relação de suporte com `Monitor` — são produzidos a
partir dele e usados como evidência para o cálculo de disponibilidade, o histórico consultável
(F006) e a fundamentação de decisões sobre `Incident` (RN-003, RN-017, RN-025) — mas, isolados,
não têm identidade, lifecycle ou regra de negócio própria: são evidência, não conceitos da
linguagem ubíqua. O domínio permanece completo descrevendo-os apenas como dado de suporte de
`Monitor`.

## Relacionamentos Conceituais

### Monitor — Incident
Cada `Monitor` admite, no máximo, um `Incident` em estado Open simultaneamente (RN-002, RN-028) —
um vínculo de cardinalidade/identidade entre os dois conceitos, não uma descrição de quando ou
como esse vínculo se forma (isso é comportamento, descrito em `workflow.md`).

### Monitor — Dados Operacionais
Dados operacionais existem em função de um `Monitor` — são a evidência de suporte que o
caracteriza ao longo do tempo (RN-003, RN-025). A relação é de origem/suporte, não de processo.

### Incident — Comunicação ao Cliente
Um `Incident` é o fato de domínio que justifica a existência de uma comunicação ao cliente. A
natureza exata dessa relação — quem a possui, como se expressa — depende da definição de
ownership de `Notification`, ainda em aberto (ver Perguntas em Aberto).

## Perguntas em Aberto

1. Uptime, Checkout e Dependency são especializações de `Monitor`, ou conceitos próprios com
   atributos compartilhados? **Decisão atual**: Monitor permanece único (`Type` como valor de
   configuração); especialização formal adiada — não há evidência de cluster suficiente de
   atributos/invariantes exclusivos por tipo além de `AcceptableResponseTime` em Checkout.
   Revisitar se novos atributos tipo-exclusivos surgirem.
2. `Dependency` é conceito de primeira classe ou configuração de um Monitor? (ver Conceitos
   Ambíguos) — questão pendente: "o que muda no comportamento de `Incident` ou `Notification`
   quando Dependency é Critical vs Informational?" Não modelar Severity até resposta clara.
3. `DailySummary` é conceito próprio ou apenas agregação/consulta sobre `Incident` e dados
   operacionais? (ver Conceitos Ambíguos)
4. Por quanto tempo os dados operacionais (verificações/medições) devem ser retidos? (PA-023 do
   documento-fonte)
5. Quem é o dono de `Notification` — este contexto, um contexto dedicado de comunicações, ou
   infraestrutura compartilhada? (questão já registrada no `overview.md`, ainda sem resposta)
6. ~~Semântica de `Success` no `Collector`~~ — **Fechada**: sucesso = 2xx ou 3xx; 4xx e 5xx =
   falha (`Sample.Success = false`). Decisão deliberada para MVP; revisitar se surgir necessidade
   de configurar expected status codes por Monitor.
