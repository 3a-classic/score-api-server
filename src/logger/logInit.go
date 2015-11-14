package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
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
	Fatal = "fatal"
	Error = "error"
	Info  = "info"
	Debug = "debug"

	ErrMsg   = "Error Message"
	TraceMsg = "Stack Trace"

	E_Nil         = "This variable is nil"
	E_WrongData   = "This is not correct data"
	E_TooManyData = "There are too many data"
	E_MakeHash    = "Can not make hash string"

	E_M_FindEntireCol  = "Can not find entire colection : mongo"
	E_M_FindCol        = "Can not find colection : mongo"
	E_M_Upsert         = "Can not upsert data : mongo"
	E_M_Insert         = "Can not insert data : mongo"
	E_M_Update         = "Can not update data : mongo"
	E_M_RegisterThread = "Can not register thread score : mongo"
	//	E_M_RegisterUser    = "Can not register user"
	E_M_SearchPhotoTask = "Can not search picture task : mongo"

	I_M_GetPage     = "Get page data : mongo"
	I_M_PostPage    = "Post page data : mongo"
	I_M_RegisterCol = "Register collection data : mongo"

	E_R_PostPage    = "Can not post page data : route"
	E_R_RegisterCol = "Can not register collection data : route"
	E_R_Upsert      = "Can not upsert data : route"
	E_R_WriteJSON   = "Can not write JSON : route"
	E_R_PingMsg     = "Can not ping message : route"
	E_R_Upgrader    = "Can not upgrader webdocket : route"
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

func PutFata(err error, trace string, msg string, obj interface{}) {

	objType := fmt.Sprintf("%+v\n", reflect.ValueOf(obj).Type())
	field := &logrus.Fields{
		ErrMsg:   fmt.Errorf("%v", err),
		TraceMsg: trace,
		objType:  fmt.Sprintf("%+v\n", obj),
	}
	mongoLog.WithFields(*field).Fatal(msg)
	slackLog.WithFields(*field).Fatal(msg)
}

func PutErr(err error, trace string, msg string, obj interface{}) {

	objType := fmt.Sprintf("%+v\n", reflect.ValueOf(obj).Type())
	field := &logrus.Fields{
		ErrMsg:   Errorf(err),
		TraceMsg: trace,
		objType:  Sprintf(obj),
	}
	mongoLog.WithFields(*field).Error(msg)
	slackLog.WithFields(*field).Error(msg)
}

func Output(field logrus.Fields, msg string, level string) {

	switch level {
	case Info:
		mongoLog.WithFields(field).Info(msg)
		slackLog.WithFields(field).Info(msg)
	case Debug:
		mongoLog.WithFields(field).Debug(msg)
		slackLog.WithFields(field).Debug(msg)
	}
}

func Errorf(err error) string {
	return fmt.Sprintf("%v", err)
}

func Sprintf(obj interface{}) string {
	return fmt.Sprintf("%+v\n", obj)
}

func Trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s\n", file, line, f.Name())
}
