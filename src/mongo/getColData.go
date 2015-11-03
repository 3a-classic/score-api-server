package mongo

import (
	"labix.org/v2/mgo/bson"
)

//for debug
//return all collection data
//
// this method is not userd web app
// becouse return is interface(map)
// and map is not defined  order
func GetAllColData(collectionName string) (*[]interface{}, error) {
	db, session := mongoInit()
	col := db.C(collectionName)
	defer session.Close()

	var results []interface{}
	err := col.Find(nil).All(&results)
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func SetAllUserCol() {
	db, session := mongoInit()
	col := db.C("user")
	defer session.Close()
	usersCol := []UserCol{}
	if err = col.Find(nil).All(&usersCol); err != nil {
		panic(err)
	}

	for _, userCol := range usersCol {
		users[userCol.UserId] = userCol
	}
}

func SetUserCol(userIds []string) {
	db, session := mongoInit()
	col := db.C("user")
	defer session.Close()

	for _, userId := range userIds {

		userCol := UserCol{}
		findQuery := bson.M{"userId": userId}

		if err = col.Find(findQuery).One(&userCol); err != nil {
			panic(err)
		}

		users[userId] = userCol
	}
}

func SetAllPlayerCol() {
	db, session := mongoInit()
	col := db.C("player")
	defer session.Close()
	playersCol := []PlayerCol{}
	if err = col.Find(nil).All(&playersCol); err != nil {
		panic(err)
	}

	for _, playerCol := range playersCol {
		players[playerCol.UserId] = playerCol
	}
}

func SetPlayerCol(userIds []string) {
	db, session := mongoInit()
	col := db.C("player")
	defer session.Close()

	for _, userId := range userIds {

		playerCol := PlayerCol{}
		findQuery := bson.M{"userrId": userId}

		if err = col.Find(findQuery).One(&playerCol); err != nil {
			panic(err)
		}

		players[userId] = playerCol
	}
}

func SetAllFieldCol() {
	db, session := mongoInit()
	col := db.C("field")
	defer session.Close()
	fieldsCol := []FieldCol{}
	if err = col.Find(nil).All(&fieldsCol); err != nil {
		panic(err)
	}

	for _, fieldCol := range fieldsCol {
		fields[fieldCol.Hole] = fieldCol
	}
}

func SetFieldCol(hole int) {
	db, session := mongoInit()
	col := db.C("field")
	defer session.Close()

	fieldCol := FieldCol{}
	findQuery := bson.M{"hole": hole}

	if err = col.Find(findQuery).One(&fieldCol); err != nil {
		panic(err)
	}

	fields[hole] = fieldCol
}

func SetAllTeamCol() {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()
	teamsCol := []TeamCol{}
	if err = col.Find(nil).All(&teamsCol); err != nil {
		panic(err)
	}

	for _, teamCol := range teamsCol {
		teams[teamCol.Name] = teamCol
	}
}

func SetTeamCol(teamName string) {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()

	teamCol := TeamCol{}
	findQuery := bson.M{"name": teamName}

	if err = col.Find(findQuery).One(&teamCol); err != nil {
		panic(err)
	}

	teams[teamName] = teamCol
}

func SetAllThreadCol() {
	db, session := mongoInit()
	col := db.C("thread")
	defer session.Close()
	threadsCol := []ThreadCol{}
	if err = col.Find(nil).All(&threadsCol); err != nil {
		panic(err)
	}

	for _, threadCol := range threadsCol {
		threads[threadCol.ThreadId] = threadCol
	}
}

func SetThreadCol(threadId string) {
	db, session := mongoInit()
	col := db.C("thread")
	defer session.Close()

	threadCol := ThreadCol{}
	findQuery := bson.M{"threadId": threadId}

	if err = col.Find(findQuery).One(&threadCol); err != nil {
		panic(err)
	}

	threads[threadId] = threadCol
}
