package mongo

import "labix.org/v2/mgo"

// set env name existing mongo server
// future or home
var environment = "future"
var mongoDbName = "testa"

func mongoInit() (*mgo.Database, *mgo.Session) {
	var mongoIp string

	switch environment {
	case "future":
		mongoIp = "172.17.0.2"
	case "home":
		mongoIp = "172.17.0.19"
	}
	session, err := mgo.Dial(mongoIp)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(mongoDbName)

	return db, session
}
