package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

//Los endpoints funcionana como RPCs
type Endpoints struct {
	CreateUser    endpoint.Endpoint
	GetUser       endpoint.Endpoint
	GetUsers      endpoint.Endpoint
	ValidateUser  endpoint.Endpoint
	ValidateToken endpoint.Endpoint
	Close         endpoint.Endpoint
	NewPassword   endpoint.Endpoint
	GetId         endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:    makeCreateUserEndpoint(s),
		GetUser:       makeGetUserEndpoint(s),
		GetUsers:      makeGetUsersEndpoint(s),
		ValidateUser:  makeValidateUserEndpoint(s),
		ValidateToken: makeValidateTokenEndpoint(s),
		Close:         makeCloseEndpoint(s),
		NewPassword:   makeNewPasswordEndpoint(s),
		GetId:         makeGetIdEndpoint(s),
	}
}

// Estas funciones son un wrapper que nos permite usar nuestras funciones en las request
func makeCreateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateUserRequest)
		token, err := s.CreateUser(ctx, req.Email, req.Password, req.Status, req.UserName)
		return CreateUserResponse{Token: token}, err
	}
}

func makeGetUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetUserRequest)
		email, status, username, err := s.GetUser(ctx, req.Id)

		return GetUserResponse{
			Email: email, Status: status, UserName: username,
		}, err
	}
}

func makeGetIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetIdRequest)
		id, err := s.GetId(ctx, req.UserName)

		return GetIdResponse{
			Id: id,
		}, err
	}
}

func makeGetUsersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		username1, username2, username3, username4, username5, err := s.GetUsers(ctx)
		return GetUsersResponse{
			UserName1: username1, UserName2: username2, UserName3: username3, UserName4: username4, UserName5: username5,
		}, err
	}
}
func makeValidateUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateUserRequest)
		token, err := s.ValidateUser(ctx, req.Email, req.Password)

		return ValidateUserResponse{
			Token: token,
		}, err
	}
}

func makeValidateTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ValidateTokenRequest)
		token, err := s.ValidateToken(ctx, req.Email, req.Token)
		return ValidateTokenResponse{Token: token}, err
	}
}
func makeCloseEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CloseRequest)
		ok, err := s.CloseSesion(ctx, req.Ok)
		return CloseResponse{Ok: ok}, err
	}
}
func makeNewPasswordEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(NewPasswordRequest)
		ok, err := s.NewPassword(ctx, req.Email, req.Password, req.RePassword)
		return NewPasswordResponse{Ok: ok}, err
	}
}
