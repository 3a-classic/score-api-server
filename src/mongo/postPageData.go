package mongo

import (
	"errors"
	"fmt"
	"strconv"

	"labix.org/v2/mgo/bson"
)

func PostLoginPageData(loginInfo *PostLogin) (*Status, error) {

	fmt.Println(loginInfo.Name + "さんがアクセスしました。")
	canLogin := false
	for _, player := range players {
		if player.Name == loginInfo.Name {
			canLogin = true
		}
	}

	if canLogin {
		fmt.Println(loginInfo.Name + "さんがログインしました。")
		return &Status{"success"}, nil
	}

	fmt.Println(loginInfo.Name + "さんがログインに失敗しました。")
	return &Status{"failed"}, nil

}

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

	targetTeam := TeamCol{}
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

	targetTeam := TeamCol{}
	for _, team := range teams {
		if team.Team == teamName {
			targetTeam = team
		}
	}

	if updatedTeamScore.Excnt != targetTeam.Excnt[holeIndex] {
		return &Status{"other updated"}, nil
	}

	fmt.Println("Team : " + teamName + ", Hole : " + holeString + "にデータを挿入します。")

	db, session := mongoInit()
	playerCol := db.C("player")
	teamCol := db.C("team")
	defer session.Close()

	targetTeam.Excnt[holeIndex] += 1

	teamQuery := bson.M{"_id": targetTeam.Id}
	if err := teamCol.Update(teamQuery, targetTeam); err != nil {
		return &Status{"failed update excnt"}, err
	}

	teamPlayers := GetPlayersDataInTheTeam(teamName)

	for playerIndex, player := range teamPlayers {
		total, putt := updatedTeamScore.Total[playerIndex], updatedTeamScore.Putt[playerIndex]
		player.Score[holeIndex]["putt"] = putt
		player.Score[holeIndex]["total"] = total

		playerQuery := bson.M{"_id": player.Id}
		if err := playerCol.Update(playerQuery, player); err != nil {
			return &Status{"failed update score"}, err
		}
	}

	//更新情報をGlobal変数に格納する
	players = GetAllPlayerCol()
	teams = GetAllTeamCol()
	return &Status{"success"}, nil
}
