package mongo

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"labix.org/v2/mgo/bson"
)

var (
	players []Player
	fields  []Field
	teams   []Team
)

func init() {
	players = GetAllPlayerCol()
	fields = GetAllFieldCol()
	teams = GetAllTeamCol()
}

func GetIndexPageData() (*Index, error) {

	teamArray := make([]string, len(teams))
	for i, team := range teams {
		teamArray[i] = team.Team
	}

	idx := &Index{
		Team:   teamArray,
		Length: len(teams),
	}
	return idx, nil
}

func GetLeadersBoardPageData() (*LeadersBoard, error) {

	var leadersBoard LeadersBoard
	var userScore UserScore

	for _, player := range players {
		var totalPar, totalScore, passedHoleCnt int
		for holeIndex, playerScore := range player.Score {
			if playerScore["total"] != 0 {
				totalPar += fields[holeIndex].Par
				totalScore += playerScore["total"].(int)
				passedHoleCnt += 1
			}
		}
		userScore = UserScore{
			Score: totalScore - totalPar,
			Total: totalScore,
			Name:  player.Name,
			Hole:  passedHoleCnt,
		}
		leadersBoard.Ranking = append(leadersBoard.Ranking, userScore)
	}
	sort.Sort(sortByScore(leadersBoard.Ranking))

	return &leadersBoard, nil
}

func GetScoreEntrySheetPageData(teamName string, holeString string) (*ScoreEntrySheet, error) {
	if len(holeString) == 0 {
		return nil, nil
	}
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	field := fields[holeIndex]
	playersInTheTeam := GetPlayersDataInTheTeam(teamName)

	member := make([]string, 4)
	stroke := make([]int, 4)
	putt := make([]int, 4)
	for playerIndex, player := range playersInTheTeam {
		member[playerIndex] = player.Name
		stroke[playerIndex] = player.Score[holeIndex]["stroke"].(int)
		putt[playerIndex] = player.Score[holeIndex]["putt"].(int)
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

	member := make([]string, len(playersInTheTeam))
	apply := make([]int, len(playersInTheTeam))
	holes := make([]Hole, len(fields))
	totalScore := make([]int, len(playersInTheTeam))
	var totalPar int
	for holeIndex, field := range fields {

		totalPar += field.Par
		score := make([]int, len(playersInTheTeam))
		for playerIndex, player := range playersInTheTeam {
			score[playerIndex] = player.Score[holeIndex]["total"].(int)
			totalScore[playerIndex] += score[playerIndex]
			if holeIndex == 0 {
				member[playerIndex] = player.Name
				apply[playerIndex] = player.Apply
			}
		}
		holes[holeIndex] = Hole{
			Hole:  holeIndex + 1,
			Par:   field.Par,
			Yard:  field.Yard,
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

func GetEntireScorePageData() (*EntireScore, error) {

}

func PostScoreEntrySheetPageData(teamName string, holeString string, updatedTeamScore *PostTeamScore) (*Status, error) {

	if len(holeString) == 0 {
		return &Status{"failed"}, errors.New("hole is not string")
	}
	holeNum, _ := strconv.Atoi(holeString)
	holeIndex := holeNum - 1

	fmt.Println("Team : " + teamName + ", Hole : " + holeString + "にデータを挿入します。")
	db, session := mongoInit()
	playerCol := db.C("player")
	defer session.Close()

	teamPlayers := GetPlayersDataInTheTeam(teamName)

	for playerIndex, player := range teamPlayers {
		stroke, putt := updatedTeamScore.Stroke[playerIndex], updatedTeamScore.Putt[playerIndex]
		player.Score[holeIndex]["stroke"] = stroke
		player.Score[holeIndex]["putt"] = putt
		player.Score[holeIndex]["total"] = stroke + putt

		query := bson.M{"_id": player.Id}
		if err := playerCol.Update(query, player); err != nil {
			return &Status{"failed"}, err
		}
	}

	players = GetAllPlayerCol()
	return &Status{"success"}, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
