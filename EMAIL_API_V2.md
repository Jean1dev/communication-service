# Email API v2

API de envio de emails com templates dinâmicos armazenados no banco de dados.
Ao contrário da v1, nenhum redeploy é necessário para adicionar suporte a um novo sistema — basta registrar o template via API e começar a enviar.

---

## Base URL

```
https://<host>/v2/email
```

---

## Fluxo de integração

```
1. Registrar template  →  POST /v2/email/templates
2. Enviar emails       →  POST /v2/email/send
```

---

## Endpoints

### 1. Registrar template

```
POST /v2/email/templates
Content-Type: application/json
```

**Body**

| Campo          | Tipo   | Obrigatório | Descrição                                                     |
|----------------|--------|-------------|---------------------------------------------------------------|
| `slug`         | string | sim         | Identificador único do template (ex: `"confirmacao-cadastro"`) |
| `name`         | string | não         | Nome legível do template                                      |
| `fromEmail`    | string | não         | Endereço remetente. Padrão: `info@jeanconsultoria.com`        |
| `htmlTemplate` | string | sim         | HTML completo do email com variáveis no formato `{{.nomeVar}}` |

**Exemplo**

```json
{
  "slug": "confirmacao-cadastro",
  "name": "Confirmação de Cadastro",
  "fromEmail": "noreply@meuapp.com",
  "htmlTemplate": "<html><body><h1>Olá, {{.username}}!</h1><p>{{.mensagem}}</p><a href='{{.link}}'>Confirmar</a></body></html>"
}
```

**Resposta — 201 Created**

```json
{
  "status": "created",
  "message": "template created successfully",
  "slug": "confirmacao-cadastro"
}
```

---

### 2. Enviar email

```
POST /v2/email/send
Content-Type: application/json
```

**Body**

| Campo          | Tipo              | Obrigatório | Descrição                                             |
|----------------|-------------------|-------------|-------------------------------------------------------|
| `templateSlug` | string            | sim         | Slug do template cadastrado                           |
| `to`           | string            | sim         | Endereço de email do destinatário                     |
| `subject`      | string            | sim         | Assunto do email                                      |
| `variables`    | map[string]string | não         | Valores para substituir as variáveis do template HTML |

**Exemplo**

```json
{
  "templateSlug": "confirmacao-cadastro",
  "to": "usuario@email.com",
  "subject": "Confirme seu cadastro",
  "variables": {
    "username": "Maria",
    "mensagem": "Seja bem-vinda à plataforma!",
    "link": "https://meuapp.com/confirmar?token=abc123"
  }
}
```

**Resposta — 200 OK**

```json
{
  "status": "pending",
  "message": "email queued successfully"
}
```

O envio é assíncrono — a resposta confirma que o email foi enfileirado, não que foi entregue.

---

### 3. Listar templates

```
GET /v2/email/templates
```

**Resposta — 200 OK**

```json
[
  {
    "slug": "confirmacao-cadastro",
    "name": "Confirmação de Cadastro",
    "fromEmail": "noreply@meuapp.com",
    "htmlTemplate": "...",
    "createdAt": "2026-04-20T10:00:00Z",
    "updatedAt": "2026-04-20T10:00:00Z"
  }
]
```

---

### 4. Buscar template por slug

```
GET /v2/email/templates/{slug}
```

**Resposta — 200 OK**

```json
{
  "slug": "confirmacao-cadastro",
  "name": "Confirmação de Cadastro",
  "fromEmail": "noreply@meuapp.com",
  "htmlTemplate": "...",
  "createdAt": "2026-04-20T10:00:00Z",
  "updatedAt": "2026-04-20T10:00:00Z"
}
```

**Resposta — 404 Not Found** quando o slug não existe.

---

### 5. Atualizar template

```
PUT /v2/email/templates/{slug}
Content-Type: application/json
```

**Body** — mesmos campos do registro (exceto `slug`, que vem na URL)

```json
{
  "name": "Confirmação de Cadastro v2",
  "fromEmail": "noreply@meuapp.com",
  "htmlTemplate": "<html><body><h1>Olá {{.username}}, tudo bem?</h1></body></html>"
}
```

**Resposta — 200 OK**

```json
{
  "status": "updated",
  "message": "template updated successfully"
}
```

---

### 6. Deletar template

```
DELETE /v2/email/templates/{slug}
```

**Resposta — 204 No Content**

---

## Sintaxe de variáveis no template

Os templates utilizam a sintaxe nativa do Go (`html/template`).

| Sintaxe              | Descrição                                    |
|----------------------|----------------------------------------------|
| `{{.nomeVar}}`       | Substitui pelo valor da variável             |
| `{{if .nomeVar}}...{{end}}` | Bloco condicional (renderiza se não vazio) |
| `{{range .lista}}...{{end}}` | Iteração sobre lista                  |

**Exemplo com condicional:**

```html
<p>Olá, {{.username}}!</p>
{{if .link}}
  <a href="{{.link}}">Clique aqui</a>
{{end}}
```

> **Atenção:** se o HTML do template contiver `{{` fora de variáveis (ex: JavaScript inline), escape com `{{ "{{" }}`.

---

## Códigos de resposta

| Código | Situação                                       |
|--------|------------------------------------------------|
| 200    | Sucesso                                        |
| 201    | Template criado                                |
| 204    | Template deletado                              |
| 400    | Dados inválidos ou template não encontrado     |
| 405    | Método HTTP não permitido                      |
| 500    | Erro interno                                   |

---

## Exemplo completo: nova integração

### Passo 1 — Registrar o template uma única vez

```bash
curl -X POST https://<host>/v2/email/templates \
  -H "Content-Type: application/json" \
  -d '{
    "slug": "boas-vindas-meuapp",
    "name": "Boas-vindas MeuApp",
    "fromEmail": "noreply@meuapp.com",
    "htmlTemplate": "<html><body><h2>Bem-vindo, {{.nome}}!</h2><p>{{.descricao}}</p></body></html>"
  }'
```

### Passo 2 — Enviar quando necessário

```bash
curl -X POST https://<host>/v2/email/send \
  -H "Content-Type: application/json" \
  -d '{
    "templateSlug": "boas-vindas-meuapp",
    "to": "novo@usuario.com",
    "subject": "Bem-vindo ao MeuApp!",
    "variables": {
      "nome": "Carlos",
      "descricao": "Sua conta foi criada com sucesso."
    }
  }'
```

---

## Diferenças em relação à v1

| Aspecto            | v1 (`/email`)                  | v2 (`/v2/email`)                         |
|--------------------|-------------------------------|------------------------------------------|
| Templates          | Hardcoded no código Go         | Armazenados no banco, gerenciáveis via API |
| Variáveis          | Struct tipada (`customBodyProps`) | `map[string]string` livre               |
| Nova integração    | Requer redeploy                | Apenas um POST para registrar o template |
| Compatibilidade    | Mantida, sem alterações        | —                                        |
