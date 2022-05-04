package account

import "context"

type Service interface {
	CreateUser(ctx context.Context, email string, password string, status string, username string) (string, error)
	GetUser(ctx context.Context, id string) (string, string, string, error)
	ValidateUser(ctx context.Context, email string, password string) (string, error)
	ValidateToken(ctx context.Context, token string, email string) (string, error)
	NewPassword(ctx context.Context, email string, password string, repassword string) (string, error)
	CloseSesion(ctx context.Context, ok string) (string, error)
}
