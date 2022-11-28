package zaplogger

import (
	"go.uber.org/zap"
	"testing"
)

func Test_Development(t *testing.T) {

	zcfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         EncodingJson,
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	Reset(zcfg)

	Sugar().Debug("1111")
	Sugar().Warn("2222")
	Sugar().Info("3333")
}
