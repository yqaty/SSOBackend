package config

import (
	"net/url"

	"github.com/UniqueStudio/UniqueSSOBackend/internal/constants"
	"github.com/gin-contrib/sessions"
	"github.com/spf13/viper"
	"github.com/xylonx/zapx"
	"go.uber.org/zap"
)

var Config *Setting

var TraefikRedirectURI *url.URL

var SessionOptions sessions.Options

func Setup(cfgFile string) error {
	Config = new(Setting)
	TraefikRedirectURI = new(url.URL)

	v := viper.New()

	v.SetConfigFile(cfgFile)

	setDefaultValue(v)

	if err := v.ReadInConfig(); err != nil {
		zapx.Error("read config failed", zap.Error(err))
		return err
	}

	if err := v.Unmarshal(Config); err != nil {
		zapx.Error("unmarshal config file failed", zap.Error(err))
		return err
	}

	redirectURI, err := url.Parse(Config.Application.TraefikRedirectUri)
	if err != nil {
		zapx.Error("parse traefik redirect uri failed", zap.Error(err), zap.String("redirectURI", Config.Application.TraefikRedirectUri))
		return err
	}

	TraefikRedirectURI = redirectURI
	SessionOptions = sessions.Options{
		MaxAge:   constants.SessionMaxAgeSeconds,
		Domain:   Config.Application.SessionDomain,
		Secure:   true,
		HttpOnly: true,
		Path:     Config.Application.SessionPath,
	}
	return nil
}

func setDefaultValue(v *viper.Viper) {
	v.SetDefault("application.grpc_port", 30000)
}
