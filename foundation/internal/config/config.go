package config

import (
	"embed"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/fs"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
)

//go:embed default-config.yml
var embedConfig embed.FS

func GetConfig(logger *zap.Logger) *koanf.Koanf {
	k := koanf.New(".")
	err := k.Load(fs.Provider(embedConfig, "default-config.yml"), yaml.Parser())
	if err != nil {
		logger.Panic("Failed to load default config", zap.Error(err))
	}

	err = k.Load(env.Provider(".", env.Opt{
		Prefix: "APP_",
		TransformFunc: func(k, v string) (string, any) {
			k = strings.ToLower(strings.TrimPrefix(k, "APP_"))
			k = strings.ReplaceAll(k, "__", "-")
			k = strings.ReplaceAll(k, "_", ".")
			return k, v
		},
	}), nil)
	if err != nil {
		logger.Panic("Failed to load env config", zap.Error(err))
	}

	return k
}
