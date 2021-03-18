package zaplogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	var err error
	logger, err = NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}
	sugar = logger.Sugar()
}

func Sync() {
	_ = logger.Sync()
}

func Sugar() *zap.SugaredLogger {
	return sugar
}

func Logger() *zap.Logger {
	return logger
}

const EnvZapEncoding = "ENV_ZAP_ENCODING"

const (
	EncodingJson    = "json"
	EncodingConsole = "console"
)

// NewProductionConfig is a reasonable production logging configuration.
// Logging is enabled at InfoLevel and above.
//
// It uses a JSON encoder, writes to standard error, and enables sampling.
// Stacktraces are automatically included on logs of ErrorLevel and above.
func NewProductionConfig() zap.Config {
	encoding := os.Getenv(EnvZapEncoding)
	switch encoding {
	case "", EncodingConsole:
		encoding = EncodingConsole
	case EncodingJson:
		encoding = EncodingJson
	default:
		encoding = EncodingConsole
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         encoding,
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// NewProductionEncoderConfig returns an opinionated EncoderConfig for
// production environments.
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// EpochTimeEncoder serializes a time.Time to a floating-point number of seconds
// since the Unix epoch.
func EpochTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	//nanos := t.UnixNano()
	//sec := float64(nanos) / float64(time.Second)
	//enc.AppendFloat64(sec)
	enc.AppendString(t.Format("2006-01-02T15:04:05.000000Z"))
}
