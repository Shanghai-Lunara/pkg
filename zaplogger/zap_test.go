package zaplogger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
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

func Test_getLevelFromEnv(t *testing.T) {
	err := os.Setenv(EnvZapLevel, "-1")
	assert.NoError(t, err)

	lv := getLevelFromEnv()
	assert.Equal(t, zap.DebugLevel, lv)
	assert.NotEqual(t, zap.InfoLevel, lv)

	err = os.Setenv(EnvZapLevel, "2")
	assert.NoError(t, err)
	lv = getLevelFromEnv()
	assert.Equal(t, zap.ErrorLevel, lv)
	assert.NotEqual(t, zap.InfoLevel, lv)
}
