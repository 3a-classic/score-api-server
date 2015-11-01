package mongo

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"labix.org/v2/mgo"
)

// set env name existing mongo server
// future or home

var (
	conf    *Config
	players []PlayerCol
	fields  []FieldCol
	teams   []TeamCol
	threads []ThreadCol
	err     error
)

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	if _, err := toml.DecodeFile(path.Join(dir, "../config/config.tml"), &conf); err != nil {
		panic(err)
	}
	players = GetAllPlayerCol()
	fields = GetAllFieldCol()
	teams = GetAllTeamCol()
	threads = GetAllThreadCol()

	const location = "Asia/Tokyo"

	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
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
