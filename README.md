# Todo App with Microservices Architecture

A microservices-based Todo application built with Go, featuring JWT authentication, MongoDB storage, gRPC inter-service communication, and comprehensive testing.

## 🏗️ Architecture

This application follows a **microservices architecture** with two independent services:

### Services

1. **Auth Service** (Port: 50051)
   - User authentication and JWT token management
   - gRPC server for internal communication
   - Handles Signup, Login, and Token Validation

2. **Todo Service** (Ports: 50052 gRPC, 8080 HTTP)
   - Todo CRUD operations
   - Dual interface: gRPC server + HTTP REST API
   - Validates tokens via Auth Service gRPC client
   - Client-facing REST API endpoints

### Communication Flow
```
Client (HTTP) → Todo Service (HTTP REST) → Auth Service (gRPC)
                     ↓
                  MongoDB (Shared Database)
```

## Features

- ✅ **Microservices Architecture**: Separate Auth and Todo services
- ✅ **gRPC Communication**: Efficient inter-service communication
- ✅ **JWT Authentication**: Secure user signup and login with JWT tokens
- ✅ **MongoDB Integration**: Shared persistent data storage
- ✅ **RESTful API**: Clean HTTP REST API for clients
- ✅ **Swagger Documentation**: Interactive API documentation
- ✅ **Standardized Responses**: Consistent API response format with i18n error codes
- ✅ **User-Todo Relationship**: One user can have multiple todos
- ✅ **CRUD Operations**: Full create, read, update, delete functionality
- ✅ **Pagination Support**: Efficient listing with pagination
- ✅ **Unit Tests**: Comprehensive service layer testing (14 tests, 100% pass)
- ✅ **Clean Architecture**: Domain-driven design with proper layering
- ✅ **Structured Logging**: Zap logger integration

## Project Structure

```
.
├── services/
│   ├── auth-service/          # Authentication microservice
│   │   ├── main.go            # Auth service entry point
│   │   └── grpc/              # gRPC server implementation
│   └── todo-service/          # Todo microservice
│       ├── main.go            # Todo service entry point
│       ├── grpc/              # gRPC server + Auth client
│       └── http/              # HTTP REST API handlers
├── proto/
│   ├── auth/                  # Auth service proto definitions
│   └── todo/                  # Todo service proto definitions
├── internal/
│   ├── app/todoapp/
│   │   ├── domain/            # Domain entities and interfaces
│   │   ├── usecase/           # Business logic layer
│   │   ├── repository/        # Data access layer
│   │   └── delivery/          # HTTP handlers (legacy monolith)
│   └── pkg/
│       ├── config/            # Configuration management
│       ├── database/          # MongoDB connection
│       └── logger/            # Logging setup
└── configs/
    └── app.yaml               # Application configuration
```

## API Endpoints

### Authentication (via Todo Service HTTP API)
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User login

### Todos (Protected)
- `POST /api/v1/todos` - Create a new todo
- `GET /api/v1/todos` - List todos (with pagination and filtering)
- `GET /api/v1/todos/:id` - Get a specific todo
- `PUT /api/v1/todos/:id` - Update a todo
- `DELETE /api/v1/todos/:id` - Delete a todo

### Other
- `GET /health` - Health check endpoint
- `GET /swagger/*any` - Swagger documentation

## API Response Format

### Success Response
```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": 10001,
    "message": "Error message"
  }
}
```

### Paginated Response
```json
{
  "success": true,
  "data": {
    "page": 1,
    "pageSize": 10,
    "totalElement": 100,
    "data": [...]
  }
}
```

## Error Codes

| Code  | Description               |
|-------|---------------------------|
| 10001 | Unauthorized (general)    |
| 10002 | Bad Request               |
| 10003 | Forbidden                 |
| 10004 | Not Found                 |
| 10005 | Internal Error            |
| 10006 | Validation Error          |
| 10007 | Conflict                  |
| 10008 | Missing Auth Header       |
| 10009 | Invalid Auth Header       |
| 10010 | Invalid Token             |
| 10011 | Invalid Token Claims      |
| 10012 | Missing User ID           |

## Prerequisites

- Go 1.24 or higher
- MongoDB 4.0 or higher

## Configuration

Edit `configs/app.yaml`:

```yaml
server:
  port: 8080       # HTTP REST API port (Todo Service)
  mode: debug      # debug, release

grpc:
  auth_service:
    port: 50051    # Auth Service gRPC port
    host: localhost
  todo_service:
    port: 50052    # Todo Service gRPC port
    host: localhost

database:
  mongodb:
    uri: mongodb://localhost:27017
    database: todoapp
    timeout: 10

jwt:
  secret: your-secret-key-change-in-production
  expiration: 24 # hours

logger:
  level: debug # debug, info, warn, error
  encoding: json # json, console
```

## Installation and Running

### Running Microservices

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-project-structure
```

2. **Install dependencies**
```bash
go mod download
```

3. **Start MongoDB** (if not already running)
```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or using Docker Compose
docker-compose up -d mongodb
```

4. **Start Auth Service** (Terminal 1)
```bash
cd services/auth-service
go run main.go
```
Output: `Auth Service listening on :50051`

5. **Start Todo Service** (Terminal 2)
```bash
cd services/todo-service
go run main.go
```
Output: 
- `Todo Service gRPC listening on :50052`
- `Todo Service HTTP listening on :8080`

### Running Legacy Monolithic Application

Alternatively, you can still run the original monolithic version:
```bash
go run cmd/todoapp/main.go
```

The server will start on `http://localhost:8080`

## Testing

### Run Unit Tests
```bash
go test -mod=mod ./internal/app/todoapp/services/... -v
```

### Run Integration Tests
```bash
go test -mod=mod ./test/integration/... -v
```

### Run All Tests
```bash
go test -mod=mod ./... -v
```

## API Documentation

Once the server is running, access the Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

## Example Usage

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

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "...",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

### 3. Create Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "Complete project",
    "description": "Finish the todo app"
  }'
```

### 4. List Todos
```bash
curl -X GET "http://localhost:8080/api/v1/todos?page=1&pageSize=10" \
  -H "Authorization: Bearer <your-token>"
```

### 5. Update Todo
```bash
curl -X PUT http://localhost:8080/api/v1/todos/<todo-id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "Updated title",
    "completed": true
  }'
```

### 6. Delete Todo
```bash
curl -X DELETE http://localhost:8080/api/v1/todos/<todo-id> \
  -H "Authorization: Bearer <your-token>"
```

## Building for Production

```bash
# Build the binary
go build -o todoapp cmd/todoapp/main.go

# Run the binary
./todoapp
```

## License

MIT

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
