package account

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Status   string `json:"status"`
		UserName string `json:"username"`
	}
	CreateUserResponse struct {
		Token string `json:"token"`
	}

	GetUserRequest struct {
		Id string `json:"id"`
	}
	GetUserResponse struct {
		Email    string `json:"email"`
		Status   string `json:"status"`
		UserName string `json:"username"`
	}

	ValidateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	ValidateUserResponse struct {
		Token string `json:"token"`
	}

	NewPasswordRequest struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RePassword string `json:"repassword"`
	}
	NewPasswordResponse struct {
		Ok string `json:"ok"`
	}

	ValidateTokenRequest struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	ValidateTokenResponse struct {
		Token string `json:"token"`
	}

	CloseRequest struct {
		Ok string `json:"ok"`
	}
	CloseResponse struct {
		Ok string `json:"ok"`
	}
)

// Funciones que nos permites codificar/decodificar los JSON
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Las decodifica en nuestras interfaces
func decodeUserReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Crea la interface a partir de la respuesta del servidor
func decodeEmailReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GetUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ValidateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateTokenReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req ValidateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeCloseReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CloseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeNewPasswordReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req NewPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
