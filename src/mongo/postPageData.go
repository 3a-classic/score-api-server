package mongo

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"
)

func PostLoginPageData(loginInfo *PostLogin) (*LoginStatus, error) {

	fmt.Println(loginInfo.UserId + "さんがアクセスしました。")
	_, ok := users[loginInfo.UserId]
	if ok {
		fmt.Println(users[loginInfo.UserId].Name + "さんがログインしました。")
		loginStatus := &LoginStatus{
			Status:   "success",
			UserId:   loginInfo.UserId,
			UserName: users[loginInfo.UserId].Name,
			Admin:    players[loginInfo.UserId].Admin,
		}
		return loginStatus, nil
	} else {
		fmt.Println(loginInfo.UserId + "さんがログインに失敗しました。")
		return &LoginStatus{Status: "failed"}, nil
	}
}

func PostApplyScoreData(teamName string, ApplyScore *PostApplyScore) (*Status, error) {

	//更新情報をGlobal変数に格納する
	defer SetPlayerCol(ApplyScore.UserIds)

	fmt.Println("Team : " + teamName + " 申請スコアを登録します。")

	AUserIdInTheTeam := teams[teamName].UserIds[0]
	if players[AUserIdInTheTeam].Apply != 0 {
		return &Status{"already registered"}, nil
	}
	log.Println(ApplyScore.Apply)
	log.Println(ApplyScore.UserIds)
	for playerIndex, userId := range ApplyScore.UserIds {

		log.Println(ApplyScore.Apply[playerIndex])
		findQuery := bson.M{"userid": userId}
		setQuery := bson.M{"$set": bson.M{"apply": ApplyScore.Apply[playerIndex]}}
		if err = UpdateMongoData("player", findQuery, setQuery); err != nil {
			return &Status{"failed"}, err
		}
	}

	return &Status{"success"}, nil
}

func PostScoreViewSheetPageData(teamName string, definedTeam *PostDefinedTeam) (*Status, error) {
	//更新情報をGlobal変数に格納する
	defer SetTeamCol(teamName)

	fmt.Println("Team : " + teamName + "のデータを確定します。")

	findQuery := bson.M{"name": teamName}
	setQuery := bson.M{"$set": bson.M{"defined": true}}
	if err = UpdateMongoData("team", findQuery, setQuery); err != nil {
		return &Status{"failed"}, err
	}

	return &Status{"success"}, nil
}

func PostScoreEntrySheetPageData(teamName string, holeString string, teamScore *PostTeamScore) (*RequestTakePictureStatus, error) {

	userIds := teams[teamName].UserIds
	//更新情報をGlobal変数に格納する
	defer SetPlayerCol(userIds)

	if len(holeString) == 0 {
		return &RequestTakePictureStatus{Status: "failed"}, errors.New("hole is not string")
	}

	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1
	holeIndexString := strconv.Itoa(holeIndex)

	if teamScore.Excnt != excnt[teamName][holeNum] {
		return &RequestTakePictureStatus{Status: "other updated"}, nil
	} else {
		excnt[teamName][holeNum]++
	}

	fmt.Println("Team : " + teamName + ", Hole : " + holeString + "にデータを挿入します。")

	for playerIndex, userId := range teamScore.UserIds {
		total, putt := teamScore.Total[playerIndex], teamScore.Putt[playerIndex]

		findQuery := bson.M{"userid": userId}
		setQuery := bson.M{
			"$set": bson.M{
				"score." + holeIndexString + ".total": total,
				"score." + holeIndexString + ".putt":  putt,
			},
		}
		if err = UpdateMongoData("player", findQuery, setQuery); err != nil {
			return &RequestTakePictureStatus{Status: "failed update score"}, err
		}
	}
	//	Thread登録
	if err := RegisterThreadOfScore(holeString, teamScore); err != nil {
		return nil, err
	}

	//	チーム内に写真リクエストがあるか確認する
	requestTakePictureStatus, err := RequestTakePicture(userIds)
	if err != nil {
		return nil, err
	}

	log.Println("before return : ", requestTakePictureStatus)
	return requestTakePictureStatus, nil
}

func UpsertNewTimeLine(thread *Thread) error {
	targetThreadId := thread.ThreadId
	//更新情報をGlobal変数に格納する
	defer SetAllThreadCol()

	colorFeeling := make(map[string]string)
	colorFeeling["default"] = "#FFFFFF"
	colorFeeling["angry"] = "#FF0000"
	colorFeeling["great"] = "#FFFF00"
	colorFeeling["sad"] = "#0000FF"
	colorFeeling["vexing"] = "#00FF00"

	db, session := mongoInit()
	threadCol := db.C("thread")
	defer session.Close()

	//新規スレッドの時
	if len(targetThreadId) == 0 {
		log.Println("insert thread")

		thread.ThreadId = make20lengthHashString()
		thread.CreatedAt = time.Now().Format(datetimeFormat)
		thread.ColorCode = colorFeeling["default"]
		if err = threadCol.Insert(thread); err != nil {
			return err
		}

		//既存スレッドに対する反応の時
	} else {
		if len(thread.ColorCode) == 0 {
			return errors.New("current colorCode is not contain in posted thread")
		}
		if len(thread.Reactions) > 1 {
			return errors.New("reactions is not 1")
		}
		log.Println("update reaction of the thread")

		currentFeeling := ""
		currentColor := threads[targetThreadId].ColorCode
		postedFeeling := getFeelingFromAWSUrl(thread.Reactions[0].Content)
		postedColor := colorFeeling[postedFeeling]

		for feeling, code := range colorFeeling {
			if currentColor == code {
				currentFeeling = feeling
			}
		}

		thread.Reactions[0].DateTime = time.Now().Format(datetimeFormat)
		findQuery := bson.M{"threadid": targetThreadId}
		pushQuery := bson.M{"$push": bson.M{"reactions": thread.Reactions[0]}}
		if err = threadCol.Update(findQuery, pushQuery); err != nil {
			return err
		}

		//投稿された感情と、現在の感情に相違がある場合
		if currentFeeling != postedFeeling {
			var setColor string
			var currentFeelingCount, postedFeelingCount int

			if currentColor == colorFeeling["default"] {
				setColor = postedColor
			} else {

				for _, r := range threads[targetThreadId].Reactions {
					switch getFeelingFromAWSUrl(r["content"].(string)) {
					case currentFeeling:
						currentFeelingCount++
					case postedFeeling:
						postedFeelingCount++
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
	log.Println("upsert thread done")
	return nil
}
