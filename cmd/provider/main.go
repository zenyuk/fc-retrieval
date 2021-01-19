package main

import (
  "github.com/ConsenSys/fc-retrieval-provider/internal/config"
  "github.com/ConsenSys/fc-retrieval-provider/internal/logger"
	"github.com/ConsenSys/fc-retrieval-provider/internal/services/provider"
	_ "github.com/joho/godotenv/autoload"
  "github.com/rs/zerolog/log"
  "go.uber.org/fx"
)

func main() {
  conf := config.NewConfig()
  logger.InitLogger(conf)

  log.Debug().Msg("Running app ...")
  app := fx.New(
    config.Module,
    provider.Module,
    fx.Invoke(provider.Start),
  )
  app.Run()
}
