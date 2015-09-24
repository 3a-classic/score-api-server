package mongo

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

func PostApplyScoreData(teamName string, registeredApplyScore *PostApplyScore) (*Status, error) {

	fmt.Println("Team : " + teamName + " 申請スコアを登録します。")
	db, session := mongoInit()
	playerCol := db.C("player")
	defer session.Close()

	teamPlayers := GetPlayersDataInTheTeam(teamName)
	if teamPlayers[0].Apply != 0 {
		return &Status{"already registered"}, nil
	}
	for playerIndex, player := range teamPlayers {
		player.Apply = registeredApplyScore.Apply[playerIndex]

		query := bson.M{"_id": player.Id}
		if err := playerCol.Update(query, player); err != nil {
			return &Status{"failed"}, err
		}
	}

	//更新情報をGlobal変数に格納する
	players = GetAllPlayerCol()
	return &Status{"success"}, nil

}

func PostScoreViewSheetPageData(teamName string, definedTeam *PostDefinedTeam) (*Status, error) {

	fmt.Println("Team : " + teamName + "のデータを確定します。")

	targetTeam := Team{}
	for _, team := range teams {
		if team.Team == teamName {
			targetTeam = team
		}
	}

	if targetTeam.Defined {
		return &Status{"already defined"}, nil
	}

	targetTeam.Defined = true

	db, session := mongoInit()
	teamCol := db.C("team")
	defer session.Close()

	query := bson.M{"_id": targetTeam.Id}
	if err := teamCol.Update(query, targetTeam); err != nil {
		return &Status{"failed"}, err
	}

	//更新情報をGlobal変数に格納する
	teams = GetAllTeamCol()
	return &Status{"success"}, nil
}

func PostScoreEntrySheetPageData(teamName string, holeString string, updatedTeamScore *PostTeamScore) (*Status, error) {

	if len(holeString) == 0 {
		return &Status{"failed"}, errors.New("hole is not string")
	}
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	fmt.Println("Team : " + teamName + ", Hole : " + holeString + "にデータを挿入します。")

	db, session := mongoInit()
	playerCol := db.C("player")
	defer session.Close()

	teamPlayers := GetPlayersDataInTheTeam(teamName)

	for playerIndex, player := range teamPlayers {
		stroke, putt := updatedTeamScore.Stroke[playerIndex], updatedTeamScore.Putt[playerIndex]
		player.Score[holeIndex]["stroke"] = stroke
		player.Score[holeIndex]["putt"] = putt
		player.Score[holeIndex]["total"] = stroke + putt

		query := bson.M{"_id": player.Id}
		if err := playerCol.Update(query, player); err != nil {
			return &Status{"failed"}, err
		}
	}

	//更新情報をGlobal変数に格納する
	players = GetAllPlayerCol()
	return &Status{"success"}, nil
}
