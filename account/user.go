package account

import "context"

type User struct {
	// agregar JSONS
	// email string 'json:"email"'
	ID       string `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

//- UUID
//- email
//- password
//- status

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, id string) (string, string, string, error)
	GetUsers(ctx context.Context) ([5]string, error)
	ValidateUser(ctx context.Context, email string) (string, error)
	NewPassword(ctx context.Context, email string, newpassword string) (string, error)
	ValidateToken(ctx context.Context, email string) (string, error)
	UpdateToken(ctx context.Context, email string, token string) (string, error)
}
