package page

import (
	//	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"reflect"
)

type Index struct {
  team []string
	length string
}

type LeaderBoard struct {
	ranking  map[]
}

type ScoreEntrySheet struct {
 team: string
  hole: int
  member: map[]
  par: int
  yard: int
  stroke: map[]
  putt  map[]
  excnt int
}

type ScoreViewSheet struct {
  team string
  member map[]
  applay map[]
  hole: map[]
}

var (
	mongoIp = "172.17.0.19"
)

func mongoInit(col string) (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial(mongoIp)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("3a-test").C(col)
	return c, session

}

func GetAllPlayerJson() (*[]Player, error) {
	collectionName := "player"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Player{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}

func GetAllFieldJson() (*[]Field, error) {
	collectionName := "field"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Field{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	//	fmt.Println(result)
	return &result, nil
}

func GetAllTeamJson() (*[]Team, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}
