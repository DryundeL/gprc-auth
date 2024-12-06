package app

import (
	grpcapp "grpc-auth/internal/app/grpc"
	"grpc-auth/internal/services/auth"
	"grpc-auth/internal/storage/pgsql"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	dbConn string,
	tokenTTL time.Duration,
) *App {
	storage, err := pgsql.New(dbConn)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
