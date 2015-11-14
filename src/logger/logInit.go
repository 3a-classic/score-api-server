package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/weekface/mgorus"
)

var (
	mongoLog = logrus.New()
	slackLog = logrus.New()
	conf     *Config
)

const (
	Fatal    = "fatal"
	Error    = "error"
	Info     = "info"
	Debug    = "debug"
	ErrMsg   = "Error Message"
	TraceMsg = "Stack Trace"
)

type Config struct {
	Mongo struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Database string `toml:"database"`
	} `toml:"mongo"`
}

func init() {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(path.Join(dir, "../config/config.tml"), &conf); err != nil {
		panic(err)
	}

	mongoLog.Formatter = new(logrus.JSONFormatter)
	slackLog.Formatter = new(logrus.JSONFormatter)

	mongoLog.Out = os.Stderr
	slackLog.Out = os.Stderr

	mongoLog.Level = logrus.DebugLevel
	slackLog.Level = logrus.ErrorLevel

	slackLog.Hooks.Add(&slackrus.SlackrusHook{
		HookURL:        "https://hooks.slack.com/services/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",
		AcceptedLevels: slackrus.LevelThreshold(logrus.ErrorLevel),
		Channel:        "#3a-classic",
		Username:       "3a-classic-error-log",
	})

	hooker, err := mgorus.NewHooker(conf.Mongo.Host, conf.Mongo.Database, "log")
	if err == nil {
		mongoLog.Hooks.Add(hooker)
	}

}

func Output(field logrus.Fields, msg string, level string) {

	switch level {
	case Fatal:
		mongoLog.WithFields(field).Fatal(msg)
		slackLog.WithFields(field).Fatal(msg)
	case Error:
		mongoLog.WithFields(field).Error(msg)
		slackLog.WithFields(field).Error(msg)
	case Info:
		mongoLog.WithFields(field).Info(msg)
		slackLog.WithFields(field).Info(msg)
	case Debug:
		mongoLog.WithFields(field).Debug(msg)
		slackLog.WithFields(field).Debug(msg)
	}
}

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s\n", file, line, f.Name())
}
