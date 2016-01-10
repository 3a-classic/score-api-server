package logger

import (
	c "config"

	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
)

func PutFata(err error, trace map[string]string, msg string, obj interface{}) {

	objType := SprintfType(obj)
	field := &logrus.Fields{
		ErrMsg:             ErrToStr(err),
		TraceMsg:           TraceToStr(trace),
		objType:            Sprintf(obj),
		"GitRemoteCodeUrl": GitRemoteCodeUrl(trace),
	}
	log.WithFields(*field).Fatal(msg)
}

func PutErr(err error, trace map[string]string, msg string, obj interface{}) {

	objType := SprintfType(obj)
	field := &logrus.Fields{
		ErrMsg:             ErrToStr(err),
		TraceMsg:           TraceToStr(trace),
		objType:            Sprintf(obj),
		"GitRemoteCodeUrl": GitRemoteCodeUrl(trace),
	}
	log.WithFields(*field).Error(msg)
}

func PutInfo(msg string, obj1 interface{}, obj2 interface{}) {
	field := &logrus.Fields{
		SprintfType(obj1): Sprintf(obj1),
		SprintfType(obj2): Sprintf(obj2),
	}

	log.WithFields(*field).Info(msg)
}

func Output(field logrus.Fields, msg string, level string) {

	switch level {
	case Info:
		log.WithFields(field).Info(msg)
	case Debug:
		log.WithFields(field).Debug(msg)
	}
}

func TraceToStr(trace map[string]string) string {
	if trace != nil {
		return fmt.Sprintf("%s:%s %s\n", trace["file"], trace["line"], trace["name"])
	} else {
		return ""
	}

}

func ErrToStr(err error) string {
	if err != nil {
		return fmt.Sprintf("%v", err)
	} else {
		return ""
	}
}

func Sprintf(obj interface{}) string {
	if obj != nil {
		return fmt.Sprintf("%+v\n", obj)
	} else {
		return ""
	}
}

func SprintfType(obj interface{}) string {
	if obj != nil {
		return fmt.Sprintf("%+v\n", reflect.ValueOf(obj).Type())
	} else {
		return ""
	}
}

func Trace() map[string]string {
	trace := make(map[string]string)
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	trace["file"], trace["line"] = file, strconv.Itoa(line)
	trace["name"] = f.Name()
	return trace
}

func GitRemoteCodeUrl(trace map[string]string) string {

	gitService := c.Conf.GitRemote.Service
	gitUrl := c.Conf.GitRemote.Url
	gitCodePath := "blob/master"

	if len(gitService) == 0 || len(gitUrl) == 0 {
		return ""
	}
	if gitService != "github" && gitService != "gitlab" {
		return ""
	}

	// delete prefix and suffix "/"
	if strings.HasPrefix(gitUrl, "/") {
		gitUrl = gitUrl[1:]
	}
	if strings.HasSuffix(gitUrl, "/") {
		cutOffLastCharLen := len(gitUrl) - 1
		gitUrl = gitUrl[:cutOffLastCharLen]
	}

	u, err := url.Parse(gitUrl)
	if err != nil {
		panic(err)
	}

	u.RawPath, u.RawQuery, u.Fragment = "", "", ""
	gitUrlPathArray := strings.Split(u.Path, "/")
	repos := gitUrlPathArray[len(gitUrlPathArray)-1]
	traceChildPath := strings.SplitAfter(trace["file"], repos)[1]
	u.Path = filepath.Join(u.Path, gitCodePath, traceChildPath)
	u.Fragment = "L" + trace["line"]

	return u.String()
}
