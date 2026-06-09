# Workflow — Account

## GetAccount

Retorna as preferências da account autenticada. Cria a account se ainda não existir (não aplicável quando trigger está ativo — mantido como fallback).

```mermaid
flowchart TD
    A[Request GET /account\nAuthorization: Bearer JWT] --> B[Middleware valida JWT\nSupabase JWT Secret]
    B --> C{JWT válido?}
    C -- Não --> D[401 Unauthorized]
    C -- Sim --> E[Extrai user_id do JWT]
    E --> F[Busca Account pelo user_id]
    F --> G{Account existe?}
    G -- Não --> H[Cria Account\nwhatsapp_number: null]
    G -- Sim --> I[Retorna Account]
    H --> I
```

---

## UpdateAccount

Influenciador atualiza o número de WhatsApp.

```mermaid
flowchart TD
    A[Request PATCH /account\nAuthorization: Bearer JWT\nbody: whatsapp_number] --> B[Middleware valida JWT]
    B --> C{JWT válido?}
    C -- Não --> D[401 Unauthorized]
    C -- Sim --> E[Extrai user_id do JWT]
    E --> F{whatsapp_number válido?\nE.164 ou null}
    F -- Não --> G[400 Bad Request]
    F -- Sim --> H[Atualiza Account\nwhatsapp_number + updated_at]
    H --> I[Retorna Account atualizada]
```

---

## ResolveRecipient (interno — contexto Notification)

Quando Notification precisa do número de WhatsApp para criar uma `Notification`, carrega a Account.

```mermaid
flowchart TD
    A[Evento recebido\nincident_opened | incident_resolved] --> B[Carrega Account]
    B --> C{WhatsAppNumber disponível?}
    C -- Não --> D[Ignorar — RN-AC004\nsem destinatário, sem Notification]
    C -- Sim --> E[Cria Notification com recipient\n= Account.WhatsAppNumber]
```

---

## Regras aplicadas nos workflows

| Regra | Onde se aplica |
|---|---|
| JWT obrigatório em todo request | GetAccount, UpdateAccount |
| Account criada no trigger `on_auth_user_created` | Criação automática |
| Fallback lazy: cria se não existir no `GET /account` | GetAccount |
| `WhatsAppNumber` em E.164 ou null | UpdateAccount |
| Account sem número → notificações ignoradas | ResolveRecipient |
| `updated_at` atualizado a cada `PATCH` | UpdateAccount |
