package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strconv"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Player struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Apply    int           `json:"apply"`
	Editable bool          `json:"editable"`
	Score    []struct {
		Hole   int `json:"hole"`
		Stroke int `json:"stroke"`
		Putt   int `json:"putt"`
		Total  int `json:"total"`
	} `json:"score"`
	Team mgo.DBRef
}

type Member struct {
	Player mgo.DBRef `json:"player"`
}

type Team struct {
	Id     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Member []Member      `json:"member"`
	Team   string
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
	mongoIp = "172.17.0.2"
	//	mongoIp = "172.17.0.19"
	dbName = "testa"
)

func mongoInit(col string) (*mgo.Collection, *mgo.Session) {
	session, err := mgo.Dial(mongoIp)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	c := session.DB(dbName).C(col)
	return c, session

}

const FILE_ROOT = "/var/www/3acses/post_data/sample_data"

func main() {
	err := upsertField()
	if err != nil {
		panic(err)
	}
	err = insertTeamAndPlayer()
	if err != nil {
		panic(err)
	}

}

func upsertField() error {

	collectionName := "field"
	fmt.Println(collectionName + "にデータを挿入します。")
	c, session := mongoInit(collectionName)
	defer session.Close()

	var fileNum = 18
	var createCnt = 0
	var updateCnt = 0
	for i := 1; i <= fileNum; i++ {
		fileName := collectionName + strconv.Itoa(i) + ".json"
		filePath := path.Join(FILE_ROOT, fileName)
		jsonString, err := ioutil.ReadFile(filePath) // ReadFileの戻り値は []byte
		if err != nil {
			return err
		}

		var field Field
		err = json.Unmarshal(jsonString, &field)
		if err != nil {
			return err
		}
		colQuerier := bson.M{"hole": field.Hole}

		//type ChangeInfo struct {
		//    Updated    int
		//    Removed    int
		//    UpsertedId interface{}
		//}
		//		field.Id = bson.NewObjectId()
		change, err := c.Upsert(colQuerier, field)
		if err != nil {
			return err
		}
		if change.Updated == 0 {
			createCnt += 1
		} else {
			updateCnt += 1
		}

	}

	fmt.Println(strconv.Itoa(createCnt) + "件作成されました。")
	fmt.Println(strconv.Itoa(updateCnt) + "件更新されました。")
	return nil
}

func insertTeamAndPlayer() error {
	playerNum := 20
	// playerNum / 4 の切り上げ
	//	teamNum := int(math.Ceil(float64(playerNum / 4)))
	playerNumPerTeam := 4
	collectionName1 := "team"
	collectionName2 := "player"
	fmt.Println(collectionName1 + "にデータを挿入します。")
	fmt.Println(collectionName2 + "にデータを挿入します。")
	c1, session1 := mongoInit(collectionName1)
	c2, session2 := mongoInit(collectionName2)
	defer session1.Close()
	defer session2.Close()

	var arrayNum int
	var teamObjectId bson.ObjectId
	//	playerArray := make(map[string]mgo.DBRef, 4)
	var OneBeforebyteOfA byte = 64
	alphabet := make([]byte, 1)
	alphabet[0] = OneBeforebyteOfA
	var team Team
	members := make([]Member, 4)
	for i := playerNum; i >= 0; i-- {

		var member Member
		if i%playerNumPerTeam == 0 {
			arrayNum = playerNumPerTeam
		} else {
			arrayNum = i % playerNumPerTeam
		}
		if (arrayNum == 4 && i != playerNum) || (i == 0) {
			//			team.Member = playerArray
			//			clearStruct(member)
			//			member.Player = new mgo.DBRef
			//	members = append(nil, member)
			alphabet[0] += 1
			team.Id = teamObjectId
			team.Team = string(alphabet)
			team.Member = members
			team.Definded = false
			c1.Insert(team)
			if i == 0 {
				break
			}
		}

		if arrayNum == 4 || i == playerNum {
			teamObjectId = bson.NewObjectId()
		}

		fileName := collectionName2 + strconv.Itoa(i) + ".json"
		filePath := path.Join(FILE_ROOT, fileName)
		jsonString, err := ioutil.ReadFile(filePath) // ReadFileの戻り値は []byte

		if err != nil {
			return err
		}
		var player Player
		err = json.Unmarshal(jsonString, &player)
		playerObjectId := bson.NewObjectId()
		member.Player = makeDBRef(playerObjectId, collectionName2)
		//		members = append(members, member)
		members[arrayNum-1] = member
		player.Id = playerObjectId
		player.Team = makeDBRef(teamObjectId, collectionName1)
		err = c2.Insert(player)
		if err != nil {
			return err
		}
	}
	fmt.Println(strconv.Itoa(playerNum) + "件作成されました。")
	return nil
}

func makeDBRef(objectId bson.ObjectId, collectionName string) mgo.DBRef {
	dbRef := mgo.DBRef{
		Collection: collectionName,
		Id:         objectId,
		Database:   dbName,
	}
	return dbRef
}

func clearStruct(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func p(str string) {
	fmt.Println(str)
}
