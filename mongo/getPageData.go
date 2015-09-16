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

	columnSize := len(players) + 2
	holeSize := len(fields)
	holeRows := make([][]string, holeSize)
	for i := 0; i < holeSize; i++ {
		holeRows[i] = make([]string, columnSize)
	}

	//	for i := 0; i < holeSize; i++ {
	//		grossRow[i] := make([]string, holeSize)
	//		for j := 0; j < len(grossRow); j++ {
	//			grossRow[i][j] = make([]string, columnSize)
	//		}
	//	}
	//	rows := make([]string, 25)
	var teamRow []string
	mainColumnNum := 2
	main_row := make([]string, len(players)+2)
	//	holeRows := make([]string, len(players)+2)
	//	var grossRow [holeSize][columnSize]string
	//  grossRow make()[holeSize][columnSize]string
	grossRow := make([]string, len(players)+2)
	netRow := make([]string, len(players)+2)
	applyRow := make([]string, len(players)+2)
	diffRow := make([]string, len(players)+2)
	orderRow := make([]string, len(players)+2)

	main_row[0] = "ホール"
	main_row[1] = "パー"
	grossRow[0] = "Gross"
	grossRow[1] = "-"
	netRow[0] = "Net"
	netRow[1] = "-"
	applyRow[0] = "申請"
	applyRow[1] = "-"
	diffRow[0] = "スコア差"
	diffRow[1] = "-"
	orderRow[0] = "順位"
	orderRow[1] = "-"

	var passedPlayerNum int
	for _, team := range teams {
		teamRow = append(teamRow, strconv.Itoa(len(team.Member)))
		teamRow = append(teamRow, team.Team)
		playerInTheTeam := GetPlayersDataInTheTeam(team.Team)
		for playerIndex, player := range playerInTheTeam {
			userDataIndex := playerIndex + passedPlayerNum + mainColumnNum
			main_row[userDataIndex] = player.Name

			var gross int
			var net int
			for holeIndex, field := range fields {
				if playerIndex == 0 {
					if field.Ignore {
						holeRows[holeIndex][0] = "-i" + strconv.Itoa(field.Hole)
					} else {
						holeRows[holeIndex][0] = strconv.Itoa(field.Hole)
					}
					holeRows[holeIndex][1] = strconv.Itoa(field.Par)
				}
				gross += player.Score[holeIndex]["total"].(int)
				if field.Ignore == false {
					net += player.Score[holeIndex]["total"].(int)
				}
				holeRows[holeIndex][userDataIndex] = strconv.Itoa(player.Score[holeIndex]["total"].(int))
			}
			grossRow[userDataIndex] = strconv.Itoa(gross)
			netRow[userDataIndex] = strconv.Itoa(net)
			applyRow[userDataIndex] = strconv.Itoa(player.Apply)
			diff := player.Apply - net
			if diff < 0 {
				diff = diff * -1
			}
			diffRow[userDataIndex] = strconv.Itoa(diff)
		}
		passedPlayerNum += len(playerInTheTeam)
	}

	fmt.Println(teamRow)
	fmt.Println(main_row)
	fmt.Println(holeRows)
	fmt.Println(grossRow)
	fmt.Println(netRow)
	fmt.Println(applyRow)
	fmt.Println(diffRow)
	fmt.Println(orderRow)
	//	fmt.Println(rows)

	fmt.Println("appendddddddddddddddddd")
	fmt.Println(append(teamRow, main_row))

	EntireScore := EntireScore{
		Rows: nil,
	}

	return &EntireScore, nil
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
