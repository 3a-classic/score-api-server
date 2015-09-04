package mongo

import (
	//	"encoding/json"
	"fmt"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"reflect"
)

type Player struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string
	Apply    string
	Editable bool
	Score    bson.M
}

type Team struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Member bson.M
	Team   string
}

type Field struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Hole           string
	DrivingContest bool
	Ignore         string
	NearPin        bool
	Par            int
	Yard           int
}

//type Member struct {
//	member string
//}
//
var (
	mongoIp = "172.17.0.2"
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

func GetAllPlayerCol() (*[]Player, error) {
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

func GetAllFieldCol() (*[]Field, error) {
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

func GetAllTeamCol() (*[]Team, error) {
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
