package mongo

import (
	//	"encoding/json"
	"fmt"
	//	"labix.org/v2/mgo"
	"sort"

	"labix.org/v2/mgo/bson"
	"reflect"
	"strconv"
)

//Children

type UserScore struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
	Hole  int    `json:"hole"`
	Total int    `json:"team"`
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
}

type sortByScore []UserScore

// Parents
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
	Stroke []int    `json:"stroke"`
	Putt   []int    `json:"putt"`
	Excnt  int      `json:"excnt"`
}

type ScoreViewSheet struct {
	Team   string   `json:"team"`
	Member []string `json:"member"`
	Apply  []int    `json:"apply"`
	Hole   []Hole   `json:"hole"`
	Sum    Sum      `json:"sum"`
}

type PostTeamScore struct {
	Member []string `json:"member"`
	Stroke []int    `json:"stroke"`
	Putt   []int    `json:"putt"`
	Excnt  int      `json:"excnt"`
}

type Status struct {
	Status string `json:"status"`
}

func GetIndexPageData() (*Index, error) {
	collectionName := "team"
	col, session := mongoInit(collectionName)
	defer session.Close()
	result := []Team{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}

	teamArr := make([]string, len(result))
	for i := 0; i < len(result); i++ {
		teamArr[i] = result[i].Team
	}

	idx := &Index{teamArr, len(result)}
	return idx, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {
	collectionName := "player"
	collectionName1 := "field"
	col, session := mongoInit(collectionName)
	col1, session1 := mongoInit(collectionName1)
	defer session.Close()
	defer session1.Close()
	result := []Player{}
	result1 := []Field{}
	err := col.Find(bson.M{}).All(&result)
	if err != nil {
		panic(err)
	}
	err = col1.Find(bson.M{}).All(&result1)
	if err != nil {
		panic(err)
	}

	var leadersBoard LeadersBoard
	usersScore := make([]UserScore, len(result))
	var userScore UserScore

	for personNum := 0; personNum < len(result); personNum++ {
		totalPar := 0
		totalScore := 0
		passedHoleCnt := 0
		for holeNum := 0; holeNum < len(result[personNum].Score); holeNum++ {
			if result[personNum].Score[holeNum]["total"] != 0 {
				holeIndex := result1[holeNum].Hole - 1
				totalPar += result1[holeIndex].Par
				totalScore += result[personNum].Score[holeNum]["total"].(int)
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore - totalPar,
			Total: totalScore,
			Name:  result[personNum].Name,
			Hole:  passedHoleCnt,
		}
		usersScore[personNum] = userScore
	}
	sort.Sort(sortByScore(usersScore))
	leadersBoard.Ranking = usersScore

	return &leadersBoard, nil
}

func GetScoreEntrySheetPageData(teamName string, holeString string) (*ScoreEntrySheet, error) {
	if len(holeString) == 0 {
		return nil, nil
	}
	holeNum, err := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	collectionName1 := "player"
	collectionName2 := "field"
	collectionName3 := "team"
	col1, session1 := mongoInit(collectionName1)
	col2, session2 := mongoInit(collectionName2)
	col3, session3 := mongoInit(collectionName3)
	defer session1.Close()
	defer session2.Close()
	defer session3.Close()
	result1 := []Player{}
	result2 := []Field{}
	result3 := []Team{}
	result4 := []Player{}
	err = col1.Find(bson.M{}).All(&result1)
	if err != nil {
		panic(err)
	}
	err = col2.Find(bson.M{"hole": holeNum}).All(&result2)
	if err != nil {
		panic(err)
	}
	err = col3.Find(bson.M{"team": teamName}).All(&result3)
	if err != nil {
		panic(err)
	}

	err = col3.Find(bson.M{"team": teamName}).All(&result3)
	if err != nil {
		panic(err)
	}

	err = session1.FindRef(&result1[0].Team).All(&result3)
	if err != nil {
		panic(err)
	}
	teamMember := make([]string, 4)
	strokeScore := make([]int, 4)
	puttScore := make([]int, 4)
	for i := 0; i < len(result3[0].Member); i++ {
		err = session1.FindRef(&result3[0].Member[i].Player).All(&result4)
		if err != nil {
			panic(err)
		}
		teamMember[i] = result4[0].Name
		strokeScore[i] = result4[0].Score[holeIndex]["stroke"].(int)
		puttScore[i] = result4[0].Score[holeIndex]["putt"].(int)
	}

	scoreEntrySheet := ScoreEntrySheet{
		Team:   teamName,
		Hole:   holeNum,
		Member: teamMember,
		Par:    result2[0].Par,
		Yard:   result2[0].Yard,
		Stroke: strokeScore,
		Putt:   puttScore,
		Excnt:  0,
	}
	//get team name from dbref
	err = session1.FindRef(&result1[0].Team).All(&result3)
	if err != nil {
		panic(err)
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {

	holeNum := 18

	collectionName1 := "player"
	collectionName2 := "field"
	collectionName3 := "team"
	col1, session1 := mongoInit(collectionName1)
	col2, session2 := mongoInit(collectionName2)
	col3, session3 := mongoInit(collectionName3)
	defer session1.Close()
	defer session2.Close()
	defer session3.Close()
	result1 := []Player{}
	result2 := []Field{}
	result3 := []Team{}
	result4 := []Player{}
	err := col1.Find(bson.M{}).All(&result1)
	if err != nil {
		panic(err)
	}
	err = col2.Find(bson.M{}).All(&result2)
	if err != nil {
		panic(err)
	}
	err = col3.Find(bson.M{"team": teamName}).All(&result3)
	if err != nil {
		panic(err)
	}

	err = col3.Find(bson.M{"team": teamName}).All(&result3)
	if err != nil {
		panic(err)
	}

	err = session1.FindRef(&result1[0].Team).All(&result3)
	if err != nil {
		panic(err)
	}
	teamMember := make([]string, len(result3[0].Member))
	apply := make([]int, len(result3[0].Member))
	for i := 0; i < len(result3[0].Member); i++ {
		err = session1.FindRef(&result3[0].Member[i].Player).All(&result4)
		if err != nil {
			panic(err)
		}
		teamMember[i] = result4[0].Name
		apply[i] = result4[0].Apply
	}
	holes := make([]Hole, holeNum)
	totalScore := make([]int, len(result3[0].Member))
	totalPar := 0
	for holeIndex := 0; holeIndex < holeNum; holeIndex++ {

		totalPar += result2[holeIndex].Par
		score := make([]int, len(result3[0].Member))
		for playerIndex := 0; playerIndex < len(result3[0].Member); playerIndex++ {
			score[playerIndex] = result4[0].Score[holeIndex]["total"].(int)
			totalScore[playerIndex] += score[playerIndex]
		}
		holes[holeIndex] = Hole{
			Hole:  holeIndex + 1,
			Par:   result2[holeIndex].Par,
			Yard:  result2[holeIndex].Yard,
			Score: score,
		}
	}

	sum := Sum{
		Par:   totalPar,
		Score: totalScore,
	}
	scoreViewSheet := ScoreViewSheet{
		Team:   teamName,
		Member: teamMember,
		Apply:  apply,
		Hole:   holes,
		Sum:    sum,
	}

	return &scoreViewSheet, nil
}

func PostScoreEntrySheetPageData(teamName string, holeString string, updatedTeamScore *PostTeamScore) (*Status, error) {

	fmt.Println(teamName)
	fmt.Println(holeString)
	fmt.Println(updatedTeamScore)

	holeNum, err := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	fmt.Println(reflect.ValueOf(holeIndex).Type())

	collectionName := "player"
	fmt.Println(collectionName + "にデータを挿入します。")
	col, session := mongoInit(collectionName)
	defer session.Close()

	for i := 0; i < len(updatedTeamScore.Member); i++ {
		result := Player{}

		query := bson.M{"name": updatedTeamScore.Member[i]}
		err = col.Find(query).One(&result)
		if err != nil {
			panic(err)
		}
		fmt.Println(result)

		stroke := updatedTeamScore.Stroke[i]
		putt := updatedTeamScore.Putt[i]
		total := stroke + putt
		result.Score[holeIndex]["stroke"] = stroke
		result.Score[holeIndex]["putt"] = putt
		result.Score[holeIndex]["total"] = total

		fmt.Println(result)

		err = col.Update(query, result)
		if err != nil {
			return nil, err
		}
	}

	status := Status{
		Status: "successs!!!",
	}
	return &status, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
