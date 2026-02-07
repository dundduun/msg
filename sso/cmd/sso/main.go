package main

import (
	"fmt"
	"github.com/dundduun/msg/sso/internal/config"
)

func main() {
	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: запустить логгер

	// TODO: запустить приложение

	// TODO: запустить gRPC-сервер
}
