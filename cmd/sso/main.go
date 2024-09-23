package main

import (
	"fmt"
	"grpc-auth/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: иницализоровать логгер

	// TODO: иницализация приложения (app)

	// TODO: запустить gRPC сервер приложения
}
