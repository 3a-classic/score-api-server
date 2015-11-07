package mongo

import (
	"labix.org/v2/mgo/bson"
)

type Config struct {
	Mongo struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		Database string `toml:"database"`
	} `toml:"mongo"`
}

// Structs for Collections

//type Score struct {
//	Hole  int `json:"hole"`
//	Putt  int `json:"putt"`
//	Total int `json:"total"`
//}

type UserCol struct {
	Id            bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserId        string        `json:"userId"`
	Name          string        `json:"name"`
	CreatedAt     string        `json:"createdAt"`
	ImgUrl        string        `json:"imgUrl"`
	Participation []bson.M      `json:"participation"`
}

type PlayerCol struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserId   string        `json:"userId"`
	Apply    int           `json:"apply"`
	Editable bool          `json:"editable"`
	Score    []bson.M      `json:"score"`
	//	Team     mgo.DBRef     `json:"team"`
	Date  string `json:"Date"`
	Admin bool   `json:"admin"`
}

type TeamCol struct {
	Id      bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserIds []string      `json:"userIds"`
	Name    string        `json:"name"`
	Defined bool          `json:"defined"`
	Date    string        `json:"Date"`
}

type FieldCol struct {
	Id             bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Hole           int           `json:"hole"`
	DrivingContest bool          `json:"drivingContest"`
	Ignore         bool          `json:"ignore"`
	Image          string        `json:"image"`
	NearPin        bool          `json:"nearPin"`
	Par            int           `json:"par"`
	Yard           int           `json:"yard"`
	Date           string        `json:"Date"`
}

type ThreadCol struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserId    string        `json:"userId"`
	ThreadId  string        `json:"threadId"`
	Msg       string        `json:"msg"`
	ImgUrl    string        `json:"imgUrl"`
	ColorCode string        `json:"colorCode"`
	Positive  bool          `json:"positive"`
	Reactions []bson.M      `json:"reactions"`
	CreatedAt string        `json:"createdAt"`
	Date      string        `json:"Date"`
}

// Structs for Page
type UserScore struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
	Hole  int    `json:"hole"`
	Total int    `json:"total"`
}

type Hole struct {
	Hole  int   `json:"hole"`
	Par   int   `json:"par"`
	Yard  int   `json:"yard"`
	Score []int `json:"score"`
}

type Sum struct {
	Par   int   `json:"par"`
	Score []int `json:"score"`
	Putt  []int `json:"putt"`
}

type Reaction struct {
	Name        string `json:"name"`
	ContentType int    `json:"contentType"`
	Content     string `json:"content"`
	DateTime    string `json:"dateTime"`
}

type Thread struct {
	ThreadId  string     `json:"threadId"`
	UserId    string     `json:"userId"`
	Msg       string     `json:"msg"`
	ImgUrl    string     `json:"imgUrl"`
	ColorCode string     `json:"colorCode"`
	Positive  bool       `json:"positive"`
	Reactions []Reaction `json:"reactions"`
	CreatedAt string     `json:"createdAt"`
}

type sortByScore []UserScore

type Index struct {
	Team   []string `json:"team"`
	Length int      `json:"length"`
}

type LeadersBoard struct {
	Ranking []UserScore `json:"ranking"`
}

type ScoreEntrySheet struct {
	Team   string   `json:"team"`
	Hole   int      `json:"hole"`
	Member []string `json:"member"`
	Par    int      `json:"par"`
	Yard   int      `json:"yard"`
	Total  []int    `json:"total"`
	Putt   []int    `json:"putt"`
	Excnt  int      `json:"excnt"`
}

type ScoreViewSheet struct {
	Team    string   `json:"team"`
	Member  []string `json:"member"`
	UserIds []string `json:"userIds"`
	Apply   []int    `json:"apply"`
	Hole    []Hole   `json:"hole"`
	OutSum  Sum      `json:"outSum"`
	InSum   Sum      `json:"inSum"`
	Sum     Sum      `json:"sum"`
	Defined bool     `json:"defined"`
}

type EntireScore struct {
	Rows [][]string `json:"rows"`
}

type TimeLine struct {
	Threads []Thread `json:"threads"`
}

type PostLogin struct {
	UserId string `json:"userId"`
}

type PostApplyScore struct {
	UserIds []string `json:"userIds"`
	Apply   []int    `json:"apply"`
}

type PostDefinedTeam struct {
	Team string `json:"team"`
}

type PostTeamScore struct {
	UserIds []string `json:"userIds"`
	Total   []int    `json:"total"`
	Putt    []int    `json:"putt"`
	Excnt   int      `json:"excnt"`
}

// reponse status

type Status struct {
	Status string `json:"status"`
}

type RequestTakePictureStatus struct {
	Status    string `json:"status"`
	UserId    string `json:"userId"`
	Name      string `json:"name"`
	ThreadMsg string `json:"threadMsg"`
	PhotoUrl  string `json:"photoUrl"`
}

// websocket struct

type TimeLineWs struct {
	Msg string `json:"msg"`
}
