# Microservices Architecture

This directory contains the microservices for the Todo application.

## Services

### Auth Service (`services/auth-service/`)

**Purpose:** Handles user authentication and JWT token management.

**Port:** 50051 (gRPC)

**Endpoints:**
- `Signup` - Creates a new user account
- `Login` - Authenticates user and returns JWT token
- `ValidateToken` - Validates JWT token and returns user information

**Dependencies:**
- MongoDB (shared database)
- User Repository
- Auth Usecase

**How to run:**
```bash
cd services/auth-service
go run main.go
```

### Todo Service (`services/todo-service/`)

**Purpose:** Handles todo CRUD operations with both gRPC and HTTP REST interfaces.

**Ports:** 
- 50052 (gRPC)
- 8080 (HTTP REST API)

**gRPC Endpoints:**
- `CreateTodo` - Creates a new todo item
- `GetTodo` - Retrieves a todo by ID
- `ListTodos` - Lists todos with pagination
- `UpdateTodo` - Updates a todo item
- `DeleteTodo` - Deletes a todo item

**HTTP REST Endpoints:**
- `POST /api/v1/todos` - Create todo
- `GET /api/v1/todos` - List todos (with pagination)
- `GET /api/v1/todos/:id` - Get todo by ID
- `PUT /api/v1/todos/:id` - Update todo
- `DELETE /api/v1/todos/:id` - Delete todo

**Dependencies:**
- MongoDB (shared database)
- Todo Repository
- Todo Usecase
- Auth Service (gRPC client for token validation)

**How to run:**
```bash
cd services/todo-service
go run main.go
```

## Architecture

```
┌─────────────────┐
│   HTTP Client   │
└────────┬────────┘
         │ REST API
         ▼
┌─────────────────────┐        gRPC         ┌─────────────────┐
│   Todo Service      │◄───────────────────►│  Auth Service   │
│  - HTTP Server      │  Token Validation   │  - gRPC Server  │
│  - gRPC Server      │                     └────────┬────────┘
└─────────┬───────────┘                              │
          │                                          │
          │            ┌──────────────┐             │
          └───────────►│   MongoDB    │◄────────────┘
                       │ (Shared DB)  │
                       └──────────────┘
```

## Communication Flow

1. **Client → Todo Service (HTTP)**: Client sends HTTP requests to Todo Service
2. **Todo Service → Auth Service (gRPC)**: Todo Service validates JWT tokens via Auth Service
3. **Todo Service → MongoDB**: Todo Service performs CRUD operations on todos
4. **Auth Service → MongoDB**: Auth Service manages user authentication

## Configuration

Configuration is managed via `configs/app.yaml`:

```yaml
grpc:
  auth_service:
    port: 50051
    host: localhost
  todo_service:
    port: 50052
    host: localhost

server:
  port: 8080  # HTTP REST API port for Todo Service

database:
  mongodb:
    uri: mongodb://localhost:27017
    database: todoapp
    timeout: 10

jwt:
  secret: your-secret-key-change-in-production
  expiration: 24
```

## Running All Services

1. Start MongoDB:
```bash
docker-compose up -d mongodb
```

2. Start Auth Service:
```bash
cd services/auth-service
go run main.go
```

3. Start Todo Service (in another terminal):
```bash
cd services/todo-service
go run main.go
```

## Testing the Services

### 1. Signup
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123","name":"John Doe"}'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 3. Create Todo (with JWT token)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"My Todo","description":"Todo description"}'
```

### 4. List Todos
```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Shared Components

Both services share the following components:
- `internal/app/todoapp/domain/` - Domain entities and interfaces
- `internal/app/todoapp/repository/` - MongoDB repositories
- `internal/app/todoapp/usecase/` - Business logic
- `internal/pkg/config/` - Configuration management
- `internal/pkg/database/` - Database connection
- `internal/pkg/logger/` - Logging utilities

## Proto Files

gRPC service definitions are in:
- `proto/auth/auth.proto` - Auth service definition
- `proto/todo/todo.proto` - Todo service definition

Generated gRPC code:
- `proto/auth/auth.pb.go` and `proto/auth/auth_grpc.pb.go`
- `proto/todo/todo.pb.go` and `proto/todo/todo_grpc.pb.go`

To regenerate proto files:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/auth/auth.proto proto/todo/todo.proto
```
