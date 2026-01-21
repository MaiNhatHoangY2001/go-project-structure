# Microservices Architecture

This directory contains the microservices for the Todo application with an API Gateway pattern.

## Services

### API Gateway (`services/api-gateway/`)

**Purpose:** Single entry point for all client requests. Routes HTTP REST requests to backend services via gRPC.

**Port:** 8080 (HTTP REST API)

**Features:**
- HTTP REST API for clients
- Protocol translation (HTTP → gRPC)
- JWT token validation
- Request routing
- Swagger documentation at `/swagger/*`
- Health check at `/health`

**HTTP REST Endpoints:**
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/todos` - Create todo (requires auth)
- `GET /api/v1/todos` - List todos (requires auth)
- `GET /api/v1/todos/:id` - Get todo (requires auth)
- `PUT /api/v1/todos/:id` - Update todo (requires auth)
- `DELETE /api/v1/todos/:id` - Delete todo (requires auth)

**Dependencies:**
- Auth Service (gRPC client)
- Todo Service (gRPC client)

**How to run:**
```bash
cd services/api-gateway
go run main.go
```

### Auth Service (`services/auth-service/`)

**Purpose:** Handles user authentication and JWT token management.

**Port:** 50051 (gRPC only - internal service)

**gRPC Endpoints:**
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

**Purpose:** Handles todo CRUD operations via gRPC.

**Port:** 50052 (gRPC only - internal service)

**gRPC Endpoints:**
- `CreateTodo` - Creates a new todo item
- `GetTodo` - Retrieves a todo by ID
- `ListTodos` - Lists todos with pagination
- `UpdateTodo` - Updates a todo item
- `DeleteTodo` - Deletes a todo item

**Dependencies:**
- MongoDB (shared database)
- Todo Repository
- Todo Usecase

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
         │ REST API (HTTP)
         ▼
┌──────────────────────┐
│   API Gateway        │
│  - HTTP Server       │
│  - Request Router    │
└─────┬──────────┬─────┘
      │          │
      │ gRPC     │ gRPC
      ▼          ▼
┌──────────────┐  ┌─────────────────┐
│ Auth Service │  │  Todo Service   │
│ - gRPC       │  │  - gRPC         │
└──────┬───────┘  └────────┬────────┘
       │                   │
       │  ┌──────────────┐ │
       └─►│   MongoDB    │◄┘
          │ (Shared DB)  │
          └──────────────┘
```

## Communication Flow

1. **Client → API Gateway (HTTP)**: Client sends HTTP REST requests to API Gateway
2. **API Gateway → Auth Service (gRPC)**: Gateway validates JWT tokens and routes auth requests
3. **API Gateway → Todo Service (gRPC)**: Gateway routes todo requests after token validation
4. **Auth Service → MongoDB**: Auth Service manages user authentication
5. **Todo Service → MongoDB**: Todo Service performs CRUD operations on todos

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

### Prerequisites
1. MongoDB running on port 27017

### Start Services (in order)

1. **Start MongoDB:**
```bash
docker-compose up -d mongodb
# Or if MongoDB is already running locally, skip this step
```

2. **Start Auth Service:**
```bash
cd services/auth-service
go run main.go
# Service starts on port 50051 (gRPC)
```

3. **Start Todo Service (in another terminal):**
```bash
cd services/todo-service
go run main.go
# Service starts on port 50052 (gRPC)
```

4. **Start API Gateway (in another terminal):**
```bash
cd services/api-gateway
go run main.go
# Gateway starts on port 8080 (HTTP REST)
```

### Access Points

- **API Gateway (for clients):** http://localhost:8080
- **Swagger Documentation:** http://localhost:8080/swagger/index.html
- **Health Check:** http://localhost:8080/health
- **Auth Service:** localhost:50051 (gRPC - internal only)
- **Todo Service:** localhost:50052 (gRPC - internal only)

## Testing the Services

All client requests go through the API Gateway at `http://localhost:8080`.

### 1. Signup
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "507f1f77bcf86cd799439011",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 3. Create Todo (with JWT token)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

### 4. List Todos
```bash
curl -X GET "http://localhost:8080/api/v1/todos?page=1&pageSize=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. Get Todo by ID
```bash
curl -X GET http://localhost:8080/api/v1/todos/TODO_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Update Todo
```bash
curl -X PUT http://localhost:8080/api/v1/todos/TODO_ID \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy groceries and cook",
    "description": "Milk, eggs, bread, chicken",
    "completed": true
  }'
```

### 7. Delete Todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/TODO_ID \
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
