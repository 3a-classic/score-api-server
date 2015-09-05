package mongo

import (
	//	"encoding/json"
	//	"fmt"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"reflect"
)

type Score struct {
	Hole   int `json:"hole"`
	Strole int `json:"stroke"`
	Putt   int `json:"putt"`
	Total  int `json:"total"`
}

type Member struct {
	Player mgo.DBRef `json:"player"`
}

type Player struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Apply    int           `json:"apply"`
	Editable bool          `json:"editable"`
	Score    []bson.M      `json:"score"`
	Team     mgo.DBRef     `json:"team"`
}

type Team struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Member []Member      `json:"member"`
	Team   string        `json:"team"`
}

type Field struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Hole           int           `json:"hole"`
	DrivingContest bool          `json:"drivingContest"`
	Ignore         bool          `json:"ignore"`
	Image          string        `json:"image"`
	NearPin        bool          `json:"nearPin"`
	Par            int           `json:"par"`
	Yard           int           `json:"yard"`
}

var (
	//	mongoIp = "172.17.0.2"
	mongoIp = "172.17.0.19"
)

func mongoInit(col string) (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial(mongoIp)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("testa").C(col)
	return c, session

}

func GetAllPlayerCol() (*[]Player, error) {
	collectionName := "player"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Player{}
	err := col.Find(nil).All(&result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}

func GetAllFieldCol() (*[]Field, error) {
	collectionName := "field"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Field{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}

func GetAllTeamCol() (*[]Team, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	return &result, nil
}
