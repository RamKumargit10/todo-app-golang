#  Todo List API – Go + Chi

A basic Todo List API built in Go with support for creating, reading, updating, and deleting tasks.
Uses in-memory storage and basic authentication.

---

## Setup Instructions

```bash
# Navigate to the project directory
cd todo-app

# Run the server
go run main.go
```

---

##  Authentication

All API requests require **Basic Auth**:

```
Username: user
Password: password
```

Send it in the header:

```
Authorization: Basic dXNlcjpwYXNzd29yZA==
```

---

##  API Endpoints

| Method | Endpoint       | Description              |
|--------|----------------|--------------------------|
| POST   | `/tasks`       | Add a new task           |
| GET    | `/tasks`       | Fetch all tasks          |
| PUT    | `/tasks/{id}`  | Update a task by ID      |
| DELETE | `/tasks/{id}`  | Delete a task by ID      |

---

## Sample Request/Response

### Add Task

**POST** `/tasks`

**Body:**

```json
{
  "title": "Learn Go"
}
```

**Response:**

```json
{
  "message": "Task added successfully"
}
```

---

### ✅ Update Task

**PUT** `/tasks/1`

**Body:**

```json
{
  "title": "Learn Go",
  "completed": true
}
```

---

##  Testing
```bash
go test ./handlers
```

Test includes:
- Adding a task
- Getting all tasks
- Updating task
- Delete a task
- Handling invalid inputs
- Handling deletion of non-existent task

