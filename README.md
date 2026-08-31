# Task Management Service

A simple HTTP service written in Go for creating and retrieving tasks.

Tasks are stored in memory and are not persisted after the server is stopped.

## Running the Service

Run the application with:

```bash
go run .
```

The server will start on port `8080`.

## API Endpoints

### Create a Task

**POST** `/tasks`

Example request body:

```json
{
  "title": "Learn Go",
  "description": "Practice building an HTTP service",
  "completed": false
}
```

The server assigns an ID to each new task.

Example response:

```json
{
  "ID": 1,
  "Title": "Learn Go",
  "Description": "Practice building an HTTP service",
  "Completed": false
}
```

### Get All Tasks

**GET** `/tasks`

Returns all tasks currently stored in memory.

Example response:

```json
[
  {
    "ID": 1,
    "Title": "Learn Go",
    "Description": "Practice building an HTTP service",
    "Completed": false
  }
]
```

## Storage

Tasks are stored in memory using a Go map. Restarting the service clears all existing tasks.
