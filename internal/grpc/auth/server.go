package auth

import (
	"context"
	ssov1 "github.com/DryundeL/protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appId int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userId int, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, err error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	validator *validator.Validate
	auth      Auth
}

func NewServerAPI(auth Auth) *ServerAPI {
	return &ServerAPI{
		validator: validator.New(),
		auth:      auth,
	}
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, NewServerAPI(auth))
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
	AppId    int    `validate:"required"`
}

func (s *ServerAPI) Login(
	ctx context.Context, request *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	loginRequest := LoginRequest{
		Email:    request.Email,
		Password: request.Password,
		AppId:    int(request.AppId),
	}

	if err := s.validator.Struct(loginRequest); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, request.GetEmail(), request.GetPassword(), int(request.GetAppId()))
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		Token: token,
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
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
