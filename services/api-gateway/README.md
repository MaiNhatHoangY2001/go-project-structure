# API Gateway Service

The API Gateway is the single entry point for all client requests in the microservices architecture. It routes requests to the appropriate backend services (Auth Service and Todo Service) via gRPC.

## Architecture

```
Client (HTTP) → API Gateway (REST) → Auth Service (gRPC)
                      ↓
                   Todo Service (gRPC)
```

## Features

- **Unified HTTP REST API** for clients
- **gRPC communication** with backend services
- **JWT token validation** via Auth Service
- **Request routing** to appropriate services
- **Standardized response format** with error codes
- **Swagger/OpenAPI documentation**
- **Health check endpoint**
- **Graceful shutdown**

## Endpoints

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User login

### Todos (requires authentication)
- `POST /api/v1/todos` - Create todo
- `GET /api/v1/todos` - List todos (with pagination)
- `GET /api/v1/todos/:id` - Get todo by ID
- `PUT /api/v1/todos/:id` - Update todo
- `DELETE /api/v1/todos/:id` - Delete todo

### System
- `GET /health` - Health check
- `GET /swagger/*` - Swagger documentation

## Running the Service

### Prerequisites
- Auth Service running on port 50051
- Todo Service running on port 50052

### Start the Gateway
```bash
cd services/api-gateway
go run main.go
```

The gateway will start on port 8080 (default).

## Configuration

The gateway uses the shared configuration from `configs/app.yaml`:

```yaml
environment: development
server:
  port: "8080"
  mode: debug
auth_service:
  host: localhost
  port: "50051"
todo_service:
  host: localhost
  port: "50052"
```

## Response Format

### Success Response
```json
{
  "success": true,
  "data": {
    "user_id": "123",
    "token": "eyJ..."
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": 10001,
    "message": "Unauthorized"
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

| Code | Description |
|------|-------------|
| 10001 | Unauthorized |
| 10002 | Bad Request |
| 10003 | Forbidden |
| 10004 | Not Found |
| 10005 | Internal Error |
| 10006 | Validation Error |
| 10007 | Conflict |
| 10008 | Missing Auth Header |
| 10009 | Invalid Auth Header |
| 10010 | Invalid Token |
| 10011 | Invalid Token Claims |
| 10012 | Missing User ID |

## Testing

### Example: Signup
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Example: Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Example: Create Todo
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

## Benefits

1. **Single Entry Point** - Clients only need to know one address
2. **Protocol Translation** - HTTP REST for clients, gRPC for services
3. **Centralized Auth** - Token validation happens at the gateway
4. **Service Abstraction** - Backend services can change without affecting clients
5. **Easy Scaling** - Scale gateway independently from backend services
6. **Monitoring** - Centralized logging and metrics collection point
