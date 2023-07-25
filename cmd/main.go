package main

import (
	"github.com/AlexCorn999/server-social-network/internal/apiserver"
	"github.com/AlexCorn999/server-social-network/internal/logger"
)

func main() {
	server := apiserver.New()
	if err := server.Start(); err != nil {
		logger.ForError(err)
	}
}
