package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskService struct {
	DB          *sql.DB
	TaskChannel chan Task
}

// ======================
// Database Operations
// ======================

func (t *TaskService) AddTask(ts *Task) error {
	query := `
		INSERT INTO tasks (title, description, status, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := t.DB.Exec(query, ts.Title, ts.Description, ts.Status, ts.CreatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	ts.ID = int(id)
	return nil
}

func (t *TaskService) UpdateTaskStatus(ts Task) error {
	_, err := t.DB.Exec("UPDATE tasks SET status = ? WHERE id = ?", ts.Status, ts.ID)
	return err
}

func (t *TaskService) ListTasks() ([]Task, error) {
	rows, err := t.DB.Query(`
		SELECT id, title, description, status, created_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// ======================
// Background Worker
// ======================

func (t *TaskService) ProcessTasks() {
	for task := range t.TaskChannel {
		log.Printf("🔧 Processando tarefa: %s", task.Title)

		time.Sleep(5 * time.Second) // Simula processamento

		task.Status = "completed"
		if err := t.UpdateTaskStatus(task); err != nil {
			log.Printf("❌ Erro ao atualizar: %v", err)
			continue
		}

		log.Printf("✅ Tarefa concluída: %s", task.Title)
	}
}

// ======================
// HTTP Handlers
// ======================

func (t *TaskService) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	task.Status = "pending"
	task.CreatedAt = time.Now()

	if err := t.AddTask(&task); err != nil {
		http.Error(w, "Erro ao salvar tarefa", http.StatusInternalServerError)
		return
	}

	t.TaskChannel <- task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (t *TaskService) HandleListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.ListTasks()
	if err != nil {
		http.Error(w, "Erro ao listar tarefas", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

// ======================
// Setup / Main
// ======================

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT,
			created_at DATETIME
		)
	`)
	return err
}

func main() {
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := createTable(db); err != nil {
		log.Fatalf("Erro criando tabela: %v", err)
	}

	taskService := TaskService{
		DB:          db,
		TaskChannel: make(chan Task, 10),
	}

	go taskService.ProcessTasks()

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case http.MethodPost:
			taskService.HandleCreateTask(w, r)
		case http.MethodGet:
			taskService.HandleListTasks(w, r)
		default:
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		log.Println("🚀 Servidor rodando em http://localhost:8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("⏳ Encerrando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	close(taskService.TaskChannel)

	log.Println("👋 Servidor finalizado")
}
