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
	storage *pgsql.Storage,
	tokenTTL time.Duration,
) *App {
	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
