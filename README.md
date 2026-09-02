# Task Management Service

A simple HTTP service written in Go for managing tasks.

Tasks are stored in memory and are not persisted after the server is stopped.

## Running the Service

Run the application locally:

```bash
go run .
```

The server runs on port `8080`.

## API Endpoints

### Create a Task

**POST** `/tasks`

Example request:

```json
{
  "title": "Learn Go",
  "description": "Practice building an HTTP service"
}
```

New tasks are created with `completed` set to `false`.

### Get All Tasks

**GET** `/tasks`

Returns all tasks currently stored in memory.

### Update a Task

**PATCH** `/tasks/{id}`

Updates the provided fields of an existing task.

Example request:

```json
{
  "title": "Learn Go concurrency",
  "completed": true
}
```

### Delete a Task

**DELETE** `/tasks/{id}`

Deletes the task with the given ID.

## Running with Docker

Build the image:

```bash
podman build -t task-manager .
```

Run the container:

```bash
podman run -p 8080:8080 task-manager
```

The service will be available on port `8080`.

## Storage

Tasks are stored in memory using a Go map. Access to the shared task storage is protected for concurrent requests.

Restarting the service clears all existing tasks.