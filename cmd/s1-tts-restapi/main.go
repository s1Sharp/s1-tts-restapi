package main

import (
	"context"
	"github.com/s1Sharp/s1-tts-restapi/routes"
	"os"

	"github.com/s1Sharp/s1-tts-restapi/internal/config"
	"github.com/s1Sharp/s1-tts-restapi/internal/logger"
)

var (
	log = logger.GetLogger()
)

func main() {
	defer log.Close()

	// read config from env
	cfg := config.ReadEnv()
	ctx := context.TODO()

	s1TtsServer := routes.NewServer(ctx, cfg)
	if err := s1TtsServer.Run(); err != nil {
		log.Error(err)
	}

	os.Exit(0)
}
