package logger

import (
	c "config"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
	"github.com/weekface/mgorus"
)

var (
	log = logrus.New()
)

func init() {

	log.Formatter = new(logrus.JSONFormatter)
	log.Level = logrus.DebugLevel
	log.Out = os.Stderr
	hooker, err := mgorus.NewHooker(c.Conf.Mongo.Host, c.Conf.Mongo.Database, c.Conf.Mongo.LogCollection)
	if err == nil {
		log.Hooks.Add(hooker)
	}

	if len(c.Conf.Slack.HookUrl) != 0 {
		log.Hooks.Add(&slackrus.SlackrusHook{
			HookURL:        c.Conf.Slack.HookUrl,
			AcceptedLevels: slackrus.LevelThreshold(logrus.ErrorLevel),
			Channel:        c.Conf.Slack.Channel,
			Username:       c.Conf.Slack.Username,
		})
	}
}
