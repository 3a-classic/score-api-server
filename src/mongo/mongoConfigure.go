package mongo

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"labix.org/v2/mgo"
)

var (
	conf           *Config
	users          map[string]UserCol
	players        map[string]PlayerCol
	fields         map[int]FieldCol
	teams          map[string]TeamCol
	threads        map[string]ThreadCol
	datetimeFormat string
	excnt          map[string]map[int]int
	err            error
	ThreadChan     chan *Thread
	ErrChan        chan string
	FinChan        chan bool
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(path.Join(dir, "../config/config.tml"), &conf); err != nil {
		panic(err)
	}
	initMap()
	SetAllUserCol()
	SetAllPlayerCol()
	SetAllFieldCol()
	SetAllTeamCol()
	SetAllThreadCol()
	setLocalTime()
	initExcnt()
	datetimeFormat = "2006/01/02 15:04:05 MST"

	ThreadChan = make(chan *Thread, 2)
}

func mongoInit() (*mgo.Database, *mgo.Session) {

	session, err := mgo.Dial(conf.Mongo.Host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(conf.Mongo.Database)

	return db, session
}

func setLocalTime() {
	const location = "Asia/Tokyo"

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func initMap() {
	users = map[string]UserCol{}
	players = map[string]PlayerCol{}
	fields = map[int]FieldCol{}
	teams = map[string]TeamCol{}
	threads = map[string]ThreadCol{}
}

func initExcnt() {
	excnt = map[string]map[int]int{}
	for _, team := range teams {
		excnt[team.Name] = map[int]int{}
		for _, field := range fields {
			excnt[team.Name][field.Hole] = 0
		}
	}
}
