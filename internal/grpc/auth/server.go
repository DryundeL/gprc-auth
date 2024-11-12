package auth

import (
	"context"
	"errors"
	ssov1 "github.com/DryundeL/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	validator *validator.Validate
}

func NewServerAPI() *ServerAPI {
	return &ServerAPI{
		validator: validator.New(),
	}
}

func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, NewServerAPI())
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func (s *ServerAPI) Login(
	ctx context.Context, request *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	loginRequest := LoginRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	if err := s.validator.Struct(loginRequest); err != nil {
		return nil, errors.New("invalid request data: " + err.Error())
	}

	return &ssov1.LoginResponse{
		Token: "some_generated_token",
	}, nil
}

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

func (s *ServerAPI) Register(
	ctx context.Context, request *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	registerRequest := RegisterRequest{
		Email:    request.Email,
		Password: request.Password,
	}

	if err := s.validator.Struct(registerRequest); err != nil {
		return nil, errors.New("invalid registration data: " + err.Error())
	}

	return &ssov1.RegisterResponse{
		UserId: 2,
	}, nil
}

func (s *ServerAPI) IsAdmin(
	ctx context.Context, request *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	return &ssov1.IsAdminResponse{
		IsAdmin: true,
	}, nil
}
