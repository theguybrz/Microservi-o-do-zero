# Microserviço Task Processor (Go + SQLite)

Pequeno microserviço em Go para criação e listagem de tarefas com processamento assíncrono via worker.

## Tecnologias
- Go
- SQLite (modernc.org/sqlite)
- REST

## Como rodar

```bash
go mod tidy
go run main.go
```

Servidor: `http://localhost:8081`

## Endpoints

- POST /tasks
- GET /tasks

### Exemplo JSON

{
  "title": "Test task",
  "description": "This is a test task"
}
