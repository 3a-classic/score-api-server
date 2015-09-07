package mongo

import (
	"fmt"
	"sort"
	"strconv"

	"labix.org/v2/mgo/bson"
)

func GetIndexPageData() (*Index, error) {
	teams := GetAllTeamCol()

	teamArray := make([]string, len(teams))
	for i := 0; i < len(teams); i++ {
		teamArray[i] = teams[i].Team
	}

	idx := &Index{
		Team:   teamArray,
		Length: len(teams),
	}
	return idx, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {
	players := GetAllPlayerCol()
	fields := GetAllFieldCol()

	var leadersBoard LeadersBoard
	usersScore := make([]UserScore, len(players))
	var userScore UserScore

	for personNum := 0; personNum < len(players); personNum++ {
		totalPar := 0
		totalScore := 0
		passedHoleCnt := 0
		for holeNum := 0; holeNum < len(players[personNum].Score); holeNum++ {
			if players[personNum].Score[holeNum]["total"] != 0 {
				holeIndex := fields[holeNum].Hole - 1
				totalPar += fields[holeIndex].Par
				totalScore += players[personNum].Score[holeNum]["total"].(int)
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore - totalPar,
			Total: totalScore,
			Name:  players[personNum].Name,
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
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	field := GetOneFieldByQuery(bson.M{"hole": holeNum})
	playersInTheTeam := GetPlayersDataInTheTeam(teamName)

	member := make([]string, 4)
	stroke := make([]int, 4)
	putt := make([]int, 4)
	for playerIndex := 0; playerIndex < len(playersInTheTeam); playerIndex++ {
		member[playerIndex] = playersInTheTeam[playerIndex].Name
		stroke[playerIndex] = playersInTheTeam[playerIndex].Score[holeIndex]["stroke"].(int)
		putt[playerIndex] = playersInTheTeam[playerIndex].Score[holeIndex]["putt"].(int)
	}

	scoreEntrySheet := ScoreEntrySheet{
		Team:   teamName,
		Hole:   holeNum,
		Member: member,
		Par:    field.Par,
		Yard:   field.Yard,
		Stroke: stroke,
		Putt:   putt,
		Excnt:  0,
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {

	playersInTheTeam := GetPlayersDataInTheTeam(teamName)
	fields := GetAllFieldCol()

	member := make([]string, len(playersInTheTeam))
	apply := make([]int, len(playersInTheTeam))
	holes := make([]Hole, len(fields))
	totalScore := make([]int, len(playersInTheTeam))
	totalPar := 0
	for holeIndex := 0; holeIndex < len(fields); holeIndex++ {

		totalPar += fields[holeIndex].Par
		score := make([]int, len(playersInTheTeam))
		for playerIndex := 0; playerIndex < len(playersInTheTeam); playerIndex++ {
			score[playerIndex] = playersInTheTeam[playerIndex].Score[holeIndex]["total"].(int)
			totalScore[playerIndex] += score[playerIndex]
			if holeIndex == 0 {
				member[playerIndex] = playersInTheTeam[playerIndex].Name
				apply[playerIndex] = playersInTheTeam[playerIndex].Apply
			}
		}
		holes[holeIndex] = Hole{
			Hole:  holeIndex + 1,
			Par:   fields[holeIndex].Par,
			Yard:  fields[holeIndex].Yard,
			Score: score,
		}
	}

	sum := Sum{
		Par:   totalPar,
		Score: totalScore,
	}
	scoreViewSheet := ScoreViewSheet{
		Team:   teamName,
		Member: member,
		Apply:  apply,
		Hole:   holes,
		Sum:    sum,
	}

	return &scoreViewSheet, nil
}

func PostScoreEntrySheetPageData(teamName string, holeString string, updatedTeamScore *PostTeamScore) (*Status, error) {

	if len(holeString) == 0 {
		return nil, nil
	}
	holeNum, err := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	fmt.Println("Team : " + teamName + ", Hole : " + holeString + "にデータを挿入します。")
	db, session := mongoInit()
	playerCol := db.C("player")
	defer session.Close()

	for i := 0; i < len(updatedTeamScore.Member); i++ {
		player := Player{}

		query := bson.M{"name": updatedTeamScore.Member[i]}
		err = playerCol.Find(query).One(&player)
		if err != nil {
			return nil, err
		}

		stroke, putt := updatedTeamScore.Stroke[i], updatedTeamScore.Putt[i]
		player.Score[holeIndex]["stroke"] = stroke
		player.Score[holeIndex]["putt"] = putt
		player.Score[holeIndex]["total"] = stroke + putt

		err = playerCol.Update(query, player)
		if err != nil {
			return nil, err
		}
	}

	status := Status{
		Status: "success",
	}
	return &status, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
