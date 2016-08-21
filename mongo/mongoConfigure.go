package mongo

import (
	c "github.com/3a-classic/score-api-server/config"

	"labix.org/v2/mgo"
)

var (
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
	initColMap()
	SetAllUserCol()
	SetAllPlayerCol()
	SetAllFieldCol()
	SetAllTeamCol()
	SetAllThreadCol()

	initExcntMap()
	ThreadChan = make(chan *Thread, 2)
}

func mongoConn() (*mgo.Database, *mgo.Session) {

	session, err := mgo.Dial(c.Conf.Mongo.Host)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(c.Conf.Mongo.Database)

	return db, session
}

func initColMap() {
	users = map[string]UserCol{}
	players = map[string]PlayerCol{}
	fields = map[int]FieldCol{}
	teams = map[string]TeamCol{}
	threads = map[string]ThreadCol{}
}

func initExcntMap() {
	excnt = map[string]map[int]int{}
	for _, team := range teams {
		excnt[team.Name] = map[int]int{}
		for _, field := range fields {
			excnt[team.Name][field.Hole] = 0
		}
	}
}
