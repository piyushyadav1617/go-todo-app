package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/piyushyadav1617/go-todo-app/internal/database"
	"github.com/piyushyadav1617/go-todo-app/internal/model"
)

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	rows, err := db.Query("SELECT id, title, description, completed FROM todos")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
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

func GetTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	var todo model.Todo
	var description *string
	err = db.QueryRow("SELECT id, title, description, completed FROM todos WHERE id = $1", id).Scan(&todo.ID, &todo.Title, &description, &todo.Completed)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Todo not found"})
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}

	todo.Description = description
	json.NewEncoder(w).Encode(todo)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	db := database.GetDB()
	err := db.QueryRow("INSERT INTO todos (title, description, completed) VALUES ($1, $2, $3) RETURNING id",
		todo.Title, todo.Description, todo.Completed).Scan(&todo.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}
	var updates model.TodoUpdate
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

	result, err := db.Exec(`UPDATE todos SET title = COALESCE($1, title), 
		description = COALESCE($2, description), 
		completed = COALESCE($3, completed) 
		WHERE id = $4`,
		updates.Title, updates.Description, updates.Completed, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", err.Error())})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Internal server error: %v", err.Error())})
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Todo not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}

func DeleteTodoWithGivenId(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	db := database.GetDB()
	result, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("%v", err.Error())})
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Todo not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted successfully"})
}
