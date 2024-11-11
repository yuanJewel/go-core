package task

import (
	"github.com/sirupsen/logrus"
	"github.com/yuanJewel/go-core/logger"
)

type wrapper struct {
	level logrus.Level
}

func (w *wrapper) Print(v ...interface{}) {
	logger.Log.Log(w.level, v...)
}

func (w *wrapper) Printf(format string, v ...interface{}) {
	logger.Log.Logf(w.level, format, v...)
}

func (w *wrapper) Println(v ...interface{}) {
	logger.Log.Logln(w.level, v...)
}

func (w *wrapper) Fatal(v ...interface{}) {
	logger.Log.Fatal(v...)
}

func (w *wrapper) Fatalf(format string, v ...interface{}) {
	logger.Log.Fatalf(format, v...)
}

func (w *wrapper) Fatalln(v ...interface{}) {
	logger.Log.Fatalln(v...)
}

func (w *wrapper) Panic(v ...interface{}) {
	logger.Log.Panic(v...)
}

func (w *wrapper) Panicf(format string, v ...interface{}) {
	logger.Log.Panicf(format, v...)
}

func (w *wrapper) Panicln(v ...interface{}) {
	logger.Log.Panicln(v...)
}
