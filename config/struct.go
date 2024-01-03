package config

import "time"

type (
	AppConfiguration struct {
		OpenID OpenID `fig:"openid" validate:"required"`
		Server Server `fig:"server" validate:"required"`
	}

	OpenID struct {
		Audience string `fig:"audience" validate:"required"`
		Issuer   string `fig:"issuer" validate:"required"`
		KeyPath  string `fig:"keyPath" validate:"required"`
	}

	Server struct {
		GracefulTimeout time.Duration `fig:"gracefulTimeout" default:"5s"`
		Name            string        `fig:"name"`
		Addr            string        `fig:"addr" default:":8080"`
		TLS             bool          `fig:"tls"`
	}
)
