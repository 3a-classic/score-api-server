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

func GetAllPlayerCol() []PlayerCol {
	db, session := mongoInit()
	col := db.C("player")
	defer session.Close()
	players := []PlayerCol{}
	err := col.Find(nil).All(&players)
	if err != nil {
		panic(err)
	}
	return players
}

func GetAllFieldCol() []FieldCol {
	db, session := mongoInit()
	col := db.C("field")
	defer session.Close()
	fields := []FieldCol{}
	err := col.Find(nil).All(&fields)
	if err != nil {
		panic(err)
	}
	return fields
}

func GetAllTeamCol() []TeamCol {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()
	teams := []TeamCol{}
	err := col.Find(nil).All(&teams)
	if err != nil {
		panic(err)
	}
	return teams
}

//func GetAllThreadCol() []TeamCol {
//	db, session := mongoInit()
//	col := db.C("team")
//	defer session.Close()
//	teams := []Team{}
//	err := col.Find(nil).All(&teams)
//	if err != nil {
//		panic(err)
//	}
//	return teams
//}

func GetPlayersDataInTheTeam(teamName string) []PlayerCol {
	db, session := mongoInit()
	col := db.C("team")
	defer session.Close()
	player := PlayerCol{}
	team := TeamCol{}

	if err := col.Find(bson.M{"team": teamName}).One(&team); err != nil {
		panic(err)
	}
	players := make([]PlayerCol, len(team.Member))
	for i, teamPlayer := range team.Member {
		if err := session.FindRef(&teamPlayer.Player).One(&player); err != nil {
			panic(err)
		}
		players[i] = player
	}
	return players
}
