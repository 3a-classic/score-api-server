package mongo

import (
	"logger"

	"errors"
	"log"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"labix.org/v2/mgo/bson"
)

func RegisterUserColData(userCols []UserCol) (*Status, error) {
	logger.Output(
		logrus.Fields{"User Collection": userCols},
		"Register User",
		logger.Info,
	)

	defer SetAllUserCol()

	db, session := mongoInit()
	col := db.C("user")
	defer session.Close()

	var createCnt, updateCnt int

	for _, userCol := range userCols {
		if len(userCol.Name) == 0 {
			logger.Output(
				logrus.Fields{logger.TraceMsg: logger.Trace(), "User Name": userCol.Name},
				"User name is not exist",
				logger.Error,
			)
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
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:     err,
					logger.TraceMsg:   logger.Trace(),
					"Find Query":      findQuery,
					"User Collection": userCol,
				},
				"can not upsert users",
				logger.Error,
			)
			return &Status{"can not upsert"}, err
		}
		if change.Updated == 0 {
			createCnt += 1
		} else {
			updateCnt += 1
		}
	}

	return &Status{"success"}, nil
}

func RegisterTeamColData(date string, teamCols []TeamCol) (*Status, error) {
	logger.Output(
		logrus.Fields{"Date": date, "Team Collection": teamCols},
		"Register Team and Player",
		logger.Info,
	)

	defer SetAllPlayerCol()
	defer SetAllTeamCol()

	db, session := mongoInit()
	playerC := db.C("player")
	teamC := db.C("team")
	defer session.Close()

	var OneBeforebyteOfA byte = 64
	alphabet := make([]byte, 1)
	alphabet[0] = OneBeforebyteOfA

	totalHoleNum := 18

	for _, teamCol := range teamCols {
		if len(teamCol.UserIds) == 0 {
			logger.Output(
				logrus.Fields{logger.TraceMsg: logger.Trace(), "User IDs": teamCol.UserIds},
				"User IDs are not exist",
				logger.Error,
			)
			return &Status{"this team do not have user id"}, nil
		}

		alphabet[0] += 1
		teamCol.Name = string(alphabet)
		teamCol.Defined = false
		teamCol.Date = date

		if err := teamC.Insert(teamCol); err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:     err,
					logger.TraceMsg:   logger.Trace(),
					"Team Collection": teamCol,
				},
				"can not insert team",
				logger.Error,
			)
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
				logger.Output(
					logrus.Fields{
						logger.ErrMsg:       err,
						logger.TraceMsg:     logger.Trace(),
						"Player Collection": player,
					},
					"can not insert player",
					logger.Error,
				)
				return &Status{"can not insert"}, err
			}
		}
	}
	return &Status{"success"}, nil
}

func RegisterFieldColData(date string, fieldCols []FieldCol) (*Status, error) {
	logger.Output(
		logrus.Fields{"Date": date, "Field Collection": fieldCols},
		"Register Field",
		logger.Info,
	)

	defer SetAllFieldCol()

	db, session := mongoInit()
	fieldC := db.C("field")
	defer session.Close()
	log.Println("field playerの登録を開始します。")

	var createCnt, updateCnt int
	for _, fieldCol := range fieldCols {
		if fieldCol.Hole > 18 || fieldCol.Hole < 0 {
			logger.Output(
				logrus.Fields{logger.TraceMsg: logger.Trace(), "Hole": fieldCol.Hole},
				"this is not hole number",
				logger.Error,
			)
			return &Status{"this is not hole number"}, err
		}

		fieldCol.Ignore = false
		fieldCol.Date = date
		findQuery := bson.M{"hole": fieldCol.Hole}

		change, err := fieldC.Upsert(findQuery, fieldCol)
		if err != nil {
			logger.Output(
				logrus.Fields{
					logger.ErrMsg:      err,
					logger.TraceMsg:    logger.Trace(),
					"Find Query":       findQuery,
					"Field Collection": fieldCol,
				},
				"can not upsert field",
				logger.Error,
			)
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
	logger.Output(
		logrus.Fields{"Request Take Picture Status": r},
		"Register Thread Image",
		logger.Info,
	)

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
		logger.Output(
			logrus.Fields{
				logger.ErrMsg:   err,
				logger.TraceMsg: logger.Trace(),
				"Find Query":    threadFindQuery,
				"Set Query":     threadSetQuery,
			},
			"can not update thread",
			logger.Error,
		)
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
		logger.Output(
			logrus.Fields{
				logger.ErrMsg:   err,
				logger.TraceMsg: logger.Trace(),
				"Find Query":    playerFindQuery,
				"Set Query":     playerSetQuery,
			},
			"can not update player",
			logger.Error,
		)
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
