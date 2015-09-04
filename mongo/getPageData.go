package mongo

import (
	//	"encoding/json"
	"fmt"
	//	"labix.org/v2/mgo"
	"sort"

	"labix.org/v2/mgo/bson"
	//	"reflect"
)

//Children

type TeamScore struct {
	Hole1 int
	Hole2 int
	Hole3 int
	Hole4 int
}

type UserScore struct {
	Score int
	Name  string
	Hole  int
}

type TeamMember struct {
	Member1 string
	Member2 string
	Member3 string
	Member4 string
	Length  int
}

type Apply struct {
	Member1 int
	Member2 int
	Member3 int
	Member4 int
}

type Hole struct {
	Hole  int
	Par   int
	Yard  int
	Score TeamScore
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
	Team   []string
	Length int
}

type LeaderBoard struct {
	Ranking []*UserScore
}

type ScoreEntrySheet struct {
	Team   string
	Hole   int
	Member TeamMember
	Par    int
	Yard   int
	Stroke TeamScore
	Putt   TeamScore
	Excnt  int
}

type ScoreViewSheet struct {
	Team   string
	Member Member
	Applay Apply
	Hole   []*Hole
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

	teamArr := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		teamArr[i] = result[i].Team
	}

	idx := &Index{teamArr, len(result)}

	//	return &result, nil
	return idx, nil
}

func GetLeaderBoardPageData() (*LeaderBoard, error) {
	collectionName := "player"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Player{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(result))

	var leadersBoard []UserScore
	var userScore UserScore

	//	fmt.Println(result[1].Score[8])
	//	fmt.Println(reflect.ValueOf(result[0].Score[0]["total"]).Type())
	for i := 0; i < len(result); i++ {
		totalScore := 0
		passedHoleCnt := 0
		for j := 0; j < len(result[i].Score); j++ {
			//			fmt.Println(result[i].Score[j]["hole"])
			if result[i].Score[j]["total"] != 0 {
				//				fmt.Println(reflect.ValueOf(result[0].Score[0]["total"]).Type())
				//				fmt.Println(reflect.ValueOf(totalScore).Type())

				totalScore += result[i].Score[j]["total"].(int)
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore,
			Name:  result[i].Name,
			Hole:  passedHoleCnt,
		}
		leadersBoard = append(leadersBoard, userScore)
	}
	fmt.Println(leadersBoard)

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

// https://gist.github.com/ikbear/4038654
// sortのメソッドパクってきた
type sortedMap struct {
	m map[string]int
	s []string
}

func (sm *sortedMap) Len() int {
	return len(sm.m)
}

func (sm *sortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *sortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func sortedKeys(m map[string]int) []string {
	sm := new(sortedMap)
	sm.m = m
	sm.s = make([]string, len(m))
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}
