package mongo

import "labix.org/v2/mgo/bson"

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

func GetAllPlayerCol() []Player {
	db, session := mongoInit()
	col := db.C("player")
	defer session.Close()
	players := []Player{}
	err := col.Find(nil).All(&players)
	if err != nil {
		panic(err)
	}
	return players
}

func GetAllFieldCol() []Field {
	db, session := mongoInit()
	col := db.C("field")
	defer session.Close()
	fields := []Field{}
	err := col.Find(nil).All(&fields)
	if err != nil {
		panic(err)
	}
	return fields
}

func GetAllTeamCol() []Team {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()
	teams := []Team{}
	err := col.Find(nil).All(&teams)
	if err != nil {
		panic(err)
	}
	return teams
}

func GetOnePlayerByQuery(query bson.M) Player {
	db, session := mongoInit()
	col := db.C("player")
	defer session.Close()
	player := Player{}
	err := col.Find(query).One(&player)
	if err != nil {
		panic(err)
	}
	return player
}

func GetOneFieldByQuery(query bson.M) Field {
	db, session := mongoInit()
	col := db.C("field")
	defer session.Close()
	field := Field{}
	err := col.Find(query).One(&field)
	if err != nil {
		panic(err)
	}
	return field
}

func GetOneTeamByQuery(query bson.M) Team {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()
	team := Team{}
	err := col.Find(query).One(&team)
	if err != nil {
		panic(err)
	}
	return team
}

func GetPlayersDataInTheTeam(teamName string) []Player {
	_, session := mongoInit()
	defer session.Close()
	player := Player{}
	players := make([]Player, 4)

	team := GetOneTeamByQuery(bson.M{"team": teamName})
	for i := 0; i < len(team.Member); i++ {
		err := session.FindRef(&team.Member[i].Player).One(&player)
		if err != nil {
			panic(err)
		}
		players[i] = player
	}
	return players
}
