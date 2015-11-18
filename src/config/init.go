package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"log"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	Conf *Config
)

type Mongo struct {
	Host          string `toml:"host"`
	Port          string `toml:"port"`
	Database      string `toml:"database"`
	LogCollection string `toml:"logcollection"`
}

type Slack struct {
	HookUrl  string `toml:"hookurl"`
	Username string `toml:"username"`
	Channel  string `toml:"channel"`
}

type GitRemote struct {
	Service string `toml:"service"`
	Url     string `toml:"url"`
}

type Config struct {
	Slack     Slack     `toml:"slack"`
	Mongo     Mongo     `toml:"mongo"`
	GitRemote GitRemote `toml:"gitremote"`
}

const (
	datetimeLocation = "Asia/Tokyo"
	DatetimeFormat   = "2006/01/02 15:04:05 MST"

	templateFilePath = "../src/config/config.tmpl"
	configFilePath   = "../src/config/config.tml"
)

func init() {

	setLocalTime()
	setConfig()

}

func setLocalTime() {

	loc, err := time.LoadLocation(datetimeLocation)
	if err != nil {
		loc = time.FixedZone(datetimeLocation, 9*60*60)
	}
	time.Local = loc
}

func checkTemplateFile(path string) {

}

func setConfig() {

	// get current dir
	pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	templateFileAbsPath := path.Join(pwd, templateFilePath)
	configFileAbsPath := path.Join(pwd, configFilePath)

	// check if file exists
	if _, err := os.Stat(templateFileAbsPath); os.IsNotExist(err) {
		log.Printf("you have to set template file at : %s", templateFileAbsPath)
		panic(err)
	}

	// read template file
	tmplString, err := ioutil.ReadFile(templateFileAbsPath)
	if err != nil {
		panic(err)
	}

	// parse tempalate file
	tmpl, err := template.New("pro").Funcs(funcMap()).Parse(string(tmplString))
	if err != nil {
		panic(err)
	}

	// create or open config file
	configToml, err := os.Create(configFileAbsPath)
	if err != nil {
		panic(err)
	}

	// get all env as map
	envMap := getAllEnv()

	// convert template to config file and write
	err = tmpl.Execute(configToml, envMap)
	if err != nil {
		panic(err)
	}

	// raed config file and set struct
	if _, err := toml.DecodeFile(configFileAbsPath, &Conf); err != nil {
		panic(err)
	}

}

func funcMap() template.FuncMap {
	funcMap := template.FuncMap{
		"_": func(i interface{}, d ...string) (interface{}, error) {
			if i == nil {
				if len(d) > 0 {
					return d[0], nil
				} else {
					return "", fmt.Errorf("nil value with no default")
				}

			}
			return i, nil
		},
	}
	return funcMap
}
