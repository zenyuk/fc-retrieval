package main

import (
	_ "github.com/joho/godotenv/autoload"
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"

	"github.com/ConsenSys/fc-retrieval-provider/config"
	"github.com/ConsenSys/fc-retrieval-provider/internal/logging"
	"github.com/ConsenSys/fc-retrieval-provider/internal/services/provider"
)

func main() {
  conf := config.NewConfig()
  logging.InitLogger(conf)

	log.Info("Running app ...")
	p := provider.NewProvider(conf)
	provider.Start(p)
}
