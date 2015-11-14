package mongo

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"labix.org/v2/mgo/bson"
)

const (
	holeInOne         = "holeInOne"
	albatross         = "albatross"
	eagle             = "eagle"
	birdie            = "birdie"
	par               = "par"
	twoPointFiveTimes = "twoPointFiveTimes"
	bestTheHole       = "bestTheHole"
	worstTheHole      = "worstTheHole"

	bestACourse    = "bestACourse"
	worstACourse   = "worstACourse"
	bestOutCourse  = "bestOutCourse"
	worstOutCourse = "worstOutCourse"
	bestInCourse   = "bestInCourse"
	worstInCourse  = "worstInCourse"
	bestAllCourse  = "bestAllCourse"
	worstAllCourse = "worstAllCourse"
)

// utils
func RegisterThreadOfTotal() error {
	defer SetAllThreadCol()
	//	var threadMsg map[string]string
	//	var threadPositive map[string]bool

	for _, team := range teams {
		if !team.Defined {
			return nil
		}
	}

	//		https://s3-ap-northeast-1.amazonaws.com/3a-classic/emotion-img/positive.png
	//		https://s3-ap-northeast-1.amazonaws.com/3a-classic/emotion-img/negative.png

	return nil
}

func RegisterThreadOfScore(holeString string, teamScore *PostTeamScore) error {

	defer SetAllThreadCol()
	holeNum, _ := strconv.Atoi(holeString)
	//	holeIndex := holeNum - 1

	threadMsg := make(map[string]string)
	threadPositive := make(map[string]bool)

	threadMsg[holeInOne] = "ホールインワン"
	threadMsg[albatross] = "アルバトロス"
	threadMsg[eagle] = "イーグル"
	threadMsg[birdie] = "バーディー"
	threadMsg[par] = "パー"
	threadMsg[twoPointFiveTimes] = "２．５倍以上のスコア"

	threadPositive[holeInOne] = true
	threadPositive[albatross] = true
	threadPositive[eagle] = true
	threadPositive[birdie] = true
	threadPositive[par] = true
	threadPositive[twoPointFiveTimes] = false

	parInThisHole := fields[holeNum].Par

	for playerIndex, userId := range teamScore.UserIds {
		var threadKey, imgUrl, inOut string
		var holeInOutNum int

		switch teamScore.Total[playerIndex] {
		case 1:
			threadKey = holeInOne
		case parInThisHole - 3:
			threadKey = albatross
		case parInThisHole - 2:
			threadKey = eagle
		case parInThisHole - 1:
			threadKey = birdie
		case parInThisHole:
			threadKey = par
		default:
			if teamScore.Total[playerIndex] > int(float64(parInThisHole)*2.5) {
				threadKey = twoPointFiveTimes
			} else {
				continue
			}
		}

		if holeNum > 9 {
			holeInOutNum = holeNum - 9
			inOut = "IN"
		} else {
			holeInOutNum = holeNum
			inOut = "OUT"
		}
		holeInOutString := strconv.Itoa(holeInOutNum)

		if threadPositive[threadKey] {
			if len(players[userId].PositivePhotoUrl) != 0 {
				imgUrl = players[userId].PositivePhotoUrl
			}
		} else {
			if len(players[userId].NegativePhotoUrl) != 0 {
				imgUrl = players[userId].NegativePhotoUrl
			}
		}

		msg := makeScoreThreadMsg(
			threadPositive[threadKey],
			inOut,
			holeInOutString,
			users[userId].Name,
			threadMsg[threadKey],
		)

		thread := &Thread{
			UserId:   userId,
			UserName: users[userId].Name,
			Msg:      msg,
			ImgUrl:   imgUrl,
			Positive: threadPositive[threadKey],
		}

		log.Println("thread : ", thread)
		if err := UpsertNewTimeLine(thread); err != nil {
			return err
		}
	}

	log.Println("thread insert done")
	//	holeThreadScore := make(map[string]int)
	//	holeThreadUserId := make(map[string]string)
	//	holeThreadMsg := make(map[string]string)
	//	holeThreadPositive := make(map[string]string)
	//
	//	holeThreadScore[bestTheHole] = parInThisHole * 3
	//	holeThreadScore[worstTheHole] = 0
	//	holeThreadPositive[bestTheHole] = true
	//	holeThreadPositive[worstTheHole] = false
	//
	//
	//	for userId, player := range players {
	//		total := player.Score[holeIndex]["total"].(int)
	//		if total == 0 {
	//			return nil
	//		}
	//		if total < holeThreadScore[bestTheHole] {
	//			holeThreadScore[bestTheHole] = total
	//			holeThreadUserId[bestTheHole] = userId
	//		} else if total < holeThreadScore[worstTheHole] {
	//			holeThreadScore[worstTheHole] = total
	//			holeThreadUserId[worstTheHole] = userId
	//		}
	//	}
	//
	//	 msg := makeMsg(
	//		 holeThreadPositive[bestTheHole],
	//		 inOut,
	//		 holeString,
	//		 users[holeThreadUserId[bestTheHole]].Name,
	//		 strconv.Itoa(holeThreadScore[bestTheHole]),
	//	 )

	return nil
}

func makeScoreThreadMsg(positive bool, inOut string, holeString string, playerName string, threadMsg string) (msg string) {

	//msg = playerName + "さんが" + inOut + "の" + holeString + "番ホールで" + threadMsg
	msg = inOut + "の" + holeString + "番ホールで" + threadMsg
	if positive {
		msg = msg + "を取りました！！"
	} else {
		msg = msg + "を取ってしまいました。。"
	}

	return
}

func makeHoleThreadMsg(positive bool, inOut string, holeString string, playerName string, score string) (msg string) {

	//	msg = playerName + "さんが" + inOut + "の" + holeString + "番ホールでスコア" + score + "を出して"
	msg = inOut + "の" + holeString + "番ホールでスコア" + score + "を出して"
	if positive {
		msg = msg + "1番でした！！"
	} else {
		msg = msg + "ビリでした。。"
	}

	return
}

func RequestTakePicture(userIds []string) (*RequestTakePictureStatus, error) {

	db, session := mongoInit()
	col := db.C("thread")
	defer session.Close()

	threadCol := ThreadCol{}
	for _, userId := range userIds {
		findQuery := bson.M{"imgurl": "", "userid": userId}
		if err = col.Find(findQuery).One(&threadCol); err != nil {
			log.Println("request is completed!", userId)
		}

		if len(threadCol.ThreadId) != 0 {
			requestTakePictureStatus := &RequestTakePictureStatus{
				Status:      "take a picture",
				UserId:      threadCol.UserId,
				TeamUserIds: userIds,
				Name:        users[threadCol.UserId].Name,
				Positive:    threadCol.Positive,
				ThreadId:    threadCol.ThreadId,
				ThreadMsg:   threadCol.Msg,
				PhotoUrl:    "",
			}

			return requestTakePictureStatus, nil
		}
	}
	return &RequestTakePictureStatus{Status: "success"}, nil
}

func getFeelingFromAWSUrl(url string) string {
	regexpString := "https://s3-ap-northeast-1.amazonaws.com/3a-classic/reaction-icon/(.+).png"
	re := regexp.MustCompile(regexpString)
	return re.FindStringSubmatch(url)[1]
}

func make20lengthHashString() string {
	b := make([]byte, 32)
	_, err = io.ReadFull(rand.Reader, b)

	if err != nil {
		return err.Error()
	}
	longHash := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")

	return string([]rune(longHash)[:20])
}

func UpdateMongoData(collection string, findQuery bson.M, updateQuery bson.M) error {

	db, session := mongoInit()
	c := db.C(collection)
	defer session.Close()

	if err = c.Update(findQuery, updateQuery); err != nil {
		return err
	}
	return nil
}
