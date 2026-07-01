package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RefreshRequest is now technically optional/empty since the token is read from the cookie,
// but we keep it here if you need to attach future payload data (like device info).
type RefreshRequest struct{}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

// AuthResponse is used internally between your Service and Handler.
// 💡 Added json:"-" so tokens are physically impossible to send in the JSON body.
type AuthResponse struct {
	AccessToken  string       `json:"-"`
	RefreshToken string       `json:"-"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	LastLogin string `json:"last_login,omitempty"`
}
