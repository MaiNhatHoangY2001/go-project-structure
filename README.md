# Todo App with MongoDB, JWT Authentication, and Gin Framework

A RESTful Todo application built with Go, featuring JWT authentication, MongoDB storage, Swagger documentation, and comprehensive testing.

## Features

- ✅ **JWT Authentication**: Secure user signup and login with JWT tokens
- ✅ **MongoDB Integration**: Persistent data storage with MongoDB
- ✅ **RESTful API**: Clean and organized API endpoints
- ✅ **Swagger Documentation**: Interactive API documentation
- ✅ **Standardized Responses**: Consistent API response format with i18n error codes
- ✅ **User-Todo Relationship**: One user can have multiple todos
- ✅ **CRUD Operations**: Full create, read, update, delete functionality
- ✅ **Pagination Support**: Efficient listing with pagination
- ✅ **Unit Tests**: Comprehensive service layer testing
- ✅ **Integration Tests**: End-to-end API testing
- ✅ **Structured Logging**: Zap logger integration

## Project Structure

```
.
├── cmd/
│   └── todoapp/
│       └── main.go           # Application entry point
├── internal/
│   ├── app/todoapp/
│   │   ├── handlers/         # HTTP request handlers
│   │   ├── middleware/       # JWT authentication middleware
│   │   ├── models/           # Data models and DTOs
│   │   ├── repositories/     # Data access layer
│   │   ├── routes/           # Route definitions
│   │   └── services/         # Business logic layer
│   └── pkg/
│       ├── config/           # Configuration management
│       ├── database/         # MongoDB connection
│       └── logger/           # Logging setup
├── configs/
│   └── app.yaml              # Application configuration
├── test/
│   └── integration/          # Integration tests
└── api/
    └── docs/                 # Swagger documentation
```

## API Endpoints

### Authentication
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
  port: 8080
  mode: debug # debug, release

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

# Or using local MongoDB installation
mongod
```

4. **Run the application**
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
