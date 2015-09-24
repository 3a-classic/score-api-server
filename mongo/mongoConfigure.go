package mongo

import (
	"github.com/BurntSushi/toml"
	"labix.org/v2/mgo"
)

// set env name existing mongo server
// future or home

var (
	conf    *Config
	players []Player
	fields  []Field
	teams   []Team
)

func init() {
	if _, err := toml.DecodeFile("config/config.tml", &conf); err != nil {
		panic(err)
	}
	players = GetAllPlayerCol()
	fields = GetAllFieldCol()
	teams = GetAllTeamCol()
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
