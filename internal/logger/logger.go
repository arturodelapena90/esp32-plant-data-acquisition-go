package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// initializes global sugared logger
func Init() error {
	zl, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	Log = zl.Sugar()
	return nil
}

// flushes  logger before exit
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
