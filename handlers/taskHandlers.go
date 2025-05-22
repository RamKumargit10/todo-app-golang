package taskHandlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"todo-app/models"
	"todo-app/store"

	"github.com/go-chi/chi/v5"
)

var users = map[string]string{
	"user": "password",
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header value
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Decode the Base64 part of the Authorization header
		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
			return
		}

		// Split the decoded string into username and password
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 || users[parts[0]] != parts[1] {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func AddTask(w http.ResponseWriter, r *http.Request) {
	writeToFile("AddTask called")
	var input map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if title, ok := input["title"].(string); ok {
		store.TL.AddTask(title)
	} else {
		http.Error(w, "Invalid input! Input should be of the type string", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := models.Response{
		Message: "Task added successfully",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	writeToFile("GetTasks called")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(FetchTasksConcurrently(&store.TL)); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	writeToFile("UpdateTask called")
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for i, task := range store.TL.Tasks {
		if task.ID == taskID {
			if title, ok := input["title"].(string); ok {
				store.TL.Tasks[i].Title = title
			} else {
				http.Error(w, "Invalid input! Input should be of the type string", http.StatusBadRequest)
				return
			}
			if completed, ok := input["completed"].(bool); ok {
				store.TL.Tasks[i].Completed = completed
			} else {
				http.Error(w, "Invalid input! Input should be of the type bool", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			resp := models.Response{
				Message: "Task updated successfully",
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, "Error encoding response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	writeToFile("DeleteTask called")
	id := chi.URLParam(r, "id")
	taskID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	err = store.TL.DeleteTask(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := models.Response{
		Message: "Task deleted successfully",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// sample channel usage
func FetchTasksConcurrently(taskManager *models.TaskList) []models.ToDo {
	var consumerWG sync.WaitGroup
	consumerWG.Add(1)
	var wg sync.WaitGroup
	var channel = make(chan models.ToDo)
	var tasks []models.ToDo = []models.ToDo{}

	// Consumer goroutine
	go func() {
		defer consumerWG.Done()
		for task := range channel {
			tasks = append(tasks, task)
		}
		fmt.Println("All tasks fetched")
	}()

	// Producer goroutine
	go func() {
		for _, task := range taskManager.GetAllTasks() {
			wg.Add(1)
			go func(t models.ToDo) {
				defer wg.Done()
				channel <- t
			}(task)
		}
		wg.Wait()
		close(channel)
	}()

	// This part will wait for all tasks to be sent to the channel
	// and the consumer to process them.
	consumerWG.Wait()

	return tasks
}

// sample task generator
func TaskGenerator(count int, taskManager *models.TaskList) {
	var taskId int
	if len(taskManager.Tasks) > 0 {
		taskId = taskManager.Tasks[len(taskManager.Tasks)-1].ID + 1
	} else {
		taskId = 1
	}

	for i := 0; i < count; i++ {
		taskName := fmt.Sprintf("Task-%d", taskId)
		taskManager.AddTask(taskName)
		taskId++
	}
}

func writeToFile(msg string) {
	filePath := "/data/output.txt"

	// Open the file (or create it if it doesn't exist)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return
	}
	defer file.Close()

	// Write a message to the file
	if _, err := file.WriteString(msg + "\n"); err != nil {
		log.Printf("Failed to write to file: %v", err)
	} else {
		log.Println("Successfully wrote to /data/output.txt")
	}
}
