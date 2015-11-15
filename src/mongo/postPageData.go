package mongo

import (
	c "config"
	l "logger"

	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
)

func PostLoginPageData(loginInfo *PostLogin) (*LoginStatus, error) {

	l.PutInfo(l.I_M_PostPage, loginInfo, nil)

	log.Println("users")
	log.Println("users")
	log.Println("users")
	log.Println(users)
	log.Println("users")
	log.Println("users")
	log.Println("users")
	_, ok := users[loginInfo.UserId]
	if ok {
		var teamName string
		for _, team := range teams {
			for _, userid := range team.UserIds {
				if userid == loginInfo.UserId {
					teamName = team.Name
				}
			}
		}

		l.Output(
			logrus.Fields{
				"User ID":   loginInfo.UserId,
				"User Name": users[loginInfo.UserId].Name,
			},
			"login success",
			l.Debug,
		)
		loginStatus := &LoginStatus{
			Status:   "success",
			UserId:   loginInfo.UserId,
			UserName: users[loginInfo.UserId].Name,
			Team:     teamName,
			Admin:    players[loginInfo.UserId].Admin,
		}
		return loginStatus, nil
	} else {
		l.Output(
			logrus.Fields{
				"User ID": loginInfo.UserId,
			},
			"login faild",
			l.Debug,
		)
		return &LoginStatus{Status: "failed"}, nil
	}
}

func PostApplyScoreData(teamName string, ApplyScore *PostApplyScore) (*Status, error) {
	l.PutInfo(l.I_M_PostPage, teamName, ApplyScore)

	//更新情報をGlobal変数に格納する
	defer SetPlayerCol(ApplyScore.UserIds)

	AUserIdInTheTeam := teams[teamName].UserIds[0]
	if players[AUserIdInTheTeam].Apply != 0 {
		l.Output(
			logrus.Fields{
				"User Apply": l.Sprintf(players[AUserIdInTheTeam].Apply),
			},
			"Apply score is already registered",
			l.Debug,
		)
		return &Status{"already registered"}, nil
	}
	for playerIndex, userId := range ApplyScore.UserIds {

		findQuery := bson.M{"userid": userId}
		setQuery := bson.M{"$set": bson.M{"apply": ApplyScore.Apply[playerIndex]}}
		if err = UpdateMongoData("player", findQuery, setQuery); err != nil {
			l.PutErr(err, l.Trace(), l.E_M_Update, ApplyScore.Apply[playerIndex])
			return &Status{"failed"}, err
		}
	}

	return &Status{"success"}, nil
}

func PostScoreViewSheetPageData(teamName string, definedTeam *PostDefinedTeam) (*Status, error) {
	l.PutInfo(l.I_M_PostPage, teamName, definedTeam)

	//更新情報をGlobal変数に格納する
	defer SetTeamCol(teamName)

	findQuery := bson.M{"name": teamName}
	setQuery := bson.M{"$set": bson.M{"defined": true}}
	if err = UpdateMongoData("team", findQuery, setQuery); err != nil {
		l.PutErr(err, l.Trace(), l.E_M_Update, teamName)
		return &Status{"failed"}, err
	}

	return &Status{"success"}, nil
}

func PostScoreEntrySheetPageData(teamName string, holeString string, teamScore *PostTeamScore) (*RequestTakePictureStatus, error) {
	l.PutInfo(l.I_M_PostPage, teamName, teamScore)

	userIds := teams[teamName].UserIds
	//更新情報をGlobal変数に格納する
	defer SetPlayerCol(userIds)

	if len(holeString) == 0 {
		l.PutErr(nil, l.Trace(), l.E_Nil, teamName)
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
			l.PutErr(err, l.Trace(), l.E_M_Update, userId)
			return &RequestTakePictureStatus{Status: "failed update score"}, err
		}
	}
	//	Thread登録
	if err := RegisterThreadOfScore(holeString, teamScore); err != nil {
		l.PutErr(err, l.Trace(), l.E_M_RegisterThread, teamScore)
		return nil, err
	}

	//	チーム内に写真リクエストがあるか確認する
	requestTakePictureStatus, err := RequestTakePicture(userIds)
	if err != nil {
		l.PutErr(err, l.Trace(), l.E_M_SearchPhotoTask, userIds)
		return nil, err
	}

	return requestTakePictureStatus, nil
}

func UpsertNewTimeLine(thread *Thread) error {
	l.PutInfo(l.I_M_PostPage, thread, nil)

	//更新情報をGlobal変数に格納する
	defer SetAllThreadCol()

	defaultColor := "#c0c0c0"

	if len(thread.ThreadId) != 0 {
		l.PutErr(nil, l.Trace(), l.E_WrongData, thread)
		return errors.New("thread id exists")
	}

	db, session := mongoConn()
	threadCol := db.C("thread")
	defer session.Close()

	thread.ThreadId = make20lengthHashString()
	thread.CreatedAt = time.Now().Format(c.DatetimeFormat)
	thread.ColorCode = defaultColor
	if err = threadCol.Insert(thread); err != nil {
		l.PutErr(err, l.Trace(), l.E_M_Insert, thread)
		return err
	}

	ThreadChan <- thread

	return nil
}

func UpdateExistingTimeLine(thread *Thread) (*Thread, error) {
	l.PutInfo(l.I_M_PostPage, thread, nil)

	targetThreadId := thread.ThreadId
	//更新情報をGlobal変数に格納する
	defer SetAllThreadCol()

	colorFeeling := make(map[string]string)
	colorFeeling["default"] = "#c0c0c0"
	colorFeeling["angry"] = "#ff7f7f"
	colorFeeling["great"] = "#ffff7f"
	colorFeeling["sad"] = "#7fbfff"
	colorFeeling["vexing"] = "#7fff7f"

	if len(targetThreadId) == 0 {
		l.PutErr(nil, l.Trace(), l.E_Nil, thread)
		return nil, errors.New("thread id do not exist")
	}

	if len(thread.ColorCode) == 0 {
		l.PutErr(err, l.Trace(), l.E_Nil, thread.ColorCode)
		return nil, errors.New("current colorCode do not contain in posted thread")
	}
	if len(thread.Reactions) > 1 {
		l.PutErr(err, l.Trace(), l.E_TooManyData, thread.Reactions)
		return nil, errors.New("reactions is not 1")
	}

	currentFeeling := ""
	currentColor := threads[targetThreadId].ColorCode
	postedFeeling := getFeelingFromAWSUrl(thread.Reactions[0].Content)
	postedColor := colorFeeling[postedFeeling]

	for feeling, code := range colorFeeling {
		if currentColor == code {
			currentFeeling = feeling
		}
	}

	thread.Reactions[0].DateTime = time.Now().Format(c.DatetimeFormat)
	findQuery := bson.M{"threadid": targetThreadId}
	pushQuery := bson.M{"$push": bson.M{"reactions": thread.Reactions[0]}}
	if err = UpdateMongoData("thread", findQuery, pushQuery); err != nil {
		l.PutErr(err, l.Trace(), l.E_M_Update, thread.Reactions[0])
		return nil, err
	}

	//投稿された感情と、現在の感情に相違がある場合
	if currentFeeling != postedFeeling {
		var setColor string
		var currentFeelingCount, postedFeelingCount int

		if currentColor == colorFeeling["default"] {
			setColor = postedColor
		} else {

			//直前の更新を反映させる
			SetAllThreadCol()

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

		if err = UpdateMongoData("thread", findQuery, setQuery); err != nil {
			l.PutErr(err, l.Trace(), l.E_M_Update, setColor)
			return nil, err
		}
		thread.ColorCode = setColor
	}
	return thread, nil
}
