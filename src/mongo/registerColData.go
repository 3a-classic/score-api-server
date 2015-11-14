package mongo

import (
	"errors"
	"log"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"
)

func RegisterUserColData(userCols []UserCol) (*Status, error) {
	defer SetAllUserCol()

	db, session := mongoInit()
	col := db.C("user")
	defer session.Close()
	log.Println("User登録を開始します。")

	var createCnt, updateCnt int

	for _, userCol := range userCols {
		if len(userCol.Name) == 0 {
			return &Status{"this user do not have name"}, nil
		}
		if len(userCol.UserId) == 0 {
			userCol.UserId = make20lengthHashString()
		}
		if len(userCol.CreatedAt) == 0 {
			userCol.CreatedAt = time.Now().Format(datetimeFormat)
		}
		findQuery := bson.M{"userid": userCol.UserId}
		change, err := col.Upsert(findQuery, userCol)
		if err != nil {
			return &Status{"can not upsert"}, err
		}
		if change.Updated == 0 {
			createCnt += 1
		} else {
			updateCnt += 1
		}
	}
	log.Println(strconv.Itoa(createCnt) + "件作成されました。")
	log.Println(strconv.Itoa(updateCnt) + "件更新されました。")

	return &Status{"success"}, nil
}

func RegisterTeamColData(date string, teamCols []TeamCol) (*Status, error) {

	defer SetAllPlayerCol()
	defer SetAllTeamCol()

	db, session := mongoInit()
	playerC := db.C("player")
	teamC := db.C("team")
	defer session.Close()
	log.Println("team playerの登録を開始します。")

	var OneBeforebyteOfA byte = 64
	alphabet := make([]byte, 1)
	alphabet[0] = OneBeforebyteOfA

	totalHoleNum := 18

	for _, teamCol := range teamCols {
		if len(teamCol.UserIds) == 0 {
			continue
		}

		alphabet[0] += 1
		teamCol.Name = string(alphabet)
		teamCol.Defined = false
		teamCol.Date = date

		if err := teamC.Insert(teamCol); err != nil {
			return &Status{"can not insert"}, err
		}

		for _, userId := range teamCol.UserIds {
			scores := []bson.M{}
			for holeNum := 1; holeNum <= totalHoleNum; holeNum++ {
				score := bson.M{
					"hole":  holeNum,
					"putt":  0,
					"total": 0,
				}
				scores = append(scores, score)

			}

			player := PlayerCol{
				UserId: userId,
				Score:  scores,
				Date:   date,
			}
			if err := playerC.Insert(player); err != nil {
				return &Status{"can not insert"}, err
			}
		}
	}
	return &Status{"success"}, nil
}

func RegisterFieldColData(date string, fieldCols []FieldCol) (*Status, error) {
	defer SetAllFieldCol()

	db, session := mongoInit()
	fieldC := db.C("field")
	defer session.Close()
	log.Println("field playerの登録を開始します。")

	var createCnt, updateCnt int
	for _, fieldCol := range fieldCols {
		if fieldCol.Hole > 18 || fieldCol.Hole < 0 {
			continue
		}

		fieldCol.Ignore = false
		fieldCol.Date = date
		findQuery := bson.M{"hole": fieldCol.Hole}

		change, err := fieldC.Upsert(findQuery, fieldCol)
		if err != nil {
			return &Status{"can not upsert"}, err
		}
		if change.Updated == 0 {
			createCnt += 1
		} else {
			updateCnt += 1
		}
	}
	log.Println(strconv.Itoa(createCnt) + "件作成されました。")
	log.Println(strconv.Itoa(updateCnt) + "件更新されました。")

	return &Status{"success"}, nil
}

func RegisterThreadImg(r *RequestTakePictureStatus) (*RequestTakePictureStatus, error) {

	defer SetPlayerCol([]string{r.UserId})

	if len(r.ThreadId) == 0 {
		return nil, errors.New("there is not thread id")
	}

	if len(r.PhotoUrl) == 0 {
		return nil, errors.New("there is not pthoto url ")
	}

	threadFindQuery := bson.M{"threadid": r.ThreadId}
	threadSetQuery := bson.M{"$set": bson.M{"imgurl": r.PhotoUrl}}
	if err = UpdateMongoData("thread", threadFindQuery, threadSetQuery); err != nil {
		return &RequestTakePictureStatus{Status: "failed"}, err
	}

	var photoKey string
	if r.Positive {
		photoKey = "positivephotourl"
	} else {
		photoKey = "negativephotourl"
	}
	playerFindQuery := bson.M{"userid": r.UserId}
	playerSetQuery := bson.M{"$set": bson.M{photoKey: r.PhotoUrl}}

	if err = UpdateMongoData("player", playerFindQuery, playerSetQuery); err != nil {
		return &RequestTakePictureStatus{Status: "failed"}, err
	}

	SetAllThreadCol()
	thread := &Thread{
		ThreadId:  r.ThreadId,
		UserId:    threads[r.ThreadId].UserId,
		UserName:  threads[r.ThreadId].UserName,
		Msg:       threads[r.ThreadId].Msg,
		ImgUrl:    threads[r.ThreadId].ImgUrl,
		ColorCode: threads[r.ThreadId].ColorCode,
		Positive:  threads[r.ThreadId].Positive,
		CreatedAt: threads[r.ThreadId].CreatedAt,
	}

	ThreadChan <- thread

	requestTakePictureStatus, err := RequestTakePicture(r.TeamUserIds)
	if err != nil {
		return nil, err
	}

	if len(requestTakePictureStatus.ThreadId) != 0 {
		return requestTakePictureStatus, nil
	}

	return &RequestTakePictureStatus{Status: "success"}, nil

}
