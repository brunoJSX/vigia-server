# Overview — Account

## Problema

O sistema não tem dono. Monitors, Incidents e Notifications existem sem vínculo com quem usa o produto. O número de WhatsApp para notificações está fixo em variável de ambiente — o influenciador não consegue alterá-lo sem acesso ao servidor. Não há proteção de acesso: qualquer pessoa com a URL da API consegue criar, pausar ou deletar monitors.

---

## Objetivo

Estabelecer o influenciador como entidade gerenciável no sistema — com autenticação e preferências configuráveis via UI. No MVP: proteger todos os endpoints da API com autenticação e permitir que o influenciador configure seu número de WhatsApp sem tocar no servidor.

---

## Usuários

- **Influenciador** — único usuário no MVP. Contrata, configura e usa o produto.

---

## Processos

### Autenticar
O influenciador faz login via Supabase Auth (email/senha ou Google). O Supabase devolve um JWT que é enviado em todo request à API do Vigia.

- **Pré-condição:** usuário registrado no Supabase Auth.
- **Pós-condição:** JWT válido disponível para o frontend.
- **Side effects:** nenhum no domínio do Vigia — autenticação é responsabilidade do Supabase.

### Consultar Account
Retorna as preferências da account autenticada.

- **Pré-condição:** request autenticado com JWT válido; account existe.
- **Pós-condição:** nenhuma alteração de estado.

### Atualizar Account
Influenciador altera o número de WhatsApp configurado.

- **Pré-condição:** request autenticado; número em formato E.164 (ex: +5511999999999).
- **Pós-condição:** account atualizada; próximas notificações usarão o novo número.
- **Side effects:** nenhum retroativo — notificações já criadas mantêm o destinatário original.

---

## Regras de Negócio

- **RN-AC001** — Autenticação delegada ao Supabase Auth. O Vigia valida o JWT em todo request; não gerencia credenciais.
- **RN-AC002** — Account é criada automaticamente no primeiro acesso autenticado (trigger no Supabase ou criação lazy no `GET /account`).
- **RN-AC003** — Uma account possui no máximo um número de WhatsApp no MVP.
- **RN-AC004** — Account sem número de WhatsApp configurado: notificações não são enviadas (sem erro — comportamento silencioso, RN-N006).
- **RN-AC005** — Número de WhatsApp deve estar em formato E.164 (`+` + código do país + número).
- **RN-AC006** — No MVP, single-tenant: uma instância do Vigia serve exatamente um influenciador.

---

## Dependências

- **Supabase Auth** — fornece `user_id` (UUID) e JWT. O `user_id` é a chave de vínculo entre a identidade e a `Account` no domínio.
- **Notification** — consome `WhatsAppNumber` da Account para determinar o destinatário. Atualmente lê de variável de ambiente; após este contexto, lê da Account.

---

## Fora de Escopo (MVP)

- Multi-usuário por conta.
- Multi-tenancy.
- Permissões e roles.
- Outros canais de notificação (email, SMS).
- Preferências por tipo de evento (silenciar incident_resolved, etc.).
- Exclusão de conta.

---

## Perguntas em Aberto

**Bloqueiam o Discovery:**
- nenhuma.

**Backlog:**
- PA-AC01 — Multi-usuário: quando um time de ops precisar de múltiplos logins.
- PA-AC02 — `account_id` como FK em `monitors` (necessário para multi-tenancy).
- PA-AC03 — Preferências de notificação por canal e por tipo de evento.
- PA-AC04 — Exclusão de conta e seus dados.
