package mongo

import (
	//	"encoding/json"
	"fmt"
	//	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	//	"reflect"
)

//Children

type Score struct {
	hole1 int
	hole2 int
	hole3 int
	hole4 int
}

type UserScore struct {
	score int
	name  string
	hole  int
}

type Member struct {
	member1 string
	member2 string
	member3 string
	member4 string
	length  int
}

type Apply struct {
	member1 int
	member2 int
	member3 int
	member4 int
}

type Hole struct {
	hole  int
	par   int
	yard  int
	score Score
}

//type Stroke struct {
//  hole1 int
//  hole2 int
//  hole3 int
//  hole4 int
//}
//
//type Putt struct {
//  hole1 int
//  hole2 int
//  hole3 int
//  hole4 int
//}

// Parents

type Index struct {
	team   []string
	length int
}

type LeaderBoard struct {
	ranking []*UserScore
}

type ScoreEntrySheet struct {
	team   string
	hole   int
	member Member
	par    int
	yard   int
	stroke Score
	putt   Score
	excnt  int
}

type ScoreViewSheet struct {
	team   string
	member Member
	applay Apply
	hole   []*Hole
}

//var (
//	mongoIp = "172.17.0.19"
//)

//func mongoInit(col string) (*mgo.Collection, *mgo.Session) {
//	session, err := mgo.Dial(mongoIp)
//	if err != nil {
//		panic(err)
//	}
//
//	session.SetMode(mgo.Monotonic, true)
//	c := session.DB("3a-test").C(col)
//	return c, session
//
//}

func GetIndexPageData() (*Index, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}

	teamArr := make([]string, 0, 10)
	for i := 0; i < len(result); i++ {
		teamArr = append(teamArr, result[i].Team)
	}

	idx := &Index{teamArr, len(result)}

	fmt.Println(idx)
	fmt.Println(result)
	fmt.Println(len(result))
	//	return &result, nil
	return idx, nil
}

func GetLeaderBoardPageData() (*LeaderBoard, error) {
	collectionName := "field"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Field{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	//	fmt.Println(result)
	//	return &result, nil
	return nil, nil
}

func GetScoreEntrySheetPageData() (*ScoreEntrySheet, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
	//	return &result, nil
	return nil, nil
}

func GetScoreViewSheetPageData() (*ScoreViewSheet, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
	//	return &result, nil
	return nil, nil
}
