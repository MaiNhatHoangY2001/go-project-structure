package dto

// UserSignupRequest represents a signup request
type UserSignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// UserLoginRequest represents a login request
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse represents a login response
type UserLoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// TodoCreateRequest represents a request to create a todo
type TodoCreateRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// TodoUpdateRequest represents a request to update a todo
type TodoUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

// TodoListQuery represents query parameters for listing todos
type TodoListQuery struct {
	Page      int   `form:"page" binding:"min=1"`
	PageSize  int   `form:"pageSize" binding:"min=1,max=100"`
	Completed *bool `form:"completed"`
}

// TodoResponse represents a todo response
type TodoResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
