package provider

import (
  "github.com/rs/zerolog/log"
  "github.com/spf13/viper"
  "go.uber.org/fx"
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
  log.Printf("Provider started at %s://%s:%s", scheme, host, port)
}

var Module = fx.Options(
  fx.Provide(
    NewProvider,
  ),
)
