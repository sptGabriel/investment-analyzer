package logging

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger returns a new logger that writes to provided writer.
// If the provided writer implements zapcore.WriterSyncer, you need to call
// (*zap.Logger).Sync() before exiting.
func New(w io.Writer) *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				MessageKey:  "message",
				TimeKey:     "timestamp",
				LevelKey:    "level",
				EncodeLevel: zapcore.LowercaseLevelEncoder,
				EncodeTime:  zapcore.RFC3339NanoTimeEncoder,
			}),
			zapcore.AddSync(w),
			zapcore.InfoLevel,
		),
	)
}
