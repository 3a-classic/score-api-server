package logger

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/Sirupsen/logrus"
)

func PutFata(err error, trace string, msg string, obj interface{}) {

	objType := SprintfType(obj)
	field := &logrus.Fields{
		ErrMsg:   fmt.Errorf("%v", err),
		TraceMsg: trace,
		objType:  fmt.Sprintf("%+v\n", obj),
	}
	log.WithFields(*field).Fatal(msg)
}

func PutErr(err error, trace string, msg string, obj interface{}) {

	fmt.Println("outErr")
	objType := SprintfType(obj)
	field := &logrus.Fields{
		ErrMsg:   Errorf(err),
		TraceMsg: trace,
		objType:  Sprintf(obj),
	}
	log.WithFields(*field).Error(msg)
	fmt.Println("outErr end")
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

func Errorf(err error) string {
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

func Trace() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d %s\n", file, line, f.Name())
}
