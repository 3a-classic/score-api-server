package mongo

import (
	"sort"
	"strconv"
)

func GetIndexPageData() (*Index, error) {

	teamArray := make([]string, len(teams))
	for i, team := range teams {
		teamArray[i] = team.Team
	}

	sort.Strings(teamArray)
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
				total, _ := playerScore["total"].(int)
				totalScore += total
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

	targetTeam := Team{}
	for _, team := range teams {
		if team.Team == teamName {
			targetTeam = team
		}
	}

	member := make([]string, len(playersInTheTeam))
	total := make([]int, len(playersInTheTeam))
	putt := make([]int, len(playersInTheTeam))
	for playerIndex, player := range playersInTheTeam {
		member[playerIndex] = player.Name
		total[playerIndex], _ = player.Score[holeIndex]["total"].(int)
		putt[playerIndex], _ = player.Score[holeIndex]["putt"].(int)
	}

	scoreEntrySheet := ScoreEntrySheet{
		Team:   teamName,
		Hole:   holeNum,
		Member: member,
		Par:    field.Par,
		Yard:   field.Yard,
		Total:  total,
		Putt:   putt,
		Excnt:  targetTeam.Excnt[holeIndex],
	}
	return &scoreEntrySheet, nil
}

func GetScoreViewSheetPageData(teamName string) (*ScoreViewSheet, error) {

	playersInTheTeam := GetPlayersDataInTheTeam(teamName)

	member := make([]string, len(playersInTheTeam))
	apply := make([]int, len(playersInTheTeam))
	holes := make([]Hole, len(fields))
	totalScore := make([]int, len(playersInTheTeam))
	var defined bool

	for _, team := range teams {
		if team.Team == teamName {
			defined = team.Defined
		}
	}
	var totalPar int
	for holeIndex, field := range fields {

		totalPar += field.Par
		score := make([]int, len(playersInTheTeam))
		for playerIndex, player := range playersInTheTeam {
			score[playerIndex], _ = player.Score[holeIndex]["total"].(int)
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
		Team:    teamName,
		Member:  member,
		Apply:   apply,
		Hole:    holes,
		Sum:     sum,
		Defined: defined,
	}

	return &scoreViewSheet, nil
}

func GetEntireScorePageData() (*EntireScore, error) {

	rowSize := 27
	columnSize := len(players) + 2
	rankMap := make(map[string]bool)

	rows := make([][]string, rowSize)
	for i := 0; i < rowSize; i++ {
		if i == 0 {
			rows[i] = make([]string, len(teams)*2)
		} else {
			rows[i] = make([]string, columnSize)
		}
	}
	mainColumnNum := 2

	rows[1][0] = "ホール"
	rows[1][1] = "パー"
	rows[2][1] = "OUT"
	rows[12][1] = "IN"
	rows[20][0] = "Gross"
	rows[20][1] = "-"
	rows[21][0] = "Net"
	rows[21][1] = "-"
	rows[22][0] = "申請"
	rows[22][1] = "-"
	rows[23][0] = "スコア差"
	rows[23][1] = "-"
	rows[24][0] = "順位"
	rows[24][1] = "-"

	var passedPlayerNum int
	for teamIndex, team := range teams {
		rows[0][teamIndex*2] = strconv.Itoa(len(team.Member))
		rows[0][teamIndex*2+1] = "TEAM " + team.Team
		playerInTheTeam := GetPlayersDataInTheTeam(team.Team)
		for playerIndex, player := range playerInTheTeam {
			userDataIndex := playerIndex + passedPlayerNum + mainColumnNum
			rows[1][userDataIndex] = player.Name

			var gross int
			var net int
			for holeIndex, field := range fields {
				holeRowNum := holeIndex + 3
				if holeRowNum > 11 {
					holeRowNum = holeIndex + 4
				}
				if playerIndex == 0 {
					holeNum := field.Hole
					if holeNum > 9 {
						holeNum = holeNum - 9
					}
					if field.Ignore {
						rows[holeRowNum][0] = "-i" + strconv.Itoa(holeNum)
					} else {
						rows[holeRowNum][0] = strconv.Itoa(holeNum)
					}
					rows[holeRowNum][1] = strconv.Itoa(field.Par)
				}
				playerTotal, _ := player.Score[holeIndex]["total"].(int)
				gross += playerTotal
				if field.Ignore {
					net += field.Par
				} else {
					net += playerTotal

				}
				rows[holeRowNum][userDataIndex] = strconv.Itoa(playerTotal)
			}
			rows[20][userDataIndex] = strconv.Itoa(gross)
			rows[21][userDataIndex] = strconv.Itoa(net)
			rows[22][userDataIndex] = strconv.Itoa(player.Apply)
			diff := player.Apply - net
			if diff < 0 {
				diff = diff * -1
			}
			rows[23][userDataIndex] = strconv.Itoa(diff)
			rankMap[strconv.Itoa(diff)] = true
		}
		passedPlayerNum += len(playerInTheTeam)
	}

	var rank []int
	for k := range rankMap {
		intK, _ := strconv.Atoi(k)
		rank = append(rank, intK)
	}

	sort.Ints(rank)
	for i := 0; i < len(rows[24])-2; i++ {
		userDataIndex := i + 2
		for j := 0; j < len(rank); j++ {
			if rows[23][userDataIndex] != strconv.Itoa(rank[j]) {
				continue
			} else {
				rows[24][userDataIndex] = strconv.Itoa(j + 1)
			}
		}
	}

	EntireScore := EntireScore{
		Rows: rows,
	}

	return &EntireScore, nil
}

func GetTimeLinePageData() (*TimeLine, error) {

	var timeLine TimeLine

	reaction1 := Reaction{
		Name:        "matsuno",
		ContentType: 1,
		Content:     "https://s3-ap-northeast-1.amazonaws.com/3a-classic/reaction-icon/angry.png",
	}

	reaction2 := Reaction{
		Name:        "kiyota",
		ContentType: 1,
		Content:     "https://s3-ap-northeast-1.amazonaws.com/3a-classic/reaction-icon/like.png",
	}

	var reactions []Reaction
	reactions = append(reactions, reaction1)
	reactions = append(reactions, reaction2)

	thread1 := Thread{
		ThreadId:  "GSHDLKFJSDLK",
		Msg:       "kiyotaさんがホール13でアルバトロスを出しました！",
		ImgUrl:    "https://s3-ap-northeast-1.amazonaws.com/3a-classic/test/emotion-img/positive-dummy.jpg",
		ColorCode: "#FF0000",
		Reactions: reactions,
		Positive:  true,
	}

	thread2 := Thread{
		ThreadId:  "GSHDGSFJSGDS",
		Msg:       "matsunoさんがホール2で+12を出しました。。",
		ImgUrl:    "https://s3-ap-northeast-1.amazonaws.com/3a-classic/test/emotion-img/negative-dummy.jpg",
		ColorCode: "#FFFF00",
		Reactions: reactions,
		Positive:  false,
	}

	timeLine.Threads = append(timeLine.Threads, thread1)
	timeLine.Threads = append(timeLine.Threads, thread2)

	return &timeLine, nil
}

// sort
// http://grokbase.com/t/gg/golang-nuts/132d2rt3hh/go-nuts-how-to-sort-an-array-of-struct-by-field
func (s sortByScore) Len() int           { return len(s) }
func (s sortByScore) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortByScore) Less(i, j int) bool { return s[i].Score < s[j].Score }
