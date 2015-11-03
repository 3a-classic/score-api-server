package mongo

import (
	"log"
	"strconv"
	"time"

	"labix.org/v2/mgo/bson"
)

func RegisterUserColData(userCols []UserCol) (*Status, error) {

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
