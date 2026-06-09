# Overview — Notification

## Problema

O influenciador não está sempre dentro do app. Quando um Incident abre ou resolve, o sistema tem a informação mas não tem como avisar o influenciador fora da interface. O impacto aparece nas vendas antes de ele saber que havia um problema.

---

## Objetivo

Entregar ao influenciador uma mensagem no WhatsApp no momento em que um Incident é aberto ou resolvido — independentemente de ele estar com o app aberto.

---

## Usuários

- **Influenciador** — destinatário das notificações; número de WhatsApp configurado no perfil.

---

## Processos

### Notificar Incident Aberto
Incident aberto no contexto Observability → `Notification` criada e entregue via WhatsApp.

- **Pré-condição:** evento `incident_opened` recebido; número de WhatsApp disponível no perfil.
- **Pós-condição:** `Notification` com status `delivered` ou `dead` após retries esgotados.
- **Side effects:** `Notification` persistida com histórico de tentativas.

### Notificar Incident Resolvido
Incident resolvido → influenciador recebe confirmação com duração do problema.

- **Pré-condição:** evento `incident_resolved` recebido; número de WhatsApp disponível.
- **Pós-condição:** `Notification` com status `delivered` ou `dead`.

---

## Regras de Negócio

- Canal único: WhatsApp.
- 1 destinatário por conta: número configurado no perfil do influenciador.
- Todos os eventos notificados por padrão — sem configuração de preferência no MVP.
- Templates fixos por tipo de evento.
- Falha de entrega: até 3 tentativas com backoff crescente. Esgotado → `dead`.
- `Notification` com status `dead` é auditável mas não reenviada manualmente no MVP.

---

## Dependências

- **Observability** — fonte dos eventos `incident_opened` e `incident_resolved`.
- **Profile** — fornece `WhatsAppNumber` do destinatário. Ownership de Profile não modelado ainda — tratado como dependência externa no MVP (PA-N02).
- **Provedor WhatsApp** — infraestrutura de entrega. Provedor específico a definir (PA-N01).

---

## Fora de Escopo (MVP)

- Notificação de Daily Summary → backlog.
- Configuração de preferências de notificação.
- Múltiplos destinatários.
- Múltiplos canais.
- Silenciar notificações.
- Escalação por tempo sem resolução.

---

## Perguntas em Aberto

**Bloqueiam o Discovery:**
- nenhuma.

**Backlog:**
- PA-N01 — Provedor WhatsApp (Twilio, Z-API, Meta API).
- PA-N02 — Ownership de Profile (onde vive `WhatsAppNumber`).
- PA-N03 — Múltiplos destinatários.
- PA-N04 — Preferências de notificação por canal/evento.
- PA-N05 — Silenciar por período.
- PA-N06 — Múltiplos canais.
- PA-N07 — Escalação.
- PA-N08 — Templates configuráveis.
