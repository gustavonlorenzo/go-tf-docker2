package loggerator

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitializeLogger initializes a Zap logger with file output.
func InitializeLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Change the output to a file
	config.OutputPaths = []string{"logfile.log"} // specify the file path here

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
