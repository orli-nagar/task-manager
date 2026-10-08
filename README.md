# Task Management Service

A REST API built with Go and PostgreSQL for creating, retrieving, updating, and deleting tasks.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/tasks` | Create a task |
| GET | `/tasks` | Retrieve all tasks |
| PATCH | `/tasks/{id}` | Update a task |
| DELETE | `/tasks/{id}` | Delete a task |

Example POST request:

```json
{
  "title": "Learn Go",
  "description": "Practice building an HTTP service"
}
```

New tasks are created with `completed: false`. PATCH requests update only the provided fields.

## Running Locally

The application requires PostgreSQL. Configure the connection in a `.env` file:

```env
DATABASE_URL=postgresql://taskmanager:your_password@localhost:5432/taskmanager?sslmode=disable
```

Start the application:

```bash
go run .
```

Alternatively, run the application and database together using Podman Compose:

```bash
podman compose up --build
```

The API is available at `http://localhost:8080`.

## Deployment to OpenShift

The `deploy/` directory contains the OpenShift manifests for the application and PostgreSQL, including Deployments, Services, and a Route.

The application runs with **two replicas**, both connected to the same PostgreSQL database.

### Deploy

With access to an OpenShift cluster and the `oc` CLI configured, run:

```bash
./deploy.sh
```

The script creates the necessary resources, builds the application image, deploys the database and application, and prints the public API URL. On first deployment, it prompts for a database password.

### Verify

Check that PostgreSQL and both application replicas are running:

```bash
oc get deployments,services,routes
```

Expected Deployment readiness: `db` at `1/1` and `task-manager` at `2/2`.

Retrieve the API hostname:

```bash
oc get route task-manager -o jsonpath='{.spec.host}{"\n"}'
```

Use `http://<route-hostname>/tasks` to test the API with Postman or curl.

The deployment was tested on a temporary OpenShift cluster provisioned through Red Hat Cluster Bot.

## Storage

Tasks are stored in PostgreSQL and shared across both application replicas. The OpenShift database currently uses ephemeral storage, so data may be lost if its Pod is replaced.