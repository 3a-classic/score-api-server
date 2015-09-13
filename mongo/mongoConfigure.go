package mongo

import (
	"github.com/BurntSushi/toml"

	"labix.org/v2/mgo"
)

// set env name existing mongo server
// future or home

func mongoInit() (*mgo.Database, *mgo.Session) {

	var conf *Config
	_, err := toml.DecodeFile("config/config.tml", &conf)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial(conf.Mongo.Host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(conf.Mongo.Database)

	return db, session
}
