package mongo

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"regexp"
	"strings"

	"labix.org/v2/mgo/bson"
)

// utils
func RequestTakePicture(userIds []string) *RequestTakePictureStatus {

	return &RequestTakePictureStatus{Status: "success"}
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
