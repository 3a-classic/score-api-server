package mongo

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

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
		if err = playerCol.Update(query, player); err != nil {
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
	if err = teamCol.Update(query, targetTeam); err != nil {
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
	if err = teamCol.Update(teamQuery, targetTeam); err != nil {
		return &Status{"failed update excnt"}, err
	}

	teamPlayers := GetPlayersDataInTheTeam(teamName)

	for playerIndex, player := range teamPlayers {
		total, putt := updatedTeamScore.Total[playerIndex], updatedTeamScore.Putt[playerIndex]
		player.Score[holeIndex]["putt"] = putt
		player.Score[holeIndex]["total"] = total

		playerQuery := bson.M{"_id": player.Id}
		if err = playerCol.Update(playerQuery, player); err != nil {
			return &Status{"failed update score"}, err
		}
	}

	//更新情報をGlobal変数に格納する
	players = GetAllPlayerCol()
	teams = GetAllTeamCol()
	return &Status{"success"}, nil
}

func UpsertNewTimeLine(thread *Thread) error {

	var datetimeFormat = "2006/01/02 15:04:05 MST"
	colorFeeling := make(map[string]string)
	colorFeeling["default"] = "#FFFFFF"
	colorFeeling["angry"] = "#FF0000"
	colorFeeling["great"] = "#FFFF00"
	colorFeeling["sad"] = "#0000FF"
	colorFeeling["vexing"] = "#00FF00"

	db, session := mongoInit()
	threadCol := db.C("thread")
	defer session.Close()

	if len(thread.ThreadId) == 0 {
		log.Println("insert thread")

		thread.ThreadId = make20lengthHashString()
		thread.CreatedAt = time.Now().Format(datetimeFormat)
		thread.ColorCode = colorFeeling["default"]
		if err = threadCol.Insert(thread); err != nil {
			return err
		}

	} else {
		log.Println("update reaction of the thread")

		var currentFeeling, postedFeeling string
		var currentColor, postedColor string

		postedFeeling = getFeelingFromAWSUrl(thread.Reactions[0].Content)
		postedColor = colorFeeling[postedFeeling]

		if len(thread.ColorCode) == 0 {
			return errors.New("current colorCode is not contain in posted thread")
		}
		currentColor = thread.ColorCode
		for feeling, code := range colorFeeling {
			if currentColor == code {
				currentFeeling = feeling
			}
		}

		thread.Reactions[0].DateTime = time.Now().Format(datetimeFormat)
		if len(thread.Reactions) > 1 {
			return errors.New("reactions is not 1")
		}
		findQuery := bson.M{"threadid": thread.ThreadId}
		pushQuery := bson.M{"$push": bson.M{"reactions": thread.Reactions[0]}}
		if err = threadCol.Update(findQuery, pushQuery); err != nil {
			return err
		}
		if currentFeeling != postedColor {
			var setColor string
			var currentFeelingCount, postedFeelingCount int

			if currentColor == colorFeeling["default"] {
				setColor := postedColor
			} else {

				for _, t := range threads {
					if t.ThreadId == thread.ThreadId {
						for _, r := range t.Reactions {
							switch getFeelingFromAWSUrl(r["content"].(string)) {
							case currentFeeling:
								currentFeelingCount++
							case postedFeeling:
								postedFeelingCount++
							}
						}
					}
				}

				if currentFeelingCount >= postedFeelingCount {
					setColor = currentColor
				} else {
					setColor = postedColor
				}
			}

			setQuery := bson.M{"$set": bson.M{"colorcode": setColor}}

			if err = threadCol.Update(findQuery, setQuery); err != nil {
				return err
			}
		}
	}
	//更新情報をGlobal変数に格納する
	threads = GetAllThreadCol()
	return nil
}

// utils
func getFeelingFromAWSUrl(url string) string {
	regexpString := "https://s3-ap-northeast-1.amazonaws.com/3a-classic/reaction-icon/(.+).png"
	re := regexp.MustCompile(regexpString)
	return re.FindStringSubmatch(url)[1]
}

func make20lengthHashString() string {
	b := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, b)

	if err != nil {
		return err
	}
	longHash := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	return string([]rune(longHash)[:20])
}
