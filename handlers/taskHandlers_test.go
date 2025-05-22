package taskHandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"todo-app/models"
	"todo-app/store"

	"github.com/go-chi/chi/v5"
)

func resetTaskList() {
	store.TL.Tasks = []models.ToDo{}
}

func TestAddTask(t *testing.T) {
	resetTaskList()

	body := `{"title": "Test Task"}`
	req := httptest.NewRequest("POST", "/tasks", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	AddTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}
	if resp.Message != "Task added successfully" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestGetTasks(t *testing.T) {
	resetTaskList()
	store.TL.AddTask("Task 1")
	store.TL.AddTask("Task 2")

	req := httptest.NewRequest("GET", "/tasks", nil)
	rr := httptest.NewRecorder()
	GetTasks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var tasks []models.ToDo
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatal("Failed to decode tasks:", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestUpdateTask(t *testing.T) {
	resetTaskList()
	store.TL.AddTask("Old Title")

	body := `{"title": "New Title", "completed": true}`
	req := httptest.NewRequest("PUT", "/tasks/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// add chi URL param manually
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	UpdateTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Message != "Task updated successfully" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestDeleteTask(t *testing.T) {
	resetTaskList()
	store.TL.AddTask("Task to delete")

	req := httptest.NewRequest("DELETE", "/tasks/2", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	DeleteTask(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var resp models.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if resp.Message != "Task deleted successfully" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestInvalidUpdateTaskID(t *testing.T) {
	resetTaskList()
	store.TL.AddTask("Task to delete")

	body := `{"title": "New Title", "completed": true}`
	req := httptest.NewRequest("PUT", "/tasks/a", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "a")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	UpdateTask(rr, req)
	fmt.Println(rr.Body.String())

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}
}

func TestDeleteNonexistentTask(t *testing.T) {
	resetTaskList()

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(99))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	DeleteTask(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", rr.Code)
	}
}
