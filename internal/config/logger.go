package config

import "go.uber.org/zap"

var Logger *zap.Logger = zap.NewNop()

func InitLogger(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Logger = zl
	return nil
}
