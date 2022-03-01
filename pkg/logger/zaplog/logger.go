package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(name, level, env, version string, out ...string) (*zap.Logger, error) {
	var config zap.Config
	switch env {
	case "prod":
		config = zap.NewProductionConfig()
	default:
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	if len(out) == 0 {
		out = append(out, "stdout")
	}

	config.OutputPaths = out

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	err := config.Level.UnmarshalText([]byte(level))
	if err != nil || len(level) == 0 {
		config.Level.SetLevel(zap.DebugLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("name", name),
	)

	return logger, nil
}

func MustNew(name, level, env, version string, out ...string) *zap.Logger {
	logger, err := New(name, level, env, version, out...)
	if err != nil {
		panic(err)
	}

	return logger
}

func FromEnv(name string, out ...string) *zap.Logger {
	env := os.Getenv("ENV")
	lvl := os.Getenv("LVL")
	ver := os.Getenv("VER")

	return MustNew(name, lvl, env, ver, out...)
}
