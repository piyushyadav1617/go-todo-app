package route

import (
	"net/http"

	"github.com/piyushyadav1617/go-todo-app/internal/handler"
)

func TodoRoutes() {
	http.HandleFunc("GET /todos", handler.GetAllTodos)
	http.HandleFunc("GET /todos/{id}", handler.GetTodoWithGivenId)
	http.HandleFunc("POST /todos", handler.CreateTodo)
	http.HandleFunc("PATCH /todos/{id}", handler.UpdateTodoWithGivenId)
	http.HandleFunc("DELETE /todos/{id}", handler.DeleteTodoWithGivenId)
}
