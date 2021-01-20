package provider

import (
	log "github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
  "github.com/spf13/viper"
)

type Provider struct {
  conf *viper.Viper
}

func NewProvider(conf *viper.Viper) *Provider {
  return &Provider{
    conf: conf,
  }
}

func Start(p *Provider) {
  scheme := p.conf.GetString("SERVICE_SCHEME")
  host := p.conf.GetString("SERVICE_HOST")
  port := p.conf.GetString("SERVICE_PORT")
  log.Info("Provider started at %s://%s:%s", scheme, host, port)
}