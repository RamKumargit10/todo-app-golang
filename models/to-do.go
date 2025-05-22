package models

import (
	"errors"
	"fmt"
)

type toDoManager interface {
	AddTask(task string) int
	DeleteTask(id int) error
	GetAllTasks() []ToDo
}

type ToDo struct {
	ID        int
	Title     string
	Completed bool
}

type TaskList struct {
	Tasks []ToDo
}

type Response struct {
	Message string `json:"message"`
}

func (t *ToDo) MarkAsCompleted() {
	t.Completed = true
}

func (t *TaskList) AddTask(task string) {
	var currentTask ToDo = ToDo{
		ID:        len(t.Tasks) + 1,
		Title:     task,
		Completed: false,
	}
	t.Tasks = append(t.Tasks, currentTask)
}

func (t *TaskList) DeleteTask(id int) error {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			return nil
		}
	}
	return errors.New(fmt.Sprintf("Task with id %v not found", id))
}

func (t *TaskList) GetAllTasks() []ToDo {
	return t.Tasks
}
