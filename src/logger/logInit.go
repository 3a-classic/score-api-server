package logger

import (
	"os"
	"path"
	"path/filepath"

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
