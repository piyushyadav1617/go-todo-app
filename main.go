package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/piyushyadav1617/go-todo-app/database"
)

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Completed   bool    `json:"completed"`
}
type TodoUpdate struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	rows, err := db.Query("SELECT id, title, description, completed FROM todos")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		var description *string
		err := rows.Scan(&todo.ID, &todo.Title, &description, &todo.Completed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
			return
		}
		todo.Description = description
		todos = append(todos, todo)
	}

	json.NewEncoder(w).Encode(todos)
}

func getTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	var todo Todo
	var description *string
	row := db.QueryRow("SELECT id, title, description, completed FROM todos WHERE id = $1", id)
	err = row.Scan(&todo.ID, &todo.Title, &description, &todo.Completed)
	todo.Description = description
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	} else if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", err.Error())})
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	var todo Todo
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid todo data"})
		return
	}
	row := db.QueryRow("INSERT INTO todos (title, description, completed) VALUES ($1, $2, $3) RETURNING id", todo.Title, todo.Description, todo.Completed)
	var id int
	err := row.Scan(&id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", err.Error())})
		return
	}
	todo.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func updateTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	var updates TodoUpdate
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid todo data: %v", err.Error()),
		})
		return
	}
	if updates.Title == nil && updates.Description == nil && updates.Completed == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No updates provided",
		})
		return
	}
	_, er := db.Exec(`UPDATE todos SET title = COALESCE($1, title), description = COALESCE($2, description), completed = COALESCE($3, completed) WHERE id = $4`, updates.Title, updates.Description, updates.Completed, id)

	if er == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", er.Error())})
		return
	} else if er != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", er.Error())})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}

func deleteTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	_, er := db.Exec(`DELETE FROM todos WHERE id = $1`, id)
	if er == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", er.Error())})
		return
	} else if er != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", er.Error())})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}

func main() {
	connStr := os.Getenv("DATABASE_URL")
	err := database.InitDB(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	http.HandleFunc("GET /todos", getAllTodos)
	http.HandleFunc("GET /todos/{id}", getTodoWithGivenId)
	http.HandleFunc("POST /todos", createTodo)
	http.HandleFunc("PATCH /todos/{id}", updateTodoWithGivenId)
	http.HandleFunc("DELETE /todos/{id}", deleteTodoWithGivenId)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
