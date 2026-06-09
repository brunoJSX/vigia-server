# Domain Discovery — Account

## Conceitos Centrais

### Account

Representa o influenciador como entidade do domínio — o dono da instalação, dos monitors e das preferências de comunicação. Não é a identidade (credenciais, sessão) — isso pertence ao Supabase Auth. `Account` é o que o domínio do Vigia conhece sobre quem usa o produto.

**Responsabilidade:** armazenar e fornecer as preferências do influenciador — no MVP, o número de WhatsApp para notificações.

**Ownership:** contexto Account.

**Invariantes:**
- `Account.ID` é o `user_id` do Supabase Auth — vínculo 1:1 imutável.
- `WhatsAppNumber`, quando presente, está em formato E.164 (`+` + código do país + número).
- `WhatsAppNumber` é nullable — account pode existir sem número configurado.
- Uma `Account` é criada no primeiro acesso autenticado; nunca duplicada para o mesmo `user_id`.

**Atributos:**

| Atributo | Tipo | Descrição |
|---|---|---|
| `ID` | string (UUID) | `user_id` do Supabase Auth |
| `WhatsAppNumber` | string (nullable) | Número de WhatsApp em E.164 |
| `CreatedAt` | timestamp | Quando a account foi criada |
| `UpdatedAt` | timestamp | Última atualização de preferências |

---

### JWT (AuthToken)

Token emitido pelo Supabase Auth após login. Não é conceito de domínio — é infraestrutura de transporte. O Vigia valida o JWT em todo request e extrai o `user_id` para identificar a `Account`. O domínio nunca manipula o token diretamente.

**Responsabilidade:** transportar a identidade do usuário autenticado até o middleware do Vigia.

**Ownership:** Supabase Auth — fora do domínio do Vigia.

---

## Fluxo de Dados

```
Influenciador faz login no frontend (Supabase Auth)
        ↓
Supabase Auth emite JWT com user_id
        ↓
Frontend envia request: Authorization: Bearer <JWT>
        ↓
Middleware valida JWT (SUPABASE_JWT_SECRET)
        ↓
Extrai user_id → injeta no contexto do request
        ↓
Use case carrega Account pelo user_id
        ↓
Processa e responde
```

---

## Responsabilidades

### Account

Responsável por:
- Representar o influenciador como entidade do domínio
- Armazenar preferências configuráveis (WhatsAppNumber)
- Ser a fonte de verdade do número de WhatsApp para o contexto Notification

Não responsável por:
- Autenticação (responsabilidade do Supabase Auth)
- Validação de sessão / expiração de token
- Gerenciar credenciais

### Middleware de Autenticação

Responsável por:
- Validar o JWT em todo request
- Extrair e injetar o `user_id` no contexto do request
- Rejeitar requests sem JWT válido (401)

Não responsável por:
- Lógica de negócio
- Carregar a `Account`

---

## Relacionamentos

### Supabase Auth → Account
`Account.ID` = `user_id` do Supabase Auth. O vínculo é criado no momento do primeiro acesso — o domínio não cria usuários no Supabase, apenas reage ao primeiro request autenticado.

### Account → Notification
Notification consome `WhatsAppNumber` da Account para determinar o destinatário. Atualmente o número vem de variável de ambiente (`WHATSAPP_RECIPIENT`); após este contexto, o contexto Notification carrega a Account para resolver o número. O acoplamento é unidirecional: Notification depende de Account, Account não conhece Notification.

---

## Lifecycle

`Account` não possui lifecycle com estados relevantes no MVP — nasce na criação e suas preferências são atualizadas. Sem estados active/inactive/suspended.

Criação: lazy no primeiro `GET /account` autenticado — se não existir account para o `user_id`, cria com `WhatsAppNumber = null` e retorna.

---

## Perguntas em Aberto

- Nenhuma bloqueia o Workflow.
- PA-AC01 — Multi-usuário: quando surgir, `Account` pode evoluir para `Organization` com múltiplos membros.
- PA-AC02 — `account_id` como FK em `monitors`: necessário para multi-tenancy; desnecessário no MVP single-tenant.
