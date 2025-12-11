
# Microserviço Go — Task Manager 🧩

Um microserviço simples em **Go** usando **SQLite** e **goroutines** para criar, listar e processar tarefas em background.

---

## 🚀 Como rodar o projeto

### **1. Instale as dependências**
Este projeto usa o driver SQLite puro:

```bash
go get modernc.org/sqlite
```

---

## **2. Rodar o servidor**

Na raiz do projeto:

```bash
go run main.go
```

Se tudo der certo, você verá:

```
Server running on port 8081
```

---

# 🧪 Testando a API

Você pode testar de 3 formas:

---

# ✅ 1. Usando REST Client no VS Code  
Certifique-se de instalar a extensão:

```
REST Client — by Huachao Mao
```

Crie um arquivo `requests.http` com:

```
### Criar tarefa
POST http://localhost:8081/tasks
Content-Type: application/json

{
    "title": "Estudar Go",
    "description": "Microserviço brabo!"
}

### Listar tarefas
GET http://localhost:8081/tasks
```

Depois é só clicar em **Send Request**.

---

# ✅ 2. Usando curl

### Criar tarefa:
```bash
curl -X POST http://localhost:8081/tasks     -H "Content-Type: application/json"     -d '{"title":"Minha Task","description":"Testando!"}'
```

### Listar tarefas:
```bash
curl http://localhost:8081/tasks
```

---

# 🗂 Estrutura do projeto

```
microservico/
│── main.go
│── tasks.db
│── README.md
│── .gitignore
```

---

# ⚙️ Como funciona

### ✔ POST /tasks  
Cria uma nova tarefa, armazena no banco e envia para o canal para ser processada.

### ✔ GET /tasks  
Retorna todas as tarefas cadastradas.

### ✔ Processamento em background  
Cada tarefa criada vai para `TaskChannel` e é processada por uma goroutine:

- Aguarda 5 segundos (simulação)
- Atualiza status para **completed**
- Loga o processamento

---

# 🐳 Rodando com Docker (opcional)

Se quiser usar Docker, crie um arquivo `Dockerfile`:

```Dockerfile
FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]
```

Build:

```bash
docker build -t microservico-go .
```

Rodar:

```bash
docker run -p 8081:8081 microservico-go
```

---

# 📄 Licença

MIT License — fique livre para usar, modificar e compartilhar.

---

# 🙌 Contribuindo

Pull requests são bem-vindos!  
Se achar um bug, abre uma issue que eu te ajudo a resolver.

---

Feito com 💙 e bastante café ☕ — por **Guy**
