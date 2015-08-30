package mongo

import (
	//	"encoding/json"
	//	"fmt"
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
