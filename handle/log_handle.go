package handle

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:   "message",
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		EncodeLevel:  CustomLevelEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	})
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func logErrorWriter() (zapcore.WriteSyncer, error) {
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	), nil
}

func logInfoWriter() (zapcore.WriteSyncer, error) {
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
	), nil
}

// NewLog to console by level log
func NewLog() (*zap.Logger, error) {
	highWriteSyncer, errorWriter := logErrorWriter()
	if errorWriter != nil {
		return nil, errorWriter
	}
	lowWriteSyncer, errorInfoWriter := logInfoWriter()
	if errorInfoWriter != nil {
		return nil, errorInfoWriter
	}

	encoder := getEncoder()

	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})

	lowCore := zapcore.NewCore(encoder, lowWriteSyncer, lowPriority)
	highCore := zapcore.NewCore(encoder, highWriteSyncer, highPriority)

	logger := zap.New(zapcore.NewTee(lowCore, highCore), zap.AddCaller())
	return logger, nil
}
