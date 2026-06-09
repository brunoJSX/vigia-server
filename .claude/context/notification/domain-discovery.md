# Domain Discovery — Notification

## Conceitos Centrais

### Notification

Representa a entrega de uma comunicação relevante ao influenciador. Possui lifecycle próprio — nasce como `pending`, progride para `delivered` ou, após retries esgotados, para `dead`. Não é o evento que a originou; é a resposta deste contexto a ele.

**Responsabilidade:** registrar a intenção de entrega, rastrear tentativas e consolidar o resultado final.

**Ownership:** contexto Notification.

**Invariantes:**
- Toda `Notification` nasce com status `pending`.
- `delivered` e `dead` são estados finais — não retrocedem.
- `Attempts` só cresce; nunca é decrementado.
- Máximo de 3 tentativas antes de transitar para `dead`.
- `DeliveredAt` só é preenchido quando status é `delivered`.

**Atributos mínimos:**

| Atributo | Descrição |
|---|---|
| `ID` | Identidade da notificação |
| `Type` | `incident_opened` ou `incident_resolved` |
| `Recipient` | Número de WhatsApp do destinatário |
| `Payload` | Dados do evento para renderizar o template |
| `Status` | `pending`, `delivered`, `failed`, `dead` |
| `Attempts` | Contador de tentativas de entrega |
| `CreatedAt` | Quando foi enfileirada |
| `DeliveredAt` | Quando foi entregue (nullable) |

---

### NotificationTemplate

Regra de transformação de `Payload` em mensagem legível por tipo de evento. Não é entidade persistida — é lógica do domínio.

**Responsabilidade:** produzir o texto da mensagem a partir do tipo e do payload.

**Ownership:** contexto Notification.

**Templates MVP:**

| Type | Mensagem |
|---|---|
| `incident_opened` | `"🔴 Problema detectado — {MonitorName} está fora do ar."` |
| `incident_resolved` | `"✅ Problema resolvido — {MonitorName} voltou ao normal. Duração: {Duration}."` |

---

### WhatsAppProvider

Abstração da infraestrutura de entrega. Não é conceito de domínio — é uma porta (interface) que isola o contexto do provedor concreto (Twilio, Z-API, Meta API). O domínio só conhece `Send(recipient, message) error`.

**Responsabilidade:** entregar a mensagem no canal. Retorna sucesso ou erro.

**Ownership:** infraestrutura — implementação fora do domínio.

---

## Fluxo de Dados

```
Observability emite evento (incident_opened | incident_resolved)
        ↓
Notification recebe evento → cria Notification{status: pending}
        ↓
Worker processa Notification pendente
        ↓
Renderiza template → chama WhatsAppProvider.Send()
        ↓
Sucesso → status: delivered, DeliveredAt preenchido
Falha   → Attempts++ → status: failed
        ↓ (se Attempts == 3)
status: dead
```

---

## Responsabilidades

### Notification

Responsável por:
- Representar uma entrega com seu lifecycle completo
- Rastrear tentativas e resultado final
- Ser auditável (status `dead` visível)

Não responsável por:
- Renderizar a mensagem (responsabilidade do template)
- Entregar a mensagem (responsabilidade do provider)
- Conhecer o conteúdo do evento de origem

### NotificationTemplate

Responsável por:
- Transformar tipo + payload em texto legível

Não responsável por:
- Persistência
- Entrega
- Lifecycle da Notification

### WhatsAppProvider

Responsável por:
- Fazer a chamada HTTP ao provedor externo
- Retornar sucesso ou erro para o domínio decidir o próximo estado

Não responsável por:
- Retry (responsabilidade do Worker)
- Persistência
- Regras de negócio

---

## Relacionamentos

### Observability → Notification
Observability emite eventos via `NotificationPublisher` (porta já existente no contexto). Notification consome esses eventos e cria `Notification` pendentes. Os contextos são desacoplados — Observability não conhece Notification.

### Profile → Notification
Notification recebe `WhatsAppNumber` como dado de entrada no momento da criação — não acessa Profile diretamente. Quem orquestra a criação (use case) é responsável por resolver o número antes de criar a `Notification`.

---

## Lifecycle

```
pending
   ↓ entrega com sucesso
delivered  (final)

pending
   ↓ falha de entrega
failed
   ↓ nova tentativa (até 3x)
failed → failed → ...
   ↓ Attempts == 3
dead  (final)
```

`failed` é estado intermediário — indica que houve falha mas ainda há tentativas disponíveis. Só transita para `dead` quando `Attempts` atinge o limite.

---

## Perguntas em Aberto

- Nenhuma bloqueia o Workflow. PA-N01 e PA-N02 são resolvidas na implementação.
